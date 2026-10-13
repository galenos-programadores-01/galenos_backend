// Package refconreports implementa el adaptador de salida del portal REFCON
// (refcon.minsa.gob.pe) que, haciendo login por cookies con la misma sesión
// de navegador, genera y descarga la hoja de referencia institucional en PDF.
package refconreports

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
)

const (
	defaultBaseURL = "https://refcon.minsa.gob.pe/refconv02"
	defaultTimeout = 30 * time.Second
	maxResponse    = 5 << 20 // 5 MB (hoja de referencia + anexos)
)

// Config agrupa la configuración del portal REFCON y las credenciales web
// institucionales que se usan para el login.
type Config struct {
	BaseURL     string
	UserWeb     string
	PasswordWeb string
	Timeout     time.Duration
}

// client implementa output.RefConReportsClient.
type client struct {
	cfg      Config
	base     *url.URL
	jar      http.CookieJar
	follow   *http.Client
	noFollow *http.Client
}

// New crea el cliente del portal REFCON. Se usa un ÚNICO CookieJar para
// toda la sesión (login → reporte → descarga), como hace el controlador
// Laravel de referencia, y dos clientes que comparten el jar: uno que sigue
// los redirects (GET inicial y descarga) y otro que los corta en el login.
func New(cfg Config) *client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}
	base, _ := url.Parse(cfg.BaseURL)

	jar, _ := cookiejar.New(nil)

	// Transporte directo (sin proxy) forzando HTTP/1.1: el gateway del MINSA
	// descarta los headers de sesión cuando se negocia HTTP/2 o se atraviesa
	// un proxy, devolviendo 403/Cloudflare.
	transport := &http.Transport{
		Proxy:             nil,
		ForceAttemptHTTP2: false,
		MaxIdleConns:      10,
		IdleConnTimeout:   90 * time.Second,
	}

	follow := &http.Client{Transport: transport, Jar: jar, Timeout: cfg.Timeout}
	noFollow := &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   cfg.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	log.Printf("[RefConReports] cliente configurado baseURL=%s", cfg.BaseURL)

	return &client{cfg: cfg, base: base, jar: jar, follow: follow, noFollow: noFollow}
}

// GenerarHojaReferencia hace login en el portal REFCON, genera el reporte de
// la hoja de referencia y descarga el PDF, devolviéndolo codificado en base64.
func (c *client) GenerarHojaReferencia(ctx context.Context, req domain.GenerarHojaReferenciaRequest) (*domain.GenerarHojaReferenciaResult, error) {
	if c.cfg.UserWeb == "" || c.cfg.PasswordWeb == "" {
		return nil, fmt.Errorf("refcon reports credentials are required")
	}

	if err := c.login(ctx); err != nil {
		return nil, err
	}

	reportName := reporteHojaReferencia(req.EstadoReferencia)

	urlFile, err := c.preview(ctx, req, reportName)
	if err != nil {
		return nil, err
	}

	pdf, err := c.descargarPDF(ctx, urlFile)
	if err != nil {
		return nil, err
	}

	log.Printf("[RefConReports] hoja de referencia generada id=%d urlFile=%s (%d bytes)", req.IDReferencia, urlFile, len(pdf))

	return &domain.GenerarHojaReferenciaResult{
		PDFBase64:        base64.StdEncoding.EncodeToString(pdf),
		URLFile:          urlFile,
		EstadoReferencia: req.EstadoReferencia,
		NombreReporte:    reportName,
	}, nil
}

// login replica el flujo del controlador Laravel: primero un GET para obtener
// las cookies de Cloudflare/REFCON y luego el POST de login con el MISMO
// CookieJar, verificando que REFCON haya emitido la cookie de sesión "byt".
func (c *client) login(ctx context.Context) error {
	desktopURL := c.cfg.BaseURL + "/desktop"

	// 1. GET inicial (sigue redirects).
	reqGet, err := http.NewRequestWithContext(ctx, http.MethodGet, desktopURL, nil)
	if err != nil {
		return fmt.Errorf("building refcon login GET: %w", err)
	}
	c.setBrowserHeaders(reqGet)

	respGet, err := c.follow.Do(reqGet)
	if err != nil {
		return fmt.Errorf("calling refcon login GET: %w", err)
	}
	io.Copy(io.Discard, io.LimitReader(respGet.Body, maxResponse))
	respGet.Body.Close()

	// 2. POST de login (sin seguir redirects).
	form := url.Values{
		"C":       {"LOGIN"},
		"S":       {"INIT"},
		"_dcp":    {""},
		"mlkuser": {c.cfg.UserWeb},
		"mlkpass": {c.cfg.PasswordWeb},
	}

	reqLogin, err := http.NewRequestWithContext(ctx, http.MethodPost, desktopURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("building refcon login POST: %w", err)
	}
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.setBrowserHeaders(reqLogin)

	respLogin, err := c.noFollow.Do(reqLogin)
	if err != nil {
		return fmt.Errorf("calling refcon login POST: %w", err)
	}
	io.Copy(io.Discard, io.LimitReader(respLogin.Body, maxResponse))
	respLogin.Body.Close()

	if respLogin.StatusCode >= 400 {
		return fmt.Errorf("refcon rechazó el login: HTTP %d", respLogin.StatusCode)
	}

	// 3. Verificar la cookie de sesión "byt".
	for _, cookie := range c.jar.Cookies(c.base) {
		if cookie.Name == "byt" {
			return nil
		}
	}

	return fmt.Errorf("login falló: no se recibió cookie de sesión byt")
}

