package httpadapter

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

// content type del 7z que el Firmador exige en la firma por lote.
const contentType7z = "application/x-7z-compressed"

// maxFirmaDocumentSize limita el tamaño de los documentos del proceso (25 MB).
const maxFirmaDocumentSize = 25 << 20

// FirmaPeruHandler expone el puerto de entrada input.FirmaPeruService.
// publicBaseURL, si se configura (FIRMAPERU_PUBLIC_URL), se usa como base para
// construir documentToSign/uploadDocumentSigned; si es vacío se deriva del
// Host de cada request.
type FirmaPeruHandler struct {
	service       input.FirmaPeruService
	publicBaseURL string
}

// NewFirmaPeruHandler inyecta el caso de uso de Firma Perú en el adaptador
// HTTP y la URL pública opcional para construir los callbacks del Firmador.
func NewFirmaPeruHandler(service input.FirmaPeruService, publicBaseURL string) *FirmaPeruHandler {
	return &FirmaPeruHandler{service: service, publicBaseURL: strings.TrimSuffix(publicBaseURL, "/")}
}

// Firmar maneja POST /api/v1/firmaperu/firmar: recibe el documento y los
// parámetros de firma y devuelve el param base64 para firmaperu.min.js.
//
// @Summary Inicia un proceso de firma digital con el Firmador Web
// @Description Almacena el documento por firmar y devuelve la cadena base64 que el frontend debe pasar a startSignature(port, param) de firmaperu.min.js. El Firmador local descargará el documento (documentToSign), el usuario firmará con su DNIe y devolverá el firmado a uploadDocumentSigned (ambas URLs internas al param).
// @Tags FirmaPeru
// @Accept multipart/form-data
// @Produce json
// @Param document formData file true "Documento a firmar (PDF para PAdES, XML para XAdES, cualquier archivo para CAdES)"
// @Param imageEstampado formData string false "Imagen PNG (base64 o data-URL) del estampado visible de la firma; por defecto se usa el logo por omisión"
// @Param signatureFormat formData string false "Formato de firma: PAdES (default) | XAdES | CAdES"
// @Param signatureLevel formData string false "Nivel de firma: B (default) | T | LTA"
// @Param signaturePackaging formData string false "Empaquetado: enveloped (default) | enveloping | detached | internallydetached"
// @Param webTsa formData string false "URL del servicio de sello de tiempo TSA (obligatorio si signatureLevel=T/LTA)"
// @Param userTsa formData string false "Usuario de la TSA"
// @Param passwordTsa formData string false "Password de la TSA"
// @Param contactInfo formData string false "Identificador de la firma en el documento"
// @Param signatureReason formData string false "Motivo de la firma"
// @Param imageToStamp formData string false "URL de la imagen PNG del estampado a usar en la firma visible"
// @Param signatureStyle formData int false "Estilo de la representación gráfica: 0 sin representación, 1 horizontal (default), 2 vertical, 3 solo estampado, 4 solo descripción"
// @Param stampTextSize formData int false "Tamaño de fuente del estampado (default 14)"
// @Param stampWordWrap formData int false "Ancho de línea del estampado (default 37)"
// @Param stampPage formData int false "Página donde se coloca la firma visible (default 1)"
// @Param positionx formData int false "Posición horizontal de la firma"
// @Param positiony formData int false "Posición vertical de la firma"
// @Param visiblePosition formData bool false "Muestra el panel de posición de la firma en el Firmador (default false)"
// @Param oneByOne formData bool false "Firma los documentos uno a uno (default false)"
// @Param role formData string false "Rol del firmante"
// @Param certificationSignature formData bool false "Firma de certificación (limita cambios posteriores)"
// @Success 200 {object} apiResponse{data=firmaPeruInitResponse}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros o documento inválidos"
// @Failure 413 {object} apiResponse{error=apiError} "Documento excede el tamaño máximo"
// @Failure 502 {object} apiResponse{error=apiError} "Error al obtener token o preparar el proceso"
// @Router /firmaperu/firmar [post]
func (h *FirmaPeruHandler) Firmar(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFirmaDocumentSize+1<<20)

	documentBytes, documentName, err := leerFormFile(c, "document")
	if err != nil {
		respondError(c, http.StatusBadRequest, "FIRMAPERU_DOCUMENT_REQUIRED", "el campo 'document' es obligatorio: "+err.Error())
		return
	}
	if len(documentBytes) > maxFirmaDocumentSize {
		respondError(c, http.StatusRequestEntityTooLarge, "FIRMAPERU_DOCUMENT_TOO_LARGE", "el documento excede el límite de 25 MB")
		return
	}

	params := firmaParamsDesdeForm(c)

	stampImage, err := leerImagenEstampado(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, "FIRMAPERU_STAMP_IMAGE_INVALID", "la imagen del estampado es inválida: "+err.Error())
		return
	}

	result, err := h.service.IniciarFirma(c.Request.Context(), domain.FirmaPeruInitRequest{
		Document:      documentBytes,
		DocumentName:  documentName,
		PublicBaseURL: h.baseURLPublica(c),
		StampImage:    stampImage,
		Params:        params,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFirmaPeruDocumentRequired),
			errors.Is(err, domain.ErrInvalidFirmaPeruFormat),
			errors.Is(err, domain.ErrInvalidFirmaPeruLevel):
			respondError(c, http.StatusBadRequest, "INVALID_FIRMAPERU_REQUEST", err.Error())
		default:
			respondError(c, http.StatusBadGateway, "FIRMAPERU_UNAVAILABLE", err.Error())
		}
		return
	}

	log.Printf("FirmaPeru: iniciada firma uuid=%s base_url=%s", result.DocumentNameUUID, h.baseURLPublica(c))
	respondSuccess(c, http.StatusOK, firmaPeruInitResponse{
		DocumentNameUUID: result.DocumentNameUUID,
		ParamBase64:      result.ParamBase64,
	})
}

