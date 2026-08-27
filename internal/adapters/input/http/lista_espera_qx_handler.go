package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type ListaEsperaQxHandler struct {
	service input.ListaEsperaQxService
}

func NewListaEsperaQxHandler(service input.ListaEsperaQxService) *ListaEsperaQxHandler {
	return &ListaEsperaQxHandler{service: service}
}

// @Summary Listar lista de espera quirurgica
// @Description Retorna la lista de espera de cirugias filtrando por rango de fechas, nombre del paciente y especialidad
// @Tags ListaEsperaQx
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fecha query string true "Fecha inicio del filtro (YYYY-MM-DD)"
// @Param fechaFin query string true "Fecha fin del filtro (YYYY-MM-DD)"
// @Param paciente query string false "Nombre del paciente"
// @Param idEspecialidad query int false "ID de especialidad"
// @Success 200 {object} httpadapter.apiResponse{data=[]domain.ListaEsperaQx}
// @Router /lista-espera-qx [get]
func (h *ListaEsperaQxHandler) HandleListar(c *gin.Context) {
	fecha := c.Query("fecha")
	fechaFin := c.Query("fechaFin")
	paciente := c.Query("paciente")
	idEspecialidad := 0
	if v := c.Query("idEspecialidad"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			idEspecialidad = parsed
		}
	}

	lista, err := h.service.Listar(c.Request.Context(), fecha, fechaFin, paciente, idEspecialidad)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "LISTA_ESPERA_QX_ERR", "Error obteniendo lista de espera quirurgica")
		return
	}

	respondSuccess(c, http.StatusOK, lista)
}

// @Summary Obtener lista de espera quirurgica por ID
// @Description Retorna los datos de un paciente en lista de espera quirurgica
// @Tags ListaEsperaQx
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del registro"
// @Success 200 {object} httpadapter.apiResponse{data=domain.ListaEsperaQxPaciente}
// @Router /lista-espera-qx/{id} [get]
func (h *ListaEsperaQxHandler) HandleObtenerPorId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	item, err := h.service.ObtenerPorId(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "LISTA_ESPERA_QX_GET_ERR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, item)
}

type crearListaEsperaQxRequest struct {
	IdPaciente       int    `json:"idPaciente" binding:"required"`
	IdMedico         int    `json:"idMedico" binding:"required"`
	FechaOrden       string `json:"fechaOrden" binding:"required"`
	Diagnostico      int    `json:"diagnostico"`
	IdEspecialidad   int    `json:"idEspecialidad"`
	FechaLaboratorio string `json:"fechaLaboratorio"`
	FechaICCardio    string `json:"fechaICCardio"`
	FechaICNeumo     string `json:"fechaICNeumo"`
	FechaICAnestesio string `json:"fechaICAnestesio"`
	Observacion      string `json:"observacion"`
}

// @Summary Crear lista de espera quirurgica
// @Description Registra un nuevo paciente en la lista de espera quirurgica
// @Tags ListaEsperaQx
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body crearListaEsperaQxRequest true "Datos del paciente"
// @Success 200 {object} httpadapter.apiResponse
// @Router /lista-espera-qx [post]
func (h *ListaEsperaQxHandler) HandleCrear(c *gin.Context) {
	var req crearListaEsperaQxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	idEmpleado := 0
	if val, exists := c.Get("idEmpleado"); exists {
		switch v := val.(type) {
		case int:
			idEmpleado = v
		case int64:
			idEmpleado = int(v)
		case float64:
			idEmpleado = int(v)
		}
	}

	item := domain.ListaEsperaQxCrear{
		IdPaciente:       req.IdPaciente,
		IdMedico:         req.IdMedico,
		FechaOrden:       req.FechaOrden,
		Diagnostico:      req.Diagnostico,
		IdEspecialidad:   req.IdEspecialidad,
		FechaLab:         req.FechaLaboratorio,
		FechaICCardio:    req.FechaICCardio,
		FechaICNeumo:     req.FechaICNeumo,
		FechaICAnestesio: req.FechaICAnestesio,
		Observacion:      req.Observacion,
	}

	if err := h.service.Crear(c.Request.Context(), item, idEmpleado); err != nil {
		respondError(c, http.StatusInternalServerError, "LISTA_ESPERA_QX_SAVE_ERR", "Error guardando lista de espera quirurgica")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Paciente registrado en lista de espera quirurgica correctamente"})
}

type modificarListaEsperaQxRequest struct {
	FechaOrden       string `json:"fechaOrden"`
	Diagnostico      int    `json:"diagnostico"`
	IdEspecialidad   int    `json:"idEspecialidad"`
	FechaLaboratorio string `json:"fechaLaboratorio"`
	FechaICCardio    string `json:"fechaICCardio"`
	FechaICNeumo     string `json:"fechaICNeumo"`
	FechaICAnestesio string `json:"fechaICAnestesio"`
	Observacion      string `json:"observacion"`
}

// @Summary Modificar lista de espera quirurgica
// @Description Actualiza los datos de un paciente en lista de espera quirurgica
// @Tags ListaEsperaQx
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del registro"
// @Param request body modificarListaEsperaQxRequest true "Datos a actualizar"
// @Success 200 {object} httpadapter.apiResponse
// @Router /lista-espera-qx/{id} [put]
func (h *ListaEsperaQxHandler) HandleModificar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	var req modificarListaEsperaQxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	item := domain.ListaEsperaQxModificar{
		Id:               id,
		FechaOrden:       req.FechaOrden,
		Diagnostico:      req.Diagnostico,
		IdEspecialidad:   req.IdEspecialidad,
		FechaLab:         req.FechaLaboratorio,
		FechaICCardio:    req.FechaICCardio,
		FechaICNeumo:     req.FechaICNeumo,
		FechaICAnestesio: req.FechaICAnestesio,
		Observacion:      req.Observacion,
	}

	if err := h.service.Modificar(c.Request.Context(), item); err != nil {
		respondError(c, http.StatusInternalServerError, "LISTA_ESPERA_QX_UPDATE_ERR", "Error actualizando lista de espera quirurgica")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Paciente actualizado en lista de espera quirurgica correctamente"})
}

// @Summary Reporte lista de espera quirurgica
// @Description Retorna el reporte de lista de espera quirurgica filtrando por rango de fechas y especialidad
// @Tags ListaEsperaQx
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fecha query string true "Fecha inicio del filtro (YYYY-MM-DD)"
// @Param fechaFin query string true "Fecha fin del filtro (YYYY-MM-DD)"
// @Param idEspecialidad query int false "ID de especialidad"
// @Success 200 {object} httpadapter.apiResponse{data=[]domain.ListaEsperaQxReporte}
// @Router /lista-espera-qx/reporte [get]
func (h *ListaEsperaQxHandler) HandleReporte(c *gin.Context) {
	fecha := c.Query("fecha")
	fechaFin := c.Query("fechaFin")
	idEspecialidad := 0
	if v := c.Query("idEspecialidad"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			idEspecialidad = parsed
		}
	}

	lista, err := h.service.Reporte(c.Request.Context(), fecha, fechaFin, idEspecialidad)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "LISTA_ESPERA_QX_REP_ERR", "Error generando reporte de lista espera quirurgica")
		return
	}

	respondSuccess(c, http.StatusOK, lista)
}
