package domain

// DiasMaximosBandejaRecetas es el tope entre fecha inicio y fecha fin.
const DiasMaximosBandejaRecetas = 7

// BandejaRecetasParams son los argumentos de usp_go_Recetas_Bandeja.
// Audita_Dispensacion e IdUsuario se envían NULL.
type BandejaRecetasParams struct {
	FechaInicio string
	FechaFin    string
	IdEstado    *int
}

// DetalleRecetaParams son los argumentos de usp_go_Receta_Bandeja_Detalle.
// Audita_Dispensacion se envía NULL. IdUsuario es el idEmpleado de la sesión.
type DetalleRecetaParams struct {
	IdReceta  int
	IdUsuario int
}

// HistoricoMedicamentoParams son los argumentos de usp_go_Recetas_Bandeja_Detalle_Historico.
type HistoricoMedicamentoParams struct {
	IdPaciente int
	IdProducto int
}
