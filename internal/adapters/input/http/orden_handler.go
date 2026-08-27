package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type OrdenHandler struct {
	service input.OrdenService
}

func NewOrdenHandler(service input.OrdenService) *OrdenHandler {
	return &OrdenHandler{service: service}
}

// @Summary Get ordenes for a patient account
// @Description Returns the medical orders for a given account ID
// @Tags Ordenes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idRegAtencion path int true "ID de la Atencion"
// @Router /ordenes/cuenta/{idRegAtencion} [get]
func (h *OrdenHandler) HandleListOrdenes(c *gin.Context) {
	idStr := c.Param("idCuentaAtencion")
	idRegAtencion, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de atención inválido")
		return
	}

	ordenes, err := h.service.ListarPorCuenta(c.Request.Context(), idRegAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ORDEN_GET_ERR", "Error obteniendo órdenes médicas")
		return
	}

	respondSuccess(c, http.StatusOK, ordenes)
}

type CreateOrdenRequest struct {
	IdRegAtencion int                   `json:"idRegAtencion" binding:"required"`
	Observacion   string                `json:"observacion"`
	Detalles      []domain.DetalleOrden `json:"detalles" binding:"required"`
}

// @Summary Create an orden medica
// @Description Creates a new medical order in a patient's evolution
// @Tags Ordenes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrdenRequest true "Orden data"
// @Router /ordenes [post]
func (h *OrdenHandler) HandleCreateOrden(c *gin.Context) {
	var req CreateOrdenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	idEmpleado := 0
	if val, exists := c.Get("idEmpleado"); exists {
		switch v := val.(type) {
		case int:
			idEmpleado = v
		case int64:
			idEmpleado = int(v)
		case float64:
			idEmpleado = int(v)
		}
	}
	if idEmpleado == 0 {
		respondError(c, http.StatusUnauthorized, "NO_EMPLOYEE", "No se pudo identificar al empleado autenticado")
		return
	}

	orden := domain.OrdenMedica{
		IdRegAtencion: req.IdRegAtencion,
		Observacion:   req.Observacion,
	}

	err := h.service.CrearOrden(c.Request.Context(), orden, req.Detalles, idEmpleado)
	if err != nil {
		respondError(c, http.StatusBadRequest, "ORDEN_SAVE_ERR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Orden médica guardada correctamente"})
}

// @Summary Buscar productos del catálogo
// @Description Busca productos/medicamentos en el catálogo de bienes e insumos con su precio de venta vigente
// @Tags Ordenes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string false "Filtro por nombre, código o nombre comercial"
// @Param limite query int false "Máximo de resultados (default 50)"
// @Router /ordenes/productos [get]
func (h *OrdenHandler) HandleBuscarProductos(c *gin.Context) {
	filtro := c.Query("q")
	limite, _ := strconv.Atoi(c.DefaultQuery("limite", "50"))

	productos, err := h.service.BuscarProductos(c.Request.Context(), filtro, limite)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "CATALOGO_GET_ERR", "Error buscando productos del catálogo")
		return
	}

	respondSuccess(c, http.StatusOK, productos)
}
