package usecase

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

// ArchivoRemotoService fachada reutilizable para subir/bajar PDFs firmados.
// Los módulos resuelven rutas con patrones env vía SubirPorPatron / helpers.
type ArchivoRemotoService struct {
	store output.AlmacenamientoRemoto
	rutas domain.RutasArchivoRemoto
}

// NewArchivoRemotoService cablea el adaptador (SFTP/FTP/noop) y las rutas configurables.
func NewArchivoRemotoService(
	store output.AlmacenamientoRemoto,
	rutas domain.RutasArchivoRemoto,
) *ArchivoRemotoService {
	if store == nil {
		store = noopAlmacenamiento{}
	}
	return &ArchivoRemotoService{store: store, rutas: rutas}
}

// Enabled true si el adaptador tiene credenciales/host.
func (s *ArchivoRemotoService) Enabled() bool {
	return s != nil && s.store != nil && s.store.Enabled()
}

// Rutas expone el builder de paths (ubicación lógica sin subir).
func (s *ArchivoRemotoService) Rutas() domain.RutasArchivoRemoto {
	return s.rutas
}

// SubirPorPatron resuelve SFTP_PATTERN_* + vars, sube el PDF y devuelve
// ubicacionLogica (misma cadena que debe persistirse en BD si aplica).
// Reutilizable por cualquier módulo: solo pasa clave de patrón y variables.
func (s *ArchivoRemotoService) SubirPorPatron(
	ctx context.Context,
	clavePatron string,
	vars map[string]string,
	pdf []byte,
) (ubicacionLogica string, err error) {
	if !s.Enabled() {
		return "", domain.ErrAlmacenamientoNoConfigurado
	}
	if len(pdf) == 0 {
		return "", fmt.Errorf("archivo remoto: pdf vacío")
	}
	ruta, ubi, err := s.rutas.Construir(clavePatron, vars)
	if err != nil {
		return "", err
	}
	if err := s.store.Subir(ctx, ruta, pdf); err != nil {
		return "", err
	}
	return ubi, nil
}

// SubirPrimeraAtencion sube el PDF principal de PA.
// Devuelve ubicacionLogica para Tab_PA.UbicacionArchivo.
func (s *ArchivoRemotoService) SubirPrimeraAtencion(
	ctx context.Context,
	idPrimeraAtencion int,
	pdf []byte,
) (ubicacionLogica string, err error) {
	if idPrimeraAtencion <= 0 {
		return "", fmt.Errorf("archivo remoto PA: idPrimeraAtencion inválido")
	}
	return s.SubirPorPatron(ctx, domain.PatronArchivoPrimeraAtencion, map[string]string{
		"idPrimeraAtencion": fmt.Sprintf("%d", idPrimeraAtencion),
		"ext":               "pdf",
	}, pdf)
}

// SubirPrimeraAtencionReceta sube un PDF satélite (examen/meds) por IdReceta.
// No se guarda en Tab_PA; la ruta se reconstruye con el patrón.
// SubirPrimeraAtencionConNombre sube el PDF con el basename ya decidido
// (incluye -N de una edición) en la misma carpeta del PDF principal.
func (s *ArchivoRemotoService) SubirPrimeraAtencionConNombre(
	ctx context.Context,
	idPrimeraAtencion int,
	nombreArchivo string,
	pdf []byte,
) (ubicacionLogica string, err error) {
	if !s.Enabled() {
		return "", domain.ErrAlmacenamientoNoConfigurado
	}
	_, base, err := s.rutas.PrimeraAtencionPDF(idPrimeraAtencion)
	if err != nil {
		return "", err
	}
	nombreArchivo = path.Base(strings.ReplaceAll(nombreArchivo, "\\", "/"))
	if nombreArchivo == "" || nombreArchivo == "." {
		return "", domain.ErrAlmacenamientoRutaInvalida
	}
	dir := path.Dir(base)
	logica := path.Join(dir, nombreArchivo)
	if err := s.store.Subir(ctx, logica, pdf); err != nil {
		return "", err
	}
	return logica, nil
}

func (s *ArchivoRemotoService) SubirPrimeraAtencionReceta(
	ctx context.Context,
	idPrimeraAtencion, idReceta int,
	pdf []byte,
) (ubicacionLogica string, err error) {
	if idPrimeraAtencion <= 0 || idReceta <= 0 {
		return "", fmt.Errorf("archivo remoto PA receta: idPrimeraAtencion e idReceta son obligatorios")
	}
	return s.SubirPorPatron(ctx, domain.PatronArchivoPrimeraAtencionReceta, map[string]string{
		"idPrimeraAtencion": fmt.Sprintf("%d", idPrimeraAtencion),
		"idReceta":          fmt.Sprintf("%d", idReceta),
		"ext":               "pdf",
	}, pdf)
}

// SubirEvolucion sube el PDF firmado de evolución médica.
func (s *ArchivoRemotoService) SubirEvolucion(
	ctx context.Context,
	idRegAtencion int,
	uuidFirma string,
	pdf []byte,
) (ubicacionLogica string, err error) {
	if !s.Enabled() {
		return "", domain.ErrAlmacenamientoNoConfigurado
	}
	if idRegAtencion <= 0 || uuidFirma == "" || len(pdf) == 0 {
		return "", fmt.Errorf("archivo remoto evolución: idRegAtencion, uuidFirma y pdf son obligatorios")
	}
	return s.SubirPorPatron(ctx, domain.PatronArchivoEvolucion, map[string]string{
		"id":   fmt.Sprintf("%d", idRegAtencion),
		"uuid": uuidFirma,
		"ext":  "pdf",
	}, pdf)
}

// SubirModulo sube un PDF con patrón genérico {base}/{modulo}/{id}/{uuid}.pdf.
func (s *ArchivoRemotoService) SubirModulo(
	ctx context.Context,
	modulo string,
	id int,
	uuidFirma string,
	pdf []byte,
) (ubicacionLogica string, err error) {
	if !s.Enabled() {
		return "", domain.ErrAlmacenamientoNoConfigurado
	}
	ruta, ubi, err := s.rutas.ModuloPDF(modulo, id, uuidFirma)
	if err != nil {
		return "", err
	}
	if err := s.store.Subir(ctx, ruta, pdf); err != nil {
		return "", err
	}
	return ubi, nil
}

// DescargarPorUbicacion baja un archivo usando la ruta lógica/remota guardada en BD.
func (s *ArchivoRemotoService) DescargarPorUbicacion(ctx context.Context, ubicacion string) ([]byte, error) {
	if !s.Enabled() {
		return nil, domain.ErrAlmacenamientoNoConfigurado
	}
	return s.store.Descargar(ctx, ubicacion)
}

// noop local para constructor nil-safe (mismo contrato que sftp.noopClient).
type noopAlmacenamiento struct{}

func (noopAlmacenamiento) Enabled() bool { return false }

func (noopAlmacenamiento) Subir(context.Context, string, []byte) error {
	return domain.ErrAlmacenamientoNoConfigurado
}

func (noopAlmacenamiento) Descargar(context.Context, string) ([]byte, error) {
	return nil, domain.ErrAlmacenamientoNoConfigurado
}

func (noopAlmacenamiento) Existe(context.Context, string) (bool, error) {
	return false, domain.ErrAlmacenamientoNoConfigurado
}

func (noopAlmacenamiento) Eliminar(context.Context, string) error {
	return domain.ErrAlmacenamientoNoConfigurado
}
