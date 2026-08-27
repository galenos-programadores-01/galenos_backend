// Package sis implementa el adaptador de salida que consulta el servicio web
// SOAP del SIS (app.sis.gob.pe/sisWSAFI). Replica la lógica del proyecto
// FastAPI: primero obtiene la sesión con GetSession(strUsuario, strClave) y
// luego consulta el afiliado con ConsultarAfiliadoFuaE pasando el token como
// strAutorizacion.
package sis

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// Config agrupa las credenciales y endpoint del servicio SIS.
type Config struct {
	Usuario       string
	Clave         string
	URL           string
	DNIAutorizado string
	Timeout       time.Duration
}

const (
	soapNS    = "http://schemas.xmlsoap.org/soap/envelope/"
	sisNS     = "http://sis.gob.pe/"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"
)

// opcionConsulta es el modo de búsqueda del afiliado; 1 = por DNI.
const opcionConsultaDefault = 1

// Tipos de documento soportados por la consulta de afiliado SIS.
const (
	TipoDocumentoDNI               = 1
	TipoDocumentoCarnetExtranjeria = 3
)

// operacionConsulta es el método SOAP que replica el proyecto FastAPI.
const operacionConsulta = "ConsultarAfiliadoFuaE"
const elementoResultado = "ConsultarAfiliadoFuaEResult"

// operacionBuscarPorAfiliacion es la operación SOAP para buscar afiliados
// por DISA/Lote/NroContrato/Correlativo/CodTabla. El SOAPAction usa
// http://www.sis.gob.pe/ (con www) como lo indica el frontend VB.
const operacionBuscarPorAfiliacion = "BuscarAsegurados"
const soapActionBuscarPorAfiliacion = "http://www.sis.gob.pe/BuscarAsegurados"
const elementoResultadoBuscar = "BuscarAseguradosResult"

type client struct {
	cfg  Config
	http *http.Client
}

// New crea el cliente SOAP del SIS.
func New(cfg Config) *client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &client{cfg: cfg, http: &http.Client{Timeout: cfg.Timeout}}
}

// ConsultarAfiliado obtiene el token de sesión y consulta el paciente
// afiliado por su número de documento.
func (c *client) ConsultarAfiliado(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error) {
	if c.cfg.DNIAutorizado == "" {
		return domain.SisAfiliado{}, fmt.Errorf("sis DNI autorizado is not configured")
	}
	if params.Opcion <= 0 {
		params.Opcion = opcionConsultaDefault
	}
	token, err := c.getSession(ctx)
	if err != nil {
		return domain.SisAfiliado{}, err
	}

	body := c.construirConsulta(token, params)
	log.Printf("[SIS] ConsultarAfiliado params: Opcion=%d, DocumentNumber=%q, Disa=%q, TipoFormato=%q, NroContrato=%q, Correlativo=%q, TipoDocumento=%d",
		params.Opcion, params.DocumentNumber, params.Disa, params.TipoFormato, params.NroContrato, params.Correlativo, params.TipoDocumento)
	log.Printf("[SIS] SOAP body:\n%s", string(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.URL, bytes.NewBuffer(body))
	if err != nil {
		return domain.SisAfiliado{}, fmt.Errorf("building sis request: %w", err)
	}
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("SOAPAction", sisNS+operacionConsulta)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/xml, application/xml, */*")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.SisAfiliado{}, fmt.Errorf("calling sis service: %w", err)
	}
	defer resp.Body.Close()

	contenido, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if readErr != nil {
		return domain.SisAfiliado{}, fmt.Errorf("reading sis response: %w", readErr)
	}

	if resp.StatusCode >= 400 {
		return domain.SisAfiliado{}, fmt.Errorf("sis service responded %d: %s", resp.StatusCode, truncate(contenido, 800))
	}

	var raw struct {
		XMLName xml.Name   `xml:"ConsultarAfiliadoFuaEResult"`
		Fields  []rawField `xml:",any"`
	}
	// Se decodifica el elemento por su nombre local para no depender del
	// prefijo del namespace del sobre SOAP.
	if err := decodeElementByLocalName(contenido, elementoResultado, &raw); err != nil {
		// El SIS responde sin el elemento cuando la sesión/autorización no es
		// válida o el DNI no tiene acceso al tipo de consulta.
		return domain.SisAfiliado{}, fmt.Errorf("parsing sis response: %w: %s", err, truncate(contenido, 400))
	}

	return mapAfiliado(raw.Fields), nil
}

// getSession llama a GetSession(strUsuario, strClave) y devuelve el token de
// autorización que exige la consulta de afiliado.
func (c *client) getSession(ctx context.Context) (string, error) {
	body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="%s">
  <soap:Body>
    <GetSession xmlns="%s">
      <strUsuario>%s</strUsuario>
      <strClave>%s</strClave>
    </GetSession>
  </soap:Body>
</soap:Envelope>`,
		soapNS,
		sisNS,
		escapeXML(c.cfg.Usuario),
		escapeXML(c.cfg.Clave),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.URL, bytes.NewBufferString(body))
	if err != nil {
		return "", fmt.Errorf("building sis session request: %w", err)
	}
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("SOAPAction", sisNS+"GetSession")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/xml, application/xml, */*")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling sis session service: %w", err)
	}
	defer resp.Body.Close()

	contenido, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if readErr != nil {
		return "", fmt.Errorf("reading sis session response: %w", readErr)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("sis session service responded %d: %s", resp.StatusCode, truncate(contenido, 800))
	}

	var result struct {
		XMLName xml.Name `xml:"GetSessionResult"`
		Value   string   `xml:",chardata"`
	}
	if err := decodeElementByLocalName(contenido, "GetSessionResult", &result); err != nil {
		return "", fmt.Errorf("parsing sis session response: %w", err)
	}

	token := result.Value
	if token == "" || !esNumerico(token) {
		return "", fmt.Errorf("sis session rejected: %q", truncate([]byte(token), 120))
	}

	return token, nil
}

