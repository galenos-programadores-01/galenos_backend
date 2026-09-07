package httpadapter

import (
	"log"
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type DiagnosticoHandler struct {
	useCase input.DiagnosticoUseCase
}

func NewDiagnosticoHandler(useCase input.DiagnosticoUseCase) *DiagnosticoHandler {
	return &DiagnosticoHandler{useCase: useCase}
}

// SearchDiagnosticos maneja GET /api/v1/diagnosticos/search.
// @Summary Buscar diagnósticos
// @Description Busca diagnósticos en base a un texto (filtro), idAtencion e idPaciente usando el SP usp_go_SelectDiagnosticos
// @Tags Diagnosticos
// @Accept json
// @Produce json
// @Param filtro query string false "Texto a buscar"
// @Param idAtencion query int false "ID de Atención"
// @Param idPaciente query int false "ID de Paciente"
// @Success 200 {object} apiResponse{data=[]domain.DiagnosticoBusqueda}
// @Failure 400 {object} apiResponse{error=apiError}
// @Failure 500 {object} apiResponse{error=apiError}
// @Router /diagnosticos/search [get]
func (h *DiagnosticoHandler) SearchDiagnosticos(c *gin.Context) {
	filtro := c.Query("filtro")
	idAtencionStr := c.Query("idAtencion")
	idPacienteStr := c.Query("idPaciente")

	idAtencion, _ := strconv.Atoi(idAtencionStr)
	idPaciente, _ := strconv.Atoi(idPacienteStr)

	log.Printf("Buscando diagnosticos con: filtro=%q, idAtencion=%d, idPaciente=%d", filtro, idAtencion, idPaciente)

	results, err := h.useCase.SearchDiagnosticos(c.Request.Context(), filtro, idAtencion, idPaciente)
	if err != nil {
		log.Printf("Internal error searching diagnosticos: %v", err)
		respondError(c, http.StatusInternalServerError, "ERR_SEARCH_DIAG", "Error buscando diagnósticos: "+err.Error())
		return
	}

	if results == nil {
		results = make([]domain.DiagnosticoBusqueda, 0)
	}

	respondSuccess(c, http.StatusOK, results)
}

// @Summary Listar diagnosticos CIE10
// @Description Busca diagnosticos por filtro (codigo o descripcion)
// @Tags Diagnosticos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param filtro query string true "Filtro de busqueda"
// @Success 200 {object} httpadapter.apiResponse{data=[]domain.DiagnosticoSimple}
// @Router /diagnosticos/listar [get]
func (h *DiagnosticoHandler) HandleListarDiagnosticos(c *gin.Context) {
	filtro := c.Query("filtro")

	results, err := h.useCase.ListarDiagnosticos(c.Request.Context(), filtro)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERR_LIST_DIAG", "Error listando diagnósticos")
		return
	}

	if results == nil {
		results = make([]domain.DiagnosticoSimple, 0)
	}

	respondSuccess(c, http.StatusOK, results)
}

// HandleObtenerDiagnosticosAtencion maneja GET /api/v1/diagnosticos/atencion/:idAtencion
// @Summary Obtener diagnósticos de una atención
// @Description Obtiene los diagnósticos de una atención médica utilizando usp_go_ObtenerDiagnosticosAtencion
// @Tags Diagnosticos
// @Accept json
// @Produce json
// @Param idAtencion path int true "ID de Atención"
// @Param idPrimeraAtencion query int false "ID de Primera Atención"
// @Param idEvolucion query int false "ID de Evolución"
// @Success 200 {object} httpadapter.apiResponse{data=[]domain.DiagnosticoAtencion}
// @Router /diagnosticos/atencion/{idAtencion} [get]
func (h *DiagnosticoHandler) HandleObtenerDiagnosticosAtencion(c *gin.Context) {
	idAtencionStr := c.Param("idAtencion")
	idAtencion, err := strconv.Atoi(idAtencionStr)
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "idAtencion inválido")
		return
	}

	var idPrimPtr *int
	if idPrimStr := c.Query("idPrimeraAtencion"); idPrimStr != "" {
		if idPrim, err := strconv.Atoi(idPrimStr); err == nil {
			idPrimPtr = &idPrim
		}
	}

	var idEvPtr *int
	if idEvStr := c.Query("idEvolucion"); idEvStr != "" {
		if idEv, err := strconv.Atoi(idEvStr); err == nil {
			idEvPtr = &idEv
		}
	}

	results, err := h.useCase.ObtenerDiagnosticosAtencion(c.Request.Context(), idAtencion, idPrimPtr, idEvPtr)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERR_GET_DIAG_ATENCION", "Error obteniendo diagnósticos de la atención: "+err.Error())
		return
	}

	if results == nil {
		results = make([]domain.DiagnosticoAtencion, 0)
	}

	respondSuccess(c, http.StatusOK, results)
}

// HandleAgregarDiagnosticoAtencion maneja POST /api/v1/diagnosticos/atencion
// @Summary Agregar diagnóstico a una atención médica
// @Description Registra un diagnóstico para una atención médica utilizando usp_go_AtencionesDiagnosticosAgregar
// @Tags Diagnosticos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.AgregarDiagnosticoAtencionRequest true "Datos del diagnóstico a registrar"
// @Success 200 {object} httpadapter.apiResponse{data=domain.AgregarDiagnosticoAtencionResponse}
// @Failure 400 {object} httpadapter.apiResponse{error=httpadapter.apiError}
// @Failure 500 {object} httpadapter.apiResponse{error=httpadapter.apiError}
// @Router /diagnosticos/atencion [post]
func (h *DiagnosticoHandler) HandleAgregarDiagnosticoAtencion(c *gin.Context) {
	var req domain.AgregarDiagnosticoAtencionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Datos de entrada inválidos: "+err.Error())
		return
	}

	res, err := h.useCase.AgregarDiagnosticoAtencion(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERR_ADD_DIAG", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, res)
}
