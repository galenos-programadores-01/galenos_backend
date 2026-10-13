package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

// ArchivoRemotoService fachada reutilizable para subir/bajar PDFs firmados (PA hoy; evolución/módulo listos para otros flujos).
type ArchivoRemotoService struct {
	store output.AlmacenamientoRemoto
	rutas domain.RutasArchivoRemoto
}

// NewArchivoRemotoService cablea el adaptador (SFTP/noop) y las rutas configurables.
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

// Rutas expone el builder de paths (por si un modulo necesita la ubicación lógica sin subir).
func (s *ArchivoRemotoService) Rutas() domain.RutasArchivoRemoto {
	return s.rutas
}

// SubirPrimeraAtencion sube el PDF firmado de PA.
// Devuelve ubicacionLogica para persistir en BD (p. ej. /ftp/primera-atencion/...).
func (s *ArchivoRemotoService) SubirPrimeraAtencion(
	ctx context.Context,
	idAtencion int,
	uuidFirma string,
	pdf []byte,
) (ubicacionLogica string, err error) {
	if !s.Enabled() {
		return "", domain.ErrAlmacenamientoNoConfigurado
	}
	if idAtencion <= 0 || uuidFirma == "" || len(pdf) == 0 {
		return "", fmt.Errorf("archivo remoto PA: idAtencion, uuidFirma y pdf son obligatorios")
	}
	ruta, ubi, err := s.rutas.PrimeraAtencionPDF(idAtencion, uuidFirma)
	if err != nil {
		return "", err
	}
	if err := s.store.Subir(ctx, ruta, pdf); err != nil {
		return "", err
	}
	return ubi, nil
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
	ruta, ubi, err := s.rutas.EvolucionPDF(idRegAtencion, uuidFirma)
	if err != nil {
		return "", err
	}
	if err := s.store.Subir(ctx, ruta, pdf); err != nil {
		return "", err
	}
	return ubi, nil
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
