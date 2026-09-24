package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type AuditoriaRepository struct {
	db *sql.DB
}

func NewAuditoriaRepository(db *sql.DB) *AuditoriaRepository {
	return &AuditoriaRepository{db: db}
}

func (r *AuditoriaRepository) Registrar(ctx context.Context, reg domain.AuditoriaRegistro) error {
	query := `EXEC dbo.AuditoriaAgregarV @IdEmpleado = @p1, @Accion = @p2, @IdRegistro = @p3, @Tabla = @p4, @idListItem = @p5, @nombrePC = @p6, @observaciones = @p7`
	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", reg.IdEmpleado),
		sql.Named("p2", reg.Accion),
		sql.Named("p3", reg.IdRegistro),
		sql.Named("p4", reg.Tabla),
		sql.Named("p5", reg.IdListItem),
		sql.Named("p6", reg.NombrePC),
		sql.Named("p7", reg.Observaciones),
	)
	if err != nil {
		return fmt.Errorf("error registrando auditoría: %w", err)
	}
	return nil
}
