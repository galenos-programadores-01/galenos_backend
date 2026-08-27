package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type resultadoService struct {
	repo output.ResultadoRepository
}

func NewResultadoService(repo output.ResultadoRepository) *resultadoService {
	return &resultadoService{repo: repo}
}

func (s *resultadoService) ListarResultadosLaboratorio(ctx context.Context, idPaciente int) ([]domain.Resultado, error) {
	return s.repo.ListarLaboratorioPorPaciente(ctx, idPaciente)
}

func (s *resultadoService) ListarResultadosImagenes(ctx context.Context, idPaciente int) ([]domain.Resultado, error) {
	return s.repo.ListarImagenesPorPaciente(ctx, idPaciente)
}

func (s *resultadoService) ObtenerDetalleLaboratorio(ctx context.Context, idOrden, idProducto int) ([]domain.DetalleResultadoLab, error) {
	return s.repo.ObtenerDetalleLaboratorio(ctx, idOrden, idProducto)
}

func (s *resultadoService) ObtenerDetalleImagen(ctx context.Context, idOrden, idProducto int) (*domain.DetalleResultadoImagen, error) {
	return s.repo.ObtenerDetalleImagen(ctx, idOrden, idProducto)
}
