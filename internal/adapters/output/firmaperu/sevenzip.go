package firmaperu

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// sevenZipTool crea y extrae archivos 7z delegando en el ejecutable 7-Zip
// (SEVENZIP_PATH). El Firmador Web de Firma Perú exige 7z en la firma por
// lote: documentToSign devuelve un 7z con todos los documentos y
// uploadDocumentSigned recibe el 7z firmado.
type sevenZipTool struct {
	exePath string
}

// NewSevenZipTool construye el adaptador de archivos 7z resolviendo la ruta
// válida del ejecutable de 7-Zip.
func NewSevenZipTool(exePath string) *sevenZipTool {
	return &sevenZipTool{exePath: resolveSevenZipExe(exePath)}
}

func resolveSevenZipExe(exePath string) string {
	if exePath != "" {
		if _, err := os.Stat(exePath); err == nil {
			return exePath
		}
	}
	if p, err := exec.LookPath("7z"); err == nil {
		return p
	}
	if p, err := exec.LookPath("7za"); err == nil {
		return p
	}
	fallbacks := []string{
		`C:\Program Files\7-Zip\7z.exe`,
		`C:\Program Files (x86)\7-Zip\7z.exe`,
		`C:\7-Zip\7z.exe`,
		`C:\tools\7z.exe`,
	}
	for _, f := range fallbacks {
		if _, err := os.Stat(f); err == nil {
			return f
		}
	}
	if exePath != "" {
		return exePath
	}
	return `C:\Program Files\7-Zip\7z.exe`
}

// Build7z crea un archivo 7z con los documentos (nombre original como nombre
// de entrada) y devuelve sus bytes.
func (t *sevenZipTool) Build7z(ctx context.Context, files []domain.FirmaPeruDocument) ([]byte, error) {
	if len(files) == 0 {
		return nil, domain.ErrFirmaPeruDocumentRequired
	}

	tmpDir, err := os.MkdirTemp("", "galenos_firma7z_*")
	if err != nil {
		return nil, fmt.Errorf("7z tmp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	args := []string{"a", "-t7z", "-mx=3", "-y", "lote.7z"}
	for i, f := range files {
		name := nombreSeguro(f.Name, i)
		if err := os.WriteFile(filepath.Join(tmpDir, name), f.Content, 0o644); err != nil {
			return nil, fmt.Errorf("7z writing %s: %w", name, err)
		}
		args = append(args, name)
	}

	cmd := exec.CommandContext(ctx, t.exePath, args...)
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("7z create: %w: %s", err, truncarSalida(string(out)))
	}

	return os.ReadFile(filepath.Join(tmpDir, "lote.7z"))
}

// Extract7z extrae todos los archivos de un 7z (recorre subdirectorios) y los
// devuelve con su nombre de entrada.
func (t *sevenZipTool) Extract7z(ctx context.Context, archive []byte) ([]domain.FirmaPeruDocument, error) {
	tmpDir, err := os.MkdirTemp("", "galenos_extract7z_*")
	if err != nil {
		return nil, fmt.Errorf("7z tmp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	arcPath := filepath.Join(tmpDir, "lote.7z")
	if err := os.WriteFile(arcPath, archive, 0o644); err != nil {
		return nil, fmt.Errorf("7z writing archive: %w", err)
	}
	outDir := filepath.Join(tmpDir, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("7z out dir: %w", err)
	}

	cmd := exec.CommandContext(ctx, t.exePath, "x", "-t7z", "-y", "-o"+outDir, arcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("7z extract: %w: %s", err, truncarSalida(string(out)))
	}

	var files []domain.FirmaPeruDocument
	err = filepath.WalkDir(outDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(outDir, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files = append(files, domain.FirmaPeruDocument{Name: rel, Content: content})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("7z listing output: %w", err)
	}
	if len(files) == 0 {
		return nil, domain.ErrFirmaPeruDocumentRequired
	}
	return files, nil
}

// nombreSeguro normaliza el nombre de un archivo para usarlo como nombre de
// entrada del 7z: sin separadores de ruta y con prefijo de índice para evitar
// colisiones y conservar el orden.
func nombreSeguro(name string, idx int) string {
	base := filepath.Base(strings.ReplaceAll(name, "\\", "/"))
	base = strings.TrimSpace(base)
	if base == "." || base == "" || base == "/" {
		base = fmt.Sprintf("doc%d.pdf", idx)
	}
	return fmt.Sprintf("%03d_%s", idx, base)
}

func truncarSalida(s string) string {
	if len(s) > 400 {
		return s[:400] + "..."
	}
	return s
}
