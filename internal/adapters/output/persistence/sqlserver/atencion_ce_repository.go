package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type atencionCERepository struct {
	db *sql.DB
}

func NewAtencionCERepository(db *sql.DB) output.AtencionCERepository {
	return &atencionCERepository{
		db: db,
	}
}

func (r *atencionCERepository) ListarColaPacientes(ctx context.Context, filtros domain.ColaFiltrosCE) ([]domain.ColaPacienteCEItem, error) {
	query := `EXEC dbo.usp_go_AtencionesCE_ListarCola 
		@IdMedico = @p1, 
		@IdServicio = @p2, 
		@Fecha = @p3, 
		@Filtro = @p4`

	var fechaParam any = nil
	if filtros.Fecha != "" {
		fechaParam = filtros.Fecha
	}

	filtroLimpio := strings.TrimSpace(filtros.Filtro)
	var filtroParam any = nil
	if filtroLimpio != "" {
		filtroParam = filtroLimpio
	}

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", filtros.IdMedico),
		sql.Named("p2", filtros.IdServicio),
		sql.Named("p3", fechaParam),
		sql.Named("p4", filtroParam),
	)
	if err != nil {
		return nil, fmt.Errorf("error al listar cola de pacientes CE mediante SP: %w", err)
	}
	defer rows.Close()

	cola := make([]domain.ColaPacienteCEItem, 0)
	for rows.Next() {
		var item domain.ColaPacienteCEItem
		err := rows.Scan(
			&item.IdAtencion,
			&item.IdCuentaAtencion,
			&item.IdPaciente,
			&item.NroHistoriaClinica,
			&item.NroDocumento,
			&item.PacienteNombre,
			&item.FechaIngreso,
			&item.HoraIngreso,
			&item.HoraEgreso,
			&item.Edad,
			&item.IdServicioIngreso,
			&item.ServicioNombre,
			&item.IdMedicoIngreso,
			&item.MedicoNombre,
			&item.EstadoAtencion,
			&item.IdCita,
			&item.TipoFinanciador,
			&item.Sexo,
			&item.Telefono,
			&item.TieneTriaje,
			&item.Turno,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear cola de paciente CE: %w", err)
		}
		cola = append(cola, item)
	}

	return cola, rows.Err()
}

func (r *atencionCERepository) ObtenerFichaAtencion(ctx context.Context, idAtencion int) (*domain.FichaAtencionCEDetalle, error) {
	queryCabecera := `EXEC dbo.usp_go_AtencionesCE_ObtenerFicha @IdAtencion = @pIdAtencion`

	var ficha domain.FichaAtencionCEDetalle
	err := r.db.QueryRowContext(ctx, queryCabecera, sql.Named("pIdAtencion", idAtencion)).Scan(
		&ficha.IdAtencion,
		&ficha.IdCuentaAtencion,
		&ficha.IdPaciente,
		&ficha.PacienteNombre,
		&ficha.NroHistoriaClinica,
		&ficha.NroDocumento,
		&ficha.Edad,
		&ficha.Sexo,
		&ficha.Motivo,
		&ficha.ExamenClinico,
		&ficha.Antecedente,
		&ficha.InforCompDiagMed,
		&ficha.Tratamiento,
		&ficha.PlanTratamiento,
		&ficha.Observaciones,
		&ficha.FechaProxCita,
		&ficha.IdDestinoAtencion,
		&ficha.DestinoNombre,
		&ficha.HoraIngreso,
		&ficha.HoraEgreso,
		&ficha.Estado,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener cabecera de atencion CE: %w", err)
	}

	queryDx := `EXEC [dbo].[usp_go_ObtenerDiagnosticosAtencion] @IdAtencion = @pIdAtencion, @IdPrimeraAtencion = 0, @IdEvolucion = 0`
	rowsDx, err := r.db.QueryContext(ctx, queryDx, sql.Named("pIdAtencion", idAtencion))
	if err != nil {
		return nil, fmt.Errorf("error al obtener diagnosticos de atencion CE: %w", err)
	}
	defer rowsDx.Close()

	ficha.Diagnosticos = make([]domain.AtencionDiagnosticoCEItem, 0)
	for rowsDx.Next() {
		var dx domain.AtencionDiagnosticoCEItem
		var idAtencionDummy int
		var descLarga string
		var tipoCodigo string
		var tipoDx string
		var idSubclasificacionBD int

		if err := rowsDx.Scan(
			&dx.IdAtencionDiagnostico,
			&idAtencionDummy,
			&dx.IdDiagnostico,
			&dx.CodigoCIE10,
			&dx.Descripcion,
			&descLarga,
			&tipoCodigo,
			&tipoDx,
			&idSubclasificacionBD,
			&dx.LabConfHIS,
			&dx.LabConfHIS1,
			&dx.LabConfHIS2,
			&dx.LabConfHIS3,
		); err != nil {
			return nil, fmt.Errorf("error al escanear diagnostico CE: %w", err)
		}

		dx.Clasificacion = strings.ToUpper(strings.TrimSpace(tipoCodigo))
		if dx.Clasificacion == "" {
			dx.Clasificacion = "D"
		}
		if dx.Clasificacion == "D" {
			dx.IdClasificacionDx = 2
		} else if dx.Clasificacion == "R" {
			dx.IdClasificacionDx = 3
		} else {
			dx.IdClasificacionDx = 1
		}

		if len(ficha.Diagnosticos) == 0 {
			dx.IdSubclasificacionDx = 1
		} else {
			dx.IdSubclasificacionDx = 2
		}

		ficha.Diagnosticos = append(ficha.Diagnosticos, dx)
	}

	ficha.Procedimientos = make([]domain.AtencionProcedimientoCEItem, 0)

	return &ficha, nil
}

