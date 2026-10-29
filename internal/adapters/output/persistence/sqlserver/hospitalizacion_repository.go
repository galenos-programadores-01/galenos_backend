package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type hospitalizacionRepository struct {
	db *sql.DB
}

// NewHospitalizacionRepository construye el adaptador que implementa el
// puerto de salida output.HospitalizacionRepository contra el SP
// usp_go_AdmisionHospitalizacion.
func NewHospitalizacionRepository(db *sql.DB) output.HospitalizacionRepository {
	return &hospitalizacionRepository{db: db}
}

// Admisionar invoca el procedimiento almacenado usp_go_AdmisionHospitalizacion
// con todos los parámetros nombrados, enviados como llamada RPC por el driver
// mssql, por lo que no se concatena SQL. Los parámetros sin valor se envían
// como NULL. El parámetro @Resultado se declara con sql.Out para capturar el
// valor de salida que el SP devuelve.
//
// Los nombres de parámetro conservan exactamente los declarados en el SP,
// incluyendo las variantes con mayúsculas que la base refleja
// (madretipodoc, madrenrodoc, idProvincia, Idnacionalidad, IdfuenteFinanciamiento,
// IdEstadollego, etc.), ya que el driver los resuelve contra el nombre formal
// del parámetro.
func (r *hospitalizacionRepository) Admisionar(ctx context.Context, admision *domain.AdmisionHospitalizacion) (string, error) {
	const procedure = `usp_go_AdmisionHospitalizacion`

	var resultado string

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdPaciente", admision.IDPaciente),
		sql.Named("IdDocIdentidad", admision.IDDocIdentidad),
		sql.Named("NroDocumento", admision.NroDocumento),
		sql.Named("ApellidoPaterno", admision.ApellidoPaterno),
		sql.Named("ApellidoMaterno", admision.ApellidoMaterno),
		sql.Named("PrimerNombre", admision.PrimerNombre),
		sql.Named("SegundoNombre", admision.SegundoNombre),
		sql.Named("TercerNombre", admision.TercerNombre),
		sql.Named("FechaNacimiento", admision.FechaNacimiento),
		sql.Named("IdTipoSexo", admision.IDTipoSexo),
		sql.Named("EstadoCivil", admision.EstadoCivil),
		sql.Named("GradoInstruccion", admision.GradoInstruccion),
		sql.Named("Ocupacion", admision.Ocupacion),
		sql.Named("IdEtnia", admision.IDEtnia),
		sql.Named("IdIdioma", admision.IDIdioma),
		sql.Named("madretipodoc", admision.MadreTipoDoc),
		sql.Named("madrenrodoc", admision.MadreNroDoc),
		sql.Named("madreapepaterno", admision.MadreApePaterno),
		sql.Named("madreapematerno", admision.MadreApeMaterno),
		sql.Named("madreprimernombre", admision.MadrePrimerNombre),
		sql.Named("madresegundonombre", admision.MadreSegundoNombre),
		sql.Named("NombrePadre", admision.NombrePadre),
		sql.Named("Telefono", admision.Telefono),
		sql.Named("Correo", admision.Correo),
		sql.Named("IdDepartamento", admision.IDDepartamento),
		sql.Named("idProvincia", admision.IDProvincia),
		sql.Named("IdDistrito", admision.IDDistrito),
		sql.Named("IdComunidadDomicilio", admision.IDComunidadDomicilio),
		sql.Named("Direccion", admision.Direccion),
		sql.Named("Idnacionalidad", admision.IDNacionalidad),
		sql.Named("IdOrigen", admision.IDOrigen),
		sql.Named("IdServicio", admision.IDServicio),
		sql.Named("IdCama", admision.IDCama),
		sql.Named("FechaIngreso", admision.FechaIngreso),
		sql.Named("HoraIngreso", admision.HoraIngreso),
		sql.Named("IdMedicoIngreso", admision.IDMedicoIngreso),
		sql.Named("IdTipoPrioridad", admision.IDTipoPrioridad),
		sql.Named("IdEstadollego", admision.IDEstadollego),
		sql.Named("IdfuenteFinanciamiento", admision.IDFuenteFinanciamiento),
		sql.Named("IdEstablecimientoReferenciaOrigen", admision.IDEstablecimientoReferenciaOrigen),
		sql.Named("NroReferencia", admision.NroReferencia),
		sql.Named("IdEmpleado", admision.IDEmpleado),
		sql.Named("Resultado", sql.Out{Dest: &resultado}),
	)
	if err != nil {
		return "", fmt.Errorf("calling usp_go_AdmisionHospitalizacion: %w", err)
	}

	return resultado, nil
}

// ObtenerReporteTicket invoca el procedimiento almacenado
// usp_go_GenererReporteTicketHospi con el parámetro @IdCuenta. Las columnas se
// leen por nombre para no depender del orden del SELECT. Retorna nil si el SP
// no devuelve filas.
func (r *hospitalizacionRepository) ObtenerReporteTicket(ctx context.Context, idCuenta int64) (*domain.ReporteTicketHospitalizacion, error) {
	const procedure = `usp_go_GenererReporteTicketHospi`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdCuenta", idCuenta))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_GenererReporteTicketHospi: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading usp_go_GenererReporteTicketHospi rows: %w", err)
	}
	if len(maps) == 0 {
		return nil, nil
	}

	row := maps[0]

	reporte := &domain.ReporteTicketHospitalizacion{
		Paciente:           derefString(rowString(row, "Paciente")),
		NroHistoriaClinica: derefString(rowString(row, "NroHistoriaClinica")),
		NroDocumento:       derefString(rowString(row, "NroDocumento")),
		Servicio:           derefString(rowString(row, "Servicio")),
		Medico:             derefString(rowString(row, "Medico")),
		Fecha:              rowTime(row, "Fecha"),
		IAFA:               derefString(rowString(row, "IAFA")),
	}
	if v := rowInt64(row, "IdCuentaAtencion"); v != nil {
		reporte.IDCuentaAtencion = *v
	}
	if hora := rowTime(row, "Hora"); hora != nil {
		reporte.Hora = hora.Format("15:04")
	} else {
		reporte.Hora = derefString(rowString(row, "Hora"))
	}

	return reporte, nil
}
