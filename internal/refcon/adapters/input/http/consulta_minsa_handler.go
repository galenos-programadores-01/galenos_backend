package refconhttp

import (
	"log"
	"net/http"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
	"github.com/gin-gonic/gin"
)

// @Summary Consultar detalle de referencia en MINSA
// @Description Consulta el detalle de una referencia en el servicio REST del MINSA por documento del paciente
// @Tags DashRefCon
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param numerodocumento query string true "Número de documento del paciente"
// @Param tipodocumento query string true "Tipo de documento (1 = DNI)"
// @Param pagina query string true "Número de página a consultar"
// @Param establecimientoDestino query string false "Código de establecimiento destino (por defecto el configurado)"
// @Param limite query string false "Límite de resultados por página (por defecto el configurado)"
// @Success 200 {object} apiResponse{data=domain.ConsultaMinsaResponse}
// @Router /dashrefcon/consulta-referencia-detalle [get]
func (h *RefConHandler) HandleConsultarReferenciaDetalle(c *gin.Context) {
	req := domain.ConsultaMinsaRequest{
		EstablecimientoDestino: c.Query("establecimientoDestino"),
		Limite:                 c.Query("limite"),
		NumeroDocumento:        c.Query("numerodocumento"),
		Pagina:                 c.Query("pagina"),
		TipoDocumento:          c.Query("tipodocumento"),
	}

	if req.NumeroDocumento == "" || req.TipoDocumento == "" || req.Pagina == "" {
		respondError(c, http.StatusBadRequest, "INVALID_PARAMS", "Los parámetros numerodocumento, tipodocumento y pagina son obligatorios")
		return
	}

	resultado, err := h.service.ConsultarReferenciaDetalle(c.Request.Context(), req)
	if err != nil {
		log.Printf("[DashRefCon] Error consultando detalle de referencia en MINSA (doc=%s): %v", req.NumeroDocumento, err)
		respondError(c, http.StatusBadGateway, "MINSA_REFCON_ERR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}
