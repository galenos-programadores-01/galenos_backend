package sqlserver

import (
	"context"
	"database/sql"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type MedicoListaEsperaRepository struct {
	db *sql.DB
}

func NewMedicoListaEsperaRepository(db *sql.DB) *MedicoListaEsperaRepository {
	return &MedicoListaEsperaRepository{db: db}
}

func (r *MedicoListaEsperaRepository) Listar(ctx context.Context) ([]domain.MedicoListaEspera, error) {
	query := "EXEC usp_go_ListarMedicosListaEspera"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.MedicoListaEspera
	for rows.Next() {
		var m domain.MedicoListaEspera
		var apellidoPaterno sql.NullString
		var apellidoMaterno sql.NullString
		var nombres sql.NullString
		var idMedico sql.NullInt64
		var dmedico sql.NullString

		if err := rows.Scan(&apellidoPaterno, &apellidoMaterno, &nombres, &idMedico, &dmedico); err != nil {
			continue
		}
		if apellidoPaterno.Valid {
			m.ApellidoPaterno = apellidoPaterno.String
		}
		if apellidoMaterno.Valid {
			m.ApellidoMaterno = apellidoMaterno.String
		}
		if nombres.Valid {
			m.Nombres = nombres.String
		}
		if idMedico.Valid {
			m.IdMedico = int(idMedico.Int64)
		}
		if dmedico.Valid {
			m.Dmedico = &dmedico.String
		}
		lista = append(lista, m)
	}
	return lista, rows.Err()
}
