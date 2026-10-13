package domain

import (
	"encoding/json"
	"strconv"
)

type ColaPacienteCEItem struct {
	IdAtencion         int    `json:"idAtencion"`
	IdCuentaAtencion   int    `json:"idCuentaAtencion"`
	IdPaciente         int    `json:"idPaciente"`
	NroHistoriaClinica int    `json:"nroHistoriaClinica"`
	NroDocumento       string `json:"nroDocumento"`
	PacienteNombre     string `json:"pacienteNombre"`
	FechaIngreso       string `json:"fechaIngreso"`
	HoraIngreso        string `json:"horaIngreso"`
	HoraEgreso         string `json:"horaEgreso"`
	Edad               int    `json:"edad"`
	IdServicioIngreso  int    `json:"idServicioIngreso"`
	ServicioNombre     string `json:"servicioNombre"`
	IdMedicoIngreso    int    `json:"idMedicoIngreso"`
	MedicoNombre       string `json:"medicoNombre"`
	EstadoAtencion     string `json:"estadoAtencion"`
	IdCita             int    `json:"idCita"`
	TipoFinanciador    string `json:"tipoFinanciador"`
	Sexo               string `json:"sexo"`
	Telefono           string `json:"telefono"`
	TieneTriaje        bool   `json:"tieneTriaje"`
	Turno              string `json:"turno"`
}

type AtencionDiagnosticoCEItem struct {
	IdAtencionDiagnostico int    `json:"idAtencionDiagnostico"`
	IdDiagnostico         int    `json:"idDiagnostico"`
	CodigoCIE10           string `json:"codigoCIE10"`
	Descripcion           string `json:"descripcion"`
	IdClasificacionDx     int    `json:"idClasificacionDx"`
	Clasificacion         string `json:"clasificacion"`
	IdSubclasificacionDx  int    `json:"idSubclasificacionDx"`
	LabConfHIS            string `json:"labConfHIS"`
	LabConfHIS1           string `json:"labConfHIS1"`
	LabConfHIS2           string `json:"labConfHIS2"`
	LabConfHIS3           string `json:"labConfHIS3"`
}

type FichaAtencionCEDetalle struct {
	IdAtencion         int                           `json:"idAtencion"`
	IdCuentaAtencion   int                           `json:"idCuentaAtencion"`
	IdPaciente         int                           `json:"idPaciente"`
	PacienteNombre     string                        `json:"pacienteNombre"`
	NroHistoriaClinica int                           `json:"nroHistoriaClinica"`
	NroDocumento       string                        `json:"nroDocumento"`
	Edad               int                           `json:"edad"`
	Sexo               string                        `json:"sexo"`
	Motivo             string                        `json:"motivo"`
	ExamenClinico      string                        `json:"examenClinico"`
	Antecedente        string                        `json:"antecedente"`
	InforCompDiagMed   string                        `json:"inforCompDiagMed"`
	Tratamiento        string                        `json:"tratamiento"`
	PlanTratamiento    string                        `json:"planTratamiento"`
	Observaciones      string                        `json:"observaciones"`
	FechaProxCita      string                        `json:"fechaProxCita"`
	IdDestinoAtencion  int                           `json:"idDestinoAtencion"`
	DestinoNombre      string                        `json:"destinoNombre"`
	HoraIngreso        string                        `json:"horaIngreso"`
	HoraEgreso         string                        `json:"horaEgreso"`
	Estado             int                           `json:"estado"`
	Diagnosticos       []AtencionDiagnosticoCEItem   `json:"diagnosticos"`
	Procedimientos     []AtencionProcedimientoCEItem `json:"procedimientos"`
}

type AtencionProcedimientoCEItem struct {
	IdProcedimiento                int     `json:"idProcedimiento"`
	IdProducto                     int     `json:"idProducto"`
	CodigoCPT                      string  `json:"codigoCPT"`
	NombreProcedimiento            string  `json:"nombreProcedimiento"`
	IdDiagnosticoAsociado          int     `json:"idDiagnosticoAsociado"`
	CodigoCIE10Asociado            string  `json:"codigoCIE10Asociado"`
	DiagnosticoDescripcionAsociado string  `json:"diagnosticoDescripcionAsociado"`
	Cantidad                       int     `json:"cantidad"`
	Observaciones                  string  `json:"observaciones"`
	Precio                         float64 `json:"precio"`
}

type AtencionDiagnosticoGuardarItem struct {
	IdDiagnostico        int    `json:"idDiagnostico"`
	IdClasificacionDx    int    `json:"idClasificacionDx"`
	IdSubclasificacionDx int    `json:"idSubclasificacionDx"`
	Clasificacion        string `json:"clasificacion"`
	LabConfHIS           string `json:"labConfHIS"`
	LabConfHIS1          string `json:"labConfHIS1"`
	LabConfHIS2          string `json:"labConfHIS2"`
	LabConfHIS3          string `json:"labConfHIS3"`
}

