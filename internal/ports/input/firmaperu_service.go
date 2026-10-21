package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// FirmaPeruService es el puerto de entrada del proceso de firma digital con
// el Firmador Web de la Plataforma Firma Perú. El flujo es asíncrono e
// interactivo: IniciarFirma prepara el documento, almacena la solicitud y
// devuelve el wrapper base64 (param_url/param_token) para startSignature;
// el Firmador consulta ObtenerParametrosFirma (endpoint param_url), descarga
// el documento, el usuario firma con su DNIe y lo devuelve, y el backend lo
// expone con ObtenerDocumentoFirmado.
type FirmaPeruService interface {
	// IniciarFirma valida el documento y los parámetros, almacena el original
	// y la solicitud, y devuelve el wrapper base64 (startSignature) y el UUID
	// del proceso.
	IniciarFirma(ctx context.Context, req domain.FirmaPeruInitRequest) (domain.FirmaPeruInitResult, error)
	// ObtenerParametrosFirma valida el param_token de un único uso y devuelve
	// el JSON base64 con los parámetros de firma (documentToSign,
	// uploadDocumentSigned y token de la plataforma) que el Firmador local
	// consume. Es el endpoint expuesto en param_url.
	ObtenerParametrosFirma(ctx context.Context, uuid string, paramToken string) (string, error)
	// ObtenerDocumento sirve el documento original por firmar (documentToSign).
	ObtenerDocumento(ctx context.Context, uuid string) (domain.FirmaPeruDocument, error)
	// RecibirDocumentoFirmado almacena el documento firmado devuelto por el
	// Firmador al callback uploadDocumentSigned.
	RecibirDocumentoFirmado(ctx context.Context, uuid string, content []byte) error
	// ObtenerDocumentoFirmado devuelve el documento firmado cuando está listo.
	ObtenerDocumentoFirmado(ctx context.Context, uuid string) ([]byte, error)
	// ObtenerImagenEstampado devuelve la imagen PNG del estampado del proceso
	// (si se definió una propia) o un error de dominio si no existe.
	ObtenerImagenEstampado(ctx context.Context, uuid string) ([]byte, error)
	// IniciarFirmaLote inicia una firma por lote (varios documentos con un
	// solo PIN): almacena todos los documentos y devuelve el wrapper base64.
	// El Firmador descargará un 7z con los documentos (documentToSign) y
	// devolverá el 7z firmado (uploadDocumentSigned).
	IniciarFirmaLote(ctx context.Context, req domain.FirmaPeruLoteInitRequest) (domain.FirmaPeruInitResult, error)
	// ObtenerDocumentosLote sirve el archivo 7z con todos los documentos del
	// lote por firmar (documentToSign en modo lote).
	ObtenerDocumentosLote(ctx context.Context, uuid string) ([]byte, error)
	// RecibirDocumentoFirmadoLote recibe el 7z firmado del lote, lo extrae y
	// guarda cada PDF firmado en FIRMAPERU_SIGNED_DIR con su nombre original.
	RecibirDocumentoFirmadoLote(ctx context.Context, uuid string, archive []byte) error
	// ObtenerLoteFirmado devuelve el 7z firmado del lote cuando está listo.
	ObtenerLoteFirmado(ctx context.Context, uuid string) ([]byte, error)
	ObtenerLoteDocumentosFirmados(ctx context.Context, uuid string) ([]domain.FirmaPeruDocument, error)
	ObtenerLoteArchivoFirmado(ctx context.Context, uuid string, name string) ([]byte, error)
}
