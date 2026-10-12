package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type serviciosCatalogoUseCase struct {
	repositorio output.ServiciosCatalogoRepository
}

func NewServiciosCatalogoUseCase(repositorio output.ServiciosCatalogoRepository) input.ServiciosCatalogoService {
	return &serviciosCatalogoUseCase{repositorio: repositorio}
}

func (u *serviciosCatalogoUseCase) Listar(ctx context.Context, filtros domain.ServicioCatalogoFiltros) (domain.PaginacionServiciosCatalogo, error) {
	if filtros.Page <= 0 {
		filtros.Page = 1
	}
	if filtros.PageSize <= 0 || filtros.PageSize > 100 {
		filtros.PageSize = 20
	}
	filtros.Filtro = strings.TrimSpace(filtros.Filtro)

	return u.repositorio.Listar(ctx, filtros)
}

func (u *serviciosCatalogoUseCase) ObtenerPorID(ctx context.Context, idProducto int) (*domain.ServicioCatalogoDetalle, error) {
	if idProducto <= 0 {
		return nil, errors.New("identificador de servicio o procedimiento no valido")
	}
	return u.repositorio.ObtenerPorID(ctx, idProducto)
}

func (u *serviciosCatalogoUseCase) Crear(ctx context.Context, payload domain.ServicioCatalogoGuardarPayload) (int, error) {
	if err := validarServicioPayload(&payload); err != nil {
		return 0, err
	}
	return u.repositorio.Crear(ctx, payload)
}

func (u *serviciosCatalogoUseCase) Actualizar(ctx context.Context, idProducto int, payload domain.ServicioCatalogoGuardarPayload) error {
	if idProducto <= 0 {
		return errors.New("identificador de servicio o procedimiento no valido")
	}
	if err := validarServicioPayload(&payload); err != nil {
		return err
	}
	return u.repositorio.Actualizar(ctx, idProducto, payload)
}

func validarServicioPayload(payload *domain.ServicioCatalogoGuardarPayload) error {
	payload.Nombre = strings.TrimSpace(payload.Nombre)
	if payload.Nombre == "" {
		return errors.New("la descripcion del servicio o procedimiento es requerida (Pestaña 1)")
	}
	if len(payload.Nombre) < 3 {
		return errors.New("la descripcion debe contener al menos 3 caracteres (Pestaña 1)")
	}
	if len(payload.Nombre) > 250 {
		return errors.New("la descripcion no puede exceder los 250 caracteres (Pestaña 1)")
	}

	if payload.IdServicioGrupo <= 0 {
		return errors.New("debe seleccionar un Grupo de Servicio válido (Pestaña 1)")
	}

	if payload.PrecioParticular < 0 || payload.PrecioParticular > 999999.99 {
		return errors.New("la tarifa particular debe ser un valor válido entre 0 y 999,999.99 (Pestaña 3)")
	}
	if payload.PrecioSis < 0 || payload.PrecioSis > 999999.99 {
		return errors.New("la tarifa SIS debe ser un valor válido entre 0 y 999,999.99 (Pestaña 3)")
	}

	return nil
}

func (u *serviciosCatalogoUseCase) Eliminar(ctx context.Context, idProducto int) error {
	if idProducto <= 0 {
		return errors.New("identificador de servicio o procedimiento no valido")
	}
	return u.repositorio.Eliminar(ctx, idProducto)
}

func (u *serviciosCatalogoUseCase) ListarGrupos(ctx context.Context) ([]domain.ServicioGrupo, error) {
	return u.repositorio.ListarGrupos(ctx)
}
