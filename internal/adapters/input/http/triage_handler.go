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

// @Summary Modifica un triaje
// @Description Modifica los datos clínicos de un triaje de emergencia invocando el SP usp_go_ModificarTriaje
// @Accept json
// @Produce json
// @Param id path int true "Id del triaje"
// @Param triaje body createTriajeRequest true "Datos clínicos editables del triaje"
// @Success 200 {object} map[string]string "Mensaje de resultado del SP"
// @Failure 400 {object} object "Error de validación"
// @Failure 500 {object} object "Error interno"
// @Router /triaje/{id} [put]
func (h *TriageHandler) UpdateTriaje(c *gin.Context) {
	var req createTriajeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	raw := c.Param("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id debe ser un entero positivo")
		return
	}

	domainObj := req.toDomain()
	domainObj.IDTriaje = &id
	if idEmpleado := c.GetInt("idEmpleado"); idEmpleado != 0 {
		empID := int64(idEmpleado)
		domainObj.EmployeeID = &empID
	}

	result, err := h.service.UpdateTriage(c.Request.Context(), domainObj)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_UPDATE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}

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

// @Summary Datos de un triaje por id
// @Description Devuelve los datos del paciente (solo lectura) y del triaje de emergencia del id indicado (SP usp_go_Triaje_EmergeciaPorId)
// @Accept json
// @Produce json
// @Param id path int true "Id del triaje"
// @Success 200 {object} map[string]interface{} "Datos del paciente y del triaje"
// @Failure 400 {object} object "Error de validación"
// @Failure 404 {object} object "Triaje no encontrado"
// @Router /triaje/{id} [get]
func (h *TriageHandler) GetTriajePorId(c *gin.Context) {
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

	item, err := h.service.GetTriajePorId(c.Request.Context(), int(id))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_GET_FAILED", err.Error())
		return
	}
	if item == nil {
		respondError(c, http.StatusNotFound, "TRIAGE_NOT_FOUND", "No se encontró el triaje")
		return
	}

	respondSuccess(c, http.StatusOK, item)
}

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

// @Summary Actualiza la IAFA de un triaje
// @Description Actualiza la fuente de financiamiento (IAFA) de un triaje de emergencia a SIS invocando el SP usp_go_Triaje_EmergeciaActualizarIAFA
// @Tags Triaje
// @Produce json
// @Param id path int true "Id del triaje"
// @Success 200 {object} map[string]bool "Operación exitosa"
// @Failure 400 {object} object "Error de validación"
// @Failure 500 {object} object "Error interno"
// @Router /triaje/{id}/iafa [put]
func (h *TriageHandler) UpdateIafa(c *gin.Context) {
	raw := c.Param("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id debe ser un entero positivo")
		return
	}

	if err := h.service.UpdateIafa(c.Request.Context(), int(id)); err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_IAFA_UPDATE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]bool{"ok": true})
}

// @Summary Reporte de triajes por empleado
// @Description Cantidad de triajes de un empleado por servicio (tópico) entre fechas con hora y minuto (SP usp_go_ReporteTriaje)
// @Accept json
// @Produce json
// @Param IdEmpleado query int true "Id del empleado"
// @Param fechaini query string true "Fecha inicial (YYYY-MM-DD o YYYY-MM-DD HH:MM)"
// @Param fechafin query string true "Fecha final (YYYY-MM-DD o YYYY-MM-DD HH:MM)"
// @Success 200 {array} map[string]interface{} "Cantidad de triajes por tópico"
// @Failure 400 {object} object "Error de validación"
// @Router /triaje/reporte-por-empleado [get]
func (h *TriageHandler) GetReporteTriaje(c *gin.Context) {
	rawEmpleado := c.Query("IdEmpleado")
	if rawEmpleado == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "IdEmpleado es obligatorio")
		return
	}
	idEmpleado, err := strconv.ParseInt(rawEmpleado, 10, 64)
	if err != nil || idEmpleado <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "IdEmpleado debe ser un entero positivo")
		return
	}

	fechaIni := c.Query("fechaini")
	fechaFin := c.Query("fechafin")
	if fechaIni == "" || fechaFin == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fechaini y fechafin son obligatorias (YYYY-MM-DD HH:MM)")
		return
	}
	layouts := []string{
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006-01-02T15:04:05.000",
	}
	if !esFechaValida(fechaIni, layouts) || !esFechaValida(fechaFin, layouts) {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fechaini y fechafin deben ser fechas válidas (YYYY-MM-DD HH:MM)")
		return
	}

	items, err := h.service.ReporteTriajePorEmpleado(c.Request.Context(), shared.ReporteTriajeParams{
		IdEmpleado: int(idEmpleado),
		FechaIni:   fechaIni,
		FechaFin:   fechaFin,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_REPORT_BY_EMPLOYEE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// @Summary Bandeja de referencias de consulta externa
// @Description Lista las referencias hechas desde consulta externa invocando el SP usp_go_ListarBandejaReferencia
// @Accept json
// @Produce json
// @Param fini query string true "Fecha inicial (YYYY-MM-DD)"
// @Param ffin query string true "Fecha final (YYYY-MM-DD)"
// @Param filtro query string false "Texto de búsqueda por paciente o número de cuenta"
// @Success 200 {array} map[string]interface{} "Referencias de consulta externa"
// @Failure 400 {object} object "Error de validación"
// @Failure 500 {object} object "Error interno"
// @Router /triaje/referencias [get]
func (h *TriageHandler) ListReferenciasConsultaExterna(c *gin.Context) {
	params := shared.ReferenciaParams{
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

	items, err := h.service.ListReferenciasConsultaExterna(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFERENCIA_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// @Summary Datos de una referencia
// @Description Devuelve los datos del destino de una referencia (establecimiento, UPS y especialidad) invocando el SP usp_go_ListarDatosReferencia
// @Accept json
// @Produce json
// @Param idAtencion path int true "Id de la atención"
// @Success 200 {object} map[string]interface{} "Datos de la referencia"
// @Failure 400 {object} object "Error de validación"
// @Failure 500 {object} object "Error interno"
// @Router /triaje/referencias/{idAtencion} [get]
func (h *TriageHandler) ListarDatosReferencia(c *gin.Context) {
	raw := c.Param("idAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion es obligatorio")
		return
	}
	idAtencion, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion debe ser un entero positivo")
		return
	}

	item, err := h.service.ListarDatosReferencia(c.Request.Context(), int(idAtencion))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFERENCIA_GET_FAILED", err.Error())
		return
	}
	if item == nil {
		respondError(c, http.StatusNotFound, "REFERENCIA_NOT_FOUND", "No se encontraron datos de la referencia")
		return
	}

	respondSuccess(c, http.StatusOK, *item)
}

// @Summary Cabecera JSON de una referencia
// @Description Devuelve la cabecera con los datos completos de una referencia (paciente, responsable, personal que registra, cita, diagnósticos y destino) invocando el SP usp_go_webReferenciaJSONCab
// @Accept json
// @Produce json
// @Param idCuentaAtencion path int true "Id de la cuenta de atención"
// @Success 200 {object} map[string]interface{} "Cabecera con los datos de la referencia"
// @Failure 400 {object} object "Error de validación"
// @Failure 404 {object} object "Referencia no encontrada"
// @Failure 500 {object} object "Error interno"
// @Router /triaje/referencias-cabecera/{idCuentaAtencion} [get]
func (h *TriageHandler) WebReferenciaJSONCab(c *gin.Context) {
	raw := c.Param("idCuentaAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idCuentaAtencion, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || idCuentaAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion debe ser un entero positivo")
		return
	}

	item, err := h.service.WebReferenciaJSONCab(c.Request.Context(), int(idCuentaAtencion))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFERENCIA_GET_FAILED", err.Error())
		return
	}
	if item == nil {
		respondError(c, http.StatusNotFound, "REFERENCIA_NOT_FOUND", "No se encontraron datos de la referencia")
		return
	}

	respondSuccess(c, http.StatusOK, *item)
}

