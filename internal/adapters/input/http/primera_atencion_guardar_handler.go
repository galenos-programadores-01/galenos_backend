package httpadapter

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/adapters/output/persistence/sqlserver"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type primeraAtencionIniciarRequest struct {
	IdAtencion    int     `json:"idAtencion" binding:"required"`
	FechaAtencion *string `json:"fechaAtencion"`
	TipoAnamnesis *string `json:"tipoAnamnesis"`
}

type primeraAtencionGuardarRequest struct {
	FechaAtencion                   *string `json:"fechaAtencion"`
	TipoAnamnesis                   *string `json:"tipoAnamnesis"`
	TiempoEnfermedad                *int    `json:"tiempoEnfermedad"`
	IdTiempoEnfermedad              *int    `json:"idTiempoEnfermedad"`
	SignosSintomas                  *string `json:"signosSintomas"`
	RelatoCronologico               *string `json:"relatoCronologico"`
	Antecedentes                    *string `json:"antecedentes"`
	IdEstadoGeneral                 *int    `json:"idEstadoGeneral"`
	IdEstadoHidratacion             *int    `json:"idEstadoHidratacion"`
	IdNivelConciencia               *int    `json:"idNivelConciencia"`
	PielAnexos                      *string `json:"pielAnexos"`
	ExamenClinicoRegional           *string `json:"examenClinicoRegional"`
	ReferenciaOportuna              *string `json:"referenciaOportuna"`
	Evolucion                       *string `json:"evolucion"`
	EvaluacionEnfermeriaObstetricia *string `json:"evaluacionEnfermeriaObstetricia"`
	IdCondicionEgreso               *int    `json:"idCondicionEgreso"`
	CodCondicionEgreso              *string `json:"codCondicionEgreso"`
	FechaEgreso                     *string `json:"fechaEgreso"`
}

type primeraAtencionGuardarResponse struct {
	IdPrimeraAtencion int    `json:"idPrimeraAtencion"`
	Mensaje           string `json:"mensaje"`
}

func idUsuarioDesdeJWT(c *gin.Context) int {
	return c.GetInt("idEmpleado")
}

// IniciarBorrador POST /api/v1/primera-atencion/iniciar
func (h *PrimeraAtencionHandler) IniciarBorrador(c *gin.Context) {
	idUsuario := idUsuarioDesdeJWT(c)
	if idUsuario == 0 {
		respondError(c, http.StatusUnauthorized, "MISSING_EMPLOYEE", "el token no contiene idEmpleado")
		return
	}

	var req primeraAtencionIniciarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if req.FechaAtencion != nil || req.TipoAnamnesis != nil {
		in, err := buildUpsertFromIniciar(req, idUsuario)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		result, err := h.service.GuardarParcial(c.Request.Context(), in)
		if err != nil {
			respondPrimeraAtencionWriteError(c, "PRIMERA_ATENCION_INICIAR_FAILED", err)
			return
		}
		respondSuccess(c, http.StatusOK, primeraAtencionGuardarResponse{
			IdPrimeraAtencion: result.IdPrimeraAtencion,
			Mensaje:           result.Mensaje,
		})
		return
	}

	result, err := h.service.IniciarBorrador(c.Request.Context(), req.IdAtencion, idUsuario)
	if err != nil {
		respondPrimeraAtencionWriteError(c, "PRIMERA_ATENCION_INICIAR_FAILED", err)
		return
	}
	respondSuccess(c, http.StatusOK, primeraAtencionGuardarResponse{
		IdPrimeraAtencion: result.IdPrimeraAtencion,
		Mensaje:           result.Mensaje,
	})
}

// ObtenerDetalle GET /api/v1/primera-atencion/detalle/:idAtencion
func (h *PrimeraAtencionHandler) ObtenerDetalle(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	row, err := h.service.ObtenerDetalle(c.Request.Context(), idAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "PRIMERA_ATENCION_DETALLE_FAILED", err.Error())
		return
	}
	if row == nil {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "no existe borrador para esta atención")
		return
	}
	respondSuccess(c, http.StatusOK, row)
}

// GuardarParcial PATCH /api/v1/primera-atencion/:idAtencion
func (h *PrimeraAtencionHandler) GuardarParcial(c *gin.Context) {
	idUsuario := idUsuarioDesdeJWT(c)
	if idUsuario == 0 {
		respondError(c, http.StatusUnauthorized, "MISSING_EMPLOYEE", "el token no contiene idEmpleado")
		return
	}

	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req primeraAtencionGuardarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	in, err := buildUpsertFromGuardar(idAtencion, idUsuario, req)
	if err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := h.service.GuardarParcial(c.Request.Context(), in)
	if err != nil {
		respondPrimeraAtencionWriteError(c, "PRIMERA_ATENCION_GUARDAR_FAILED", err)
		return
	}
	respondSuccess(c, http.StatusOK, primeraAtencionGuardarResponse{
		IdPrimeraAtencion: result.IdPrimeraAtencion,
		Mensaje:           result.Mensaje,
	})
}

type primeraAtencionCerrarFirmaRequest struct {
	UuidFirma        string `json:"uuidFirma"`
	UbicacionArchivo string `json:"ubicacionArchivo"`
}

type primeraAtencionCerrarFirmaResponse struct {
	IdPrimeraAtencion int    `json:"idPrimeraAtencion"`
	UbicacionArchivo  string `json:"ubicacionArchivo"`
	Mensaje           string `json:"mensaje"`
}

