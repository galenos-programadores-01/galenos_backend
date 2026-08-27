package httpadapter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

const errIdCuentaAtencion = "idCuentaAtencion debe ser un entero positivo"

// SisHandler expone el puerto de entrada input.SisService.
type SisHandler struct {
	service input.SisService
}

// NewSisHandler inyecta el caso de uso de SIS en el adaptador HTTP.
func NewSisHandler(service input.SisService) *SisHandler {
	return &SisHandler{service: service}
}

// ConsultarAfiliado maneja GET /api/v1/sis/afiliado/:nrodoc.
//
// @Summary Consulta un afiliado en SIS por número de documento
// @Description Invoca el servicio SOAP de SIS para consultar la afiliación de una persona.
// @Tags SIS
// @Produce json
// @Param nrodoc path string true "Número de documento (DNI)"
// @Param strDisa query string false "Código DISA"
// @Param strTipoFormato query string false "Tipo de formato"
// @Param strNroContrato query string false "Número de contrato"
// @Param strCorrelativo query string false "Correlativo"
// @Param strTipoDocumento query int false "Tipo de documento (1=DNI, 3=CE)"
// @Param intOpcion query int false "Opción de consulta"
// @Success 200 {object} apiResponse{data=object} "Afiliado encontrado"
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 502 {object} apiResponse{error=apiError} "Servicio SIS no disponible"
// @Router /sis/afiliado/{nrodoc} [get]
func (h *SisHandler) ConsultarAfiliado(c *gin.Context) {
	params, err := parseSisAfiliadoParams(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_SIS_REQUEST", err.Error())
		return
	}

	result, err := h.service.ConsultarAfiliado(c.Request.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDocumentNumber),
			errors.Is(err, domain.ErrInvalidDocumentType):
			respondError(c, http.StatusBadRequest, "INVALID_SIS_REQUEST", err.Error())
		default:
			respondError(c, http.StatusBadGateway, "SIS_UNAVAILABLE", err.Error())
		}
		return
	}

	respondSuccess(c, http.StatusOK, toSisAfiliadoResponse(result))
}

// GestionarAfiliacion maneja POST /api/v1/sis/filiaciones.
//
// @Summary Registra o actualiza una afiliación SIS
// @Description Invoca el SP webSisFiliacionesGestionar para guardar los datos de afiliación de un paciente.
// @Tags SIS
// @Accept json
// @Produce json
// @Param request body sisAfiliacionRequest true "Datos de afiliación SIS"
// @Success 200 {object} apiResponse{data=object} "Afiliación registrada"
// @Failure 400 {object} apiResponse{error=apiError} "Cuerpo inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al guardar la afiliación"
// @Router /sis/filiaciones [post]
func (h *SisHandler) GestionarAfiliacion(c *gin.Context) {
	var req sisAfiliacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if idEmpleado := c.GetInt("idEmpleado"); idEmpleado != 0 {
		idEmp64 := int64(idEmpleado)
		req.IdUsuarioAuditoria = &idEmp64
	}

	if err := h.service.GestionarAfiliacion(c.Request.Context(), req.toDomain()); err != nil {
		respondError(c, http.StatusInternalServerError, "SIS_AFILIACION_SAVE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"estado": "ok"})
}

// BuscarAfiliacion maneja GET /api/v1/sis/filiaciones.
//
// @Summary Busca una afiliación SIS por DISA, TipoFormato, NroContrato, Correlativo y CodigoTabla
// @Description Consulta el servicio SIS (ConsultarAfiliadoFuaE) usando los parámetros de afiliación. El nrodoc no es obligatorio.
// @Tags SIS
// @Produce json
// @Param strDisa query string true "Código DISA (ej. 035)"
// @Param strTipoFormato query string true "Tipo formato / Lote (ej. E)"
// @Param strNroContrato query string true "Número de contrato (ej. 11427904)"
// @Param strCorrelativo query string false "Correlativo"
// @Param strCodigoTabla query string false "Código de tabla"
// @Success 200 {object} apiResponse{data=object} "Afiliado encontrado"
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 502 {object} apiResponse{error=apiError} "Servicio SIS no disponible"
// @Router /sis/filiaciones [get]
func (h *SisHandler) BuscarAfiliacion(c *gin.Context) {
	dis := c.Query("strDisa")
	tipoFormato := c.Query("strTipoFormato")
	nroContrato := c.Query("strNroContrato")
	codigo := c.Query("strCodTabla")
	if codigo == "" {
		codigo = c.Query("strCodigoTabla")
	}

	if dis == "" || tipoFormato == "" || nroContrato == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "strDisa, strTipoFormato y strNroContrato son requeridos")
		return
	}

	params := shared.SISAfiliadoParams{
		TipoDocumento: 1,
		Opcion:        1,
		Disa:          dis,
		TipoFormato:   tipoFormato,
		NroContrato:   nroContrato,
		Correlativo:   c.Query("strCorrelativo"),
		CodTabla:      codigo,
		Lote:          tipoFormato,
		Tabla:         codigo,
	}

	result, err := h.service.BuscarPorAfiliacion(c.Request.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDocumentNumber),
			errors.Is(err, domain.ErrInvalidDocumentType):
			respondError(c, http.StatusBadRequest, "INVALID_SIS_REQUEST", err.Error())
		default:
			respondError(c, http.StatusBadGateway, "SIS_UNAVAILABLE", err.Error())
		}
		return
	}

	respondSuccess(c, http.StatusOK, toSisAfiliadoResponse(result))
}

