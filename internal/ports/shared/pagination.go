// Package shared contiene tipos reutilizables entre puertos de entrada y
// salida (paginación) que no pertenecen a ningún agregado de dominio en
// particular.
package shared

// PageRequest normaliza los parámetros de paginación recibidos desde un
// adaptador de entrada.
type PageRequest struct {
	Page     int
	PageSize int
}

const maxPageSize = 100

// NewPageRequest aplica valores por defecto y límites razonables.
func NewPageRequest(page, pageSize int) PageRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > maxPageSize {
		pageSize = 20
	}
	return PageRequest{Page: page, PageSize: pageSize}
}

// Offset calcula el desplazamiento SQL (OFFSET) para esta página.
func (p PageRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// PageResponse es el sobre estándar de una respuesta paginada hacia Angular.
type PageResponse[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// PatientSearchParams agrupa los filtros de búsqueda de pacientes que
// alimentan el SP usp_go_listarpacientes. Un valor vacío en cualquier campo
// implica que ese filtro no se aplica.
type PatientSearchParams struct {
	DocumentNumber  string
	HistoryNumber   string
	PaternalSurname string
	MaternalSurname string
	Names           string
}

// NewPageResponse construye la respuesta paginada a partir de los items de
// la página actual y el total de registros que reportó el repositorio.
func NewPageResponse[T any](items []T, req PageRequest, totalItems int) PageResponse[T] {
	totalPages := 0
	if req.PageSize > 0 {
		totalPages = (totalItems + req.PageSize - 1) / req.PageSize
	}
	return PageResponse[T]{
		Items:      items,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// TriageListParams agrupa los filtros del SP ListarTriaje_Emergencia.
// Los filtros opcionales se envían como vacíos o -100 para que el SP no
// los aplique, replicando el comportamiento del frontend.
type TriageListParams struct {
	FechaInicio       string
	FechaFin          string
	Filtro            string
	DerivadoAServicio int
	IdEmpleado        int
}

// TriageAdmisionParams agrupa los filtros del SP
// webGestionAtencion_E_H_BusquedaFiltrar, que devuelve los pacientes con
// triaje que aún no han sido admisionados. Los filtros numéricos se envían
// como 0 para que el SP no los aplique.
type TriageAdmisionParams struct {
	Fecha          string
	Filtro         string
	NroCta         int
	IdDepartamento int
	IdEspecialidad int
	IdServicio     int
	IdTipoServicio int
}

// TriageReporteParams agrupa los filtros del SP WebSelectReporteTriaje.
// Los filtros se envían como -100 para que el SP no los aplique.
type TriageReporteParams struct {
	IDTriaje   int
	IDPaciente int
}

// FichaAdmisionParams agrupa los parámetros del SP webFichaEmergencia,
// que retorna los datos del paciente y adicionales para generar la ficha
// de admisión de una cuenta de atención.
type FichaAdmisionParams struct {
	IdCuentaAtencion int64
}

// TriajeConsultaParams agrupa los filtros de la bandeja de triaje de
// consulta externa (SP AtencionesTriajeFiltro). El filtro se construye
// como fragmento WHERE validado en el repositorio; el texto de búsqueda
// se escapa para evitar inyección SQL.
type TriajeConsultaParams struct {
	FechaInicio string
	FechaFin    string
	Filtro      string
	IdServicio  int
}

// TriajeConsultaEstadoParams agrupa los parámetros del SP
// AtencionesTriajeEstado, que actualiza el estado de un triaje de
// consulta externa.
type TriajeConsultaEstadoParams struct {
	IdTriaje int64
	Estado   string
}
