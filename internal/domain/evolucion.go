package domain

// PatientListItem representa a un paciente en la bandeja
type PatientListItem struct {
	IdRegAtencion int    `json:"idRegAtencion"`
	IdPaciente    int    `json:"idPaciente"`
	Historia      string `json:"historia"`
	Nombre        string `json:"nombre"`
	Edad          string `json:"edad"`
	Sexo          string `json:"sexo"`
	Ubicacion     string `json:"ubicacion"`
	Cama          string `json:"cama"`
	Estado        string `json:"estado"`
}

// EvolucionFirma representa una evolución guardada en formato JSON / Base64
type EvolucionFirma struct {
	IdRegAtencion      int    `json:"idRegAtencion"`
	IdFirma            int    `json:"idFirma"`
	NombreDocumento    string `json:"nombreDocumento"`
	NombreArchivo      string `json:"nombreArchivo"`
	RutaBase           string `json:"rutaBase"`
	DataB64            string `json:"dataB64"`
	IdEmpleadoRegistra int    `json:"idEmpleadoRegistra"`
	FechaRegistro      string `json:"fechaRegistro"`
	Estado             int    `json:"estado"`
}

type EvolucionBandejaItem struct {
	IdEvolucion            int      `json:"idEvolucion"`
	IdEpisodio             *int     `json:"idEpisodio"`
	NroAtencion            int      `json:"nroAtencion"`
	FechaAtencion          string   `json:"fechaAtencion"`
	IdPaciente             int      `json:"idPaciente"`
	Paciente               string   `json:"paciente"`
	Documento              string   `json:"documento"`
	IdCuentaAtencion       int      `json:"idCuentaAtencion"`
	Motivo                 string   `json:"motivo"`
	EscalaDolor            *int     `json:"escalaDolor"`
	Glasgow                *int     `json:"glasgow"`
	PASistolica            *int     `json:"paSistolica"`
	PADiastolica           *int     `json:"paDiastolica"`
	FrecuenciaCardiaca     *int     `json:"frecuenciaCardiaca"`
	FrecuenciaRespiratoria *int     `json:"frecuenciaRespiratoria"`
	Temperatura            *float64 `json:"temperatura"`
	SaturacionOxigeno      *int     `json:"saturacionOxigeno"`
	IdEstadoClinico        *int     `json:"idEstadoClinico"`
	IdPronostico           *int     `json:"idPronostico"`
	EstadoFirma            *int     `json:"estadoFirma"`
	FechaFirma             *string  `json:"fechaFirma"`
	UsuarioCreacion        *int     `json:"usuarioCreacion"`
	FechaCreacion          *string  `json:"fechaCreacion"`
	EquipoCreacion         *string  `json:"equipoCreacion"`
	EstadoRegistro         *int     `json:"estadoRegistro"`
}

type EvolucionMedicaInsert struct {
	IdAtencion                            int      `json:"idAtencion"`
	IdPaciente                            int      `json:"idPaciente"`
	IdMedico                              int      `json:"idMedico"`
	FechaAtencion                         *string  `json:"fechaAtencion"`
	IdTipoGravedad                        *int     `json:"idTipoGravedad"`
	MotivoConsulta                        *string  `json:"motivoConsulta"`
	TiempoEnfermedad                      *string  `json:"tiempoEnfermedad"`
	Anamnesis                             *string  `json:"anamnesis"`
	EscalaDolor                           *int     `json:"escalaDolor"`
	Glasgow                               *int     `json:"glasgow"`
	PASistolica                           *int     `json:"paSistolica"`
	PADiastolica                          *int     `json:"paDiastolica"`
	FrecuenciaCardiaca                    *int     `json:"frecuenciaCardiaca"`
	FrecuenciaRespiratoria                *int     `json:"frecuenciaRespiratoria"`
	Temperatura                           *float64 `json:"temperatura"`
	SaturacionOxigeno                     *int     `json:"saturacionOxigeno"`
	Peso                                  *float64 `json:"peso"`
	Talla                                 *float64 `json:"talla"`
	IMC                                   *float64 `json:"imc"`
	Glicemia                              *float64 `json:"glicemia"`
	ExamenFisicoGeneral                   *string  `json:"examenFisicoGeneral"`
	ExamenFisicoPiel                      *string  `json:"examenFisicoPiel"`
	ExamenFisicoCabezaCuello              *string  `json:"examenFisicoCabezaCuello"`
	ExamenFisicoToraxPulmon               *string  `json:"examenFisicoToraxPulmon"`
	ExamenFisicoCorazon                   *string  `json:"examenFisicoCorazon"`
	ExamenFisicoAbdomen                   *string  `json:"examenFisicoAbdomen"`
	ExamenFisicoGenitourinario            *string  `json:"examenFisicoGenitourinario"`
	ExamenFisicoExtremidadesOsteomuscular *string  `json:"examenFisicoExtremidadesOsteomuscular"`
	ExamenFisicoNeurologicoMental         *string  `json:"examenFisicoNeurologicoMental"`
	IdEstadoClinico                       *int     `json:"idEstadoClinico"`
	IdPronostico                          *int     `json:"idPronostico"`
	IndicacionDieta                       *string  `json:"indicacionDieta"`
	IndicacionReposo                      *string  `json:"indicacionReposo"`
	IndicacionHidratacion                 *string  `json:"indicacionHidratacion"`
	IndicacionOxigeno                     *string  `json:"indicacionOxigeno"`
	IndicacionRestriccion                 *string  `json:"indicacionRestriccion"`
	Sugerencia                            *string  `json:"sugerencia"`
	UsuarioCreacion                       int      `json:"usuarioCreacion"`
	EquipoCreacion                        *string  `json:"equipoCreacion"`
	EstadoRegistro                        *int     `json:"estadoRegistro"`
	EstadoFirma                           *int     `json:"estadoFirma"`
}
