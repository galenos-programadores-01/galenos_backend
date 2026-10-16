package httpadapter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

// AuditoriaQfHandler expone la bandeja de recetas de farmacia.
type AuditoriaQfHandler struct {
	service input.AuditoriaQfService
}

// NewAuditoriaQfHandler construye el handler HTTP.
func NewAuditoriaQfHandler(service input.AuditoriaQfService) *AuditoriaQfHandler {
	return &AuditoriaQfHandler{service: service}
}

// ListarBandeja maneja GET /api/v1/farmacia/auditoria-qf/recetas.
func (h *AuditoriaQfHandler) ListarBandeja(c *gin.Context) {
	fechaInicio := c.Query("fechaInicio")
	fechaFin := c.Query("fechaFin")
	if fechaInicio == "" || fechaFin == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fechaInicio y fechaFin son obligatorias")
		return
	}

	params := domain.BandejaRecetasParams{
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
	}
	if raw := c.Query("idEstado"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idEstado debe ser un entero")
			return
		}
		params.IdEstado = &v
	}

	items, err := h.service.ListarBandeja(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err, domain.ErrFechaBandejaInvalida) ||
			errors.Is(err, domain.ErrFechaBandejaOrden) ||
			errors.Is(err, domain.ErrRangoFechasBandeja) {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "AUDITORIA_QF_LIST_FAILED", err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, items)
}

// ListarEstados maneja GET /api/v1/farmacia/auditoria-qf/estados.
func (h *AuditoriaQfHandler) ListarEstados(c *gin.Context) {
	items, err := h.service.ListarEstados(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "AUDITORIA_QF_ESTADOS_FAILED", err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, items)
}

// ListarDetalle maneja GET /api/v1/farmacia/auditoria-qf/recetas/:idReceta/detalle.
func (h *AuditoriaQfHandler) ListarDetalle(c *gin.Context) {
	idReceta, err := strconv.Atoi(c.Param("idReceta"))
	if err != nil || idReceta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idReceta debe ser un entero")
		return
	}
	idUsuario := c.GetInt("idEmpleado")
	if idUsuario <= 0 {
		respondError(c, http.StatusUnauthorized, "MISSING_USER", "no se encontró el usuario de la sesión")
		return
	}

	items, err := h.service.ListarDetalle(c.Request.Context(), domain.DetalleRecetaParams{
		IdReceta:  idReceta,
		IdUsuario: idUsuario,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "AUDITORIA_QF_DETALLE_FAILED", err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, items)
}

// ListarHistorico maneja GET /api/v1/farmacia/auditoria-qf/historico.
func (h *AuditoriaQfHandler) ListarHistorico(c *gin.Context) {
	idPaciente, okPaciente := enteroPositivo(c.Query("idPaciente"))
	idProducto, okProducto := enteroPositivo(c.Query("idProducto"))
	if !okPaciente || !okProducto {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPaciente e idProducto deben ser enteros")
		return
	}

	items, err := h.service.ListarHistorico(c.Request.Context(), domain.HistoricoMedicamentoParams{
		IdPaciente: idPaciente,
		IdProducto: idProducto,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "AUDITORIA_QF_HISTORICO_FAILED", err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, items)
}

func enteroPositivo(raw string) (int, bool) {
	valor, err := strconv.Atoi(raw)
	if err != nil || valor <= 0 {
		return 0, false
	}
	return valor, true
}
