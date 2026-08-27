package domain

type ListaEsperaQx struct {
	NroHistoriaClinica int     `json:"nroHistoriaClinica"`
	NroDocumento       string  `json:"nroDocumento"`
	Paciente           string  `json:"paciente"`
	FechaNacimiento    string  `json:"fechaNacimiento"`
	FechaOrden         string  `json:"fechaOrden"`
	Observacion        *string `json:"observacion"`
}

type ListaEsperaQxCrear struct {
	IdPaciente       int    `json:"idPaciente"`
	IdMedico         int    `json:"idMedico"`
	IdTipoDocumento  int    `json:"idTipoDocumento"`
	NroDocumento     string `json:"nroDocumento"`
	ApellidoPaterno  string `json:"apellidoPaterno"`
	ApellidoMaterno  string `json:"apellidoMaterno"`
	PrimerNombre     string `json:"primerNombre"`
	SegundoNombre    string `json:"segundoNombre"`
	FechaNacimiento  string `json:"fechaNacimiento"`
	IdSexo           int    `json:"idSexo"`
	Telefono         string `json:"telefono"`
	Direccion        string `json:"direccion"`
	FechaOrden       string `json:"fechaOrden"`
	Diagnostico      string `json:"diagnostico"`
	FechaLaboratorio string `json:"fechaLaboratorio"`
	FechaICCardio    string `json:"fechaICCardio"`
	FechaICNeumo     string `json:"fechaICNeumo"`
	FechaICAnestesio string `json:"fechaICAnestesio"`
	Medico           string `json:"medico"`
	Observacion      string `json:"observacion"`
}
