package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type interconsultaDxRequest struct {
	IdDiagnostico        int `json:"idDiagnostico" binding:"required"`
	IdSubclasificacionDx int `json:"idSubclasificacionDx" binding:"required"`
}

type interconsultaRegistrarRequest struct {
	IdEpisodio         int                      `json:"idEpisodio" binding:"required"`
	IdPrimeraAtencion  int                      `json:"idPrimeraAtencion" binding:"required"`
	IdServicio         int                      `json:"idServicio" binding:"required"`
	IdEspecialidad     int                      `json:"idEspecialidad" binding:"required"`
	IdPrioridad        int                      `json:"idPrioridad" binding:"required"`
	IdMotivoParametro  int                      `json:"idMotivoParametro" binding:"required"`
	MotivoOrden        *int                     `json:"motivoOrden" binding:"required"` // 0..2; puntero: gin required no traga el 0
	MotivoClinico      string                   `json:"motivoClinico" binding:"required"`
	IdEmpleado         int                      `json:"idEmpleado" binding:"required"`
	Diagnosticos       []interconsultaDxRequest `json:"diagnosticos" binding:"required,min=1"`
}

type interconsultaListadoResponse struct {
	IdInterconsulta  int    `json:"idInterconsulta"`
	IdAtencionOrigen int    `json:"idAtencionOrigen"`
	IdEspecialidad   int    `json:"idEspecialidad"`
	IdMedicoDestino  int    `json:"idMedicoDestino"`
	Motivo           string `json:"motivo"`
	FechaSolicitud   string `json:"fechaSolicitud"`
	Estado           string `json:"estado"`
	Especialidad     string `json:"especialidad"`
	ResumenHC        string `json:"resumenHC"`
}

// RegistrarSolicitudInterconsulta POST /api/v1/primera-atencion/:idAtencion/interconsultas
func (h *PrimeraAtencionHandler) RegistrarSolicitudInterconsulta(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion invalido")
		return
	}

	var req interconsultaRegistrarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if req.MotivoOrden == nil || *req.MotivoOrden < 0 || *req.MotivoOrden > 2 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "motivoOrden debe ser 0, 1 o 2")
		return
	}
	if req.IdPrimeraAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion es obligatorio")
		return
	}

	dx := make([]domain.InterconsultaDiagnosticoItem, 0, len(req.Diagnosticos))
	for _, d := range req.Diagnosticos {
		dx = append(dx, domain.InterconsultaDiagnosticoItem{
			IdDiagnostico:        d.IdDiagnostico,
			IdSubclasificacionDx: d.IdSubclasificacionDx,
		})
	}

	idIC, err := h.service.RegistrarSolicitudInterconsulta(c.Request.Context(), domain.InterconsultaRegistrarSolicitud{
		IdAtencion:         idAtencion,
		IdEpisodio:         req.IdEpisodio,
		IdPrimeraAtencion:  req.IdPrimeraAtencion,
		IdServicio:         req.IdServicio,
		IdEspecialidad:     req.IdEspecialidad,
		IdPrioridad:        req.IdPrioridad,
		IdMotivoParametro:  req.IdMotivoParametro,
		MotivoOrden:       *req.MotivoOrden,
		MotivoClinico:      req.MotivoClinico,
		IdEmpleado:         req.IdEmpleado,
		Diagnosticos:       dx,
	})
	if err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_INTERC_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idInterconsulta": idIC,
		"mensaje":         "Solicitud de interconsulta registrada",
	})
}

// ListarInterconsultas GET .../interconsultas?idEpisodio=&idPrimeraAtencion=
func (h *PrimeraAtencionHandler) ListarInterconsultas(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion invalido")
		return
	}
	idEpisodio, err := strconv.Atoi(c.Query("idEpisodio"))
	if err != nil || idEpisodio <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idEpisodio es obligatorio")
		return
	}
	idPrimeraAtencion, err := strconv.Atoi(c.Query("idPrimeraAtencion"))
	if err != nil || idPrimeraAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion es obligatorio")
		return
	}

	items, err := h.service.ListarInterconsultas(
		c.Request.Context(), idAtencion, idEpisodio, idPrimeraAtencion,
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_INTERC_LIST_FAILED", err.Error())
		return
	}

	resp := make([]interconsultaListadoResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, interconsultaListadoResponse{
			IdInterconsulta:  it.IdInterconsulta,
			IdAtencionOrigen: it.IdAtencionOrigen,
			IdEspecialidad:   it.IdEspecialidad,
			IdMedicoDestino:  it.IdMedicoDestino,
			Motivo:           it.Motivo,
			FechaSolicitud:   it.FechaSolicitud,
			Estado:           it.Estado,
			Especialidad:     it.Especialidad,
			ResumenHC:        it.ResumenHC,
		})
	}
	respondSuccess(c, http.StatusOK, resp)
}

// EliminarSolicitudInterconsulta DELETE .../interconsultas/:idInterconsulta?idEpisodio=
func (h *PrimeraAtencionHandler) EliminarSolicitudInterconsulta(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion invalido")
		return
	}
	idInterconsulta, err := strconv.Atoi(c.Param("idInterconsulta"))
	if err != nil || idInterconsulta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idInterconsulta invalido")
		return
	}
	idEpisodio, err := strconv.Atoi(c.Query("idEpisodio"))
	if err != nil || idEpisodio <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idEpisodio es obligatorio")
		return
	}

	if err := h.service.EliminarSolicitudInterconsulta(
		c.Request.Context(), idAtencion, idInterconsulta, idEpisodio,
	); err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_INTERC_DELETE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idInterconsulta": idInterconsulta,
		"mensaje":         "Interconsulta eliminada",
	})
}
