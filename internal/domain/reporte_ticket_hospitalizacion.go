package domain

import "time"

// ReporteTicketHospitalizacion agrupa los datos que devuelve el SP
// usp_go_GenererReporteTicketHospi (una fila por cuenta de atención) para
// imprimir el ticket de admisión de hospitalización.
type ReporteTicketHospitalizacion struct {
	Paciente           string     `json:"paciente"`
	NroHistoriaClinica string     `json:"nroHistoriaClinica"`
	NroDocumento       string     `json:"nroDocumento"`
	IDCuentaAtencion   int64      `json:"idCuentaAtencion"`
	Servicio           string     `json:"servicio"`
	Medico             string     `json:"medico"`
	Fecha              *time.Time `json:"fecha"`
	Hora               string     `json:"hora"`
	IAFA               string     `json:"iafa"`
}
