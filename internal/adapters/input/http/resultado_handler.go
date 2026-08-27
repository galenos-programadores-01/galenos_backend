package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type ResultadoHandler struct {
	service input.ResultadoService
}

func NewResultadoHandler(service input.ResultadoService) *ResultadoHandler {
	return &ResultadoHandler{service: service}
}

// @Summary Get laboratory results for a patient
// @Description Returns laboratory results history for a given patient ID
// @Tags Resultados
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idPaciente path int true "ID of the Patient"
// @Router /resultados/laboratorio/paciente/{idPaciente} [get]
func (h *ResultadoHandler) HandleListResultadosLaboratorio(c *gin.Context) {
	idStr := c.Param("idPaciente")
	idPaciente, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de paciente inválido")
		return
	}

	resultados, err := h.service.ListarResultadosLaboratorio(c.Request.Context(), idPaciente)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "RES_LAB_ERR", "Error obteniendo resultados de laboratorio")
		return
	}

	respondSuccess(c, http.StatusOK, resultados)
}

// @Summary Get imaging results for a patient
// @Description Returns imaging results history for a given patient ID
// @Tags Resultados
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idPaciente path int true "ID of the Patient"
// @Router /resultados/imagenes/paciente/{idPaciente} [get]
func (h *ResultadoHandler) HandleListResultadosImagenes(c *gin.Context) {
	idStr := c.Param("idPaciente")
	idPaciente, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de paciente inválido")
		return
	}

	resultados, err := h.service.ListarResultadosImagenes(c.Request.Context(), idPaciente)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "RES_IMG_ERR", "Error obteniendo resultados de imágenes")
		return
	}

	respondSuccess(c, http.StatusOK, resultados)
}
