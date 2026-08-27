package domain

type Resultado struct {
	IdResultado   int    `json:"idResultado"`
	IdPaciente    int    `json:"idPaciente"`
	IdOrden       int    `json:"idOrden"`
	IdProducto    int    `json:"idProducto"`
	TipoResultado string `json:"tipoResultado"` // "Laboratorio" or "Imagen"
	NombreExamen  string `json:"nombreExamen"`
	FechaExamen   string `json:"fechaExamen"`
	Detalle       string `json:"detalle"`
	Estado        string `json:"estado"`
}

type DetalleResultadoLab struct {
	Grupo            string `json:"grupo"`
	Item             string `json:"item"`
	ValorTexto       string `json:"valorTexto"`
	Unidad           string `json:"unidad"`
	ValorReferencial string `json:"valorReferencial"`
	Metodo           string `json:"metodo"`
}

type DetalleResultadoImagen struct {
	IdOrden      int    `json:"idOrden"`
	IdProducto   int    `json:"idProducto"`
	NombreExamen string `json:"nombreExamen"`
	FechaInforme string `json:"fechaInforme"`
	InformeTexto string `json:"informeTexto"`
}
