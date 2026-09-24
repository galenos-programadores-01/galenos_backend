package domain

import "time"

// SisAfiliado es la respuesta normalizada de una consulta de afiliación al
// servicio web SOAP del SIS (ConsultarAfiliadoSIS). Mapea los campos de
// ResultQueryAsegurado que devuelve la operación.
type SisAfiliado struct {
	IdError            string
	Resultado          string
	TipoDocumento      string
	NroDocumento       string
	ApePaterno         string
	ApeMaterno         string
	Nombres            string
	FecAfiliacion      string
	EESS               string
	DescEESS           string
	EESSUbigeo         string
	DescEESSUbigeo     string
	Regimen            string
	TipoSeguro         string
	DescTipoSeguro     string
	Contrato           string
	FecCaducidad       string
	Estado             string
	Tabla              string
	IdNumReg           string
	Genero             string
	FecNacimiento      string
	IdUbigeo           string
	Direccion          string
	Disa               string
	TipoFormato        string
	NroContrato        string
	Correlativo        string
	IdPlan             string
	IdGrupoPoblacional string
	MsgConfidencial    string
}

// SisAfiliacion representa el registro de afiliación que se persiste
// invocando el SP usp_go_webSisFiliacionesGestionar. Campos opcionales a NULL.
type SisAfiliacion struct {
	IDSiasis                *int64     `json:"idSiasis"`
	Codigo                  *string    `json:"codigo"`
	AfiliacionDisa          *string    `json:"afiliacionDisa"`
	TipoFormato             *string    `json:"afiliacionTipoFormato"`
	NroFormato              *string    `json:"afiliacionNroFormato"`
	NroIntegrante           *string    `json:"afiliacionNroIntegrante"`
	DocumentoTipo           *string    `json:"documentoTipo"`
	CodigoEstablAdscripcion *string    `json:"codigoEstablAdscripcion"`
	AfiliacionFecha         *time.Time `json:"afiliacionFecha"`
	Paterno                 *string    `json:"paterno"`
	Materno                 *string    `json:"materno"`
	PNombre                 *string    `json:"pNombre"`
	ONombres                *string    `json:"oNombres"`
	Genero                  *string    `json:"genero"`
	FNacimiento             *time.Time `json:"fNacimiento"`
	IdDistritoDomicilio     *string    `json:"idDistritoDomicilio"`
	Estado                  *string    `json:"estado"`
	Fbaja                   *string    `json:"fBaja"`
	DocumentoNumero         *string    `json:"documentoNumero"`
	MotivoBaja              *string    `json:"motivoBaja"`
	FbajaOK                 *time.Time `json:"fBajaOk"`
	DescEESS                *string    `json:"descEESS"`
	DescEESSUbigeo          *string    `json:"descEessUbigeo"`
	Regimen                 *string    `json:"regimen"`
	TipoSeguro              *string    `json:"tipoSeguro"`
	DescTipoSeguro          *string    `json:"descTipoSeguro"`
	Contrato                *string    `json:"contrato"`
	IdPlan                  *string    `json:"idPlan"`
	IdGrupoPoblacional      *string    `json:"idGrupoPoblacional"`
	MsgConfidencial         *string    `json:"msgConfidencial"`
	IdUsuarioAuditoria      *int64     `json:"idUsuarioAuditoria"`
}
