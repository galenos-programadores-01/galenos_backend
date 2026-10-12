package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type AuditoriaRepository interface {
	Registrar(ctx context.Context, reg domain.AuditoriaRegistro) error
}
