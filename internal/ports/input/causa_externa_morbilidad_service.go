package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// CausaExternaMorbilidadService es el puerto de entrada para consultar el
// catálogo de causas externas de morbilidad.
type CausaExternaMorbilidadService interface {
	Listar(ctx context.Context) ([]domain.CausaExternaMorbilidad, error)
}
