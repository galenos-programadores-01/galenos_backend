package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type medicoListaEsperaService struct {
	repo output.MedicoListaEsperaRepository
}

func NewMedicoListaEsperaService(repo output.MedicoListaEsperaRepository) *medicoListaEsperaService {
	return &medicoListaEsperaService{repo: repo}
}

func (s *medicoListaEsperaService) Listar(ctx context.Context) ([]domain.MedicoListaEspera, error) {
	return s.repo.Listar(ctx)
}
