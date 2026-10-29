// Package refconsql implementa el repositorio local del módulo refcon sobre
// SQL Server, autocontenido para no acoplarse a la capa de persistencia base.
package refconsql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/galenos-pro/appointments-api/internal/refcon/domain"
)

type RefConRepository struct {
	db *sql.DB
}

func NewRefConRepository(db *sql.DB) *RefConRepository {
	return &RefConRepository{db: db}
}

// ListarReferenciasPorMes ejecuta usp_go_ListarReferenciasPorMesyAnio.
// El SP recibe un rango de fechas (@FechaInicio, @FechaFin) y, opcionalmente,
// un estado de referencia (@estadoreferencia, 0 = todos). Devuelve siempre 5
// columnas: descripción del estado, conteo, número de mes, nombre del mes y
// cantidad; el mapeo se hace por nombre de columna.
func (r *RefConRepository) ListarReferenciasPorMes(ctx context.Context, fechaInicio, fechaFin time.Time, estadoReferencia int) ([]domain.ReferenciaPorMes, error) {
	query := "EXEC dbo.usp_go_ListarReferenciasPorMesyAnio @FechaInicio = @p1, @FechaFin = @p2, @estadoreferencia = @p3"
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", fechaInicio),
		sql.Named("p2", fechaFin),
		sql.Named("p3", estadoReferencia),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarReferenciasPorMesyAnio: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading usp_go_ListarReferenciasPorMesyAnio rows: %w", err)
	}

	lista := make([]domain.ReferenciaPorMes, 0, len(maps))
	for _, m := range maps {
		lista = append(lista, domain.ReferenciaPorMes{
			Estado:    valueToString(m["Descripcion"]),
			Conteo:    valueToInt(m["Conteo"]),
			Mes:       valueToInt(m["MesNum"]),
			NombreMes: valueToString(m["Mes"]),
			Cantidad:  valueToInt(m["Cantidad"]),
		})
	}
	return lista, nil
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
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
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
