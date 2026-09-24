package domain

// AuditoriaRegistro representa una fila a insertar en la tabla Auditoria
// mediante el procedimiento dbo.AuditoriaAgregarV.
type AuditoriaRegistro struct {
	IdEmpleado    int
	Accion        string
	IdRegistro    int
	Tabla         string
	IdListItem    int
	NombrePC      string
	Observaciones string
}
