package httpadapter

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// PatientHandler expone el puerto de entrada input.PatientService.
type PatientHandler struct {
	service input.PatientService
}

// NewPatientHandler inyecta el caso de uso en el adaptador HTTP.
func NewPatientHandler(service input.PatientService) *PatientHandler {
	return &PatientHandler{service: service}
}

// List maneja GET /api/v1/pacientes?page=1&pageSize=20.
//
// @Summary Lista pacientes paginados
// @Description Devuelve una página de pacientes ordenada, con datos de paginación.
// @Tags Pacientes
// @Produce json
// @Param page query int false "Número de página (default 1)"
// @Param pageSize query int false "Tamaño de página (default 20, máximo 100)"
// @Success 200 {object} apiResponse{data=[]patientResponse}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /pacientes [get]
func (h *PatientHandler) List(c *gin.Context) {
	page := shared.NewPageRequest(
		queryInt(c, "page", 1),
		queryInt(c, "pageSize", 20),
	)

	result, err := h.service.List(c.Request.Context(), page)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	dtoResult := shared.PageResponse[patientResponse]{
		Items:      make([]patientResponse, len(result.Items)),
		TotalItems: result.TotalItems,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
	for i, p := range result.Items {
		dtoResult.Items[i] = toPatientResponse(p)
	}

	respondSuccess(c, http.StatusOK, dtoResult)
}

// Get maneja GET /api/v1/pacientes/:idOrDoc. Detecta el tipo de
// identificador del path param: si es numérico busca el detalle por id
// (SP webPacientesListarIdPaciente); si no, busca por número de documento
// (SP sp_listarPaciente). Es la fusión de los endpoints de FastAPI y del
// Go original para no romper el frontend.
//
// @Summary Busca un paciente por id o por número de documento
// @Description Si el path param es numérico invoca webPacientesListarIdPaciente (detalle); en caso contrario sp_listarPaciente.
// @Tags Pacientes
// @Produce json
// @Param idOrDoc path string true "Id o número de documento del paciente"
// @Success 200 {object} apiResponse{data=domain.PatientDetail}
// @Failure 400 {object} apiResponse{error=apiError} "Identificador inválido"
// @Failure 404 {object} apiResponse{error=apiError} "Paciente no encontrado"
// @Router /pacientes/{idOrDoc} [get]
func (h *PatientHandler) Get(c *gin.Context) {
	param := c.Param("idOrDoc")
	if param == "" {
		respondError(c, http.StatusBadRequest, "INVALID_PATIENT_ID", domain.ErrInvalidDocumentNumber.Error())
		return
	}

	if id, err := strconv.ParseInt(param, 10, 64); err == nil {
		detail, derr := h.service.GetByID(c.Request.Context(), id)
		if derr != nil {
			h.respondGetError(c, derr)
			return
		}
		respondSuccess(c, http.StatusOK, toPatientDetailResponse(detail))
		return
	}

	patient, err := h.service.GetByDocumentNumber(c.Request.Context(), param)
	if err != nil {
		h.respondGetError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, toPatientResponse(patient))
}

// respondGetError traduce los errores de los casos de uso de consulta.
func (h *PatientHandler) respondGetError(c *gin.Context, err error) {
	switch err {
	case domain.ErrInvalidDocumentNumber, domain.ErrInvalidPatientID:
		respondError(c, http.StatusBadRequest, "INVALID_PATIENT_ID", err.Error())
	case domain.ErrPatientNotFound:
		respondError(c, http.StatusNotFound, "PATIENT_NOT_FOUND", err.Error())
	default:
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}

// GetDatosAdicionales maneja GET /api/v1/pacientes/:idOrDoc/datos-adicionales.
//
// @Summary Obtiene los antecedentes y datos adicionales de un paciente
// @Description Invoca el SP usp_go_PacientesDatosAdicionalesIdPaciente para devolver antecedentes quirúrgicos, patológicos, obstétricos, etc.
// @Tags Pacientes
// @Produce json
// @Param idOrDoc path int true "Id del paciente"
// @Success 200 {object} apiResponse{data=domain.PacienteDatosAdicionales}
// @Failure 400 {object} apiResponse{error=apiError} "Identificador inválido"
// @Router /pacientes/{idOrDoc}/datos-adicionales [get]
func (h *PatientHandler) GetDatosAdicionales(c *gin.Context) {
	idParam := c.Param("idOrDoc")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_PATIENT_ID", domain.ErrInvalidPatientID.Error())
		return
	}

	datos, err := h.service.GetDatosAdicionales(c.Request.Context(), id)
	if err != nil {
		h.respondGetError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, datos)
}

// UpdateDatosAdicionales maneja PUT /api/v1/pacientes/:idOrDoc/datos-adicionales.
//
// @Summary Actualiza los antecedentes y datos adicionales de un paciente
// @Description Invoca los SPs PacientesDatosAdicionalesModificar / PacientesDatosAdicionalesAgregar
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param idOrDoc path int true "Id del paciente"
// @Param request body domain.PacienteDatosAdicionales true "Antecedentes a guardar"
// @Success 200 {object} apiResponse{data=domain.PacienteDatosAdicionales}
// @Failure 400 {object} apiResponse{error=apiError} "Identificador o payload inválido"
// @Router /pacientes/{idOrDoc}/datos-adicionales [put]
func (h *PatientHandler) UpdateDatosAdicionales(c *gin.Context) {
	idParam := c.Param("idOrDoc")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_PATIENT_ID", domain.ErrInvalidPatientID.Error())
		return
	}

	var req domain.PacienteDatosAdicionales
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Cuerpo de solicitud inválido")
		return
	}

	idUsuario := 1
	if val, exists := c.Get("idEmpleado"); exists {
		switch v := val.(type) {
		case int:
			idUsuario = v
		case float64:
			idUsuario = int(v)
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				idUsuario = parsed
			}
		}
	}

	updated, err := h.service.UpdateDatosAdicionales(c.Request.Context(), id, req, idUsuario)
	if err != nil {
		h.respondGetError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, updated)
}

