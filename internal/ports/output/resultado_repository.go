package output

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ResultadoRepository interface {
	ListarLaboratorioPorPaciente(ctx context.Context, idPaciente int) ([]domain.Resultado, error)
	ListarImagenesPorPaciente(ctx context.Context, idPaciente int) ([]domain.Resultado, error)
	ObtenerDetalleLaboratorio(ctx context.Context, idOrden, idProducto int) ([]domain.DetalleResultadoLab, error)
	ObtenerDetalleImagen(ctx context.Context, idOrden, idProducto int) (*domain.DetalleResultadoImagen, error)
}