// CerrarFirma POST /api/v1/primera-atencion/:idAtencion/cerrar-firma
func (h *PrimeraAtencionHandler) CerrarFirma(c *gin.Context) {
	idUsuario := idUsuarioDesdeJWT(c)
	if idUsuario == 0 {
		respondError(c, http.StatusUnauthorized, "MISSING_EMPLOYEE", "el token no contiene idEmpleado")
		return
	}

	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req primeraAtencionCerrarFirmaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	uuidFirma := strings.TrimSpace(req.UuidFirma)
	ubicacion := strings.TrimSpace(req.UbicacionArchivo)
	if ubicacion == "" {
		if uuidFirma == "" {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "uuidFirma es obligatorio")
			return
		}
		ubicacion = h.ubicacionLogicaPrimeraAtencion(idAtencion, uuidFirma)
	}

	// Si el almacenamiento remoto está activo: subir PDF(s) firmados y usar la ruta del PA principal.
	if h.archivos != nil && h.archivos.Enabled() {
		if uuidFirma == "" {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "uuidFirma es obligatorio para subir al almacenamiento remoto")
			return
		}
		if h.firmaStore == nil {
			respondError(c, http.StatusInternalServerError, "REMOTE_FIRMA_STORE_MISSING", "almacén de Firma Perú no disponible")
			return
		}

		loteFiles, errLote := h.firmaStore.GetLoteSignedFiles(c.Request.Context(), uuidFirma)
		if errLote == nil && len(loteFiles) > 0 {
			ubiPrincipal := ""
			ubiFallback := ""
			for _, f := range loteFiles {
				ubi, errUp := h.archivos.SubirPrimeraAtencionNombrado(c.Request.Context(), idAtencion, f.Name, f.Content)
				if errUp != nil {
					respondError(c, http.StatusBadGateway, "REMOTE_UPLOAD_FAILED", "no se pudo subir el PDF del lote al remoto: "+errUp.Error())
					return
				}
				if ubiFallback == "" {
					ubiFallback = ubi
				}
				if esPdfPrincipalPA(idAtencion, f.Name) {
					ubiPrincipal = ubi
				}
			}
			if ubiPrincipal != "" {
				ubicacion = ubiPrincipal
			} else {
				ubicacion = ubiFallback
			}
		} else {
			pdf, errPdf := h.firmaStore.GetSigned(c.Request.Context(), uuidFirma)
			if errPdf != nil {
				respondError(c, http.StatusConflict, "FIRMA_PDF_NO_DISPONIBLE", "no se encontró el PDF firmado para subir al remoto: "+errPdf.Error())
				return
			}
			ubi, errUp := h.archivos.SubirPrimeraAtencion(c.Request.Context(), idAtencion, uuidFirma, pdf)
			if errUp != nil {
				respondError(c, http.StatusBadGateway, "REMOTE_UPLOAD_FAILED", "no se pudo subir el PDF al almacenamiento remoto: "+errUp.Error())
				return
			}
			ubicacion = ubi
		}
	}

	result, err := h.service.CerrarFirma(c.Request.Context(), domain.PrimeraAtencionCerrarFirma{
		IdAtencion:        idAtencion,
		IdUsuarioModifica: idUsuario,
		UbicacionArchivo:  ubicacion,
	})
	if err != nil {
		respondPrimeraAtencionWriteError(c, "PRIMERA_ATENCION_CERRAR_FIRMA_FAILED", err)
		return
	}

	respondSuccess(c, http.StatusOK, primeraAtencionCerrarFirmaResponse{
		IdPrimeraAtencion: result.IdPrimeraAtencion,
		UbicacionArchivo:  ubicacion,
		Mensaje:           result.Mensaje,
	})
}

// ubicacionLogicaPrimeraAtencion usa el patrón SFTP configurado o el default /ftp/...
func (h *PrimeraAtencionHandler) ubicacionLogicaPrimeraAtencion(idAtencion int, uuidFirma string) string {
	if h.archivos != nil {
		// Preferir nombre canónico del PDF principal del lote.
		nombrePA := fmt.Sprintf("PA_%d.pdf", idAtencion)
		if _, ubi, err := h.archivos.Rutas().PrimeraAtencionPDFNombrado(idAtencion, nombrePA); err == nil && ubi != "" {
			return ubi
		}
		if _, ubi, err := h.archivos.Rutas().PrimeraAtencionPDF(idAtencion, uuidFirma); err == nil && ubi != "" {
			return ubi
		}
	}
	return fmt.Sprintf("/ftp/primera-atencion/%d/PA_%d.pdf", idAtencion, idAtencion)
}

// esPdfPrincipalPA identifica PA_{id}.pdf (no satélites _EX_ / _MED_).
func esPdfPrincipalPA(idAtencion int, nombre string) bool {
	base := strings.ToLower(filepath.Base(strings.TrimSpace(nombre)))
	want := strings.ToLower(fmt.Sprintf("pa_%d.pdf", idAtencion))
	return base == want
}

func respondPrimeraAtencionWriteError(c *gin.Context, code string, err error) {
	switch {
	case errors.Is(err, domain.ErrPrimeraAtencionYaFirmada):
		respondError(c, http.StatusConflict, "PRIMERA_ATENCION_FIRMADA", err.Error())
	case errors.Is(err, domain.ErrPrimeraAtencionSinBorrador):
		respondError(c, http.StatusNotFound, "NOT_FOUND", err.Error())
	default:
		respondError(c, http.StatusInternalServerError, code, err.Error())
	}
}

