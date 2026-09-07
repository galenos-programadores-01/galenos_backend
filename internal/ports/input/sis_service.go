package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// SisService es el puerto de entrada para consultar el servicio SIS.
type SisService interface {
	// ConsultarAfiliado trae el paciente afiliado por su número de documento.
	ConsultarAfiliado(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error)

	// GestionarAfiliacion guarda o actualiza una afiliación SIS invocando
	// el SP webSisFiliacionesGestionar.
	GestionarAfiliacion(ctx context.Context, afiliacion *domain.SisAfiliacion) error

	// ForzarGuardadoFua fuerza el guardado del FUA de una cuenta de
	// atención invocando el SP webSisFuaAtencionForzarGuardado.
	ForzarGuardadoFua(ctx context.Context, idCuentaAtencion int64) error

	// AgregarFua agrega el FUA de una cuenta de atención invocando el SP
	// usp_go_webFUAgregar. Retorna el @Respuesta del SP.
	AgregarFua(ctx context.Context, idCuentaAtencion, idEmpleado int64, nombrePc string) (string, error)

	// GetFuaImprimir consulta los datos del FUA para imprimir invocando el
	// SP webFuaImprimirIdCuentaAtencion. Retorna nil si el SP no devuelve filas.
	GetFuaImprimir(ctx context.Context, idCuentaAtencion int64) (*map[string]any, error)

	// ListDiagnosticos consulta los diagnósticos de una atención invocando el
	// SP webAtencionesDiagnosticosIdAtencion.
	ListDiagnosticos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error)

	// ListMedicamentos consulta los medicamentos de una cuenta de atención
	// invocando el SP webMedicamentosListarIdCuentaAtencion.
	ListMedicamentos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error)

	// ListProcedimientos consulta los procedimientos de una cuenta de atención
	// invocando el SP webProcedimientoListarIdCuentaAtencion.
	ListProcedimientos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error)

	// ListConsumo consulta el detalle de consumo (orden de servicio) de una
	// cuenta de atención invocando el SP webFactOrdenServicioDetaDesFinaListarIdcuenta.
	ListConsumo(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error)
}
