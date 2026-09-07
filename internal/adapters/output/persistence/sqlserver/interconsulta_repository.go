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
	query := "EXEC usp_go_ListarInterconsultaPorId @IdInterconsulta = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", id))
	if err != nil {
		return nil, fmt.Errorf("consultando interconsulta %d: %w", id, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colMap := make(map[string]int)
	for i, c := range cols {
		colMap[c] = i
	}

	vals := make([]interface{}, len(cols))
	valPtrs := make([]interface{}, len(cols))
	for i := range vals {
		valPtrs[i] = &vals[i]
	}

	if err := rows.Scan(valPtrs...); err != nil {
		return nil, fmt.Errorf("escaneando interconsulta %d: %w", id, err)
	}

	var ic domain.Interconsulta
	ic.IdInterconsulta = id

	if idx, ok := colMap["IdRegistroAtencion"]; ok && vals[idx] != nil {
		if v, ok := vals[idx].(int64); ok {
			ic.IdAtencionOrigen = int(v)
		}
	}
	if idx, ok := colMap["IdEspecialidadSolicitada"]; ok && vals[idx] != nil {
		if v, ok := vals[idx].(int64); ok {
			ic.IdEspecialidad = int(v)
		}
	}
	if idx, ok := colMap["IdEmpleadoRespuesta"]; ok && vals[idx] != nil {
		if v, ok := vals[idx].(int64); ok {
			ic.IdMedicoDestino = int(v)
		}
	}
	if idx, ok := colMap["MotivoClinico"]; ok && vals[idx] != nil {
		if v, ok := vals[idx].(string); ok {
			ic.Motivo = v
		}
	}
	if idx, ok := colMap["FechaSolicitud"]; ok && vals[idx] != nil {
		if t, ok := vals[idx].(time.Time); ok {
			ic.FechaSolicitud = t.Format("2006-01-02T15:04:05")
		} else if s, ok := vals[idx].(string); ok {
			ic.FechaSolicitud = s
		}
	}
	if idx, ok := colMap["IdEstado"]; ok && vals[idx] != nil {
		if v, ok := vals[idx].(int64); ok {
			if v == 1 {
				ic.Estado = "SOLICITUD"
			} else {
				ic.Estado = "RESPUESTA"
			}
		}
	}

	return &ic, nil
}

func (r *InterconsultaRepository) ListarPorServicio(ctx context.Context, tipoServicio string) ([]domain.Interconsulta, error) {
	query := `EXEC usp_go_ListarInterconsultasSegunTipoServicio 
		@IdCuentaAtencion = -100,
		@IdTipoServicio = @p1,
		@NroDocumento = '-100',
		@NroHistoriaClinica = '-100',
		@ApePaterno = '-100',
		@ApeMaterno = '-100',
		@fini = '2000-01-01',
		@ffin = '2030-12-31'`

	idTipo := 1
	if tipoServicio != "" {
		fmt.Sscanf(tipoServicio, "%d", &idTipo)
	}

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idTipo))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.Interconsulta
	for rows.Next() {
		var idInterconsulta, idCuentaAtencion, idEspecialidad int
		var idRegAtencion, nroDoc, paciente, estadoInterconsulta, servSol, espRpta, estado, medSol, motivo, opDiag, manPac, trasPac, prioridad, especialidad sql.NullString
		var fechaSol sql.NullTime

		if err := rows.Scan(
			&idInterconsulta, &idRegAtencion, &idCuentaAtencion, &nroDoc, &fechaSol,
			&paciente, &estadoInterconsulta, &idEspecialidad, &servSol, &espRpta,
			&estado, &medSol, &motivo, &opDiag, &manPac, &trasPac, &prioridad, &especialidad,
		); err != nil {
			log.Printf("escaneando interconsulta por servicio: %v", err)
			continue
		}

		var ic domain.Interconsulta
		ic.IdInterconsulta = idInterconsulta
		ic.IdEspecialidad = idEspecialidad
		ic.Motivo = motivo.String
		ic.Estado = estado.String
		if fechaSol.Valid {
			ic.FechaSolicitud = fechaSol.Time.Format("2006-01-02T15:04:05")
		}
		lista = append(lista, ic)
	}
	return lista, rows.Err()
}

func (r *InterconsultaRepository) ListarPorAtencion(ctx context.Context, idAtencion int) ([]domain.Interconsulta, error) {
	query := "EXEC usp_go_ListarInterconsultasPorAtencion @IdAtencion = @p1"
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
			ic.FechaSolicitud = fecha.Time.Format("2006-01-02T15:04:05")
		}
		if estado.Valid {
			ic.Estado = estado.String
		}
		lista = append(lista, ic)
	}
	return lista, rows.Err()
}

func (r *InterconsultaRepository) Guardar(ctx context.Context, ic domain.Interconsulta) error {
	query := "EXEC dbo.usp_go_InterconsultaCrear @IdAtencionOrigen = @p1, @IdEspecialidad = @p2, @IdMedicoDestino = @p3, @Motivo = @p4"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", ic.IdAtencionOrigen), sql.Named("p2", ic.IdEspecialidad), sql.Named("p3", ic.IdMedicoDestino), sql.Named("p4", ic.Motivo))
	return err
}

func (r *InterconsultaRepository) ActualizarEstado(ctx context.Context, id int, estado string) error {
	tipo := 1
	if estado == "RESPUESTA" || estado == "2" {
		tipo = 2
	}
	query := "EXEC usp_go_InterconsultaActualizaEstadoFirmado @IdInterconsulta = @p1, @Tipo = @p2"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", id), sql.Named("p2", tipo))
	return err
}

func (r *InterconsultaRepository) GuardarFirma(ctx context.Context, firma domain.FirmaInterconsulta) error {
	query := "EXEC usp_go_InterconsultaActualizaEstadoFirmado @IdInterconsulta = @p1, @Tipo = @p2"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", firma.IdInterconsulta), sql.Named("p2", 1))
	return err
}

func (r *InterconsultaRepository) ListarEspecialidades(ctx context.Context) ([]domain.EspecialidadInterconsulta, error) {
	query := "EXEC dbo.usp_go_EspecialidadesXidTipoServicio @IdTipoServicio = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", 1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.EspecialidadInterconsulta
	for rows.Next() {
		var esp domain.EspecialidadInterconsulta
		var nombre sql.NullString
		var descLarga sql.NullString
		if err := rows.Scan(&esp.IdEspecialidad, &nombre, &descLarga); err != nil {
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
	query := "EXEC usp_go_MedicosFiltrarPorIdEspecialidad @IdEspecialialidad = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", IdEspecialidad))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.MedicoInterconsulta
	for rows.Next() {
		var m domain.MedicoInterconsulta
		var idMedico int
		var cod, apePat, apeMat, nom, esp, col, rne sql.NullString
		if err := rows.Scan(&idMedico, &cod, &apePat, &apeMat, &nom, &esp, &col, &rne); err != nil {
			log.Printf("escaneando médico por especialidad: %v", err)
			continue
		}
		m.IdMedico = idMedico
		if cod.Valid {
			m.CodigoPlanilla = &cod.String
		}
		nombreCompleto := fmt.Sprintf("%s %s, %s", apePat.String, apeMat.String, nom.String)
		m.Medico = &nombreCompleto
		lista = append(lista, m)
	}
	return lista, rows.Err()
}
