package domain

import "time"

// PatientCreateHistoria transporta los datos de un paciente nuevo hacia el SP
// usp_go_PacienteHistoriaClinicaAgregar, que ademas de insertar al paciente
// genera su historia clinica. A diferencia de domain.PatientCreate (SP
// WebPacienteAgregar_E_H), el SP recibe un unico IdDistrito, el pais como
// IdPais, la ocupacion como IdOcupacion y la fuente de financiamiento como
// IdFuenteFinanciamiento.
//
// Los campos opcionales son punteros; un puntero nil se envia como NULL al
// procedimiento. EmployeeID no es opcional: lo toma el adaptador HTTP del
// claim idEmpleado del JWT, nunca del cuerpo de la peticion.
type PatientCreateHistoria struct {
	PaternalSurname   *string
	MaternalSurname   *string
	FirstName         *string
	SecondName        *string
	DateOfBirth       *time.Time
	DocIdentityID     *int64
	DocumentNumber    *string
	Phone             *string
	HomeAddress       *string
	SexTypeID         *int64
	MaritalStatusID   *int64
	HomeDistrictID    *int64
	HomeCountryID     *int64
	EthnicityID       *int64
	LanguageID        *int64
	OccupationID      *int64
	EducationDegreeID *int64
	// InsuranceTypeID corresponde a @IdFuenteFinanciamiento del SP.
	InsuranceTypeID *int64
	// Datos de la madre tutor. Ojo: el nombre del parametro del documento es
	// @MadreNroDocumento en este SP y @MadreDocumento en usp_go_ModificarPaciente.
	MotherDocumentNumber  *string
	MotherPaternalSurname *string
	MotherMaternalSurname *string
	MotherFirstName       *string
	MotherSecondName      *string
	DisabilityID          *int64
	IncapacityID          *int64
	// EmployeeID corresponde a @IdEmpleado, que el SP usa como
	// @IdUsuarioAuditoria al crear la historia clinica.
	EmployeeID int64
}
