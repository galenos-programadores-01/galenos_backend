package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sisRepository struct {
	db *sql.DB
}

// NewSisRepository construye el adaptador que implementa el puerto de
// salida output.SisRepository contra el SP usp_go_webSisFiliacionesGestionar.
func NewSisRepository(db *sql.DB) output.SisRepository {
	return &sisRepository{db: db}
}

// GestionarAfiliacion invoca el procedimiento almacenado
// usp_go_webSisFiliacionesGestionar con todos los parámetros de la afiliación SIS.
// Los parámetros van nombrados, por lo que el driver mssql los despacha
// como llamada RPC, sin concatenar SQL. Los campos sin valor se envían
// como NULL.
func (r *sisRepository) GestionarAfiliacion(ctx context.Context, afiliacion *domain.SisAfiliacion) error {
	const procedure = `usp_go_webSisFiliacionesGestionar`

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("idSiasis", afiliacion.IDSiasis),
		sql.Named("Codigo", afiliacion.Codigo),
		sql.Named("AfiliacionDisa", afiliacion.AfiliacionDisa),
		sql.Named("AfiliacionTipoFormato", afiliacion.TipoFormato),
		sql.Named("AfiliacionNroFormato", afiliacion.NroFormato),
		sql.Named("AfiliacionNroIntegrante", afiliacion.NroIntegrante),
		sql.Named("DocumentoTipo", afiliacion.DocumentoTipo),
		sql.Named("CodigoEstablAdscripcion", afiliacion.CodigoEstablAdscripcion),
		sql.Named("AfiliacionFecha", afiliacion.AfiliacionFecha),
		sql.Named("Paterno", afiliacion.Paterno),
		sql.Named("Materno", afiliacion.Materno),
		sql.Named("Pnombre", afiliacion.PNombre),
		sql.Named("Onombres", afiliacion.ONombres),
		sql.Named("Genero", afiliacion.Genero),
		sql.Named("Fnacimiento", afiliacion.FNacimiento),
		sql.Named("IdDistritoDomicilio", afiliacion.IdDistritoDomicilio),
		sql.Named("Estado", afiliacion.Estado),
		sql.Named("Fbaja", afiliacion.Fbaja),
		sql.Named("DocumentoNumero", afiliacion.DocumentoNumero),
		sql.Named("MotivoBaja", afiliacion.MotivoBaja),
		sql.Named("FbajaOK", afiliacion.FbajaOK),
		sql.Named("DescEESS", afiliacion.DescEESS),
		sql.Named("DescEESSUbigeo", afiliacion.DescEESSUbigeo),
		sql.Named("Regimen", afiliacion.Regimen),
		sql.Named("TipoSeguro", afiliacion.TipoSeguro),
		sql.Named("DescTipoSeguro", afiliacion.DescTipoSeguro),
		sql.Named("Contrato", afiliacion.Contrato),
		sql.Named("IdPlan", afiliacion.IdPlan),
		sql.Named("IdGrupoPoblacional", afiliacion.IdGrupoPoblacional),
		sql.Named("MsgConfidencial", afiliacion.MsgConfidencial),
		sql.Named("IdUsuarioAuditoria", afiliacion.IdUsuarioAuditoria),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_webSisFiliacionesGestionar: %w", err)
	}

	return nil
}

// ForzarGuardadoFua invoca el procedimiento almacenado
// webSisFuaAtencionForzarGuardado con el id de la cuenta de atención,
// forzando la generación/guardado del FUA en el SIS.
func (r *sisRepository) ForzarGuardadoFua(ctx context.Context, idCuentaAtencion int64) error {
	const procedure = `webSisFuaAtencionForzarGuardado`

	_, err := r.db.ExecContext(ctx, procedure, sql.Named("IdcuentaAtencion", idCuentaAtencion))
	if err != nil {
		return fmt.Errorf("calling webSisFuaAtencionForzarGuardado: %w", err)
	}

	return nil
}

