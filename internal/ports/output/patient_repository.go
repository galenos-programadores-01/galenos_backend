package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// PatientRepository es el puerto de salida para lectura de pacientes.
type PatientRepository interface {
	// List retorna la página solicitada de pacientes y el total de
	// registros existentes (para calcular TotalPages en el llamador).
	List(ctx context.Context, page shared.PageRequest) ([]domain.Patient, int, error)

	// GetByDocumentNumber busca un paciente por número de documento
	// invocando el procedimiento almacenado sp_listarPaciente. Retorna
	// domain.ErrPatientNotFound si no existe ningún registro.
	GetByDocumentNumber(ctx context.Context, documentNumber string) (*domain.Patient, error)

	// GetByDocumentNumberAndType busca un paciente por número y tipo de
	// documento invocando el procedimiento almacenado
	// usp_go_ListarPacientePorNroDocyTipo. Retorna domain.ErrPatientNotFound
	// si no existe ningún registro.
	GetByDocumentNumberAndType(ctx context.Context, documentNumber string, documentTypeID int64) (*domain.Patient, error)

	// Search lista pacientes que coinciden con los filtros opcionales
	// invocando el procedimiento almacenado usp_go_listarpacientes.
	Search(ctx context.Context, params shared.PatientSearchParams) ([]domain.Patient, error)

	// GetByID devuelve el detalle completo de un paciente invocando el
	// procedimiento almacenado webPacientesListarIdPaciente. Retorna
	// domain.ErrPatientNotFound si no existe ningún registro.
	GetByID(ctx context.Context, id int64) (*domain.PatientDetail, error)

	// Update modifica un paciente invocando el procedimiento almacenado
	// usp_go_ModificarPaciente. Los campos nil se envían como NULL.
	Update(ctx context.Context, id int64, update domain.PatientUpdate) error

	// Create registra un paciente nuevo invocando el procedimiento almacenado
	// WebPacienteAgregar_E_H y retorna el IdPaciente generado (SCOPE_IDENTITY).
	Create(ctx context.Context, create domain.PatientCreate) (int64, error)

	// Delete elimina un paciente invocando PacientesSePuedeEliminar primero;
	// si el paciente tiene registros asociados retorna
	// domain.ErrPatientCannotBeDeleted y no elimina nada.
	Delete(ctx context.Context, id int64) error

	// GetDatosAdicionales retorna los antecedentes de un paciente invocando
	// el procedimiento almacenado usp_go_PacientesDatosAdicionalesIdPaciente.
	GetDatosAdicionales(ctx context.Context, idPaciente int64) (domain.PacienteDatosAdicionales, error)

	// UpdateDatosAdicionales actualiza o inserta los antecedentes de un paciente
	// invocando PacientesDatosAdicionalesModificar / PacientesDatosAdicionalesAgregar.
	UpdateDatosAdicionales(ctx context.Context, idPaciente int64, datos domain.PacienteDatosAdicionales, idUsuario int) error
}