// ParametrosFirma maneja POST /api/v1/firmaperu/params/:uuid: valida el
// param_token de único uso y devuelve el JSON base64 con los parámetros de
// firma. Es la URL interna del wrapper que publica el Firmador local.
//
// @Summary Devuelve los parámetros de firma al Firmador local
// @Description Endpoint interno expuesto en param_url del wrapper: el Firmador local lo llama (form-urlencoded) con el param_token de un solo uso y responde con el JSON base64 de los parámetros de firma (documentToSign, uploadDocumentSigned, token, etc.).
// @Tags FirmaPeru
// @Accept application/x-www-form-urlencoded
// @Produce text/plain
// @Param uuid path string true "UUID del proceso de firma"
// @Param param_token formData string true "Token de un único uso entregado en el wrapper"
// @Success 200 {string} string "JSON base64 con los parámetros de firma"
// @Failure 400 {object} apiResponse{error=apiError} "param_token inválido"
// @Failure 404 {object} apiResponse{error=apiError} "Proceso de firma no encontrado"
// @Failure 502 {object} apiResponse{error=apiError} "Error al obtener el token de la plataforma"
// @Router /firmaperu/params/{uuid} [post]
func (h *FirmaPeruHandler) ParametrosFirma(c *gin.Context) {
	paramToken := c.PostForm("param_token")
	paramBase64, err := h.service.ObtenerParametrosFirma(c.Request.Context(), c.Param("uuid"), paramToken)
	if err != nil {
		log.Printf("FirmaPeru: params uuid=%s error=%v (token_valido=false)", c.Param("uuid"), err)
		switch {
		case errors.Is(err, domain.ErrFirmaPeruDocumentNotFound):
			respondError(c, http.StatusNotFound, "FIRMAPERU_DOCUMENT_NOT_FOUND", err.Error())
		case errors.Is(err, domain.ErrFirmaPeruInvalidParamToken):
			respondError(c, http.StatusBadRequest, "FIRMAPERU_INVALID_PARAM_TOKEN", err.Error())
		default:
			respondError(c, http.StatusBadGateway, "FIRMAPERU_UNAVAILABLE", err.Error())
		}
		return
	}

	log.Printf("FirmaPeru: params uuid=%s entregados al Firmador (token_valido=true, base64=%d chars)", c.Param("uuid"), len(paramBase64))
	c.Data(http.StatusOK, "text/plain", []byte(paramBase64))
}

