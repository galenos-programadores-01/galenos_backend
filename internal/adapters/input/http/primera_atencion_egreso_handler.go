package httpadapter

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/adapters/output/persistence/sqlserver"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type altaDomicilioRequest struct {
	IdEpisodio       int     `json:"idEpisodio" binding:"required"`
	IdMedicoEgreso   int     `json:"idMedicoEgreso" binding:"required"`
	FechaEgreso      string  `json:"fechaEgreso" binding:"required"`
	IdTipoAlta       int     `json:"idTipoAlta" binding:"required"`
	IdCondicionAlta  int     `json:"idCondicionAlta" binding:"required"`
	IdPronostico     *int    `json:"idPronostico"`
	Recomendacion    string  `json:"recomendacion"`
	IdServicioEgreso int     `json:"idServicioEgreso" binding:"required"`
	IdDestinoAtencion *int   `json:"idDestinoAtencion"`
}

type altaHospitalizacionRequest struct {
	IdEpisodio        int    `json:"idEpisodio" binding:"required"`
	IdMedicoEgreso    int    `json:"idMedicoEgreso" binding:"required"`
	FechaEgreso       string `json:"fechaEgreso" binding:"required"`
	IdServicioEgreso  int    `json:"idServicioEgreso" binding:"required"` // servicio que recibe
	IdDestinoAtencion int    `json:"idDestinoAtencion" binding:"required"` // catálogo tarjeta HOS
	MotivoTraslado    string `json:"motivoTraslado"`
}

type altaFugaRequest struct {
	IdEpisodio         int    `json:"idEpisodio" binding:"required"`
	IdMedicoEgreso     int    `json:"idMedicoEgreso" binding:"required"`
	FechaEgreso        string `json:"fechaEgreso" binding:"required"` // fecha/hora de la fuga
	FechaFuga          string `json:"fechaFuga"`                      // último contacto (opcional)
	IdServicioEgreso   int    `json:"idServicioEgreso" binding:"required"`
	IdDestinoAtencion  int    `json:"idDestinoAtencion" binding:"required"` // catálogo FUG
	ReportaFuga        string `json:"reportaFuga"`
	IdNotificacionFuga *int   `json:"idNotificacionFuga"` // 1|2|3
	Observaciones      string `json:"observaciones"`
}

type altaFallecidoRequest struct {
	IdEpisodio             int    `json:"idEpisodio" binding:"required"`
	IdMedicoEgreso         int    `json:"idMedicoEgreso" binding:"required"`
	FechaEgreso            string `json:"fechaEgreso" binding:"required"` // fecha/hora fallecimiento
	IdServicioEgreso       int    `json:"idServicioEgreso" binding:"required"` // bandeja
	IdDestinoAtencion      int    `json:"idDestinoAtencion" binding:"required"` // catálogo FAL
	IdDestinoServicio      *int   `json:"idDestinoServicio"`                   // destino cadáver (grupo 11)
	IdTipoCertificado      *int   `json:"idTipoCertificado"`                   // grupo 12
	ActDefuncion           string `json:"actDefuncion"`                        // N° certificado
	IdtiempoEstanciaDeseso *int   `json:"idtiempoEstanciaDeseso"`              // grupo 13
	RequiereAutopsia       *bool  `json:"requiereAutopsia"`                    // necropsia
	Observaciones          string `json:"observaciones"`                       // familiar informado
}

type altaHospitalariaResponse struct {
	Resultado int    `json:"resultado"`
	Mensaje   string `json:"mensaje"`
}