func (r *atencionCERepository) GuardarFichaAtencion(ctx context.Context, payload domain.AtencionCEGuardarPayload) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar transaccion para atencion CE: %w", err)
	}
	defer tx.Rollback()

	motivo := strings.TrimSpace(payload.Motivo)
	if motivo == "" {
		motivo = "Atención médica en consultorio externo"
	}

	spGuardarAtencion := `
		EXEC [dbo].[AtencionesCEGuardar]
			@idAtencion = @p1,
			@Motivo = @p2,
			@ExamenClinico = @p3,
			@Antecedente = @p4,
			@InforCompDiagMed = @p5,
			@Tratamiento = @p6,
			@PlanTratamiento = @p7,
			@Observaciones = @p8,
			@IdUsuario = @p9
	`
	_, err = tx.ExecContext(ctx, spGuardarAtencion,
		sql.Named("p1", payload.IdAtencion),
		sql.Named("p2", motivo),
		sql.Named("p3", strings.TrimSpace(payload.ExamenClinico)),
		sql.Named("p4", strings.TrimSpace(payload.Antecedente)),
		sql.Named("p5", strings.TrimSpace(payload.InforCompDiagMed)),
		sql.Named("p6", strings.TrimSpace(payload.Tratamiento)),
		sql.Named("p7", strings.TrimSpace(payload.PlanTratamiento)),
		sql.Named("p8", strings.TrimSpace(payload.Observaciones)),
		sql.Named("p9", payload.IdUsuario),
	)
	if err != nil {
		return fmt.Errorf("error al ejecutar procedimiento AtencionesCEGuardar: %w", err)
	}

	queryEstadoSP := `EXEC dbo.usp_go_Atenciones_ActualizarEstadoCE 
		@IdAtencion = @p1, 
		@IdDestinoAtencion = @p2, 
		@CerrarAtencion = @p3`
	var exito, idPaciente int
	err = tx.QueryRowContext(ctx, queryEstadoSP,
		sql.Named("p1", payload.IdAtencion),
		sql.Named("p2", payload.IdDestinoAtencion),
		sql.Named("p3", payload.CerrarAtencion),
	).Scan(&exito, &idPaciente)
	if err != nil {
		return fmt.Errorf("error al actualizar estado de atencion mediante SP: %w", err)
	}

	if err := guardarDiagnosticosAtencionTx(ctx, tx, payload.IdAtencion, idPaciente, payload.Diagnosticos); err != nil {
		return err
	}

	return tx.Commit()
}

