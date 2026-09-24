package domain

import "encoding/json"

// SaveReferenciaRequest es la estructura exacta que espera el servicio
// saveReferencia del MINSA para registrar una referencia.
type SaveReferenciaRequest struct {
	Cita                   CitaReferencia          `json:"cita"`
	CPT                    CPT                     `json:"cpt"`
	DatosReferencia        DatosReferenciaSave     `json:"datos_referencia"`
	Diagnostico            []DiagnosticoReferencia `json:"diagnostico"`
	Paciente               PacienteReferencia      `json:"paciente"`
	PersonaAcompana        PersonaAcompana         `json:"persona_acompana"`
	PersonalRegistra       PersonalRegistra        `json:"personal_registra"`
	PersonaEstablecimiento PersonaEstablecimiento  `json:"persona_establecimiento"`
	ResponsableReferencia  ResponsableReferencia   `json:"responsable_referencia"`
	Tratamiento            []TratamientoReferencia `json:"tratamiento"`
}

// CitaReferencia son los signos vitales y la información clínica de la cita.
type CitaReferencia struct {
	FechaVencimientoSis       string `json:"fecha_vencimiento_sis"`
	FrecuenciaCardiaca        string `json:"frecuencia_cardiaca"`
	FrecuenciaRespiratoria    string `json:"frecuencia_respiratoria"`
	IDFinanciador             string `json:"id_financiador"`
	NumAfil                   string `json:"num_afil"`
	Peso                      string `json:"peso"`
	PresionArterialDiastolica string `json:"presion_arterial_diastolica"`
	PresionArterialSistolica  string `json:"presion_arterial_sistolica"`
	ResumeAnamnesis           string `json:"resumeanamnesis"`
	ResumeExfisico            string `json:"resumeexfisico"`
	Talla                     string `json:"talla"`
	Temperatura               string `json:"temperatura"`
}

// CPT son los códigos CPT de la referencia.
type CPT struct {
	CPT1 string `json:"cpt_1"`
	CPT2 string `json:"cpt_2"`
	CPT3 string `json:"cpt_3"`
}

// DatosReferenciaSave son los datos generales del destino de la referencia.
type DatosReferenciaSave struct {
	CodEspecialidad     string           `json:"codEspecialidad"`
	Condicion           string           `json:"condicion"`
	DescCarteraServicio string           `json:"desc_Cartera_servicio"`
	FechaReferencia     string           `json:"fechaReferencia"`
	FgRegistro          string           `json:"fgRegistro"`
	HoraReferencia      string           `json:"horaReferencia"`
	IDCarteraServicio   string           `json:"idCarteraServicio"`
	IDEnvio             string           `json:"idEnvio"`
	IDTipoAtencion      string           `json:"idTipoAtencion"`
	IDTipoTransporte    string           `json:"idTipoTransporte"`
	IDEstabDestino      string           `json:"idestabDestino"`
	IDEstabOrigen       string           `json:"idestabOrigen"`
	IDUpsOrigen         string           `json:"idupsOrigen"`
	IDUpsDestino        string           `json:"idupsdestino"`
	MotivoReferencia    MotivoReferencia `json:"motivo_referencia"`
	NotasObs            string           `json:"notasobs"`
}

// MotivoReferencia es el motivo clínico de la referencia.
type MotivoReferencia struct {
	IDMotivoRef  string `json:"idmotivoref"`
	ObsMotivoRef string `json:"obsmotivoref"`
}

// DiagnosticoReferencia es un diagnóstico de la referencia.
type DiagnosticoReferencia struct {
	Diagnostico     string `json:"diagnostico"`
	NroDiagnostico  string `json:"nro_diagnostico"`
	TipoDiagnostico string `json:"tipo_diagnostico"`
}

// PacienteReferencia son los datos del paciente de la referencia.
type PacienteReferencia struct {
	ApellidoMaternoPaciente string `json:"apelmatpac"`
	ApellidoPaternoPaciente string `json:"apelpatpac"`
	CelularPaciente         string `json:"celularpac"`
	CorreoPaciente          string `json:"correopac"`
	Direccion               string `json:"direccion"`
	FechaNacimientoPaciente string `json:"fechnacpac"`
	IDSexo                  string `json:"idsexo"`
	IDTipoDoc               string `json:"idtipodoc"`
	NombresPaciente         string `json:"nombpac"`
	NumeroHistoria          string `json:"nrohis"`
	NumeroDocumento         string `json:"numdoc"`
	TelefonoPaciente        string `json:"telefonopac"`
	UbigeoActual            string `json:"ubigeoactual"`
	UbigeoReniec            string `json:"ubigeoreniec"`
}

