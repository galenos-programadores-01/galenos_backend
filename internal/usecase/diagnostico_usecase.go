package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type diagnosticoUseCase struct {
	repo output.DiagnosticoRepository
}

func NewDiagnosticoUseCase(repo output.DiagnosticoRepository) input.DiagnosticoUseCase {
	return &diagnosticoUseCase{repo: repo}
}

func (uc *diagnosticoUseCase) SearchDiagnosticos(ctx context.Context, filtro string, idAtencion, idPaciente int) ([]domain.DiagnosticoBusqueda, error) {
	return uc.repo.SearchDiagnosticos(ctx, filtro, idAtencion, idPaciente)
}

func (uc *diagnosticoUseCase) ListarDiagnosticos(ctx context.Context, filtro string) ([]domain.DiagnosticoSimple, error) {
	return uc.repo.ListarDiagnosticos(ctx, filtro)
}
