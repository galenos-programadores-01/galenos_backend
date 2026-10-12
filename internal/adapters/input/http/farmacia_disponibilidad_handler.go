package httpadapter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type FarmaciaDisponibilidadHandler struct {
	servicio input.FarmaciaDisponibilidadService
}

func NewFarmaciaDisponibilidadHandler(servicio input.FarmaciaDisponibilidadService) *FarmaciaDisponibilidadHandler {
	return &FarmaciaDisponibilidadHandler{
		servicio: servicio,
	}
}

func (h *FarmaciaDisponibilidadHandler) ListarAlmacenes(c *gin.Context) {
	almacenes, err := h.servicio.ListarAlmacenes(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_ALMACENES", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, almacenes)
}

func (h *FarmaciaDisponibilidadHandler) ObtenerKpis(c *gin.Context) {
	idAlmacen, _ := strconv.Atoi(c.DefaultQuery("idAlmacen", "0"))
	canastaTipo := c.DefaultQuery("canastaTipo", "TODOS")

	kpis, err := h.servicio.ObtenerKpis(c.Request.Context(), idAlmacen, canastaTipo)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_OBTENER_KPIS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, kpis)
}

func (h *FarmaciaDisponibilidadHandler) ListarDisponibilidad(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	idAlmacen, _ := strconv.Atoi(c.DefaultQuery("idAlmacen", "0"))
	filtro := c.Query("filtro")
	estadoSemaforo := c.Query("estadoSemaforo")
	canastaTipo := c.DefaultQuery("canastaTipo", "TODOS")
	fechaInicio := c.Query("fechaInicio")
	fechaFin := c.Query("fechaFin")

	filtros := domain.FarmaciaDisponibilidadFiltros{
		IdAlmacen:      idAlmacen,
		Filtro:         filtro,
		EstadoSemaforo: estadoSemaforo,
		CanastaTipo:    canastaTipo,
		FechaInicio:    fechaInicio,
		FechaFin:       fechaFin,
		Page:           page,
		PageSize:       pageSize,
	}

	if stockQuery := c.Query("soloConStock"); stockQuery != "" {
		conStock := stockQuery == "true" || stockQuery == "1"
		filtros.SoloConStock = &conStock
	}

	if esencialesQuery := c.Query("soloEsenciales"); esencialesQuery != "" {
		esenciales := esencialesQuery == "true" || esencialesQuery == "1"
		filtros.SoloEsenciales = &esenciales
	}

	if movQuery := c.Query("soloConMovimiento"); movQuery != "" {
		conMov := movQuery == "true" || movQuery == "1"
		filtros.SoloConMovimiento = &conMov
	}

	resultado, err := h.servicio.ListarDisponibilidad(c.Request.Context(), filtros)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_DISPONIBILIDAD", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}

func (h *FarmaciaDisponibilidadHandler) ListarLotesPorProducto(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, "ID_INVALIDO", "El identificador del producto no es válido")
		return
	}

	idAlmacen, _ := strconv.Atoi(c.DefaultQuery("idAlmacen", "0"))

	lotes, err := h.servicio.ListarLotesPorProducto(c.Request.Context(), idProducto, idAlmacen)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_LOTES", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, lotes)
}

func (h *FarmaciaDisponibilidadHandler) ListarLotesPorVencer(c *gin.Context) {
	idAlmacen, _ := strconv.Atoi(c.DefaultQuery("idAlmacen", "0"))
	diasMax, _ := strconv.Atoi(c.DefaultQuery("diasMax", "60"))

	lotes, err := h.servicio.ListarLotesPorVencer(c.Request.Context(), idAlmacen, diasMax)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_LOTES_VENCER", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, lotes)
}

func (h *FarmaciaDisponibilidadHandler) ConsultarKardex(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, "ID_INVALIDO", "El identificador del producto no es válido")
		return
	}

	idAlmacen, _ := strconv.Atoi(c.DefaultQuery("idAlmacen", "0"))

	fechaInicioStr := c.Query("fechaInicio")
	fechaFinStr := c.Query("fechaFin")

	var fechaInicio, fechaFin time.Time

	if fechaInicioStr != "" {
		if t, err := time.Parse("2006-01-02", fechaInicioStr); err == nil {
			fechaInicio = t
		}
	}
	if fechaInicio.IsZero() {
		fechaInicio = time.Now().AddDate(0, -1, 0)
	}

	if fechaFinStr != "" {
		if t, err := time.Parse("2006-01-02", fechaFinStr); err == nil {
			fechaFin = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
		}
	}
	if fechaFin.IsZero() {
		fechaFin = time.Now()
	}

	movimientos, err := h.servicio.ConsultarKardex(c.Request.Context(), idProducto, idAlmacen, fechaInicio, fechaFin)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_CONSULTAR_KARDEX", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, movimientos)
}

func (h *FarmaciaDisponibilidadHandler) ObtenerTendenciasSalidas(c *gin.Context) {
	idAlmacen, _ := strconv.Atoi(c.DefaultQuery("idAlmacen", "0"))

	tendencias, err := h.servicio.ObtenerTendenciasSalidas(c.Request.Context(), idAlmacen)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_OBTENER_TENDENCIAS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, tendencias)
}

func (h *FarmaciaDisponibilidadHandler) ObtenerEstadisticas(c *gin.Context) {
	idAlmacen, _ := strconv.Atoi(c.DefaultQuery("idAlmacen", "0"))
	canastaTipo := c.DefaultQuery("canastaTipo", "TODOS")
	fechaInicio := c.Query("fechaInicio")
	fechaFin := c.Query("fechaFin")

	estadisticas, err := h.servicio.ObtenerEstadisticasDashboard(c.Request.Context(), idAlmacen, canastaTipo, fechaInicio, fechaFin)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_OBTENER_ESTADISTICAS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, estadisticas)
}
