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
