package domain

type Interconsulta struct {
	IdInterconsulta  int    `json:"idInterconsulta"`
	IdAtencionOrigen int    `json:"idAtencionOrigen"`
	IdEspecialidad   int    `json:"IdEspecialidad"`
	IdMedicoDestino  int    `json:"idMedicoDestino"`
	Motivo           string `json:"motivo"`
	FechaSolicitud   string `json:"fechaSolicitud"`
	Estado           string `json:"estado"`
}

type FirmaInterconsulta struct {
	IdInterconsulta int    `json:"idInterconsulta"`
	DataB64         string `json:"dataB64"`
	IdEmpleadoFirma int    `json:"idEmpleadoFirma"`
}

type EspecialidadInterconsulta struct {
	IdEspecialidad   int     `json:"IdEspecialidad"`
	Nombre           *string `json:"nombre"`
	DescripcionLarga *string `json:"descripcionLarga"`
}

type MedicoInterconsulta struct {
	IdMedico       int     `json:"idMedico"`
	IdEmpleado     int     `json:"idEmpleado"`
	CodigoPlanilla *string `json:"codigoPlanilla"`
	Medico         *string `json:"medico"`
}
