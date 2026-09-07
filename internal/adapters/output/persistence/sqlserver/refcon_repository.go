package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type RefConRepository struct {
	db *sql.DB
}

func NewRefConRepository(db *sql.DB) *RefConRepository {
	return &RefConRepository{db: db}
}

// ListarReferenciasPorMes ejecuta usp_go_ListarReferenciasPorMesyAnio.
// El SP devuelve 3 columnas para opcion=1 (mes, nombre del mes, cantidad) y
// 5 columnas para opcion 2, 3 y 4 (estado, conteo, mes, nombre mes, cantidad);
// por eso el mapeo se hace por posición según la cantidad de columnas.
func (r *RefConRepository) ListarReferenciasPorMes(ctx context.Context, mes, anio, estado, opcion int, ups string) ([]domain.ReferenciaPorMes, error) {
	query := "EXEC usp_go_ListarReferenciasPorMesyAnio @mes = @p1, @año = @p2, @estado = @p3, @ups = @p4, @opcion = @p5"
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", mes),
		sql.Named("p2", anio),
		sql.Named("p3", estado),
		sql.Named("p4", ups),
		sql.Named("p5", opcion),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarReferenciasPorMesyAnio: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("reading columns: %w", err)
	}

	var lista []domain.ReferenciaPorMes
	for rows.Next() {
		values, err := rowValues(rows, len(columns))
		if err != nil {
			return nil, err
		}

		switch len(columns) {
		case 3:
			// columnas: month(FechaReferencia), Mes, Cantidad
			lista = append(lista, domain.ReferenciaPorMes{
				Mes:       valueToInt(values[0]),
				NombreMes: valueToString(values[1]),
				Cantidad:  valueToInt(values[2]),
			})
		case 5:
			// columnas: descripcion, count(...), month(FechaReferencia), Mes, Cantidad
			estadoDesc := valueToString(values[0])
			item := domain.ReferenciaPorMes{
				Mes:       valueToInt(values[2]),
				NombreMes: valueToString(values[3]),
				Cantidad:  valueToInt(values[4]),
			}
			if estadoDesc != "" {
				item.Estado = &estadoDesc
			}
			lista = append(lista, item)
		}
	}
	return lista, rows.Err()
}

func (r *RefConRepository) ListarUps(ctx context.Context) ([]domain.Ups, error) {
	query := "EXEC usp_go_ListarUPs"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarUPs: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading ups: %w", err)
	}

	lista := make([]domain.Ups, 0, len(maps))
	for _, m := range maps {
		var item domain.Ups
		if v := rowString(m, "Codigo", "codigo"); v != nil {
			item.Codigo = *v
		}
		if v := rowString(m, "Descripcion", "descripcion"); v != nil {
			item.Descripcion = *v
		}
		lista = append(lista, item)
	}
	return lista, nil
}

// rowValues escanea una fila completa del result set en un slice de valores.
func rowValues(rows *sql.Rows, n int) ([]any, error) {
	values := make([]any, n)
	pointers := make([]any, n)
	for i := range values {
		pointers[i] = &values[i]
	}
	if err := rows.Scan(pointers...); err != nil {
		return nil, err
	}
	return values, nil
}

// valueToInt convierte un valor crudo de database/sql a entero (0 si es nil).
func valueToInt(v any) int {
	switch n := v.(type) {
	case int64:
		return int(n)
	case int32:
		return int(n)
	case int:
		return n
	case float64:
		return int(n)
	case []byte:
		if i, err := strconv.Atoi(string(n)); err == nil {
			return i
		}
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i
		}
	}
	return 0
}

// valueToString convierte un valor crudo de database/sql a cadena ("" si es nil).
func valueToString(v any) string {
	switch s := v.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", s)
	}
}
