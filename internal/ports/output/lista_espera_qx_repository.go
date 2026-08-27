package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ListaEsperaQxRepository interface {
	Listar(ctx context.Context, fecha string, paciente string) ([]domain.ListaEsperaQx, error)
	Crear(ctx context.Context, item domain.ListaEsperaQxCrear) error
}
