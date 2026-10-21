package domain

type FormatoImpresionReceta struct {
	IdentificadorReceta           int                            `json:"idReceta"`
	NumeroReceta                  string                         `json:"numeroReceta"`
	IdentificadorCuentaAtencion   int                            `json:"idCuentaAtencion"`
	IdentificadorAtencion         int                            `json:"idAtencion"`
	PacienteNombreCompleto        string                         `json:"pacienteNombre"`
	DireccionDomicilio            string                         `json:"direccion"`
	NumeroFormatoUnicoAtencion    string                         `json:"numeroFua"`
	TipoAtencion                  string                         `json:"tipoAtencion"`
	TipoFinanciamiento            string                         `json:"tipoFinanciamiento"`
	NumeroTelefono                string                         `json:"telefono"`
	EspecialidadMedica            string                         `json:"especialidadMedico"`
	Consultorio                   string                         `json:"consultorio"`
	NumeroHistoriaClinica         string                         `json:"nroHistoriaClinica"`
	EdadPaciente                  string                         `json:"edad"`
	CuentaAtencion                string                         `json:"cuentaAtencion"`
	FechaAtencion                 string                         `json:"fechaAtencion"`
	DocumentoNacionalIdentidad    string                         `json:"documentoIdentidad"`
	PesoKilogramos                string                         `json:"peso"`
	TallaCentimetros              string                         `json:"talla"`
	MedicoNombreCompleto          string                         `json:"medicoNombre"`
	MedicoColegiatura             string                         `json:"medicoColegiatura"`
	MedicoRegistroEspecialidad    string                         `json:"medicoRNE"`
	UsuarioResponsable            string                         `json:"usuario"`
	Diagnosticos                  []DiagnosticoImpresionReceta   `json:"diagnosticos"`
	DiagnosticoTextoResumido      string                         `json:"diagnosticoTexto"`
	Medicamentos                  []MedicamentoImpresionReceta   `json:"medicamentos"`
}

type DiagnosticoImpresionReceta struct {
	CodigoCIE10                   string `json:"codigoCIE10"`
	Descripcion                   string `json:"descripcion"`
	TipoDiagnostico               string `json:"tipo"`
}

type MedicamentoImpresionReceta struct {
	NumeroItem                    int    `json:"item"`
	CodigoProducto                string `json:"codigo"`
	NombreProductoCompleto        string `json:"nombre"`
	FormaFarmaceutica             string `json:"formaFarmaceutica"`
	DosisPosologiaCompleta        string `json:"dosis"`
	CantidadTotal                 int    `json:"cantidad"`
}
