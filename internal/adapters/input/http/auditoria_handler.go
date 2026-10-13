package httpadapter

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

type AuditoriaHandler struct {
	service input.AuditoriaService
}

func NewAuditoriaHandler(service input.AuditoriaService) *AuditoriaHandler {
	return &AuditoriaHandler{service: service}
}

type AgregarAuditoriaRequest struct {
	Accion        string `json:"accion"`
	IdRegistro    int    `json:"idRegistro"`
	Tabla         string `json:"tabla"`
	IdListItem    int    `json:"idListItem"`
	NombrePC      string `json:"nombrePC"`
	Observaciones string `json:"observaciones"`
}

// @Summary Registrar auditoría
// @Description Inserta un registro en la tabla Auditoria mediante dbo.AuditoriaAgregarV
// @Tags Auditoria
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AgregarAuditoriaRequest true "Registro de auditoría a agregar"
// @Router /auditoria [post]
func (h *AuditoriaHandler) HandleRegistrar(c *gin.Context) {
	var req AgregarAuditoriaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	accion := strings.ToUpper(strings.TrimSpace(req.Accion))
	if accion == "" || len(accion) > 1 {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "El campo accion debe ser un único carácter")
		return
	}

	reg := domain.AuditoriaRegistro{
		IdEmpleado:    idUsuarioDesdeContexto(c),
		Accion:        accion,
		IdRegistro:    req.IdRegistro,
		Tabla:         strings.TrimSpace(req.Tabla),
		IdListItem:    req.IdListItem,
		NombrePC:      strings.TrimSpace(req.NombrePC),
		Observaciones: strings.TrimSpace(req.Observaciones),
	}

	if err := h.service.Registrar(c.Request.Context(), reg); err != nil {
		respondError(c, http.StatusInternalServerError, "AUDITORIA_ADD_ERR", "Error registrando auditoría")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Auditoría registrada correctamente"})
}
