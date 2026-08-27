package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

const (
	msgInvalidID   = "ID inválido"
	msgInvalidBody = "Cuerpo de petición inválido"
)

type InterconsultaHandler struct {
	service input.InterconsultaService
}

func NewInterconsultaHandler(service input.InterconsultaService) *InterconsultaHandler {
	return &InterconsultaHandler{service: service}
}

// @Summary Get interconsulta by ID
// @Description Returns the interconsulta details for a given ID
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la Interconsulta"
// @Router /interconsultas/{id} [get]
func (h *InterconsultaHandler) HandleObtenerPorId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", msgInvalidID)
		return
	}

	ic, err := h.service.ObtenerPorId(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_GET_ERR", "Error obteniendo interconsulta")
		return
	}

	respondSuccess(c, http.StatusOK, ic)
}

// @Summary List interconsultas by service
// @Description Returns the interconsultas for a given service type
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tipoServicio path string true "Tipo de Servicio"
// @Router /interconsultas/servicio/{tipoServicio} [get]
func (h *InterconsultaHandler) HandleListarPorServicio(c *gin.Context) {
	tipoServicio := c.Param("tipoServicio")
	lista, err := h.service.ListarPorServicio(c.Request.Context(), tipoServicio)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_LIST_ERR", "Error listando interconsultas")
		return
	}
	respondSuccess(c, http.StatusOK, lista)
}

// @Summary List interconsultas by attention ID
// @Description Returns the interconsultas associated with a specific attention/encounter
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idAtencion path int true "ID de la Atención"
// @Router /interconsultas/atencion/{idAtencion} [get]
func (h *InterconsultaHandler) HandleListarPorAtencion(c *gin.Context) {
	idStr := c.Param("idAtencion")
	idAtencion, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", msgInvalidID)
		return
	}

	lista, err := h.service.ListarPorAtencion(c.Request.Context(), idAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_LIST_ERR", "Error listando interconsultas por atención")
		return
	}
	respondSuccess(c, http.StatusOK, lista)
}

// @Summary List interconsulta specialties
// @Description Returns the specialties available for interconsulta requests
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /interconsultas/especialidades [get]
func (h *InterconsultaHandler) HandleListarEspecialidades(c *gin.Context) {
	lista, err := h.service.ListarEspecialidades(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_ESP_ERR", "Error listando especialidades de interconsulta")
		return
	}
	respondSuccess(c, http.StatusOK, lista)
}

// @Summary List doctors by specialty
// @Description Returns the doctors available for a given interconsulta specialty
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param IdEspecialidad path int true "ID de la Especialidad"
// @Router /interconsultas/medicos/{IdEspecialidad} [get]
func (h *InterconsultaHandler) HandleListarMedicosPorEspecialidad(c *gin.Context) {
	idStr := c.Param("IdEspecialidad")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", msgInvalidID)
		return
	}

	lista, err := h.service.ListarMedicosPorEspecialidad(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_MED_ERR", "Error listando médicos por especialidad")
		return
	}
	respondSuccess(c, http.StatusOK, lista)
}

type CreateInterconsultaRequest struct {
	IdAtencionOrigen int    `json:"idAtencionOrigen" binding:"required"`
	IdEspecialidad   int    `json:"IdEspecialidad" binding:"required"`
	IdMedicoDestino  int    `json:"idMedicoDestino" binding:"required"`
	Motivo           string `json:"motivo" binding:"required"`
}

// @Summary Create an interconsulta
// @Description Creates a new interconsulta request
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateInterconsultaRequest true "Interconsulta data"
// @Router /interconsultas [post]
func (h *InterconsultaHandler) HandleCrear(c *gin.Context) {
	var req CreateInterconsultaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", msgInvalidBody)
		return
	}

	ic := domain.Interconsulta{
		IdAtencionOrigen: req.IdAtencionOrigen,
		IdEspecialidad:   req.IdEspecialidad,
		IdMedicoDestino:  req.IdMedicoDestino,
		Motivo:           req.Motivo,
	}

	err := h.service.Crear(c.Request.Context(), ic)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_CREATE_ERR", "Error creando interconsulta")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Interconsulta creada correctamente"})
}

type UpdateEstadoRequest struct {
	Estado string `json:"estado" binding:"required"`
}

// @Summary Update interconsulta state
// @Description Updates the state of an existing interconsulta
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la Interconsulta"
// @Param request body UpdateEstadoRequest true "Estado data"
// @Router /interconsultas/{id}/estado [put]
func (h *InterconsultaHandler) HandleActualizarEstado(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", msgInvalidID)
		return
	}

	var req UpdateEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", msgInvalidBody)
		return
	}

	err = h.service.ActualizarEstado(c.Request.Context(), id, req.Estado)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_UPDATE_ERR", "Error actualizando estado")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Estado actualizado correctamente"})
}

type SignInterconsultaRequest struct {
	DataB64 string `json:"dataB64" binding:"required"`
}

// @Summary Sign an interconsulta
// @Description Saves the signature for an interconsulta
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la Interconsulta"
// @Param request body SignInterconsultaRequest true "Firma data"
// @Router /interconsultas/{id}/firma [post]
func (h *InterconsultaHandler) HandleGuardarFirma(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", msgInvalidID)
		return
	}

	var req SignInterconsultaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", msgInvalidBody)
		return
	}

	idEmpleado := c.GetInt("idEmpleado")
	if idEmpleado == 0 {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "No se identificó al empleado en el token")
		return
	}

	firma := domain.FirmaInterconsulta{
		IdInterconsulta: id,
		DataB64:         req.DataB64,
		IdEmpleadoFirma: idEmpleado,
	}

	err = h.service.GuardarFirma(c.Request.Context(), firma)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_SIGN_ERR", "Error guardando firma")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Firma guardada correctamente"})
}
