package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type SintomaService interface {
	ListarCatalogo(ctx context.Context) ([]domain.SintomaCatalogo, error)
	AgregarCatalogo(ctx context.Context, sistema, sintoma string, idUsuario int) error
	GuardarEvolucionSintomas(ctx context.Context, idRegAtencion int, sintomas []domain.SintomaSeleccionado, idUsuario int) error
	ObtenerAtencionSintomas(ctx context.Context, idRegAtencion int) ([]domain.SintomaSeleccionado, error)
}
