package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type AuditoriaService interface {
	Registrar(ctx context.Context, reg domain.AuditoriaRegistro) error
}
