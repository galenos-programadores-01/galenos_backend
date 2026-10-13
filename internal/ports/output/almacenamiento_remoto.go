package output

import "context"

// AlmacenamientoRemoto puerto de salida para subir/bajar archivos
// (implementaciones: SFTP o FTP). Los casos de uso dependen solo de esta interfaz.
type AlmacenamientoRemoto interface {
	// Enabled indica si hay configuración usable (host/credenciales).
	Enabled() bool
	// Subir crea directorios padres si hace falta y escribe el archivo.
	Subir(ctx context.Context, rutaRemota string, contenido []byte) error
	// Descargar lee el archivo remoto completo.
	Descargar(ctx context.Context, rutaRemota string) ([]byte, error)
	// Existe consulta si la ruta remota está presente.
	Existe(ctx context.Context, rutaRemota string) (bool, error)
	// Eliminar borra el archivo remoto (no elimina directorios).
	Eliminar(ctx context.Context, rutaRemota string) error
}
