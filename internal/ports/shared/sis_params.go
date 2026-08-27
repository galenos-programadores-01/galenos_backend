package shared

// SISAfiliadoParams agrupa los parámetros de la consulta de afiliado al
// servicio web del SIS (operación ConsultarAfiliadoFuaE / BuscarAsegurados).
type SISAfiliadoParams struct {
	// DocumentNumber es el número de documento del paciente a consultar
	// (strNroDocumento).
	DocumentNumber string
	// TipoDocumento es el tipo de documento (1 = DNI, 3 = Carnet de
	// Extranjería). Se envía como strTipoDocumento.
	TipoDocumento int
	// Opcion es el intOpcion del WS; 0 o ausente equivale a 1.
	Opcion int
	// Disa, TipoFormato, NroContrato y Correlativo son los parámetros
	// opcionales del WS de ConsultarAfiliadoFuaE.
	Disa        string
	TipoFormato string
	NroContrato string
	Correlativo string
	// Lote es el número de lote de la afiliación; en la operación
	// BuscarAsegurados se envía como TipoFormato.
	Lote string
	// CodTabla es el código de tabla de la afiliación (BuscarAsegurados).
	CodTabla string
	// Tabla es el código de tabla para la búsqueda de afiliación
	// (parámetro adicional en ConsultarAfiliadoFuaE).
	Tabla string
}
