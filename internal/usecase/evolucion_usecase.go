package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type evolucionUseCase struct {
	repo output.EvolucionRepository
}

func NewEvolucionUseCase(repo output.EvolucionRepository) input.EvolucionService {
	return &evolucionUseCase{
		repo: repo,
	}
}

func (uc *evolucionUseCase) GetPatientTray(ctx context.Context, fini, ffin string, idUsuario int) ([]domain.PatientListItem, error) {
	return uc.repo.ListPatients(ctx, fini, ffin, idUsuario)
}

func (uc *evolucionUseCase) GetEvoluciones(ctx context.Context, idRegAtencion int) ([]domain.EvolucionFirma, error) {
	return uc.repo.ListEvoluciones(ctx, idRegAtencion)
}

func (uc *evolucionUseCase) SaveEvolucion(ctx context.Context, idRegAtencion int, idEmpleadoRegistra int, dataB64 string) error {
	evolution := domain.EvolucionFirma{
		IdRegAtencion:      idRegAtencion,
		NombreDocumento:    "EvolucionMedica",
		DataB64:            dataB64,
		IdEmpleadoRegistra: idEmpleadoRegistra,
	}
	return uc.repo.SaveEvolucion(ctx, evolution)
}

func (uc *evolucionUseCase) GetBandeja(ctx context.Context, fechaInicio, fechaFin, filtro string) ([]domain.EvolucionBandejaItem, error) {
	return uc.repo.ListBandeja(ctx, fechaInicio, fechaFin, filtro)
}

func (uc *evolucionUseCase) InsertEvolucionMedica(ctx context.Context, item domain.EvolucionMedicaInsert) (int, string, error) {
	return uc.repo.InsertEvolucionMedica(ctx, item)
}
