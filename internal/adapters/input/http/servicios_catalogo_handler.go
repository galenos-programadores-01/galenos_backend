package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

const (
	errCodeIdServicioInvalido = "ID_INVALIDO"
	msgIdServicioInvalido     = "El identificador de servicio no es válido"
)

type ServiciosCatalogoHandler struct {
	servicio input.ServiciosCatalogoService
}

func NewServiciosCatalogoHandler(servicio input.ServiciosCatalogoService) *ServiciosCatalogoHandler {
	return &ServiciosCatalogoHandler{servicio: servicio}
}

func (h *ServiciosCatalogoHandler) Listar(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	idGrupo, _ := strconv.Atoi(c.DefaultQuery("idGrupo", "0"))
	filtro := c.Query("filtro")

	filtros := domain.ServicioCatalogoFiltros{
		Filtro:          filtro,
		IdServicioGrupo: idGrupo,
		Page:            page,
		PageSize:        pageSize,
	}

	if cptQuery := c.Query("cpt"); cptQuery != "" {
		esCpt := cptQuery == "true" || cptQuery == "1"
		filtros.SoloCpt = &esCpt
	}

	if qxQuery := c.Query("quirurgico"); qxQuery != "" {
		esQx := qxQuery == "true" || qxQuery == "1"
		filtros.SoloQuirurgico = &esQx
	}

	if estadoQuery := c.Query("estado"); estadoQuery != "" {
		if estadoVal, err := strconv.Atoi(estadoQuery); err == nil {
			filtros.Estado = &estadoVal
		}
	}

	resultado, err := h.servicio.Listar(c.Request.Context(), filtros)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_SERVICIOS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}

func (h *ServiciosCatalogoHandler) ObtenerPorID(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, errCodeIdServicioInvalido, msgIdServicioInvalido)
		return
	}

	detalle, err := h.servicio.ObtenerPorID(c.Request.Context(), idProducto)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_OBTENER_SERVICIO", err.Error())
		return
	}
	if detalle == nil {
		respondError(c, http.StatusNotFound, "NO_ENCONTRADO", "No se encontró el servicio solicitado")
		return
	}

	respondSuccess(c, http.StatusOK, detalle)
}

func (h *ServiciosCatalogoHandler) Crear(c *gin.Context) {
	var payload domain.ServicioCatalogoGuardarPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "PAYLOAD_INVALIDO", err.Error())
		return
	}

	nuevoID, err := h.servicio.Crear(c.Request.Context(), payload)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_CREAR_SERVICIO", err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, gin.H{"idProducto": nuevoID})
}

func (h *ServiciosCatalogoHandler) Actualizar(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, errCodeIdServicioInvalido, msgIdServicioInvalido)
		return
	}

	var payload domain.ServicioCatalogoGuardarPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "PAYLOAD_INVALIDO", err.Error())
		return
	}

	if err := h.servicio.Actualizar(c.Request.Context(), idProducto, payload); err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_ACTUALIZAR_SERVICIO", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"actualizado": true})
}

func (h *ServiciosCatalogoHandler) ListarGrupos(c *gin.Context) {
	grupos, err := h.servicio.ListarGrupos(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_GRUPOS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, grupos)
}

func (h *ServiciosCatalogoHandler) Eliminar(c *gin.Context) {
	idParam := c.Param("id")
	idProducto, err := strconv.Atoi(idParam)
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, errCodeIdServicioInvalido, msgIdServicioInvalido)
		return
	}

	if err := h.servicio.Eliminar(c.Request.Context(), idProducto); err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_ELIMINAR_SERVICIO", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"eliminado": true})
}
