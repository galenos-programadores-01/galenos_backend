package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type triageUseCase struct {
	repo output.TriageRepository
}

// NewTriageUseCase construye el caso de uso de registro de triaje.
func NewTriageUseCase(repo output.TriageRepository) input.TriageService {
	return &triageUseCase{repo: repo}
}

// CreateTriage delega la persistencia en el repositorio y devuelve el
// resultado informado por el procedimiento almacenado.
func (uc *triageUseCase) CreateTriage(ctx context.Context, triage *domain.Triage) (string, error) {
	result, err := uc.repo.Create(ctx, triage)
	if err != nil {
		return "", fmt.Errorf("registering triage: %w", err)
	}
	return result, nil
}

// ListTriage delega el listado (SP ListarTriaje_Emergencia) en el
// repositorio y devuelve los registros crudos.
func (uc *triageUseCase) ListTriage(ctx context.Context, params shared.TriageListParams) ([]map[string]any, error) {
	items, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing triages: %w", err)
	}
	return items, nil
}

// ListPendingAdmission delega en el repositorio (SP
// webGestionAtencion_E_H_BusquedaFiltrar) y devuelve los registros crudos.
func (uc *triageUseCase) ListPendingAdmission(ctx context.Context, params shared.TriageAdmisionParams) ([]map[string]any, error) {
	items, err := uc.repo.ListPendingAdmission(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing triages pending admission: %w", err)
	}
	return items, nil
}

// CreateAdmission delega la admisión (SP WebCrearAtencionDesdeTriaje) en
// el repositorio y devuelve el resultado informado por el SP.
func (uc *triageUseCase) CreateAdmission(ctx context.Context, admision *domain.AdmisionDesdeTriaje) (string, error) {
	result, err := uc.repo.CreateAdmission(ctx, admision)
	if err != nil {
		return "", fmt.Errorf("creating admission from triage: %w", err)
	}
	return result, nil
}

// GetReporte delega en el repositorio (SP WebSelectReporteTriaje) y
// devuelve los registros crudos del reporte.
func (uc *triageUseCase) GetReporte(ctx context.Context, params shared.TriageReporteParams) ([]map[string]any, error) {
	items, err := uc.repo.GetReporte(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("getting triage report: %w", err)
	}
	return items, nil
}

// GetFichaAdmision delega en el repositorio (SP webFichaEmergencia) y
// devuelve los datos crudos de la ficha.
func (uc *triageUseCase) GetFichaAdmision(ctx context.Context, params shared.FichaAdmisionParams) (*map[string]any, error) {
	item, err := uc.repo.GetFichaAdmision(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("getting admission record: %w", err)
	}
	return item, nil
}

// ListarMedicosPorEspecialidad delega en el repositorio (SP
// usp_go_MedicosFiltrarPorIdEspecialidad) y devuelve los médicos de la
// especialidad.
func (uc *triageUseCase) ListarMedicosPorEspecialidad(ctx context.Context, IdEspecialidad int) ([]domain.MedicoFila, error) {
	items, err := uc.repo.ListarMedicosPorEspecialidad(ctx, IdEspecialidad)
	if err != nil {
		return nil, fmt.Errorf("listing doctors by specialty: %w", err)
	}
	return items, nil
}

// ListTriajeConsulta delega en el repositorio (SP AtencionesTriajeFiltro)
// y devuelve la bandeja de triaje de consulta externa.
func (uc *triageUseCase) ListTriajeConsulta(ctx context.Context, params shared.TriajeConsultaParams) ([]map[string]any, error) {
	items, err := uc.repo.ListTriajeConsulta(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing outpatient triage: %w", err)
	}
	return items, nil
}

// CreateTriajeConsulta delega en el repositorio (SP
// AtencionesTriajeAgregar) y devuelve el @Resultado informado por el SP.
func (uc *triageUseCase) CreateTriajeConsulta(ctx context.Context, triaje *domain.TriajeConsulta) (string, error) {
	result, err := uc.repo.CreateTriajeConsulta(ctx, triaje)
	if err != nil {
		return "", fmt.Errorf("registering outpatient triage: %w", err)
	}
	return result, nil
}

// GetTriajeConsultaPorAtencion delega en el repositorio y devuelve el
// triaje de consulta externa vigente de la atención.
func (uc *triageUseCase) GetTriajeConsultaPorAtencion(ctx context.Context, idAtencion int64) (*map[string]any, error) {
	item, err := uc.repo.GetTriajeConsultaPorAtencion(ctx, idAtencion)
	if err != nil {
		return nil, fmt.Errorf("getting outpatient triage by attention: %w", err)
	}
	return item, nil
}

// UpdateEstadoTriajeConsulta delega en el repositorio (SP
// AtencionesTriajeEstado) para actualizar el estado del triaje.
func (uc *triageUseCase) UpdateEstadoTriajeConsulta(ctx context.Context, params shared.TriajeConsultaEstadoParams) error {
	if err := uc.repo.UpdateEstadoTriajeConsulta(ctx, params); err != nil {
		return fmt.Errorf("updating outpatient triage state: %w", err)
	}
	return nil
}
