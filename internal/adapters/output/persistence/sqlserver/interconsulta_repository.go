package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type InterconsultaRepository struct {
	db *sql.DB
}

func NewInterconsultaRepository(db *sql.DB) *InterconsultaRepository {
	return &InterconsultaRepository{db: db}
}

func (r *InterconsultaRepository) ObtenerPorId(ctx context.Context, id int) (*domain.Interconsulta, error) {
	query := "EXEC webListarInterconsultaPorId @Id = @p1"
	row := r.db.QueryRowContext(ctx, query, sql.Named("p1", id))

	var ic domain.Interconsulta
	var fecha sql.NullTime
	var estado sql.NullString

	err := row.Scan(&ic.IdInterconsulta, &ic.IdAtencionOrigen, &ic.IdEspecialidad, &ic.IdMedicoDestino, &ic.Motivo, &fecha, &estado)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("escaneando interconsulta %d: %w", id, err)
	}
	if fecha.Valid {
		ic.FechaSolicitud = fecha.Time.Format(time.RFC3339)
	}
	if estado.Valid {
		ic.Estado = estado.String
	}
	return &ic, nil
}

func (r *InterconsultaRepository) ListarPorServicio(ctx context.Context, tipoServicio string) ([]domain.Interconsulta, error) {
	query := "EXEC WebListarInterconsultasSegunTipoServicio @TipoServicio = @p1"
	return r.ejecutarConsultaInterconsultas(ctx, query, sql.Named("p1", tipoServicio))
}

func (r *InterconsultaRepository) ListarPorAtencion(ctx context.Context, idAtencion int) ([]domain.Interconsulta, error) {
	query := "EXEC WebListarInterconsultasPorAtencion @IdAtencion = @p1"
	return r.ejecutarConsultaInterconsultas(ctx, query, sql.Named("p1", idAtencion))
}

func (r *InterconsultaRepository) ejecutarConsultaInterconsultas(ctx context.Context, query string, arg interface{}) ([]domain.Interconsulta, error) {
	rows, err := r.db.QueryContext(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.Interconsulta
	for rows.Next() {
		var ic domain.Interconsulta
		var fecha sql.NullTime
		var estado sql.NullString

		if err := rows.Scan(&ic.IdInterconsulta, &ic.IdAtencionOrigen, &ic.IdEspecialidad, &ic.IdMedicoDestino, &ic.Motivo, &fecha, &estado); err != nil {
			log.Printf("escaneando interconsulta: %v", err)
			continue
		}
		if fecha.Valid {
			ic.FechaSolicitud = fecha.Time.Format(time.RFC3339)
		}
		if estado.Valid {
			ic.Estado = estado.String
		}
		lista = append(lista, ic)
	}
	return lista, rows.Err()
}

func (r *InterconsultaRepository) Guardar(ctx context.Context, ic domain.Interconsulta) error {
	query := "EXEC webInterconsultasGrabar @IdAtencionOrigen = @p1, @IdEspecialidad = @p2, @IdMedicoDestino = @p3, @Motivo = @p4"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", ic.IdAtencionOrigen), sql.Named("p2", ic.IdEspecialidad), sql.Named("p3", ic.IdMedicoDestino), sql.Named("p4", ic.Motivo))
	return err
}

func (r *InterconsultaRepository) ActualizarEstado(ctx context.Context, id int, estado string) error {
	query := "EXEC webInterconsultaActualizaEstadoFirmado @IdInterconsulta = @p1, @Estado = @p2"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", id), sql.Named("p2", estado))
	return err
}

func (r *InterconsultaRepository) GuardarFirma(ctx context.Context, firma domain.FirmaInterconsulta) error {
	query := "EXEC webInterconsultaConsultarFirma @IdInterconsulta = @p1, @IdEmpleado = @p2, @DataB64 = @p3"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", firma.IdInterconsulta), sql.Named("p2", firma.IdEmpleadoFirma), sql.Named("p3", firma.DataB64))
	return err
}

func (r *InterconsultaRepository) ListarEspecialidades(ctx context.Context) ([]domain.EspecialidadInterconsulta, error) {
	query := "EXEC webEspecialidadesListarInterConsulta"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.EspecialidadInterconsulta
	for rows.Next() {
		var esp domain.EspecialidadInterconsulta
		var nombre sql.NullString
		var descLarga sql.NullString
		if err := rows.Scan(&esp.IdEspecialidad, &descLarga, &nombre); err != nil {
			log.Printf("escaneando especialidad: %v", err)
			continue
		}
		if nombre.Valid {
			esp.Nombre = &nombre.String
		}
		if descLarga.Valid {
			esp.DescripcionLarga = &descLarga.String
		}
		lista = append(lista, esp)
	}
	return lista, rows.Err()
}

func (r *InterconsultaRepository) ListarMedicosPorEspecialidad(ctx context.Context, IdEspecialidad int) ([]domain.MedicoInterconsulta, error) {
	query := "EXEC InterconsultaFiltrrMedicoXIdEspecialidad @IdEspecialidad = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", IdEspecialidad))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.MedicoInterconsulta
	for rows.Next() {
		var m domain.MedicoInterconsulta
		var codigo sql.NullString
		var medico sql.NullString
		if err := rows.Scan(&m.IdMedico, &m.IdEmpleado, &codigo, &medico); err != nil {
			log.Printf("escaneando médico por especialidad: %v", err)
			continue
		}
		if codigo.Valid {
			m.CodigoPlanilla = &codigo.String
		}
		if medico.Valid {
			m.Medico = &medico.String
		}
		lista = append(lista, m)
	}
	return lista, rows.Err()
}
