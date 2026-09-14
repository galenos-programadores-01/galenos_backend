package domain

// GenerarHojaReferenciaRequest son los datos que el portal REFCON necesita
// para generar la hoja de referencia institucional en PDF.
type GenerarHojaReferenciaRequest struct {
	IDEstablecimiento int    `json:"idestablecimiento"`
	IDReferencia      int    `json:"idreferencia"`
	EstadoReferencia  string `json:"estadoreferencia"`
}

// GenerarHojaReferenciaResult es la hoja de referencia oficial descargada
// desde el portal REFCON, codificada en base64.
type GenerarHojaReferenciaResult struct {
	PDFBase64        string `json:"archivoB64"`
	URLFile          string `json:"urlFile"`
	EstadoReferencia string `json:"estadoreferencia"`
	NombreReporte    string `json:"nombreReporte,omitempty"`
}
