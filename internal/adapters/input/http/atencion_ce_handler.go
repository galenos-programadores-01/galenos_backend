package httpadapter

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type AtencionCEHandler struct {
	servicio input.AtencionCEService
}

func NewAtencionCEHandler(servicio input.AtencionCEService) *AtencionCEHandler {
	return &AtencionCEHandler{
		servicio: servicio,
	}
}

func (h *AtencionCEHandler) ListarCola(c *gin.Context) {
	idMedico, _ := strconv.Atoi(c.DefaultQuery("idMedico", "0"))
	idServicio, _ := strconv.Atoi(c.DefaultQuery("idServicio", "0"))
	fecha := c.Query("fecha")
	filtro := c.DefaultQuery("filtro", c.Query("q"))

	filtros := domain.ColaFiltrosCE{
		IdMedico:   idMedico,
		IdServicio: idServicio,
		Fecha:      fecha,
		Filtro:     filtro,
	}

	cola, err := h.servicio.ListarColaPacientes(c.Request.Context(), filtros)
	if err != nil {
		if errors.Is(c.Request.Context().Err(), context.Canceled) {
			return
		}
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_COLA_CE", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, cola)
}

func (h *AtencionCEHandler) ObtenerFicha(c *gin.Context) {
	idParam := c.Param("id")
	idAtencion, err := strconv.Atoi(idParam)
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "ID_INVALIDO", "El identificador de atencion no es valido")
		return
	}

	ficha, err := h.servicio.ObtenerFichaAtencion(c.Request.Context(), idAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_OBTENER_FICHA_CE", err.Error())
		return
	}
	if ficha == nil {
		respondError(c, http.StatusNotFound, "NO_ENCONTRADO", "No se encontro la atencion especificada")
		return
	}

	respondSuccess(c, http.StatusOK, ficha)
}

func (h *AtencionCEHandler) GuardarFicha(c *gin.Context) {
	var payload domain.AtencionCEGuardarPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "PAYLOAD_INVALIDO", err.Error())
		return
	}

	idEmpleado := c.GetInt("idEmpleado")
	if idEmpleado > 0 && payload.IdUsuario == 0 {
		payload.IdUsuario = idEmpleado
	}

	if err := h.servicio.GuardarFichaAtencion(c.Request.Context(), payload); err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_GUARDAR_FICHA_CE", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"guardado": true})
}

func (h *AtencionCEHandler) ListarDestinos(c *gin.Context) {
	destinos, err := h.servicio.ListarDestinosAtencion(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_DESTINOS_CE", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, destinos)
}

func (h *AtencionCEHandler) ListarCatalogoLabs(c *gin.Context) {
	labs, err := h.servicio.ListarCatalogoLabsHIS(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_LISTAR_LABS_HIS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, labs)
}
