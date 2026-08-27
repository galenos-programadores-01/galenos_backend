package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ListaEsperaQxRepository interface {
	Listar(ctx context.Context, fecha string, fechaFin string, paciente string, idEspecialidad int) ([]domain.ListaEsperaQx, error)
	ObtenerPorId(ctx context.Context, id int) (domain.ListaEsperaQxPaciente, error)
	Reporte(ctx context.Context, fecha string, fechaFin string, idEspecialidad int) ([]domain.ListaEsperaQxReporte, error)
	Crear(ctx context.Context, item domain.ListaEsperaQxCrear, idUsuario int) error
	Modificar(ctx context.Context, item domain.ListaEsperaQxModificar) error
}
