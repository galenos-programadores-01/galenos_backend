package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

// HospitalizacionHandler expone el puerto de entrada
// input.HospitalizacionService.
type HospitalizacionHandler struct {
	service input.HospitalizacionService
}

// NewHospitalizacionHandler inyecta el caso de uso de admisión de
// hospitalización en el adaptador HTTP.
func NewHospitalizacionHandler(service input.HospitalizacionService) *HospitalizacionHandler {
	return &HospitalizacionHandler{service: service}
}

// CreateAdmision maneja POST /api/v1/hospitalizacion/admision.
//
// @Summary Admisiona un paciente en hospitalización
// @Description Registra la admisión de un paciente en hospitalización invocando el SP usp_go_AdmisionHospitalizacion y devuelve el @Resultado del SP. El idEmpleado se toma del token JWT.
// @Tags Hospitalizacion
// @Accept json
// @Produce json
// @Param admision body createAdmisionHospitalizacionRequest true "Datos de la admisión"
// @Success 200 {object} apiResponse{data=object}
// @Failure 400 {object} apiResponse{error=apiError} "Payload inválido"
// @Failure 401 {object} apiResponse{error=apiError} "No autenticado"
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Security BearerAuth
// @Router /hospitalizacion/admision [post]
func (h *HospitalizacionHandler) CreateAdmision(c *gin.Context) {
	var req createAdmisionHospitalizacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	idEmpleado := c.GetInt("idEmpleado")
	if idEmpleado == 0 {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "empleado no autenticado")
		return
	}

	domainObj := req.toDomain()
	empID := int64(idEmpleado)
	domainObj.IDEmpleado = &empID

	result, err := h.service.Admisionar(c.Request.Context(), domainObj)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ADMISION_HOSPITALIZACION_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}

// ObtenerReporteTicket maneja
// GET /api/v1/hospitalizacion/reporte-ticket/:idCuenta.
//
// @Summary Datos del ticket de admisión de hospitalización
// @Description Devuelve los datos para imprimir el ticket de admisión de hospitalización de una cuenta de atención (SP usp_go_GenererReporteTicketHospi).
// @Tags Hospitalizacion
// @Accept json
// @Produce json
// @Param idCuenta path int true "Id de la cuenta de atención"
// @Success 200 {object} apiResponse{data=domain.ReporteTicketHospitalizacion}
// @Failure 400 {object} apiResponse{error=apiError} "Id inválido"
// @Failure 401 {object} apiResponse{error=apiError} "No autenticado"
// @Failure 404 {object} apiResponse{error=apiError} "Cuenta no encontrada"
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Security BearerAuth
// @Router /hospitalizacion/reporte-ticket/{idCuenta} [get]
func (h *HospitalizacionHandler) ObtenerReporteTicket(c *gin.Context) {
	raw := c.Param("idCuenta")
	idCuenta, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || idCuenta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuenta debe ser un entero positivo")
		return
	}

	reporte, err := h.service.ObtenerReporteTicket(c.Request.Context(), idCuenta)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "REPORTE_TICKET_HOSPITALIZACION_FAILED", err.Error())
		return
	}
	if reporte == nil {
		respondError(c, http.StatusNotFound, "REPORTE_TICKET_NOT_FOUND", "no se encontró la cuenta de atención")
		return
	}

	respondSuccess(c, http.StatusOK, reporte)
}