// ForzarGuardadoFua maneja POST /api/v1/sis/fua.
//
// @Summary Fuerza el guardado del FUA de una cuenta de atención
// @Description Invoca el SP webSisFuaAtencionForzarGuardado para forzar la generación/guardado del FUA.
// @Tags SIS
// @Accept json
// @Produce json
// @Param request body sisFuaForzarGuardadoRequest true "Id de la cuenta de atención"
// @Success 200 {object} apiResponse{data=object} "FUA guardado"
// @Failure 400 {object} apiResponse{error=apiError} "Cuerpo inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al guardar el FUA"
// @Router /sis/fua [post]
func (h *SisHandler) ForzarGuardadoFua(c *gin.Context) {
	var req sisFuaForzarGuardadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.service.ForzarGuardadoFua(c.Request.Context(), req.IdCuentaAtencion); err != nil {
		respondError(c, http.StatusInternalServerError, "SIS_FUA_SAVE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"estado": "ok"})
}

// AgregarFua maneja POST /api/v1/sis/fua/agregar.
//
// @Summary Agrega el FUA de una cuenta de atención
// @Description Invoca el SP usp_go_webFUAgregar para registrar el FUA. Devuelve el @Respuesta del procedimiento.
// @Tags SIS
// @Accept json
// @Produce json
// @Param request body sisFuaAgregarRequest true "Datos para agregar el FUA"
// @Success 200 {object} apiResponse{data=object} "Respuesta del SP"
// @Failure 400 {object} apiResponse{error=apiError} "Cuerpo inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al agregar el FUA"
// @Router /sis/fua/agregar [post]
func (h *SisHandler) AgregarFua(c *gin.Context) {
	var req sisFuaAgregarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := h.service.AgregarFua(c.Request.Context(), req.IdCuentaAtencion, req.IdEmpleado, req.NombrePc)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SIS_FUA_ADD_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"respuesta": result})
}

// GetFuaImprimir maneja GET /api/v1/sis/fua/imprimir.
//
// @Summary Consulta los datos del FUA para imprimir
// @Description Invoca el SP webFuaImprimirIdCuentaAtencion con el id de la cuenta de atención. Devuelve los datos del FUA.
// @Tags SIS
// @Produce json
// @Param idCuentaAtencion query int true "Id de la cuenta de atención"
// @Success 200 {object} apiResponse{data=object} "Datos del FUA"
// @Failure 400 {object} apiResponse{error=apiError} "Parámetro inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al consultar el FUA"
// @Router /sis/fua/imprimir [get]
func (h *SisHandler) GetFuaImprimir(c *gin.Context) {
	var params sisFuaImprimirParams
	if err := c.ShouldBindQuery(&params); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	data, err := h.service.GetFuaImprimir(c.Request.Context(), params.IdCuentaAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SIS_FUA_PRINT_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, data)
}

