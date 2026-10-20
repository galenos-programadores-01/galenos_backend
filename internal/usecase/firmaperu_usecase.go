package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type firmaPeruUseCase struct {
	client  output.FirmaPeruClient
	store   output.FirmaPeruStore
	archive output.SevenZipArchive
}

// NewFirmaPeruUseCase construye el caso de uso de firma digital con el
// Firmador Web de Firma Perú.
func NewFirmaPeruUseCase(client output.FirmaPeruClient, store output.FirmaPeruStore, archive output.SevenZipArchive) input.FirmaPeruService {
	return &firmaPeruUseCase{client: client, store: store, archive: archive}
}

// IniciarFirma valida los parámetros, almacena el original y la solicitud, y
// construye el wrapper base64 (param_url/param_token/document_extension) que
// el frontend pasa a startSignature(port, param). El Firmador local consulta
// el endpoint param_url para obtener los parámetros de firma.
func (uc *firmaPeruUseCase) IniciarFirma(ctx context.Context, req domain.FirmaPeruInitRequest) (domain.FirmaPeruInitResult, error) {
	if len(req.Document) == 0 {
		return domain.FirmaPeruInitResult{}, domain.ErrFirmaPeruDocumentRequired
	}
	if !validarFormatoFirmaPeru(req.Params.SignatureFormat) {
		return domain.FirmaPeruInitResult{}, domain.ErrInvalidFirmaPeruFormat
	}
	if !validarNivelFirmaPeru(req.Params.SignatureLevel) {
		return domain.FirmaPeruInitResult{}, domain.ErrInvalidFirmaPeruLevel
	}

	uuid, err := nuevoUUID()
	if err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("generating document uuid: %w", err)
	}

	if err := uc.store.SaveDocument(ctx, uuid, domain.FirmaPeruDocument{
		Name:    req.DocumentName,
		Content: req.Document,
	}); err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("storing firma peru document: %w", err)
	}

	if len(req.StampImage) > 0 {
		if err := uc.store.SaveStampImage(ctx, uuid, req.StampImage); err != nil {
			return domain.FirmaPeruInitResult{}, fmt.Errorf("storing firma peru stamp image: %w", err)
		}
	}

	paramToken, err := nuevoUUID()
	if err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("generating firma peru param_token: %w", err)
	}

	if err := uc.store.SaveSolicitud(ctx, uuid, domain.FirmaPeruSolicitud{
		ParamToken:    paramToken,
		PublicBaseURL: strings.TrimSuffix(req.PublicBaseURL, "/"),
		Params:        aplicarDefaults(req.Params),
	}); err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("storing firma peru solicitud: %w", err)
	}

	baseURL := strings.TrimSuffix(req.PublicBaseURL, "/")
	wrapper := domain.FirmaPeruParamWrapper{
		ParamURL:          baseURL + "/api/v1/firmaperu/params/" + uuid,
		ParamToken:        paramToken,
		DocumentExtension: extensionDocumento(req.DocumentName),
	}

	wrapperBase64, err := uc.client.BuildWrapperBase64(wrapper)
	if err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("building firma peru wrapper param: %w", err)
	}

	return domain.FirmaPeruInitResult{
		ParamBase64:      wrapperBase64,
		DocumentNameUUID: uuid,
	}, nil
}

// ObtenerParametrosFirma valida el param_token de único uso y devuelve el JSON
// base64 con los parámetros de firma (documentToSign, uploadDocumentSigned y
// el JWT de la plataforma). Es el endpoint al que llama el Firmador local
// (param_url).
func (uc *firmaPeruUseCase) ObtenerParametrosFirma(ctx context.Context, uuid string, paramToken string) (string, error) {
	sol, err := uc.store.GetSolicitud(ctx, uuid)
	if err != nil {
		return "", err
	}
	if paramToken == "" || sol.ParamToken != paramToken {
		return "", domain.ErrFirmaPeruInvalidParamToken
	}

	token, err := uc.client.GetToken(ctx)
	if err != nil {
		return "", fmt.Errorf("getting firma peru token: %w", err)
	}

	params := sol.Params
	baseURL := sol.PublicBaseURL
	if params.BatchOperation {
		params.DocumentToSign = baseURL + "/api/v1/firmaperu/documentos/" + uuid + "/lote"
		params.UploadDocumentSigned = baseURL + "/api/v1/firmaperu/documentos/" + uuid + "/lote"
		params.OneByOne = false
	} else {
		params.DocumentToSign = baseURL + "/api/v1/firmaperu/documentos/" + uuid
		params.UploadDocumentSigned = baseURL + "/api/v1/firmaperu/documentos/" + uuid
	}
	params.Token = token
	if _, err := uc.store.GetStampImage(ctx, uuid); err == nil {
		params.ImageToStamp = baseURL + "/api/v1/firmaperu/estampado/" + uuid
	} else if params.ImageToStamp == "" {
		params.ImageToStamp = baseURL + "/api/v1/firmaperu/estampado.png"
	}

	paramBase64, err := uc.client.BuildParamBase64(params)
	if err != nil {
		return "", fmt.Errorf("building firma peru sign params: %w", err)
	}

	return paramBase64, nil
}

