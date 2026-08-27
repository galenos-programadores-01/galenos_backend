package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sqlServerEvolucionRepository struct {
	db *sql.DB
}

func NewSqlServerEvolucionRepository(db *sql.DB) output.EvolucionRepository {
	return &sqlServerEvolucionRepository{db: db}
}

func (r *sqlServerEvolucionRepository) ListPatients(ctx context.Context, fini, ffin string, idUsuario int) ([]domain.PatientListItem, error) {
	// Se usa el SP usp_go_ListarPacientesSegunTipoServicio, que requiere
	// @IdTipoServicio (2 = Emergencia), @Fecha (datetime), @Filtro y @IdUsuario.
	// El handler ya validó el formato de fini antes de llegar aquí.
	fecha, err := time.Parse("2006-01-02", fini)
	if err != nil {
		return nil, fmt.Errorf("fini inválido: %w", err)
	}
	query := `EXEC [dbo].[usp_go_ListarPacientesSegunTipoServicio] @IdTipoServicio = @p1, @Fecha = @p2, @Filtro = @p3, @IdUsuario = @p4`
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", 2),
		sql.Named("p2", fecha),
		sql.Named("p3", ""),
		sql.Named("p4", idUsuario),
	)
	if err != nil {
		return nil, fmt.Errorf("error querying patients for tray: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("error reading patient tray maps: %w", err)
	}

	var patients []domain.PatientListItem
	for _, m := range maps {
		patients = append(patients, mapToPatientListItem(m))
	}
	return patients, nil
}

func mapToPatientListItem(m map[string]any) domain.PatientListItem {
	var p domain.PatientListItem

	// El SP usp_go_ListarPacientesSegunTipoServicio retorna: IdEpisodio,
	// IdAtencion, IdCuentaAtencion, Paciente, NroHistoriaClinica, Servicio,
	// Sexo y cama.
	p.IdRegAtencion = getIntFallback(m, "idRegAtencion", "IdAtencion")
	p.IdPaciente = getIntFallback(m, "IdPaciente", "IdEpisodio")
	p.Historia = getStringFallback(m, "historia", "NroHistoriaClinica", "N/A")
	p.Nombre = getNombrePaciente(m)
	p.Edad = getStringFallback(m, "edad", "", "N/A")
	p.Sexo = getSexoFallback(m)
	p.Ubicacion = getStringFallback(m, "ubicacion", "Servicio", "Emergencia")
	p.Cama = getStringFallback(m, "cama", "", "NS")
	p.Estado = getStringFallback(m, "estado", "", "Pendiente")

	return p
}

func getIntFallback(m map[string]any, key1, key2 string) int {
	if val, ok := m[key1]; ok && val != nil {
		return int(val.(int64))
	}
	if key2 != "" {
		if val, ok := m[key2]; ok && val != nil {
			return int(val.(int64))
		}
	}
	return 0
}

func getStringFallback(m map[string]any, key1, key2, fallback string) string {
	if ptr := rowString(m, key1); ptr != nil && *ptr != "" {
		return *ptr
	}
	if key2 != "" {
		if ptr := rowString(m, key2); ptr != nil && *ptr != "" {
			return *ptr
		}
	}
	return fallback
}

func getNombrePaciente(m map[string]any) string {
	if ptr := rowString(m, "nombre", "Paciente"); ptr != nil && *ptr != "" {
		return *ptr
	}
	pat := getStringFallback(m, "ApellidoPaterno", "", "")
	mat := getStringFallback(m, "ApellidoMaterno", "", "")
	nom := getStringFallback(m, "PrimerNombre", "", "")
	return fmt.Sprintf("%s %s, %s", pat, mat, nom)
}

func getSexoFallback(m map[string]any) string {
	if ptr := rowString(m, "sexo"); ptr != nil && *ptr != "" {
		return *ptr
	}
	if ptr := rowInt64(m, "IdTipoSexo"); ptr != nil {
		return fmt.Sprintf("%d", *ptr)
	}
	return "0"
}

