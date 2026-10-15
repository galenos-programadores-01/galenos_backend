package domain

import (
	"errors"
	"path"
	"strconv"
	"strings"
)

// Errores del almacenamiento remoto (SFTP u otro adaptador).
var (
	ErrAlmacenamientoNoConfigurado       = errors.New("almacenamiento remoto no configurado")
	ErrAlmacenamientoArchivoNoEncontrado = errors.New("archivo remoto no encontrado")
	ErrAlmacenamientoPatronInvalido      = errors.New("patrón de ruta remota inválido")
	ErrAlmacenamientoRutaInvalida        = errors.New("ruta remota inválida")
	ErrAlmacenamientoVarsIncompletas     = errors.New("faltan variables del patrón de ruta")
)

// Claves de patrón configurables (env SFTP_PATTERN_*).
const (
	PatronArchivoPrimeraAtencion       = "primera_atencion"
	PatronArchivoPrimeraAtencionReceta = "primera_atencion_receta"
	PatronArchivoEvolucion             = "evolucion"
)

// RutasArchivoRemoto construye rutas remotas y la ubicación lógica de BD
// a partir de patrones configurables. Es puro (sin I/O): lo usan casos de uso
// y el adaptador SFTP/FTP. Otros módulos reutilizan Construir / helpers.
type RutasArchivoRemoto struct {
	// BasePath prefijo remoto, p. ej. "/galenos-pro".
	BasePath string
	// Patterns mapa clave → plantilla relativa, p. ej.
	// "primera_atencion" → "emergencia/primera-atencion/{idPrimeraAtencion}.pdf"
	Patterns map[string]string
}

// DefaultsRutasArchivoRemoto patrones por defecto (PA + receta satélite + evolución).
func DefaultsRutasArchivoRemoto() map[string]string {
	return map[string]string{
		PatronArchivoPrimeraAtencion:       "emergencia/primera-atencion/{idPrimeraAtencion}.pdf",
		PatronArchivoPrimeraAtencionReceta: "emergencia/primera-atencion/{idPrimeraAtencion}_{idReceta}.pdf",
		PatronArchivoEvolucion:             "evoluciones/{id}/{uuid}.pdf",
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
// Devuelve ruta remota (SFTP/FTP) y ubicación lógica (BD).
// Falla si el patrón no existe, queda algún {placeholder} sin resolver o hay "..".
func (r RutasArchivoRemoto) Construir(clavePatron string, vars map[string]string) (rutaRemota, ubicacionLogica string, err error) {
	clavePatron = strings.TrimSpace(clavePatron)
	patron, ok := r.Patterns[clavePatron]
	if !ok || patron == "" {
		return "", "", ErrAlmacenamientoPatronInvalido
	}
	rel, err := applyPattern(patron, vars)
	if err != nil {
		return "", "", err
	}
	if rel == "" || containsDotDot(rel) {
		return "", "", ErrAlmacenamientoRutaInvalida
	}
	rutaRemota = joinRemote(r.BasePath, rel)
	ubicacionLogica = ensureLeadingSlash(rutaRemota)
	return rutaRemota, ubicacionLogica, nil
}

// PrimeraAtencionPDF ruta del PDF principal de PA (la que se guarda en Tab_PA.UbicacionArchivo).
func (r RutasArchivoRemoto) PrimeraAtencionPDF(idPrimeraAtencion int) (rutaRemota, ubicacionLogica string, err error) {
	if idPrimeraAtencion <= 0 {
		return "", "", ErrAlmacenamientoRutaInvalida
	}
	return r.Construir(PatronArchivoPrimeraAtencion, map[string]string{
		"idPrimeraAtencion": strconv.Itoa(idPrimeraAtencion),
		"ext":               "pdf",
	})
}

// PrimeraAtencionRecetaPDF ruta del PDF satélite (examen/meds) por IdReceta.
// No se persiste en Tab_PA; se reconstruye con el mismo patrón + IDs.
func (r RutasArchivoRemoto) PrimeraAtencionRecetaPDF(idPrimeraAtencion, idReceta int) (rutaRemota, ubicacionLogica string, err error) {
	if idPrimeraAtencion <= 0 || idReceta <= 0 {
		return "", "", ErrAlmacenamientoRutaInvalida
	}
	return r.Construir(PatronArchivoPrimeraAtencionReceta, map[string]string{
		"idPrimeraAtencion": strconv.Itoa(idPrimeraAtencion),
		"idReceta":          strconv.Itoa(idReceta),
		"ext":               "pdf",
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

// ModuloPDF patrón genérico {modulo}/{id}/{uuid}.pdf para módulos nuevos sin env propio.
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

// ArchivoLotePA clasifica un nombre del lote Firma Perú:
//   - "42.pdf"        → principal PA
//   - "42_98765.pdf"  → satélite receta
type ArchivoLotePA struct {
	EsPrincipal       bool
	IdPrimeraAtencion int
	IdReceta          int
}

// ParseNombreArchivoLotePA interpreta el basename del PDF del lote.
func ParseNombreArchivoLotePA(nombre string) (ArchivoLotePA, error) {
	base := strings.ToLower(path.Base(strings.ReplaceAll(strings.TrimSpace(nombre), "\\", "/")))
	if !strings.HasSuffix(base, ".pdf") {
		return ArchivoLotePA{}, ErrAlmacenamientoRutaInvalida
	}
	stem := strings.TrimSuffix(base, ".pdf")
	if stem == "" {
		return ArchivoLotePA{}, ErrAlmacenamientoRutaInvalida
	}
	if i := strings.IndexByte(stem, '_'); i >= 0 {
		idPA, err1 := strconv.Atoi(stem[:i])
		idReceta, err2 := strconv.Atoi(stem[i+1:])
		if err1 != nil || err2 != nil || idPA <= 0 || idReceta <= 0 {
			return ArchivoLotePA{}, ErrAlmacenamientoRutaInvalida
		}
		return ArchivoLotePA{EsPrincipal: false, IdPrimeraAtencion: idPA, IdReceta: idReceta}, nil
	}
	idPA, err := strconv.Atoi(stem)
	if err != nil || idPA <= 0 {
		return ArchivoLotePA{}, ErrAlmacenamientoRutaInvalida
	}
	return ArchivoLotePA{EsPrincipal: true, IdPrimeraAtencion: idPA}, nil
}

// NombreArchivoLotePA nombre canónico del PDF principal en el 7z.
func NombreArchivoLotePA(idPrimeraAtencion int) string {
	return strconv.Itoa(idPrimeraAtencion) + ".pdf"
}

// NombreArchivoLotePAReceta nombre canónico del PDF satélite en el 7z.
func NombreArchivoLotePAReceta(idPrimeraAtencion, idReceta int) string {
	return strconv.Itoa(idPrimeraAtencion) + "_" + strconv.Itoa(idReceta) + ".pdf"
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

func applyPattern(patron string, vars map[string]string) (string, error) {
	out := patron
	for k, v := range vars {
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	if strings.Contains(out, "{") {
		return "", ErrAlmacenamientoVarsIncompletas
	}
	return trimSlashes(out), nil
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