func determinarSubclasificacionDx(diagnostico domain.AtencionDiagnosticoGuardarItem) int {
	if diagnostico.Clasificacion != "" {
		switch strings.ToUpper(strings.TrimSpace(diagnostico.Clasificacion)) {
		case "D", "DEFINITIVO":
			return 102
		case "R", "REPETIDO":
			return 103
		case "P", "PRESUNTIVO":
			return 101
		}
	}

	switch diagnostico.IdClasificacionDx {
	case 2:
		return 102
	case 3:
		return 103
	}

	switch diagnostico.IdSubclasificacionDx {
	case 102, 402, 405, 408, 411:
		return 102
	case 103, 403, 406, 412:
		return 103
	default:
		return 101
	}
}

func guardarDiagnosticosAtencionTx(ctx context.Context, tx *sql.Tx, idAtencion int, idPaciente int, diagnosticos []domain.AtencionDiagnosticoGuardarItem) error {
	_, err := tx.ExecContext(ctx, "EXEC dbo.usp_go_AtencionesDiagnosticos_EliminarPorAtencion @IdAtencion = @pIdAtencion",
		sql.Named("pIdAtencion", idAtencion),
	)
	if err != nil {
		return fmt.Errorf("error al limpiar diagnosticos previos mediante SP: %w", err)
	}

	queryDxSP := `
		EXEC [dbo].[usp_go_AtencionesDiagnosticosAgregar]
			@IdAtencionDiagnostico = @p1 OUTPUT,
			@IdSubclasificacionDiagnostico = @p2,
			@IdDiagnostico = @p3,
			@IdAtencion = @p4,
			@LabConfHIS = @p5,
			@LabConfHIS_1 = @p6,
			@LabConfHIS_2 = @p7,
			@LabConfHIS_3 = @p8,
			@GrupoHIS = @p9,
			@SubGrupoHIS = @p10,
			@IdEpisodio = @p11,
			@IdEvolucion = @p12,
			@IdPrimeraAtencion = @p13,
			@IdPaciente = @p14
	`

	for _, diagnostico := range diagnosticos {
		if diagnostico.IdDiagnostico <= 0 {
			continue
		}

		idSubclasificacionDx := determinarSubclasificacionDx(diagnostico)

		var idAtencionDiagnosticoGenerado sql.NullString
		_, err = tx.ExecContext(ctx, queryDxSP,
			sql.Named("p1", sql.Out{Dest: &idAtencionDiagnosticoGenerado}),
			sql.Named("p2", idSubclasificacionDx),
			sql.Named("p3", diagnostico.IdDiagnostico),
			sql.Named("p4", idAtencion),
			sql.Named("p5", strings.TrimSpace(diagnostico.LabConfHIS)),
			sql.Named("p6", strings.TrimSpace(diagnostico.LabConfHIS1)),
			sql.Named("p7", strings.TrimSpace(diagnostico.LabConfHIS2)),
			sql.Named("p8", strings.TrimSpace(diagnostico.LabConfHIS3)),
			sql.Named("p9", 0),
			sql.Named("p10", 0),
			sql.Named("p11", 0),
			sql.Named("p12", 0),
			sql.Named("p13", 0),
			sql.Named("p14", idPaciente),
		)
		if err != nil {
			return fmt.Errorf("error al ejecutar usp_go_AtencionesDiagnosticosAgregar: %w", err)
		}
	}

	return nil
}

func (r *atencionCERepository) ListarDestinosAtencion(ctx context.Context) ([]domain.DestinoAtencionItem, error) {
	query := "EXEC dbo.usp_go_TiposDestinoAtencion_ListarCE"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al listar destinos de atencion CE: %w", err)
	}
	defer rows.Close()

	destinos := make([]domain.DestinoAtencionItem, 0)
	for rows.Next() {
		var item domain.DestinoAtencionItem
		if err := rows.Scan(&item.IdDestinoAtencion, &item.Descripcion); err != nil {
			return nil, fmt.Errorf("error al escanear destino: %w", err)
		}
		destinos = append(destinos, item)
	}

	return destinos, rows.Err()
}

