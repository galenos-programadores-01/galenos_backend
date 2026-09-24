// Package minsa implementa el adaptador de salida del servicio REST de
// referencias y contrarreferencias del MINSA (consultaReferenciaDetalle y
// listadoUps).
package minsa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
)

const (
	defaultURL               = "https://servicios.minsa.gob.pe/mcs-servicios-refcon/servicio/v1.0.0/consultaReferenciaDetalle"
	defaultUpsURL            = "https://servicios.minsa.gob.pe/mcs-servicios-refcon/servicio/v1.0.0/listadoUps"
	defaultEspecialidadesURL = "https://servicios.minsa.gob.pe/mcs-referencia-interoperabilidad/refcon-interoperabilidad/v1.0/listadoEspecialidades"
	defaultSaveReferenciaURL = "https://servicios.minsa.gob.pe/mcs-servicios-refcon/servicio/v1.0.0/saveReferencia"
	maxResponse              = 5 << 20 // 5 MB (incluye anamnesis y exámenes largos)
	defaultTimeout           = 30 * time.Second
	defaultDestino           = "7634"
	defaultLimite            = "11"
	defaultUserAgent         = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
)

// Config agrupa la URL del servicio, las credenciales institucionales que van
// en el header y los valores por defecto de la petición.
type Config struct {
	URL                    string
	UpsURL                 string
	EspecialidadesURL      string
	SaveReferenciaURL      string
	Username               string
	Password               string
	IPClient               string
	EstablecimientoDestino string
	Limite                 string
	Timeout                time.Duration
}

// client implementa output.MinsaRefConClient.
type client struct {
	cfg  Config
	http *http.Client
}

// New crea el cliente del servicio de referencias de MINSA. Los valores por
// defecto garantizan que siempre se envíen establecimientoDestino y limite aun
// si la configuración llegara vacía.
func New(cfg Config) *client {
	if cfg.URL == "" {
		cfg.URL = defaultURL
	}
	if cfg.UpsURL == "" {
		cfg.UpsURL = defaultUpsURL
	}
	if cfg.EspecialidadesURL == "" {
		cfg.EspecialidadesURL = defaultEspecialidadesURL
	}
	if cfg.SaveReferenciaURL == "" {
		cfg.SaveReferenciaURL = defaultSaveReferenciaURL
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}
	if cfg.EstablecimientoDestino == "" {
		cfg.EstablecimientoDestino = defaultDestino
	}
	if cfg.Limite == "" {
		cfg.Limite = defaultLimite
	}
	log.Printf("[DashRefCon] cliente configurado ipclient=%s establecimientoDestino=%s limite=%s", cfg.IPClient, cfg.EstablecimientoDestino, cfg.Limite)

	// Transporte directo (sin proxy) forzando HTTP/1.1: el gateway del MINSA
	// descarta los headers username/password/ipclient cuando se negocia HTTP/2
	// o se atraviesa un proxy, devolviendo el error de autorización 6000.
	transport := &http.Transport{
		Proxy:             nil,
		ForceAttemptHTTP2: false,
		MaxIdleConns:      10,
		IdleConnTimeout:   90 * time.Second,
	}

	return &client{cfg: cfg, http: &http.Client{Transport: transport, Timeout: cfg.Timeout}}
}

// ConsultarReferenciaDetalle envía la petición de detalle de referencia al
// servicio de MINSA con las credenciales institucionales en el header.
func (c *client) ConsultarReferenciaDetalle(ctx context.Context, req domain.ConsultaMinsaRequest) (*domain.ConsultaMinsaResponse, error) {
	if c.cfg.Username == "" || c.cfg.Password == "" || c.cfg.IPClient == "" {
		return nil, fmt.Errorf("minsa refcon credentials are required")
	}

	if req.EstablecimientoDestino == "" {
		req.EstablecimientoDestino = c.cfg.EstablecimientoDestino
	}
	if req.Limite == "" {
		req.Limite = c.cfg.Limite
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling minsa refcon request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building minsa refcon request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("username", c.cfg.Username)
	httpReq.Header.Set("password", c.cfg.Password)
	httpReq.Header.Set("ipclient", c.cfg.IPClient)
	httpReq.Header.Set("User-Agent", defaultUserAgent)

	log.Printf("[DashRefCon] consultaReferenciaDetalle body: %s", body)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling minsa refcon consultaReferenciaDetalle: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse))
	if err != nil {
		return nil, fmt.Errorf("reading minsa refcon response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("minsa refcon responded %d: %s", resp.StatusCode, truncate(respBody, 400))
	}

	var out domain.ConsultaMinsaResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decoding minsa refcon response: %w", err)
	}

	if out.Codigo != "" && out.Codigo != "0000" {
		log.Printf("[DashRefCon] respuesta MINSA status=%d codigo=%s mensaje=%q body=%s", resp.StatusCode, out.Codigo, out.Mensaje, truncate(respBody, 500))
	}

	return &out, nil
}

