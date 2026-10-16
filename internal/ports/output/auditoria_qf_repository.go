package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// AuditoriaQfRepository llama a los SP de la bandeja de farmacia.
type AuditoriaQfRepository interface {
	ListarBandeja(ctx context.Context, params domain.BandejaRecetasParams) ([]map[string]any, error)
	ListarEstados(ctx context.Context) ([]map[string]any, error)
	ListarDetalle(ctx context.Context, params domain.DetalleRecetaParams) ([]map[string]any, error)
	ListarHistorico(ctx context.Context, params domain.HistoricoMedicamentoParams) ([]map[string]any, error)
}