// RegistrarAltaDomicilio POST .../egreso/alta-domicilio
func (h *PrimeraAtencionHandler) RegistrarAltaDomicilio(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req altaDomicilioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	fecha, err := sqlserver.ParseFechaHoraLocal(req.FechaEgreso)
	if err != nil || fecha == nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fechaEgreso inválida")
		return
	}

	idEpisodio := req.IdEpisodio
	idTipo := req.IdTipoAlta
	idCond := req.IdCondicionAlta
	idUsuario := idUsuarioDesdeJWT(c)
	var idUsuarioPtr *int
	if idUsuario > 0 {
		idUsuarioPtr = &idUsuario
	}

	res, err := h.service.RegistrarAltaHospitalaria(c.Request.Context(), domain.AltaHospitalariaRegistrar{
		IdAtencion:        idAtencion,
		IdEpisodio:        &idEpisodio,
		IdDestinoAtencion: req.IdDestinoAtencion,
		IdDestinoServicio: nil,
		IdMedicoEgreso:    req.IdMedicoEgreso,
		FechaEgreso:       *fecha,
		IdCondicionAlta:   &idCond,
		IdTipoAlta:        &idTipo,
		IdPronostico:      req.IdPronostico,
		Recomendacion:     strings.TrimSpace(req.Recomendacion),
		IdServicioEgreso:  req.IdServicioEgreso,
		IdUsuario:         idUsuarioPtr,
		ActDefuncion:      "",
		MotivoTraslado:    "",
	})
	if err != nil {
		respondPrimeraAtencionWriteError(c, "PRIMERA_ATENCION_ALTA_DOMICILIO_FAILED", err)
		return
	}

	respondSuccess(c, http.StatusOK, altaHospitalariaResponse{
		Resultado: res.Resultado,
		Mensaje:   res.Mensaje,
	})
}

// RegistrarAltaHospitalizacion POST .../egreso/alta-hospitalizacion
// IdCondicionAlta / IdTipoAlta / IdDestinoServicio = NULL; MotivoTraslado se envía al SP.
func (h *PrimeraAtencionHandler) RegistrarAltaHospitalizacion(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req altaHospitalizacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	fecha, err := sqlserver.ParseFechaHoraLocal(req.FechaEgreso)
	if err != nil || fecha == nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fechaEgreso inválida")
		return
	}

	idEpisodio := req.IdEpisodio
	idDestino := req.IdDestinoAtencion
	idUsuario := idUsuarioDesdeJWT(c)
	var idUsuarioPtr *int
	if idUsuario > 0 {
		idUsuarioPtr = &idUsuario
	}

	res, err := h.service.RegistrarAltaHospitalaria(c.Request.Context(), domain.AltaHospitalariaRegistrar{
		IdAtencion:        idAtencion,
		IdEpisodio:        &idEpisodio,
		IdDestinoAtencion: &idDestino,
		IdDestinoServicio: nil,
		IdMedicoEgreso:    req.IdMedicoEgreso,
		FechaEgreso:       *fecha,
		IdCondicionAlta:   nil,
		IdTipoAlta:        nil,
		IdPronostico:      nil,
		Recomendacion:     "",
		IdServicioEgreso:  req.IdServicioEgreso,
		IdUsuario:         idUsuarioPtr,
		ActDefuncion:      "",
		MotivoTraslado:    strings.TrimSpace(req.MotivoTraslado),
	})
	if err != nil {
		respondPrimeraAtencionWriteError(c, "PRIMERA_ATENCION_ALTA_HOSP_FAILED", err)
		return
	}

	respondSuccess(c, http.StatusOK, altaHospitalariaResponse{
		Resultado: res.Resultado,
		Mensaje:   res.Mensaje,
	})
}

