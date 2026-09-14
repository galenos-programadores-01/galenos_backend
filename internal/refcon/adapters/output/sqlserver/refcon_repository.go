// Package refconsql implementa el repositorio local del módulo refcon sobre
// SQL Server, autocontenido para no acoplarse a la capa de persistencia base.
package refconsql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
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

// ListarTodasLasUps ejecuta el procedimiento ListarUPS, que devuelve el
// catálogo completo de Unidades Productoras de Servicios (código y descripción)
// directamente desde la tabla SuSalud_ups.
func (r *RefConRepository) ListarTodasLasUps(ctx context.Context) ([]domain.Ups, error) {
	query := "EXEC ListarUPS"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("calling ListarUPS: %w", err)
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

// ListarEstablecimientos ejecuta el procedimiento usp_go_ListarEstablecimientos,
// que devuelve el catálogo de establecimientos de salud (código IPRESS y nombre).
func (r *RefConRepository) ListarEstablecimientos(ctx context.Context) ([]domain.Establecimiento, error) {
	query := "EXEC usp_go_ListarEstablecimientos"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarEstablecimientos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading establecimientos: %w", err)
	}

	lista := make([]domain.Establecimiento, 0, len(maps))
	for _, m := range maps {
		var item domain.Establecimiento
		if v := rowString(m, "Codigo", "codigo"); v != nil {
			item.Codigo = *v
		}
		if v := rowString(m, "Nombre", "nombre"); v != nil {
			item.Nombre = *v
		}
		lista = append(lista, item)
	}
	return lista, nil
}

// ListarDistritosPorIdReniec ejecuta el procedimiento
// usp_go_ListarDistritosPorIdReniec, que devuelve el distrito asociado a un
// código RENIEC (IdReniec).
func (r *RefConRepository) ListarDistritosPorIdReniec(ctx context.Context, idReniec int) ([]domain.DistritoReniec, error) {
	query := "EXEC usp_go_ListarDistritosPorIdReniec @IdReniec = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idReniec))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarDistritosPorIdReniec: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading distritos por idreniec: %w", err)
	}

	lista := make([]domain.DistritoReniec, 0, len(maps))
	for _, m := range maps {
		lista = append(lista, domain.DistritoReniec{
			IdDistrito:  int64(valueToInt(m["IdDistrito"])),
			Nombre:      valueToString(m["Nombre"]),
			IdReniec:    int64(valueToInt(m["IdReniec"])),
			IdProvincia: int64(valueToInt(m["IdProvincia"])),
		})
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

// rowsToMaps convierte el result set en un slice de mapas columna->valor.
func rowsToMaps(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var lista []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(columns))
		for i, col := range columns {
			row[col] = values[i]
		}
		lista = append(lista, row)
	}
	return lista, rows.Err()
}

// rowString devuelve el primer valor no nulo de las columnas dadas como
// cadena, o nil si ninguna existe (o todas son nulas).
func rowString(m map[string]any, names ...string) *string {
	for _, name := range names {
		if v, ok := m[name]; ok {
			if s := valueToString(v); s != "" {
				return &s
			}
		}
	}
	return nil
}
