package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type MedicoListaEsperaRepository interface {
	Listar(ctx context.Context) ([]domain.MedicoListaEspera, error)
}