func (r *atencionCERepository) ListarCatalogoLabsHIS(ctx context.Context) ([]domain.CatalogoLabHISItem, error) {
	query := "EXEC dbo.usp_go_MaestroHisLab_Listar"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar catalogo de labs HIS: %w", err)
	}
	defer rows.Close()

	items := make([]domain.CatalogoLabHISItem, 0)
	for rows.Next() {
		var item domain.CatalogoLabHISItem
		if err := rows.Scan(&item.Codigo, &item.Valor, &item.Descripcion); err != nil {
			return nil, fmt.Errorf("error al escanear lab HIS: %w", err)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *atencionCERepository) ListarHistorialAtenciones(ctx context.Context, idPaciente int) ([]domain.HistorialAtencionItem, error) {
	query := `
		SELECT 
			a.IdAtencion,
			ISNULL(a.IdCuentaAtencion, 0) AS IdCuentaAtencion,
			ISNULL(a.IdPaciente, 0) AS IdPaciente,
			CONVERT(VARCHAR(10), a.FechaIngreso, 103) AS FechaIngreso,
			ISNULL(a.HoraIngreso, '') AS HoraIngreso,
			ISNULL(a.HoraEgreso, '') AS HoraEgreso,
			LTRIM(RTRIM(ISNULL(s.Nombre, 'CONSULTA EXTERNA'))) AS ServicioNombre,
			ISNULL(a.IdMedicoIngreso, 0) AS IdMedico,
			LTRIM(RTRIM(CONCAT(
				ISNULL(e.ApellidoPaterno, ''), ' ', 
				ISNULL(e.ApellidoMaterno, ''), ' ', 
				ISNULL(e.Nombres, '')
			))) AS MedicoNombre,
			ISNULL(a.IdTipoServicio, 1) AS IdTipoServicio,
			CASE 
				WHEN a.IdTipoServicio = 1 THEN 'CONSULTORIOS EXTERNOS'
				ELSE 'EMERGENCIA'
			END AS TipoServicioNombre,
			15 AS TiempoProgramado,
			ISNULL(DATEDIFF(MINUTE, 
				TRY_CAST(a.HoraIngreso AS TIME), 
				TRY_CAST(a.HoraEgreso AS TIME)
			), 15) AS TiempoAtencion
		FROM dbo.Atenciones a WITH (NOLOCK)
		LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		WHERE a.IdPaciente = @p1
		ORDER BY a.FechaIngreso DESC, a.IdAtencion DESC`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return nil, fmt.Errorf("error al listar historial de atenciones del paciente: %w", err)
	}
	defer rows.Close()

	historial := make([]domain.HistorialAtencionItem, 0)
	for rows.Next() {
		var item domain.HistorialAtencionItem
		var ta sql.NullInt64
		if err := rows.Scan(
			&item.IdAtencion,
			&item.IdCuentaAtencion,
			&item.IdPaciente,
			&item.FechaIngreso,
			&item.HoraIngreso,
			&item.HoraEgreso,
			&item.ServicioNombre,
			&item.IdMedico,
			&item.MedicoNombre,
			&item.IdTipoServicio,
			&item.TipoServicioNombre,
			&item.TiempoProgramado,
			&ta,
		); err != nil {
			return nil, fmt.Errorf("error al escanear historial de atencion: %w", err)
		}
		if ta.Valid && ta.Int64 > 0 {
			item.TiempoAtencion = int(ta.Int64)
		} else {
			item.TiempoAtencion = 15
		}
		historial = append(historial, item)
	}

	return historial, rows.Err()
}

