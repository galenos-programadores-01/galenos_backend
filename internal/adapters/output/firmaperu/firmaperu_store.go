package firmaperu

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type memoryStore struct {
	mu        sync.Mutex
	records   map[string]*registro
	signedDir string
}

type registro struct {
	original     domain.FirmaPeruDocument
	signed       []byte
	solicitud    *domain.FirmaPeruSolicitud
	stamp        []byte
	lote         []domain.FirmaPeruDocument
	loteSign     []byte
	loteFirmados []domain.FirmaPeruDocument
}

// NewMemoryStore crea el almacén de documentos del proceso de firma. Si
// signedDir no es vacío, los documentos firmados recibidos también se escriben
// en ese directorio (p. ej. la carpeta public del frontend).
func NewMemoryStore(signedDir string) *memoryStore {
	return &memoryStore{records: make(map[string]*registro), signedDir: signedDir}
}

func (s *memoryStore) SaveDocument(ctx context.Context, uuid string, doc domain.FirmaPeruDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil {
		r = &registro{}
		s.records[uuid] = r
	}
	r.original = doc
	r.signed = nil
	return nil
}

func (s *memoryStore) GetDocument(ctx context.Context, uuid string) (domain.FirmaPeruDocument, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil || len(r.original.Content) == 0 {
		return domain.FirmaPeruDocument{}, domain.ErrFirmaPeruDocumentNotFound
	}
	return r.original, nil
}

func (s *memoryStore) SaveSigned(ctx context.Context, uuid string, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("FirmaPeru Store: SaveSigned uuid=%s bytes=%d", uuid, len(content))

	r := s.records[uuid]
	if r == nil {
		return domain.ErrFirmaPeruDocumentNotFound
	}
	if len(content) == 0 {
		return domain.ErrFirmaPeruDocumentRequired
	}

	r.signed = append([]byte(nil), content...)

	return s.guardarEnDisco(r, uuid, r.original.Name, content)
}

// guardarEnDisco escribe el documento firmado en FIRMAPERU_SIGNED_DIR (si
// está configurado) con el nombre indicado.
func (s *memoryStore) guardarEnDisco(r *registro, uuid string, nombreOriginal string, content []byte) error {
	if s.signedDir == "" {
		return nil
	}
	if err := os.MkdirAll(s.signedDir, 0o755); err != nil {
		log.Printf("FirmaPeru Store: error creando dir %s: %v", s.signedDir, err)
		return err
	}
	nombre := filepath.Base(nombreOriginal)
	if nombre == "." || nombre == "" {
		nombre = uuid + ".pdf"
	}
	ruta := filepath.Join(s.signedDir, nombre)
	if err := os.WriteFile(ruta, content, 0o644); err != nil {
		log.Printf("FirmaPeru Store: error escribiendo %s: %v", ruta, err)
		return err
	}
	log.Printf("FirmaPeru Store: firmado guardado en %s (%d bytes)", ruta, len(content))
	return nil
}

func (s *memoryStore) GetSigned(ctx context.Context, uuid string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil || len(r.signed) == 0 {
		return nil, domain.ErrFirmaPeruSignedNotAvailable
	}
	return r.signed, nil
}

func (s *memoryStore) SaveSolicitud(ctx context.Context, uuid string, sol domain.FirmaPeruSolicitud) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil {
		return domain.ErrFirmaPeruDocumentNotFound
	}
	solicitud := sol
	r.solicitud = &solicitud
	return nil
}

func (s *memoryStore) GetSolicitud(ctx context.Context, uuid string) (domain.FirmaPeruSolicitud, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil || r.solicitud == nil {
		return domain.FirmaPeruSolicitud{}, domain.ErrFirmaPeruDocumentNotFound
	}
	return *r.solicitud, nil
}

func (s *memoryStore) SaveStampImage(ctx context.Context, uuid string, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil {
		return domain.ErrFirmaPeruDocumentNotFound
	}
	r.stamp = append([]byte(nil), content...)
	return nil
}

func (s *memoryStore) GetStampImage(ctx context.Context, uuid string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil || len(r.stamp) == 0 {
		return nil, domain.ErrFirmaPeruDocumentNotFound
	}
	return r.stamp, nil
}

func (s *memoryStore) SaveLote(ctx context.Context, uuid string, docs []domain.FirmaPeruDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil {
		r = &registro{}
		s.records[uuid] = r
	}
	r.lote = append([]domain.FirmaPeruDocument(nil), docs...)
	r.loteSign = nil
	r.loteFirmados = nil
	return nil
}

func (s *memoryStore) GetLote(ctx context.Context, uuid string) ([]domain.FirmaPeruDocument, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil || len(r.lote) == 0 {
		return nil, domain.ErrFirmaPeruDocumentNotFound
	}
	return r.lote, nil
}

func (s *memoryStore) SaveLoteSigned(ctx context.Context, uuid string, archive []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil {
		return domain.ErrFirmaPeruDocumentNotFound
	}
	if len(archive) == 0 {
		return domain.ErrFirmaPeruDocumentRequired
	}
	r.loteSign = append([]byte(nil), archive...)
	return nil
}

func (s *memoryStore) GetLoteSigned(ctx context.Context, uuid string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil || len(r.loteSign) == 0 {
		return nil, domain.ErrFirmaPeruSignedNotAvailable
	}
	return r.loteSign, nil
}

// SaveLoteSignedFile guarda en memoria (y disco si FIRMAPERU_SIGNED_DIR) uno
// de los PDFs firmados del lote, con el nombre original del documento.
func (s *memoryStore) SaveLoteSignedFile(ctx context.Context, uuid string, name string, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(content) == 0 {
		return domain.ErrFirmaPeruDocumentRequired
	}
	r := s.records[uuid]
	if r == nil {
		return domain.ErrFirmaPeruDocumentNotFound
	}

	nombre := filepath.Base(name)
	if nombre == "." || nombre == "" {
		nombre = uuid + ".pdf"
	}
	copia := append([]byte(nil), content...)
	actualizado := false
	for i := range r.loteFirmados {
		if strings.EqualFold(filepath.Base(r.loteFirmados[i].Name), nombre) {
			r.loteFirmados[i] = domain.FirmaPeruDocument{Name: nombre, Content: copia}
			actualizado = true
			break
		}
	}
	if !actualizado {
		r.loteFirmados = append(r.loteFirmados, domain.FirmaPeruDocument{
			Name:    nombre,
			Content: copia,
		})
	}
	return s.guardarEnDisco(r, uuid, nombre, content)
}

func (s *memoryStore) GetLoteSignedFiles(ctx context.Context, uuid string) ([]domain.FirmaPeruDocument, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil || len(r.loteFirmados) == 0 {
		return nil, domain.ErrFirmaPeruSignedNotAvailable
	}
	salida := make([]domain.FirmaPeruDocument, len(r.loteFirmados))
	for i, f := range r.loteFirmados {
		salida[i] = domain.FirmaPeruDocument{
			Name:    f.Name,
			Content: append([]byte(nil), f.Content...),
		}
	}
	return salida, nil
}

func (s *memoryStore) GetLoteSignedFile(ctx context.Context, uuid string, name string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.records[uuid]
	if r == nil || len(r.loteFirmados) == 0 {
		return nil, domain.ErrFirmaPeruSignedNotAvailable
	}
	nombreBuscado := strings.ToLower(filepath.Base(name))
	for _, doc := range r.loteFirmados {
		if strings.ToLower(filepath.Base(doc.Name)) == nombreBuscado {
			return append([]byte(nil), doc.Content...), nil
		}
	}
	return nil, domain.ErrFirmaPeruDocumentNotFound
}
