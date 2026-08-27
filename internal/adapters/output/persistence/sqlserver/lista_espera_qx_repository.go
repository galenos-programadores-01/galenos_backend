package sqlserver

import (
	"context"
	"database/sql"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ListaEsperaQxRepository struct {
	db *sql.DB
}

func NewListaEsperaQxRepository(db *sql.DB) *ListaEsperaQxRepository {
	return &ListaEsperaQxRepository{db: db}
}

func (r *ListaEsperaQxRepository) Listar(ctx context.Context, fecha string, paciente string) ([]domain.ListaEsperaQx, error) {
	query := "EXEC usp_go_ListaEsperaQxListar @Fecha = @p1, @Paciente = @p2"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", fecha), sql.Named("p2", paciente))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.ListaEsperaQx
	for rows.Next() {
		var item domain.ListaEsperaQx
		var nroHistoriaClinica sql.NullInt64
		var nroDocumento sql.NullString
		var pacienteNombre sql.NullString
		var fechaNacimiento sql.NullTime
		var fechaOrden sql.NullTime
		var observacion sql.NullString

		if err := rows.Scan(&nroHistoriaClinica, &nroDocumento, &pacienteNombre, &fechaNacimiento, &fechaOrden, &observacion); err != nil {
			continue
		}
		if nroHistoriaClinica.Valid {
			item.NroHistoriaClinica = int(nroHistoriaClinica.Int64)
		}
		if nroDocumento.Valid {
			item.NroDocumento = nroDocumento.String
		}
		if pacienteNombre.Valid {
			item.Paciente = pacienteNombre.String
		}
		if fechaNacimiento.Valid {
			item.FechaNacimiento = fechaNacimiento.Time.Format("2006-01-02")
		}
		if fechaOrden.Valid {
			item.FechaOrden = fechaOrden.Time.Format("2006-01-02")
		}
		if observacion.Valid {
			item.Observacion = &observacion.String
		}
		lista = append(lista, item)
	}
	return lista, rows.Err()
}

func (r *ListaEsperaQxRepository) Crear(ctx context.Context, item domain.ListaEsperaQxCrear) error {
	query := `EXEC usp_go_ListaEsperaQxGrabar
		@IdPaciente = @p0, @IdMedico = @p19,
		@IdTipoDocumento = @p1, @NroDocumento = @p2, @ApellidoPaterno = @p3, @ApellidoMaterno = @p4,
		@PrimerNombre = @p5, @SegundoNombre = @p6, @FechaNacimiento = @p7, @IdSexo = @p8,
		@Telefono = @p9, @Direccion = @p10, @FechaOrden = @p11, @Diagnostico = @p12,
		@FechaLaboratorio = @p13, @FechaICCardio = @p14, @FechaICNeumo = @p15,
		@FechaICAnestesio = @p16, @Medico = @p17, @Observacion = @p18`
	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p0", item.IdPaciente),
		sql.Named("p1", item.IdTipoDocumento),
		sql.Named("p2", item.NroDocumento),
		sql.Named("p3", item.ApellidoPaterno),
		sql.Named("p4", item.ApellidoMaterno),
		sql.Named("p5", item.PrimerNombre),
		sql.Named("p6", item.SegundoNombre),
		sql.Named("p7", item.FechaNacimiento),
		sql.Named("p8", item.IdSexo),
		sql.Named("p9", item.Telefono),
		sql.Named("p10", item.Direccion),
		sql.Named("p11", item.FechaOrden),
		sql.Named("p12", item.Diagnostico),
		sql.Named("p13", item.FechaLaboratorio),
		sql.Named("p14", item.FechaICCardio),
		sql.Named("p15", item.FechaICNeumo),
		sql.Named("p16", item.FechaICAnestesio),
		sql.Named("p17", item.Medico),
		sql.Named("p18", item.Observacion),
		sql.Named("p19", item.IdMedico),
	)
	return err
}
