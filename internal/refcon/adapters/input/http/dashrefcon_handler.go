package refconhttp

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
	"github.com/galenos-pro/appointments-api/internal/refcon/ports/input"
	"github.com/gin-gonic/gin"
)

type RefConHandler struct {
	service input.RefConService
}

func NewRefConHandler(service input.RefConService) *RefConHandler {
	return &RefConHandler{service: service}
}

// @Summary Listar referencias por rango de fechas
// @Description Retorna la cantidad de referencias/contrarreferencias agrupadas por mes y estado para el dashboard de barras (SP usp_go_ListarReferenciasPorMesyAnio)
// @Tags DashRefCon
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fechaInicio query string true "Fecha de inicio (yyyy-mm-dd)"
// @Param fechaFin query string true "Fecha de fin (yyyy-mm-dd)"
// @Param estado query int false "Estado de la referencia (0 = todos; 3 ACEPTADO, 4 RECHAZADO, 5 PACIENTE RECIBIDO, 7 PACIENTE CITADO, 8 CONTRAREFERIDO)"
// @Success 200 {object} apiResponse{data=[]domain.ReferenciaPorMes}
// @Failure 400 {object} apiResponse
// @Router /dashrefcon/referencias [get]
func (h *RefConHandler) HandleListarReferencias(c *gin.Context) {
	fechaInicio, errInicio := queryDateParam(c, "fechaInicio")
	if errInicio != nil {
		respondError(c, http.StatusBadRequest, "INVALID_FECHA_INICIO", "El parámetro fechaInicio es obligatorio y debe tener formato yyyy-mm-dd")
		return
	}
	fechaFin, errFin := queryDateParam(c, "fechaFin")
	if errFin != nil {
		respondError(c, http.StatusBadRequest, "INVALID_FECHA_FIN", "El parámetro fechaFin es obligatorio y debe tener formato yyyy-mm-dd")
		return
	}
	if fechaInicio.After(fechaFin) {
		respondError(c, http.StatusBadRequest, "INVALID_RANGO_FECHAS", "La fecha de inicio no puede ser mayor que la fecha de fin")
		return
	}
	estado := queryIntParam(c, "estado", 0)

	resultado, err := h.service.ListarReferenciasPorMes(c.Request.Context(), fechaInicio, fechaFin, estado)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "DASHREFCON_REF_ERR", "Error obteniendo las referencias del dashboard")
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}

// @Summary Listar UPS
// @Description Retorna el catálogo de Unidades Productoras de Servicios (UPS)
// @Tags DashRefCon
// @Produce json
// @Security BearerAuth
// @Success 200 {object} apiResponse{data=[]domain.Ups}
// @Router /dashrefcon/ups [get]
func (h *RefConHandler) HandleListarUps(c *gin.Context) {
	ups, err := h.service.ListarUps(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "DASHREFCON_UPS_ERR", "Error obteniendo las UPS")
		return
	}

	respondSuccess(c, http.StatusOK, ups)
}

// @Summary Listar todas las UPS
// @Description Retorna el catálogo completo de Unidades Productoras de Servicios (UPS) con código y descripción desde la tabla SuSalud_ups
// @Tags RefCon
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.Ups}
// @Router /refcon/ups [get]
func (h *RefConHandler) HandleListarTodasLasUps(c *gin.Context) {
	ups, err := h.service.ListarTodasLasUps(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFCON_UPS_ERR", "Error obteniendo las UPS")
		return
	}

	respondSuccess(c, http.StatusOK, ups)
}

// @Summary Listar establecimientos
// @Description Retorna el catálogo de establecimientos de salud con código IPRESS y nombre desde la tabla Establecimientos
// @Tags RefCon
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.Establecimiento}
// @Router /refcon/establecimientos [get]
func (h *RefConHandler) HandleListarEstablecimientos(c *gin.Context) {
	establecimientos, err := h.service.ListarEstablecimientos(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFCON_ESTABLECIMIENTOS_ERR", "Error obteniendo los establecimientos")
		return
	}

	respondSuccess(c, http.StatusOK, establecimientos)
}

