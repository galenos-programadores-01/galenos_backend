package output

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type OrdenRepository interface {
	ListarPorCuenta(ctx context.Context, idRegAtencion int) ([]domain.OrdenMedica, error)
	CrearOrden(ctx context.Context, orden domain.OrdenMedica, detalles []domain.DetalleOrden, idEmpleado int) error
	BuscarProductos(ctx context.Context, filtro string, limite int) ([]domain.ProductoCatalogo, error)
	ObtenerFormatoImpresionReceta(ctx context.Context, identificadorReceta int, identificadorCuentaAtencion int, identificadorAtencion int) (*domain.FormatoImpresionReceta, error)
}

