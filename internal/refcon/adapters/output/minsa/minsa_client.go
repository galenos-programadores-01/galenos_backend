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
	defaultURL     = "https://servicios.minsa.gob.pe/mcs-servicios-refcon/servicio/v1.0.0/consultaReferenciaDetalle"
	defaultUpsURL  = "https://servicios.minsa.gob.pe/mcs-servicios-refcon/servicio/v1.0.0/listadoUps"
	maxResponse    = 5 << 20 // 5 MB (incluye anamnesis y exámenes largos)
	defaultTimeout = 30 * time.Second
	defaultDestino = "7634"
	defaultLimite  = "11"
)

// Config agrupa la URL del servicio, las credenciales institucionales que van
// en el header y los valores por defecto de la petición.
type Config struct {
	URL                    string
	UpsURL                 string
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

func truncate(b []byte, max int) string {
	s := string(b)
	if len(s) > max {
		s = s[:max]
	}
	return s
}
