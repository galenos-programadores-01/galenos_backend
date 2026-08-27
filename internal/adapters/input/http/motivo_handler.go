package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type MotivoHandler struct {
	service input.MotivoService
}

func NewMotivoHandler(service input.MotivoService) *MotivoHandler {
	return &MotivoHandler{service: service}
}

// @Summary Get motivos de atencion for a patient evolution
// @Description Returns the saved reasons for attention for a given registration ID
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idRegAtencion path int true "ID of the Registration / Encounter"
// @Router /evoluciones/{idRegAtencion}/motivos [get]
func (h *MotivoHandler) HandleListMotivos(c *gin.Context) {
	idStr := c.Param("idRegAtencion")
	idRegAtencion, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de atención inválido")
		return
	}

	motivos, err := h.service.ListarMotivos(c.Request.Context(), idRegAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "MOTIVOS_GET_ERR", "Error obteniendo motivos de atención")
		return
	}

	respondSuccess(c, http.StatusOK, motivos)
}

type SaveMotivoRequest struct {
	Motivo string `json:"motivo" binding:"required"`
}

// @Summary Save a motivo de atencion
// @Description Saves a new reason for attention in a patient's evolution
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idRegAtencion path int true "ID of the Registration / Encounter"
// @Param request body SaveMotivoRequest true "Motivo data"
// @Router /evoluciones/{idRegAtencion}/motivos [post]
func (h *MotivoHandler) HandleCreateMotivo(c *gin.Context) {
	idStr := c.Param("idRegAtencion")
	idRegAtencion, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de atención inválido")
		return
	}

	var req SaveMotivoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	err = h.service.GuardarMotivo(c.Request.Context(), idRegAtencion, req.Motivo)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "MOTIVOS_SAVE_ERR", "Error guardando el motivo de atención")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Motivo guardado correctamente"})
}
