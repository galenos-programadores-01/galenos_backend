package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type ordenService struct {
	repo output.OrdenRepository
}

func NewOrdenService(repo output.OrdenRepository) *ordenService {
	return &ordenService{repo: repo}
}

func (s *ordenService) ListarPorCuenta(ctx context.Context, idRegAtencion int) ([]domain.OrdenMedica, error) {
	return s.repo.ListarPorCuenta(ctx, idRegAtencion)
}

func (s *ordenService) CrearOrden(ctx context.Context, orden domain.OrdenMedica, detalles []domain.DetalleOrden, idEmpleado int) error {
	return s.repo.CrearOrden(ctx, orden, detalles, idEmpleado)
}

func (s *ordenService) BuscarProductos(ctx context.Context, filtro string, limite int) ([]domain.ProductoCatalogo, error) {
	return s.repo.BuscarProductos(ctx, filtro, limite)
}

func (s *ordenService) ObtenerFormatoImpresionReceta(ctx context.Context, identificadorReceta int, identificadorCuentaAtencion int, identificadorAtencion int) (*domain.FormatoImpresionReceta, error) {
	return s.repo.ObtenerFormatoImpresionReceta(ctx, identificadorReceta, identificadorCuentaAtencion, identificadorAtencion)
}

