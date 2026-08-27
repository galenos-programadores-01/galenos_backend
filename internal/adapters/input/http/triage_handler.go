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

	result, err := h.service.CreateTriage(c.Request.Context(), req.toDomain())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_REGISTER_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}


func (h *TriageHandler) List(c *gin.Context) {
	params := shared.TriageListParams{
		FechaInicio:       c.Query("fini"),
		FechaFin:          c.Query("ffin"),
		Filtro:            c.Query("filtro"),
		DerivadoAServicio: -100,
		IdEstado:          -100,
	}

	if params.FechaInicio == "" || params.FechaFin == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fini y ffin son obligatorios (YYYY-MM-DD)")
		return
	}

	if raw := c.Query("derivadoAServicio"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "derivadoAServicio debe ser un entero")
			return
		}
		params.DerivadoAServicio = int(v)
	}

	if raw := c.Query("idEstado"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idEstado debe ser un entero")
			return
		}
		params.IdEstado = int(v)
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

	// Validar el formato de las fechas antes de armar el fragmento WHERE.
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
