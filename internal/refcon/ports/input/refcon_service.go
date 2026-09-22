// Package input define los contratos que el módulo refcon expone a la capa
// HTTP (puertos de entrada).
package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
)

type RefConService interface {
	ListarReferenciasPorMes(ctx context.Context, mes, anio, estado, opcion int, ups string) ([]domain.ReferenciaPorMes, error)
	ListarUps(ctx context.Context) ([]domain.Ups, error)
	ListarTodasLasUps(ctx context.Context) ([]domain.Ups, error)
	ListarEstablecimientos(ctx context.Context) ([]domain.Establecimiento, error)
	ListarDistritosPorIdReniec(ctx context.Context, idReniec int) ([]domain.DistritoReniec, error)
	ConsultarReferenciaDetalle(ctx context.Context, req domain.ConsultaMinsaRequest) (*domain.ConsultaMinsaResponse, error)
	ListarUpssMinsa(ctx context.Context, codigoRenipress string) (*domain.ListadoUpsResponse, error)
	GenerarHojaReferencia(ctx context.Context, req domain.GenerarHojaReferenciaRequest) (*domain.GenerarHojaReferenciaResult, error)
}
