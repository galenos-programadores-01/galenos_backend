package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// TriageRepository es el puerto de salida para persistir y listar los
// registros de triaje.
type TriageRepository interface {
	// Create persiste un triaje invocando el procedimiento almacenado
	// webTab_PacienteTriajeAgregar. Retorna el valor del parámetro de
	// salida @Resultado, que el SP llena con el código/mensaje del
	// resultado de la operación.
	Create(ctx context.Context, triage *domain.Triage) (string, error)

	// List invoca el SP ListarTriaje_Emergencia con los filtros de rango
	// de fechas, texto libre, derivado a servicio y estado. Devuelve los
	// registros como mapas columna -> valor, respetando los nombres que
	// el SP devuelve en runtime.
	List(ctx context.Context, params shared.TriageListParams) ([]map[string]any, error)

	// ListPendingAdmission invoca el SP webGestionAtencion_E_H_BusquedaFiltrar,
	// que devuelve los pacientes con triaje que aún no han sido admisionados.
	ListPendingAdmission(ctx context.Context, params shared.TriageAdmisionParams) ([]map[string]any, error)

	// CreateAdmission crea la atención (admisión) de un paciente desde su
	// triaje invocando el SP WebCrearAtencionDesdeTriaje. Retorna el
	// valor del parámetro de salida @Resultado.
	CreateAdmission(ctx context.Context, admision *domain.AdmisionDesdeTriaje) (string, error)

	// GetReporte invoca el SP WebSelectReporteTriaje con los filtros por
	// id de triaje e id de paciente. Devuelve los registros como mapas
	// columna -> valor.
	GetReporte(ctx context.Context, params shared.TriageReporteParams) ([]map[string]any, error)

	// GetFichaAdmision invoca el SP webFichaEmergencia y devuelve los
	// datos del paciente y adicionales para generar la ficha de admisión
	// de la cuenta de atención indicada.
	GetFichaAdmision(ctx context.Context, params shared.FichaAdmisionParams) (*map[string]any, error)

	// ListarMedicosPorEspecialidad invoca el SP
	// usp_go_MedicosFiltrarPorIdEspecialidad y devuelve los médicos de la
	// especialidad indicada.
	ListarMedicosPorEspecialidad(ctx context.Context, IdEspecialidad int) ([]domain.MedicoFila, error)

	// ListTriajeConsulta invoca el SP AtencionesTriajeFiltro con el
	// fragmento WHERE construido a partir de los filtros validados.
	// Devuelve los registros como mapas columna -> valor.
	ListTriajeConsulta(ctx context.Context, params shared.TriajeConsultaParams) ([]map[string]any, error)

	// CreateTriajeConsulta invoca el SP AtencionesTriajeAgregar, que
	// registra un triaje nuevo o actualiza el vigente de la atención.
	// Retorna el valor del parámetro de salida @Resultado.
	CreateTriajeConsulta(ctx context.Context, triaje *domain.TriajeConsulta) (string, error)

	// GetTriajeConsultaPorAtencion consulta el triaje de consulta externa
	// vigente (UltimoTriaje = 1) de una atención. Devuelve nil si no hay.
	GetTriajeConsultaPorAtencion(ctx context.Context, idAtencion int64) (*map[string]any, error)

	// UpdateEstadoTriajeConsulta invoca el SP AtencionesTriajeEstado para
	// actualizar el estado de un triaje de consulta externa.
	UpdateEstadoTriajeConsulta(ctx context.Context, params shared.TriajeConsultaEstadoParams) error
}
