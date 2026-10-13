package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
	"github.com/galenos-pro/appointments-api/internal/usecase"
)

// PrimeraAtencionHandler expone los endpoints HTTP de Primera Atención.
type PrimeraAtencionHandler struct {
	service    input.PrimeraAtencionService
	archivos   *usecase.ArchivoRemotoService // SFTP reutilizable
	firmaStore output.FirmaPeruStore         // PDF firmado en memoria/disco
}

// NewPrimeraAtencionHandler construye el handler HTTP.
// archivos / firmaStore pueden ser nil (sin SFTP o sin store).
func NewPrimeraAtencionHandler(
	service input.PrimeraAtencionService,
	archivos *usecase.ArchivoRemotoService,
	firmaStore output.FirmaPeruStore,
) *PrimeraAtencionHandler {
	return &PrimeraAtencionHandler{
		service:    service,
		archivos:   archivos,
		firmaStore: firmaStore,
	}
}

// ListarBandeja maneja GET /api/v1/primera-atencion.
//
// @Summary Lista la bandeja de Primera Atención (emergencia)
// @Description Invoca web_go_GestionPrimeraAtencionFiltrar. Marca "Atendido" si ya se firmó y "En curso" si hay borrador.
// @Tags PrimeraAtencion
// @Produce json
// @Param filtro query string false "Texto de búsqueda (nombre, documento, cuenta…)"
// @Param idTipoServicio query int false "Id tipo de servicio (p. ej. 2 = emergencia). Default 2"
// @Success 200 {object} apiResponse{data=[]object}
// @Failure 401 {object} apiResponse{error=apiError} "Token inválido o ausente"
// @Failure 500 {object} apiResponse{error=apiError} "Error al listar la bandeja"
// @Router /primera-atencion [get]
// @Security BearerAuth
func (h *PrimeraAtencionHandler) ListarBandeja(c *gin.Context) {
	idEmpleado := c.GetInt("idEmpleado")
	if idEmpleado == 0 {
		respondError(c, http.StatusUnauthorized, "MISSING_EMPLOYEE", "el token no contiene idEmpleado")
		return
	}

	params := shared.PrimeraAtencionParams{
		Filtro:         c.Query("filtro"),
		IdEmpleado:     idEmpleado,
		IdTipoServicio: 2, // emergencia por defecto
	}

	if raw := c.Query("idTipoServicio"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idTipoServicio debe ser un entero")
			return
		}
		params.IdTipoServicio = int(v)
	}

	items, err := h.service.ListarBandeja(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "PRIMERA_ATENCION_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}
