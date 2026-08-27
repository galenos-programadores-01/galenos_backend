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

func (s *listaEsperaQxService) Listar(ctx context.Context, fecha string, fechaFin string, paciente string, idEspecialidad int) ([]domain.ListaEsperaQx, error) {
	return s.repo.Listar(ctx, fecha, fechaFin, paciente, idEspecialidad)
}

func (s *listaEsperaQxService) ObtenerPorId(ctx context.Context, id int) (domain.ListaEsperaQxPaciente, error) {
	return s.repo.ObtenerPorId(ctx, id)
}

func (s *listaEsperaQxService) Crear(ctx context.Context, item domain.ListaEsperaQxCrear, idUsuario int) error {
	return s.repo.Crear(ctx, item, idUsuario)
}

func (s *listaEsperaQxService) Modificar(ctx context.Context, item domain.ListaEsperaQxModificar) error {
	return s.repo.Modificar(ctx, item)
}

func (s *listaEsperaQxService) Reporte(ctx context.Context, fecha string, fechaFin string, idEspecialidad int) ([]domain.ListaEsperaQxReporte, error) {
	return s.repo.Reporte(ctx, fecha, fechaFin, idEspecialidad)
}