// RegistrarAltaFuga POST .../egreso/alta-fuga
// FechaEgreso = fecha/hora fuga; FechaFuga = último contacto; IdNotificacionFuga 1|2|3.
func (h *PrimeraAtencionHandler) RegistrarAltaFuga(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req altaFugaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	fechaEgreso, err := sqlserver.ParseFechaHoraLocal(req.FechaEgreso)
	if err != nil || fechaEgreso == nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fechaEgreso inválida")
		return
	}

	var fechaFuga *time.Time
	if strings.TrimSpace(req.FechaFuga) != "" {
		ff, errFf := sqlserver.ParseFechaHoraLocal(req.FechaFuga)
		if errFf != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fechaFuga inválida")
			return
		}
		fechaFuga = ff
	}

	idEpisodio := req.IdEpisodio
	idDestino := req.IdDestinoAtencion
	idUsuario := idUsuarioDesdeJWT(c)
	var idUsuarioPtr *int
	if idUsuario > 0 {
		idUsuarioPtr = &idUsuario
	}

	res, err := h.service.RegistrarAltaHospitalaria(c.Request.Context(), domain.AltaHospitalariaRegistrar{
		IdAtencion:         idAtencion,
		IdEpisodio:         &idEpisodio,
		IdDestinoAtencion:  &idDestino,
		IdDestinoServicio:  nil,
		IdMedicoEgreso:     req.IdMedicoEgreso,
		FechaEgreso:        *fechaEgreso,
		IdCondicionAlta:    nil,
		IdTipoAlta:         nil,
		IdPronostico:       nil,
		Recomendacion:      "",
		IdServicioEgreso:   req.IdServicioEgreso,
		IdUsuario:          idUsuarioPtr,
		ActDefuncion:       "",
		MotivoTraslado:     "",
		FechaFuga:          fechaFuga,
		ReportaFuga:        strings.TrimSpace(req.ReportaFuga),
		IdNotificacionFuga: req.IdNotificacionFuga,
		Observaciones:      strings.TrimSpace(req.Observaciones),
	})
	if err != nil {
		respondPrimeraAtencionWriteError(c, "PRIMERA_ATENCION_ALTA_FUGA_FAILED", err)
		return
	}

	respondSuccess(c, http.StatusOK, altaHospitalariaResponse{
		Resultado: res.Resultado,
		Mensaje:   res.Mensaje,
	})
}

// RegistrarAltaFallecido POST .../egreso/alta-fallecido
// IdDestinoServicio = destino cadáver; Observaciones = familiar informado;
// ActDefuncion = N° certificado; DX en FE con clasif. 4.
func (h *PrimeraAtencionHandler) RegistrarAltaFallecido(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req altaFallecidoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	fechaEgreso, err := sqlserver.ParseFechaHoraLocal(req.FechaEgreso)
	if err != nil || fechaEgreso == nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fechaEgreso inválida")
		return
	}

	idEpisodio := req.IdEpisodio
	idDestino := req.IdDestinoAtencion
	idUsuario := idUsuarioDesdeJWT(c)
	var idUsuarioPtr *int
	if idUsuario > 0 {
		idUsuarioPtr = &idUsuario
	}

	res, err := h.service.RegistrarAltaHospitalaria(c.Request.Context(), domain.AltaHospitalariaRegistrar{
		IdAtencion:             idAtencion,
		IdEpisodio:             &idEpisodio,
		IdDestinoAtencion:      &idDestino,
		IdDestinoServicio:      req.IdDestinoServicio,
		IdMedicoEgreso:         req.IdMedicoEgreso,
		FechaEgreso:            *fechaEgreso,
		IdCondicionAlta:        nil,
		IdTipoAlta:             nil,
		IdPronostico:           nil,
		Recomendacion:          "",
		IdServicioEgreso:       req.IdServicioEgreso,
		IdUsuario:              idUsuarioPtr,
		IdTipoCertificado:      req.IdTipoCertificado,
		ActDefuncion:           strings.TrimSpace(req.ActDefuncion),
		MotivoTraslado:         "",
		Observaciones:          strings.TrimSpace(req.Observaciones),
		IdtiempoEstanciaDeseso: req.IdtiempoEstanciaDeseso,
		RequiereAutopsia:       req.RequiereAutopsia,
	})
	if err != nil {
		respondPrimeraAtencionWriteError(c, "PRIMERA_ATENCION_ALTA_FALLECIDO_FAILED", err)
		return
	}

	respondSuccess(c, http.StatusOK, altaHospitalariaResponse{
		Resultado: res.Resultado,
		Mensaje:   res.Mensaje,
	})
}
