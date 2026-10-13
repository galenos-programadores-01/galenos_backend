package domain

import (
	"errors"
	"strconv"
	"strings"
)

// Errores del almacenamiento remoto (SFTP u otro adaptador).
var (
	ErrAlmacenamientoNoConfigurado       = errors.New("almacenamiento remoto no configurado")
	ErrAlmacenamientoArchivoNoEncontrado = errors.New("archivo remoto no encontrado")
	ErrAlmacenamientoPatronInvalido      = errors.New("patrón de ruta remota inválido")
	ErrAlmacenamientoRutaInvalida        = errors.New("ruta remota inválida")
)

// Claves de patrón configurables (env SFTP_PATTERN_*).
const (
	PatronArchivoPrimeraAtencion = "primera_atencion"
	PatronArchivoEvolucion       = "evolucion"
)

// RutasArchivoRemoto construye rutas remotas y la ubicación lógica de BD
// a partir de patrones configurables ({id}, {uuid}, {ext}).
// Es puro (sin I/O): lo usan casos de uso y el adaptador SFTP.
type RutasArchivoRemoto struct {
	// BasePath prefijo remoto, p. ej. "/ftp".
	BasePath string
	// Patterns mapa clave → plantilla relativa, p. ej.
	// "primera_atencion" → "primera-atencion/{id}/{uuid}.pdf"
	Patterns map[string]string
}

// DefaultsRutasArchivoRemoto patrones alineados con UbicacionArchivo actuales.
func DefaultsRutasArchivoRemoto() map[string]string {
	return map[string]string{
		PatronArchivoPrimeraAtencion: "primera-atencion/{id}/{uuid}.pdf",
		PatronArchivoEvolucion:       "evoluciones/{id}/{uuid}.pdf",
	}
}

// NuevaRutasArchivoRemoto normaliza base y completa patrones faltantes con defaults.
func NuevaRutasArchivoRemoto(basePath string, patterns map[string]string) RutasArchivoRemoto {
	base := normalizeBasePath(basePath)
	merged := DefaultsRutasArchivoRemoto()
	for k, v := range patterns {
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if k == "" || v == "" {
			continue
		}
		merged[k] = trimSlashes(v)
	}
	return RutasArchivoRemoto{BasePath: base, Patterns: merged}
}

// Construir resuelve un patrón por clave y variables.
// Devuelve ruta remota (para SFTP) y ubicación lógica (para BD, p. ej. /ftp/...).
func (r RutasArchivoRemoto) Construir(clavePatron string, vars map[string]string) (rutaRemota, ubicacionLogica string, err error) {
	clavePatron = strings.TrimSpace(clavePatron)
	patron, ok := r.Patterns[clavePatron]
	if !ok || patron == "" {
		return "", "", ErrAlmacenamientoPatronInvalido
	}
	rel := applyPattern(patron, vars)
	if rel == "" || containsDotDot(rel) {
		return "", "", ErrAlmacenamientoRutaInvalida
	}
	rutaRemota = joinRemote(r.BasePath, rel)
	ubicacionLogica = ensureLeadingSlash(rutaRemota)
	return rutaRemota, ubicacionLogica, nil
}

// PrimeraAtencionPDF ruta para el PDF firmado de PA.
func (r RutasArchivoRemoto) PrimeraAtencionPDF(idAtencion int, uuidFirma string) (rutaRemota, ubicacionLogica string, err error) {
	return r.Construir(PatronArchivoPrimeraAtencion, map[string]string{
		"id":   strconv.Itoa(idAtencion),
		"uuid": strings.TrimSpace(uuidFirma),
		"ext":  "pdf",
	})
}

// EvolucionPDF ruta para el PDF firmado de evolución médica.
func (r RutasArchivoRemoto) EvolucionPDF(idRegAtencion int, uuidFirma string) (rutaRemota, ubicacionLogica string, err error) {
	return r.Construir(PatronArchivoEvolucion, map[string]string{
		"id":   strconv.Itoa(idRegAtencion),
		"uuid": strings.TrimSpace(uuidFirma),
		"ext":  "pdf",
	})
}

// ModuloPDF patrón genérico {modulo}/{id}/{uuid}.pdf para módulos nuevos.
func (r RutasArchivoRemoto) ModuloPDF(modulo string, id int, uuidFirma string) (rutaRemota, ubicacionLogica string, err error) {
	modulo = trimSlashes(strings.TrimSpace(modulo))
	uuidFirma = strings.TrimSpace(uuidFirma)
	if modulo == "" || uuidFirma == "" || id <= 0 {
		return "", "", ErrAlmacenamientoRutaInvalida
	}
	rel := modulo + "/" + strconv.Itoa(id) + "/" + uuidFirma + ".pdf"
	if containsDotDot(rel) {
		return "", "", ErrAlmacenamientoRutaInvalida
	}
	rutaRemota = joinRemote(r.BasePath, rel)
	ubicacionLogica = ensureLeadingSlash(rutaRemota)
	return rutaRemota, ubicacionLogica, nil
}

func normalizeBasePath(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return "/ftp"
	}
	return "/" + trimSlashes(base)
}

func joinRemote(base, rel string) string {
	base = trimSlashes(base)
	rel = trimSlashes(rel)
	if base == "" {
		return "/" + rel
	}
	return "/" + base + "/" + rel
}

func ensureLeadingSlash(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return p
	}
	if !strings.HasPrefix(p, "/") {
		return "/" + p
	}
	return p
}

func applyPattern(patron string, vars map[string]string) string {
	out := patron
	for k, v := range vars {
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	return trimSlashes(out)
}

func containsDotDot(p string) bool {
	for _, part := range strings.Split(p, "/") {
		if part == ".." {
			return true
		}
	}
	return false
}

func trimSlashes(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\\", "/")
	s = strings.Trim(s, "/")
	for strings.Contains(s, "//") {
		s = strings.ReplaceAll(s, "//", "/")
	}
	return s
}
