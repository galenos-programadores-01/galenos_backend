package input

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type InterconsultaService interface {
	ObtenerPorId(ctx context.Context, id int) (*domain.Interconsulta, error)
	ListarPorServicio(ctx context.Context, tipoServicio string) ([]domain.Interconsulta, error)
	ListarPorAtencion(ctx context.Context, idAtencion int) ([]domain.Interconsulta, error)
	Crear(ctx context.Context, interconsulta domain.Interconsulta) error
	ActualizarEstado(ctx context.Context, id int, estado string) error
	GuardarFirma(ctx context.Context, firma domain.FirmaInterconsulta) error
	ListarEspecialidades(ctx context.Context) ([]domain.EspecialidadInterconsulta, error)
	ListarMedicosPorEspecialidad(ctx context.Context, IdEspecialidad int) ([]domain.MedicoInterconsulta, error)
}
