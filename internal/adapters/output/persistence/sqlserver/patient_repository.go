package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type patientRepository struct {
	db *sql.DB
}

// NewPatientRepository construye el adaptador que implementa el puerto de
// salida output.PatientRepository contra la tabla pacientes en SQL Server.
func NewPatientRepository(db *sql.DB) output.PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) List(ctx context.Context, page shared.PageRequest) ([]domain.Patient, int, error) {
	const countQuery = `SELECT COUNT(1) FROM pacientes`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting patients: %w", err)
	}

	const listQuery = `
		SELECT IdPaciente, NroDocumento, NroHistoriaClinica, ApellidoPaterno,
		       ApellidoMaterno, PrimerNombre, SegundoNombre, TercerNombre, FechaNacimiento, IdTipoSexo
		FROM pacientes 
		WHERE ISNULL(NroDocumento, '') <> '' 
		ORDER BY NroDocumento
		OFFSET @Offset ROWS FETCH NEXT @PageSize ROWS ONLY`

	rows, err := r.db.QueryContext(ctx, listQuery,
		sql.Named("Offset", page.Offset()),
		sql.Named("PageSize", page.PageSize),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying patients: %w", err)
	}
	defer rows.Close()

	patients := make([]domain.Patient, 0, page.PageSize)
	for rows.Next() {
		var (
			patient         domain.Patient
			historyNumber   sql.NullString
			maternalSurname sql.NullString
			secondName      sql.NullString
			thirdName       sql.NullString
			birthDate       sql.NullTime
		)
		if err := rows.Scan(
			&patient.PatientID,
			&patient.DocumentNumber,
			&historyNumber,
			&patient.PaternalSurname,
			&maternalSurname,
			&patient.FirstName,
			&secondName,
			&thirdName,
			&birthDate,
			&patient.SexTypeID,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning patient row: %w", err)
		}
		patient.HistoryNumber = historyNumber.String
		patient.MaternalSurname = maternalSurname.String
		patient.SecondName = secondName.String
		patient.ThirdName = thirdName.String
		if birthDate.Valid {
			d := birthDate.Time
			patient.DateOfBirth = &d
		}
		patients = append(patients, patient)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating patient rows: %w", err)
	}

	return patients, total, nil
}

// GetByDocumentNumber invoca el procedimiento almacenado sp_listarPaciente,
// que recibe @NumDocumento varchar(18) y retorna ApellidoPaterno,
// ApellidoMaterno, PrimerNombre, SegundoNombre y TercerNombre del paciente
// cuyo número de documento coincide (el SP no repite el número de
// documento en el resultado, ya que es el propio parámetro de entrada).
// El nombre del procedimiento se pasa como texto de consulta y el
// parámetro va nombrado (@NumDocumento -> "NumDocumento"), lo que hace que
// el driver mssql lo despache como una llamada RPC en vez de concatenar
// SQL, evitando cualquier inyección.
func (r *patientRepository) GetByDocumentNumber(ctx context.Context, documentNumber string) (*domain.Patient, error) {
	const procedure = `sp_listarPaciente`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("NumDocumento", documentNumber))
	if err != nil {
		return nil, fmt.Errorf("calling sp_listarPaciente: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("iterating sp_listarPaciente result: %w", err)
		}
		return nil, domain.ErrPatientNotFound
	}

	var (
		patient         domain.Patient
		maternalSurname sql.NullString
		secondName      sql.NullString
		thirdName       sql.NullString
	)
	if err := rows.Scan(
		&patient.PaternalSurname,
		&maternalSurname,
		&patient.FirstName,
		&secondName,
		&thirdName,
	); err != nil {
		return nil, fmt.Errorf("scanning sp_listarPaciente row: %w", err)
	}
	patient.DocumentNumber = documentNumber
	patient.MaternalSurname = maternalSurname.String
	patient.SecondName = secondName.String
	patient.ThirdName = thirdName.String

	return &patient, nil
}

