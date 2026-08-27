package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type SintomaHandler struct {
	service input.SintomaService
}

func NewSintomaHandler(service input.SintomaService) *SintomaHandler {
	return &SintomaHandler{service: service}
}

func idUsuarioDesdeContexto(c *gin.Context) int {
	return c.GetInt("idEmpleado")
}

// @Summary Listar catálogo de síntomas
// @Description Devuelve el catálogo de síntomas agrupado por sistema (Tab_Sintomas_Catalogo)
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /sintomas/catalogo [get]
func (h *SintomaHandler) HandleListarCatalogo(c *gin.Context) {
	sintomas, err := h.service.ListarCatalogo(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SINTOMAS_GET_ERR", "Error obteniendo el catálogo de síntomas")
		return
	}
	respondSuccess(c, http.StatusOK, sintomas)
}

type AgregarSintomaRequest struct {
	Sistema string `json:"sistema" binding:"required"`
	Sintoma string `json:"sintoma" binding:"required"`
}

// @Summary Agregar síntoma al catálogo
// @Description Inserta un síntoma nuevo en el catálogo si no existe
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AgregarSintomaRequest true "Síntoma a agregar"
// @Router /sintomas/catalogo [post]
func (h *SintomaHandler) HandleAgregarCatalogo(c *gin.Context) {
	var req AgregarSintomaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	idUsuario := idUsuarioDesdeContexto(c)
	err := h.service.AgregarCatalogo(c.Request.Context(), req.Sistema, req.Sintoma, idUsuario)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SINTOMAS_ADD_ERR", "Error agregando síntoma al catálogo")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Síntoma agregado correctamente"})
}

type GuardarSintomasRequest struct {
	Sintomas []domain.SintomaSeleccionado `json:"sintomas" binding:"required"`
}

// @Summary Guardar síntomas seleccionados de la evolución
// @Description Reemplaza los síntomas seleccionados de la evolución (sp_go_InsertarEvolucionSintomas)
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idRegAtencion path int true "ID of the Registration / Encounter"
// @Param request body GuardarSintomasRequest true "Síntomas seleccionados"
// @Router /evoluciones/{idRegAtencion}/sintomas [post]
func (h *SintomaHandler) HandleGuardarSintomas(c *gin.Context) {
	idStr := c.Param("idRegAtencion")
	idRegAtencion, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de atención inválido")
		return
	}

	var req GuardarSintomasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	idUsuario := idUsuarioDesdeContexto(c)
	err = h.service.GuardarEvolucionSintomas(c.Request.Context(), idRegAtencion, req.Sintomas, idUsuario)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SINTOMAS_SAVE_ERR", "Error guardando síntomas de la evolución")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Síntomas guardados correctamente"})
}
