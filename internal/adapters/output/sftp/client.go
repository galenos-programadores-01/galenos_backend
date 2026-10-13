package sftp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

// client implementa output.AlmacenamientoRemoto vía SFTP/SSH.
type client struct {
	cfg Config
}

// New construye el adaptador SFTP. Si Host o User están vacíos, devuelve un
// cliente noop (Enabled=false) para no romper módulos que aún no usan SFTP.
func New(cfg Config) (output.AlmacenamientoRemoto, error) {
	cfg = normalizeConfig(cfg)
	if strings.TrimSpace(cfg.Host) == "" || strings.TrimSpace(cfg.User) == "" {
		return noopClient{}, nil
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("sftp: SFTP_PASSWORD is required when SFTP_HOST is set")
	}
	return &client{cfg: cfg}, nil
}

func normalizeConfig(cfg Config) Config {
	if cfg.Port <= 0 {
		cfg.Port = 22
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return cfg
}

func (c *client) Enabled() bool { return true }

func (c *client) Subir(ctx context.Context, rutaRemota string, contenido []byte) error {
	ruta, err := sanitizeRemotePath(rutaRemota)
	if err != nil {
		return err
	}
	return c.withSession(ctx, func(sc *sftp.Client) error {
		dir := path.Dir(ruta)
		if dir != "" && dir != "." && dir != "/" {
			if err := sc.MkdirAll(dir); err != nil {
				return fmt.Errorf("sftp mkdir %s: %w", dir, err)
			}
		}
		f, err := sc.OpenFile(ruta, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		if err != nil {
			return fmt.Errorf("sftp open %s: %w", ruta, err)
		}
		defer f.Close()
		if _, err := io.Copy(f, bytes.NewReader(contenido)); err != nil {
			return fmt.Errorf("sftp write %s: %w", ruta, err)
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
	err = c.withSession(ctx, func(sc *sftp.Client) error {
		f, err := sc.Open(ruta)
		if err != nil {
			if os.IsNotExist(err) {
				return domain.ErrAlmacenamientoArchivoNoEncontrado
			}
			return fmt.Errorf("sftp open %s: %w", ruta, err)
		}
		defer f.Close()
		out, err = io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("sftp read %s: %w", ruta, err)
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
	err = c.withSession(ctx, func(sc *sftp.Client) error {
		_, err := sc.Stat(ruta)
		if err == nil {
			exists = true
			return nil
		}
		if os.IsNotExist(err) {
			exists = false
			return nil
		}
		return fmt.Errorf("sftp stat %s: %w", ruta, err)
	})
	return exists, err
}

func (c *client) Eliminar(ctx context.Context, rutaRemota string) error {
	ruta, err := sanitizeRemotePath(rutaRemota)
	if err != nil {
		return err
	}
	return c.withSession(ctx, func(sc *sftp.Client) error {
		if err := sc.Remove(ruta); err != nil {
			if os.IsNotExist(err) {
				return domain.ErrAlmacenamientoArchivoNoEncontrado
			}
			return fmt.Errorf("sftp remove %s: %w", ruta, err)
		}
		return nil
	})
}

func (c *client) withSession(ctx context.Context, fn func(*sftp.Client) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	hostKeyCallback, err := c.hostKeyCallback()
	if err != nil {
		return err
	}

	sshCfg := &ssh.ClientConfig{
		User: c.cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.cfg.Password),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         c.cfg.Timeout,
	}

	addr := net.JoinHostPort(c.cfg.Host, fmt.Sprintf("%d", c.cfg.Port))
	dialer := net.Dialer{Timeout: c.cfg.Timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("sftp dial %s: %w", addr, err)
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, sshCfg)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("sftp ssh handshake %s: %w", addr, err)
	}
	sshClient := ssh.NewClient(sshConn, chans, reqs)
	defer sshClient.Close()

	sc, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("sftp client: %w", err)
	}
	defer sc.Close()

	done := make(chan error, 1)
	go func() { done <- fn(sc) }()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (c *client) hostKeyCallback() (ssh.HostKeyCallback, error) {
	if c.cfg.KnownHostsPath != "" {
		cb, err := knownhosts.New(c.cfg.KnownHostsPath)
		if err != nil {
			return nil, fmt.Errorf("sftp known_hosts %s: %w", c.cfg.KnownHostsPath, err)
		}
		return cb, nil
	}
	if c.cfg.InsecureIgnoreHostKey {
		return ssh.InsecureIgnoreHostKey(), nil
	}
	return nil, fmt.Errorf("sftp: configure SFTP_KNOWN_HOSTS or SFTP_INSECURE_IGNORE_HOST_KEY=true")
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
