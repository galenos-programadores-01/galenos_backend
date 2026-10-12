package refconhttp

import (
	"log"
	"net/http"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
	"github.com/gin-gonic/gin"
)

// @Summary Generar hoja de referencia PDF del MINSA
// @Description Hace login en el portal REFCON y genera/descarga la hoja de referencia oficial en PDF, devuelta en base64
// @Tags RefCon
// @Accept json
// @Produce json
// @Param body body domain.GenerarHojaReferenciaRequest true "Datos de la referencia a generar"
// @Success 200 {object} apiResponse{data=domain.GenerarHojaReferenciaResult}
// @Router /refcon/hoja-referencia [post]
func (h *RefConHandler) HandleGenerarHojaReferencia(c *gin.Context) {
	var req domain.GenerarHojaReferenciaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "El cuerpo debe ser JSON con idreferencia, idestablecimiento y estadoreferencia")
		return
	}

	if req.IDReferencia <= 0 || req.IDEstablecimiento <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_PARAMS", "idreferencia e idestablecimiento deben ser mayores a 0")
		return
	}

	resultado, err := h.service.GenerarHojaReferencia(c.Request.Context(), req)
	if err != nil {
		log.Printf("[RefCon] Error generando hoja de referencia (id=%d): %v", req.IDReferencia, err)
		respondError(c, http.StatusBadGateway, "REFCON_HOJA_ERR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}
