package httpadapter

import (
	"net/http"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type MedicoListaEsperaHandler struct {
	service input.MedicoListaEsperaService
}

func NewMedicoListaEsperaHandler(service input.MedicoListaEsperaService) *MedicoListaEsperaHandler {
	return &MedicoListaEsperaHandler{service: service}
}

// @Summary Listar medicos para lista de espera
// @Description Retorna la lista de medicos activos para la lista de espera quirurgica
// @Tags MedicosListaEspera
// @Produce json
// @Success 200 {object} apiResponse{data=object}
// @Router /medicos-lista-espera [get]
func (h *MedicoListaEsperaHandler) HandleListar(c *gin.Context) {
	lista, err := h.service.Listar(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "MEDICOS_LISTA_ESPERA_ERR", "Error obteniendo medicos")
		return
	}
	respondSuccess(c, http.StatusOK, lista)
}
