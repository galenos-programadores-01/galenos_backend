// Package usecase implementa los casos de uso del módulo de referencias.
package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
	"github.com/galenos-pro/appointments-api/internal/refcon/ports/output"
)

type refconService struct {
	repo    output.RefConRepository
	minsa   output.MinsaRefConClient
	reports output.RefConReportsClient
}

func NewRefConService(repo output.RefConRepository, minsa output.MinsaRefConClient, reports output.RefConReportsClient) *refconService {
	return &refconService{repo: repo, minsa: minsa, reports: reports}
}

func (s *refconService) ListarReferenciasPorMes(ctx context.Context, mes, anio, estado, opcion int, ups string) ([]domain.ReferenciaPorMes, error) {
	return s.repo.ListarReferenciasPorMes(ctx, mes, anio, estado, opcion, ups)
}

func (s *refconService) ListarUps(ctx context.Context) ([]domain.Ups, error) {
	return s.repo.ListarUps(ctx)
}

func (s *refconService) ListarTodasLasUps(ctx context.Context) ([]domain.Ups, error) {
	return s.repo.ListarTodasLasUps(ctx)
}

func (s *refconService) ListarEstablecimientos(ctx context.Context) ([]domain.Establecimiento, error) {
	return s.repo.ListarEstablecimientos(ctx)
}

func (s *refconService) ListarDistritosPorIdReniec(ctx context.Context, idReniec int) ([]domain.DistritoReniec, error) {
	return s.repo.ListarDistritosPorIdReniec(ctx, idReniec)
}

func (s *refconService) ConsultarReferenciaDetalle(ctx context.Context, req domain.ConsultaMinsaRequest) (*domain.ConsultaMinsaResponse, error) {
	return s.minsa.ConsultarReferenciaDetalle(ctx, req)
}

func (s *refconService) ListarUpssMinsa(ctx context.Context, codigoRenipress string) (*domain.ListadoUpsResponse, error) {
	return s.minsa.ListadoUpss(ctx, codigoRenipress)
}

func (s *refconService) ListarEspecialidadesMinsa(ctx context.Context) (*domain.ListadoEspecialidadesResponse, error) {
	return s.minsa.ListadoEspecialidades(ctx)
}

func (s *refconService) GuardarReferenciaMinsa(ctx context.Context, req domain.SaveReferenciaRequest) (*domain.SaveReferenciaResponse, error) {
	return s.minsa.SaveReferencia(ctx, req)
}

func (s *refconService) GenerarHojaReferencia(ctx context.Context, req domain.GenerarHojaReferenciaRequest) (*domain.GenerarHojaReferenciaResult, error) {
	return s.reports.GenerarHojaReferencia(ctx, req)
}