// construirConsulta arma el sobre SOAP de ConsultarAfiliadoFuaE. El SIS exige
// el DNI autorizado de la llamada en strDni y el documento a consultar se
// envía como strTipoDocumento + strNroDocumento: strTipoDocumento 1 (DNI) o 2
// (Carnet de Extranjería), con el número en strNroDocumento. Los demás
// parámetros (strDisa, strTipoFormato, strNroContrato, strCorrelativo) son
// opcionales y se replican los valores que llegan del request.
func (c *client) construirConsulta(token string, params shared.SISAfiliadoParams) []byte {
	tipo := TipoDocumentoDNI
	if params.TipoDocumento == TipoDocumentoCarnetExtranjeria {
		tipo = TipoDocumentoCarnetExtranjeria
	}
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="%s">
  <soap:Body>
    <ConsultarAfiliadoFuaE xmlns="%s">
      <intOpcion>%d</intOpcion>
      <strAutorizacion>%s</strAutorizacion>
      <strDni>%s</strDni>
      <strTipoDocumento>%d</strTipoDocumento>
      <strNroDocumento>%s</strNroDocumento>
      <strDisa>%s</strDisa>
      <strTipoFormato>%s</strTipoFormato>
      <strNroContrato>%s</strNroContrato>
<strCorrelativo>%s</strCorrelativo>
    </ConsultarAfiliadoFuaE>
  </soap:Body>
</soap:Envelope>`,
		soapNS,
		sisNS,
		params.Opcion,
		escapeXML(token),
		escapeXML(c.cfg.DNIAutorizado),
		tipo,
		escapeXML(params.DocumentNumber),
		escapeXML(params.Disa),
		escapeXML(params.TipoFormato),
		escapeXML(params.NroContrato),
		escapeXML(params.Correlativo),
	))
}

// BuscarPorAfiliacion consulta el paciente afiliado por los parámetros de
// afiliación (DISA, Lote, Contrato, Correlativo, CodTabla) usando la
// operación BuscarAsegurados del SIS.
// NOTA: el VB original NO llama GetSession para esta operación; el token
// de sisWSAFI no aplica para RecepcionTrama.asmx.
func (c *client) BuscarPorAfiliacion(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error) {
	body := c.construirBuscarPorAfiliacion(params)
	log.Printf("[SIS] BuscarPorAfiliacion params: Disa=%q, Lote=%q, Contrato=%q, Correlativo=%q, CodTabla=%q",
		params.Disa, params.Lote, params.NroContrato, params.Correlativo, params.CodTabla)
	log.Printf("[SIS] SOAP body:\n%s", string(body))

	buscarURL := "http://app.sis.gob.pe/edi/RecepcionTrama.asmx"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, buscarURL, bytes.NewBuffer(body))
	if err != nil {
		return domain.SisAfiliado{}, fmt.Errorf("building sis request: %w", err)
	}
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("SOAPAction", soapActionBuscarPorAfiliacion)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/xml, application/xml, */*")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.SisAfiliado{}, fmt.Errorf("calling sis service: %w", err)
	}
	defer resp.Body.Close()

	contenido, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if readErr != nil {
		return domain.SisAfiliado{}, fmt.Errorf("reading sis response: %w", readErr)
	}

	if resp.StatusCode >= 400 {
		return domain.SisAfiliado{}, fmt.Errorf("sis service responded %d: %s", resp.StatusCode, truncate(contenido, 800))
	}

	log.Printf("[SIS] BuscarPorAfiliacion response (%d bytes):\n%s", len(contenido), truncate(contenido, 3000))

	var rawResult struct {
		XMLName  xml.Name `xml:"BuscarAseguradosResult"`
		InnerXML string   `xml:",innerxml"`
	}
	if err := decodeElementByLocalName(contenido, elementoResultadoBuscar, &rawResult); err != nil {
		return domain.SisAfiliado{}, fmt.Errorf("parsing sis response: %w: %s", err, truncate(contenido, 400))
	}

	log.Printf("[SIS] BuscarAseguradosResult innerXML:\n%s", rawResult.InnerXML)

	var raw struct {
		XMLName xml.Name   `xml:"result"`
		Fields  []rawField `xml:",any"`
	}
	if err := xml.Unmarshal([]byte("<result>"+rawResult.InnerXML+"</result>"), &raw); err != nil {
		log.Printf("[SIS] BuscarPorAfiliacion: raw innerXML is not XML, treating as delimited text")
		return domain.SisAfiliado{
			Resultado: rawResult.InnerXML,
		}, nil
	}

	log.Printf("[SIS] BuscarPorAfiliacion parsed fields: %d", len(raw.Fields))
	for _, f := range raw.Fields {
		log.Printf("[SIS]   %s = %q", f.XMLName.Local, f.Value)
	}

	return mapAfiliado(raw.Fields), nil
}

// construirBuscarPorAfiliacion arma el sobre SOAP de BuscarAsegurados.
// Replica la lógica del frontend VB: usa Disa, TipoFormato (=Lote),
// Contrato, Correlativo y CodTabla. No usa namespace en el body porque
// el template VB original (SIS_Buscar_AfiliadoxNroAfiliado) no lo incluye.
func (c *client) construirBuscarPorAfiliacion(params shared.SISAfiliadoParams) []byte {
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="%s">
  <soap:Body>
    <BuscarAsegurados>
      <Disa>%s</Disa>
      <TipoFormato>%s</TipoFormato>
      <Contrato>%s</Contrato>
      <Correlativo>%s</Correlativo>
      <CodTabla>%s</CodTabla>
    </BuscarAsegurados>
  </soap:Body>
</soap:Envelope>`,
		soapNS,
		escapeXML(params.Disa),
		escapeXML(params.Lote),
		escapeXML(params.NroContrato),
		escapeXML(params.Correlativo),
		escapeXML(params.CodTabla),
	))
}

