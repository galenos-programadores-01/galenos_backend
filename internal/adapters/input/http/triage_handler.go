package httpadapter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// TriageHandler expone el puerto de entrada input.TriageService.
type TriageHandler struct {
	service input.TriageService
}

// NewTriageHandler inyecta el caso de uso de triaje en el adaptador HTTP.
func NewTriageHandler(service input.TriageService) *TriageHandler {
	return &TriageHandler{service: service}
}

// Create registra un nuevo triaje de emergencia.
//
//	@Summary	 Registrar triaje
//	@Description Registra un nuevo triaje de emergencia con signos vitales y clasificación.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		body body	 createTriajeRequest	true	"Datos del triaje"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje [post]
func (h *TriageHandler) Create(c *gin.Context) {
	var req createTriajeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	domainObj := req.toDomain()
	if idEmpleado := c.GetInt("idEmpleado"); idEmpleado != 0 {
		empID := int64(idEmpleado)
		domainObj.EmployeeID = &empID
	}

	result, err := h.service.CreateTriage(c.Request.Context(), domainObj)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_REGISTER_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}

// List retorna la lista de triajes filtrada por fechas.
//
//	@Summary	 Listar triajes
//	@Description Lista triajes de emergencia con filtros de fecha y estado.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		fini			query		string	true	"Fecha inicio (YYYY-MM-DD)"
//	@Param		ffin			query		string	true	"Fecha fin (YYYY-MM-DD)"
//	@Param		filtro			query		string	false	"Filtro por paciente o documento"
//	@Param		derivadoAServicio	query		int		false	"ID de servicio derivado"
//	@Param		idEmpleado		query		int		false	"ID del empleado logueado"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje [get]
func (h *TriageHandler) List(c *gin.Context) {
	idEmpleado := c.GetInt("idEmpleado")
	if idEmpleado == 0 {
		if raw := c.Query("idEmpleado"); raw != "" {
			if v, err := strconv.ParseInt(raw, 10, 64); err == nil {
				idEmpleado = int(v)
			}
		}
	}

	params := shared.TriageListParams{
		FechaInicio:       c.Query("fini"),
		FechaFin:          c.Query("ffin"),
		Filtro:            c.Query("filtro"),
		DerivadoAServicio: -100,
		IdEmpleado:        idEmpleado,
	}

	if params.FechaInicio == "" || params.FechaFin == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fini y ffin son obligatorios (YYYY-MM-DD)")
		return
	}

	if raw := c.Query("derivadoAServicio"); raw != "" && raw != "-100" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "derivadoAServicio debe ser un entero")
			return
		}
		params.DerivadoAServicio = int(v)
	}

	items, err := h.service.ListTriage(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// ListPendingAdmission lista pacientes pendientes de admisión desde triaje.
//
//	@Summary	 Pendientes de admisión
//	@Description Lista pacientes en triaje pendientes de admisión hospitalaria.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		fecha			query		string	true	"Fecha de atención (YYYY-MM-DD)"
//	@Param		filtro			query		string	false	"Filtro por paciente o documento"
//	@Param		nroCta			query		int		false	"Número de cuenta"
//	@Param		idDepartamento	query		int		false	"ID de departamento"
//	@Param		IdEspecialidad	query		int		false	"ID de especialidad"
//	@Param		idServicio		query		int		false	"ID de servicio"
//	@Param		idTipoServicio	query		int		false	"ID de tipo de servicio"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/pendientes-admision [get]
func (h *TriageHandler) ListPendingAdmission(c *gin.Context) {
	params := shared.TriageAdmisionParams{
		Fecha:  c.Query("fecha"),
		Filtro: c.Query("filtro"),
	}

	if params.Fecha == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fecha es obligatoria (YYYY-MM-DD)")
		return
	}

	intFields := []struct {
		raw   string
		dest  *int
		label string
	}{
		{c.Query("nroCta"), &params.NroCta, "nroCta"},
		{c.Query("idDepartamento"), &params.IdDepartamento, "idDepartamento"},
		{c.Query("IdEspecialidad"), &params.IdEspecialidad, "IdEspecialidad"},
		{c.Query("idServicio"), &params.IdServicio, "idServicio"},
		{c.Query("idTipoServicio"), &params.IdTipoServicio, "idTipoServicio"},
	}
	for _, f := range intFields {
		if f.raw == "" {
			continue
		}
		v, err := strconv.ParseInt(f.raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", f.label+" debe ser un entero")
			return
		}
		*f.dest = int(v)
	}

	items, err := h.service.ListPendingAdmission(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_PENDING_ADMISSION_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// CreateAdmission crea una admisión hospitalaria desde un triaje.
//
//	@Summary	 Crear admisión desde triaje
//	@Description Registra la admisión de un paciente derivado desde triaje de emergencia.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		body body	 createAdmissionFromTriageRequest	true	"Datos de la admisión"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/admision [post]
func (h *TriageHandler) CreateAdmission(c *gin.Context) {
	var req createAdmissionFromTriageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	domainObj := req.toDomain()
	if idEmpleado := c.GetInt("idEmpleado"); idEmpleado != 0 {
		empID := int64(idEmpleado)
		domainObj.IDEmpleado = &empID
	}

	result, err := h.service.CreateAdmission(c.Request.Context(), domainObj)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ADMISSION_FROM_TRIAGE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}

// GetReporte retorna el reporte de un triaje.
//
//	@Summary	 Reporte de triaje
//	@Description Obtiene los datos del reporte de un triaje específico.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		id			query		int	false	"ID del triaje"
//	@Param		idPaciente	query		int	false	"ID del paciente"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/reporte [get]
func (h *TriageHandler) GetReporte(c *gin.Context) {
	params := shared.TriageReporteParams{
		IDTriaje:   -100,
		IDPaciente: -100,
	}

	if raw := c.Query("id"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id debe ser un entero")
			return
		}
		params.IDTriaje = int(v)
	}

	if raw := c.Query("idPaciente"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPaciente debe ser un entero")
			return
		}
		params.IDPaciente = int(v)
	}

	items, err := h.service.GetReporte(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_REPORT_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// GetFichaAdmision retorna la ficha de admisión de una cuenta de atención.
//
//	@Summary	 Ficha de admisión
//	@Description Obtiene los datos de la ficha de admisión por número de cuenta.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		idCuentaAtencion	query		int	true	"ID de la cuenta de atención"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/ficha-admision [get]
func (h *TriageHandler) GetFichaAdmision(c *gin.Context) {
	raw := c.Query("idCuentaAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion debe ser un entero positivo")
		return
	}

	item, err := h.service.GetFichaAdmision(c.Request.Context(), shared.FichaAdmisionParams{IdCuentaAtencion: id})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ADMISSION_RECORD_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, item)
}