func (r *sqlServerEvolucionRepository) ListEvoluciones(ctx context.Context, idRegAtencion int) ([]domain.EvolucionFirma, error) {
	query := `EXEC [dbo].[webEvolucionesFirmaListarIdRegAtencion] @IdRegAtencion = @p1, @NombreDocumento = 'EvolucionMedica'`
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idRegAtencion))
	if err != nil {
		return nil, fmt.Errorf("error querying evolutions: %w", err)
	}
	defer rows.Close()

	var evolutions []domain.EvolucionFirma
	for rows.Next() {
		var e domain.EvolucionFirma
		var rBase, nArchivo, doc, data, fRegistro sql.NullString
		var idEmpReg, idEmpMod, idEmpAnula, estado sql.NullInt64
		var fMod, fAnul sql.NullTime

		if err := rows.Scan(
			&e.IdRegAtencion,
			&e.IdFirma,
			&rBase,
			&nArchivo,
			&doc,
			&data,
			&idEmpReg,
			&idEmpMod,
			&idEmpAnula,
			&fRegistro,
			&fMod,
			&fAnul,
			&estado,
		); err != nil {
			return nil, fmt.Errorf("error scanning evolution: %w", err)
		}

		if rBase.Valid {
			e.RutaBase = rBase.String
		}
		if nArchivo.Valid {
			e.NombreArchivo = nArchivo.String
		}
		if doc.Valid {
			e.NombreDocumento = doc.String
		}
		if data.Valid {
			e.DataB64 = data.String
		}
		if idEmpReg.Valid {
			e.IdEmpleadoRegistra = int(idEmpReg.Int64)
		}
		if fRegistro.Valid {
			e.FechaRegistro = fRegistro.String
		}
		if estado.Valid {
			e.Estado = int(estado.Int64)
		}

		evolutions = append(evolutions, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando evoluciones: %w", err)
	}

	return evolutions, nil
}

func (r *sqlServerEvolucionRepository) SaveEvolucion(ctx context.Context, evolution domain.EvolucionFirma) error {
	query := `
		EXEC [dbo].[Web_sp_GuardarEvolucionFirma] 
			@IdRegAtencion = @p1, 
			@RutaBase = @p2, 
			@NombreArchivo = @p3, 
			@NombreDocumento = @p4, 
			@DataB64 = @p5, 
			@IdEmpleadoRegistra = @p6, 
			@IdEmpleadoModifica = @p7, 
			@Estado = @p8
	`
	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", evolution.IdRegAtencion),
		sql.Named("p2", "evols/"),
		sql.Named("p3", fmt.Sprintf("evol_%d_%d.json", evolution.IdRegAtencion, time.Now().Unix())),
		sql.Named("p4", evolution.NombreDocumento),
		sql.Named("p5", evolution.DataB64),
		sql.Named("p6", evolution.IdEmpleadoRegistra),
		sql.Named("p7", evolution.IdEmpleadoRegistra),
		sql.Named("p8", 1),
	)
	if err != nil {
		return fmt.Errorf("error saving evolution: %w", err)
	}
	return nil
}