// PersonaAcompana es la persona que acompaña al paciente en la referencia.
type PersonaAcompana struct {
	ApellidoMaterno string `json:"apelmatacomp"`
	ApellidoPaterno string `json:"apelpatacomp"`
	FechaNacimiento string `json:"fechanacacomp"`
	IDColegio       string `json:"idcolegioacomp"`
	IDProfesion     string `json:"idprofesionacomp"`
	IDSexo          string `json:"idsexoacomp"`
	IDTipoDoc       string `json:"idtipodocacmop"`
	Nombres         string `json:"nombperacomp"`
	NumeroDocumento string `json:"numdocacomp"`
}

// PersonalRegistra es el personal de salud que registra la referencia.
type PersonalRegistra struct {
	TipoDocumento   string `json:"tipoDocumento"`
	NumeroDocumento string `json:"nroDocumento"`
	ApellidoPaterno string `json:"apellidoPaterno"`
	ApellidoMaterno string `json:"apellidoMaterno"`
	Nombres         string `json:"nombres"`
	FechaNacimiento string `json:"fechaNacimiento"`
	IDColegio       string `json:"idcolegio"`
	IDProfesion     string `json:"idprofesion"`
	Sexo            string `json:"sexo"`
}

// PersonaEstablecimiento es la persona de contacto del establecimiento.
type PersonaEstablecimiento struct {
	ApellidoMaterno string `json:"apelmata"`
	ApellidoPaterno string `json:"apelpata"`
	FechaNacimiento string `json:"fechanac"`
	IDColegio       string `json:"idcolegio"`
	IDProfesion     string `json:"idprofesion"`
	IDSexo          string `json:"idsexo"`
	IDTipoDoc       string `json:"idtipodoc"`
	Nombres         string `json:"nombper"`
	NumeroDocumento string `json:"numdoc"`
}

// ResponsableReferencia es quien firma la referencia.
type ResponsableReferencia struct {
	ApellidoMaterno string `json:"apelmatrefiere"`
	ApellidoPaterno string `json:"apelpatrefiere"`
	FechaNacimiento string `json:"fechanacrefiere"`
	IDColegio       string `json:"idcolegioref"`
	IDProfesion     string `json:"idprofesionref"`
	IDSexo          string `json:"idsexorefiere"`
	IDTipoDoc       string `json:"idtipodocref"`
	Nombres         string `json:"nombperrefiere"`
	NumeroDocumento string `json:"numdocref"`
}

// TratamientoReferencia es un tratamiento indicado en la referencia.
type TratamientoReferencia struct {
	Cantidad          string `json:"cantidad"`
	CodigoMedicamento string `json:"codigo_medicamento"`
	Frecuencia        string `json:"frecuencia"`
	NroDiagnostico    string `json:"nro_diagnostico"`
	NroTratamiento    string `json:"nro_tratamiento"`
	Periodo           string `json:"periodo"`
	UnidadTiempo      string `json:"unidad_tiempo"`
}

// SaveReferenciaResponse es la respuesta del servicio saveReferencia del MINSA.
// "datos" puede venir como objeto (éxito) o como cadena (error de validación),
// por eso SaveReferenciaDatos implementa UnmarshalJSON tolerante.
type SaveReferenciaResponse struct {
	Codigo  string              `json:"codigo"`
	Mensaje *string             `json:"mensaje"`
	Datos   SaveReferenciaDatos `json:"datos"`
}

// SaveReferenciaDatos es el bloque "datos" de la respuesta del MINSA.
type SaveReferenciaDatos struct {
	FgEstado   string `json:"fg_estado"`
	DescEstado string `json:"desc estado"`
}

func (d *SaveReferenciaDatos) UnmarshalJSON(b []byte) error {
	raw := string(b)
	if raw == "null" || (len(raw) > 0 && raw[0] == '"') {
		*d = SaveReferenciaDatos{}
		return nil
	}
	var v SaveReferenciaDatos
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*d = v
	return nil
}