// GetByDocumentNumberAndType invoca el procedimiento almacenado
// usp_go_ListarPacientePorNroDocyTipo, que recibe @NroDocumento y
// @IdTipoDocIdentidad y retorna IdPaciente, NroDocumento,
// NroHistoriaClinica, ApellidoPaterno, ApellidoMaterno, PrimerNombre y
// SegundoNombre del paciente cuyo número de documento coincide. Los
// parámetros van nombrados, por lo que el driver mssql los despacha como
// una llamada RPC, sin concatenar SQL.
func (r *patientRepository) GetByDocumentNumberAndType(ctx context.Context, documentNumber string, documentTypeID int64) (*domain.Patient, error) {
	const procedure = `usp_go_ListarPacientePorNroDocyTipo`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("NroDocumento", documentNumber),
		sql.Named("IdTipoDocIdentidad", documentTypeID),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarPacientePorNroDocyTipo: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("iterating usp_go_ListarPacientePorNroDocyTipo result: %w", err)
		}
		return nil, domain.ErrPatientNotFound
	}

	var (
		patient           domain.Patient
		historyNumber     sql.NullString
		maternalSurname   sql.NullString
		secondName        sql.NullString
		birthDate         sql.NullTime
		homeDistrictID    sql.NullInt64
		homeCenterID      sql.NullInt64
		sexTypeID         sql.NullInt64
		maritalStatusID   sql.NullInt64
		educationDegreeID sql.NullInt64
		homeAddress       sql.NullString
		phone             sql.NullString
	)
	if err := rows.Scan(
		&patient.PatientID,
		&patient.DocumentNumber,
		&historyNumber,
		&patient.PaternalSurname,
		&maternalSurname,
		&patient.FirstName,
		&secondName,
		&birthDate,
		&homeDistrictID,
		&homeCenterID,
		&sexTypeID,
		&maritalStatusID,
		&educationDegreeID,
		&homeAddress,
		&phone,
	); err != nil {
		return nil, fmt.Errorf("scanning usp_go_ListarPacientePorNroDocyTipo row: %w", err)
	}
	patient.HistoryNumber = historyNumber.String
	patient.MaternalSurname = maternalSurname.String
	patient.SecondName = secondName.String
	if birthDate.Valid {
		d := birthDate.Time
		patient.DateOfBirth = &d
	}
	patient.HomeDistrictID = nullableInt64Ptr(homeDistrictID)
	patient.HomeCenterID = nullableInt64Ptr(homeCenterID)
	patient.SexTypeID = nullableInt64Ptr(sexTypeID)
	patient.MaritalStatusID = nullableInt64Ptr(maritalStatusID)
	patient.EducationDegreeID = nullableInt64Ptr(educationDegreeID)
	if homeAddress.Valid {
		patient.HomeAddress = &homeAddress.String
	}
	if phone.Valid {
		patient.Phone = &phone.String
	}

	return &patient, nil
}

// Search invoca el procedimiento almacenado usp_go_listarpacientes con los
// filtros opcionales del buscador. Los filtros vacíos se ignoran y el SP
// retorna el listado con las columnas TipoDocumento y SegundoNombre entre
// las habituales. Se mapea por nombre de columna para resistir cambios del
// resultado del SP sin modificar el adaptador.
func (r *patientRepository) Search(ctx context.Context, params shared.PatientSearchParams) ([]domain.Patient, error) {
	const procedure = `usp_go_listarpacientes`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("Nrodoc", params.DocumentNumber),
		sql.Named("NroHc", params.HistoryNumber),
		sql.Named("ApellidoPaterno", params.PaternalSurname),
		sql.Named("ApellidoMaterno", params.MaternalSurname),
		sql.Named("Nombres", params.Names),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_listarpacientes: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading patients search: %w", err)
	}

	patients := make([]domain.Patient, 0, len(maps))
	for _, m := range maps {
		p := domain.Patient{
			PatientID:       derefInt64(rowInt64(m, "IdPaciente")),
			DocumentNumber:  derefString(rowString(m, "NroDocumento")),
			HistoryNumber:   derefString(rowString(m, "NroHistoriaClinica")),
			PaternalSurname: derefString(rowString(m, "ApellidoPaterno")),
			MaternalSurname: derefString(rowString(m, "ApellidoMaterno")),
			FirstName:       derefString(rowString(m, "PrimerNombre")),
			SecondName:      derefString(rowString(m, "SegundoNombre")),
			ThirdName:       derefString(rowString(m, "TercerNombre")),
			DocumentType:    rowString(m, "TipoDocumento"),
			DocIdentityID:   rowInt64(m, "IdDocIdentidad"),
		}
		if v := rowTime(m, "FechaNacimiento"); v != nil {
			p.DateOfBirth = v
		}
		if v := rowInt64(m, "IdTipoSexo"); v != nil {
			p.SexTypeID = v
		}
		patients = append(patients, p)
	}

	return patients, nil
}