// DescargarDocumento maneja GET /api/v1/firmaperu/documentos/:uuid: sirve el
// documento original que el Firmador descargará (documentToSign).
//
// @Summary Descarga el documento original por firmar
// @Description URL interna usada por el Firmador (documentToSign) para descargar el documento que será firmado. Requiere el UUID devuelto en POST /firmaperu/firmar.
// @Tags FirmaPeru
// @Produce octet-stream
// @Param uuid path string true "UUID del proceso de firma"
// @Success 200 {file} binary "Documento original"
// @Failure 404 {object} apiResponse{error=apiError} "Documento no encontrado"
// @Router /firmaperu/documentos/{uuid} [get]
func (h *FirmaPeruHandler) DescargarDocumento(c *gin.Context) {
	doc, err := h.service.ObtenerDocumento(c.Request.Context(), c.Param("uuid"))
	if err != nil {
		respondError(c, http.StatusNotFound, "FIRMAPERU_DOCUMENT_NOT_FOUND", err.Error())
		return
	}

	c.Data(http.StatusOK, contentTypePorNombre(doc.Name), doc.Content)
}

// RecibirDocumentoFirmado maneja POST /api/v1/firmaperu/documentos/:uuid:
// guarda el documento firmado que el Firmador envía (uploadDocumentSigned).
//
// @Summary Recibe el documento firmado del Firmador
// @Description URL interna a la que el Firmador (uploadDocumentSigned) devuelve el documento firmado. Acepta el archivo como multipart (campo file) o como cuerpo plano.
// @Tags FirmaPeru
// @Accept multipart/form-data
// @Param uuid path string true "UUID del proceso de firma"
// @Param file formData file false "Documento firmado (multipart)"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse{error=apiError} "Documento vacío o inválido"
// @Failure 404 {object} apiResponse{error=apiError} "Proceso de firma no encontrado"
// @Router /firmaperu/documentos/{uuid} [post]
func (h *FirmaPeruHandler) RecibirDocumentoFirmado(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFirmaDocumentSize+1<<20)

	content, err := leerCuerpoDocumento(c)
	if err != nil {
		log.Printf("FirmaPeru: recibo firma uuid=%s error al leer cuerpo: %v", c.Param("uuid"), err)
		respondError(c, http.StatusBadRequest, "FIRMAPERU_DOCUMENT_INVALID", "no se pudo leer el documento firmado: "+err.Error())
		return
	}

	if err := h.service.RecibirDocumentoFirmado(c.Request.Context(), c.Param("uuid"), content); err != nil {
		log.Printf("FirmaPeru: recibo firma uuid=%s error=%v bytes=%d", c.Param("uuid"), err, len(content))
		switch {
		case errors.Is(err, domain.ErrFirmaPeruDocumentNotFound):
			respondError(c, http.StatusNotFound, "FIRMAPERU_DOCUMENT_NOT_FOUND", err.Error())
		default:
			respondError(c, http.StatusBadRequest, "INVALID_FIRMAPERU_REQUEST", err.Error())
		}
		return
	}

	log.Printf("FirmaPeru: documento firmado recibido uuid=%s bytes=%d content_type=%s", c.Param("uuid"), len(content), c.GetHeader("Content-Type"))
	respondSuccess(c, http.StatusOK, nil)
}

// DescargarDocumentoFirmado maneja GET /api/v1/firmaperu/documentos/:uuid/firmado:
// devuelve el documento firmado cuando el Firmador ya lo subió.
//
// @Summary Descarga el documento firmado
// @Description Devuelve el documento ya firmado (PAdES/XAdES/CAdES) o 404 si el Firmador aún no lo subió.
// @Tags FirmaPeru
// @Produce octet-stream
// @Param uuid path string true "UUID del proceso de firma"
// @Success 200 {file} binary "Documento firmado"
// @Failure 404 {object} apiResponse{error=apiError} "Documento no disponible aún o no encontrado"
// @Router /firmaperu/documentos/{uuid}/firmado [get]
func (h *FirmaPeruHandler) DescargarDocumentoFirmado(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	signed, err := h.service.ObtenerDocumentoFirmado(c.Request.Context(), c.Param("uuid"))
	if err != nil {
		log.Printf("FirmaPeru: descarga firmado uuid=%s aun no disponible: %v", c.Param("uuid"), err)
		respondError(c, http.StatusNotFound, "FIRMAPERU_SIGNED_NOT_AVAILABLE", err.Error())
		return
	}

	log.Printf("FirmaPeru: descarga firmado uuid=%s bytes=%d", c.Param("uuid"), len(signed))
	c.Data(http.StatusOK, "application/pdf", signed)
}