func buildUpsertFromIniciar(req primeraAtencionIniciarRequest, idUsuario int) (domain.PrimeraAtencionUpsert, error) {
	now := time.Now()
	fecha := &now
	if req.FechaAtencion != nil {
		parsed, err := sqlserver.ParseFechaHoraLocal(*req.FechaAtencion)
		if err != nil {
			return domain.PrimeraAtencionUpsert{}, err
		}
		if parsed != nil {
			fecha = parsed
		}
	}
	tipo := "D"
	if req.TipoAnamnesis != nil && *req.TipoAnamnesis != "" {
		tipo = *req.TipoAnamnesis
	}
	return domain.PrimeraAtencionUpsert{
		IdAtencion:        req.IdAtencion,
		FechaAtencion:     fecha,
		TipoAnamnesis:     &tipo,
		IdUsuarioRegistro: idUsuario,
	}, nil
}

func buildUpsertFromGuardar(idAtencion, idUsuario int, req primeraAtencionGuardarRequest) (domain.PrimeraAtencionUpsert, error) {
	in := domain.PrimeraAtencionUpsert{
		IdAtencion:                      idAtencion,
		TiempoEnfermedad:                req.TiempoEnfermedad,
		IdTiempoEnfermedad:              req.IdTiempoEnfermedad,
		SignosSintomas:                  req.SignosSintomas,
		RelatoCronologico:               req.RelatoCronologico,
		Antecedentes:                    req.Antecedentes,
		IdEstadoGeneral:                 req.IdEstadoGeneral,
		IdEstadoHidratacion:             req.IdEstadoHidratacion,
		IdNivelConciencia:               req.IdNivelConciencia,
		PielAnexos:                      req.PielAnexos,
		ExamenClinicoRegional:           req.ExamenClinicoRegional,
		ReferenciaOportuna:              req.ReferenciaOportuna,
		Evolucion:                       req.Evolucion,
		EvaluacionEnfermeriaObstetricia: req.EvaluacionEnfermeriaObstetricia,
		IdCondicionEgreso:               req.IdCondicionEgreso,
		CodCondicionEgreso:              req.CodCondicionEgreso,
		IdUsuarioRegistro:               idUsuario,
		IdUsuarioModifica:               &idUsuario,
	}

	if req.FechaAtencion != nil {
		parsed, err := sqlserver.ParseFechaHoraLocal(*req.FechaAtencion)
		if err != nil {
			return domain.PrimeraAtencionUpsert{}, err
		}
		in.FechaAtencion = parsed
	}
	if req.TipoAnamnesis != nil && *req.TipoAnamnesis != "" {
		in.TipoAnamnesis = req.TipoAnamnesis
	}
	if req.FechaEgreso != nil {
		parsed, err := sqlserver.ParseFechaHoraLocal(*req.FechaEgreso)
		if err != nil {
			return domain.PrimeraAtencionUpsert{}, err
		}
		in.FechaEgreso = parsed
	}
	now := time.Now()
	in.FechaModifica = &now

	return in, nil
}

type primeraAtencionDiagnosticoItemRequest struct {
	IdDiagnostico        int    `json:"idDiagnostico" binding:"required"`
	IdSubclasificacionDx int    `json:"idSubclasificacionDx"`
	Tipo                 string `json:"tipo"` // legado P|D|R si no viene idSubclasificacionDx
}

type primeraAtencionRegistrarDiagnosticosRequest struct {
	IdPaciente         int                                     `json:"idPaciente" binding:"required"`
	IdPrimeraAtencion  *int                                    `json:"idPrimeraAtencion"`
	IdEpisodio         int                                     `json:"idEpisodio" binding:"required"`
	IdClasificacionDx  int                                     `json:"idClasificacionDx" binding:"required"`
	Diagnosticos       []primeraAtencionDiagnosticoItemRequest `json:"diagnosticos" binding:"required,min=1"`
}

type primeraAtencionDiagnosticoResultadoResponse struct {
	IdDiagnostico            int    `json:"idDiagnostico"`
	IdAtencionDiagnostico    int    `json:"idAtencionDiagnostico"`
	SubclasificacionAnterior string `json:"subclasificacionAnterior"`
}

type primeraAtencionRegistrarDiagnosticosResponse struct {
	Resultados []primeraAtencionDiagnosticoResultadoResponse `json:"resultados"`
	Mensaje    string                                        `json:"mensaje"`
}

// RegistrarDiagnosticos POST /api/v1/primera-atencion/:idAtencion/diagnosticos
func (h *PrimeraAtencionHandler) RegistrarDiagnosticos(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req primeraAtencionRegistrarDiagnosticosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	items := make([]domain.DiagnosticoRegistrarItem, 0, len(req.Diagnosticos))
	for _, d := range req.Diagnosticos {
		tipo := domain.TipoDiagnosticoUI(strings.ToUpper(strings.TrimSpace(d.Tipo)))
		items = append(items, domain.DiagnosticoRegistrarItem{
			IdDiagnostico:        d.IdDiagnostico,
			IdSubclasificacionDx: d.IdSubclasificacionDx,
			Tipo:                 tipo,
		})
	}

	idPrimera := 0
	if req.IdPrimeraAtencion != nil {
		idPrimera = *req.IdPrimeraAtencion
	}
	if req.IdEpisodio <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idEpisodio es obligatorio")
		return
	}
	if req.IdClasificacionDx <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idClasificacionDx es obligatorio")
		return
	}

	resultados, err := h.service.RegistrarDiagnosticos(
		c.Request.Context(),
		idAtencion,
		req.IdPaciente,
		idPrimera,
		req.IdEpisodio,
		req.IdClasificacionDx,
		items,
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_DX_FAILED", err.Error())
		return
	}

	respItems := make([]primeraAtencionDiagnosticoResultadoResponse, 0, len(resultados))
	for i, r := range resultados {
		idDx := 0
		if i < len(req.Diagnosticos) {
			idDx = req.Diagnosticos[i].IdDiagnostico
		}
		respItems = append(respItems, primeraAtencionDiagnosticoResultadoResponse{
			IdDiagnostico:            idDx,
			IdAtencionDiagnostico:    r.IdAtencionDiagnostico,
			SubclasificacionAnterior: r.SubclasificacionAnterior,
		})
	}

	respondSuccess(c, http.StatusOK, primeraAtencionRegistrarDiagnosticosResponse{
		Resultados: respItems,
		Mensaje:    "Diagnósticos registrados",
	})
}

