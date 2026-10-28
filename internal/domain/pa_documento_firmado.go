package domain

// Tipos de PDF firmado de Primera Atención. No incluye procedimientos ni interconsulta.
const (
	TipoDocumentoFirmadoPA           = "PA"
	TipoDocumentoFirmadoExamen       = "EXAMEN"
	TipoDocumentoFirmadoMedicamento  = "MEDICAMENTO"
)

// DocumentoFirmadoPA es una fila inmutable: un PDF ya firmado.
type DocumentoFirmadoPA struct {
	IdPrimeraAtencion int
	NroEdicion        int
	TipoDocumento     string
	IdPuntoCarga      *int
	IdReceta          *int
	NombreArchivo     string
	UbicacionArchivo  string
	EstadoFirma       int
}

// DocumentoFirmadoRegistro es lo que el cierre de firma envía para persistir.
type DocumentoFirmadoRegistro struct {
	TipoDocumento string
	IdPuntoCarga  *int
	IdReceta      *int
	NombreArchivo string
}