type AtencionCEGuardarPayload struct {
	IdAtencion        int                              `json:"idAtencion"`
	Motivo            string                           `json:"motivo"`
	ExamenClinico     string                           `json:"examenClinico"`
	Antecedente       string                           `json:"antecedente"`
	InforCompDiagMed  string                           `json:"inforCompDiagMed"`
	Tratamiento       string                           `json:"tratamiento"`
	PlanTratamiento   string                           `json:"planTratamiento"`
	Observaciones     string                           `json:"observaciones"`
	FechaProxCita     string                           `json:"fechaProxCita"`
	IdDestinoAtencion int                              `json:"idDestinoAtencion"`
	IdUsuario         int                              `json:"idUsuario"`
	CerrarAtencion    bool                             `json:"cerrarAtencion"`
	Diagnosticos      []AtencionDiagnosticoGuardarItem `json:"diagnosticos"`
	Procedimientos    []AtencionProcedimientoCEItem    `json:"procedimientos"`
}

type DestinoAtencionItem struct {
	IdDestinoAtencion int    `json:"idDestinoAtencion"`
	Descripcion       string `json:"descripcion"`
}

type CatalogoLabHISItem struct {
	Codigo      int    `json:"codigo"`
	Valor       string `json:"valor"`
	Descripcion string `json:"descripcion"`
}

type ColaFiltrosCE struct {
	IdMedico   int
	IdServicio int
	Fecha      string
	Filtro     string
}

type HistorialAtencionItem struct {
	IdAtencion         int    `json:"idAtencion"`
	IdCuentaAtencion   int    `json:"idCuentaAtencion"`
	IdPaciente         int    `json:"idPaciente"`
	FechaIngreso       string `json:"fechaIngreso"`
	HoraIngreso        string `json:"horaIngreso"`
	HoraEgreso         string `json:"horaEgreso"`
	ServicioNombre     string `json:"servicioNombre"`
	IdMedico           int    `json:"idMedico"`
	MedicoNombre       string `json:"medicoNombre"`
	IdTipoServicio     int    `json:"idTipoServicio"`
	TipoServicioNombre string `json:"tipoServicioNombre"`
	TiempoProgramado   int    `json:"tiempoProgramado"`
	TiempoAtencion     int    `json:"tiempoAtencion"`
}

type HistorialTriajeItem struct {
	PresionArterial        string  `json:"presionArterial"`
	FrecuenciaCardiaca     int     `json:"frecuenciaCardiaca"`
	FrecuenciaRespiratoria int     `json:"frecuenciaRespiratoria"`
	Temperatura            float64 `json:"temperatura"`
	Saturacion             int     `json:"saturacion"`
	Peso                   float64 `json:"peso"`
	Talla                  float64 `json:"talla"`
	IMC                    float64 `json:"imc"`
	Hemoglobina            string  `json:"hemoglobina"`
	FiO2                   string  `json:"fiO2"`
	FechaRegistro          string  `json:"fechaRegistro"`
	Observacion            string  `json:"observacion"`
}

type PacienteHistorialTriajeItem struct {
	IdTriaje               int     `json:"idTriaje"`
	IdAtencion             int     `json:"idAtencion"`
	FechaRegistro          string  `json:"fechaRegistro"`
	HoraRegistro           string  `json:"horaRegistro"`
	PresionArterial        string  `json:"presionArterial"`
	Pulso                  string  `json:"pulso"`
	FrecuenciaCardiaca     int     `json:"frecuenciaCardiaca"`
	FrecuenciaRespiratoria int     `json:"frecuenciaRespiratoria"`
	Temperatura            float64 `json:"temperatura"`
	Saturacion             int     `json:"saturacion"`
	Peso                   float64 `json:"peso"`
	Talla                  float64 `json:"talla"`
	IMC                    float64 `json:"imc"`
	Hemoglobina            string  `json:"hemoglobina"`
	Observacion            string  `json:"observacion"`
}

type HistorialAtencionDetalle struct {
	IdAtencion         int                         `json:"idAtencion"`
	IdCuentaAtencion   int                         `json:"idCuentaAtencion"`
	IdPaciente         int                         `json:"idPaciente"`
	PacienteNombre     string                      `json:"pacienteNombre"`
	NroHistoriaClinica int                         `json:"nroHistoriaClinica"`
	NroDocumento       string                      `json:"nroDocumento"`
	Edad               int                         `json:"edad"`
	Sexo               string                      `json:"sexo"`
	ServicioNombre     string                      `json:"servicioNombre"`
	MedicoNombre       string                      `json:"medicoNombre"`
	FechaAtencion      string                      `json:"fechaAtencion"`
	HoraInicio         string                      `json:"horaInicio"`
	HoraFin            string                      `json:"horaFin"`
	Triaje             *HistorialTriajeItem        `json:"triaje"`
	Diagnosticos       []AtencionDiagnosticoCEItem `json:"diagnosticos"`
	InforCompDiagMed   string                      `json:"inforCompDiagMed"`
	Motivo             string                      `json:"motivo"`
	Tratamiento        string                      `json:"tratamiento"`
	PlanTratamiento    string                      `json:"planTratamiento"`
	Observaciones      string                      `json:"observaciones"`
	ExamenClinico      string                      `json:"examenClinico"`
	ConsumoEnServicio  string                      `json:"consumoEnServicio"`
	ServicioApoyoAlDx  string                      `json:"servicioApoyoAlDx"`
}

