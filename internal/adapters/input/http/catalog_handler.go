package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

// CatalogHandler expone el puerto de entrada input.CatalogService.
type CatalogHandler struct {
	service input.CatalogService
}

// NewCatalogHandler inyecta el caso de uso de catálogos en el adaptador HTTP.
func NewCatalogHandler(service input.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

// ListEtnias maneja GET /api/v1/etnias.
//
// @Summary Lista etnias
// @Description Devuelve el catálogo de etnias (SP ups_go_ListarEtnias).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.Etnia}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /etnias [get]
func (h *CatalogHandler) ListEtnias(c *gin.Context) {
	items, err := h.service.ListEtnias(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toEtniaResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListIdiomas maneja GET /api/v1/idiomas.
//
// @Summary Lista idiomas
// @Description Devuelve el catálogo de idiomas (SP ups_go_ListarIdiomas).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.Idioma}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /idiomas [get]
func (h *CatalogHandler) ListIdiomas(c *gin.Context) {
	items, err := h.service.ListIdiomas(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toIdiomaResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListTipoSexos maneja GET /api/v1/tipos-sexo.
//
// @Summary Lista tipos de sexo
// @Description Devuelve el catálogo de sexos (SP usp_go_ListarTiposSexos).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.TipoSexo}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /tipos-sexo [get]
func (h *CatalogHandler) ListTipoSexos(c *gin.Context) {
	items, err := h.service.ListTipoSexos(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toTipoSexoResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListEstadosCivil maneja GET /api/v1/estados-civil.
//
// @Summary Lista estados civiles
// @Description Devuelve el catálogo de estados civiles (SP usp_go_ListarEstadosCivil).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.TipoEstadoCivil}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /estados-civil [get]
func (h *CatalogHandler) ListEstadosCivil(c *gin.Context) {
	items, err := h.service.ListEstadosCivil(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toTipoEstadoCivilResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListGradosInstruccion maneja GET /api/v1/grados-instruccion.
//
// @Summary Lista grados de instrucción
// @Description Devuelve el catálogo de grados de instrucción (SP usp_go_ListarGradoInstruccion).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.TipoGradoInstruccion}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /grados-instruccion [get]
func (h *CatalogHandler) ListGradosInstruccion(c *gin.Context) {
	items, err := h.service.ListGradosInstruccion(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toTipoGradoInstruccionResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListOcupaciones maneja GET /api/v1/ocupaciones.
//
// @Summary Lista ocupaciones
// @Description Devuelve el catálogo de ocupaciones (SP usp_go_ListarOcupaciones).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.TipoOcupacion}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /ocupaciones [get]
func (h *CatalogHandler) ListOcupaciones(c *gin.Context) {
	items, err := h.service.ListOcupaciones(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toTipoOcupacionResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListTiposDocumentos maneja GET /api/v1/tipos-documentos.
//
// @Summary Lista tipos de documento de identidad
// @Description Devuelve el catálogo de tipos de documento (SP usp_go_ListarTiposDocumentos).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.TipoDocumento}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /tipos-documentos [get]
func (h *CatalogHandler) ListTiposDocumentos(c *gin.Context) {
	items, err := h.service.ListTiposDocumentos(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toTipoDocumentoResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListDepartamentos maneja GET /api/v1/departamentos.
//
// @Summary Lista departamentos
// @Description Devuelve el catálogo de departamentos (SP usp_go_ListarDepartamentos).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.Departamento}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /departamentos [get]
func (h *CatalogHandler) ListDepartamentos(c *gin.Context) {
	items, err := h.service.ListDepartamentos(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toDepartamentoResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListProvincias maneja GET /api/v1/provincias/:idDepartamento.
//
// @Summary Lista provincias por departamento
// @Description Devuelve las provincias de un departamento (SP usp_go_ListarProvincias).
// @Tags Catalogos
// @Produce json
// @Param idDepartamento path int true "Id del departamento"
// @Success 200 {object} apiResponse{data=[]domain.Provincia}
// @Failure 400 {object} apiResponse{error=apiError} "Id inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /provincias/{idDepartamento} [get]
func (h *CatalogHandler) ListProvincias(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	items, err := h.service.ListProvincias(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toProvinciaResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListDistritos maneja GET /api/v1/distritos/:idProvincia.
//
// @Summary Lista distritos por provincia
// @Description Devuelve los distritos de una provincia (SP usp_go_ListarDistritos).
// @Tags Catalogos
// @Produce json
// @Param idProvincia path int true "Id de la provincia"
// @Success 200 {object} apiResponse{data=[]domain.Distrito}
// @Failure 400 {object} apiResponse{error=apiError} "Id inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /distritos/{idProvincia} [get]
func (h *CatalogHandler) ListDistritos(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	items, err := h.service.ListDistritos(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toDistritoResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListCentrosPoblados maneja GET /api/v1/centros-poblados/:idDistrito.
//
// @Summary Lista centros poblados por distrito
// @Description Devuelve los centros poblados de un distrito (SP usp_go_ListarCentrosPoblados).
// @Tags Catalogos
// @Produce json
// @Param idDistrito path int true "Id del distrito"
// @Success 200 {object} apiResponse{data=[]domain.CentroPoblado}
// @Failure 400 {object} apiResponse{error=apiError} "Id inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /centros-poblados/{idDistrito} [get]
func (h *CatalogHandler) ListCentrosPoblados(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	items, err := h.service.ListCentrosPoblados(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toCentroPobladoResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListPaises maneja GET /api/v1/paises.
//
// @Summary Lista países
// @Description Devuelve el catálogo de países (SP usp_go_ListarPaises).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.Pais}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /paises [get]
func (h *CatalogHandler) ListPaises(c *gin.Context) {
	items, err := h.service.ListPaises(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toPaisResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListEstadosLlegoPaciente maneja GET /api/v1/estados-llego-paciente.
//
// @Summary Lista estados de llegada del paciente
// @Description Devuelve el catálogo de estados de llegada del paciente (SP usp_go_listarEstadosLlegoPaciente).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.EstadoLlegoPaciente}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /estados-llego-paciente [get]
func (h *CatalogHandler) ListEstadosLlegoPaciente(c *gin.Context) {
	items, err := h.service.ListEstadosLlegoPaciente(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toEstadoLlegoPacienteResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListFuentesFinanciamiento maneja GET /api/v1/fuentes-financiamiento.
//
// @Summary Lista fuentes de financiamiento
// @Description Devuelve el catálogo de fuentes de financiamiento (SP usp_go_ListarFuentesFinanciamiento).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.FuenteFinanciamiento}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /fuentes-financiamiento [get]
func (h *CatalogHandler) ListFuentesFinanciamiento(c *gin.Context) {
	items, err := h.service.ListFuentesFinanciamiento(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toFuenteFinanciamientoResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// ListServicios maneja GET /api/v1/servicios/:idTipoServicio.
//
// @Summary Lista servicios por tipo de servicio
// @Description Devuelve los servicios de un tipo (SP usp_go_ListarServicios).
// @Tags Catalogos
// @Produce json
// @Param idTipoServicio path int true "Id del tipo de servicio"
// @Success 200 {object} apiResponse{data=[]domain.Servicio}
// @Failure 400 {object} apiResponse{error=apiError} "Id inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /servicios/{idTipoServicio} [get]
func (h *CatalogHandler) ListServicios(c *gin.Context) {
	idTipoServicio, err := strconv.ParseInt(c.Param("idTipoServicio"), 10, 64)
	if err != nil || idTipoServicio <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "id inválido")
		return
	}
	items, err := h.service.ListServicios(c.Request.Context(), idTipoServicio)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toServicioResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// GetDatosInstitucion maneja GET /api/v1/datos-institucion.
//
// @Summary Obtiene los datos de la institución
// @Description Devuelve los datos de la institución (SP webParametrosDatosInstitucion).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=domain.DatosInstitucion}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /datos-institucion [get]
func (h *CatalogHandler) GetDatosInstitucion(c *gin.Context) {
	item, err := h.service.GetDatosInstitucion(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, toDatosInstitucionResponse(*item))
}

// ListEspecialidades maneja GET /api/v1/especialidades.
//
// @Summary Lista especialidades
// @Description Devuelve el catálogo de especialidades (SP usp_go_ListarEspecialidades).
// @Tags Catalogos
// @Produce json
// @Success 200 {object} apiResponse{data=[]domain.Especialidad}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /especialidades [get]
func (h *CatalogHandler) ListEspecialidades(c *gin.Context) {
	items, err := h.service.ListEspecialidades(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	dtoItems := make([]any, 0, len(items))
	for _, it := range items {
		dtoItems = append(dtoItems, toEspecialidadResponse(it))
	}
	respondSuccess(c, http.StatusOK, dtoItems)
}

// GetParametro maneja GET /api/v1/parametros/:idParametro.
//
// @Summary Obtiene un parámetro por id
// @Description Devuelve los valores de un parámetro (SP usp_go_webParametroSeleccionarPorId).
// @Tags Catalogos
// @Produce json
// @Param idParametro path int true "Id del parámetro"
// @Success 200 {object} apiResponse{data=domain.Parametro}
// @Failure 400 {object} apiResponse{error=apiError} "Id inválido"
// @Failure 404 {object} apiResponse{error=apiError} "Parámetro no encontrado"
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /parametros/{idParametro} [get]
func (h *CatalogHandler) GetParametro(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("idParametro"), 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "id inválido")
		return
	}
	item, err := h.service.GetParametro(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if item == nil {
		respondError(c, http.StatusNotFound, "PARAMETRO_NOT_FOUND", "parámetro no encontrado")
		return
	}
	respondSuccess(c, http.StatusOK, toParametroResponse(*item))
}

func (h *CatalogHandler) ListRecetaFrecuencias(c *gin.Context) {
	items, err := h.service.ListRecetaFrecuencias(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if items == nil {
		items = make([]domain.CatalogItem, 0)
	}
	respondSuccess(c, http.StatusOK, items)
}

func (h *CatalogHandler) ListRecetaUnidadesDosis(c *gin.Context) {
	items, err := h.service.ListRecetaUnidadesDosis(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if items == nil {
		items = make([]domain.CatalogItem, 0)
	}
	respondSuccess(c, http.StatusOK, items)
}

func (h *CatalogHandler) ListRecetaViasAdministracion(c *gin.Context) {
	items, err := h.service.ListRecetaViasAdministracion(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if items == nil {
		items = make([]domain.CatalogItem, 0)
	}
	respondSuccess(c, http.StatusOK, items)
}

func (h *CatalogHandler) BuscarMedicamentosReceta(c *gin.Context) {
	filtro := c.Query("q")
	idPaciente, _ := strconv.Atoi(c.DefaultQuery("idPaciente", "908637"))
	items, err := h.service.BuscarMedicamentosReceta(c.Request.Context(), filtro, idPaciente)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if items == nil {
		items = make([]domain.MedicamentoBusqueda, 0)
	}
	respondSuccess(c, http.StatusOK, items)
}

// parseID extrae y valida el path param "id" como entero.
func parseID(c *gin.Context) (int64, bool) {
	raw := c.Param("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "id inválido")
		return 0, false
	}
	return id, true
}

// Referencia en blanco para que swag resuelva los tipos de dominio en los
// comentarios de los handlers.
var _ = domain.Etnia{}
