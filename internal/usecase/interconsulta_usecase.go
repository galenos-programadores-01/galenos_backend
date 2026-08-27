package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type interconsultaService struct {
	repo output.InterconsultaRepository
}

func NewInterconsultaService(repo output.InterconsultaRepository) *interconsultaService {
	return &interconsultaService{repo: repo}
}

func (s *interconsultaService) ObtenerPorId(ctx context.Context, id int) (*domain.Interconsulta, error) {
	return s.repo.ObtenerPorId(ctx, id)
}

func (s *interconsultaService) ListarPorServicio(ctx context.Context, tipoServicio string) ([]domain.Interconsulta, error) {
	return s.repo.ListarPorServicio(ctx, tipoServicio)
}

func (s *interconsultaService) ListarPorAtencion(ctx context.Context, idAtencion int) ([]domain.Interconsulta, error) {
	return s.repo.ListarPorAtencion(ctx, idAtencion)
}

func (s *interconsultaService) Crear(ctx context.Context, interconsulta domain.Interconsulta) error {
	return s.repo.Guardar(ctx, interconsulta)
}

func (s *interconsultaService) ActualizarEstado(ctx context.Context, id int, estado string) error {
	return s.repo.ActualizarEstado(ctx, id, estado)
}

func (s *interconsultaService) GuardarFirma(ctx context.Context, firma domain.FirmaInterconsulta) error {
	return s.repo.GuardarFirma(ctx, firma)
}

func (s *interconsultaService) ListarEspecialidades(ctx context.Context) ([]domain.EspecialidadInterconsulta, error) {
	return s.repo.ListarEspecialidades(ctx)
}

func (s *interconsultaService) ListarMedicosPorEspecialidad(ctx context.Context, IdEspecialidad int) ([]domain.MedicoInterconsulta, error) {
	return s.repo.ListarMedicosPorEspecialidad(ctx, IdEspecialidad)
}
