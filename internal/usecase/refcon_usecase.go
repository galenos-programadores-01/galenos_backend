package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type refconService struct {
	repo output.RefConRepository
}

func NewRefConService(repo output.RefConRepository) *refconService {
	return &refconService{repo: repo}
}

func (s *refconService) ListarReferenciasPorMes(ctx context.Context, mes, anio, estado, opcion int, ups string) ([]domain.ReferenciaPorMes, error) {
	return s.repo.ListarReferenciasPorMes(ctx, mes, anio, estado, opcion, ups)
}

func (s *refconService) ListarUps(ctx context.Context) ([]domain.Ups, error) {
	return s.repo.ListarUps(ctx)
}
