// Package output define los contratos que el módulo refcon necesita de sus
// adaptadores de salida (repositorio local y cliente externo MINSA).
package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
)

type RefConRepository interface {
	ListarReferenciasPorMes(ctx context.Context, mes, anio, estado, opcion int, ups string) ([]domain.ReferenciaPorMes, error)
	ListarUps(ctx context.Context) ([]domain.Ups, error)
	ListarTodasLasUps(ctx context.Context) ([]domain.Ups, error)
	ListarEstablecimientos(ctx context.Context) ([]domain.Establecimiento, error)
	ListarDistritosPorIdReniec(ctx context.Context, idReniec int) ([]domain.DistritoReniec, error)
}

// MinsaRefConClient consume el servicio REST de referencias del MINSA
// (consultaReferenciaDetalle y listadoUps).
type MinsaRefConClient interface {
	ConsultarReferenciaDetalle(ctx context.Context, req domain.ConsultaMinsaRequest) (*domain.ConsultaMinsaResponse, error)
	ListadoUpss(ctx context.Context, codigoRenipress string) (*domain.ListadoUpsResponse, error)
}

// RefConReportsClient consume el portal REFCON (refcon.minsa.gob.pe) haciendo
// login por cookies para generar y descargar la hoja de referencia en PDF.
type RefConReportsClient interface {
	GenerarHojaReferencia(ctx context.Context, req domain.GenerarHojaReferenciaRequest) (*domain.GenerarHojaReferenciaResult, error)
}
