package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type auditoriaQfUseCase struct {
	repo output.AuditoriaQfRepository
}

// NewAuditoriaQfUseCase construye el caso de uso de Auditoría QF.
func NewAuditoriaQfUseCase(repo output.AuditoriaQfRepository) input.AuditoriaQfService {
	return &auditoriaQfUseCase{repo: repo}
}

func (uc *auditoriaQfUseCase) ListarBandeja(ctx context.Context, params domain.BandejaRecetasParams) ([]map[string]any, error) {
	if err := validarRangoBandeja(params.FechaInicio, params.FechaFin); err != nil {
		return nil, err
	}
	items, err := uc.repo.ListarBandeja(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing auditoria qf bandeja: %w", err)
	}
	return items, nil
}

func (uc *auditoriaQfUseCase) ListarEstados(ctx context.Context) ([]map[string]any, error) {
	items, err := uc.repo.ListarEstados(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing receta estados: %w", err)
	}
	return items, nil
}

func (uc *auditoriaQfUseCase) ListarDetalle(ctx context.Context, params domain.DetalleRecetaParams) ([]map[string]any, error) {
	items, err := uc.repo.ListarDetalle(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing auditoria qf detalle: %w", err)
	}
	return items, nil
}

func (uc *auditoriaQfUseCase) ListarHistorico(ctx context.Context, params domain.HistoricoMedicamentoParams) ([]map[string]any, error) {
	items, err := uc.repo.ListarHistorico(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing auditoria qf historico: %w", err)
	}
	return items, nil
}

func validarRangoBandeja(inicioRaw, finRaw string) error {
	inicio, errInicio := time.Parse("2006-01-02T15:04:05", inicioRaw)
	fin, errFin := time.Parse("2006-01-02T15:04:05", finRaw)
	if errInicio != nil || errFin != nil {
		return domain.ErrFechaBandejaInvalida
	}
	inicioDia := time.Date(inicio.Year(), inicio.Month(), inicio.Day(), 0, 0, 0, 0, time.UTC)
	finDia := time.Date(fin.Year(), fin.Month(), fin.Day(), 0, 0, 0, 0, time.UTC)
	if finDia.Before(inicioDia) {
		return domain.ErrFechaBandejaOrden
	}
	dias := int(finDia.Sub(inicioDia).Hours() / 24)
	if dias > domain.DiasMaximosBandejaRecetas {
		return domain.ErrRangoFechasBandeja
	}
	return nil
}