// GetByID invoca el procedimiento almacenado webPacientesListarIdPaciente
// con el id del paciente y mapea sus columnas (nombres conocidos en runtime)
// al detalle completo. Retorna domain.ErrPatientNotFound si no hay registros.
func (r *patientRepository) GetByID(ctx context.Context, id int64) (*domain.PatientDetail, error) {
	const procedure = `webPacientesListarIdPaciente`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdPaciente", id))
	if err != nil {
		return nil, fmt.Errorf("calling webPacientesListarIdPaciente: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading patient detail: %w", err)
	}
	if len(maps) == 0 {
		return nil, domain.ErrPatientNotFound
	}

	m := maps[0]
	detail := domain.PatientDetail{
		PatientID:             derefInt64(rowInt64(m, "IdPaciente")),
		DocumentNumber:        rowString(m, "NroDocumento"),
		HistoryNumber:         rowString(m, "NroHistoriaClinica"),
		PaternalSurname:       rowString(m, "ApellidoPaterno"),
		MaternalSurname:       rowString(m, "ApellidoMaterno"),
		FirstName:             rowString(m, "PrimerNombre"),
		SecondName:            rowString(m, "SegundoNombre"),
		ThirdName:             rowString(m, "TercerNombre"),
		Phone:                 rowString(m, "Telefono"),
		HomeAddress:           rowString(m, "DireccionDomicilio"),
		AutoGenerated:         rowString(m, "Autogenerado"),
		SexTypeID:             rowInt64(m, "IdTipoSexo"),
		OriginID:              rowInt64(m, "IdProcedencia"),
		EducationDegreeID:     rowInt64(m, "IdGradoInstruccion"),
		MaritalStatusID:       rowInt64(m, "IdEstadoCivil"),
		DocIdentityID:         rowInt64(m, "IdDocIdentidad"),
		OccupationTypeID:      rowInt64(m, "IdTipoOcupacion"),
		BirthCenterID:         rowInt64(m, "IdCentroPobladoNacimiento"),
		HomeCenterID:          rowInt64(m, "IdCentroPobladoDomicilio"),
		FatherName:            rowString(m, "NombrePadre"),
		MotherName:            rowString(m, "NombreMadre"),
		NumberingTypeID:       rowInt64(m, "IdTipoNumeracion"),
		OriginCenterID:        rowInt64(m, "IdCentroPobladoProcedencia"),
		Observation:           rowString(m, "Observacion"),
		HomeCountryID:         rowInt64(m, "IdPaisDomicilio"),
		OriginCountryID:       rowInt64(m, "IdPaisProcedencia"),
		BirthCountryID:        rowInt64(m, "IdPaisNacimiento"),
		OriginDistrictID:      rowInt64(m, "IdDistritoProcedencia"),
		HomeDistrictID:        rowInt64(m, "IdDistritoDomicilio"),
		BirthDistrictID:       rowInt64(m, "IdDistritoNacimiento"),
		FamilyRecord:          rowString(m, "FichaFamiliar"),
		EthnicityID:           rowString(m, "IdEtnia"),
		BloodType:             rowString(m, "GrupoSanguineo"),
		RhFactor:              rowString(m, "FactorRh"),
		LanguageID:            rowInt64(m, "IdIdioma"),
		Email:                 rowString(m, "Email"),
		MotherDocument:        rowString(m, "madreDocumento"),
		MotherPaternalSurname: rowString(m, "madreApellidoPaterno"),
		MotherMaternalSurname: rowString(m, "madreApellidoMaterno"),
		MotherFirstName:       rowString(m, "madrePrimerNombre"),
		MotherSecondName:      rowString(m, "madreSegundoNombre"),
		ChildOrderNumber:      rowInt64(m, "NroOrdenHijo"),
		MotherDocType:         rowString(m, "madreTipoDocumento"),
		Sector:                rowString(m, "Sector"),
		Sectorist:             rowString(m, "Sectorista"),
		BirthRecordNumber:     rowString(m, "NumPartida"),
		EcoGenCaseNumber:      rowString(m, "NumCasoEcoGen"),
		XRayCaseNumber:        rowString(m, "NumCasoRayosX"),
		EcogObsCaseNumber:     rowString(m, "NumCasoEcogObs"),
		StateID:               rowInt64(m, "IdEstado"),
		Cellphone:             rowString(m, "celular"),
		InsuranceTypeID:       rowInt64(m, "IdTipoSeguro"),
		PatientPhoto:          rowString(m, "FotoPaciente"),
		PatientSignature:      rowString(m, "FirmaPaciente"),
	}

	if v := rowTime(m, "FechaNacimiento"); v != nil {
		detail.DateOfBirth = v
	}
	if v := rowBool(m, "UsoWebReniec"); v != nil {
		detail.UsesWebReniec = v
	}
	if v := boolToInt64(rowBool(m, "Dispacidad")); v != nil {
		detail.DisabilityID = v
	}
	if v := boolToInt64(rowBool(m, "Incapacidad")); v != nil {
		detail.IncapacityID = v
	}

	return &detail, nil
}

