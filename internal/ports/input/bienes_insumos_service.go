package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type BienesInsumosService interface {
	Listar(ctx context.Context, filtros domain.BienInsumoFiltros) (domain.PaginacionBienesInsumos, error)
	ObtenerPorID(ctx context.Context, idProducto int) (*domain.BienInsumoDetalle, error)
	Crear(ctx context.Context, payload domain.BienInsumoGuardarPayload) (int, error)
	Actualizar(ctx context.Context, idProducto int, payload domain.BienInsumoGuardarPayload) error
	Eliminar(ctx context.Context, idProducto int) error
	ListarGruposFarmacologicos(ctx context.Context) ([]domain.GrupoFarmacologico, error)
	ListarSubGruposFarmacologicos(ctx context.Context, idGrupo int) ([]domain.SubGrupoFarmacologico, error)
	ObtenerAuxiliares(ctx context.Context) (*domain.CatalogoBienesInsumosAuxiliares, error)
	ListarPreciosPorPlan(ctx context.Context, idProducto int) ([]domain.BienInsumoPrecioPlan, error)
	GuardarPreciosPorPlan(ctx context.Context, idProducto int, precios []domain.BienInsumoPrecioPlan) error
}