type primeraAtencionDiagnosticoListItemResponse struct {
	IdAtencionDiagnostico int    `json:"idAtencionDiagnostico"`
	IdDiagnostico         int    `json:"idDiagnostico"`
	IdPrimeraAtencion     int    `json:"idPrimeraAtencion"`
	Descripcion           string `json:"descripcion"`
	CodigoCIE10           string `json:"codigoCIE10"`
	Tipo                  string `json:"tipo"`
	IdSubclasificacionDx  int    `json:"idSubclasificacionDx"`
	NumOrdenCorrelativo   int    `json:"numOrdenCorrelativo"`
}

// ListarDiagnosticos GET /api/v1/primera-atencion/:idAtencion/diagnosticos
func (h *PrimeraAtencionHandler) ListarDiagnosticos(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	idPrimera := 0
	if raw := c.Query("idPrimeraAtencion"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion inválido")
			return
		}
		idPrimera = v
	}

	items, err := h.service.ListarDiagnosticos(c.Request.Context(), idAtencion, idPrimera)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "PRIMERA_ATENCION_DX_LIST_FAILED", err.Error())
		return
	}

	resp := make([]primeraAtencionDiagnosticoListItemResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, primeraAtencionDiagnosticoListItemResponse{
			IdAtencionDiagnostico: it.IdAtencionDiagnostico,
			IdDiagnostico:         it.IdDiagnostico,
			IdPrimeraAtencion:     it.IdPrimeraAtencion,
			Descripcion:           it.Descripcion,
			CodigoCIE10:           it.CodigoCIE10,
			Tipo:                  string(it.Tipo),
			IdSubclasificacionDx:  it.IdSubclasificacionDx,
			NumOrdenCorrelativo:   it.NumOrdenCorrelativo,
		})
	}
	respondSuccess(c, http.StatusOK, resp)
}

// EliminarDiagnostico DELETE /api/v1/primera-atencion/:idAtencion/diagnosticos/:idAtencionDiagnostico?idPrimeraAtencion=
func (h *PrimeraAtencionHandler) EliminarDiagnostico(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}
	idAtencionDiagnostico, err := strconv.Atoi(c.Param("idAtencionDiagnostico"))
	if err != nil || idAtencionDiagnostico <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencionDiagnostico inválido")
		return
	}
	idPrimeraAtencion, err := strconv.Atoi(c.Query("idPrimeraAtencion"))
	if err != nil || idPrimeraAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion es obligatorio")
		return
	}

	if err := h.service.EliminarDiagnostico(c.Request.Context(), idAtencion, idPrimeraAtencion, idAtencionDiagnostico); err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_DX_DELETE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idAtencionDiagnostico": idAtencionDiagnostico,
		"idPrimeraAtencion":     idPrimeraAtencion,
		"mensaje":               "Diagnóstico eliminado",
	})
}

type planTrabajoItemRequest struct {
	IdPrimeraAtencion int     `json:"idPrimeraAtencion" binding:"required"`
	Tipo              string  `json:"tipo" binding:"required"`
	IdProducto        *int    `json:"idProducto"`
	IdEspecialidad    *int    `json:"idEspecialidad"`
	Codigo            string  `json:"codigo"`
	Nombre            string  `json:"nombre"`
}

type planTrabajoItemResponse struct {
	IdPlanTrabajoItem int     `json:"idPlanTrabajoItem"`
	IdAtencion        int     `json:"idAtencion"`
	IdPrimeraAtencion int     `json:"idPrimeraAtencion"`
	Tipo              string  `json:"tipo"`
	IdProducto        *int    `json:"idProducto,omitempty"`
	IdEspecialidad    *int    `json:"idEspecialidad,omitempty"`
	Codigo            string  `json:"codigo"`
	Nombre            string  `json:"nombre"`
	NumOrden          int     `json:"numOrden"`
}

// RegistrarPlanTrabajoItem POST /api/v1/primera-atencion/:idAtencion/plan-trabajo
func (h *PrimeraAtencionHandler) RegistrarPlanTrabajoItem(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req planTrabajoItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	res, err := h.service.RegistrarPlanTrabajoItem(
		c.Request.Context(),
		idAtencion,
		req.IdPrimeraAtencion,
		domain.PlanTrabajoItemRegistrarItem{
			Tipo:           req.Tipo,
			IdProducto:     req.IdProducto,
			IdEspecialidad: req.IdEspecialidad,
			Codigo:         req.Codigo,
			Nombre:         req.Nombre,
		},
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_PLAN_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idPlanTrabajoItem": res.IdPlanTrabajoItem,
		"mensaje":           "Ítem de plan de trabajo registrado",
	})
}