// Update invoca el procedimiento almacenado usp_go_ModificarPaciente con los
// campos editables. El parámetro IdPaciente se fija desde la ruta; los
// campos no enviados se pasan como NULL.
func (r *patientRepository) Update(ctx context.Context, id int64, update domain.PatientUpdate) error {
	const procedure = `usp_go_ModificarPaciente`

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdPaisNacimiento", update.BirthCountryID),
		sql.Named("ApellidoMaterno", update.MaternalSurname),
		sql.Named("DireccionDomicilio", update.HomeAddress),
		sql.Named("IdPaisProcedencia", update.OriginCountryID),
		sql.Named("IdPaciente", id),
		sql.Named("ApellidoPaterno", update.PaternalSurname),
		sql.Named("PrimerNombre", update.FirstName),
		sql.Named("SegundoNombre", update.SecondName),
		sql.Named("TercerNombre", update.ThirdName),
		sql.Named("FechaNacimiento", update.DateOfBirth),
		sql.Named("NroDocumento", update.DocumentNumber),
		sql.Named("Telefono", update.Phone),
		sql.Named("celular", update.Cellphone),
		sql.Named("Autogenerado", update.AutoGenerated),
		sql.Named("IdTipoSexo", update.SexTypeID),
		sql.Named("IdProcedencia", update.OriginID),
		sql.Named("IdGradoInstruccion", update.EducationDegreeID),
		sql.Named("IdEstadoCivil", update.MaritalStatusID),
		sql.Named("IdDocIdentidad", update.DocIdentityID),
		sql.Named("IdTipoOcupacion", update.OccupationTypeID),
		sql.Named("IdCentroPobladoDomicilio", update.HomeCenterID),
		sql.Named("NombrePadre", update.FatherName),
		sql.Named("NombreMadre", update.MotherName),
		sql.Named("IdPaisDomicilio", update.HomeCountryID),
		sql.Named("NroHistoriaClinica", update.HistoryNumber),
		sql.Named("IdCentroPobladoNacimiento", update.BirthCenterID),
		sql.Named("IdCentroPobladoProcedencia", update.OriginCenterID),
		sql.Named("IdDistritoProcedencia", update.OriginDistrictID),
		sql.Named("IdDistritoDomicilio", update.HomeDistrictID),
		sql.Named("IdDistritoNacimiento", update.BirthDistrictID),
		sql.Named("IdEtnia", update.EthnicityID),
		sql.Named("IdIdioma", update.LanguageID),
		sql.Named("Email", update.Email),
		sql.Named("Dispacidad", update.DisabilityID),
		sql.Named("Incapacidad", update.IncapacityID),
		sql.Named("IdUsuarioAuditoria", update.AuditUserID),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_ModificarPaciente: %w", err)
	}

	return nil
}

