package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type farmaciaDisponibilidadUseCase struct {
	repositorio output.FarmaciaDisponibilidadRepository
}

func NewFarmaciaDisponibilidadUseCase(repositorio output.FarmaciaDisponibilidadRepository) input.FarmaciaDisponibilidadService {
	return &farmaciaDisponibilidadUseCase{repositorio: repositorio}
}

func (u *farmaciaDisponibilidadUseCase) ListarAlmacenes(ctx context.Context) ([]domain.AlmacenFarmaciaItem, error) {
	return u.repositorio.ListarAlmacenes(ctx)
}

func (u *farmaciaDisponibilidadUseCase) ObtenerKpis(ctx context.Context, idAlmacen int, canastaTipo string) (domain.DisponibilidadKpis, error) {
	canastaTipo = strings.TrimSpace(canastaTipo)
	if canastaTipo == "" {
		canastaTipo = "TODOS"
	}
	return u.repositorio.ObtenerKpis(ctx, idAlmacen, canastaTipo)
}

func (u *farmaciaDisponibilidadUseCase) ListarDisponibilidad(ctx context.Context, filtros domain.FarmaciaDisponibilidadFiltros) (domain.PaginacionDisponibilidad, error) {
	if filtros.Page <= 0 {
		filtros.Page = 1
	}
	if filtros.PageSize <= 0 || filtros.PageSize > 100 {
		filtros.PageSize = 20
	}
	filtros.Filtro = strings.TrimSpace(filtros.Filtro)

	return u.repositorio.ListarDisponibilidad(ctx, filtros)
}

func (u *farmaciaDisponibilidadUseCase) ListarLotesPorProducto(ctx context.Context, idProducto int, idAlmacen int) ([]domain.LoteItem, error) {
	if idProducto <= 0 {
		return nil, errors.New("identificador de producto no valido")
	}
	return u.repositorio.ListarLotesPorProducto(ctx, idProducto, idAlmacen)
}

func (u *farmaciaDisponibilidadUseCase) ListarLotesPorVencer(ctx context.Context, idAlmacen int, diasMax int) ([]domain.LoteVencimientoDetalleItem, error) {
	if diasMax <= 0 {
		diasMax = 60
	}
	return u.repositorio.ListarLotesPorVencer(ctx, idAlmacen, diasMax)
}

func (u *farmaciaDisponibilidadUseCase) ConsultarKardex(ctx context.Context, idProducto int, idAlmacen int, fechaInicio, fechaFin time.Time) ([]domain.KardexMovimientoItem, error) {
	if idProducto <= 0 {
		return nil, errors.New("identificador de producto no valido")
	}
	if fechaInicio.IsZero() {
		fechaInicio = time.Now().AddDate(0, -1, 0)
	}
	if fechaFin.IsZero() {
		fechaFin = time.Now()
	}
	return u.repositorio.ConsultarKardex(ctx, idProducto, idAlmacen, fechaInicio, fechaFin)
}

func (u *farmaciaDisponibilidadUseCase) ObtenerTendenciasSalidas(ctx context.Context, idAlmacen int) (domain.TendenciasSalidasFarmacia, error) {
	return u.repositorio.ObtenerTendenciasSalidas(ctx, idAlmacen)
}

func (u *farmaciaDisponibilidadUseCase) ObtenerEstadisticasDashboard(ctx context.Context, idAlmacen int, canastaTipo string, fechaInicio, fechaFin string) (domain.EstadisticasFarmaciaDashboard, error) {
	canastaTipo = strings.TrimSpace(canastaTipo)
	if canastaTipo == "" {
		canastaTipo = "TODOS"
	}
	return u.repositorio.ObtenerEstadisticasDashboard(ctx, idAlmacen, canastaTipo, fechaInicio, fechaFin)
}