// ListarPlanTrabajoItems GET /api/v1/primera-atencion/:idAtencion/plan-trabajo
func (h *PrimeraAtencionHandler) ListarPlanTrabajoItems(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	idPrimera := 0
	if raw := c.Query("idPrimeraAtencion"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion inválido")
			return
		}
		idPrimera = v
	}
	tipo := strings.TrimSpace(c.Query("tipo"))

	items, err := h.service.ListarPlanTrabajoItems(c.Request.Context(), idAtencion, idPrimera, tipo)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "PRIMERA_ATENCION_PLAN_LIST_FAILED", err.Error())
		return
	}

	resp := make([]planTrabajoItemResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, planTrabajoItemResponse{
			IdPlanTrabajoItem: it.IdPlanTrabajoItem,
			IdAtencion:        it.IdAtencion,
			IdPrimeraAtencion: it.IdPrimeraAtencion,
			Tipo:              it.Tipo,
			IdProducto:        it.IdProducto,
			IdEspecialidad:    it.IdEspecialidad,
			Codigo:            it.Codigo,
			Nombre:            it.Nombre,
			NumOrden:          it.NumOrden,
		})
	}
	respondSuccess(c, http.StatusOK, resp)
}

// EliminarPlanTrabajoItem DELETE /api/v1/primera-atencion/:idAtencion/plan-trabajo/:idPlanTrabajoItem?idPrimeraAtencion=
func (h *PrimeraAtencionHandler) EliminarPlanTrabajoItem(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}
	idPlanTrabajoItem, err := strconv.Atoi(c.Param("idPlanTrabajoItem"))
	if err != nil || idPlanTrabajoItem <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPlanTrabajoItem inválido")
		return
	}
	idPrimeraAtencion, err := strconv.Atoi(c.Query("idPrimeraAtencion"))
	if err != nil || idPrimeraAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion es obligatorio")
		return
	}

	if err := h.service.EliminarPlanTrabajoItem(c.Request.Context(), idAtencion, idPrimeraAtencion, idPlanTrabajoItem); err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_PLAN_DELETE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idPlanTrabajoItem": idPlanTrabajoItem,
		"idPrimeraAtencion": idPrimeraAtencion,
		"mensaje":           "Ítem de plan de trabajo eliminado",
	})
}

type medicamentoRegistrarRequest struct {
	IdPrimeraAtencion   int     `json:"idPrimeraAtencion" binding:"required"`
	IdCuentaAtencion    int     `json:"idCuentaAtencion" binding:"required"`
	IdPuntoCarga        int     `json:"idPuntoCarga" binding:"required"`
	IdServicioReceta    int     `json:"idServicioReceta" binding:"required"`
	IdPaciente          int     `json:"idPaciente" binding:"required"`
	IdUsuarioAuditoria  int     `json:"idUsuarioAuditoria"`
	IdProducto          int     `json:"idProducto" binding:"required"`
	Cantidad            int     `json:"cantidad" binding:"required"`
	Precio              float64 `json:"precio"` // precioUnitario del catálogo
	Duracion            int     `json:"duracion" binding:"required"`
	IdDosis             *int    `json:"idDosis"`
	IdUNIDDosis         *int    `json:"idUNIDDosis"`
	IdFrecuencia        *int    `json:"idFrecuencia"`
	IdViaAdministracion *int    `json:"idViaAdministracion"`
	IdDiagnostico       int     `json:"idDiagnostico" binding:"required"`
	CodigoDiagnostico   string  `json:"codigoDiagnostico" binding:"required"`
}

// RegistrarMedicamento POST /api/v1/primera-atencion/:idAtencion/medicamentos
func (h *PrimeraAtencionHandler) RegistrarMedicamento(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req medicamentoRegistrarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	res, err := h.service.RegistrarMedicamento(c.Request.Context(), domain.RecetaMedicamentoRegistrar{
		IdAtencion:          idAtencion,
		IdPuntoCarga:        req.IdPuntoCarga,
		IdCuentaAtencion:    req.IdCuentaAtencion,
		IdServicioReceta:    req.IdServicioReceta,
		IdPaciente:          req.IdPaciente,
		IdUsuarioAuditoria:  req.IdUsuarioAuditoria,
		IdPrimeraAtencion:   req.IdPrimeraAtencion,
		IdProducto:          req.IdProducto,
		Cantidad:            req.Cantidad,
		Precio:              req.Precio,
		IdDosis:             req.IdDosis,
		IdUNIDDosis:         req.IdUNIDDosis,
		IdFrecuencia:        req.IdFrecuencia,
		IdViaAdministracion: req.IdViaAdministracion,
		Duracion:            req.Duracion,
		IdDiagnostico:       req.IdDiagnostico,
		CodigoDiagnostico:   req.CodigoDiagnostico,
	})
	if err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_MED_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idReceta":          res.IdReceta,
		"idProducto":        res.IdProducto,
		"cantidad":          res.Cantidad,
		"precio":            res.Precio,
		"total":             res.Total,
		"idDiagnostico":     req.IdDiagnostico,
		"codigoDiagnostico": req.CodigoDiagnostico,
		"idAtencion":        idAtencion,
		"mensaje":           res.Mensaje,
	})
}

type medicamentoCamposActualizarRequest struct {
	IdCuentaAtencion    int  `json:"idCuentaAtencion" binding:"required"`
	IdPrimeraAtencion   int  `json:"idPrimeraAtencion" binding:"required"`
	IdDosis             *int `json:"idDosis"`
	IdUNIDDosis         *int `json:"idUNIDDosis"`
	IdFrecuencia        *int `json:"idFrecuencia"`
	IdViaAdministracion *int `json:"idViaAdministracion"`
}