// derefInt64 devuelve 0 cuando el puntero es nil.
func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// derefString devuelve "" cuando el puntero es nil.
func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// Create invoca el procedimiento almacenado WebPacienteAgregar_E_H, que
// inserta el paciente y devuelve el IdPaciente generado en el parámetro de
// salida @IdPaciente (SCOPE_IDENTITY). Los parámetros sin default que el
// frontend no envía se pasan como NULL; @NroHistoriaClinica y @IdEstado
// tienen valores por defecto cuando no se informan.
func (r *patientRepository) Create(ctx context.Context, create domain.PatientCreate) (int64, error) {
	const procedure = `WebPacienteAgregar_E_H`

	historyNumber := ""
	if create.HistoryNumber != nil {
		historyNumber = *create.HistoryNumber
	}
	stateID := int64(1)
	if create.StateID != nil {
		stateID = *create.StateID
	}

	var id int64
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("ApellidoPaterno", create.PaternalSurname),
		sql.Named("ApellidoMaterno", create.MaternalSurname),
		sql.Named("PrimerNombre", create.FirstName),
		sql.Named("SegundoNombre", create.SecondName),
		sql.Named("FechaNacimiento", create.DateOfBirth),
		sql.Named("IdDocIdentidad", create.DocIdentityID),
		sql.Named("NroDocumento", create.DocumentNumber),
		sql.Named("Telefono", create.Phone),
		sql.Named("DireccionDomicilio", create.HomeAddress),
		sql.Named("IdTipoSexo", create.SexTypeID),
		sql.Named("IdEstadoCivil", create.MaritalStatusID),
		sql.Named("IdDistritoNacimiento", create.BirthDistrictID),
		sql.Named("IdDistritoDomicilio", create.HomeDistrictID),
		sql.Named("IdTipoSeguro", create.InsuranceTypeID),
		sql.Named("IdProcedencia", create.OriginID),
		sql.Named("IdGradoInstruccion", create.EducationDegreeID),
		sql.Named("IdTipoOcupacion", create.OccupationTypeID),
		sql.Named("NroHistoriaClinica", historyNumber),
		sql.Named("IdEstado", stateID),
		sql.Named("IdPaciente", sql.Out{Dest: &id}),
	)
	if err != nil {
		return 0, fmt.Errorf("calling WebPacienteAgregar_E_H: %w", err)
	}

	return id, nil
}