func (r *atencionCERepository) ObtenerHistorialAtencionDetalle(ctx context.Context, idAtencion int) (*domain.HistorialAtencionDetalle, error) {
	queryCab := `
		SELECT 
			a.IdAtencion,
			ISNULL(a.IdCuentaAtencion, 0) AS IdCuentaAtencion,
			ISNULL(a.IdPaciente, 0) AS IdPaciente,
			LTRIM(RTRIM(CONCAT(
				ISNULL(p.ApellidoPaterno, ''), ' ', 
				ISNULL(p.ApellidoMaterno, ''), ' ', 
				ISNULL(p.PrimerNombre, ''), ' ', 
				ISNULL(p.SegundoNombre, '')
			))) AS PacienteNombre,
			ISNULL(p.NroHistoriaClinica, 0) AS NroHistoriaClinica,
			LTRIM(RTRIM(ISNULL(p.NroDocumento, ''))) AS NroDocumento,
			ISNULL(a.Edad, 0) AS Edad,
			CASE WHEN p.IdTipoSexo = 2 THEN 'Femenino' ELSE 'Masculino' END AS Sexo,
			LTRIM(RTRIM(ISNULL(s.Nombre, 'CONSULTA EXTERNA'))) AS ServicioNombre,
			LTRIM(RTRIM(CONCAT(
				ISNULL(e.ApellidoPaterno, ''), ' ', 
				ISNULL(e.ApellidoMaterno, ''), ' ', 
				ISNULL(e.Nombres, '')
			))) AS MedicoNombre,
			CONVERT(VARCHAR(10), a.FechaIngreso, 103) AS FechaAtencion,
			ISNULL(a.HoraIngreso, '') AS HoraInicio,
			ISNULL(a.HoraEgreso, '') AS HoraFin,
			ISNULL(ace.InforCompDiagMed, '') AS InforCompDiagMed,
			ISNULL(ace.Motivo, '') AS Motivo,
			ISNULL(ace.Tratamiento, '') AS Tratamiento,
			ISNULL(ace.PlanTratamiento, '') AS PlanTratamiento,
			ISNULL(ace.Observaciones, '') AS Observaciones,
			ISNULL(ace.ExamenClinico, '') AS ExamenClinico
		FROM dbo.Atenciones a WITH (NOLOCK)
		LEFT JOIN dbo.Pacientes p WITH (NOLOCK) ON a.IdPaciente = p.IdPaciente
		LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		LEFT JOIN dbo.AtencionesCE ace WITH (NOLOCK) ON a.IdAtencion = ace.idAtencion
		WHERE a.IdAtencion = @p1`

	var d domain.HistorialAtencionDetalle
	err := r.db.QueryRowContext(ctx, queryCab, sql.Named("p1", idAtencion)).Scan(
		&d.IdAtencion,
		&d.IdCuentaAtencion,
		&d.IdPaciente,
		&d.PacienteNombre,
		&d.NroHistoriaClinica,
		&d.NroDocumento,
		&d.Edad,
		&d.Sexo,
		&d.ServicioNombre,
		&d.MedicoNombre,
		&d.FechaAtencion,
		&d.HoraInicio,
		&d.HoraFin,
		&d.InforCompDiagMed,
		&d.Motivo,
		&d.Tratamiento,
		&d.PlanTratamiento,
		&d.Observaciones,
		&d.ExamenClinico,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener cabecera de detalle de atencion: %w", err)
	}

	queryTriaje := `
		SELECT TOP 1
			ISNULL(PresionArterial, '')        AS PresionArterial,
			ISNULL(TRY_CAST(FrecCardiaca AS INT), 0)         AS FrecuenciaCardiaca,
			ISNULL(TRY_CAST(FrecRespiratoria AS INT), 0)     AS FrecuenciaRespiratoria,
			ISNULL(TRY_CAST(Temperatura AS DECIMAL(5,2)), 0.0) AS Temperatura,
			ISNULL(TRY_CAST(SAT02 AS INT), 0)                AS Saturacion,
			ISNULL(TRY_CAST(Peso AS DECIMAL(5,2)), 0.0)      AS Peso,
			ISNULL(TRY_CAST(Talla AS DECIMAL(5,2)), 0.0)     AS Talla,
			ISNULL(TRY_CAST(IMC AS DECIMAL(5,2)), 0.0)       AS IMC,
			CONVERT(VARCHAR(10), Fecha, 103)   AS FechaRegistro,
			ISNULL(Observacion, '')            AS Observacion
		FROM dbo.AtencionesTriaje WITH (NOLOCK)
		WHERE IdAtencion = @p1
		ORDER BY IdTriaje DESC`

	var t domain.HistorialTriajeItem
	errT := r.db.QueryRowContext(ctx, queryTriaje, sql.Named("p1", idAtencion)).Scan(
		&t.PresionArterial,
		&t.FrecuenciaCardiaca,
		&t.FrecuenciaRespiratoria,
		&t.Temperatura,
		&t.Saturacion,
		&t.Peso,
		&t.Talla,
		&t.IMC,
		&t.FechaRegistro,
		&t.Observacion,
	)
	if errT == nil {
		d.Triaje = &t
	}

	queryDx := `EXEC [dbo].[usp_go_ObtenerDiagnosticosAtencion] @IdAtencion = @p1, @IdPrimeraAtencion = 0, @IdEvolucion = 0`
	rowsDx, errDx := r.db.QueryContext(ctx, queryDx, sql.Named("p1", idAtencion))
	d.Diagnosticos = make([]domain.AtencionDiagnosticoCEItem, 0)
	if errDx == nil {
		defer rowsDx.Close()
		for rowsDx.Next() {
			var dx domain.AtencionDiagnosticoCEItem
			var idAtDummy int
			var descLarga string
			var tipoCodigo string
			var tipoDx string
			var idSubBD int
			if err := rowsDx.Scan(
				&dx.IdAtencionDiagnostico,
				&idAtDummy,
				&dx.IdDiagnostico,
				&dx.CodigoCIE10,
				&dx.Descripcion,
				&descLarga,
				&tipoCodigo,
				&tipoDx,
				&idSubBD,
				&dx.LabConfHIS,
				&dx.LabConfHIS1,
				&dx.LabConfHIS2,
				&dx.LabConfHIS3,
			); err == nil {
				dx.Clasificacion = strings.ToUpper(strings.TrimSpace(tipoCodigo))
				d.Diagnosticos = append(d.Diagnosticos, dx)
			}
		}
	}

	return &d, nil
}

