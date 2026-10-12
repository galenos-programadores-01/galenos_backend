package output

import (
	"context"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type FarmaciaDisponibilidadRepository interface {
	ListarAlmacenes(ctx context.Context) ([]domain.AlmacenFarmaciaItem, error)
	ObtenerKpis(ctx context.Context, idAlmacen int, canastaTipo string) (domain.DisponibilidadKpis, error)
	ListarDisponibilidad(ctx context.Context, filtros domain.FarmaciaDisponibilidadFiltros) (domain.PaginacionDisponibilidad, error)
	ListarLotesPorProducto(ctx context.Context, idProducto int, idAlmacen int) ([]domain.LoteItem, error)
	ListarLotesPorVencer(ctx context.Context, idAlmacen int, diasMax int) ([]domain.LoteVencimientoDetalleItem, error)
	ConsultarKardex(ctx context.Context, idProducto int, idAlmacen int, fechaInicio, fechaFin time.Time) ([]domain.KardexMovimientoItem, error)

	ObtenerTendenciasSalidas(ctx context.Context, idAlmacen int) (domain.TendenciasSalidasFarmacia, error)
	ObtenerEstadisticasDashboard(ctx context.Context, idAlmacen int, canastaTipo string, fechaInicio, fechaFin string) (domain.EstadisticasFarmaciaDashboard, error)
}
