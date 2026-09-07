package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type DiagnosticoRepository interface {
	SearchDiagnosticos(ctx context.Context, filtro string, idAtencion, idPaciente int) ([]domain.DiagnosticoBusqueda, error)
	ListarDiagnosticos(ctx context.Context, filtro string) ([]domain.DiagnosticoSimple, error)
	ObtenerDiagnosticosAtencion(ctx context.Context, idAtencion int, idPrimeraAtencion *int, idEvolucion *int) ([]domain.DiagnosticoAtencion, error)
	AgregarDiagnosticoAtencion(ctx context.Context, req domain.AgregarDiagnosticoAtencionRequest) (*domain.AgregarDiagnosticoAtencionResponse, error)
}
