package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type appointmentRepository struct {
	db *sql.DB
}

// NewAppointmentRepository construye el adaptador que implementa el puerto
// de salida output.AppointmentRepository contra SQL Server.
func NewAppointmentRepository(db *sql.DB) output.AppointmentRepository {
	return &appointmentRepository{db: db}
}

// Create persiste la cita dentro de una transacción: primero verifica
// disponibilidad del médico con un bloqueo de lectura (UPDLOCK, HOLDLOCK)
// para evitar condiciones de carrera con inserciones concurrentes sobre el
// mismo médico/horario, y luego inserta. Todos los parámetros se envían
// como parámetros nombrados (@Nombre), nunca concatenados, para prevenir
// inyección SQL.
func (r *appointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	query := `EXEC dbo.usp_go_Appointments_Crear
		@ID = @p1,
		@PatientID = @p2,
		@DoctorID = @p3,
		@StartsAt = @p4,
		@EndsAt = @p5,
		@Status = @p6,
		@Reason = @p7,
		@CreatedAt = @p8,
		@UpdatedAt = @p9`

	var codigoResultado int
	var mensaje string
	err := r.db.QueryRowContext(ctx, query,
		sql.Named("p1", string(appointment.ID)),
		sql.Named("p2", string(appointment.PatientID)),
		sql.Named("p3", string(appointment.DoctorID)),
		sql.Named("p4", appointment.Slot.StartsAt),
		sql.Named("p5", appointment.Slot.EndsAt),
		sql.Named("p6", string(appointment.Status)),
		sql.Named("p7", appointment.Reason),
		sql.Named("p8", appointment.CreatedAt),
		sql.Named("p9", appointment.UpdatedAt),
	).Scan(&codigoResultado, &mensaje)

	if err != nil {
		return fmt.Errorf("calling usp_go_Appointments_Crear: %w", err)
	}

	if codigoResultado == -1 {
		return domain.ErrDoctorNotAvailable
	}

	return nil
}

// GetByID recupera una cita por id usando el procedimiento almacenado correspondiente.
func (r *appointmentRepository) GetByID(ctx context.Context, id domain.AppointmentID) (*domain.Appointment, error) {
	const query = `EXEC dbo.usp_go_Appointments_ObtenerPorId @ID = @p1`

	row := r.db.QueryRowContext(ctx, query, sql.Named("p1", string(id)))

	var (
		appointment domain.Appointment
		rawID       string
		patientID   string
		doctorID    string
		status      string
	)

	err := row.Scan(
		&rawID,
		&patientID,
		&doctorID,
		&appointment.Slot.StartsAt,
		&appointment.Slot.EndsAt,
		&status,
		&appointment.Reason,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAppointmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning appointment row: %w", err)
	}

	appointment.ID = domain.AppointmentID(rawID)
	appointment.PatientID = domain.PatientID(patientID)
	appointment.DoctorID = domain.DoctorID(doctorID)
	appointment.Status = domain.AppointmentStatus(status)

	return &appointment, nil
}
