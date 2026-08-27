package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sqlServerDiagnosticoRepository struct {
	db *sql.DB
}

func NewSqlServerDiagnosticoRepository(db *sql.DB) output.DiagnosticoRepository {
	return &sqlServerDiagnosticoRepository{db: db}
}

func (r *sqlServerDiagnosticoRepository) SearchDiagnosticos(ctx context.Context, filtro string, idAtencion, idPaciente int) ([]domain.DiagnosticoBusqueda, error) {
	query := "EXEC usp_go_SelectDiagnosticos @Filtro = @p1, @IdAtencion = @p2, @IdPaciente = @p3"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", filtro), sql.Named("p2", idAtencion), sql.Named("p3", idPaciente))
	if err != nil {
		return nil, fmt.Errorf("error querying diagnosticos: %w", err)
	}
	defer rows.Close()

	var results []domain.DiagnosticoBusqueda
	for rows.Next() {
		var d domain.DiagnosticoBusqueda
		var eMax, eMin, idSexo, yaReg sql.NullInt64
		var intra, activo, cancer sql.NullBool
		var desc, cie10, descL sql.NullString

		if err := rows.Scan(
			&d.IdDiagnostico,
			&intra,
			&desc,
			&cie10,
			&activo,
			&descL,
			&eMax,
			&eMin,
			&idSexo,
			&cancer,
			&yaReg,
		); err != nil {
			return nil, fmt.Errorf("error scanning diagnostico: %w", err)
		}

		if intra.Valid {
			if intra.Bool {
				d.Intrahospitalario = 1
			} else {
				d.Intrahospitalario = 0
			}
		}
		if desc.Valid {
			d.Descripcion = desc.String
		}
		if cie10.Valid {
			d.CodigoCIE10 = cie10.String
		}
		if activo.Valid {
			if activo.Bool {
				d.EsActivo = 1
			} else {
				d.EsActivo = 0
			}
		}
		if cancer.Valid {
			if cancer.Bool {
				d.Cancer = 1
			} else {
				d.Cancer = 0
			}
		}
		if descL.Valid {
			d.DescripcionLarga = descL.String
		}
		if eMax.Valid {
			d.EdadMaxDias = int(eMax.Int64)
		}
		if eMin.Valid {
			d.EdadMinDias = int(eMin.Int64)
		}
		if idSexo.Valid {
			d.IdTipoSexo = int(idSexo.Int64)
		}
		if yaReg.Valid {
			d.YaRegistrado = int(yaReg.Int64)
		}

		results = append(results, d)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando diagnosticos: %w", err)
	}

	log.Printf("Repository: Returning %d rows for filtro=%q, idAtencion=%d, idPaciente=%d", len(results), filtro, idAtencion, idPaciente)
	return results, nil
}

func (r *sqlServerDiagnosticoRepository) ListarDiagnosticos(ctx context.Context, filtro string) ([]domain.DiagnosticoSimple, error) {
	query := "EXEC usp_go_ListarDiagnosticos @Filtro = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", filtro))
	if err != nil {
		return nil, fmt.Errorf("error listing diagnosticos: %w", err)
	}
	defer rows.Close()

	var results []domain.DiagnosticoSimple
	for rows.Next() {
		var d domain.DiagnosticoSimple
		if err := rows.Scan(&d.IdDiagnostico, &d.CodigoCIE10, &d.Descripcion); err != nil {
			return nil, fmt.Errorf("error scanning diagnostico simple: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}