// FirmarLote maneja POST /api/v1/firmaperu/lote: recibe varios documentos
// (campos 'document') y los parámetros de firma, y devuelve el param base64
// para startSignature. El Firmador firmará todos los documentos con un solo
// PIN del DNIe (firma por lote, bachtOperation=true).
//
// @Summary Inicia una firma por lote (varios documentos, un solo PIN)
// @Description Recibe múltiples documentos (campo 'document') y devuelve el param base64 para startSignature. El Firmador descargará un 7z con todos los documentos (documentToSign) y devolverá el 7z firmado (uploadDocumentSigned); el usuario firma todos con un solo PIN del DNIe.
// @Tags FirmaPeru
// @Accept multipart/form-data
// @Produce json
// @Param document formData file true "Documentos a firmar (repetir el campo por cada PDF)"
// @Param imageEstampado formData string false "Imagen PNG (base64 o data-URL) del estampado visible de la firma"
// @Param signatureFormat formData string false "Formato de firma: PAdES (default)"
// @Param signatureLevel formData string false "Nivel de firma: B (default) | T | LTA"
// @Param signatureReason formData string false "Motivo de la firma"
// @Param signatureStyle formData int false "Estilo de la representación gráfica (default 1)"
// @Param stampTextSize formData int false "Tamaño de fuente del estampado (default 14)"
// @Param stampWordWrap formData int false "Ancho de línea del estampado (default 37)"
// @Param stampPage formData int false "Página donde se coloca la firma visible (default 1)"
// @Param positionx formData int false "Posición horizontal de la firma"
// @Param positiony formData int false "Posición vertical de la firma"
// @Param role formData string false "Rol del firmante"
// @Success 200 {object} apiResponse{data=firmaPeruInitResponse}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros o documentos inválidos"
// @Failure 413 {object} apiResponse{error=apiError} "Documento excede el tamaño máximo"
// @Failure 502 {object} apiResponse{error=apiError} "Error al obtener token o preparar el lote"
// @Router /firmaperu/lote [post]
func (h *FirmaPeruHandler) FirmarLote(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFirmaDocumentSize*30+1<<20)

	documentos, err := leerFormFiles(c, "document")
	if err != nil {
		respondError(c, http.StatusBadRequest, "FIRMAPERU_DOCUMENT_REQUIRED", "el campo 'document' es obligatorio: "+err.Error())
		return
	}
	if len(documentos) == 0 {
		respondError(c, http.StatusBadRequest, "FIRMAPERU_DOCUMENT_REQUIRED", "no se recibieron documentos para firmar")
		return
	}
	for _, d := range documentos {
		if len(d.Content) > maxFirmaDocumentSize {
			respondError(c, http.StatusRequestEntityTooLarge, "FIRMAPERU_DOCUMENT_TOO_LARGE", "un documento excede el límite de 25 MB")
			return
		}
	}

	stampImage, err := leerImagenEstampado(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, "FIRMAPERU_STAMP_IMAGE_INVALID", "la imagen del estampado es inválida: "+err.Error())
		return
	}

	result, err := h.service.IniciarFirmaLote(c.Request.Context(), domain.FirmaPeruLoteInitRequest{
		Documents:     documentos,
		PublicBaseURL: h.baseURLPublica(c),
		StampImage:    stampImage,
		Params:        firmaParamsDesdeForm(c),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFirmaPeruDocumentRequired),
			errors.Is(err, domain.ErrInvalidFirmaPeruFormat),
			errors.Is(err, domain.ErrInvalidFirmaPeruLevel):
			respondError(c, http.StatusBadRequest, "INVALID_FIRMAPERU_REQUEST", err.Error())
		default:
			respondError(c, http.StatusBadGateway, "FIRMAPERU_UNAVAILABLE", err.Error())
		}
		return
	}

	log.Printf("FirmaPeru: lote iniciado uuid=%s documentos=%d base_url=%s", result.DocumentNameUUID, len(documentos), h.baseURLPublica(c))
	respondSuccess(c, http.StatusOK, firmaPeruInitResponse{
		DocumentNameUUID: result.DocumentNameUUID,
		ParamBase64:      result.ParamBase64,
	})
}