// preview solicita la generación del reporte y extrae el urlFile del header
// X-JSON de la respuesta (mismo contrato que el preview() de Laravel).
func (c *client) preview(ctx context.Context, req domain.GenerarHojaReferenciaRequest, reportName string) (string, error) {
	paramStore, _ := json.Marshal(map[string]int{
		"idreferencia":      req.IDReferencia,
		"idestablecimiento": req.IDEstablecimiento,
	})

	form := url.Values{
		"outputtype":     {"PDF"},
		"reportname":     {reportName},
		"paramstoreport": {string(paramStore)},
	}

	reqReport, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.BaseURL+"/reports", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("building refcon report request: %w", err)
	}
	reqReport.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqReport.Header.Set("Accept", "application/json, text/plain, */*")
	c.setBrowserHeaders(reqReport)

	resp, err := c.follow.Do(reqReport)
	if err != nil {
		return "", fmt.Errorf("calling refcon report generation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
		return "", fmt.Errorf("refcon/cloudflare bloqueó la generación del reporte: HTTP 403 %s", truncate(body, 400))
	}
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
		return "", fmt.Errorf("refcon devolvió un error HTTP %d: %s", resp.StatusCode, truncate(body, 400))
	}

	xJSON := resp.Header.Get("X-JSON")
	if xJSON == "" {
		return "", fmt.Errorf("refcon no devolvió el header X-JSON")
	}

	var header struct {
		URLFile string `json:"urlFile"`
	}
	if err := json.Unmarshal([]byte(xJSON), &header); err != nil {
		return "", fmt.Errorf("el header X-JSON no contiene un JSON válido: %w", err)
	}
	if header.URLFile == "" {
		return "", fmt.Errorf("no fue posible obtener urlFile del reporte")
	}

	// Validación original de Laravel: el penúltimo segmento (separado por "-")
	// no debe ser "0".
	parts := strings.Split(header.URLFile, "-")
	if len(parts) <= 2 || parts[len(parts)-2] == "0" {
		return "", fmt.Errorf("no fue posible obtener la url del reporte (urlFile=%q)", header.URLFile)
	}

	return header.URLFile, nil
}

// descargarPDF descarga el PDF generado reutilizando las cookies de la
// misma sesión y verifica el magic number %PDF.
func (c *client) descargarPDF(ctx context.Context, urlFile string) ([]byte, error) {
	downloadURL := c.cfg.BaseURL + "/reports?C=DL&f=" + url.QueryEscape(urlFile)

	reqGet, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building refcon pdf download request: %w", err)
	}
	reqGet.Header.Set("Accept", "application/pdf,text/html;q=0.9,*/*;q=0.8")
	reqGet.Header.Set("Referer", c.cfg.BaseURL+"/reports")
	c.setBrowserHeaders(reqGet)

	resp, err := c.follow.Do(reqGet)
	if err != nil {
		return nil, fmt.Errorf("calling refcon pdf download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1000))
		return nil, fmt.Errorf("cloudflare bloqueó la descarga del PDF: HTTP 403 %s", truncate(body, 400))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refcon no pudo descargar el PDF: HTTP %d", resp.StatusCode)
	}

	pdf, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse))
	if err != nil {
		return nil, fmt.Errorf("reading refcon pdf: %w", err)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF")) {
		return nil, fmt.Errorf("refcon respondió correctamente pero el contenido no parece un PDF (content-type=%s)", resp.Header.Get("Content-Type"))
	}

	return pdf, nil
}

// reporteHojaReferencia devuelve el nombre interno del reporte según el
// estado de la referencia (rechazado/pendiente/genérico), igual que Laravel.
func reporteHojaReferencia(estado string) string {
	switch strings.ToLower(strings.TrimSpace(estado)) {
	case "rechazado":
		return "../his/reports/referencia-rechazado"
	case "pendiente":
		return "../his/reports/referencia-pendiente"
	default:
		return "../his/reports/referencia"
	}
}

// setBrowserHeaders aplica los headers de navegador Chrome que REFCON espera.
func (c *client) setBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-PE,es;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Origin", c.cfg.BaseURL)
	req.Header.Set("Referer", c.cfg.BaseURL+"/desktop")
}

func truncate(b []byte, max int) string {
	s := string(b)
	if len(s) > max {
		s = s[:max]
	}
	return s
}
