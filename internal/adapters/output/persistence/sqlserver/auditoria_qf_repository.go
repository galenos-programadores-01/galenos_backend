package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type auditoriaQfRepository struct {
	db *sql.DB
}

// NewAuditoriaQfRepository crea el adaptador SQL Server de Auditoría QF.
func NewAuditoriaQfRepository(db *sql.DB) output.AuditoriaQfRepository {
	return &auditoriaQfRepository{db: db}
}

func (r *auditoriaQfRepository) ListarBandeja(ctx context.Context, params domain.BandejaRecetasParams) ([]map[string]any, error) {
	const query = `
EXEC dbo.usp_go_Recetas_Bandeja
    @FechaInicio = @FechaInicio,
    @FechaFin = @FechaFin,
    @IdEstado = @IdEstado,
    @Audita_Dispensacion = @Audita_Dispensacion,
    @IdUsuario = @IdUsuario`

	var idEstado sql.NullInt64
	if params.IdEstado != nil {
		idEstado = sql.NullInt64{Int64: int64(*params.IdEstado), Valid: true}
	}

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("FechaInicio", params.FechaInicio),
		sql.Named("FechaFin", params.FechaFin),
		sql.Named("IdEstado", idEstado),
		sql.Named("Audita_Dispensacion", sql.NullInt64{}),
		sql.Named("IdUsuario", sql.NullInt64{}),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_Recetas_Bandeja: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading recetas bandeja: %w", err)
	}
	return normalizarFilas(maps), nil
}

func (r *auditoriaQfRepository) ListarEstados(ctx context.Context) ([]map[string]any, error) {
	const query = `EXEC dbo.usp_go_RecetaEstados_Listar`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_RecetaEstados_Listar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading receta estados: %w", err)
	}
	return normalizarFilas(maps), nil
}

func (r *auditoriaQfRepository) ListarDetalle(ctx context.Context, params domain.DetalleRecetaParams) ([]map[string]any, error) {
	const query = `
EXEC dbo.usp_go_Receta_Bandeja_Detalle
    @IdReceta = @IdReceta,
    @Audita_Dispensacion = @Audita_Dispensacion,
    @IdUsuario = @IdUsuario`

	var idUsuario sql.NullInt64
	if params.IdUsuario > 0 {
		idUsuario = sql.NullInt64{Int64: int64(params.IdUsuario), Valid: true}
	}

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("IdReceta", params.IdReceta),
		sql.Named("Audita_Dispensacion", sql.NullInt64{}),
		sql.Named("IdUsuario", idUsuario),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_Receta_Bandeja_Detalle: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading receta detalle: %w", err)
	}
	return normalizarFilas(maps), nil
}

func (r *auditoriaQfRepository) ListarHistorico(ctx context.Context, params domain.HistoricoMedicamentoParams) ([]map[string]any, error) {
	const query = `
EXEC dbo.usp_go_Recetas_Bandeja_Detalle_Historico
    @IdPaciente = @IdPaciente,
    @IdProducto = @IdProducto`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("IdPaciente", params.IdPaciente),
		sql.Named("IdProducto", params.IdProducto),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_Recetas_Bandeja_Detalle_Historico: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading receta historico: %w", err)
	}
	return normalizarFilas(maps), nil
}

func normalizarFilas(filas []map[string]any) []map[string]any {
	for _, fila := range filas {
		for clave, valor := range fila {
			switch v := valor.(type) {
			case []byte:
				fila[clave] = string(v)
			case time.Time:
				fila[clave] = v.Format("2006-01-02 15:04:05")
			}
		}
	}
	return filas
}
