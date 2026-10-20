package domain

import "time"

// DiasMaximosBandejaRecetas es el tope entre fecha inicio y fecha fin.
const DiasMaximosBandejaRecetas = 7

// BandejaRecetasParams son los argumentos de usp_go_Recetas_Bandeja.
// Audita_Dispensacion e IdUsuario se envían NULL.
type BandejaRecetasParams struct {
	FechaInicio string
	FechaFin    string
	IdEstado    *int
}

// DetalleRecetaParams son los argumentos de usp_go_Recetas_Bandeja_Detalle.
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

// Ids de dbo.RecetaEstados que usa usp_go_updateRecetaCabecera.
// El texto que luego muestra la bandeja es Auditado, Observado y Rechazado.
const (
	IdEstadoRecetaAprobar  = 4
	IdEstadoRecetaObservar = 5
	IdEstadoRecetaRechazar = 6
	LargoMotivoReceta      = 200
)

// GuardarEstadoRecetaParams son los argumentos de usp_go_updateRecetaCabecera.
// IdUsuario es el idEmpleado de la sesión. Las fechas vacías se envían NULL.
type GuardarEstadoRecetaParams struct {
	IdReceta          int
	IdEstado          int
	MotivoObservacion string
	FechaObservacion  *time.Time
	MotivoRechazo     string
	FechaRechazo      *time.Time
	IdUsuario         int
}

// ResultadoRecetaError es el texto de negocio que devuelve @Resultado.
type ResultadoRecetaError struct {
	Mensaje string
}

func (e *ResultadoRecetaError) Error() string {
	return e.Mensaje
}
