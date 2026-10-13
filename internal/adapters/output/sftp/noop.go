package sftp

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// noopClient implementa AlmacenamientoRemoto cuando SFTP no está configurado.
// Enabled()=false; las operaciones devuelven ErrAlmacenamientoNoConfigurado.
type noopClient struct{}

func (noopClient) Enabled() bool { return false }

func (noopClient) Subir(context.Context, string, []byte) error {
	return domain.ErrAlmacenamientoNoConfigurado
}

func (noopClient) Descargar(context.Context, string) ([]byte, error) {
	return nil, domain.ErrAlmacenamientoNoConfigurado
}

func (noopClient) Existe(context.Context, string) (bool, error) {
	return false, domain.ErrAlmacenamientoNoConfigurado
}

func (noopClient) Eliminar(context.Context, string) error {
	return domain.ErrAlmacenamientoNoConfigurado
}