type HistorialReporteOperatorioItem struct {
	IdReporte          int    `json:"idReporte"`
	IdAtencion         int    `json:"idAtencion"`
	CodigoCPT          string `json:"codigoCPT"`
	Fecha              string `json:"fecha"`
	Descripcion        string `json:"descripcion"`
	Cirujano           string `json:"cirujano"`
	Servicio           string `json:"servicio"`
	Estado             string `json:"estado"`
	Cama               string `json:"cama"`
	HoraInicio         string `json:"horaInicio"`
	HoraFin            string `json:"horaFin"`
	TiempoQuirurgico   string `json:"tiempoQuirurgico"`
	DiagnosticoPre     string `json:"diagnosticoPre"`
	DiagnosticoPost    string `json:"diagnosticoPost"`
	OperacionIndicada  string `json:"operacionIndicada"`
	OperacionRealizada string `json:"operacionRealizada"`
	PrimerAyudante     string `json:"primerAyudante"`
	SegundoAyudante    string `json:"segundoAyudante"`
	Anestesiologo      string `json:"anestesiologo"`
	Instrumentista     string `json:"instrumentista"`
	PersonalCirculante string `json:"personalCirculante"`
	TipoAnestesia      string `json:"tipoAnestesia"`
	TipoOperacion      string `json:"tipoOperacion"`
	CuentaGasas        string `json:"cuentaGasas"`
	PiezaOperatoria    string `json:"piezaOperatoria"`
	TipoIntervencion   string `json:"tipoIntervencion"`
	RiesgoQuirurgico   string `json:"riesgoQuirurgico"`
	Hallazgos          string `json:"hallazgos"`
}

type HistorialProcedimientoItem struct {
	IdProcedimiento   int    `json:"idProcedimiento"`
	IdAtencion        int    `json:"idAtencion"`
	CodigoCPT         string `json:"codigoCPT"`
	Fecha             string `json:"fecha"`
	Descripcion       string `json:"descripcion"`
	Medico            string `json:"medico"`
	Servicio          string `json:"servicio"`
	Estado            string `json:"estado"`
	NumeroCuenta      string `json:"numeroCuenta"`
	NumeroHistoria    string `json:"numeroHistoria"`
	FechaHoraAtencion string `json:"fechaHoraAtencion"`
	Especialidad      string `json:"especialidad"`
	TipoSeguro        string `json:"tipoSeguro"`
	Procedencia       string `json:"procedencia"`
	CMP               string `json:"cmp"`
	Hallazgos         string `json:"hallazgos"`
	Diagnosticos      string `json:"diagnosticos"`
	Observacion       string `json:"observacion"`
}

func (p *AtencionCEGuardarPayload) UnmarshalJSON(data []byte) error {
	type Alias AtencionCEGuardarPayload
	aux := struct {
		IdAtencion        any `json:"idAtencion"`
		IdDestinoAtencion any `json:"idDestinoAtencion"`
		IdUsuario         any `json:"idUsuario"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.IdAtencion = parsearEnteroFlexible(aux.IdAtencion)
	p.IdDestinoAtencion = parsearEnteroFlexible(aux.IdDestinoAtencion)
	p.IdUsuario = parsearEnteroFlexible(aux.IdUsuario)
	return nil
}

func (item *AtencionDiagnosticoGuardarItem) UnmarshalJSON(data []byte) error {
	type Alias AtencionDiagnosticoGuardarItem
	aux := struct {
		IdDiagnostico        any `json:"idDiagnostico"`
		IdClasificacionDx    any `json:"idClasificacionDx"`
		IdSubclasificacionDx any `json:"idSubclasificacionDx"`
		*Alias
	}{
		Alias: (*Alias)(item),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	item.IdDiagnostico = parsearEnteroFlexible(aux.IdDiagnostico)
	item.IdClasificacionDx = parsearEnteroFlexible(aux.IdClasificacionDx)
	item.IdSubclasificacionDx = parsearEnteroFlexible(aux.IdSubclasificacionDx)
	return nil
}

func parsearEnteroFlexible(valor any) int {
	if valor == nil {
		return 0
	}
	switch v := valor.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 0
}

type ServicioMedicoItem struct {
	IdServicio              int    `json:"idServicio"`
	DescripcionHospitalaria string `json:"descripcionHospitalaria"`
	NombreServicio          string `json:"nombreServicio"`
}
