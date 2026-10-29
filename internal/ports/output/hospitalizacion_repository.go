package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// HospitalizacionRepository es el puerto de salida para registrar las
// admisiones de hospitalización.
type HospitalizacionRepository interface {
	// Admisionar persiste la admisión de un paciente en hospitalización
	// invocando el procedimiento almacenado usp_go_AdmisionHospitalizacion.
	// Retorna el valor del parámetro de salida @Resultado, que el SP llena
	// con el código/mensaje del resultado de la operación.
	Admisionar(ctx context.Context, admision *domain.AdmisionHospitalizacion) (string, error)

	// ObtenerReporteTicket devuelve los datos del ticket de admisión de
	// hospitalización para una cuenta de atención, invocando el
	// procedimiento almacenado usp_go_GenererReporteTicketHospi. Retorna nil
	// si la cuenta no existe.
	ObtenerReporteTicket(ctx context.Context, idCuenta int64) (*domain.ReporteTicketHospitalizacion, error)
}
