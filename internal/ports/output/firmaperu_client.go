package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// FirmaPeruClient es el puerto de salida del Firmador Web de la Plataforma
// Firma Perú (Web Service / cliente web). GetToken obtiene el JWT del
// servicio generate-token; BuildParamBase64 serializa los parámetros de firma
// al JSON base64 que devuelve el endpoint param_url y BuildWrapperBase64
// serializa el wrapper (param_url/param_token/document_extension) que el
// frontend pasa a startSignature(port, param).
type FirmaPeruClient interface {
	// GetToken obtiene un token de acceso del servicio generate-token de la
	// plataforma usando client_id/client_secret.
	GetToken(ctx context.Context) (string, error)
	// BuildParamBase64 construye la cadena base64 del JSON de parámetros de
	// firma (objeto que espera el endpoint param_url).
	BuildParamBase64(params domain.FirmaPeruSignParams) (string, error)
	// BuildWrapperBase64 construye la cadena base64 del wrapper que el
	// frontend pasa a startSignature(port, param).
	BuildWrapperBase64(wrapper domain.FirmaPeruParamWrapper) (string, error)
}

// FirmaPeruStore es el puerto de salida para persistir temporalmente los
// documentos del proceso de firma (original y firmado), identificados por su
// UUID.
type FirmaPeruStore interface {
	// SaveDocument guarda el documento original por firmar.
	SaveDocument(ctx context.Context, uuid string, doc domain.FirmaPeruDocument) error
	// GetDocument recupera el documento original; devuelve
	// domain.ErrFirmaPeruDocumentNotFound si no existe.
	GetDocument(ctx context.Context, uuid string) (domain.FirmaPeruDocument, error)
	// SaveSigned guarda el documento ya firmado devuelto por el Firmador.
	SaveSigned(ctx context.Context, uuid string, content []byte) error
	// GetSigned recupera el documento firmado; devuelve
	// domain.ErrFirmaPeruSignedNotAvailable si aún no se firmó.
	GetSigned(ctx context.Context, uuid string) ([]byte, error)
	// SaveSolicitud guarda el estado del proceso de firma (param_token de uso
	// único, URL pública y parámetros de firma) asociado al UUID.
	SaveSolicitud(ctx context.Context, uuid string, sol domain.FirmaPeruSolicitud) error
	// GetSolicitud recupera la solicitud asociada al UUID; devuelve
	// domain.ErrFirmaPeruDocumentNotFound si no existe.
	GetSolicitud(ctx context.Context, uuid string) (domain.FirmaPeruSolicitud, error)
	// SaveStampImage guarda la imagen PNG del estampado (logo del hospital)
	// que el Firmador descargará para la representación visible de la firma.
	SaveStampImage(ctx context.Context, uuid string, content []byte) error
	// GetStampImage recupera la imagen del estampado; devuelve
	// domain.ErrFirmaPeruDocumentNotFound si no se definió imagen propia.
	GetStampImage(ctx context.Context, uuid string) ([]byte, error)
	// SaveLote guarda los documentos de un proceso de firma por lote (firma
	// masiva con un solo PIN), identificados por su UUID.
	SaveLote(ctx context.Context, uuid string, docs []domain.FirmaPeruDocument) error
	// GetLote recupera los documentos de un lote; devuelve
	// domain.ErrFirmaPeruDocumentNotFound si no existe.
	GetLote(ctx context.Context, uuid string) ([]domain.FirmaPeruDocument, error)
	// SaveLoteSigned guarda el archivo 7z firmado devuelto por el Firmador
	// para un lote (el Firmador lo sirve al frontend vía GET lote/firmado).
	SaveLoteSigned(ctx context.Context, uuid string, archive []byte) error
	// GetLoteSigned recupera el 7z firmado de un lote; devuelve
	// domain.ErrFirmaPeruSignedNotAvailable si aún no se firmó.
	GetLoteSigned(ctx context.Context, uuid string) ([]byte, error)
	// SaveLoteSignedFile guarda en memoria (y disco si hay FIRMAPERU_SIGNED_DIR)
	// un PDF firmado de un lote con su nombre original.
	SaveLoteSignedFile(ctx context.Context, uuid string, name string, content []byte) error
	// GetLoteSignedFiles recupera los PDFs firmados individuales del lote;
	// domain.ErrFirmaPeruSignedNotAvailable si aún no hay ninguno.
	GetLoteSignedFiles(ctx context.Context, uuid string) ([]domain.FirmaPeruDocument, error)
}

// SevenZipArchive es el puerto de salida para crear y extraer archivos 7z
// (formato que exige el Firmador en la firma por lote: documentToSign devuelve
// un 7z con todos los documentos y uploadDocumentSigned recibe el 7z firmado).
type SevenZipArchive interface {
	// Build7z crea un archivo 7z con los documentos (nombre original + bytes).
	Build7z(ctx context.Context, files []domain.FirmaPeruDocument) ([]byte, error)
	// Extract7z extrae el contenido de un archivo 7z.
	Extract7z(ctx context.Context, archive []byte) ([]domain.FirmaPeruDocument, error)
}