// ObtenerDocumento devuelve el documento original pendiente de firma.
func (uc *firmaPeruUseCase) ObtenerDocumento(ctx context.Context, uuid string) (domain.FirmaPeruDocument, error) {
	return uc.store.GetDocument(ctx, uuid)
}

// RecibirDocumentoFirmado almacena el documento firmado del callback.
func (uc *firmaPeruUseCase) RecibirDocumentoFirmado(ctx context.Context, uuid string, content []byte) error {
	if len(content) == 0 {
		return domain.ErrFirmaPeruDocumentRequired
	}
	return uc.store.SaveSigned(ctx, uuid, content)
}

// ObtenerDocumentoFirmado devuelve el documento ya firmado.
func (uc *firmaPeruUseCase) ObtenerDocumentoFirmado(ctx context.Context, uuid string) ([]byte, error) {
	return uc.store.GetSigned(ctx, uuid)
}

// ObtenerImagenEstampado devuelve la imagen del estampado del proceso
// (logo del hospital) si el frontend envió una propia.
func (uc *firmaPeruUseCase) ObtenerImagenEstampado(ctx context.Context, uuid string) ([]byte, error) {
	return uc.store.GetStampImage(ctx, uuid)
}

// IniciarFirmaLote valida los documentos y parámetros, almacena todos los
// documentos y la solicitud (bachtOperation=true, oneByOne=false), y devuelve
// el wrapper base64 para startSignature. El Firmador descargará el 7z con
// todos los documentos (ObtenerDocumentosLote) y devolverá el 7z firmado.
func (uc *firmaPeruUseCase) IniciarFirmaLote(ctx context.Context, req domain.FirmaPeruLoteInitRequest) (domain.FirmaPeruInitResult, error) {
	if len(req.Documents) == 0 {
		return domain.FirmaPeruInitResult{}, domain.ErrFirmaPeruDocumentRequired
	}
	if !validarFormatoFirmaPeru(req.Params.SignatureFormat) {
		return domain.FirmaPeruInitResult{}, domain.ErrInvalidFirmaPeruFormat
	}
	if !validarNivelFirmaPeru(req.Params.SignatureLevel) {
		return domain.FirmaPeruInitResult{}, domain.ErrInvalidFirmaPeruLevel
	}

	uuid, err := nuevoUUID()
	if err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("generating lote uuid: %w", err)
	}

	if err := uc.store.SaveLote(ctx, uuid, req.Documents); err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("storing firma peru lote: %w", err)
	}

	if len(req.StampImage) > 0 {
		if err := uc.store.SaveStampImage(ctx, uuid, req.StampImage); err != nil {
			return domain.FirmaPeruInitResult{}, fmt.Errorf("storing firma peru stamp image: %w", err)
		}
	}

	paramToken, err := nuevoUUID()
	if err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("generating firma peru param_token: %w", err)
	}

	params := aplicarDefaults(req.Params)
	params.BatchOperation = true
	params.OneByOne = false

	if err := uc.store.SaveSolicitud(ctx, uuid, domain.FirmaPeruSolicitud{
		ParamToken:    paramToken,
		PublicBaseURL: strings.TrimSuffix(req.PublicBaseURL, "/"),
		Params:        params,
	}); err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("storing firma peru solicitud: %w", err)
	}

	baseURL := strings.TrimSuffix(req.PublicBaseURL, "/")
	wrapper := domain.FirmaPeruParamWrapper{
		ParamURL:          baseURL + "/api/v1/firmaperu/params/" + uuid,
		ParamToken:        paramToken,
		DocumentExtension: extensionDocumento(req.Documents[0].Name),
	}

	wrapperBase64, err := uc.client.BuildWrapperBase64(wrapper)
	if err != nil {
		return domain.FirmaPeruInitResult{}, fmt.Errorf("building firma peru wrapper param: %w", err)
	}

	return domain.FirmaPeruInitResult{
		ParamBase64:      wrapperBase64,
		DocumentNameUUID: uuid,
	}, nil
}

// ObtenerDocumentosLote construye y devuelve el 7z con todos los documentos
// del lote (documentToSign en modo lote).
func (uc *firmaPeruUseCase) ObtenerDocumentosLote(ctx context.Context, uuid string) ([]byte, error) {
	docs, err := uc.store.GetLote(ctx, uuid)
	if err != nil {
		return nil, err
	}
	archivo, err := uc.archive.Build7z(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("building firma peru lote 7z: %w", err)
	}
	return archivo, nil
}