// AgregarFua invoca el procedimiento almacenado usp_go_webFUAgregar con el id de
// la cuenta de atención, el empleado y el nombre del PC. El parámetro
// @Respuesta se declara con sql.Out para capturar el valor de salida.
func (r *sisRepository) AgregarFua(ctx context.Context, idCuentaAtencion, idEmpleado int64, nombrePc string) (string, error) {
	const procedure = `usp_go_webFUAgregar`

	var respuesta string

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("Respuesta", sql.Out{Dest: &respuesta}),
		sql.Named("IdcuentaAtencion", idCuentaAtencion),
		sql.Named("IdEmplado", idEmpleado),
		sql.Named("NombrePc", nombrePc),
	)
	if err != nil {
		return "", fmt.Errorf("calling usp_go_webFUAgregar: %w", err)
	}

	return respuesta, nil
}

// GetFuaImprimir invoca el procedimiento almacenado
// webFuaImprimirIdCuentaAtencion con el id de la cuenta de atención. El SP
// retorna los datos del FUA para imprimir. Devuelve nil si no hay registros.
func (r *sisRepository) GetFuaImprimir(ctx context.Context, idCuentaAtencion int64) (*map[string]any, error) {
	const procedure = `webFuaImprimirIdCuentaAtencion`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdcuentaAtencion", idCuentaAtencion))
	if err != nil {
		return nil, fmt.Errorf("calling webFuaImprimirIdCuentaAtencion: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading fua print data: %w", err)
	}
	if len(maps) == 0 {
		return nil, nil
	}

	m := maps[0]
	return &m, nil
}

// listRows invoca un procedimiento almacenado de solo lectura que retorna
// varias filas y las transforma a []map[string]any con la utilidad rowsToMaps.
// Recibe el nombre del SP y los parámetros nombrados a enviar.
func (r *sisRepository) listRows(ctx context.Context, procedure string, args ...any) ([]map[string]any, error) {
	rows, err := r.db.QueryContext(ctx, procedure, args...)
	if err != nil {
		return nil, fmt.Errorf("calling %s: %w", procedure, err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading rows of %s: %w", procedure, err)
	}

	return maps, nil
}

// ListDiagnosticos invoca webAtencionesDiagnosticosIdCuentaAtencion con el id de
// la cuenta de atención y devuelve los diagnósticos de esa atención.
func (r *sisRepository) ListDiagnosticos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error) {
	return r.listRows(ctx, `webAtencionesDiagnosticosIdCuentaAtencion`, sql.Named("Parametros", fmt.Sprintf("%d", idCuentaAtencion)))
}

// ListMedicamentos invoca webMedicamentosListarIdCuentaAtencion con el id
// de la cuenta de atención y el parámetro RecetaAdicional en -100
// (por defecto) y devuelve los medicamentos de esa atención.
func (r *sisRepository) ListMedicamentos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error) {
	return r.listRows(ctx, `webMedicamentosListarIdCuentaAtencion`,
		sql.Named("IdcuentaAtencion", idCuentaAtencion),
		sql.Named("RecetaAdicional", -100),
	)
}

// ListProcedimientos invoca webProcedimientoListarIdCuentaAtencion con el
// id de la cuenta de atención y devuelve los procedimientos de esa atención.
func (r *sisRepository) ListProcedimientos(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error) {
	return r.listRows(ctx, `webProcedimientoListarIdCuentaAtencion`, sql.Named("IdcuentaAtencion", idCuentaAtencion))
}

// ListConsumo invoca webFactOrdenServicioDetaDesFinaListarIdcuenta con el id
// de la cuenta de atención y el IdUsuario 2937 (por defecto) y devuelve el
// detalle de consumo de esa atención.
func (r *sisRepository) ListConsumo(ctx context.Context, idCuentaAtencion int64) ([]map[string]any, error) {
	return r.listRows(ctx, `webFactOrdenServicioDetaDesFinaListarIdcuenta`,
		sql.Named("IdCuentaAtencion", idCuentaAtencion),
		sql.Named("IdUsuario", 2937),
	)
}