func (r *sqlServerEvolucionRepository) ListBandeja(ctx context.Context, fechaInicio, fechaFin, filtro string) ([]domain.EvolucionBandejaItem, error) {
	query := `EXEC dbo.usp_go_EvolucionesMedicas_Bandeja @FechaInicio = @p1, @FechaFin = @p2, @Filtro = @p3`

	var fInicio, fFin interface{} = nil, nil
	if fechaInicio != "" {
		fInicio = fechaInicio
	}
	if fechaFin != "" {
		fFin = fechaFin
	}

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", fInicio),
		sql.Named("p2", fFin),
		sql.Named("p3", filtro),
	)
	if err != nil {
		return nil, fmt.Errorf("error al consultar bandeja evoluciones: %w", err)
	}
	defer rows.Close()

	var list []domain.EvolucionBandejaItem
	for rows.Next() {
		var item domain.EvolucionBandejaItem
		var idEpisodio, idCuentaAtencion, escalaDolor, glasgow, paSistolica, paDiastolica, fc, fr, satO2, idEstadoClinico, idPronostico, estadoFirma, usrCreacion, estRegistro sql.NullInt64
		var temp sql.NullFloat64
		var fAtencion, paciente, documento, motivo, fFirma, fCreacion, eqCreacion sql.NullString

		if err := rows.Scan(
			&item.IdEvolucion, &idEpisodio, &item.NroAtencion, &fAtencion,
			&item.IdPaciente, &paciente, &documento, &idCuentaAtencion, &motivo,
			&escalaDolor, &glasgow, &paSistolica, &paDiastolica, &fc, &fr,
			&temp, &satO2, &idEstadoClinico, &idPronostico, &estadoFirma,
			&fFirma, &usrCreacion, &fCreacion, &eqCreacion, &estRegistro,
		); err != nil {
			return nil, fmt.Errorf("error escaneando item de bandeja evoluciones: %w", err)
		}

		if idEpisodio.Valid {
			v := int(idEpisodio.Int64)
			item.IdEpisodio = &v
		}
		if fAtencion.Valid {
			item.FechaAtencion = fAtencion.String
		}
		if paciente.Valid {
			item.Paciente = paciente.String
		}
		if documento.Valid {
			item.Documento = documento.String
		}
		if idCuentaAtencion.Valid {
			item.IdCuentaAtencion = int(idCuentaAtencion.Int64)
		}
		if motivo.Valid {
			item.Motivo = motivo.String
		}
		if escalaDolor.Valid {
			v := int(escalaDolor.Int64)
			item.EscalaDolor = &v
		}
		if glasgow.Valid {
			v := int(glasgow.Int64)
			item.Glasgow = &v
		}
		if paSistolica.Valid {
			v := int(paSistolica.Int64)
			item.PASistolica = &v
		}
		if paDiastolica.Valid {
			v := int(paDiastolica.Int64)
			item.PADiastolica = &v
		}
		if fc.Valid {
			v := int(fc.Int64)
			item.FrecuenciaCardiaca = &v
		}
		if fr.Valid {
			v := int(fr.Int64)
			item.FrecuenciaRespiratoria = &v
		}
		if temp.Valid {
			v := temp.Float64
			item.Temperatura = &v
		}
		if satO2.Valid {
			v := int(satO2.Int64)
			item.SaturacionOxigeno = &v
		}
		if idEstadoClinico.Valid {
			v := int(idEstadoClinico.Int64)
			item.IdEstadoClinico = &v
		}
		if idPronostico.Valid {
			v := int(idPronostico.Int64)
			item.IdPronostico = &v
		}
		if estadoFirma.Valid {
			v := int(estadoFirma.Int64)
			item.EstadoFirma = &v
		}
		if fFirma.Valid {
			v := fFirma.String
			item.FechaFirma = &v
		}
		if usrCreacion.Valid {
			v := int(usrCreacion.Int64)
			item.UsuarioCreacion = &v
		}
		if fCreacion.Valid {
			v := fCreacion.String
			item.FechaCreacion = &v
		}
		if eqCreacion.Valid {
			v := eqCreacion.String
			item.EquipoCreacion = &v
		}
		if estRegistro.Valid {
			v := int(estRegistro.Int64)
			item.EstadoRegistro = &v
		}

		list = append(list, item)
	}

	return list, rows.Err()
}

