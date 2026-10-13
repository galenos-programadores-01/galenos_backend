package usecase

import (
	"context"
	"errors"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type atencionCEUseCase struct {
	repo output.AtencionCERepository
}

func NewAtencionCEUseCase(repo output.AtencionCERepository) input.AtencionCEService {
	return &atencionCEUseCase{
		repo: repo,
	}
}

func (uc *atencionCEUseCase) ListarColaPacientes(ctx context.Context, filtros domain.ColaFiltrosCE) ([]domain.ColaPacienteCEItem, error) {
	return uc.repo.ListarColaPacientes(ctx, filtros)
}

func (uc *atencionCEUseCase) ObtenerFichaAtencion(ctx context.Context, idAtencion int) (*domain.FichaAtencionCEDetalle, error) {
	if idAtencion <= 0 {
		return nil, errors.New("identificador de atencion invalido")
	}
	return uc.repo.ObtenerFichaAtencion(ctx, idAtencion)
}

func (uc *atencionCEUseCase) GuardarFichaAtencion(ctx context.Context, payload domain.AtencionCEGuardarPayload) error {
	if payload.IdAtencion <= 0 {
		return errors.New("identificador de atencion requerido para guardar la ficha")
	}
	return uc.repo.GuardarFichaAtencion(ctx, payload)
}

func (uc *atencionCEUseCase) ListarDestinosAtencion(ctx context.Context) ([]domain.DestinoAtencionItem, error) {
	return uc.repo.ListarDestinosAtencion(ctx)
}

func (uc *atencionCEUseCase) ListarCatalogoLabsHIS(ctx context.Context) ([]domain.CatalogoLabHISItem, error) {
	return uc.repo.ListarCatalogoLabsHIS(ctx)
}

func (uc *atencionCEUseCase) ListarHistorialAtenciones(ctx context.Context, idPaciente int) ([]domain.HistorialAtencionItem, error) {
	if idPaciente <= 0 {
		return nil, errors.New("identificador de paciente invalido")
	}
	return uc.repo.ListarHistorialAtenciones(ctx, idPaciente)
}

func (uc *atencionCEUseCase) ObtenerHistorialAtencionDetalle(ctx context.Context, idAtencion int) (*domain.HistorialAtencionDetalle, error) {
	if idAtencion <= 0 {
		return nil, errors.New("identificador de atencion invalido")
	}
	return uc.repo.ObtenerHistorialAtencionDetalle(ctx, idAtencion)
}

func (uc *atencionCEUseCase) ListarReportesOperatorios(ctx context.Context, idPaciente int) ([]domain.HistorialReporteOperatorioItem, error) {
	if idPaciente <= 0 {
		return nil, errors.New("identificador de paciente invalido")
	}
	return uc.repo.ListarReportesOperatorios(ctx, idPaciente)
}

func (uc *atencionCEUseCase) ListarProcedimientosHistorial(ctx context.Context, idPaciente int) ([]domain.HistorialProcedimientoItem, error) {
	if idPaciente <= 0 {
		return nil, errors.New("identificador de paciente invalido")
	}
	return uc.repo.ListarProcedimientosHistorial(ctx, idPaciente)
}

func (uc *atencionCEUseCase) ListarServiciosPorMedico(ctx context.Context, idTipoServicio int, idEmpleado int) ([]domain.ServicioMedicoItem, error) {
	if idTipoServicio <= 0 {
		idTipoServicio = 1
	}
	if idEmpleado == 0 {
		idEmpleado = -100
	}
	servicios, err := uc.repo.ListarServiciosPorMedico(ctx, idTipoServicio, idEmpleado)
	if err != nil {
		return nil, err
	}
	if len(servicios) == 0 && idEmpleado != -100 {
		return uc.repo.ListarServiciosPorMedico(ctx, idTipoServicio, -100)
	}
	return servicios, nil
}

func (uc *atencionCEUseCase) ListarHistorialTriajePaciente(ctx context.Context, idPaciente int) ([]domain.PacienteHistorialTriajeItem, error) {
	if idPaciente <= 0 {
		return nil, errors.New("identificador de paciente invalido")
	}
	return uc.repo.ListarHistorialTriajePaciente(ctx, idPaciente)
}

