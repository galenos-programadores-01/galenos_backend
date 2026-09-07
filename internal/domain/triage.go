package domain

import "time"

// Triage representa el registro de triaje que se persiste invocando el
// procedimiento almacenado webTab_PacienteTriajeAgregar. Los campos
// opcionales son punteros para distinguir "no enviado" (NULL) de un valor
// cero; el SP decide cómo interpretarlos.
type Triage struct {
	IDTriaje                  *int64
	DocIdentityID             *int64
	DocumentNumber            *string
	PaternalSurname           *string
	MaternalSurname           *string
	FirstName                 *string
	SecondName                *string
	ThirdName                 *string
	SexTypeID                 *int64
	DateOfBirth               *time.Time
	Phone                     *string
	HomeDepartmentID          *int64
	HomeProvinceID            *int64
	HomeDistrictID            *int64
	HomeCommunityID           *int64
	HomeAddress               *string
	IsTrafficAccident         *int64
	FundingSourceID           *int64
	Email                     *string
	MaritalStatusID           *int64
	HeartRate                 *int64
	Temperature               *float64
	BloodPressure             *string
	OxygenSaturation          *int64
	RespiratoryRate           *int64
	FiO2                      *int64
	Weight                    *float64
	Height                    *float64
	BMI                       *float64
	EvolutionTimeQuantity     *int64
	EvolutionTimeQuantityUnit *string
	PainScale                 *int64
	GlasgowScale              *int64
	PriorityTypeID            *int64
	ServiceID                 *int64
	Motivo                    *string
	IsPregnant                *int64
	ArrivalStateID            *int64
	FUR                       *time.Time
	EsGestante                *bool
	EdadGestacional           *int64
	FPP                       *time.Time
	NroControlesPrenatales    *int64
	MovimientosFetales        *int64
	Photo                     *string
	EmployeeID                *int64
}

// TriajeConsulta representa el registro de triaje de consulta externa que
// se persiste invocando el procedimiento almacenado
// AtencionesTriajeAgregar. Los signos vitales se envían como texto (el SP
// los declara VARCHAR(10)) y los opcionales como punteros para distinguir
// "no enviado" (NULL) de un valor vacío.
type TriajeConsulta struct {
	IdAtencion        int64
	IdPaciente        int64
	IdEmpleado        int64
	Talla             *string
	Peso              *string
	Temperatura       *string
	Pulso             *string
	FrecRespiratoria  *string
	FrecCardiaca      *string
	FrecCardiacaFetal *string
	PerimCefalico     *string
	Origen            *string
	PerimAbdominal    *string
	SAT02             *string
	FI02              *string
	PresionArterial   *string
	Hemoglobina       *string
	Observacion       *string
	IMC               *string
	Gestante          *string
}