// Delete verifica primero con PacientesSePuedeEliminar que el paciente no
// tenga registros asociados (Atenciones, farmPreVenta, historias, etc.);
// si Respuesta es 0 retorna domain.ErrPatientCannotBeDeleted. Si el
// paciente no existe, el SP lo marca como eliminable y el DELETE no afecta
// filas; en ese caso se retorna domain.ErrPatientNotFound.
func (r *patientRepository) Delete(ctx context.Context, id int64) error {
	const canDeleteProcedure = `PacientesSePuedeEliminar`
	const deleteProcedure = `PacientesEliminarPorIdPaciente`

	var respuesta int
	if _, err := r.db.ExecContext(ctx, canDeleteProcedure,
		sql.Named("IdPaciente", id),
		sql.Named("Respuesta", sql.Out{Dest: &respuesta}),
	); err != nil {
		return fmt.Errorf("calling PacientesSePuedeEliminar: %w", err)
	}
	if respuesta != 1 {
		return domain.ErrPatientCannotBeDeleted
	}

	result, err := r.db.ExecContext(ctx, deleteProcedure, sql.Named("IdPaciente", id))
	if err != nil {
		return fmt.Errorf("calling PacientesEliminarPorIdPaciente: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("reading delete result: %w", err)
	}
	if affected == 0 {
		return domain.ErrPatientNotFound
	}

	return nil
}

// GetDatosAdicionales invoca usp_go_PacientesDatosAdicionalesIdPaciente para
// obtener los antecedentes médicos del paciente.
func (r *patientRepository) GetDatosAdicionales(ctx context.Context, idPaciente int64) (domain.PacienteDatosAdicionales, error) {
	var datos domain.PacienteDatosAdicionales
	datos.IdPaciente = int(idPaciente)

	// Resolver IdPaciente real si el identificador recibido corresponde a un IdAtencion / IdEpisodio
	resolvedId := idPaciente
	var idPacFromAtencion int64
	errAtencion := r.db.QueryRowContext(ctx, "SELECT IdPaciente FROM dbo.Atenciones WHERE IdAtencion = @p1", sql.Named("p1", idPaciente)).Scan(&idPacFromAtencion)
	if errAtencion == nil && idPacFromAtencion > 0 {
		resolvedId = idPacFromAtencion
	}

	rows, err := r.db.QueryContext(
		ctx,
		"EXEC usp_go_PacientesDatosAdicionalesIdPaciente @IdPaciente = @p1",
		sql.Named("p1", resolvedId),
	)
	if err != nil {
		return datos, fmt.Errorf("calling usp_go_PacientesDatosAdicionalesIdPaciente: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var idPac int
		var ant, antAlerg, antObst, antQuir, antFam, antPat, otrosComorb sql.NullString
		var fNacCalc sql.NullBool
		var hta, ob, dislip, an, hig, tiro, tb, fuma, ca sql.NullInt64

		if err := rows.Scan(
			&idPac,
			&ant,
			&antAlerg,
			&antObst,
			&antQuir,
			&antFam,
			&antPat,
			&fNacCalc,
			&hta,
			&ob,
			&dislip,
			&an,
			&hig,
			&tiro,
			&tb,
			&fuma,
			&ca,
			&otrosComorb,
		); err != nil {
			return datos, fmt.Errorf("scanning usp_go_PacientesDatosAdicionalesIdPaciente: %w", err)
		}

		datos.IdPaciente = int(resolvedId)
		datos.Antecedentes = strings.TrimSpace(ant.String)
		datos.AntecedAlergico = strings.TrimSpace(antAlerg.String)
		datos.AntecedObstetrico = strings.TrimSpace(antObst.String)
		datos.AntecedQuirurgico = strings.TrimSpace(antQuir.String)
		datos.AntecedFamiliar = strings.TrimSpace(antFam.String)
		datos.AntecedPatologico = strings.TrimSpace(antPat.String)
		datos.FNacimientoCalculada = fNacCalc.Valid && fNacCalc.Bool
		datos.HipertensionArterial = int(hta.Int64)
		datos.Obesidad = int(ob.Int64)
		datos.Dislipidemia = int(dislip.Int64)
		datos.Anemia = int(an.Int64)
		datos.HigadoGraso = int(hig.Int64)
		datos.EnfTiroidea = int(tiro.Int64)
		datos.Tuberculosis = int(tb.Int64)
		datos.FumaActualmente = int(fuma.Int64)
		datos.Cancer = int(ca.Int64)
		datos.OtrosComorbilidad = strings.TrimSpace(otrosComorb.String)
	}

	return datos, nil
}

// UpdateDatosAdicionales invoca usp_go_PacientesDatosAdicionalesModificar (y si no existe
// la fila, usp_go_PacientesDatosAdicionalesAgregar) y usp_go_PacientesDatosAdicionalesActualizarComorbilidades.
func (r *patientRepository) UpdateDatosAdicionales(ctx context.Context, idPaciente int64, datos domain.PacienteDatosAdicionales, idUsuario int) error {
	resolvedId := idPaciente
	var idPacFromAtencion int64
	errAtencion := r.db.QueryRowContext(ctx, "SELECT IdPaciente FROM dbo.Atenciones WHERE IdAtencion = @p1", sql.Named("p1", idPaciente)).Scan(&idPacFromAtencion)
	if errAtencion == nil && idPacFromAtencion > 0 {
		resolvedId = idPacFromAtencion
	}

	result, err := r.db.ExecContext(
		ctx,
		"EXEC dbo.usp_go_PacientesDatosAdicionalesModificar @idPaciente = @p1, @antecedentes = @p2, @antecedAlergico = @p3, @antecedObstetrico = @p4, @antecedQuirurgico = @p5, @antecedFamiliar = @p6, @antecedPatologico = @p7, @IdUsuarioAuditoria = @p8",
		sql.Named("p1", resolvedId),
		sql.Named("p2", datos.Antecedentes),
		sql.Named("p3", datos.AntecedAlergico),
		sql.Named("p4", datos.AntecedObstetrico),
		sql.Named("p5", datos.AntecedQuirurgico),
		sql.Named("p6", datos.AntecedFamiliar),
		sql.Named("p7", datos.AntecedPatologico),
		sql.Named("p8", idUsuario),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_PacientesDatosAdicionalesModificar: %w", err)
	}

	affected, err := result.RowsAffected()
	if err == nil && affected == 0 {
		_, errInsert := r.db.ExecContext(
			ctx,
			"EXEC dbo.usp_go_PacientesDatosAdicionalesAgregar @idPaciente = @p1, @antecedentes = @p2, @antecedAlergico = @p3, @antecedObstetrico = @p4, @antecedQuirurgico = @p5, @antecedFamiliar = @p6, @antecedPatologico = @p7, @IdUsuarioAuditoria = @p8",
			sql.Named("p1", resolvedId),
			sql.Named("p2", datos.Antecedentes),
			sql.Named("p3", datos.AntecedAlergico),
			sql.Named("p4", datos.AntecedObstetrico),
			sql.Named("p5", datos.AntecedQuirurgico),
			sql.Named("p6", datos.AntecedFamiliar),
			sql.Named("p7", datos.AntecedPatologico),
			sql.Named("p8", idUsuario),
		)
		if errInsert != nil {
			return fmt.Errorf("calling usp_go_PacientesDatosAdicionalesAgregar: %w", errInsert)
		}
	}

	// Actualizar comorbilidades con el nuevo SP usp_go
	_, _ = r.db.ExecContext(
		ctx,
		"EXEC dbo.usp_go_PacientesDatosAdicionalesActualizarComorbilidades @IdPaciente = @p1, @HipertensionArterial = @p2, @Obesidad = @p3, @Dislipidemia = @p4, @Anemia = @p5, @HigadoGraso = @p6, @EnfTiroidea = @p7, @Tuberculosis = @p8, @FumaActualmente = @p9, @Cancer = @p10, @DescripcionCancer = @p11, @OtrosComorbilidad = @p12",
		sql.Named("p1", resolvedId),
		sql.Named("p2", datos.HipertensionArterial),
		sql.Named("p3", datos.Obesidad),
		sql.Named("p4", datos.Dislipidemia),
		sql.Named("p5", datos.Anemia),
		sql.Named("p6", datos.HigadoGraso),
		sql.Named("p7", datos.EnfTiroidea),
		sql.Named("p8", datos.Tuberculosis),
		sql.Named("p9", datos.FumaActualmente),
		sql.Named("p10", datos.Cancer),
		sql.Named("p11", ""),
		sql.Named("p12", datos.OtrosComorbilidad),
	)

	return nil
}

// nullableInt64Ptr devuelve nil cuando el valor SQL es NULL; de lo contrario
// un puntero al valor.
func nullableInt64Ptr(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	v := n.Int64
	return &v
}
