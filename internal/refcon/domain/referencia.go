// Package domain contiene las entidades del módulo de referencias y
// contrarreferencias (dashboard y consultas). No depende de adaptadores.
package domain

// ReferenciaPorMes agrupa la cantidad de referencias/contrarreferencias de un
// mes, discriminada por estado, para el dashboard de barras. Equivale a las
// columnas del SP usp_go_ListarReferenciasPorMesyAnio.
type ReferenciaPorMes struct {
	Estado    string `json:"estado"`
	Conteo    int    `json:"conteo"`
	Mes       int    `json:"mes"`
	NombreMes string `json:"nombreMes"`
	Cantidad  int    `json:"cantidad"`
}

// Ups identifica una Unidad Productora de Servicios de salud.
type Ups struct {
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
}

// Establecimiento identifica un establecimiento de salud (IPRESS).
type Establecimiento struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

// DistritoReniec identifica un distrito obtenido por su código RENIEC.
type DistritoReniec struct {
	IdDistrito  int64  `json:"idDistrito"`
	Nombre      string `json:"nombre"`
	IdReniec    int64  `json:"idReniec"`
	IdProvincia int64  `json:"idProvincia"`
}
