package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type EvolucionRepository interface {
	ListPatients(ctx context.Context, fini, ffin string, idUsuario int) ([]domain.PatientListItem, error)
	ListEvoluciones(ctx context.Context, idRegAtencion int) ([]domain.EvolucionFirma, error)
	SaveEvolucion(ctx context.Context, evolucion domain.EvolucionFirma) error
	ListBandeja(ctx context.Context, fechaInicio, fechaFin, filtro string) ([]domain.EvolucionBandejaItem, error)
	InsertEvolucionMedica(ctx context.Context, item domain.EvolucionMedicaInsert) (int, string, error)
}
