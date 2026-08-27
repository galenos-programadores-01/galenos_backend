package input

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ResultadoService interface {
	ListarResultadosLaboratorio(ctx context.Context, idPaciente int) ([]domain.Resultado, error)
	ListarResultadosImagenes(ctx context.Context, idPaciente int) ([]domain.Resultado, error)
	ObtenerDetalleLaboratorio(ctx context.Context, idOrden, idProducto int) ([]domain.DetalleResultadoLab, error)
	ObtenerDetalleImagen(ctx context.Context, idOrden, idProducto int) (*domain.DetalleResultadoImagen, error)
}
