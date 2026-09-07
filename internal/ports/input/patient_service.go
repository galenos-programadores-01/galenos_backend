package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// PatientService es el puerto de entrada para consulta de pacientes.
type PatientService interface {
	List(ctx context.Context, page shared.PageRequest) (shared.PageResponse[domain.Patient], error)

	// GetByDocumentNumber busca un paciente por su número de documento.
	GetByDocumentNumber(ctx context.Context, documentNumber string) (domain.Patient, error)

	// GetByDocumentNumberAndType busca un paciente por número y tipo de
	// documento invocando el SP usp_go_ListarPacientePorNroDocyTipo.
	GetByDocumentNumberAndType(ctx context.Context, documentNumber string, documentTypeID int64) (domain.Patient, error)

	// Search lista pacientes combinando filtros opcionales (documento, HC,
	// apellidos, nombres) invocando el SP usp_go_listarpacientes.
	Search(ctx context.Context, params shared.PatientSearchParams) ([]domain.Patient, error)

	// GetByID devuelve el detalle completo de un paciente por su id
	// invocando el SP webPacientesListarIdPaciente.
	GetByID(ctx context.Context, id int64) (domain.PatientDetail, error)

	// Update modifica los datos editables de un paciente invocando el SP
	// usp_go_ModificarPaciente y devuelve el detalle actualizado.
	Update(ctx context.Context, id int64, update domain.PatientUpdate) (domain.PatientDetail, error)

	// Create registra un paciente nuevo invocando el SP WebPacienteAgregar_E_H
	// y devuelve el detalle del paciente creado.
	Create(ctx context.Context, create domain.PatientCreate) (domain.PatientDetail, error)

	// Delete elimina un paciente si no tiene registros asociados (SP
	// PacientesSePuedeEliminar + PacientesEliminarPorIdPaciente).
	Delete(ctx context.Context, id int64) error

	// GetDatosAdicionales retorna los antecedentes de un paciente.
	GetDatosAdicionales(ctx context.Context, idPaciente int64) (domain.PacienteDatosAdicionales, error)

	// UpdateDatosAdicionales actualiza o inserta los antecedentes de un paciente.
	UpdateDatosAdicionales(ctx context.Context, idPaciente int64, datos domain.PacienteDatosAdicionales, idUsuario int) (domain.PacienteDatosAdicionales, error)
}
