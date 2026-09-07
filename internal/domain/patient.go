package domain

import "time"

// Patient es un modelo de lectura del agregado "Paciente", reflejando
// exactamente las columnas expuestas por la consulta a la tabla pacientes
// (NroDocumento, ApellidoPaterno, ApellidoMaterno, PrimerNombre,
// SegundoNombre, TercerNombre) así como por el SP usp_go_listarpacientes
// (IdPaciente, NroHistoriaClinica, FechaNacimiento, TipoDocumento).
type Patient struct {
	PatientID         int64
	DocumentNumber    string
	DocumentType      *string
	DocIdentityID     *int64
	HistoryNumber     string
	PaternalSurname   string
	MaternalSurname   string
	FirstName         string
	SecondName        string
	ThirdName         string
	DateOfBirth       *time.Time
	HomeDistrictID    *int64
	HomeCenterID      *int64
	SexTypeID         *int64
	MaritalStatusID   *int64
	EducationDegreeID *int64
	HomeAddress       *string
	Phone             *string
}

// PacienteDatosAdicionales representa los datos de antecedentes retornados
// por el procedimiento almacenado usp_go_PacientesDatosAdicionalesIdPaciente.
type PacienteDatosAdicionales struct {
	IdPaciente           int    `json:"idPaciente"`
	Antecedentes         string `json:"antecedentes"`
	AntecedAlergico      string `json:"antecedAlergico"`
	AntecedObstetrico    string `json:"antecedObstetrico"`
	AntecedQuirurgico    string `json:"antecedQuirurgico"`
	AntecedFamiliar      string `json:"antecedFamiliar"`
	AntecedPatologico    string `json:"antecedPatologico"`
	FNacimientoCalculada bool   `json:"fNacimientoCalculada"`
	HipertensionArterial int    `json:"hipertensionArterial"`
	Obesidad             int    `json:"obesidad"`
	Dislipidemia         int    `json:"dislipidemia"`
	Anemia               int    `json:"anemia"`
	HigadoGraso          int    `json:"higadoGraso"`
	EnfTiroidea          int    `json:"enfTiroidea"`
	Tuberculosis         int    `json:"tuberculosis"`
	FumaActualmente      int    `json:"fumaActualmente"`
	Cancer               int    `json:"cancer"`
	OtrosComorbilidad    string `json:"otrosComorbilidad"`
}