// DescargarLoteDocumentos maneja GET /api/v1/firmaperu/documentos/:uuid/lote:
// devuelve el archivo 7z con todos los documentos del lote que el Firmador
// descargará (documentToSign en modo lote).
//
// @Summary Descarga el 7z con los documentos del lote
// @Description URL interna usada por el Firmador en modo lote (documentToSign): descarga un archivo 7z con todos los documentos del lote por firmar.
// @Tags FirmaPeru
// @Produce octet-stream
// @Param uuid path string true "UUID del proceso de firma"
// @Success 200 {file} binary "Archivo 7z con los documentos"
// @Failure 404 {object} apiResponse{error=apiError} "Lote no encontrado"
// @Router /firmaperu/documentos/{uuid}/lote [get]
func (h *FirmaPeruHandler) DescargarLoteDocumentos(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	archivo, err := h.service.ObtenerDocumentosLote(c.Request.Context(), c.Param("uuid"))
	if err != nil {
		respondError(c, http.StatusNotFound, "FIRMAPERU_DOCUMENT_NOT_FOUND", err.Error())
		return
	}
	log.Printf("FirmaPeru: lote uuid=%s 7z descargado bytes=%d", c.Param("uuid"), len(archivo))
	c.Header("Content-Disposition", `attachment; filename="lote.7z"`)
	c.Data(http.StatusOK, contentType7z, archivo)
}

// RecibirLoteFirmado maneja POST /api/v1/firmaperu/documentos/:uuid/lote:
// recibe el 7z firmado que el Firmador devuelve (uploadDocumentSigned en modo
// lote), lo extrae y guarda cada PDF firmado en FIRMAPERU_SIGNED_DIR.
//
// @Summary Recibe el 7z firmado del lote
// @Description URL interna a la que el Firmador (uploadDocumentSigned) devuelve el 7z con los documentos firmados del lote. Acepta el archivo como multipart (campo signed_file) o como cuerpo plano.
// @Tags FirmaPeru
// @Accept multipart/form-data
// @Param uuid path string true "UUID del proceso de firma"
// @Param signed_file formData file false "Archivo 7z con los documentos firmados (multipart)"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse{error=apiError} "Archivo vacío o inválido"
// @Failure 404 {object} apiResponse{error=apiError} "Lote no encontrado"
// @Router /firmaperu/documentos/{uuid}/lote [post]
func (h *FirmaPeruHandler) RecibirLoteFirmado(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFirmaDocumentSize*30+1<<20)

	content, err := leerCuerpoDocumento(c)
	if err != nil {
		log.Printf("FirmaPeru: recibo lote firmado uuid=%s error al leer cuerpo: %v", c.Param("uuid"), err)
		respondError(c, http.StatusBadRequest, "FIRMAPERU_DOCUMENT_INVALID", "no se pudo leer el archivo 7z firmado: "+err.Error())
		return
	}

	if err := h.service.RecibirDocumentoFirmadoLote(c.Request.Context(), c.Param("uuid"), content); err != nil {
		log.Printf("FirmaPeru: recibo lote firmado uuid=%s error=%v bytes=%d", c.Param("uuid"), err, len(content))
		switch {
		case errors.Is(err, domain.ErrFirmaPeruDocumentNotFound):
			respondError(c, http.StatusNotFound, "FIRMAPERU_DOCUMENT_NOT_FOUND", err.Error())
		default:
			respondError(c, http.StatusBadRequest, "INVALID_FIRMAPERU_REQUEST", err.Error())
		}
		return
	}

	log.Printf("FirmaPeru: lote firmado recibido uuid=%s bytes=%d", c.Param("uuid"), len(content))
	respondSuccess(c, http.StatusOK, nil)
}