// RecibirDocumentoFirmadoLote extrae el 7z firmado devuelto por el Firmador y
// guarda cada PDF en FIRMAPERU_SIGNED_DIR con el nombre del triaje original.
func (uc *firmaPeruUseCase) RecibirDocumentoFirmadoLote(ctx context.Context, uuid string, archive []byte) error {
	if len(archive) == 0 {
		return domain.ErrFirmaPeruDocumentRequired
	}
	firmados, err := uc.archive.Extract7z(ctx, archive)
	if err != nil {
		return fmt.Errorf("extracting firma peru lote 7z: %w", err)
	}
	originales, err := uc.store.GetLote(ctx, uuid)
	if err != nil {
		return err
	}

	targets := mapearFirmadosLote(originales, firmados)
	for _, f := range targets {
		if err := uc.store.SaveLoteSignedFile(ctx, uuid, f.Name, f.Content); err != nil {
			log.Printf("FirmaPeru: lote uuid=%s no se pudo guardar %s: %v", uuid, f.Name, err)
		}
	}
	if err := uc.store.SaveLoteSigned(ctx, uuid, archive); err != nil {
		return err
	}
	log.Printf("FirmaPeru: lote uuid=%s firmado: %d/%d documentos guardados", uuid, len(targets), len(originales))
	return nil
}

// ObtenerLoteFirmado devuelve el 7z firmado del lote cuando el Firmador lo
// subió (lo consume el polling del frontend para confirmar el fin del lote).
func (uc *firmaPeruUseCase) ObtenerLoteFirmado(ctx context.Context, uuid string) ([]byte, error) {
	return uc.store.GetLoteSigned(ctx, uuid)
}

// reTriajeID extrae el número de triaje de nombres tipo TRIAJE_123.pdf.
var reTriajeID = regexp.MustCompile(`(?i)TRIAJE[_ ]?(\d+)`)

// mapearFirmadosLote empareja cada archivo firmado del 7z con el original:
// 1) basename exacto (PA_*, TRIAJE_*), 2) id de triaje, 3) orden.
func mapearFirmadosLote(originales, firmados []domain.FirmaPeruDocument) []domain.FirmaPeruDocument {
	usado := make([]bool, len(originales))
	var out []domain.FirmaPeruDocument

	for _, f := range firmados {
		firmadoBase := strings.ToLower(filepath.Base(f.Name))
		idx := -1
		for i, o := range originales {
			if usado[i] {
				continue
			}
			if strings.ToLower(filepath.Base(o.Name)) == firmadoBase {
				idx = i
				break
			}
		}
		if idx < 0 {
			id := triajeID(f.Name)
			for i, o := range originales {
				if !usado[i] && id >= 0 && triajeID(o.Name) == id {
					idx = i
					break
				}
			}
		}
		if idx < 0 {
			for i := range originales {
				if !usado[i] {
					idx = i
					break
				}
			}
		}
		if idx < 0 {
			continue
		}
		usado[idx] = true
		out = append(out, domain.FirmaPeruDocument{
			Name:    filepath.Base(originales[idx].Name),
			Content: f.Content,
		})
	}
	return out
}

// triajeID devuelve el número de triaje contenido en un nombre de archivo o
// -1 si no se puede identificar.
func triajeID(nombre string) int {
	m := reTriajeID.FindStringSubmatch(nombre)
	if len(m) != 2 {
		return -1
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return -1
	}
	return n
}

// aplicarDefaults completa los parámetros con los valores que usa la
// integración de referencia del Firmador Web cuando el cliente no los envía.
func aplicarDefaults(p domain.FirmaPeruSignParams) domain.FirmaPeruSignParams {
	if p.SignatureFormat == "" {
		p.SignatureFormat = domain.FirmaPeruFormatPAdES
	}
	if p.SignatureLevel == "" {
		p.SignatureLevel = domain.FirmaPeruLevelB
	}
	if p.SignaturePackaging == "" {
		p.SignaturePackaging = "enveloped"
		if p.SignatureFormat == domain.FirmaPeruFormatCAdES {
			p.SignaturePackaging = "detached"
		}
	}
	if p.CertificateFilter == "" {
		p.CertificateFilter = ".*"
	}
	if p.Theme == "" {
		p.Theme = "claro"
	}
	if p.SignatureStyle == 0 {
		p.SignatureStyle = 1
	}
	if p.StampTextSize == 0 {
		p.StampTextSize = 14
	}
	if p.StampWordWrap == 0 {
		p.StampWordWrap = 37
	}
	if p.StampPage == 0 {
		p.StampPage = 1
	}
	return p
}

func validarFormatoFirmaPeru(formato string) bool {
	switch formato {
	case domain.FirmaPeruFormatPAdES, domain.FirmaPeruFormatXAdES, domain.FirmaPeruFormatCAdES, "":
		return true
	}
	return false
}

func validarNivelFirmaPeru(nivel string) bool {
	switch nivel {
	case domain.FirmaPeruLevelB, domain.FirmaPeruLevelT, domain.FirmaPeruLevelLTA, "":
		return true
	}
	return false
}

// nuevoUUID genera un identificador aleatorio de 32 caracteres hexadecimales.
func nuevoUUID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// extensionDocumento devuelve la extensión del documento sin el punto (p. ej.
// "pdf"); "pdf" si no se puede inferir.
func extensionDocumento(name string) string {
	name = strings.TrimSpace(name)
	idx := strings.LastIndex(name, ".")
	if idx < 0 || idx == len(name)-1 {
		return "pdf"
	}
	ext := name[idx+1:]
	if ext == "" {
		return "pdf"
	}
	return strings.ToLower(ext)
}
