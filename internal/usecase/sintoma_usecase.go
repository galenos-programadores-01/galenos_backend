package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sintomaService struct {
	repo output.SintomaRepository
}

func NewSintomaService(repo output.SintomaRepository) *sintomaService {
	return &sintomaService{repo: repo}
}

func (s *sintomaService) ListarCatalogo(ctx context.Context) ([]domain.SintomaCatalogo, error) {
	return s.repo.ListarCatalogo(ctx)
}

func (s *sintomaService) AgregarCatalogo(ctx context.Context, sistema, sintoma string, idUsuario int) error {
	return s.repo.AgregarCatalogo(ctx, sistema, sintoma, idUsuario)
}

func (s *sintomaService) GuardarEvolucionSintomas(ctx context.Context, idRegAtencion int, sintomas []domain.SintomaSeleccionado, idUsuario int) error {
	return s.repo.GuardarEvolucionSintomas(ctx, idRegAtencion, sintomas, idUsuario)
}

func (s *sintomaService) ObtenerAtencionSintomas(ctx context.Context, idRegAtencion int) ([]domain.SintomaSeleccionado, error) {
	return s.repo.ObtenerAtencionSintomas(ctx, idRegAtencion)
}
