package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// CausaExternaMorbilidadRepository accede al catálogo de causas externas
// de morbilidad en SQL Server.
type CausaExternaMorbilidadRepository struct {
	db *sql.DB
}

func NewCausaExternaMorbilidadRepository(db *sql.DB) *CausaExternaMorbilidadRepository {
	return &CausaExternaMorbilidadRepository{db: db}
}

// Listar invoca el SP uso_go_EmergenciaCausaExternaMorbilidadListar y
// devuelve todas las causas externas de morbilidad.
func (r *CausaExternaMorbilidadRepository) Listar(ctx context.Context) ([]domain.CausaExternaMorbilidad, error) {
	query := "EXEC dbo.uso_go_EmergenciaCausaExternaMorbilidadListar"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando causas externas de morbilidad: %w", err)
	}
	defer rows.Close()

	var causas []domain.CausaExternaMorbilidad
	for rows.Next() {
		var (
			id          sql.NullInt64
			descripcion sql.NullString
		)
		if err := rows.Scan(&id, &descripcion); err != nil {
			return nil, fmt.Errorf("error escaneando causa externa de morbilidad: %w", err)
		}
		causas = append(causas, domain.CausaExternaMorbilidad{
			IdCausaExternaMorbilidad: int(id.Int64),
			Descripcion:              descripcion.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return causas, nil
}
