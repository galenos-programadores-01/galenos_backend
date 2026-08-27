package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type DiagnosticoUseCase interface {
	SearchDiagnosticos(ctx context.Context, filtro string, idAtencion, idPaciente int) ([]domain.DiagnosticoBusqueda, error)
	ListarDiagnosticos(ctx context.Context, filtro string) ([]domain.DiagnosticoSimple, error)
}