// DescargarLoteFirmado maneja GET /api/v1/firmaperu/documentos/:uuid/lote/firmado:
// devuelve el 7z firmado del lote cuando el Firmador ya lo subió (lo usa el
// polling del frontend para confirmar que el lote terminó).
//
// @Summary Descarga el 7z firmado del lote
// @Description Devuelve el archivo 7z con los documentos ya firmados del lote o 404 si el Firmador aún no lo subió.
// @Tags FirmaPeru
// @Produce octet-stream
// @Param uuid path string true "UUID del proceso de firma"
// @Success 200 {file} binary "Archivo 7z firmado"
// @Failure 404 {object} apiResponse{error=apiError} "Lote no disponible aún"
// @Router /firmaperu/documentos/{uuid}/lote/firmado [get]
func (h *FirmaPeruHandler) DescargarLoteFirmado(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	archivo, err := h.service.ObtenerLoteFirmado(c.Request.Context(), c.Param("uuid"))
	if err != nil {
		log.Printf("FirmaPeru: descarga lote firmado uuid=%s aun no disponible: %v", c.Param("uuid"), err)
		respondError(c, http.StatusNotFound, "FIRMAPERU_SIGNED_NOT_AVAILABLE", err.Error())
		return
	}
	log.Printf("FirmaPeru: descarga lote firmado uuid=%s bytes=%d", c.Param("uuid"), len(archivo))
	c.Data(http.StatusOK, contentType7z, archivo)
}

type archivoFirmadoItemRespuesta struct {
	Nombre string `json:"nombre"`
	Tamano int    `json:"tamano"`
}

func (h *FirmaPeruHandler) ListarLoteArchivosFirmados(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	documentos, err := h.service.ObtenerLoteDocumentosFirmados(c.Request.Context(), c.Param("uuid"))
	if err != nil {
		respondError(c, http.StatusNotFound, "FIRMAPERU_SIGNED_NOT_AVAILABLE", err.Error())
		return
	}

	lista := make([]archivoFirmadoItemRespuesta, 0, len(documentos))
	for _, doc := range documentos {
		lista = append(lista, archivoFirmadoItemRespuesta{
			Nombre: doc.Name,
			Tamano: len(doc.Content),
		})
	}
	respondSuccess(c, http.StatusOK, lista)
}

func (h *FirmaPeruHandler) DescargarLoteArchivoFirmado(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	nombre := c.Param("nombre")
	contenido, err := h.service.ObtenerLoteArchivoFirmado(c.Request.Context(), c.Param("uuid"), nombre)
	if err != nil {
		respondError(c, http.StatusNotFound, "FIRMAPERU_DOCUMENT_NOT_FOUND", err.Error())
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filepath.Base(nombre)))
	c.Data(http.StatusOK, contentTypePorNombre(nombre), contenido)
}

// baseURLPublica devuelve la URL pública configurada o la deriva del Host del
// request (respetando el header X-Forwarded-Proto).
func (h *FirmaPeruHandler) baseURLPublica(c *gin.Context) string {
	if h.publicBaseURL != "" {
		return h.publicBaseURL
	}
	scheme := "http"
	if proto := c.GetHeader("X-Forwarded-Proto"); proto == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// leerFormFile extrae y lee en memoria el archivo de un campo del form-data.
// Devuelve el contenido y el nombre original del archivo. Si el campo no
// existe devuelve un error descriptivo.
func leerFormFile(c *gin.Context, campo string) ([]byte, string, error) {
	fileHeader, err := c.FormFile(campo)
	if err != nil {
		return nil, "", err
	}
	f, err := fileHeader.Open()
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	contenido, err := io.ReadAll(io.LimitReader(f, maxFirmaDocumentSize+1<<20))
	if err != nil {
		return nil, "", err
	}
	return contenido, fileHeader.Filename, nil
}

// leerFormFiles extrae y lee en memoria todos los archivos de un campo de
// form-data (usado por la firma por lote, que envía varios 'document').
func leerFormFiles(c *gin.Context, campo string) ([]domain.FirmaPeruDocument, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	fileHeaders := form.File[campo]
	if len(fileHeaders) == 0 {
		return nil, errors.New("campo '" + campo + "' sin archivos")
	}
	docs := make([]domain.FirmaPeruDocument, 0, len(fileHeaders))
	for _, fh := range fileHeaders {
		f, err := fh.Open()
		if err != nil {
			return nil, err
		}
		contenido, readErr := io.ReadAll(io.LimitReader(f, maxFirmaDocumentSize+1<<20))
		f.Close()
		if readErr != nil {
			return nil, readErr
		}
		docs = append(docs, domain.FirmaPeruDocument{Name: fh.Filename, Content: contenido})
	}
	return docs, nil
}

// leerImagenEstampado lee la imagen PNG del estampado enviada por el frontend:
// como archivo (campo imageEstampado) o como campo de texto con la imagen en
// base64 (con o sin prefijo data:image/png;base64,). Devuelve nil si el
// frontend no envió imagen propia (se usará el estampado por defecto).
func leerImagenEstampado(c *gin.Context) ([]byte, error) {
	if fileHeader, err := c.FormFile("imageEstampado"); err == nil {
		f, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return io.ReadAll(io.LimitReader(f, 4<<20))
	}

	raw := strings.TrimSpace(c.PostForm("imageEstampado"))
	if raw == "" {
		return nil, nil
	}

	base64Data := raw
	if idx := strings.Index(base64Data, "base64,"); idx >= 0 {
		base64Data = base64Data[idx+len("base64,"):]
	}
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}
	if len(decoded) == 0 {
		return nil, nil
	}
	return decoded, nil
}

