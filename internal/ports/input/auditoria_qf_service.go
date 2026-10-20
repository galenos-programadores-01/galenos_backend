package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// AuditoriaQfService es el puerto de entrada de Auditoría QF.
type AuditoriaQfService interface {
	ListarBandeja(ctx context.Context, params domain.BandejaRecetasParams) ([]map[string]any, error)
	ListarEstados(ctx context.Context) ([]map[string]any, error)
	ListarDetalle(ctx context.Context, params domain.DetalleRecetaParams) ([]map[string]any, error)
	ListarHistorico(ctx context.Context, params domain.HistoricoMedicamentoParams) ([]map[string]any, error)
	GuardarEstado(ctx context.Context, params domain.GuardarEstadoRecetaParams) (string, error)
}