// ListMedicosPorEspecialidad retorna médicos de una especialidad.
//
//	@Summary	 Médicos por especialidad
//	@Description Lista los médicos disponibles de una especialidad específica.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		IdEspecialidad	path		int	true	"ID de la especialidad"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/medicos/{IdEspecialidad} [get]
func (h *TriageHandler) ListMedicosPorEspecialidad(c *gin.Context) {
	raw := c.Param("IdEspecialidad")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "IdEspecialidad es obligatorio")
		return
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "IdEspecialidad debe ser un entero positivo o 0")
		return
	}

	items, err := h.service.ListarMedicosPorEspecialidad(c.Request.Context(), int(id))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_MEDICOS_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// ListTriajeConsulta retorna la lista de triajes de consulta externa.
//
//	@Summary	 Listar triajes consulta
//	@Description Lista triajes de consulta externa con filtros de fecha y servicio.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		fini		query		string	true	"Fecha inicio (YYYY-MM-DD)"
//	@Param		ffin		query		string	true	"Fecha fin (YYYY-MM-DD)"
//	@Param		filtro		query		string	false	"Filtro por paciente o documento"
//	@Param		idServicio	query		int		false	"ID de servicio"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/consulta [get]
func (h *TriageHandler) ListTriajeConsulta(c *gin.Context) {
	params := shared.TriajeConsultaParams{
		FechaInicio: c.Query("fini"),
		FechaFin:    c.Query("ffin"),
		Filtro:      c.Query("filtro"),
	}

	if params.FechaInicio == "" || params.FechaFin == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fini y ffin son obligatorios (YYYY-MM-DD)")
		return
	}

	if _, err := time.Parse("2006-01-02", params.FechaInicio); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fini debe ser YYYY-MM-DD")
		return
	}
	if _, err := time.Parse("2006-01-02", params.FechaFin); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "ffin debe ser YYYY-MM-DD")
		return
	}

	if raw := c.Query("idServicio"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || v < 0 {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idServicio debe ser un entero no negativo")
			return
		}
		params.IdServicio = int(v)
	}

	items, err := h.service.ListTriajeConsulta(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_CONSULTA_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// CreateTriajeConsulta registra un triaje de consulta externa.
//
//	@Summary	 Registrar triaje consulta
//	@Description Registra signos vitales de un paciente en consulta externa.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		body body	 createTriajeConsultaRequest	true	"Datos del triaje de consulta"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/consulta [post]
func (h *TriageHandler) CreateTriajeConsulta(c *gin.Context) {
	var req createTriajeConsultaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	domainObj := req.toDomain()
	if idEmpleado := c.GetInt("idEmpleado"); idEmpleado != 0 {
		domainObj.IdEmpleado = int64(idEmpleado)
	}

	result, err := h.service.CreateTriajeConsulta(c.Request.Context(), domainObj)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_CONSULTA_REGISTER_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}

// GetTriajeConsultaPorAtencion retorna el triaje de consulta de una atención.
//
//	@Summary	 Triaje por atención
//	@Description Obtiene el triaje de consulta externa asociado a una atención.
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		idAtencion	path		int	true	"ID de la atención"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/consulta/atencion/{idAtencion} [get]
func (h *TriageHandler) GetTriajeConsultaPorAtencion(c *gin.Context) {
	raw := c.Param("idAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion es obligatorio")
		return
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion debe ser un entero positivo")
		return
	}

	item, err := h.service.GetTriajeConsultaPorAtencion(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_CONSULTA_GET_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, item)
}

// UpdateEstadoTriajeConsulta actualiza el estado de un triaje de consulta.
//
//	@Summary	 Actualizar estado triaje consulta
//	@Description Cambia el estado de un triaje de consulta externa (ej. atendido).
//	@Tags		triaje
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"ID del triaje"
//	@Param		body	body		triajeConsultaEstadoRequest	true	"Nuevo estado"
//	@Success	200	{object}	map[string]bool
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		 /triaje/consulta/{id}/estado [put]
func (h *TriageHandler) UpdateEstadoTriajeConsulta(c *gin.Context) {
	raw := c.Param("id")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id es obligatorio")
		return
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id debe ser un entero positivo")
		return
	}

	var req triajeConsultaEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if req.Estado == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "estado es obligatorio")
		return
	}

	err = h.service.UpdateEstadoTriajeConsulta(c.Request.Context(), shared.TriajeConsultaEstadoParams{
		IdTriaje: id,
		Estado:   req.Estado,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_CONSULTA_STATE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]bool{"ok": true})
}
