package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type hospitalizacionUseCase struct {
	repo output.HospitalizacionRepository
}

// NewHospitalizacionUseCase construye el caso de uso de admisión de
// hospitalización.
func NewHospitalizacionUseCase(repo output.HospitalizacionRepository) input.HospitalizacionService {
	return &hospitalizacionUseCase{repo: repo}
}

// Admisionar delega la persistencia en el repositorio (SP
// usp_go_AdmisionHospitalizacion) y devuelve el resultado informado por el
// procedimiento almacenado.
func (uc *hospitalizacionUseCase) Admisionar(ctx context.Context, admision *domain.AdmisionHospitalizacion) (string, error) {
	result, err := uc.repo.Admisionar(ctx, admision)
	if err != nil {
		return "", fmt.Errorf("admitting patient to hospitalization: %w", err)
	}
	return result, nil
}

// ObtenerReporteTicket delega la consulta en el repositorio (SP
// usp_go_GenererReporteTicketHospi). Retorna nil si la cuenta no existe.
func (uc *hospitalizacionUseCase) ObtenerReporteTicket(ctx context.Context, idCuenta int64) (*domain.ReporteTicketHospitalizacion, error) {
	reporte, err := uc.repo.ObtenerReporteTicket(ctx, idCuenta)
	if err != nil {
		return nil, fmt.Errorf("getting hospitalization ticket report: %w", err)
	}
	return reporte, nil
}