// leerCuerpoDocumento lee el documento firmado: si llega como multipart usa el
// primer archivo de cualquier campo (el Firmador no fija nombre de campo); si
// llega como cuerpo plano usa el body completo.
func leerCuerpoDocumento(c *gin.Context) ([]byte, error) {
	if strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
		form, err := c.MultipartForm()
		if err != nil {
			return nil, err
		}
		for _, headers := range form.File {
			if len(headers) == 0 {
				continue
			}
			f, err := headers[0].Open()
			if err != nil {
				return nil, err
			}
			defer f.Close()
			return io.ReadAll(io.LimitReader(f, maxFirmaDocumentSize+1<<20))
		}
	}
	return io.ReadAll(c.Request.Body)
}

// contentTypePorNombre infiere el Content-Type de un archivo según su extensión.
func contentTypePorNombre(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".pdf":
		return "application/pdf"
	case ".xml":
		return "application/xml"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

// firmaParamsDesdeForm lee los parámetros de firma enviados por el frontend
// en el multipart (form-data) y completa los valores por defecto.
func firmaParamsDesdeForm(c *gin.Context) domain.FirmaPeruSignParams {
	return domain.FirmaPeruSignParams{
		SignatureFormat:        c.PostForm("signatureFormat"),
		SignatureLevel:         c.PostForm("signatureLevel"),
		SignaturePackaging:     c.PostForm("signaturePackaging"),
		WebTsa:                 c.PostForm("webTsa"),
		UserTsa:                c.PostForm("userTsa"),
		PasswordTsa:            c.PostForm("passwordTsa"),
		ContactInfo:            c.PostForm("contactInfo"),
		SignatureReason:        c.PostForm("signatureReason"),
		SignatureStyle:         formIntOrDefault(c, "signatureStyle", -1),
		StampTextSize:          formIntOrDefault(c, "stampTextSize", 14),
		StampWordWrap:          formIntOrDefault(c, "stampWordWrap", 37),
		StampPage:              formIntOrDefault(c, "stampPage", 1),
		PositionX:              formIntOrDefault(c, "positionx", 0),
		PositionY:              formIntOrDefault(c, "positiony", 0),
		VisiblePosition:        formBoolValue(c, "visiblePosition"),
		OneByOne:               formBoolValue(c, "oneByOne"),
		Role:                   c.PostForm("role"),
		CertificationSignature: formBoolValue(c, "certificationSignature"),
		BatchOperation:         formBoolValue(c, "batchOperation"),
		CertificateFilter:      ".*",
		Theme:                  "claro",
		ImageToStamp:           c.PostForm("imageToStamp"),
	}
}

// formIntOrDefault parsea un campo de formulario como int con valor por defecto.
func formIntOrDefault(c *gin.Context, campo string, def int) int {
	raw := c.PostForm(campo)
	if raw == "" {
		return def
	}
	valor, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return valor
}

// formBoolValue parsea un campo de formulario como bool; false si está vacío.
func formBoolValue(c *gin.Context, campo string) bool {
	raw := c.PostForm(campo)
	if raw == "" {
		return false
	}
	valor, err := strconv.ParseBool(raw)
	if err != nil {
		return false
	}
	return valor
}

// firmaPeruInitResponse es la respuesta de POST /firmaperu/firmar.
type firmaPeruInitResponse struct {
	DocumentNameUUID string `json:"documentNameUUID"`
	ParamBase64      string `json:"paramBase64"`
}
