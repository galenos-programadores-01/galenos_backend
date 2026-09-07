package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type RefConService interface {
	ListarReferenciasPorMes(ctx context.Context, mes, anio, estado, opcion int, ups string) ([]domain.ReferenciaPorMes, error)
	ListarUps(ctx context.Context) ([]domain.Ups, error)
}