// rawField captura un elemento hijo con su nombre y su contenido de texto,
// de modo que la respuesta del SIS se mapee sin depender del orden ni del
// prefijo de namespace.
type rawField struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// mapAfiliado convierte la lista de campos XML cruda en el modelo de dominio.
func mapAfiliado(fields []rawField) domain.SisAfiliado {
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		m[f.XMLName.Local] = f.Value
	}
	val := func(k string) string { return m[k] }
	return domain.SisAfiliado{
		IdError:            val("IdError"),
		Resultado:          val("Resultado"),
		TipoDocumento:      val("TipoDocumento"),
		NroDocumento:       val("NroDocumento"),
		ApePaterno:         val("ApePaterno"),
		ApeMaterno:         val("ApeMaterno"),
		Nombres:            val("Nombres"),
		FecAfiliacion:      val("FecAfiliacion"),
		EESS:               val("EESS"),
		DescEESS:           val("DescEESS"),
		EESSUbigeo:         val("EESSUbigeo"),
		DescEESSUbigeo:     val("DescEESSUbigeo"),
		Regimen:            val("Regimen"),
		TipoSeguro:         val("TipoSeguro"),
		DescTipoSeguro:     val("DescTipoSeguro"),
		Contrato:           val("Contrato"),
		FecCaducidad:       val("FecCaducidad"),
		Estado:             val("Estado"),
		Tabla:              val("Tabla"),
		IdNumReg:           val("IdNumReg"),
		Genero:             val("Genero"),
		FecNacimiento:      val("FecNacimiento"),
		IdUbigeo:           val("IdUbigeo"),
		Direccion:          val("Direccion"),
		Disa:               val("Disa"),
		TipoFormato:        val("TipoFormato"),
		NroContrato:        val("NroContrato"),
		Correlativo:        val("Correlativo"),
		IdPlan:             val("IdPlan"),
		IdGrupoPoblacional: val("IdGrupoPoblacional"),
		MsgConfidencial:    val("MsgConfidencial"),
	}
}

// decodeElementByLocalName recorre el XML y decodifica el primer elemento
// cuyo nombre local coincida con el solicitado.
func decodeElementByLocalName(data []byte, name string, v interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			return fmt.Errorf("element %s not found in response", name)
		}
		if err != nil {
			return err
		}
		if start, ok := tok.(xml.StartElement); ok && start.Name.Local == name {
			return decoder.DecodeElement(v, &start)
		}
	}
}

func escapeXML(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

// esNumerico indica si el token de sesión del SIS es un valor válido. Los
// errores de credenciales llegan como texto ("CLAVE INCORRECTA"), mientras
// que un token válido es un número.
func esNumerico(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func truncate(b []byte, max int) string {
	s := string(b)
	if len(s) > max {
		s = s[:max]
	}
	return s
}
