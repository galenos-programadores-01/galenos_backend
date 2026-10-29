package domain

import "time"

// AdmisionHospitalizacion agrupa los datos necesarios para admisionar a un
// paciente en hospitalización, invocando el SP usp_go_AdmisionHospitalizacion.
// Los campos opcionales son punteros para distinguir "no enviado" (NULL) de un
// valor cero, ya que el SP los declara con valor por defecto.
type AdmisionHospitalizacion struct {
	IDPaciente                        *int64
	IDDocIdentidad                    *int64
	NroDocumento                      *string
	ApellidoPaterno                   *string
	ApellidoMaterno                   *string
	PrimerNombre                      *string
	SegundoNombre                     *string
	TercerNombre                      *string
	FechaNacimiento                   *time.Time
	IDTipoSexo                        *int64
	EstadoCivil                       *int64
	GradoInstruccion                  *int64
	Ocupacion                         *int64
	IDEtnia                           *int64
	IDIdioma                          *int64
	MadreTipoDoc                      *int64
	MadreNroDoc                       *string
	MadreApePaterno                   *string
	MadreApeMaterno                   *string
	MadrePrimerNombre                 *string
	MadreSegundoNombre                *string
	NombrePadre                       *string
	Telefono                          *string
	Correo                            *string
	IDDepartamento                    *int64
	IDProvincia                       *int64
	IDDistrito                        *int64
	IDComunidadDomicilio              *int64
	Direccion                         *string
	IDNacionalidad                    *int64
	IDOrigen                          *int64
	IDServicio                        *int64
	IDCama                            *int64
	FechaIngreso                      *time.Time
	HoraIngreso                       *string
	IDMedicoIngreso                   *int64
	IDTipoPrioridad                   *int64
	IDEstadollego                     *int64
	IDFuenteFinanciamiento            *int64
	IDEstablecimientoReferenciaOrigen *int64
	NroReferencia                     *string
	IDEmpleado                        *int64
}
