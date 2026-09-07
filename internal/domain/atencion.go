package domain

// AdmisionDesdeTriaje agrupa los datos necesarios para crear la atención
// (admisión) de un paciente a partir de su triaje, invocando el SP
// WebCrearAtencionDesdeTriaje. Los campos opcionales son punteros para
// distinguir "no enviado" (NULL) de un valor cero.
type AdmisionDesdeTriaje struct {
	IDTriaje            *int64
	IDPacienteTriaje    *int64
	NroDocumento        *string
	IDEmpleado          *int64
	IDMedico            *int64
	NombreAcompanante   *string
	TelefonoAcompanante *string
	DireccionPaciente   *string
	Observacion         *string
}
