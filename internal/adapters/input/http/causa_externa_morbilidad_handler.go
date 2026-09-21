package httpadapter

import (
	"net/http"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

// CausaExternaMorbilidadHandler expone el catálogo de causas externas de
// morbilidad vía HTTP.
type CausaExternaMorbilidadHandler struct {
	service input.CausaExternaMorbilidadService
}

func NewCausaExternaMorbilidadHandler(service input.CausaExternaMorbilidadService) *CausaExternaMorbilidadHandler {
	return &CausaExternaMorbilidadHandler{service: service}
}

// @Summary Listar causas externas de morbilidad
// @Description Devuelve el catálogo de causas externas de morbilidad (uso_go_EmergenciaCausaExternaMorbilidadListar)
// @Tags Triaje
// @Accept json
// @Produce json
// @Router /triaje/causas-externas-morbilidad [get]
func (h *CausaExternaMorbilidadHandler) HandleListar(c *gin.Context) {
	causas, err := h.service.Listar(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "CAUSA_EXTERNA_MORBILIDAD_GET_ERR", "Error obteniendo el catálogo de causas externas de morbilidad")
		return
	}
	respondSuccess(c, http.StatusOK, causas)
}
