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
