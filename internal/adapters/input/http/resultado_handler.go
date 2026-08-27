package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
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
// @Router /api/v1/resultados/laboratorio/paciente/{idPaciente} [get]
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

	if resultados == nil {
		resultados = make([]domain.Resultado, 0)
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
// @Router /api/v1/resultados/imagenes/paciente/{idPaciente} [get]
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

	if resultados == nil {
		resultados = make([]domain.Resultado, 0)
	}

	respondSuccess(c, http.StatusOK, resultados)
}

// @Summary Get laboratory detailed results by order and product CPT
// @Tags Resultados
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idOrden query int true "ID de Orden"
// @Param idProducto query int true "ID de Producto CPT"
// @Router /api/v1/resultados/laboratorio/detalle [get]
func (h *ResultadoHandler) HandleObtenerDetalleLaboratorio(c *gin.Context) {
	idOrden, _ := strconv.Atoi(c.Query("idOrden"))
	idProducto, _ := strconv.Atoi(c.Query("idProducto"))

	detalles, err := h.service.ObtenerDetalleLaboratorio(c.Request.Context(), idOrden, idProducto)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "RES_LAB_DET_ERR", "Error obteniendo detalle de laboratorio")
		return
	}

	if detalles == nil {
		detalles = make([]domain.DetalleResultadoLab, 0)
	}

	respondSuccess(c, http.StatusOK, detalles)
}

// @Summary Get imaging detailed results by order and product ID
// @Tags Resultados
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idOrden query int true "ID de Orden"
// @Param idProducto query int true "ID de Producto"
// @Router /api/v1/resultados/imagenes/detalle [get]
func (h *ResultadoHandler) HandleObtenerDetalleImagen(c *gin.Context) {
	idOrden, _ := strconv.Atoi(c.Query("idOrden"))
	idProducto, _ := strconv.Atoi(c.Query("idProducto"))

	detalle, err := h.service.ObtenerDetalleImagen(c.Request.Context(), idOrden, idProducto)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "RES_IMG_DET_ERR", "Error obteniendo detalle de imágenes")
		return
	}

	if detalle == nil {
		detalle = &domain.DetalleResultadoImagen{}
	}

	respondSuccess(c, http.StatusOK, detalle)
}
