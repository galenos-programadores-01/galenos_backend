package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// SisClient es el puerto de salida para consultar el servicio web SOAP del
// SIS. La implementación concreta invoca el endpoint externo y maneja la
// sesión (GetSession + ConsultarAfiliadoFuaE / BuscarAsegurados).
type SisClient interface {
	// ConsultarAfiliado trae el paciente afiliado por su número de documento
	// contra el servicio SIS.
	ConsultarAfiliado(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error)

	// BuscarPorAfiliacion trae el paciente afiliado por los parámetros de
	// afiliación (DISA, Lote, Contrato, Correlativo, CodTabla) usando la
	// operación BuscarAsegurados del SIS.
	BuscarPorAfiliacion(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error)
}