// ActualizarMedicamentoCampos PATCH .../medicamentos/:idReceta/:idProducto
// Actualiza dosis/unidad/frecuencia/vía de un detalle ya guardado.
func (h *PrimeraAtencionHandler) ActualizarMedicamentoCampos(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}
	idReceta, err := strconv.Atoi(c.Param("idReceta"))
	if err != nil || idReceta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idReceta inválido")
		return
	}
	idProducto, err := strconv.Atoi(c.Param("idProducto"))
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idProducto inválido")
		return
	}

	var req medicamentoCamposActualizarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	err = h.service.ActualizarMedicamentoCampos(c.Request.Context(), domain.RecetaMedicamentoCamposActualizar{
		IdAtencion:          idAtencion,
		IdReceta:            idReceta,
		IdProducto:          idProducto,
		IdCuentaAtencion:    req.IdCuentaAtencion,
		IdPrimeraAtencion:   req.IdPrimeraAtencion,
		IdDosis:             req.IdDosis,
		IdUNIDDosis:         req.IdUNIDDosis,
		IdFrecuencia:        req.IdFrecuencia,
		IdViaAdministracion: req.IdViaAdministracion,
	})
	if err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_MED_UPDATE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idReceta":   idReceta,
		"idProducto": idProducto,
		"mensaje":    "Campos de medicamento actualizados",
	})
}

type medicamentoListItemResponse struct {
	IdReceta            int     `json:"idReceta"`
	IdPuntoCarga        int     `json:"idPuntoCarga"`
	IdCuentaAtencion    int     `json:"idCuentaAtencion"`
	IdServicioReceta    int     `json:"idServicioReceta"`
	IdPrimeraAtencion   int     `json:"idPrimeraAtencion"`
	IdProducto          int     `json:"idProducto"`
	Cantidad            int     `json:"cantidad"`
	Precio              float64 `json:"precio"`
	IdDosis             *int    `json:"idDosis,omitempty"`
	IdUNIDDosis         *int    `json:"idUNIDDosis,omitempty"`
	IdFrecuencia        *int    `json:"idFrecuencia,omitempty"`
	IdViaAdministracion *int    `json:"idViaAdministracion,omitempty"`
	Duracion            *int    `json:"duracion,omitempty"`
	Codigo              string  `json:"codigo"`
	Nombre              string  `json:"nombre"`
	IdDiagnostico       *int    `json:"idDiagnostico"`
	CodigoDiagnostico   string  `json:"codigoDiagnostico"`
}

// ListarMedicamentos GET /api/v1/primera-atencion/:idAtencion/medicamentos
func (h *PrimeraAtencionHandler) ListarMedicamentos(c *gin.Context) {
	if _, err := strconv.Atoi(c.Param("idAtencion")); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	idCuenta, err := strconv.Atoi(c.Query("idCuentaAtencion"))
	if err != nil || idCuenta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idPrimera := 0
	if raw := c.Query("idPrimeraAtencion"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion inválido")
			return
		}
		idPrimera = v
	}
	idPunto := 0
	if raw := c.Query("idPuntoCarga"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPuntoCarga inválido")
			return
		}
		idPunto = v
	}

	items, err := h.service.ListarMedicamentos(c.Request.Context(), idCuenta, idPrimera, idPunto)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "PRIMERA_ATENCION_MED_LIST_FAILED", err.Error())
		return
	}

	resp := make([]medicamentoListItemResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, medicamentoListItemResponse{
			IdReceta:            it.IdReceta,
			IdPuntoCarga:        it.IdPuntoCarga,
			IdCuentaAtencion:    it.IdCuentaAtencion,
			IdServicioReceta:    it.IdServicioReceta,
			IdPrimeraAtencion:   it.IdPrimeraAtencion,
			IdProducto:          it.IdProducto,
			Cantidad:            it.Cantidad,
			Precio:              it.Precio,
			IdDosis:             it.IdDosis,
			IdUNIDDosis:         it.IdUNIDDosis,
			IdFrecuencia:        it.IdFrecuencia,
			IdViaAdministracion: it.IdViaAdministracion,
			Duracion:            it.Duracion,
			Codigo:              it.Codigo,
			Nombre:              it.Nombre,
			IdDiagnostico:       it.IdDiagnostico,
			CodigoDiagnostico:   it.CodigoDiagnostico,
		})
	}
	respondSuccess(c, http.StatusOK, resp)
}

// EliminarMedicamento DELETE .../medicamentos/:idReceta/:idProducto?idCuentaAtencion=&idPrimeraAtencion=
func (h *PrimeraAtencionHandler) EliminarMedicamento(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}
	idReceta, err := strconv.Atoi(c.Param("idReceta"))
	if err != nil || idReceta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idReceta inválido")
		return
	}
	idProducto, err := strconv.Atoi(c.Param("idProducto"))
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idProducto inválido")
		return
	}
	idCuenta, err := strconv.Atoi(c.Query("idCuentaAtencion"))
	if err != nil || idCuenta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idPrimera, err := strconv.Atoi(c.Query("idPrimeraAtencion"))
	if err != nil || idPrimera <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion es obligatorio")
		return
	}

	if err := h.service.EliminarMedicamento(c.Request.Context(), idAtencion, idReceta, idProducto, idCuenta, idPrimera); err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_MED_DELETE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idReceta":   idReceta,
		"idProducto": idProducto,
		"mensaje":    "Medicamento eliminado",
	})
}

