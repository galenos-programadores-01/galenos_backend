package domain

type ListaEsperaQx struct {
	IdListaEspera      int     `json:"idListaEspera"`
	NroHistoriaClinica int     `json:"nroHistoriaClinica"`
	NroDocumento       string  `json:"nroDocumento"`
	Paciente           string  `json:"paciente"`
	Edad               int     `json:"edad"`
	FechaOrden         string  `json:"fechaOrden"`
	Especialidad       string  `json:"especialidad"`
	Observacion        *string `json:"observacion"`
	DiasTranscurridos  int     `json:"diasTranscurridos"`
}

type ListaEsperaQxPaciente struct {
	NroDocumento     string  `json:"nroDocumento"`
	IdDocIdentidad   *int    `json:"idDocIdentidad"`
	ApellidoPaterno  string  `json:"apellidoPaterno"`
	ApellidoMaterno  string  `json:"apellidoMaterno"`
	PrimerNombre     string  `json:"primerNombre"`
	Direccion        *string `json:"direccion"`
	Telefono         *string `json:"telefono"`
	IdTipoSexo       *int    `json:"idTipoSexo"`
	FechaNacimiento  string  `json:"fechaNacimiento"`
	FechaOrden       string  `json:"fechaOrden"`
	Diagnostico      *string `json:"diagnostico"`
	IdDiagnostico    *int    `json:"idDiagnostico"`
	IdEspecialidad   *int    `json:"idEspecialidad"`
	FechaLab         string  `json:"fechaLab"`
	FechaICCardio    string  `json:"fechaICCardio"`
	FechaICNeumo     string  `json:"fechaICNeumo"`
	FechaICAnestesio string  `json:"fechaICAnestesio"`
	IdMedico         *int    `json:"idMedico"`
	Medico           *string `json:"medico"`
	Observacion      *string `json:"observacion"`
}

type ListaEsperaQxModificar struct {
	Id               int    `json:"id"`
	FechaOrden       string `json:"fechaOrden"`
	Diagnostico      int    `json:"diagnostico"`
	IdEspecialidad   int    `json:"idEspecialidad"`
	FechaLab         string `json:"fechaLaboratorio"`
	FechaICCardio    string `json:"fechaICCardio"`
	FechaICNeumo     string `json:"fechaICNeumo"`
	FechaICAnestesio string `json:"fechaICAnestesio"`
	Observacion      string `json:"observacion"`
}

type ListaEsperaQxCrear struct {
	IdPaciente       int    `json:"idPaciente"`
	IdMedico         int    `json:"idMedico"`
	FechaOrden       string `json:"fechaOrden"`
	Diagnostico      int    `json:"diagnostico"`
	IdEspecialidad   int    `json:"idEspecialidad"`
	FechaLab         string `json:"fechaLaboratorio"`
	FechaICCardio    string `json:"fechaICCardio"`
	FechaICNeumo     string `json:"fechaICNeumo"`
	FechaICAnestesio string `json:"fechaICAnestesio"`
	Observacion      string `json:"observacion"`
}

type ListaEsperaQxReporte struct {
	Id               int     `json:"id"`
	NroHistoriaClinica int   `json:"nroHistoriaClinica"`
	NroDocumento      string `json:"nroDocumento"`
	Paciente          string `json:"paciente"`
	Edad              int    `json:"edad"`
	Telefono          string `json:"telefono"`
	FechaOrden        string `json:"fechaOrden"`
	Especialidad      string `json:"especialidad"`
	Diagnostico       string `json:"diagnostico"`
	FechaLab          string `json:"fechaLab"`
	FechaICCardio     string `json:"fechaICCardio"`
	FechaICNeumo      string `json:"fechaICNeumo"`
	FechaICAnestesio  string `json:"fechaICAnestesio"`
	Medico            string `json:"medico"`
	Observacion       string `json:"observacion"`
	DiasTranscurridos int    `json:"diasTranscurridos"`
}
