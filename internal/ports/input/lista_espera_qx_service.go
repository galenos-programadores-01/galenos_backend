package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ListaEsperaQxService interface {
	Listar(ctx context.Context, fecha string, paciente string) ([]domain.ListaEsperaQx, error)
	Crear(ctx context.Context, item domain.ListaEsperaQxCrear) error
}
