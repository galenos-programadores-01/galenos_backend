package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type RefConHandler struct {
	service input.RefConService
}

func NewRefConHandler(service input.RefConService) *RefConHandler {
	return &RefConHandler{service: service}
}

// @Summary Listar referencias por mes y año
// @Description Retorna la cantidad de referencias/contrarreferencias agrupadas por mes y estado para el dashboard de barras
// @Tags DashRefCon
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mes query int false "Mes (0 = todos los meses)"
// @Param anio query int true "Año de la referencia"
// @Param estado query int false "Estado de la referencia (3 ACEPTADO, 4 RECHAZADO, 5 PACIENTE RECIBIDO, 7 PACIENTE CITADO, 8 CONTRAREFERIDO)"
// @Param ups query string false "Código de la UPS"
// @Param opcion query int true "Opción del reporte (1 total, 2 por estado, 3 por estado filtrado, 4 por estado y UPS)"
// @Success 200 {object} apiResponse{data=[]domain.ReferenciaPorMes}
// @Router /dashrefcon/referencias [get]
func (h *RefConHandler) HandleListarReferencias(c *gin.Context) {
	mes := queryIntParam(c, "mes", 0)
	anio := queryIntParam(c, "anio", 0)
	estado := queryIntParam(c, "estado", 0)
	opcion := queryIntParam(c, "opcion", 1)
	ups := c.Query("ups")

	if anio <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ANIO", "El parámetro anio es obligatorio y debe ser mayor a 0")
		return
	}
	if opcion < 1 || opcion > 4 {
		respondError(c, http.StatusBadRequest, "INVALID_OPCION", "El parámetro opcion debe estar entre 1 y 4")
		return
	}
	if mes < 0 || mes > 12 {
		respondError(c, http.StatusBadRequest, "INVALID_MES", "El parámetro mes debe estar entre 0 y 12")
		return
	}

	resultado, err := h.service.ListarReferenciasPorMes(c.Request.Context(), mes, anio, estado, opcion, ups)
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

func queryIntParam(c *gin.Context, name string, fallback int) int {
	if v := c.Query(name); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}
