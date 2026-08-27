package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type listaEsperaQxService struct {
	repo output.ListaEsperaQxRepository
}

func NewListaEsperaQxService(repo output.ListaEsperaQxRepository) *listaEsperaQxService {
	return &listaEsperaQxService{repo: repo}
}

func (s *listaEsperaQxService) Listar(ctx context.Context, fecha string, paciente string) ([]domain.ListaEsperaQx, error) {
	return s.repo.Listar(ctx, fecha, paciente)
}

func (s *listaEsperaQxService) Crear(ctx context.Context, item domain.ListaEsperaQxCrear) error {
	return s.repo.Crear(ctx, item)
}
