package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type causaExternaMorbilidadService struct {
	repo output.CausaExternaMorbilidadRepository
}

// NewCausaExternaMorbilidadService construye el servicio del catálogo de
// causas externas de morbilidad.
func NewCausaExternaMorbilidadService(repo output.CausaExternaMorbilidadRepository) input.CausaExternaMorbilidadService {
	return &causaExternaMorbilidadService{repo: repo}
}

// Listar delega en el repositorio (SP
// uso_go_EmergenciaCausaExternaMorbilidadListar) y devuelve el catálogo.
func (s *causaExternaMorbilidadService) Listar(ctx context.Context) ([]domain.CausaExternaMorbilidad, error) {
	items, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing causes of external morbidity: %w", err)
	}
	return items, nil
}