// ListadoUpss consulta el listado de Unidades Productoras de Servicios (UPS)
// del establecimiento indicado por su código RENIPRESS.
func (c *client) ListadoUpss(ctx context.Context, codigoRenipress string) (*domain.ListadoUpsResponse, error) {
	if c.cfg.Username == "" || c.cfg.Password == "" || c.cfg.IPClient == "" {
		return nil, fmt.Errorf("minsa refcon credentials are required")
	}
	if codigoRenipress == "" {
		return nil, fmt.Errorf("minsa refcon codigoRenipress is required")
	}

	url := strings.TrimSuffix(c.cfg.UpsURL, "/") + "/" + codigoRenipress

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building minsa refcon listadoUps request: %w", err)
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("username", c.cfg.Username)
	httpReq.Header.Set("password", c.cfg.Password)
	httpReq.Header.Set("ipclient", c.cfg.IPClient)
	httpReq.Header.Set("User-Agent", defaultUserAgent)

	log.Printf("[DashRefCon] listadoUps url=%s", url)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling minsa refcon listadoUps: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse))
	if err != nil {
		return nil, fmt.Errorf("reading minsa refcon response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("minsa refcon responded %d: %s", resp.StatusCode, truncate(respBody, 400))
	}

	var out domain.ListadoUpsResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decoding minsa refcon listadoUps response: %w", err)
	}

	if out.Codigo != "" && out.Codigo != "0000" {
		log.Printf("[DashRefCon] respuesta MINSA listadoUps status=%d codigo=%s mensaje=%q body=%s", resp.StatusCode, out.Codigo, out.Mensaje, truncate(respBody, 500))
	}

	return &out, nil
}

// ListadoEspecialidades consulta el listado de especialidades vigentes del
// servicio REST de interoperabilidad del MINSA (sin parámetros).
func (c *client) ListadoEspecialidades(ctx context.Context) (*domain.ListadoEspecialidadesResponse, error) {
	if c.cfg.Username == "" || c.cfg.Password == "" || c.cfg.IPClient == "" {
		return nil, fmt.Errorf("minsa refcon credentials are required")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.EspecialidadesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building minsa refcon listadoEspecialidades request: %w", err)
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("username", c.cfg.Username)
	httpReq.Header.Set("password", c.cfg.Password)
	httpReq.Header.Set("ipclient", c.cfg.IPClient)
	httpReq.Header.Set("User-Agent", defaultUserAgent)

	log.Printf("[DashRefCon] listadoEspecialidades url=%s", c.cfg.EspecialidadesURL)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling minsa refcon listadoEspecialidades: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse))
	if err != nil {
		return nil, fmt.Errorf("reading minsa refcon response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("minsa refcon responded %d: %s", resp.StatusCode, truncate(respBody, 400))
	}

	var out domain.ListadoEspecialidadesResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decoding minsa refcon listadoEspecialidades response: %w", err)
	}

	if out.Codigo != "" && out.Codigo != "0000" {
		log.Printf("[DashRefCon] respuesta MINSA listadoEspecialidades status=%d codigo=%s body=%s", resp.StatusCode, out.Codigo, truncate(respBody, 500))
	}

	return &out, nil
}

func truncate(b []byte, max int) string {
	s := string(b)
	if len(s) > max {
		s = s[:max]
	}
	return s
}

// SaveReferencia registra una referencia en el servicio saveReferencia del
// MINSA enviando el cuerpo JSON con las credenciales institucionales en el
// header.
func (c *client) SaveReferencia(ctx context.Context, req domain.SaveReferenciaRequest) (*domain.SaveReferenciaResponse, error) {
	if c.cfg.Username == "" || c.cfg.Password == "" || c.cfg.IPClient == "" {
		return nil, fmt.Errorf("minsa refcon credentials are required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling minsa refcon saveReferencia request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.SaveReferenciaURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building minsa refcon saveReferencia request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("username", c.cfg.Username)
	httpReq.Header.Set("password", c.cfg.Password)
	httpReq.Header.Set("ipclient", c.cfg.IPClient)
	httpReq.Header.Set("User-Agent", defaultUserAgent)

	log.Printf("[DashRefCon] saveReferencia url=%s", c.cfg.SaveReferenciaURL)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling minsa refcon saveReferencia: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse))
	if err != nil {
		return nil, fmt.Errorf("reading minsa refcon response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("minsa refcon responded %d: %s", resp.StatusCode, truncate(respBody, 400))
	}

	var out domain.SaveReferenciaResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decoding minsa refcon saveReferencia response: %w", err)
	}

	if out.Codigo != "" && out.Codigo != "0000" {
		log.Printf("[DashRefCon] respuesta MINSA saveReferencia status=%d codigo=%s mensaje=%q body=%s", resp.StatusCode, out.Codigo, derefStr(out.Mensaje), truncate(respBody, 500))
	}

	return &out, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
