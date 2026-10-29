package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// HospitalizacionService es el puerto de entrada para registrar las
// admisiones de hospitalización.
type HospitalizacionService interface {
	// Admisionar registra la admisión de un paciente en hospitalización
	// invocando el procedimiento almacenado usp_go_AdmisionHospitalizacion
	// y retorna el @Resultado del SP.
	Admisionar(ctx context.Context, admision *domain.AdmisionHospitalizacion) (string, error)

	// ObtenerReporteTicket devuelve los datos del ticket de admisión de
	// hospitalización para una cuenta de atención, invocando el
	// procedimiento almacenado usp_go_GenererReporteTicketHospi. Retorna nil
	// si la cuenta no existe.
	ObtenerReporteTicket(ctx context.Context, idCuenta int64) (*domain.ReporteTicketHospitalizacion, error)
}
