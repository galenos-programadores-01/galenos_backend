package httpadapter

import (
	"net/http"

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
// @Description Retorna la lista de espera de cirugias filtrando por fecha y nombre del paciente
// @Tags ListaEsperaQx
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fecha query string true "Fecha de filtro (YYYY-MM-DD)"
// @Param paciente query string false "Nombre del paciente"
// @Success 200 {object} httpadapter.apiResponse{data=[]domain.ListaEsperaQx}
// @Router /api/v1/lista-espera-qx [get]
func (h *ListaEsperaQxHandler) HandleListar(c *gin.Context) {
	fecha := c.Query("fecha")
	paciente := c.Query("paciente")

	lista, err := h.service.Listar(c.Request.Context(), fecha, paciente)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "LISTA_ESPERA_QX_ERR", "Error obteniendo lista de espera quirurgica")
		return
	}

	respondSuccess(c, http.StatusOK, lista)
}

type crearListaEsperaQxRequest struct {
	IdPaciente       int    `json:"idPaciente"`
	IdMedico         int    `json:"idMedico"`
	IdTipoDocumento  int    `json:"idTipoDocumento" binding:"required"`
	NroDocumento     string `json:"nroDocumento" binding:"required"`
	ApellidoPaterno  string `json:"apellidoPaterno" binding:"required"`
	ApellidoMaterno  string `json:"apellidoMaterno"`
	PrimerNombre     string `json:"primerNombre" binding:"required"`
	SegundoNombre    string `json:"segundoNombre"`
	FechaNacimiento  string `json:"fechaNacimiento" binding:"required"`
	IdSexo           int    `json:"idSexo" binding:"required"`
	Telefono         string `json:"telefono"`
	Direccion        string `json:"direccion"`
	FechaOrden       string `json:"fechaOrden" binding:"required"`
	Diagnostico      string `json:"diagnostico"`
	FechaLaboratorio string `json:"fechaLaboratorio"`
	FechaICCardio    string `json:"fechaICCardio"`
	FechaICNeumo     string `json:"fechaICNeumo"`
	FechaICAnestesio string `json:"fechaICAnestesio"`
	Medico           string `json:"medico"`
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
// @Router /api/v1/lista-espera-qx [post]
func (h *ListaEsperaQxHandler) HandleCrear(c *gin.Context) {
	var req crearListaEsperaQxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	item := domain.ListaEsperaQxCrear{
		IdPaciente:       req.IdPaciente,
		IdMedico:         req.IdMedico,
		IdTipoDocumento:  req.IdTipoDocumento,
		NroDocumento:     req.NroDocumento,
		ApellidoPaterno:  req.ApellidoPaterno,
		ApellidoMaterno:  req.ApellidoMaterno,
		PrimerNombre:     req.PrimerNombre,
		SegundoNombre:    req.SegundoNombre,
		FechaNacimiento:  req.FechaNacimiento,
		IdSexo:           req.IdSexo,
		Telefono:         req.Telefono,
		Direccion:        req.Direccion,
		FechaOrden:       req.FechaOrden,
		Diagnostico:      req.Diagnostico,
		FechaLaboratorio: req.FechaLaboratorio,
		FechaICCardio:    req.FechaICCardio,
		FechaICNeumo:     req.FechaICNeumo,
		FechaICAnestesio: req.FechaICAnestesio,
		Medico:           req.Medico,
		Observacion:      req.Observacion,
	}

	if err := h.service.Crear(c.Request.Context(), item); err != nil {
		respondError(c, http.StatusInternalServerError, "LISTA_ESPERA_QX_SAVE_ERR", "Error guardando lista de espera quirurgica")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Paciente registrado en lista de espera quirurgica correctamente"})
}
