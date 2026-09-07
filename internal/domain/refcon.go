package domain

// ReferenciaPorMes agrupa la cantidad de referencias/contrarreferencias de
// un mes, opcionalmente discriminada por estado, para el dashboard de barras.
type ReferenciaPorMes struct {
	Mes       int     `json:"mes"`
	NombreMes string  `json:"nombreMes"`
	Cantidad  int     `json:"cantidad"`
	Estado    *string `json:"estado,omitempty"`
}

// Ups identifica una Unidad Productora de Servicios de salud.
type Ups struct {
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
}
