package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type bienesInsumosUseCase struct {
	repositorio output.BienesInsumosRepository
}

func NewBienesInsumosUseCase(repositorio output.BienesInsumosRepository) input.BienesInsumosService {
	return &bienesInsumosUseCase{repositorio: repositorio}
}

func (u *bienesInsumosUseCase) Listar(ctx context.Context, filtros domain.BienInsumoFiltros) (domain.PaginacionBienesInsumos, error) {
	if filtros.Page <= 0 {
		filtros.Page = 1
	}
	if filtros.PageSize <= 0 || filtros.PageSize > 100 {
		filtros.PageSize = 20
	}
	filtros.Filtro = strings.TrimSpace(filtros.Filtro)

	return u.repositorio.Listar(ctx, filtros)
}

func (u *bienesInsumosUseCase) ObtenerPorID(ctx context.Context, idProducto int) (*domain.BienInsumoDetalle, error) {
	if idProducto <= 0 {
		return nil, errors.New("identificador de producto no valido")
	}
	return u.repositorio.ObtenerPorID(ctx, idProducto)
}

func (u *bienesInsumosUseCase) Crear(ctx context.Context, payload domain.BienInsumoGuardarPayload) (int, error) {
	if err := validarPayload(&payload); err != nil {
		return 0, err
	}
	return u.repositorio.Crear(ctx, payload)
}

func (u *bienesInsumosUseCase) Actualizar(ctx context.Context, idProducto int, payload domain.BienInsumoGuardarPayload) error {
	if idProducto <= 0 {
		return errors.New("identificador de producto no valido")
	}
	if err := validarPayload(&payload); err != nil {
		return err
	}
	return u.repositorio.Actualizar(ctx, idProducto, payload)
}

func validarPayload(payload *domain.BienInsumoGuardarPayload) error {
	payload.Nombre = strings.TrimSpace(payload.Nombre)
	if payload.Nombre == "" {
		return errors.New("el nombre oficial del bien/insumo es requerido (Pestaña 1)")
	}
	if len(payload.Nombre) > 250 {
		return errors.New("el nombre oficial no puede superar los 250 caracteres")
	}

	if payload.IdCentroCosto <= 0 {
		return errors.New("debe seleccionar un Centro de Costo válido (Pestaña 2)")
	}
	if payload.IdPartida <= 0 {
		return errors.New("debe seleccionar una Partida Presupuestal válida (Pestaña 2)")
	}

	if payload.PrecioCompra < 0 || payload.PrecioCompra > 999999.99 {
		return errors.New("el precio de compra debe ser un valor válido entre 0 y 999,999.99")
	}
	if payload.PrecioDistribucion < 0 || payload.PrecioDistribucion > 999999.99 {
		return errors.New("el precio de distribución debe ser un valor válido entre 0 y 999,999.99")
	}
	if payload.PrecioDonacion < 0 || payload.PrecioDonacion > 999999.99 {
		return errors.New("el precio de donación debe ser un valor válido entre 0 y 999,999.99")
	}

	tienePrecioValido := payload.PrecioCompra > 0 || payload.PrecioDistribucion > 0
	for _, plan := range payload.PreciosPorPlan {
		if plan.PrecioUnitario < 0 || plan.PrecioUnitario > 999999.99 {
			return errors.New("los precios por plan deben ser valores válidos entre 0 y 999,999.99 (Pestaña 3)")
		}
		if plan.PrecioUnitario > 0 {
			tienePrecioValido = true
		}
	}

	if !tienePrecioValido {
		return errors.New("debe ingresar al menos un precio de venta o adquisición mayor a cero (Pestaña 3)")
	}

	return nil
}

func (u *bienesInsumosUseCase) Eliminar(ctx context.Context, idProducto int) error {
	if idProducto <= 0 {
		return errors.New("identificador de producto no valido")
	}
	return u.repositorio.Eliminar(ctx, idProducto)
}

func (u *bienesInsumosUseCase) ListarGruposFarmacologicos(ctx context.Context) ([]domain.GrupoFarmacologico, error) {
	return u.repositorio.ListarGruposFarmacologicos(ctx)
}

func (u *bienesInsumosUseCase) ListarSubGruposFarmacologicos(ctx context.Context, idGrupo int) ([]domain.SubGrupoFarmacologico, error) {
	return u.repositorio.ListarSubGruposFarmacologicos(ctx, idGrupo)
}

func (u *bienesInsumosUseCase) ObtenerAuxiliares(ctx context.Context) (*domain.CatalogoBienesInsumosAuxiliares, error) {
	return u.repositorio.ObtenerAuxiliares(ctx)
}

func (u *bienesInsumosUseCase) ListarPreciosPorPlan(ctx context.Context, idProducto int) ([]domain.BienInsumoPrecioPlan, error) {
	return u.repositorio.ListarPreciosPorPlan(ctx, idProducto)
}

func (u *bienesInsumosUseCase) GuardarPreciosPorPlan(ctx context.Context, idProducto int, precios []domain.BienInsumoPrecioPlan) error {
	if idProducto <= 0 {
		return errors.New("identificador de producto no valido")
	}
	return u.repositorio.GuardarPreciosPorPlan(ctx, idProducto, precios)
}