func (r *sqlServerEvolucionRepository) InsertEvolucionMedica(ctx context.Context, item domain.EvolucionMedicaInsert) (int, string, error) {
	query := `
		EXEC dbo.usp_go_EvolucionesMedicas_Insertar
			@IdAtencion = @p1,
			@IdPaciente = @p2,
			@IdMedico = @p3,
			@FechaAtencion = @p4,
			@IdTipoGravedad = @p5,
			@MotivoConsulta = @p6,
			@TiempoEnfermedad = @p7,
			@Anamnesis = @p8,
			@EscalaDolor = @p9,
			@Glasgow = @p10,
			@PASistolica = @p11,
			@PADiastolica = @p12,
			@FrecuenciaCardiaca = @p13,
			@FrecuenciaRespiratoria = @p14,
			@Temperatura = @p15,
			@SaturacionOxigeno = @p16,
			@Peso = @p17,
			@Talla = @p18,
			@IMC = @p19,
			@Glicemia = @p20,
			@ExamenFisicoGeneral = @p21,
			@ExamenFisicoPiel = @p22,
			@ExamenFisicoCabezaCuello = @p23,
			@ExamenFisicoToraxPulmon = @p24,
			@ExamenFisicoCorazon = @p25,
			@ExamenFisicoAbdomen = @p26,
			@ExamenFisicoGenitourinario = @p27,
			@ExamenFisicoExtremidadesOsteomuscular = @p28,
			@ExamenFisicoNeurologicoMental = @p29,
			@IdEstadoClinico = @p30,
			@IdPronostico = @p31,
			@IndicacionDieta = @p32,
			@IndicacionReposo = @p33,
			@IndicacionHidratacion = @p34,
			@IndicacionOxigeno = @p35,
			@IndicacionRestriccion = @p36,
			@Sugerencia = @p37,
			@UsuarioCreacion = @p38,
			@EquipoCreacion = @p39,
			@EstadoRegistro = @p40,
			@EstadoFirma = @p41,
			@IdEvolucion = @p42 OUTPUT,
			@Mensaje = @p43 OUTPUT
	`

	var idEvolucion sql.NullInt64
	var mensaje string

	err := r.db.QueryRowContext(ctx, query,
		sql.Named("p1", item.IdAtencion),
		sql.Named("p2", item.IdPaciente),
		sql.Named("p3", item.IdMedico),
		sql.Named("p4", item.FechaAtencion),
		sql.Named("p5", item.IdTipoGravedad),
		sql.Named("p6", item.MotivoConsulta),
		sql.Named("p7", item.TiempoEnfermedad),
		sql.Named("p8", item.Anamnesis),
		sql.Named("p9", item.EscalaDolor),
		sql.Named("p10", item.Glasgow),
		sql.Named("p11", item.PASistolica),
		sql.Named("p12", item.PADiastolica),
		sql.Named("p13", item.FrecuenciaCardiaca),
		sql.Named("p14", item.FrecuenciaRespiratoria),
		sql.Named("p15", item.Temperatura),
		sql.Named("p16", item.SaturacionOxigeno),
		sql.Named("p17", item.Peso),
		sql.Named("p18", item.Talla),
		sql.Named("p19", item.IMC),
		sql.Named("p20", item.Glicemia),
		sql.Named("p21", item.ExamenFisicoGeneral),
		sql.Named("p22", item.ExamenFisicoPiel),
		sql.Named("p23", item.ExamenFisicoCabezaCuello),
		sql.Named("p24", item.ExamenFisicoToraxPulmon),
		sql.Named("p25", item.ExamenFisicoCorazon),
		sql.Named("p26", item.ExamenFisicoAbdomen),
		sql.Named("p27", item.ExamenFisicoGenitourinario),
		sql.Named("p28", item.ExamenFisicoExtremidadesOsteomuscular),
		sql.Named("p29", item.ExamenFisicoNeurologicoMental),
		sql.Named("p30", item.IdEstadoClinico),
		sql.Named("p31", item.IdPronostico),
		sql.Named("p32", item.IndicacionDieta),
		sql.Named("p33", item.IndicacionReposo),
		sql.Named("p34", item.IndicacionHidratacion),
		sql.Named("p35", item.IndicacionOxigeno),
		sql.Named("p36", item.IndicacionRestriccion),
		sql.Named("p37", item.Sugerencia),
		sql.Named("p38", item.UsuarioCreacion),
		sql.Named("p39", item.EquipoCreacion),
		sql.Named("p40", item.EstadoRegistro),
		sql.Named("p41", item.EstadoFirma),
		sql.Named("p42", sql.Out{Dest: &idEvolucion}),
		sql.Named("p43", sql.Out{Dest: &mensaje}),
	).Scan(&mensaje)

	if err != nil && err != sql.ErrNoRows {
		return 0, "", fmt.Errorf("error ejecutando insercion de evolucion medica: %w", err)
	}

	id := 0
	if idEvolucion.Valid {
		id = int(idEvolucion.Int64)
	}

	return id, mensaje, nil
}
