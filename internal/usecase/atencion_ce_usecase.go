package usecase

import (
	"context"
	"errors"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type atencionCEUseCase struct {
	repo output.AtencionCERepository
}

func NewAtencionCEUseCase(repo output.AtencionCERepository) input.AtencionCEService {
	return &atencionCEUseCase{
		repo: repo,
	}
}

func (uc *atencionCEUseCase) ListarColaPacientes(ctx context.Context, filtros domain.ColaFiltrosCE) ([]domain.ColaPacienteCEItem, error) {
	return uc.repo.ListarColaPacientes(ctx, filtros)
}

func (uc *atencionCEUseCase) ObtenerFichaAtencion(ctx context.Context, idAtencion int) (*domain.FichaAtencionCEDetalle, error) {
	if idAtencion <= 0 {
		return nil, errors.New("identificador de atencion invalido")
	}
	return uc.repo.ObtenerFichaAtencion(ctx, idAtencion)
}

func (uc *atencionCEUseCase) GuardarFichaAtencion(ctx context.Context, payload domain.AtencionCEGuardarPayload) error {
	if payload.IdAtencion <= 0 {
		return errors.New("identificador de atencion requerido para guardar la ficha")
	}
	return uc.repo.GuardarFichaAtencion(ctx, payload)
}

func (uc *atencionCEUseCase) ListarDestinosAtencion(ctx context.Context) ([]domain.DestinoAtencionItem, error) {
	return uc.repo.ListarDestinosAtencion(ctx)
}

func (uc *atencionCEUseCase) ListarCatalogoLabsHIS(ctx context.Context) ([]domain.CatalogoLabHISItem, error) {
	return uc.repo.ListarCatalogoLabsHIS(ctx)
}
