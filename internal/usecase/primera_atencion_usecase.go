package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type primeraAtencionUseCase struct {
	repo output.PrimeraAtencionRepository
}

// NewPrimeraAtencionUseCase construye el caso de uso de Primera Atención.
func NewPrimeraAtencionUseCase(repo output.PrimeraAtencionRepository) input.PrimeraAtencionService {
	return &primeraAtencionUseCase{repo: repo}
}

// ListarBandeja invoca web_go_GestionPrimeraAtencionFiltrar y enriquece Estado
// con "En curso" cuando existe borrador en Tab_PrimeraAtencion (hasta que el SP
// de bandeja exponga estado oficial).
func (uc *primeraAtencionUseCase) ListarBandeja(ctx context.Context, params shared.PrimeraAtencionParams) ([]map[string]any, error) {
	items, err := uc.repo.ListarBandeja(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing primera atencion bandeja: %w", err)
	}

	ids := idsAtencionDesdeBandeja(items)
	conBorrador, err := uc.repo.IdsConBorradorActivo(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("checking borradores primera atencion: %w", err)
	}
	firmados, err := uc.repo.IdsFirmados(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("checking firmas primera atencion: %w", err)
	}
	enriquecerEstadoBandeja(items, conBorrador, firmados)

	return items, nil
}
