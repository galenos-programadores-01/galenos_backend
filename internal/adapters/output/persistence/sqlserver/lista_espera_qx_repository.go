package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ListaEsperaQxRepository struct {
	db *sql.DB
}

func NewListaEsperaQxRepository(db *sql.DB) *ListaEsperaQxRepository {
	return &ListaEsperaQxRepository{db: db}
}

func (r *ListaEsperaQxRepository) Listar(ctx context.Context, fecha string, fechaFin string, paciente string, idEspecialidad int) ([]domain.ListaEsperaQx, error) {
	query := "EXEC usp_go_ListaEsperaQxListar @Fecha = @p1, @FechaFin = @p2, @Paciente = @p3, @IdEspecialidad = @p4"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", fecha), sql.Named("p2", fechaFin), sql.Named("p3", paciente), sql.Named("p4", idEspecialidad))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.ListaEsperaQx
	for rows.Next() {
		var item domain.ListaEsperaQx
		var idListaEspera sql.NullInt64
		var nroHistoriaClinica sql.NullInt64
		var nroDocumento sql.NullString
		var pacienteNombre sql.NullString
		var edad sql.NullInt64
		var fechaOrden sql.NullTime
		var especialidad sql.NullString
		var observacion sql.NullString
		var diasTranscurridos sql.NullInt64

		if err := rows.Scan(&idListaEspera, &nroHistoriaClinica, &nroDocumento, &pacienteNombre, &edad, &fechaOrden, &especialidad, &observacion, &diasTranscurridos); err != nil {
			continue
		}
		if idListaEspera.Valid {
			item.IdListaEspera = int(idListaEspera.Int64)
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
		if edad.Valid {
			item.Edad = int(edad.Int64)
		}
		if fechaOrden.Valid {
			item.FechaOrden = fechaOrden.Time.Format("2006-01-02")
		}
		if especialidad.Valid {
			item.Especialidad = especialidad.String
		}
		if observacion.Valid {
			item.Observacion = &observacion.String
		}
		if diasTranscurridos.Valid {
			item.DiasTranscurridos = int(diasTranscurridos.Int64)
		}
		lista = append(lista, item)
	}
	return lista, rows.Err()
}

func (r *ListaEsperaQxRepository) ObtenerPorId(ctx context.Context, id int) (domain.ListaEsperaQxPaciente, error) {
	query := "EXEC usp_go_ListarPacienteListaEspera @id = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", id))
	if err != nil {
		return domain.ListaEsperaQxPaciente{}, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return domain.ListaEsperaQxPaciente{}, err
	}

	if !rows.Next() {
		return domain.ListaEsperaQxPaciente{}, sql.ErrNoRows
	}

	raw := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range raw {
		ptrs[i] = &raw[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		return domain.ListaEsperaQxPaciente{}, err
	}

	colMap := make(map[string]interface{}, len(cols))
	for i, c := range cols {
		colMap[strings.ToLower(c)] = raw[i]
	}

	var item domain.ListaEsperaQxPaciente

	getStr := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := colMap[strings.ToLower(k)]; ok && v != nil {
				switch val := v.(type) {
				case string:
					return val
				case []byte:
					return string(val)
				case time.Time:
					if val.IsZero() {
						return ""
					}
					return val.Format("2006-01-02")
				}
			}
		}
		return ""
	}
	getIntPtr := func(keys ...string) *int {
		for _, k := range keys {
			if v, ok := colMap[strings.ToLower(k)]; ok && v != nil {
				switch val := v.(type) {
				case int64:
					n := int(val)
					return &n
				case int32:
					n := int(val)
					return &n
				case int16:
					n := int(val)
					return &n
				case int8:
					n := int(val)
					return &n
				case int:
					return &val
				case uint8:
					n := int(val)
					return &n
				case float64:
					n := int(val)
					return &n
				case float32:
					n := int(val)
					return &n
				case []byte:
					if n, err := strconv.Atoi(string(val)); err == nil {
						return &n
					}
				case fmt.Stringer:
					if n, err := strconv.Atoi(val.String()); err == nil {
						return &n
					}
				}
			}
		}
		return nil
	}

	item.NroDocumento = getStr("NroDocumento")
	item.IdDocIdentidad = getIntPtr("IdDocIdentidad")
	item.ApellidoPaterno = getStr("ApellidoPaterno")
	item.ApellidoMaterno = getStr("ApellidoMaterno")
	item.PrimerNombre = getStr("PrimerNombre")
	item.Direccion = ptrStr(getStr("DireccionDomicilio"))
	item.Telefono = ptrStr(getStr("Telefono"))
	item.IdTipoSexo = getIntPtr("IdTipoSexo")
	item.FechaOrden = getStr("fechaorden")
	item.Diagnostico = ptrStr(getStr("diagnostico"))
	item.FechaICCardio = getStr("FechaIC_Cardio")
	item.FechaICNeumo = getStr("FechaIC_Neumo")
	item.FechaICAnestesio = getStr("FechaIC_Anestesio")
	item.IdMedico = getIntPtr("IdMedico")
	item.Medico = ptrStr(getStr("Medico"))

	item.FechaNacimiento = getStr("FechaNacimiento")
	item.FechaLab = getStr("FechaLab")
	item.Observacion = ptrStr(getStr("Observacion"))
	item.IdDiagnostico = getIntPtr("IdDiagnostico")
	item.IdEspecialidad = getIntPtr("IdEspecialidad")

	return item, nil
}

func ptrStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (r *ListaEsperaQxRepository) Modificar(ctx context.Context, item domain.ListaEsperaQxModificar) error {
	query := `EXEC usp_go_ModificarListaEsperaQx
		@id = @p1, @observacion = @p2, @diagnostico = @p3,  @IdEspecialidad = @p4,
		@fechaorden = @p5, @fechalab = @p6, @fechaic_cardio = @p7,
		@fechaic_neumo = @p8, @fechaic_anestesio = @p9`
	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", item.Id),
		sql.Named("p2", item.Observacion),
		sql.Named("p3", item.Diagnostico),
		sql.Named("p4", item.IdEspecialidad),
		sql.Named("p5", item.FechaOrden),
		sql.Named("p6", item.FechaLab),
		sql.Named("p7", item.FechaICCardio),
		sql.Named("p8", item.FechaICNeumo),
		sql.Named("p9", item.FechaICAnestesio),
	)
	return err
}

func (r *ListaEsperaQxRepository) Crear(ctx context.Context, item domain.ListaEsperaQxCrear, idUsuario int) error {
	query := `EXEC usp_go_GrabarListaEsperaQx
		@FechaOrden = @p1, @IdPaciente = @p2, @Diagnostico = @p3, @IdEspecialidad = @p4,
		@FechaLab = @p5, @FechaIC_Cardio = @p6, @FechaIC_Neumo = @p7,
		@FechaIC_Anestesio = @p8, @IdMedico = @p9, @Observacion = @p10,
		@IdUsuarioRegistra = @p11`
	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", item.FechaOrden),
		sql.Named("p2", item.IdPaciente),
		sql.Named("p3", item.Diagnostico),
		sql.Named("p4", item.IdEspecialidad),
		sql.Named("p5", item.FechaLab),
		sql.Named("p6", item.FechaICCardio),
		sql.Named("p7", item.FechaICNeumo),
		sql.Named("p8", item.FechaICAnestesio),
		sql.Named("p9", item.IdMedico),
		sql.Named("p10", item.Observacion),
		sql.Named("p11", idUsuario),
	)
	return err
}

func (r *ListaEsperaQxRepository) Reporte(ctx context.Context, fecha string, fechaFin string, idEspecialidad int) ([]domain.ListaEsperaQxReporte, error) {
	query := "EXEC usp_go_ListaEsperaQxReporte @Fecha = @p1, @FechaFin = @p2, @IdEspecialidad = @p3"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", fecha), sql.Named("p2", fechaFin), sql.Named("p3", idEspecialidad))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.ListaEsperaQxReporte
	for rows.Next() {
		var item domain.ListaEsperaQxReporte
		var id sql.NullInt64
		var nroHistoriaClinica sql.NullInt64
		var nroDocumento sql.NullString
		var paciente sql.NullString
		var edad sql.NullInt64
		var telefono sql.NullString
		var fechaOrden sql.NullTime
		var especialidad sql.NullString
		var diagnostico sql.NullString
		var fechaLab sql.NullTime
		var fechaICCardio sql.NullTime
		var fechaICNeumo sql.NullTime
		var fechaICAnestesio sql.NullTime
		var medico sql.NullString
		var observacion sql.NullString
		var diasTranscurridos sql.NullInt64

		if err := rows.Scan(&id, &nroHistoriaClinica, &nroDocumento, &paciente, &edad, &telefono, &fechaOrden, &especialidad, &diagnostico, &fechaLab, &fechaICCardio, &fechaICNeumo, &fechaICAnestesio, &medico, &observacion, &diasTranscurridos); err != nil {
			continue
		}
		if id.Valid {
			item.Id = int(id.Int64)
		}
		if nroHistoriaClinica.Valid {
			item.NroHistoriaClinica = int(nroHistoriaClinica.Int64)
		}
		if nroDocumento.Valid {
			item.NroDocumento = nroDocumento.String
		}
		if paciente.Valid {
			item.Paciente = paciente.String
		}
		if edad.Valid {
			item.Edad = int(edad.Int64)
		}
		if telefono.Valid {
			item.Telefono = telefono.String
		}
		if fechaOrden.Valid {
			item.FechaOrden = fechaOrden.Time.Format("2006-01-02")
		}
		if especialidad.Valid {
			item.Especialidad = especialidad.String
		}
		if diagnostico.Valid {
			item.Diagnostico = diagnostico.String
		}
		if fechaLab.Valid {
			item.FechaLab = fechaLab.Time.Format("2006-01-02")
		}
		if fechaICCardio.Valid {
			item.FechaICCardio = fechaICCardio.Time.Format("2006-01-02")
		}
		if fechaICNeumo.Valid {
			item.FechaICNeumo = fechaICNeumo.Time.Format("2006-01-02")
		}
		if fechaICAnestesio.Valid {
			item.FechaICAnestesio = fechaICAnestesio.Time.Format("2006-01-02")
		}
		if medico.Valid {
			item.Medico = medico.String
		}
		if observacion.Valid {
			item.Observacion = observacion.String
		}
		if diasTranscurridos.Valid {
			item.DiasTranscurridos = int(diasTranscurridos.Int64)
		}
		lista = append(lista, item)
	}
	return lista, rows.Err()
}