type examenRegistrarRequest struct {
	IdPrimeraAtencion  int     `json:"idPrimeraAtencion" binding:"required"`
	IdCuentaAtencion   int     `json:"idCuentaAtencion" binding:"required"`
	IdPuntoCarga       int     `json:"idPuntoCarga" binding:"required"` // del catálogo (lab/rayos/…)
	IdServicioReceta   int     `json:"idServicioReceta" binding:"required"`
	IdPaciente         int     `json:"idPaciente" binding:"required"`
	IdUsuarioAuditoria int     `json:"idUsuarioAuditoria"`
	IdProducto         int     `json:"idProducto" binding:"required"`
	Precio             float64 `json:"precio"` // precioUnitario del catálogo
	IdDiagnostico      int     `json:"idDiagnostico" binding:"required"`      // → RecetaDetalle.IdDiagnostico
	CodigoDiagnostico  string  `json:"codigoDiagnostico" binding:"required"` // → RecetaDetalle.CodigoDiagnostico (CIE-10)
}

// RegistrarExamen POST /api/v1/primera-atencion/:idAtencion/examenes
// Cabecera por IdPuntoCarga del examen (crear/reusar) + detalle Cantidad=1, Precio, Total=Cantidad*Precio.
func (h *PrimeraAtencionHandler) RegistrarExamen(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req examenRegistrarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	res, err := h.service.RegistrarExamen(c.Request.Context(), domain.RecetaExamenRegistrar{
		IdAtencion:         idAtencion,
		IdPuntoCarga:       req.IdPuntoCarga,
		IdCuentaAtencion:   req.IdCuentaAtencion,
		IdServicioReceta:   req.IdServicioReceta,
		IdPaciente:         req.IdPaciente,
		IdUsuarioAuditoria: req.IdUsuarioAuditoria,
		IdPrimeraAtencion:  req.IdPrimeraAtencion,
		IdProducto:         req.IdProducto,
		Precio:             req.Precio,
		IdDiagnostico:      req.IdDiagnostico,
		CodigoDiagnostico:  req.CodigoDiagnostico,
	})
	if err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_EXA_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idReceta":           res.IdReceta,
		"idProducto":         res.IdProducto,
		"cantidad":           res.Cantidad,
		"precio":             res.Precio,
		"total":              res.Total,
		"idDiagnostico":      req.IdDiagnostico,
		"codigoDiagnostico":  req.CodigoDiagnostico,
		"idAtencion":         idAtencion,
		"mensaje":            res.Mensaje,
	})
}

type examenListItemResponse struct {
	IdReceta           int     `json:"idReceta"`
	IdPuntoCarga       int     `json:"idPuntoCarga"`
	PuntoCarga         string  `json:"puntoCarga"`
	IdCuentaAtencion   int     `json:"idCuentaAtencion"`
	IdServicioReceta   int     `json:"idServicioReceta"`
	IdPrimeraAtencion  int     `json:"idPrimeraAtencion"`
	IdProducto         int     `json:"idProducto"`
	Cantidad           int     `json:"cantidad"`
	Precio             float64 `json:"precio"`
	Total              float64 `json:"total"`
	Codigo             string  `json:"codigo"`
	Nombre             string  `json:"nombre"`
	IdDiagnostico      *int    `json:"idDiagnostico"`
	CodigoDiagnostico  string  `json:"codigoDiagnostico"`
}

// ListarExamenes GET .../examenes?idCuentaAtencion=&idPrimeraAtencion=
// Excluye farmacia (IdPuntoCarga=5) en el SP vía @ExcluirIdPuntoCarga.
func (h *PrimeraAtencionHandler) ListarExamenes(c *gin.Context) {
	if _, err := strconv.Atoi(c.Param("idAtencion")); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	idCuenta, err := strconv.Atoi(c.Query("idCuentaAtencion"))
	if err != nil || idCuenta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idPrimera := 0
	if raw := c.Query("idPrimeraAtencion"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion inválido")
			return
		}
		idPrimera = v
	}

	items, err := h.service.ListarExamenes(c.Request.Context(), idCuenta, idPrimera)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "PRIMERA_ATENCION_EXA_LIST_FAILED", err.Error())
		return
	}

	resp := make([]examenListItemResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, examenListItemResponse{
			IdReceta:          it.IdReceta,
			IdPuntoCarga:      it.IdPuntoCarga,
			PuntoCarga:        it.PuntoCarga,
			IdCuentaAtencion:  it.IdCuentaAtencion,
			IdServicioReceta:  it.IdServicioReceta,
			IdPrimeraAtencion: it.IdPrimeraAtencion,
			IdProducto:        it.IdProducto,
			Cantidad:          it.Cantidad,
			Precio:            it.Precio,
			Total:             it.Total,
			Codigo:            it.Codigo,
			Nombre:            it.Nombre,
			IdDiagnostico:     it.IdDiagnostico,
			CodigoDiagnostico: it.CodigoDiagnostico,
		})
	}
	respondSuccess(c, http.StatusOK, resp)
}

// EliminarExamen DELETE .../examenes/:idReceta/:idProducto?idCuentaAtencion=&idPrimeraAtencion=
func (h *PrimeraAtencionHandler) EliminarExamen(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}
	idReceta, err := strconv.Atoi(c.Param("idReceta"))
	if err != nil || idReceta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idReceta inválido")
		return
	}
	idProducto, err := strconv.Atoi(c.Param("idProducto"))
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idProducto inválido")
		return
	}
	idCuenta, err := strconv.Atoi(c.Query("idCuentaAtencion"))
	if err != nil || idCuenta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idPrimera, err := strconv.Atoi(c.Query("idPrimeraAtencion"))
	if err != nil || idPrimera <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion es obligatorio")
		return
	}

	if err := h.service.EliminarExamen(c.Request.Context(), idAtencion, idReceta, idProducto, idCuenta, idPrimera); err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_EXA_DELETE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idReceta":   idReceta,
		"idProducto": idProducto,
		"mensaje":    "Examen eliminado",
	})
}

