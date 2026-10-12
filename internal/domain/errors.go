package domain

import "errors"

// Errores de dominio. Los adaptadores de entrada los traducen a códigos
// HTTP; los adaptadores de salida nunca deberían generarlos, solo el dominio.
var (
	ErrInvalidPatient                    = errors.New("patient id is required")
	ErrInvalidDoctor                     = errors.New("doctor id is required")
	ErrInvalidTimeSlot                   = errors.New("appointment end time must be after start time")
	ErrInvalidTimeSlotDuration           = errors.New("appointment duration is out of the allowed range")
	ErrAppointmentInPast                 = errors.New("appointment cannot be scheduled in the past")
	ErrCannotCancelCompleted             = errors.New("a completed appointment cannot be cancelled")
	ErrAlreadyCancelled                  = errors.New("appointment is already cancelled")
	ErrCannotRescheduleClosedAppointment = errors.New("cannot reschedule a cancelled or completed appointment")
	ErrAppointmentNotFound               = errors.New("appointment not found")
	ErrDoctorNotAvailable                = errors.New("doctor is not available in the requested time slot")
	ErrInvalidAppointmentID              = errors.New("appointment id is required")
	ErrInvalidDocumentNumber             = errors.New("document number is required")
	ErrInvalidDocumentType               = errors.New("document type id is required")
	ErrPatientNotFound                   = errors.New("patient not found")
	ErrPatientCannotBeDeleted            = errors.New("patient cannot be deleted because it has associated records")
	ErrInvalidPatientID                  = errors.New("patient id is required")
	ErrInvalidReniecOperation            = errors.New("invalid reniec operation")
	ErrInvalidCredentials                = errors.New("invalid credentials")
	ErrInvalidToken                      = errors.New("invalid token")
	ErrMissingEmployeeID                 = errors.New("the authenticated employee id is required")
)
