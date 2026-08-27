package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type EvolucionService interface {
	GetPatientTray(ctx context.Context, fini, ffin string, idUsuario int) ([]domain.PatientListItem, error)
	GetEvoluciones(ctx context.Context, idRegAtencion int) ([]domain.EvolucionFirma, error)
	SaveEvolucion(ctx context.Context, idRegAtencion int, idEmpleadoRegistra int, dataB64 string) error
	GetBandeja(ctx context.Context, fechaInicio, fechaFin, filtro string) ([]domain.EvolucionBandejaItem, error)
	InsertEvolucionMedica(ctx context.Context, item domain.EvolucionMedicaInsert) (int, string, error)
}