type procedimientoRegistrarRequest struct {
	IdPrimeraAtencion  int     `json:"idPrimeraAtencion" binding:"required"`
	IdCuentaAtencion   int     `json:"idCuentaAtencion" binding:"required"`
	IdServicioReceta   int     `json:"idServicioReceta" binding:"required"`
	IdPaciente         int     `json:"idPaciente" binding:"required"`
	IdUsuarioAuditoria int     `json:"idUsuarioAuditoria"`
	IdProducto         int     `json:"idProducto" binding:"required"`
	Precio             float64 `json:"precio"`
	IdDiagnostico      int     `json:"idDiagnostico" binding:"required"`
	CodigoDiagnostico  string  `json:"codigoDiagnostico" binding:"required"`
}

// RegistrarProcedimiento POST .../procedimientos
// Cabecera fija IdPuntoCarga=26; detalle cant=1 + DX.
func (h *PrimeraAtencionHandler) RegistrarProcedimiento(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	var req procedimientoRegistrarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	res, err := h.service.RegistrarProcedimiento(c.Request.Context(), domain.RecetaProcedimientoRegistrar{
		IdAtencion:         idAtencion,
		IdCuentaAtencion:   req.IdCuentaAtencion,
		IdServicioReceta:   req.IdServicioReceta,
		IdPaciente:         req.IdPaciente,
		IdUsuarioAuditoria: req.IdUsuarioAuditoria,
		IdPrimeraAtencion:  req.IdPrimeraAtencion,
		IdProducto:         req.IdProducto,
		Precio:             req.Precio,
		IdDiagnostico:      req.IdDiagnostico,
		CodigoDiagnostico:  req.CodigoDiagnostico,
	})
	if err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_PRO_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idReceta":          res.IdReceta,
		"idProducto":        res.IdProducto,
		"cantidad":          res.Cantidad,
		"precio":            res.Precio,
		"total":             res.Total,
		"idPuntoCarga":      domain.IdPuntoCargaProcedimientos,
		"idDiagnostico":     req.IdDiagnostico,
		"codigoDiagnostico": req.CodigoDiagnostico,
		"idAtencion":        idAtencion,
		"mensaje":           res.Mensaje,
	})
}

// ListarProcedimientos GET .../procedimientos?idCuentaAtencion=&idPrimeraAtencion=
// Filtra RecetaCabecera.IdPuntoCarga = 26.
func (h *PrimeraAtencionHandler) ListarProcedimientos(c *gin.Context) {
	if _, err := strconv.Atoi(c.Param("idAtencion")); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}

	idCuenta, err := strconv.Atoi(c.Query("idCuentaAtencion"))
	if err != nil || idCuenta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idPrimera := 0
	if raw := c.Query("idPrimeraAtencion"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion inválido")
			return
		}
		idPrimera = v
	}

	items, err := h.service.ListarProcedimientos(c.Request.Context(), idCuenta, idPrimera)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "PRIMERA_ATENCION_PRO_LIST_FAILED", err.Error())
		return
	}

	resp := make([]examenListItemResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, examenListItemResponse{
			IdReceta:          it.IdReceta,
			IdPuntoCarga:      it.IdPuntoCarga,
			PuntoCarga:        it.PuntoCarga,
			IdCuentaAtencion:  it.IdCuentaAtencion,
			IdServicioReceta:  it.IdServicioReceta,
			IdPrimeraAtencion: it.IdPrimeraAtencion,
			IdProducto:        it.IdProducto,
			Cantidad:          it.Cantidad,
			Precio:            it.Precio,
			Total:             it.Total,
			Codigo:            it.Codigo,
			Nombre:            it.Nombre,
			IdDiagnostico:     it.IdDiagnostico,
			CodigoDiagnostico: it.CodigoDiagnostico,
		})
	}
	respondSuccess(c, http.StatusOK, resp)
}

// EliminarProcedimiento DELETE .../procedimientos/:idReceta/:idProducto
func (h *PrimeraAtencionHandler) EliminarProcedimiento(c *gin.Context) {
	idAtencion, err := strconv.Atoi(c.Param("idAtencion"))
	if err != nil || idAtencion <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion inválido")
		return
	}
	idReceta, err := strconv.Atoi(c.Param("idReceta"))
	if err != nil || idReceta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idReceta inválido")
		return
	}
	idProducto, err := strconv.Atoi(c.Param("idProducto"))
	if err != nil || idProducto <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idProducto inválido")
		return
	}
	idCuenta, err := strconv.Atoi(c.Query("idCuentaAtencion"))
	if err != nil || idCuenta <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	idPrimera, err := strconv.Atoi(c.Query("idPrimeraAtencion"))
	if err != nil || idPrimera <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPrimeraAtencion es obligatorio")
		return
	}

	if err := h.service.EliminarProcedimiento(c.Request.Context(), idAtencion, idReceta, idProducto, idCuenta, idPrimera); err != nil {
		respondError(c, http.StatusBadRequest, "PRIMERA_ATENCION_PRO_DELETE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"idReceta":   idReceta,
		"idProducto": idProducto,
		"mensaje":    "Procedimiento eliminado",
	})
}

