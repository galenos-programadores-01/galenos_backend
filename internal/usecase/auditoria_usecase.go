package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type auditoriaService struct {
	repo output.AuditoriaRepository
}

func NewAuditoriaService(repo output.AuditoriaRepository) *auditoriaService {
	return &auditoriaService{repo: repo}
}

func (s *auditoriaService) Registrar(ctx context.Context, reg domain.AuditoriaRegistro) error {
	return s.repo.Registrar(ctx, reg)
}
