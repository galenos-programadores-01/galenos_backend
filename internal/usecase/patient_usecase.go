package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type patientUseCase struct {
	repo output.PatientRepository
}

// NewPatientUseCase construye el caso de uso de consulta de pacientes.
func NewPatientUseCase(repo output.PatientRepository) input.PatientService {
	return &patientUseCase{repo: repo}
}

func (uc *patientUseCase) List(ctx context.Context, page shared.PageRequest) (shared.PageResponse[domain.Patient], error) {
	patients, total, err := uc.repo.List(ctx, page)
	if err != nil {
		return shared.PageResponse[domain.Patient]{}, fmt.Errorf("listing patients: %w", err)
	}
	return shared.NewPageResponse(patients, page, total), nil
}

func (uc *patientUseCase) GetByDocumentNumber(ctx context.Context, documentNumber string) (domain.Patient, error) {
	documentNumber = strings.TrimSpace(documentNumber)
	if documentNumber == "" {
		return domain.Patient{}, domain.ErrInvalidDocumentNumber
	}

	patient, err := uc.repo.GetByDocumentNumber(ctx, documentNumber)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			return domain.Patient{}, err
		}
		return domain.Patient{}, fmt.Errorf("getting patient by document number: %w", err)
	}

	return *patient, nil
}

func (uc *patientUseCase) GetByDocumentNumberAndType(ctx context.Context, documentNumber string, documentTypeID int64) (domain.Patient, error) {
	documentNumber = strings.TrimSpace(documentNumber)
	if documentNumber == "" {
		return domain.Patient{}, domain.ErrInvalidDocumentNumber
	}
	if documentTypeID <= 0 {
		return domain.Patient{}, domain.ErrInvalidDocumentType
	}

	patient, err := uc.repo.GetByDocumentNumberAndType(ctx, documentNumber, documentTypeID)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			return domain.Patient{}, err
		}
		return domain.Patient{}, fmt.Errorf("getting patient by document number and type: %w", err)
	}

	return *patient, nil
}

func (uc *patientUseCase) Search(ctx context.Context, params shared.PatientSearchParams) ([]domain.Patient, error) {
	patients, err := uc.repo.Search(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("searching patients: %w", err)
	}
	return patients, nil
}

func (uc *patientUseCase) GetByID(ctx context.Context, id int64) (domain.PatientDetail, error) {
	if id <= 0 {
		return domain.PatientDetail{}, domain.ErrInvalidPatientID
	}

	patient, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			return domain.PatientDetail{}, err
		}
		return domain.PatientDetail{}, fmt.Errorf("getting patient by id: %w", err)
	}

	return *patient, nil
}

func (uc *patientUseCase) Update(ctx context.Context, id int64, update domain.PatientUpdate) (domain.PatientDetail, error) {
	if id <= 0 {
		return domain.PatientDetail{}, domain.ErrInvalidPatientID
	}

	if err := uc.repo.Update(ctx, id, update); err != nil {
		return domain.PatientDetail{}, fmt.Errorf("updating patient: %w", err)
	}

	patient, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return domain.PatientDetail{}, fmt.Errorf("getting updated patient: %w", err)
	}

	return *patient, nil
}

func (uc *patientUseCase) Create(ctx context.Context, create domain.PatientCreate) (domain.PatientDetail, error) {
	id, err := uc.repo.Create(ctx, create)
	if err != nil {
		return domain.PatientDetail{}, fmt.Errorf("creating patient: %w", err)
	}

	patient, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return domain.PatientDetail{}, fmt.Errorf("getting created patient: %w", err)
	}

	return *patient, nil
}

func (uc *patientUseCase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidPatientID
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrPatientCannotBeDeleted || err == domain.ErrPatientNotFound {
			return err
		}
		return fmt.Errorf("deleting patient: %w", err)
	}

	return nil
}

func (uc *patientUseCase) GetDatosAdicionales(ctx context.Context, idPaciente int64) (domain.PacienteDatosAdicionales, error) {
	if idPaciente <= 0 {
		return domain.PacienteDatosAdicionales{}, domain.ErrInvalidPatientID
	}

	datos, err := uc.repo.GetDatosAdicionales(ctx, idPaciente)
	if err != nil {
		return domain.PacienteDatosAdicionales{}, fmt.Errorf("getting patient additional data: %w", err)
	}

	return datos, nil
}

func (uc *patientUseCase) UpdateDatosAdicionales(ctx context.Context, idPaciente int64, datos domain.PacienteDatosAdicionales, idUsuario int) (domain.PacienteDatosAdicionales, error) {
	if idPaciente <= 0 {
		return domain.PacienteDatosAdicionales{}, domain.ErrInvalidPatientID
	}

	if err := uc.repo.UpdateDatosAdicionales(ctx, idPaciente, datos, idUsuario); err != nil {
		return domain.PacienteDatosAdicionales{}, fmt.Errorf("updating patient additional data: %w", err)
	}

	return uc.repo.GetDatosAdicionales(ctx, idPaciente)
}
