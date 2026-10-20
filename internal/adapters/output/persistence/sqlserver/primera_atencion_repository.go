package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type primeraAtencionRepository struct {
	db *sql.DB
}

// NewPrimeraAtencionRepository crea el adaptador SQL Server de Primera Atención.
func NewPrimeraAtencionRepository(db *sql.DB) output.PrimeraAtencionRepository {
	return &primeraAtencionRepository{db: db}
}

// ListarBandeja invoca web_go_GestionPrimeraAtencionFiltrar. Las columnas
// se mapean a map[string]any porque el SP las resuelve en runtime (mismo
// patrón que la bandeja de admisión).
func (r *primeraAtencionRepository) ListarBandeja(ctx context.Context, params shared.PrimeraAtencionParams) ([]map[string]any, error) {
	const procedure = `web_go_GestionPrimeraAtencionFiltrar`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("filtro", params.Filtro),
		sql.Named("IdtipoServicio", params.IdTipoServicio),
		sql.Named("IdEmpleado", params.IdEmpleado),
	)
	if err != nil {
		return nil, fmt.Errorf("calling web_go_GestionPrimeraAtencionFiltrar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading primera atencion bandeja: %w", err)
	}

	return maps, nil
}
