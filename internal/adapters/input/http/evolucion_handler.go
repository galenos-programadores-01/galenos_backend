package httpadapter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type EvolucionHandler struct {
	service input.EvolucionService
}

func NewEvolucionHandler(service input.EvolucionService) *EvolucionHandler {
	return &EvolucionHandler{service: service}
}

// @Summary List patients for the medical evolution tray
// @Description Returns a list of patients with an attention within the date range (Atenciones.FechaIngreso). If fini/ffin are omitted, returns the most recent 50.
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fini query string false "Fecha inicial (YYYY-MM-DD)"
// @Param ffin query string false "Fecha final (YYYY-MM-DD)"
// @Router /api/v1/evoluciones/pacientes [get]
func (h *EvolucionHandler) HandleListPatients(c *gin.Context) {
	fini := c.Query("fini")
	if _, err := time.Parse("2006-01-02", fini); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_FINI", "El parámetro fini (YYYY-MM-DD) es requerido")
		return
	}
	if ffin := c.Query("ffin"); ffin != "" {
		if _, err := time.Parse("2006-01-02", ffin); err != nil {
			respondError(c, http.StatusBadRequest, "INVALID_FFIN", "El parámetro ffin (YYYY-MM-DD) es inválido")
			return
		}
	}

	idUsuario := 1
	if val, exists := c.Get("idEmpleado"); exists {
		switch id := val.(type) {
		case float64:
			idUsuario = int(id)
		case string:
			if parsed, err := strconv.Atoi(id); err == nil {
				idUsuario = parsed
			}
		case int:
			idUsuario = id
		}
	}

	patients, err := h.service.GetPatientTray(c.Request.Context(), fini, c.Query("ffin"), idUsuario)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "EVOL_PATIENTS_ERR", "Error obteniendo pacientes")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "7"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 7
	}

	totalItems := len(patients)
	totalPages := 1
	if totalItems > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}
	inicio := (page - 1) * pageSize
	if inicio > totalItems {
		inicio = totalItems
	}
	fin := inicio + pageSize
	if fin > totalItems {
		fin = totalItems
	}

	respondSuccess(c, http.StatusOK, map[string]any{
		"page":       page,
		"pageSize":   pageSize,
		"totalItems": totalItems,
		"totalPages": totalPages,
		"items":      patients[inicio:fin],
	})
}

// @Summary Get evolutions for a patient
// @Description Returns the saved medical evolutions for a given registration ID
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pacienteId path int true "ID of the Registration / Encounter"
// @Router /api/v1/evoluciones/paciente/{pacienteId} [get]
func (h *EvolucionHandler) HandleListEvoluciones(c *gin.Context) {
	idRegAtencionStr := c.Param("pacienteId")
	idRegAtencion, err := strconv.Atoi(idRegAtencionStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de atención inválido")
		return
	}

	evolutions, err := h.service.GetEvoluciones(c.Request.Context(), idRegAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "EVOL_GET_ERR", "Error obteniendo evoluciones")
		return
	}

	respondSuccess(c, http.StatusOK, evolutions)
}

type SaveEvolucionRequest struct {
	IdRegAtencion int    `json:"idRegAtencion" binding:"required"`
	DataB64       string `json:"dataB64" binding:"required"`
}

// @Summary Save an evolution
// @Description Saves a new medical evolution
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SaveEvolucionRequest true "Evolution data"
// @Router /api/v1/evoluciones [post]
func (h *EvolucionHandler) HandleCreateEvolucion(c *gin.Context) {
	var req SaveEvolucionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	idEmpleado := c.GetInt("idEmpleado")
	if idEmpleado == 0 {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "No se identificó al empleado en el token")
		return
	}

	err := h.service.SaveEvolucion(c.Request.Context(), req.IdRegAtencion, idEmpleado, req.DataB64)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "EVOL_SAVE_ERR", "Error guardando la evolución")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{
		"message":   "Evolución guardada correctamente",
		"ipCliente": c.ClientIP(),
		"fecha":     time.Now().Format("2006-01-02"),
		"hora":      time.Now().Format("15:04:05"),
	})
}

// @Summary Get evolutions for the main tray using usp_go_EvolucionesMedicas_Bandeja
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fechaInicio query string false "Fecha inicial (YYYY-MM-DD)"
// @Param fechaFin query string false "Fecha final (YYYY-MM-DD)"
// @Param filtro query string false "Texto de búsqueda (DNI, Nombre, NroAtencion)"
// @Router /api/v1/evoluciones/bandeja [get]
func (h *EvolucionHandler) HandleBandeja(c *gin.Context) {
	fechaInicio := c.Query("fechaInicio")
	fechaFin := c.Query("fechaFin")
	filtro := c.Query("filtro")

	list, err := h.service.GetBandeja(c.Request.Context(), fechaInicio, fechaFin, filtro)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "BANDEJA_ERR", "Error consultando la bandeja de evoluciones: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, list)
}

// @Summary Insert a complete medical evolution using usp_go_EvolucionesMedicas_Insertar
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /api/v1/evoluciones/registro [post]
func (h *EvolucionHandler) HandleInsertEvolucionMedica(c *gin.Context) {
	var req domain.EvolucionMedicaInsert
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido: "+err.Error())
		return
	}

	idEmpleado := c.GetInt("idEmpleado")
	if idEmpleado != 0 {
		req.UsuarioCreacion = idEmpleado
	}

	idEvolucion, mensaje, err := h.service.InsertEvolucionMedica(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "EVOL_INSERT_ERR", "Error registrando evolución médica: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]any{
		"idEvolucion": idEvolucion,
		"mensaje":     mensaje,
	})
}
