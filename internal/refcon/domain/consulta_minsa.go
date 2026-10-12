package domain

import (
	"encoding/json"
	"strconv"
)

// ConsultaMinsaRequest es la petición que se envía al servicio REST de
// referencias del MINSA (consultaReferenciaDetalle).
type ConsultaMinsaRequest struct {
	EstablecimientoDestino string `json:"establecimientoDestino"`
	Limite                 string `json:"limite"`
	NumeroDocumento        string `json:"numerodocumento"`
	Pagina                 string `json:"pagina"`
	TipoDocumento          string `json:"tipodocumento"`
}

// ConsultaMinsaResponse es la respuesta del servicio MINSA (codigo/mensaje
// que envuelven los datos paginados de la consulta).
type ConsultaMinsaResponse struct {
	Codigo  string              `json:"codigo"`
	Mensaje string              `json:"mensaje"`
	Datos   *ConsultaMinsaDatos `json:"datos,omitempty"`
}

// ConsultaMinsaDatos es el bloque "datos" de la respuesta (paginado). MINSA lo
// devuelve como objeto cuando hay resultados y como cadena (p. ej. "") cuando
// no los hay, por eso implementa UnmarshalJSON.
type ConsultaMinsaDatos struct {
	Paginas   flexibleInt       `json:"paginas"`
	PorPagina string            `json:"porPagina"`
	Total     flexibleInt       `json:"total"`
	Datos     []ReferenciaMinsa `json:"datos,omitempty"`
}

func (d *ConsultaMinsaDatos) UnmarshalJSON(b []byte) error {
	raw := string(b)
	if raw == "null" || (len(raw) > 0 && raw[0] == '"') {
		return nil
	}
	type alias ConsultaMinsaDatos
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*d = ConsultaMinsaDatos(a)
	return nil
}

// ReferenciaMinsa es una referencia devuelta por MINSA en el detalle.
type ReferenciaMinsa struct {
	Rownum           *string              `json:"rownum,omitempty"`
	Paciente         PacienteMinsa        `json:"paciente,omitempty"`
	DatosTutor       DatosTutorMinsa      `json:"datos_tutor,omitempty"`
	DatosReferencia  DatosReferenciaMinsa `json:"datos_referencia,omitempty"`
	Diagnosticos     []DiagnosticoMinsa   `json:"diagnosticos,omitempty"`
	CptProcedimiento *string              `json:"cpt_procedimiento,omitempty"`
	CptLaboratorio   *string              `json:"cpt_laboratorio,omitempty"`
	CptImagenes      *string              `json:"cpt_imagenes,omitempty"`
	Tratamiento      *string              `json:"tratamiento,omitempty"`
}

// PacienteMinsa datos del paciente de la referencia.
type PacienteMinsa struct {
	TipoDocumento       *string `json:"tipo_documento,omitempty"`
	NumeroDocumento     *string `json:"numero_documento,omitempty"`
	Nombres             *string `json:"nombres,omitempty"`
	PrimerApellido      *string `json:"primer_apellido,omitempty"`
	SegundoApellido     *string `json:"segundo_apellido,omitempty"`
	Sexo                *string `json:"sexo,omitempty"`
	Direccion           *string `json:"direccion,omitempty"`
	Ubigeo1             *string `json:"ubigeo1,omitempty"`
	Ubigeo2             *string `json:"ubigeo2,omitempty"`
	FechaNacimiento     *string `json:"fecha_nacimiento,omitempty"`
	NumeroSeguro        *string `json:"numero_seguro,omitempty"`
	FechaVencimientoSis *string `json:"fecha_vencimiento_sis,omitempty"`
	Celular             *string `json:"celular,omitempty"`
}

// DatosTutorMinsa datos del tutor del paciente (puede venir vacío).
type DatosTutorMinsa struct {
	TipoDocumento   *string `json:"tipo_documento,omitempty"`
	NumeroDocumento *string `json:"numero_documento,omitempty"`
	Nombres         *string `json:"nombres,omitempty"`
	PrimerApellido  *string `json:"primer_apellido,omitempty"`
	SegundoApellido *string `json:"segundo_apellido,omitempty"`
	Celular         *string `json:"celular,omitempty"`
	Correo          *string `json:"correo,omitempty"`
}

// DatosReferenciaMinsa datos propios de la referencia/contrarreferencia.
type DatosReferenciaMinsa struct {
	CodigoEspecialidad          *string `json:"codigo_especialidad,omitempty"`
	CodigoEstado                *string `json:"codigoEstado,omitempty"`
	Estado                      *string `json:"estado,omitempty"`
	Condicion                   *string `json:"condicion,omitempty"`
	FechaReferencia             *string `json:"fecha_referencia,omitempty"`
	HoraReferencia              *string `json:"hora_referencia,omitempty"`
	TipoTransporte              *string `json:"tipo_transporte,omitempty"`
	ServicioOrigen              *string `json:"servicio_origen,omitempty"`
	CodigoEstablecimientoOrigen *string `json:"codigo_establecimiento_origen,omitempty"`
	ServicioDestino             *string `json:"servicio_destino,omitempty"`
	NumeroReferencia            *string `json:"numero_referencia,omitempty"`
	IdReferencia                *string `json:"id_referencia,omitempty"`
	FechaEnvio                  *string `json:"fecha_envio,omitempty"`
	ResumeAnamnesis             *string `json:"resume_anamnesis,omitempty"`
	ResumeExfisico              *string `json:"resume_exfisico,omitempty"`
	MotivoReferencia            *string `json:"motivo_referencia,omitempty"`
	TipoFinanciador             *string `json:"tipo_financiador,omitempty"`
	FechaAceptacion             *string `json:"fecha_aceptacion,omitempty"`
}

// DiagnosticoMinsa un diagnóstico CIE-X de la referencia.
type DiagnosticoMinsa struct {
	Id              *string `json:"id,omitempty"`
	CodigoCiex      *string `json:"codigo_ciex,omitempty"`
	TipoDiagnostico *string `json:"tipo_diagnostico,omitempty"`
}

// flexibleInt decodifica enteros que MINSA puede devolver como número o como
// cadena, tolerando esa inconsistencia típica del servicio.
type flexibleInt int

func (f *flexibleInt) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == "" {
		return nil
	}
	if i, err := strconv.Atoi(s); err == nil {
		*f = flexibleInt(i)
		return nil
	}
	var fallback int
	if err := json.Unmarshal(b, &fallback); err != nil {
		return err
	}
	*f = flexibleInt(fallback)
	return nil
}