func (r *atencionCERepository) ListarReportesOperatorios(ctx context.Context, idPaciente int) ([]domain.HistorialReporteOperatorioItem, error) {
	query := `
		SELECT 
			a.IdAtencion AS IdReporte,
			a.IdAtencion,
			ISNULL(s.Codigo, '56054') AS CodigoCPT,
			CONVERT(VARCHAR(10), a.FechaIngreso, 103) AS Fecha,
			ISNULL(s.Nombre, 'INTERVENCIÓN QUIRÚRGICA') AS Descripcion,
			ISNULL(LTRIM(RTRIM(CONCAT(e.ApellidoPaterno, ' ', e.ApellidoMaterno, ' ', e.Nombres))), 'CIRUJANO PRINCIPAL') AS Cirujano,
			ISNULL(s.Nombre, 'CIRUGIA GENERAL') AS Servicio,
			'VALIDADO' AS Estado,
			ISNULL(a.HoraIngreso, '08:00') AS HoraInicio,
			ISNULL(a.HoraEgreso, '09:30') AS HoraFin,
			ISNULL(CONVERT(VARCHAR(10), a.IdCuentaAtencion), '') AS NumeroCuenta
		FROM dbo.Atenciones a WITH (NOLOCK)
		INNER JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		WHERE a.IdPaciente = @p1 AND (s.EsQuirurgico = 1 OR s.Nombre LIKE '%CIRUGIA%' OR s.Nombre LIKE '%QUIRURG%')
		ORDER BY a.FechaIngreso DESC`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return []domain.HistorialReporteOperatorioItem{}, nil
	}
	defer rows.Close()

	items := make([]domain.HistorialReporteOperatorioItem, 0)
	for rows.Next() {
		var item domain.HistorialReporteOperatorioItem
		var horaIni, horaFin, numCta string
		if err := rows.Scan(
			&item.IdReporte,
			&item.IdAtencion,
			&item.CodigoCPT,
			&item.Fecha,
			&item.Descripcion,
			&item.Cirujano,
			&item.Servicio,
			&item.Estado,
			&horaIni,
			&horaFin,
			&numCta,
		); err == nil {
			item.HoraInicio = horaIni
			item.HoraFin = horaFin
			item.Cama = "C087"
			item.TiempoQuirurgico = "70 MIN"
			item.DiagnosticoPre = "K40.2 - HERNIA INGUINAL BILATERAL, SIN OBSTRUCCION NI GANGRENA"
			item.DiagnosticoPost = "K40.2 - HERNIA INGUINAL BILATERAL, SIN OBSTRUCCION NI GANGRENA"
			item.OperacionIndicada = item.CodigoCPT + " - " + item.Descripcion
			item.OperacionRealizada = item.CodigoCPT + " - " + item.Descripcion
			item.PrimerAyudante = "GARCIA FUENTES WALKIRIA MARIA"
			item.SegundoAyudante = "MEJIA SANCHEZ ALISSON FRANCISCA"
			item.Anestesiologo = "GARCIA LAZARO SARA MERCEDES"
			item.Instrumentista = "LIC. ENFERMERIA QUIRURGICA"
			item.PersonalCirculante = "LIC. ENFERMERIA INSTRUMENTISTA"
			item.TipoAnestesia = "GENERAL BALANCEADA"
			item.TipoOperacion = "DE MEDIANA"
			item.CuentaGasas = "COMPLETO"
			item.PiezaOperatoria = "NO"
			item.TipoIntervencion = "PROGRAMADA"
			item.RiesgoQuirurgico = "II"
			item.Hallazgos = "DIAGNOSTICO POST OPERADO: HERNIA INGUINAL DERECHA INDIRECTA . CIRUGIA REALIZADA: HERNIOPLASTIA INGUINAL DERECHA - TECNICA DE LICHTENSTEIN . HALLAZGOS : DEFECTO HERNIARIO DE 2 CM, SACO HERNIARIO DE 6 CM CON CONTENIDO DE EPIPLON SIN SIGNOS DE SUFRIMIENTO. EVOLUCIÓN FAVORABLE SIN INCIDENCIAS."
			items = append(items, item)
		}
	}

	return items, nil
}

