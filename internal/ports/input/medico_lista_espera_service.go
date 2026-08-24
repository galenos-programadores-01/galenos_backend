package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type MedicoListaEsperaService interface {
	Listar(ctx context.Context) ([]domain.MedicoListaEspera, error)
}