// @Summary Listar distritos por IdReniec
// @Description Retorna el distrito asociado a un código RENIEC (tabla Distritos)
// @Tags RefCon
// @Produce json
// @Param idReniec path int true "IdReniec del distrito"
// @Success 200 {object} apiResponse{data=[]domain.DistritoReniec}
// @Router /refcon/distritos-reniec/{idReniec} [get]
func (h *RefConHandler) HandleListarDistritosPorIdReniec(c *gin.Context) {
	idReniec, err := strconv.Atoi(c.Param("idReniec"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID_RENIEC", "El IdReniec debe ser un número entero")
		return
	}

	distritos, err := h.service.ListarDistritosPorIdReniec(c.Request.Context(), idReniec)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REFCON_DISTRITOS_RENIEC_ERR", "Error obteniendo los distritos por IdReniec")
		return
	}

	respondSuccess(c, http.StatusOK, distritos)
}

// @Summary Listar UPS de un establecimiento en MINSA
// @Description Consulta el listado de Unidades Productoras de Servicios (UPS) del establecimiento indicado por su código RENIPRESS en el servicio REST del MINSA
// @Tags DashRefCon
// @Produce json
// @Security BearerAuth
// @Param codigoRenipress path string true "Código RENIPRESS del establecimiento"
// @Success 200 {object} apiResponse{data=domain.ListadoUpsResponse}
// @Router /dashrefcon/upss/{codigoRenipress} [get]
func (h *RefConHandler) HandleListadoUpss(c *gin.Context) {
	codigoRenipress := c.Param("codigoRenipress")
	if codigoRenipress == "" {
		respondError(c, http.StatusBadRequest, "INVALID_PARAMS", "El parámetro codigoRenipress es obligatorio")
		return
	}

	resultado, err := h.service.ListarUpssMinsa(c.Request.Context(), codigoRenipress)
	if err != nil {
		log.Printf("[DashRefCon] Error consultando UPS en MINSA (codigo=%s): %v", codigoRenipress, err)
		respondError(c, http.StatusBadGateway, "MINSA_REFCON_ERR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}

// @Summary Listar especialidades vigentes en MINSA
// @Description Consulta el listado de especialidades vigentes en el servicio REST de interoperabilidad del MINSA (sin parámetros)
// @Tags DashRefCon
// @Produce json
// @Security BearerAuth
// @Success 200 {object} apiResponse{data=domain.ListadoEspecialidadesResponse}
// @Router /dashrefcon/especialidades-minsa [get]
func (h *RefConHandler) HandleListadoEspecialidades(c *gin.Context) {
	resultado, err := h.service.ListarEspecialidadesMinsa(c.Request.Context())
	if err != nil {
		log.Printf("[DashRefCon] Error consultando especialidades en MINSA: %v", err)
		respondError(c, http.StatusBadGateway, "MINSA_REFCON_ERR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}

// @Summary Enviar (registrar) una referencia en MINSA
// @Description Registra una referencia en el servicio saveReferencia del MINSA enviando la estructura completa de la referencia
// @Tags DashRefCon
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body domain.SaveReferenciaRequest true "Estructura de la referencia a registrar"
// @Success 200 {object} apiResponse{data=domain.SaveReferenciaResponse}
// @Failure 400 {object} apiResponse
// @Failure 502 {object} apiResponse
// @Router /dashrefcon/enviar-referencia [post]
func (h *RefConHandler) HandleSaveReferencia(c *gin.Context) {
	var req domain.SaveReferenciaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "El cuerpo debe ser JSON con la estructura de la referencia")
		return
	}

	resultado, err := h.service.GuardarReferenciaMinsa(c.Request.Context(), req)
	if err != nil {
		log.Printf("[DashRefCon] Error registrando referencia en MINSA: %v", err)
		respondError(c, http.StatusBadGateway, "MINSA_REFCON_ERR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, resultado)
}

func queryIntParam(c *gin.Context, name string, fallback int) int {
	if v := c.Query(name); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}

// queryDateParam lee un parámetro de fecha (yyyy-mm-dd). Devuelve error cuando
// el parámetro está vacío o no tiene un formato válido.
func queryDateParam(c *gin.Context, name string) (time.Time, error) {
	v := c.Query(name)
	if v == "" {
		return time.Time{}, fmt.Errorf("parámetro %s vacío", name)
	}
	for _, layout := range []string{"2006-01-02", "2006-1-2"} {
		if fecha, err := time.Parse(layout, v); err == nil {
			return fecha, nil
		}
	}
	return time.Time{}, fmt.Errorf("parámetro %s con formato inválido", name)
}