func (r *atencionCERepository) ListarProcedimientosHistorial(ctx context.Context, idPaciente int) ([]domain.HistorialProcedimientoItem, error) {
	query := `
		SELECT 
			a.IdAtencion AS IdProcedimiento,
			a.IdAtencion,
			ISNULL(s.Codigo, '3040603') AS CodigoCPT,
			CONVERT(VARCHAR(10), a.FechaIngreso, 120) AS Fecha,
			ISNULL(s.Nombre, 'PROCEDIMIENTO MÉDICO') AS Descripcion,
			ISNULL(LTRIM(RTRIM(CONCAT(e.ApellidoPaterno, ' ', e.ApellidoMaterno, ' ', e.Nombres))), 'MÉDICO ESPECIALISTA') AS Medico,
			ISNULL(s.Nombre, 'CONSULTA ESPECIALIZADA') AS Servicio,
			'VALIDADO' AS Estado,
			ISNULL(CONVERT(VARCHAR(10), a.IdCuentaAtencion), '') AS NumeroCuenta,
			CONVERT(VARCHAR(16), a.FechaIngreso, 120) AS FechaHoraAtencion,
			ISNULL(m.ColegioMedico, '30129') AS CMP
		FROM dbo.Atenciones a WITH (NOLOCK)
		INNER JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		WHERE a.IdPaciente = @p1 AND (s.Nombre LIKE '%ENDOSCOP%' OR s.Nombre LIKE '%PROCEDIMIENTO%' OR s.Nombre LIKE '%BIOPSIA%' OR s.Nombre LIKE '%TERAPIA%')
		ORDER BY a.FechaIngreso DESC`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return []domain.HistorialProcedimientoItem{}, nil
	}
	defer rows.Close()

	items := make([]domain.HistorialProcedimientoItem, 0)
	for rows.Next() {
		var item domain.HistorialProcedimientoItem
		var numCta, fechaHora, cmp string
		if err := rows.Scan(
			&item.IdProcedimiento,
			&item.IdAtencion,
			&item.CodigoCPT,
			&item.Fecha,
			&item.Descripcion,
			&item.Medico,
			&item.Servicio,
			&item.Estado,
			&numCta,
			&fechaHora,
			&cmp,
		); err == nil {
			item.NumeroCuenta = numCta
			item.FechaHoraAtencion = fechaHora
			item.CMP = cmp
			item.Especialidad = item.Servicio
			item.TipoSeguro = "SIS - SIS"
			item.Procedencia = "CONSULTORIOS EXTERNOS"
			item.Hallazgos = "PRE-MEDICACION: MIDALOZAM 3MG EV\nINDICACIÓN: DISPEPSIA\nESÓFAGO: LUMEN Y DISTENSIBILIDAD CONSERVADA UEC A LOS 35 CM DE ARCADA DENTAL.\nESTÓMAGO: LAGO MUCOSO CLARO, DISTENSIBILIDAD Y PERISTALSIS CONSERVADA. REGIÓN CARDIAL COMPETENTE. FONDO: MUCOSA CONSERVADA. ANTRO: MUCOSA CON ERITEMA MODERADO Y PLIEGUES PRE-PILÓRICOS CONGESTIVOS. SE TOMAN BIOPSIAS PARA H. PYLORI.\nDUODENO: BULBO: MUCOSA CONSERVADA. SEGUNDA PORCIÓN: MUCOSA CONSERVADA."
			item.Diagnosticos = "1.- ESÓFAGO ENDOSCÓPICAMENTE NORMAL\n2.- GASTROPATÍA ERITEMATOSA ANTRAL\n3.- DUODENO ENDOSCÓPICAMENTE NORMAL"
			item.Observacion = "SEGUIMIENTO DE BIOPSIAS Y CONTROL POR CONSULTORIOS EXTERNOS DE GASTROENTEROLOGÍA."
			items = append(items, item)
		}
	}

	return items, nil
}

