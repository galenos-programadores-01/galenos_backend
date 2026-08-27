package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type DiagnosticoRepository interface {
	SearchDiagnosticos(ctx context.Context, filtro string, idAtencion, idPaciente int) ([]domain.DiagnosticoBusqueda, error)
	ListarDiagnosticos(ctx context.Context, filtro string) ([]domain.DiagnosticoSimple, error)
}
