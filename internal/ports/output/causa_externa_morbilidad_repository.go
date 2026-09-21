package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// CausaExternaMorbilidadRepository es el puerto de salida para acceder al
// catálogo de causas externas de morbilidad.
type CausaExternaMorbilidadRepository interface {
	// Listar invoca el SP uso_go_EmergenciaCausaExternaMorbilidadListar y
	// devuelve el catálogo completo de causas externas de morbilidad.
	Listar(ctx context.Context) ([]domain.CausaExternaMorbilidad, error)
}
