package sqlserver

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func rowsToMaps(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("reading columns: %w", err)
	}

	results := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
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
		results = append(results, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return results, nil
}

// rowString devuelve el valor string de una columna (búsqueda por nombre
// exacto o case-insensitive). Retorna nil si la columna no existe o es NULL.
func rowString(row map[string]any, names ...string) *string {
	for _, name := range names {
		if value, ok := lookup(row, name); ok {
			switch v := value.(type) {
			case nil:
				return nil
			case string:
				return &v
			case []byte:
				s := string(v)
				return &s
			default:
				s := fmt.Sprintf("%v", v)
				return &s
			}
		}
	}
	return nil
}

// rowInt64 devuelve el valor entero de una columna. Retorna nil si la
// columna no existe o es NULL.
func rowInt64(row map[string]any, names ...string) *int64 {
	for _, name := range names {
		if value, ok := lookup(row, name); ok {
			switch v := value.(type) {
			case nil:
				return nil
			case int64:
				return &v
			case int:
				i := int64(v)
				return &i
			case float64:
				i := int64(v)
				return &i
			}
		}
	}
	return nil
}

// rowTime devuelve el valor de fecha/hora de una columna. Retorna nil si la
// columna no existe o es NULL.
func rowTime(row map[string]any, names ...string) *time.Time {
	for _, name := range names {
		if value, ok := lookup(row, name); ok {
			switch v := value.(type) {
			case nil:
				return nil
			case time.Time:
				return &v
			case string:
				if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
					return &t
				}
				if t, err := time.Parse("2006-01-02", v); err == nil {
					return &t
				}
			}
		}
	}
	return nil
}

// rowBool devuelve el valor booleano de una columna. Retorna nil si la
// columna no existe o es NULL.
func rowBool(row map[string]any, names ...string) *bool {
	for _, name := range names {
		if value, ok := lookup(row, name); ok {
			switch v := value.(type) {
			case nil:
				return nil
			case bool:
				return &v
			case int64:
				b := v != 0
				return &b
			case int:
				b := v != 0
				return &b
			}
		}
	}
	return nil
}

// rowFloat64 devuelve el valor numérico de punto flotante de una columna.
// Retorna nil si la columna no existe o es NULL.
func rowFloat64(row map[string]any, names ...string) *float64 {
	for _, name := range names {
		if value, ok := lookup(row, name); ok {
			switch v := value.(type) {
			case nil:
				return nil
			case float64:
				return &v
			case float32:
				f := float64(v)
				return &f
			case int64:
				f := float64(v)
				return &f
			case int:
				f := float64(v)
				return &f
			case []byte:
				f, err := strconv.ParseFloat(string(v), 64)
				if err == nil {
					return &f
				}
			case string:
				f, err := strconv.ParseFloat(v, 64)
				if err == nil {
					return &f
				}
			}
		}
	}
	return nil
}

// boolToInt64 convierte un *bool en *int64 (1/0) o nil si el puntero es nil.
func boolToInt64(p *bool) *int64 {
	if p == nil {
		return nil
	}
	var v int64
	if *p {
		v = 1
	}
	return &v
}

// lookup busca una columna por nombre exacto y luego case-insensitive.
func lookup(row map[string]any, name string) (any, bool) {
	if value, ok := row[name]; ok {
		return value, true
	}
	for key, value := range row {
		if strings.EqualFold(key, name) {
			return value, true
		}
	}
	return nil, false
}
