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

	if payload.IdUsuario <= 0 {
		if idEmpleado := c.GetInt("idEmpleado"); idEmpleado > 0 {
			payload.IdUsuario = idEmpleado
		} else if claims, ok := c.Get("auth_claims"); ok {
			if authClaims, ok := claims.(domain.AuthClaims); ok && authClaims.IdEmpleado > 0 {
				payload.IdUsuario = authClaims.IdEmpleado
			}
		} else {
			payload.IdUsuario = 1
		}
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

func (h *AtencionCEHandler) HandleListarHistorialAtenciones(c *gin.Context) {
	idParam := c.Param("idPaciente")
	idPaciente, err := strconv.Atoi(idParam)
	if err != nil || idPaciente <= 0 {
		respondError(c, http.StatusBadRequest, "ID_INVALIDO", "El identificador de paciente no es valido")
		return
	}

	historial, err := h.servicio.ListarHistorialAtenciones(c.Request.Context(), idPaciente)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_HISTORIAL_ATENCIONES", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, historial)
}

func (h *AtencionCEHandler) HandleObtenerHistorialAtencionDetalle(c *gin.Context) {
	idParam := c.Param("idAtencion")
	idAtencion, err := strconv.Atoi(idParam)
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "ID_INVALIDO", "El identificador de atencion no es valido")
		return
	}

	detalle, err := h.servicio.ObtenerHistorialAtencionDetalle(c.Request.Context(), idAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_DETALLE_HISTORIAL", err.Error())
		return
	}
	if detalle == nil {
		respondError(c, http.StatusNotFound, "NO_ENCONTRADO", "No se encontro la atencion especificada")
		return
	}

	respondSuccess(c, http.StatusOK, detalle)
}

func (h *AtencionCEHandler) HandleListarReportesOperatorios(c *gin.Context) {
	idParam := c.Param("idPaciente")
	idPaciente, err := strconv.Atoi(idParam)
	if err != nil || idPaciente <= 0 {
		respondError(c, http.StatusBadRequest, "ID_INVALIDO", "El identificador de paciente no es valido")
		return
	}

	reportes, err := h.servicio.ListarReportesOperatorios(c.Request.Context(), idPaciente)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_REPORTES_OPERATORIOS", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, reportes)
}

func (h *AtencionCEHandler) HandleListarProcedimientosHistorial(c *gin.Context) {
	idParam := c.Param("idPaciente")
	idPaciente, err := strconv.Atoi(idParam)
	if err != nil || idPaciente <= 0 {
		respondError(c, http.StatusBadRequest, "ID_INVALIDO", "El identificador de paciente no es valido")
		return
	}

	procedimientos, err := h.servicio.ListarProcedimientosHistorial(c.Request.Context(), idPaciente)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_PROCEDIMIENTOS_HISTORIAL", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, procedimientos)
}

func (h *AtencionCEHandler) HandleListarHistorialTriajePaciente(c *gin.Context) {
	idParam := c.Param("idPaciente")
	idPaciente, err := strconv.Atoi(idParam)
	if err != nil || idPaciente <= 0 {
		respondError(c, http.StatusBadRequest, "ID_INVALIDO", "El identificador de paciente no es valido")
		return
	}

	triajes, err := h.servicio.ListarHistorialTriajePaciente(c.Request.Context(), idPaciente)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_HISTORIAL_TRIAJE", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, triajes)
}

func (h *AtencionCEHandler) ListarServiciosMedico(c *gin.Context) {
	idEmpleadoParam := c.Query("idEmpleado")
	var idEmpleado int
	if idEmpleadoParam != "" {
		idEmpleado, _ = strconv.Atoi(idEmpleadoParam)
	} else if claims, ok := c.Get("auth_claims"); ok {
		if authClaims, ok := claims.(domain.AuthClaims); ok && authClaims.IdEmpleado > 0 {
			idEmpleado = authClaims.IdEmpleado
		}
	}
	if idEmpleado == 0 {
		idEmpleado = -100
	}

	idTipoServicioParam := c.DefaultQuery("idTipoServicio", "1")
	idTipoServicio, err := strconv.Atoi(idTipoServicioParam)
	if err != nil || idTipoServicio <= 0 {
		idTipoServicio = 1
	}

	servicios, err := h.servicio.ListarServiciosPorMedico(c.Request.Context(), idTipoServicio, idEmpleado)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERROR_SERVICIOS_MEDICO", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, servicios)
}

