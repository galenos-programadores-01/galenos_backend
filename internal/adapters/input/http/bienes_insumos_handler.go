package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

const (
	errCodeIdProductoInvalido = "ID_INVALIDO"
	msgIdProductoInvalido     = "El identificador de producto no es válido"
)

type BienesInsumosHandler struct {
	servicio input.BienesInsumosService
}

func NewBienesInsumosHandler(servicio input.BienesInsumosService) *BienesInsumosHandler {
	return &BienesInsumosHandler{servicio: servicio}
}

func (h *BienesInsumosHandler) Listar(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	idGrupo, _ := strconv.Atoi(c.DefaultQuery("idGrupo", "0"))
	tipoProducto, _ := strconv.Atoi(c.DefaultQuery("tipoProducto", "0"))
	filtro := c.Query("filtro")

	filtros := domain.BienInsumoFiltros{
		Filtro:               filtro,
		IdGrupoFarmacologico: idGrupo,
		TipoProducto:         tipoProducto,
		Page:                 page,
		PageSize:             pageSize,
	}

	if petitorioQuery := c.Query("petitorio"); petitorioQuery != "" {
		esPetitorio := petitorioQuery == "true" || petitorioQuery == "1"
		filtros.SoloPetitorio = &esPetitorio
	}

	if indicadorQuery := c.Query("indicadorGestion"); indicadorQuery != "" {
		esIndicador := indicadorQuery == "true" || indicadorQuery == "1"
		filtros.SoloIndicadorGestion = &esIndicador
	} else if esencialQuery := c.Query("esencial"); esencialQuery != "" {
		esIndicador := esencialQuery == "true" || esencialQuery == "1"
		filtros.SoloIndicadorGestion = &esIndicador
	}

	if estadoQuery := c.Query("estado"); estadoQuery != "" {
		if estadoVal, err := strconv.Atoi(estadoQuery); err == nil {
			filtros.Estado = &estadoVal
		}
	}

	resultado, err := h.servicio.Listar(c.Request.Context(), filtros)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_BIENES", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}

func (h *BienesInsumosHandler) ObtenerPorID(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, errCodeIdProductoInvalido, msgIdProductoInvalido)
		return
	}

	detalle, err := h.servicio.ObtenerPorID(c.Request.Context(), idProducto)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_OBTENER_BIEN", err.Error())
		return
	}
	if detalle == nil {
		respondError(c, http.StatusNotFound, "NO_ENCONTRADO", "No se encontró el bien o insumo solicitado")
		return
	}

	respondSuccess(c, http.StatusOK, detalle)
}

func (h *BienesInsumosHandler) Crear(c *gin.Context) {
	var payload domain.BienInsumoGuardarPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "PAYLOAD_INVALIDO", err.Error())
		return
	}

	nuevoID, err := h.servicio.Crear(c.Request.Context(), payload)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_CREAR_BIEN", err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, gin.H{"idProducto": nuevoID})
}

func (h *BienesInsumosHandler) Actualizar(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, errCodeIdProductoInvalido, msgIdProductoInvalido)
		return
	}

	var payload domain.BienInsumoGuardarPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "PAYLOAD_INVALIDO", err.Error())
		return
	}

	if err := h.servicio.Actualizar(c.Request.Context(), idProducto, payload); err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_ACTUALIZAR_BIEN", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"actualizado": true})
}

func (h *BienesInsumosHandler) ListarGrupos(c *gin.Context) {
	grupos, err := h.servicio.ListarGruposFarmacologicos(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_GRUPOS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, grupos)
}

func (h *BienesInsumosHandler) ListarSubGrupos(c *gin.Context) {
	idGrupo, _ := strconv.Atoi(c.Param("idGrupo"))
	subgrupos, err := h.servicio.ListarSubGruposFarmacologicos(c.Request.Context(), idGrupo)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_SUBGRUPOS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, subgrupos)
}

func (h *BienesInsumosHandler) ObtenerAuxiliares(c *gin.Context) {
	auxiliares, err := h.servicio.ObtenerAuxiliares(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_OBTENER_AUXILIARES", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, auxiliares)
}

func (h *BienesInsumosHandler) ListarPrecios(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, errCodeIdProductoInvalido, msgIdProductoInvalido)
		return
	}

	precios, err := h.servicio.ListarPreciosPorPlan(c.Request.Context(), idProducto)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_PRECIOS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, precios)
}

func (h *BienesInsumosHandler) GuardarPrecios(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, errCodeIdProductoInvalido, msgIdProductoInvalido)
		return
	}

	var payload []domain.BienInsumoPrecioPlan
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "PAYLOAD_INVALIDO", err.Error())
		return
	}

	if err := h.servicio.GuardarPreciosPorPlan(c.Request.Context(), idProducto, payload); err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_GUARDAR_PRECIOS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"guardado": true})
}

func (h *BienesInsumosHandler) Eliminar(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, errCodeIdProductoInvalido, msgIdProductoInvalido)
		return
	}

	if err := h.servicio.Eliminar(c.Request.Context(), idProducto); err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_ELIMINAR_BIEN", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"eliminado": true})
}
