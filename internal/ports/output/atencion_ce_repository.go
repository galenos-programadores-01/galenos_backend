package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type AtencionCERepository interface {
	ListarColaPacientes(ctx context.Context, filtros domain.ColaFiltrosCE) ([]domain.ColaPacienteCEItem, error)
	ObtenerFichaAtencion(ctx context.Context, idAtencion int) (*domain.FichaAtencionCEDetalle, error)
	GuardarFichaAtencion(ctx context.Context, payload domain.AtencionCEGuardarPayload) error
	ListarDestinosAtencion(ctx context.Context) ([]domain.DestinoAtencionItem, error)
	ListarCatalogoLabsHIS(ctx context.Context) ([]domain.CatalogoLabHISItem, error)
}