// Update maneja PUT /api/v1/pacientes/:id.
//
// @Summary Modifica un paciente
// @Description Actualiza los datos editables de un paciente (SP usp_go_ModificarPaciente) y devuelve el detalle actualizado.
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param id path int true "Id del paciente"
// @Param request body updatePatientRequest true "Datos editables del paciente"
// @Success 200 {object} apiResponse{data=domain.PatientDetail}
// @Failure 400 {object} apiResponse{error=apiError} "Validación inválida"
// @Failure 404 {object} apiResponse{error=apiError} "Paciente no encontrado"
// @Router /pacientes/{id} [put]
func (h *PatientHandler) Update(c *gin.Context) {
	param := c.Param("idOrDoc")
	if param == "" {
		param = c.Param("id")
	}
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_PATIENT_ID", domain.ErrInvalidPatientID.Error())
		return
	}

	var req updatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	detail, err := h.service.Update(c.Request.Context(), id, req.toDomain())
	if err != nil {
		switch err {
		case domain.ErrInvalidPatientID:
			respondError(c, http.StatusBadRequest, "INVALID_PATIENT_ID", err.Error())
		case domain.ErrPatientNotFound:
			respondError(c, http.StatusNotFound, "PATIENT_NOT_FOUND", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	respondSuccess(c, http.StatusOK, toPatientDetailResponse(detail))
}

// Create maneja POST /api/v1/pacientes.
//
// @Summary Registra un paciente nuevo
// @Description Persiste el paciente invocando el SP WebPacienteAgregar_E_H y devuelve el detalle del paciente creado.
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param request body createPatientRequest true "Datos del paciente a registrar"
// @Success 201 {object} apiResponse{data=patientDetailResponse} "Paciente creado"
// @Failure 400 {object} apiResponse{error=apiError} "Cuerpo inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al registrar el paciente"
// @Router /pacientes [post]
func (h *PatientHandler) Create(c *gin.Context) {
	var req createPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	detail, err := h.service.Create(c.Request.Context(), req.toDomain())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "PATIENT_CREATE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, toPatientDetailResponse(detail))
}

// Delete maneja DELETE /api/v1/pacientes/:id.
//
// @Summary Elimina un paciente
// @Description Verifica con PacientesSePuedeEliminar que el paciente no tenga registros asociados y lo elimina (SP PacientesEliminarPorIdPaciente).
// @Tags Pacientes
// @Produce json
// @Param id path int true "Id del paciente"
// @Success 200 {object} apiResponse "Paciente eliminado"
// @Failure 400 {object} apiResponse{error=apiError} "Id inválido"
// @Failure 404 {object} apiResponse{error=apiError} "Paciente no encontrado"
// @Failure 409 {object} apiResponse{error=apiError} "Paciente con registros asociados, no se puede eliminar"
// @Router /pacientes/{id} [delete]
func (h *PatientHandler) Delete(c *gin.Context) {
	param := c.Param("idOrDoc")
	if param == "" {
		param = c.Param("id")
	}
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_PATIENT_ID", domain.ErrInvalidPatientID.Error())
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		switch err {
		case domain.ErrInvalidPatientID:
			respondError(c, http.StatusBadRequest, "INVALID_PATIENT_ID", err.Error())
		case domain.ErrPatientNotFound:
			respondError(c, http.StatusNotFound, "PATIENT_NOT_FOUND", err.Error())
		case domain.ErrPatientCannotBeDeleted:
			respondError(c, http.StatusConflict, "PATIENT_CANNOT_BE_DELETED", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	respondSuccess(c, http.StatusOK, map[string]bool{"deleted": true})
}

// GetByDocumentAndType maneja GET
// /api/v1/pacientes/por-documento?nroDocumento=&idTipoDocIdentidad=.//
// @Summary Busca un paciente por número y tipo de documento
// @Description Invoca el SP usp_go_ListarPacientePorNroDocyTipo con el número y el tipo de documento.
// @Tags Pacientes
// @Produce json
// @Param nroDocumento query string true "Número de documento del paciente"
// @Param idTipoDocIdentidad query int true "Id del tipo de documento de identidad"
// @Success 200 {object} apiResponse{data=domain.Patient}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 404 {object} apiResponse{error=apiError} "Paciente no encontrado"
// @Router /pacientes/por-documento [get]
func (h *PatientHandler) GetByDocumentAndType(c *gin.Context) {
	documentNumber := c.Query("nroDocumento")
	docTypeID, err := strconv.ParseInt(c.Query("idTipoDocIdentidad"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_DOCUMENT_TYPE", domain.ErrInvalidDocumentType.Error())
		return
	}

	patient, err := h.service.GetByDocumentNumberAndType(c.Request.Context(), documentNumber, docTypeID)
	if err != nil {
		switch err {
		case domain.ErrInvalidDocumentNumber:
			respondError(c, http.StatusBadRequest, "INVALID_DOCUMENT_NUMBER", err.Error())
		case domain.ErrInvalidDocumentType:
			respondError(c, http.StatusBadRequest, "INVALID_DOCUMENT_TYPE", err.Error())
		case domain.ErrPatientNotFound:
			respondError(c, http.StatusNotFound, "PATIENT_NOT_FOUND", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	respondSuccess(c, http.StatusOK, toPatientResponse(patient))
}

// Search maneja GET /api/v1/pacientes/buscar?documento=&hc=&paterno=&materno=&nombres=.
//
// @Summary Busca pacientes aplicando filtros opcionales
// @Description Invoca el procedimiento almacenado usp_go_listarpacientes. Los filtros vacíos se ignoran.
// @Tags Pacientes
// @Produce json
// @Param documento query string false "Número de documento del paciente"
// @Param hc query string false "Número de historia clínica"
// @Param paterno query string false "Apellido paterno"
// @Param materno query string false "Apellido materno"
// @Param nombres query string false "Primer nombre"
// @Success 200 {object} apiResponse{data=[]domain.Patient}
// @Failure 500 {object} apiResponse{error=apiError} "Error interno"
// @Router /pacientes/buscar [get]
func (h *PatientHandler) Search(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	params := shared.PatientSearchParams{
		DocumentNumber:  c.Query("documento"),
		HistoryNumber:   c.Query("hc"),
		PaternalSurname: c.Query("paterno"),
		MaternalSurname: c.Query("materno"),
		Names:           c.Query("nombres"),
	}

	result, err := h.service.Search(ctx, params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(c, http.StatusGatewayTimeout, "TIMEOUT_ERROR", "La búsqueda excedió el tiempo límite. Por favor use filtros más específicos.")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	dtoResult := make([]patientResponse, 0, len(result))
	for _, p := range result {
		dtoResult = append(dtoResult, toPatientResponse(p))
	}

	respondSuccess(c, http.StatusOK, dtoResult)
}

func queryInt(c *gin.Context, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