// ListDiagnosticos maneja GET /api/v1/sis/diagnosticos.
//
// @Summary Consulta los diagnósticos de una atención
// @Description Invoca el SP webAtencionesDiagnosticosIdCuentaAtencion con el id de la cuenta de atención. Devuelve los diagnósticos.
// @Tags SIS
// @Accept json
// @Produce json
// @Param idCuentaAtencion query int true "Id de la cuenta de atención"
// @Success 200 {object} apiResponse{data=[]object} "Lista de diagnósticos"
// @Failure 400 {object} apiResponse{error=apiError} "Error de validación"
// @Failure 500 {object} apiResponse{error=apiError} "Error al consultar los diagnósticos"
// @Router /sis/diagnosticos [get]
func (h *SisHandler) ListDiagnosticos(c *gin.Context) {
	var params sisDiagnosticosParams
	if err := c.ShouldBindQuery(&params); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if params.IdCuentaAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion debe ser un entero positivo")
		return
	}

	items, err := h.service.ListDiagnosticos(c.Request.Context(), params.IdCuentaAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SIS_DIAGNOSTICS_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// ListConsumo maneja GET /api/v1/sis/consumo.
//
// @Summary Lista el detalle de consumo de una atención
// @Description Invoca el SP webFactOrdenServicioDetaDesFinaListarIdcuenta con el id de la cuenta de atención.
// @Tags SIS
// @Produce json
// @Param idCuentaAtencion query int true "Id de la cuenta de atención"
// @Success 200 {object} apiResponse{data=[]object} "Lista de consumos"
// @Failure 400 {object} apiResponse{error=apiError} "Parámetro inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al listar el consumo"
// @Router /sis/consumo [get]
func (h *SisHandler) ListConsumo(c *gin.Context) {
	var params sisCuentaAtencionParams
	if err := c.ShouldBindQuery(&params); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if params.IdCuentaAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", errIdCuentaAtencion)
		return
	}

	items, err := h.service.ListConsumo(c.Request.Context(), params.IdCuentaAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SIS_CONSUMPTION_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// ListMedicamentos maneja GET /api/v1/sis/medicamentos.
//
// @Summary Lista los medicamentos de una atención
// @Description Invoca el SP webMedicamentosListarIdCuentaAtencion con el id de la cuenta de atención.
// @Tags SIS
// @Produce json
// @Param idCuentaAtencion query int true "Id de la cuenta de atención"
// @Success 200 {object} apiResponse{data=[]object} "Lista de medicamentos"
// @Failure 400 {object} apiResponse{error=apiError} "Parámetro inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al listar los medicamentos"
// @Router /sis/medicamentos [get]
func (h *SisHandler) ListMedicamentos(c *gin.Context) {
	var params sisCuentaAtencionParams
	if err := c.ShouldBindQuery(&params); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if params.IdCuentaAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", errIdCuentaAtencion)
		return
	}

	items, err := h.service.ListMedicamentos(c.Request.Context(), params.IdCuentaAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SIS_MEDICATIONS_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// ListProcedimientos maneja GET /api/v1/sis/procedimientos.
//
// @Summary Lista los procedimientos de una atención
// @Description Invoca el SP webProcedimientoListarIdCuentaAtencion con el id de la cuenta de atención.
// @Tags SIS
// @Produce json
// @Param idCuentaAtencion query int true "Id de la cuenta de atención"
// @Success 200 {object} apiResponse{data=[]object} "Lista de procedimientos"
// @Failure 400 {object} apiResponse{error=apiError} "Parámetro inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al listar los procedimientos"
// @Router /sis/procedimientos [get]
func (h *SisHandler) ListProcedimientos(c *gin.Context) {
	var params sisCuentaAtencionParams
	if err := c.ShouldBindQuery(&params); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if params.IdCuentaAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", errIdCuentaAtencion)
		return
	}

	items, err := h.service.ListProcedimientos(c.Request.Context(), params.IdCuentaAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "SIS_PROCEDURES_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

func parseSisAfiliadoParams(c *gin.Context) (shared.SISAfiliadoParams, error) {
	params := shared.SISAfiliadoParams{
		DocumentNumber: c.Param("nrodoc"),
		TipoDocumento:  1,
		Opcion:         1,
		Disa:           c.Query("strDisa"),
		TipoFormato:    c.Query("strTipoFormato"),
		NroContrato:    c.Query("strNroContrato"),
		Correlativo:    c.Query("strCorrelativo"),
	}
	if raw := c.Query("strTipoDocumento"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || (parsed != 1 && parsed != 3) {
			return params, errors.New("strTipoDocumento debe ser 1 (DNI) o 3 (Carnet de Extranjería)")
		}
		params.TipoDocumento = int(parsed)
	} else if raw := c.Query("tipoDocumento"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || (parsed != 1 && parsed != 3) {
			return params, errors.New("tipoDocumento debe ser 1 (DNI) o 3 (Carnet de Extranjería)")
		}
		params.TipoDocumento = int(parsed)
	}
	if raw := c.Query("intOpcion"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || parsed <= 0 {
			return params, errors.New("intOpcion debe ser un entero positivo")
		}
		params.Opcion = int(parsed)
	}
	return params, nil
}
