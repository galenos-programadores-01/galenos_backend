package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type sisUseCase struct {
	client output.SisClient
	repo   output.SisRepository
}

// NewSisUseCase construye el caso de uso de consulta al SIS.
func NewSisUseCase(client output.SisClient, repo output.SisRepository) input.SisService {
	return &sisUseCase{client: client, repo: repo}
}

// GestionarAfiliacion delega la persistencia en el repositorio.
func (uc *sisUseCase) GestionarAfiliacion(ctx context.Context, afiliacion *domain.SisAfiliacion) error {
	if err := uc.repo.GestionarAfiliacion(ctx, afiliacion); err != nil {
		return fmt.Errorf("managing sis afiliacion: %w", err)
	}
	return nil
}

// ForzarGuardadoFua delega el guardado del FUA en el repositorio.
func (uc *sisUseCase) ForzarGuardadoFua(ctx context.Context, idCuentaAtencion int64) error {
	if err := uc.repo.ForzarGuardadoFua(ctx, idCuentaAtencion); err != nil {
		return fmt.Errorf("forcing fua save: %w", err)
	}
	return nil
}

// AgregarFua delega el agregado del FUA en el repositorio y devuelve el
// resultado informado por el procedimiento almacenado.
func (uc *sisUseCase) AgregarFua(ctx context.Context, idCuentaAtencion, idEmpleado int64, nombrePc string) (string, error) {
	result, err := uc.repo.AgregarFua(ctx, idCuentaAtencion, idEmpleado, nombrePc)
	if err != nil {
		return "", fmt.Errorf("adding fua: %w", err)
	}
	return result, nil
}

// GetFuaImprimir delega la consulta del FUA para imprimir en el repositorio.
func (uc *sisUseCase) GetFuaImprimir(ctx context.Context, idCuentaAtencion int64) (*map[string]any, error) {
	data, err := uc.repo.GetFuaImprimir(ctx, idCuentaAtencion)
	if err != nil {
		return nil, fmt.Errorf("getting fua print data: %w", err)
	}
	return data, nil
}

// ListDiagnosticos delega la consulta de diagnósticos de la atención.
func (uc *sisUseCase) ListDiagnosticos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error) {
	data, err := uc.repo.ListDiagnosticos(ctx, idCuentaAtencion)
	if err != nil {
		return nil, fmt.Errorf("listing diagnostics: %w", err)
	}
	return data, nil
}

// ListMedicamentos delega la consulta de medicamentos de la cuenta de atención.
func (uc *sisUseCase) ListMedicamentos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error) {
	data, err := uc.repo.ListMedicamentos(ctx, idCuentaAtencion)
	if err != nil {
		return nil, fmt.Errorf("listing medications: %w", err)
	}
	return data, nil
}

// ListProcedimientos delega la consulta de procedimientos de la cuenta de atención.
func (uc *sisUseCase) ListProcedimientos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error) {
	data, err := uc.repo.ListProcedimientos(ctx, idCuentaAtencion)
	if err != nil {
		return nil, fmt.Errorf("listing procedures: %w", err)
	}
	return data, nil
}

// ListConsumo delega la consulta del detalle de consumo de la cuenta de atención.
func (uc *sisUseCase) ListConsumo(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error) {
	data, err := uc.repo.ListConsumo(ctx, idCuentaAtencion)
	if err != nil {
		return nil, fmt.Errorf("listing consumption: %w", err)
	}
	return data, nil
}

func (uc *sisUseCase) ConsultarAfiliado(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error) {
	if params.DocumentNumber == "" && (params.Disa == "" || params.TipoFormato == "" || params.NroContrato == "") {
		return domain.SisAfiliado{}, domain.ErrInvalidDocumentNumber
	}
	if params.TipoDocumento != 1 && params.TipoDocumento != 3 {
		return domain.SisAfiliado{}, domain.ErrInvalidDocumentType
	}
	if params.Opcion <= 0 {
		params.Opcion = 1
	}

	result, err := uc.client.ConsultarAfiliado(ctx, params)
	if err != nil {
		return domain.SisAfiliado{}, fmt.Errorf("consulting sis: %w", err)
	}

	return result, nil
}

func (uc *sisUseCase) BuscarPorAfiliacion(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error) {
	if params.Disa == "" || params.Lote == "" || params.NroContrato == "" {
		return domain.SisAfiliado{}, domain.ErrInvalidDocumentNumber
	}

	result, err := uc.client.BuscarPorAfiliacion(ctx, params)
	if err != nil {
		return domain.SisAfiliado{}, fmt.Errorf("searching sis by affiliation: %w", err)
	}

	return result, nil
}
