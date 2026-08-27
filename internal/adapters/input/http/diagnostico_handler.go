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
