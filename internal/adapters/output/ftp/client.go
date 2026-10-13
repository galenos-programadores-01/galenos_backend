package ftp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"path"
	"strings"
	"time"

	goftp "github.com/jlaffaye/ftp"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

// client implementa output.AlmacenamientoRemoto via FTP/FTPS (puerto 21 tipico).
type client struct {
	cfg Config
}

// New construye el adaptador FTP. Si Host o User estan vacios, devuelve noop
// (Enabled=false) para no romper modulos que aun no usan almacenamiento remoto.
func New(cfg Config) (output.AlmacenamientoRemoto, error) {
	cfg = normalizeConfig(cfg)
	if strings.TrimSpace(cfg.Host) == "" || strings.TrimSpace(cfg.User) == "" {
		return noopClient{}, nil
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("ftp: SFTP_PASSWORD is required when SFTP_HOST is set")
	}
	return &client{cfg: cfg}, nil
}

func normalizeConfig(cfg Config) Config {
	if cfg.Port <= 0 {
		cfg.Port = 21
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	cfg.TLSMode = normalizeTLSMode(cfg.TLSMode)
	return cfg
}

func normalizeTLSMode(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "plain", "none", "off":
		return "plain"
	case "implicit":
		return "implicit"
	default:
		// explicit: AUTH TLS antes de USER/PASS (caso IIS hospitalario).
		return "explicit"
	}
}

func (c *client) Enabled() bool { return true }

func (c *client) Subir(ctx context.Context, rutaRemota string, contenido []byte) error {
	ruta, err := sanitizeRemotePath(rutaRemota)
	if err != nil {
		return err
	}
	return c.withSession(ctx, func(conn *goftp.ServerConn) error {
		dir := path.Dir(ruta)
		if dir != "" && dir != "." && dir != "/" {
			if err := mkdirAllFTP(conn, dir); err != nil {
				return err
			}
		}
		if err := conn.Stor(ruta, bytes.NewReader(contenido)); err != nil {
			return fmt.Errorf("ftp stor %s: %w", ruta, err)
		}
		return nil
	})
}

func (c *client) Descargar(ctx context.Context, rutaRemota string) ([]byte, error) {
	ruta, err := sanitizeRemotePath(rutaRemota)
	if err != nil {
		return nil, err
	}
	var out []byte
	err = c.withSession(ctx, func(conn *goftp.ServerConn) error {
		r, err := conn.Retr(ruta)
		if err != nil {
			if isFTPNotFound(err) {
				return domain.ErrAlmacenamientoArchivoNoEncontrado
			}
			return fmt.Errorf("ftp retr %s: %w", ruta, err)
		}
		defer r.Close()
		out, err = io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("ftp read %s: %w", ruta, err)
		}
		return nil
	})
	return out, err
}

func (c *client) Existe(ctx context.Context, rutaRemota string) (bool, error) {
	ruta, err := sanitizeRemotePath(rutaRemota)
	if err != nil {
		return false, err
	}
	var exists bool
	err = c.withSession(ctx, func(conn *goftp.ServerConn) error {
		_, err := conn.FileSize(ruta)
		if err == nil {
			exists = true
			return nil
		}
		if isFTPNotFound(err) {
			exists = false
			return nil
		}
		return fmt.Errorf("ftp size %s: %w", ruta, err)
	})
	return exists, err
}

func (c *client) Eliminar(ctx context.Context, rutaRemota string) error {
	ruta, err := sanitizeRemotePath(rutaRemota)
	if err != nil {
		return err
	}
	return c.withSession(ctx, func(conn *goftp.ServerConn) error {
		if err := conn.Delete(ruta); err != nil {
			if isFTPNotFound(err) {
				return domain.ErrAlmacenamientoArchivoNoEncontrado
			}
			return fmt.Errorf("ftp delete %s: %w", ruta, err)
		}
		return nil
	})
}

func (c *client) withSession(ctx context.Context, fn func(*goftp.ServerConn) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	addr := net.JoinHostPort(c.cfg.Host, fmt.Sprintf("%d", c.cfg.Port))
	opts := []goftp.DialOption{
		goftp.DialWithTimeout(c.cfg.Timeout),
		goftp.DialWithContext(ctx),
	}

	tlsCfg := &tls.Config{
		ServerName:         c.cfg.Host,
		InsecureSkipVerify: c.cfg.InsecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
	}
	switch c.cfg.TLSMode {
	case "explicit":
		opts = append(opts, goftp.DialWithExplicitTLS(tlsCfg))
	case "implicit":
		opts = append(opts, goftp.DialWithTLS(tlsCfg))
	}

	conn, err := goftp.Dial(addr, opts...)
	if err != nil {
		return fmt.Errorf("ftp dial %s tls=%s: %w", addr, c.cfg.TLSMode, err)
	}
	defer conn.Quit()

	if err := conn.Login(c.cfg.User, c.cfg.Password); err != nil {
		return fmt.Errorf("ftp login %s: %w", c.cfg.User, err)
	}

	done := make(chan error, 1)
	go func() { done <- fn(conn) }()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// mkdirAllFTP crea directorios padres uno a uno (FTP no tiene MkdirAll nativo).
func mkdirAllFTP(conn *goftp.ServerConn, dir string) error {
	dir = path.Clean(dir)
	if dir == "/" || dir == "." {
		return nil
	}
	parts := strings.Split(strings.Trim(dir, "/"), "/")
	cur := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		cur += "/" + p
		if err := conn.MakeDir(cur); err != nil {
			if err2 := conn.ChangeDir(cur); err2 != nil {
				return fmt.Errorf("ftp mkdir %s: %w", cur, err)
			}
			_ = conn.ChangeDir("/")
		}
	}
	return nil
}

func sanitizeRemotePath(ruta string) (string, error) {
	ruta = strings.TrimSpace(ruta)
	ruta = strings.ReplaceAll(ruta, "\\", "/")
	if ruta == "" {
		return "", domain.ErrAlmacenamientoRutaInvalida
	}
	if strings.Contains(ruta, "..") {
		return "", domain.ErrAlmacenamientoRutaInvalida
	}
	if !strings.HasPrefix(ruta, "/") {
		ruta = "/" + ruta
	}
	return path.Clean(ruta), nil
}

func isFTPNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "550") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "no such file") ||
		strings.Contains(msg, "failed to change directory")
}