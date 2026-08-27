package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// TriageService es el puerto de entrada para registrar y consultar los
// triajes de los pacientes.
type TriageService interface {
	// CreateTriage registra el triaje invocando el procedimiento almacenado
	// webTab_PacienteTriajeAgregar y retorna el @Resultado del SP.
	CreateTriage(ctx context.Context, triage *domain.Triage) (string, error)

	// ListTriage lista los triajes invocando el SP ListarTriaje_Emergencia
	// con los filtros recibidos.
	ListTriage(ctx context.Context, params shared.TriageListParams) ([]map[string]any, error)

	// ListPendingAdmission lista los pacientes con triaje sin admisión
	// invocando el SP webGestionAtencion_E_H_BusquedaFiltrar.
	ListPendingAdmission(ctx context.Context, params shared.TriageAdmisionParams) ([]map[string]any, error)

	// CreateAdmission admisiona un paciente desde su triaje invocando el
	// SP WebCrearAtencionDesdeTriaje y retorna el @Resultado del SP.
	CreateAdmission(ctx context.Context, admision *domain.AdmisionDesdeTriaje) (string, error)

	// GetReporte genera el reporte de triaje invocando el SP
	// WebSelectReporteTriaje con los filtros por id de triaje e id de
	// paciente (-100 para ignorar el filtro).
	GetReporte(ctx context.Context, params shared.TriageReporteParams) ([]map[string]any, error)

	// GetFichaAdmision genera la ficha de admisión invocando el SP
	// webFichaEmergencia para la cuenta de atención indicada.
	GetFichaAdmision(ctx context.Context, params shared.FichaAdmisionParams) (*map[string]any, error)

	// ListarMedicosPorEspecialidad lista los médicos de una especialidad
	// invocando el SP usp_go_MedicosFiltrarPorIdEspecialidad.
	ListarMedicosPorEspecialidad(ctx context.Context, IdEspecialidad int) ([]domain.MedicoFila, error)

	// ListTriajeConsulta lista la bandeja de triaje de consulta externa
	// invocando el SP AtencionesTriajeFiltro.
	ListTriajeConsulta(ctx context.Context, params shared.TriajeConsultaParams) ([]map[string]any, error)

	// CreateTriajeConsulta registra o actualiza el triaje de consulta
	// externa invocando el SP AtencionesTriajeAgregar. Retorna el
	// @Resultado informado por el SP.
	CreateTriajeConsulta(ctx context.Context, triaje *domain.TriajeConsulta) (string, error)

	// GetTriajeConsultaPorAtencion devuelve el triaje de consulta externa
	// vigente (UltimoTriaje = 1) de una atención, o nil si no existe.
	GetTriajeConsultaPorAtencion(ctx context.Context, idAtencion int64) (*map[string]any, error)

	// UpdateEstadoTriajeConsulta actualiza el estado de un triaje de
	// consulta externa invocando el SP AtencionesTriajeEstado.
	UpdateEstadoTriajeConsulta(ctx context.Context, params shared.TriajeConsultaEstadoParams) error
}
