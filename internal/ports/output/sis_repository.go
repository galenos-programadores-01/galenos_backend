package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// SisRepository es el puerto de salida para persistir las afiliaciones SIS
// de los pacientes asegurados contra la base de datos.
type SisRepository interface {
	// GestionarAfiliacion guarda o actualiza una afiliación SIS invocando
	// el procedimiento almacenado usp_go_webSisFiliacionesGestionar.
	GestionarAfiliacion(ctx context.Context, afiliacion *domain.SisAfiliacion) error

	// ForzarGuardadoFua fuerza el guardado del FUA de una cuenta de
	// atención invocando el SP webSisFuaAtencionForzarGuardado.
	ForzarGuardadoFua(ctx context.Context, idCuentaAtencion int64) error

	// AgregarFua agrega el FUA de una cuenta de atención invocando el SP
	// usp_go_webFUAgregar. Retorna el valor del parámetro de salida @Respuesta.
	AgregarFua(ctx context.Context, idCuentaAtencion, idEmpleado int64, nombrePc string) (string, error)

	// GetFuaImprimir consulta los datos del FUA para imprimir invocando el
	// SP webFuaImprimirIdCuentaAtencion.
	GetFuaImprimir(ctx context.Context, idCuentaAtencion int64) (*map[string]any, error)

	// ListDiagnosticos consulta los diagnósticos de una atención invocando el
	// SP webAtencionesDiagnosticosIdCuentaAtencion.
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