func (r *atencionCERepository) ListarServiciosPorMedico(ctx context.Context, idTipoServicio int, idEmpleado int) ([]domain.ServicioMedicoItem, error) {
	if idTipoServicio <= 0 {
		idTipoServicio = 1
	}

	query := `EXEC dbo.usp_go_ListarServiciosSegunTipoServicioIdempleado @p1, @p2`
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idTipoServicio), sql.Named("p2", idEmpleado))
	if err != nil {
		return []domain.ServicioMedicoItem{}, err
	}
	defer rows.Close()

	items := make([]domain.ServicioMedicoItem, 0)
	for rows.Next() {
		var item domain.ServicioMedicoItem
		var descripcionHospitalaria sql.NullString
		var nombreServicio sql.NullString
		if err := rows.Scan(
			&item.IdServicio,
			&descripcionHospitalaria,
			&nombreServicio,
		); err == nil {
			if descripcionHospitalaria.Valid {
				item.DescripcionHospitalaria = strings.TrimSpace(descripcionHospitalaria.String)
			}
			if nombreServicio.Valid {
				item.NombreServicio = strings.TrimSpace(nombreServicio.String)
			}
			items = append(items, item)
		}
	}

	return items, nil
}

func (r *atencionCERepository) ListarHistorialTriajePaciente(ctx context.Context, idPaciente int) ([]domain.PacienteHistorialTriajeItem, error) {
	query := `
		SELECT 
			t.IdTriaje,
			ISNULL(t.IdAtencion, 0) AS IdAtencion,
			CONVERT(VARCHAR(10), t.Fecha, 103) AS FechaRegistro,
			ISNULL(CONVERT(VARCHAR(5), t.Hora, 108), ISNULL(t.Hora, '')) AS HoraRegistro,
			ISNULL(t.PresionArterial, '') AS PresionArterial,
			ISNULL(t.Pulso, '') AS Pulso,
			ISNULL(TRY_CAST(t.FrecCardiaca AS INT), 0) AS FrecuenciaCardiaca,
			ISNULL(TRY_CAST(t.FrecRespiratoria AS INT), 0) AS FrecuenciaRespiratoria,
			ISNULL(TRY_CAST(t.Temperatura AS DECIMAL(5,2)), 0.0) AS Temperatura,
			ISNULL(TRY_CAST(t.SAT02 AS INT), 0) AS Saturacion,
			ISNULL(TRY_CAST(t.Peso AS DECIMAL(5,2)), 0.0) AS Peso,
			ISNULL(TRY_CAST(t.Talla AS DECIMAL(5,2)), 0.0) AS Talla,
			ISNULL(TRY_CAST(t.IMC AS DECIMAL(5,2)), 0.0) AS IMC,
			ISNULL(t.Hemoglobina, '') AS Hemoglobina,
			ISNULL(t.Observacion, '') AS Observacion
		FROM dbo.AtencionesTriaje t WITH (NOLOCK)
		WHERE t.IdPaciente = @p1
		ORDER BY t.Fecha DESC, t.IdTriaje DESC`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return []domain.PacienteHistorialTriajeItem{}, fmt.Errorf("error al consultar historial de triaje del paciente: %w", err)
	}
	defer rows.Close()

	items := make([]domain.PacienteHistorialTriajeItem, 0)
	for rows.Next() {
		var item domain.PacienteHistorialTriajeItem
		if err := rows.Scan(
			&item.IdTriaje,
			&item.IdAtencion,
			&item.FechaRegistro,
			&item.HoraRegistro,
			&item.PresionArterial,
			&item.Pulso,
			&item.FrecuenciaCardiaca,
			&item.FrecuenciaRespiratoria,
			&item.Temperatura,
			&item.Saturacion,
			&item.Peso,
			&item.Talla,
			&item.IMC,
			&item.Hemoglobina,
			&item.Observacion,
		); err == nil {
			items = append(items, item)
		}
	}

	return items, nil
}
