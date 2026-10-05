// Package reniec implementa el adaptador de salida que consulta el servicio
// web SOAP de RENIEC (wsvmin.minsa.gob.pe), replicando la lógica del
// proyecto FastAPI.
package reniec

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// Config agrupa las credenciales y endpoint del servicio RENIEC.
type Config struct {
	App     string
	Usuario string
	Clave   string
	URL     string
	Timeout time.Duration
}

const (
	baseNS    = "http://schemas.xmlsoap.org/soap/envelope/"
	tempURINS = "http://tempuri.org/"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"
)

var operaciones = map[string]string{
	"basico":   "obtenerDatosBasicos",
	"completo": "obtenerDatosCompletos",
}

type client struct {
	cfg  Config
	http *http.Client
}

// New crea el cliente SOAP de RENIEC.
func New(cfg Config) *client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &client{cfg: cfg, http: &http.Client{Timeout: cfg.Timeout}}
}

func (c *client) Consultar(ctx context.Context, dni string, operacion string) (domain.ReniecResult, error) {
	metodo, ok := operaciones[operacion]
	if !ok {
		return domain.ReniecResult{}, fmt.Errorf("%w: %q", domain.ErrInvalidReniecOperation, operacion)
	}

	body := c.construirSOAP(metodo, dni)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.URL, bytes.NewBuffer(body))
	if err != nil {
		return domain.ReniecResult{}, fmt.Errorf("building reniec request: %w", err)
	}
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("SOAPAction", tempURINS+metodo)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/xml, application/xml, */*")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.ReniecResult{}, fmt.Errorf("calling reniec service: %w", err)
	}
	defer resp.Body.Close()

	contenido, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if readErr != nil {
		return domain.ReniecResult{}, fmt.Errorf("reading reniec response: %w", readErr)
	}

	if resp.StatusCode >= 400 {
		detalle := truncate(contenido, 800)
		return domain.ReniecResult{}, fmt.Errorf("reniec service responded %d: %s", resp.StatusCode, detalle)
	}

	resultado, parseErr := parsearStrings(contenido)
	if parseErr != nil {
		return domain.ReniecResult{}, fmt.Errorf("parsing reniec response: %w", parseErr)
	}

	if esError(resultado) {
		codigo, mensaje := interpretarError(resultado)
		log.Printf("reniec error dni=%s operacion=%s codigo=%s: %s resultado=%q",
			dni, operacion, codigo, mensaje, resultado)
		return domain.ReniecResult{}, fmt.Errorf("reniec error %s: %s", codigo, mensaje)
	}

	return domain.ReniecResult{
		DNI:       dni,
		Operacion: operacion,
		Resultado: resultado,
		Datos:     interpretarDatos(resultado, operacion),
	}, nil
}

// Layout de las respuestas SOAP de RENIEC. El servicio devuelve un arreglo
// plano de elementos <string> en un orden fijo. Estas posiciones están
// verificadas contra respuestas reales del servicio: el campo N de la
// documentación de RENIEC corresponde al índice N-1 del arreglo.
//
// Header común a las dos operaciones: [0] código de error ("0000" = consulta
// exitosa) y [1] descripción del error (vacía cuando la consulta es exitosa).
const (
	idxCodigoError = 0
	idxDetalle     = 1
)

// Layout de obtenerDatosCompletos (49 campos). El bloque de domicilio ocupa
// los índices 8-19 y el de nacimiento los 23-28; el estado civil y el grado de
// instrucción vienen como códigos, no como descripciones.
const (
	idxCompNroDoc      = 2
	idxCompPaterno     = 4
	idxCompMaterno     = 5
	idxCompNombres     = 7
	idxCompCodDeptoDom = 10
	idxCompCodProvDom  = 11
	idxCompCodDistDom  = 12
	idxCompDeptoDom    = 16
	idxCompProvDom     = 17
	idxCompDistDom     = 18
	idxCompEstadoCivil = 20
	idxCompSexo        = 22
	idxCompCodDeptoNac = 23
	idxCompCodProvNac  = 24
	idxCompCodDistNac  = 25
	idxCompDeptoNac    = 26
	idxCompProvNac     = 27
	idxCompDistNac     = 28
	idxCompFechaNac    = 29
	idxCompNombrePadre = 30
	idxCompNombreMadre = 31
	idxCompDireccion   = 36
)

// Layout de obtenerDatosBasicos (23 campos). No incluye el número de documento
// ni el dígito de verificación al inicio; el domicilio va completo y el sexo
// aparece justo después de la dirección.
const (
	idxBasPaterno     = 2
	idxBasMaterno     = 3
	idxBasNombres     = 5
	idxBasCodDeptoDom = 8
	idxBasCodProvDom  = 9
	idxBasCodDistDom  = 10
	idxBasDeptoDom    = 14
	idxBasProvDom     = 15
	idxBasDistDom     = 16
	idxBasDireccion   = 18
	idxBasSexo        = 19
	idxBasFechaNac    = 20
	idxBasNroDoc      = 22
)

// codigoEstadoCivilOK es el único valor que RENIEC devuelve en el campo de
// código de error cuando la consulta se resolvió correctamente.
const codigoEstadoCivilOK = "0000"

// interpretarDatos extrae los campos de la persona desde el arreglo crudo
// usando las posiciones verificadas de cada operación SOAP.
func interpretarDatos(resultado []string, operacion string) domain.ReniecDatos {
	datos := domain.ReniecDatos{}

	switch operacion {
	case "basico":
		datos.ApellidoPaterno = tokenEn(resultado, idxBasPaterno)
		datos.ApellidoMaterno = tokenEn(resultado, idxBasMaterno)
		datos.Nombres = tokenEn(resultado, idxBasNombres)
		datos.Sexo = sexoDesdeCodigo(tokenEn(resultado, idxBasSexo))
		datos.FechaNacimiento = convertirFecha(tokenEn(resultado, idxBasFechaNac))
		datos.Departamento = tokenEn(resultado, idxBasDeptoDom)
		datos.Provincia = tokenEn(resultado, idxBasProvDom)
		datos.Distrito = tokenEn(resultado, idxBasDistDom)
		datos.Direccion = tokenEn(resultado, idxBasDireccion)
		datos.Ubigeo = unirUbigeo(
			tokenEn(resultado, idxBasCodDeptoDom), tokenEn(resultado, idxBasCodProvDom), tokenEn(resultado, idxBasCodDistDom),
		)
	default:
		datos.ApellidoPaterno = tokenEn(resultado, idxCompPaterno)
		datos.ApellidoMaterno = tokenEn(resultado, idxCompMaterno)
		datos.Nombres = tokenEn(resultado, idxCompNombres)
		datos.Sexo = sexoDesdeCodigo(tokenEn(resultado, idxCompSexo))
		datos.FechaNacimiento = convertirFecha(tokenEn(resultado, idxCompFechaNac))
		datos.NombrePadre = tokenEn(resultado, idxCompNombrePadre)
		datos.NombreMadre = tokenEn(resultado, idxCompNombreMadre)
		// Domicilio: nombres del ubigeo en [16]/[17]/[18] y dirección en [36].
		datos.Departamento = tokenEn(resultado, idxCompDeptoDom)
		datos.Provincia = tokenEn(resultado, idxCompProvDom)
		datos.Distrito = tokenEn(resultado, idxCompDistDom)
		datos.Direccion = tokenEn(resultado, idxCompDireccion)
		datos.Ubigeo = unirUbigeo(
			tokenEn(resultado, idxCompCodDeptoDom), tokenEn(resultado, idxCompCodProvDom), tokenEn(resultado, idxCompCodDistDom),
		)
		// Nacimiento: nombres del ubigeo en [26]/[27]/[28].
		datos.DepartamentoNacimiento = tokenEn(resultado, idxCompDeptoNac)
		datos.ProvinciaNacimiento = tokenEn(resultado, idxCompProvNac)
		datos.DistritoNacimiento = tokenEn(resultado, idxCompDistNac)
		datos.EstadoCivil = estadoCivilDesdeCodigo(tokenEn(resultado, idxCompEstadoCivil))
	}

	partes := separarNombres(datos.Nombres)
	if len(partes) > 0 {
		datos.PrimerNombre = partes[0]
	}
	if len(partes) > 1 {
		datos.SegundoNombre = partes[1]
	}
	if len(partes) > 2 {
		datos.TercerNombre = partes[2]
	}

	return datos
}

// unirUbigeo concatena los códigos de departamento, provincia y distrito en el
// ubigeo de 6 dígitos que usa el catálogo. Si falta alguna piece devuelve
// vacío para no construir un código inválido.
func unirUbigeo(departamento, provincia, distrito string) string {
	if departamento == "" || provincia == "" || distrito == "" {
		return ""
	}
	ubigeo := departamento + provincia + distrito
	if !esUbigeo(ubigeo) {
		return ""
	}
	return ubigeo
}

// separarNombres divide los prenombres (e.g. "CARLOS MELECIO") en primer y
// segundo nombre; cualquier parte extra va al tercer nombre.
func separarNombres(nombres string) []string {
	partes := strings.Fields(nombres)
	if len(partes) <= 2 {
		return partes
	}
	tercero := strings.Join(partes[2:], " ")
	return []string{partes[0], partes[1], tercero}
}

// convertirFecha normaliza la fecha RENIEC (DD/MM/YYYY) a formato ISO YYYY-MM-DD
// para que el input tipo date del frontend la consuma de forma directa.
func convertirFecha(fecha string) string {
	fecha = strings.TrimSpace(fecha)
	if fecha == "" {
		return ""
	}
	partes := strings.Split(fecha, "/")
	if len(partes) != 3 {
		return fecha
	}
	dd, mm, aa := partes[0], partes[1], partes[2]
	if mm == "" || dd == "" || len(aa) != 4 || !isDigit(dd) || !isDigit(mm) {
		return fecha
	}
	return aa + "-" + mm + "-" + dd
}

// construirSOAP arma el sobre SOAP con las credenciales en el header,
// igual que el proyecto FastAPI.
func (c *client) construirSOAP(metodo string, nrodoc string) []byte {
	body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns:xsd="http://www.w3.org/2001/XMLSchema"
 xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <Credencialmq xmlns="http://tempuri.org/">
      <app>%s</app>
      <usuario>%s</usuario>
      <clave>%s</clave>
    </Credencialmq>
  </soap:Header>
  <soap:Body>
    <%s xmlns="http://tempuri.org/">
      <nrodoc>%s</nrodoc>
    </%s>
  </soap:Body>
</soap:Envelope>`,
		escapeXML(c.cfg.App),
		escapeXML(c.cfg.Usuario),
		escapeXML(c.cfg.Clave),
		metodo,
		escapeXML(nrodoc),
		metodo,
	)
	return []byte(body)
}

// parsearStrings recorre el XML de la respuesta y extrae el texto de todos
// los elementos <string>, que es el formato que RENIEC devuelve tanto para
// datos como para mensajes de error.
func parsearStrings(contenido []byte) ([]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(contenido))
	valores := make([]string, 0)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "string" {
				var valor string
				if err := decoder.DecodeElement(&valor, &t); err != nil {
					return nil, err
				}
				valores = append(valores, valor)
			}
		}
	}

	return valores, nil
}

// esError detecta el fallo de la consulta a RENIEC. El servicio no devuelve
// un código HTTP de error ni un Fault SOAP: ante una consulta inválida
// responde un arreglo cuyo primer campo es el código de error (por ejemplo
// "5114" para un documento inexistente) y el resto de campos quedan vacíos.
// Solo "0000" significa consulta exitosa.
func esError(resultado []string) bool {
	if len(resultado) <= idxDetalle {
		return false
	}
	codigo := strings.TrimSpace(resultado[idxCodigoError])
	return codigo != "" && codigo != codigoEstadoCivilOK
}

func interpretarError(resultado []string) (string, string) {
	mensajes := make([]string, 0, len(resultado)-1)
	for _, m := range resultado[idxDetalle:] {
		m = strings.TrimSpace(m)
		if m != "" {
			mensajes = append(mensajes, m)
		}
	}
	mensaje := strings.Join(mensajes, " ")
	if mensaje == "" {
		mensaje = "Error desconocido"
	}
	return resultado[idxCodigoError], mensaje
}

// sexoDesdeCodigo traduce el código de sexo de RENIEC (1 = masculino,
// 2 = femenino) al texto que consume el frontend. Acepta también el texto
// literal por si el servicio lo devolviera en ese formato.
func sexoDesdeCodigo(codigo string) string {
	switch strings.ToUpper(strings.TrimSpace(codigo)) {
	case "1", "M", "MASCULINO":
		return "MASCULINO"
	case "2", "F", "FEMENINO":
		return "FEMENINO"
	}
	return ""
}

// estadoCivilDesdeCodigo traduce el código de estado civil de RENIEC (1 =
// soltero, 2 = casado, 3 = viudo, 4 = divorciado, 5 = separado, 6 =
// conviviente) al texto que consume el frontend. Los códigos de RENIEC
// coinciden con los ids del catálogo de estados civiles.
func estadoCivilDesdeCodigo(codigo string) string {
	switch strings.TrimSpace(codigo) {
	case "1":
		return "SOLTERO"
	case "2":
		return "CASADO"
	case "3":
		return "VIUDO"
	case "4":
		return "DIVORCIADO"
	case "5":
		return "SEPARADO"
	case "6":
		return "CONVIVIENTE"
	}
	return ""
}

// esUbigeo verifica si el token es un código ubigeo (6 dígitos).
func esUbigeo(token string) bool {
	return len(token) == 6 && isDigit(token)
}

// tokenEn devuelve el token de la posición i. Si la posición está fuera de
// rango o RENIEC devolvió su nulo ("SIN DATOS", "S/D") devuelve vacío.
func tokenEn(resultado []string, i int) string {
	if i < 0 || i >= len(resultado) {
		return ""
	}
	v := strings.TrimSpace(resultado[i])
	if v == "" || esSinDatos(v) {
		return ""
	}
	return v
}

// esSinDatos verifica si el token es el nulo que usa RENIEC.
func esSinDatos(token string) bool {
	t := strings.ToUpper(strings.TrimSpace(token))
	return t == "SIN DATOS" || t == "S/D" || t == "SIN DATO"
}

func isDigit(s string) bool {
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

func escapeXML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return r.Replace(s)
}

func truncate(b []byte, max int) string {
	s := string(b)
	if len(s) > max {
		s = s[:max]
	}
	return s
}