// @Summary Detalle JSON de una referencia
// @Description Devuelve el detalle de la referencia (los diagnósticos) invocando el SP usp_go_webReferenciaJSONDet
// @Accept json
// @Produce json
// @Param idCuentaAtencion path int true "Id de la cuenta de atención"
// @Success 200 {array} map[string]interface{} "Detalle con los diagnósticos de la referencia"
// @Failure 400 {object} object "Error de validación"
// @Failure 500 {object} object "Error interno"
// @Router /triaje/referencias-detalle/{idCuentaAtencion} [get]
func (h *TriageHandler) WebReferenciaJSONDet(c *gin.Context) {
	raw := c.Param("idCuentaAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idCuentaAtencion, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || idCuentaAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion debe ser un entero positivo")
		return
	}

	items, err := h.service.WebReferenciaJSONDet(c.Request.Context(), int(idCuentaAtencion))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFERENCIA_DET_GET_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// @Summary Tratamiento JSON de una referencia
// @Description Devuelve el tratamiento de la referencia invocando el SP usp_go_WebTratamientoReferenciaJSON
// @Accept json
// @Produce json
// @Param idCuentaAtencion path int true "Id de la cuenta de atención"
// @Success 200 {array} map[string]interface{} "Tratamiento de la referencia"
// @Failure 400 {object} object "Error de validación"
// @Failure 500 {object} object "Error interno"
// @Router /triaje/referencias-tratamiento/{idCuentaAtencion} [get]
func (h *TriageHandler) WebTratamientoReferenciaJSON(c *gin.Context) {
	raw := c.Param("idCuentaAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idCuentaAtencion, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || idCuentaAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion debe ser un entero positivo")
		return
	}

	items, err := h.service.WebTratamientoReferenciaJSON(c.Request.Context(), int(idCuentaAtencion))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFERENCIA_TRATAMIENTO_GET_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// @Summary CPT de una referencia
// @Description Devuelve los códigos CPT de la referencia invocando el SP usp_go_webCPTReferenciaJSON
// @Accept json
// @Produce json
// @Param idCuentaAtencion path int true "Id de la cuenta de atención"
// @Success 200 {array} map[string]interface{} "CPT de la referencia"
// @Failure 400 {object} object "Error de validación"
// @Failure 500 {object} object "Error interno"
// @Router /triaje/referencias-cpt/{idCuentaAtencion} [get]
func (h *TriageHandler) WebCPTReferenciaJSON(c *gin.Context) {
	raw := c.Param("idCuentaAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idCuentaAtencion, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || idCuentaAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion debe ser un entero positivo")
		return
	}

	items, err := h.service.WebCPTReferenciaJSON(c.Request.Context(), int(idCuentaAtencion))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFERENCIA_CPT_GET_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

func esFechaValida(valor string, layouts []string) bool {
	for _, l := range layouts {
		if _, err := time.Parse(l, valor); err == nil {
			return true
		}
	}
	return false
}
