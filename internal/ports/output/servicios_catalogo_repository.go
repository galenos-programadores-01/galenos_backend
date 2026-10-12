package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ServiciosCatalogoRepository interface {
	Listar(ctx context.Context, filtros domain.ServicioCatalogoFiltros) (domain.PaginacionServiciosCatalogo, error)
	ObtenerPorID(ctx context.Context, idProducto int) (*domain.ServicioCatalogoDetalle, error)
	Crear(ctx context.Context, payload domain.ServicioCatalogoGuardarPayload) (int, error)
	Actualizar(ctx context.Context, idProducto int, payload domain.ServicioCatalogoGuardarPayload) error
	Eliminar(ctx context.Context, idProducto int) error
	ListarGrupos(ctx context.Context) ([]domain.ServicioGrupo, error)
}
