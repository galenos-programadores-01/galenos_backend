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

	if telefono := r.obtenerTelefonoListaEspera(ctx, id); telefono != nil {
		item.Telefono = telefono
	}

	return item, nil
}

func (r *ListaEsperaQxRepository) obtenerTelefonoListaEspera(ctx context.Context, id int) *string {
	var telefono sql.NullString
	if err := r.db.QueryRowContext(ctx, "EXEC dbo.usp_go_ListaEsperaQx_ObtenerTelefonoPorId @Id = @p1", sql.Named("p1", id)).Scan(&telefono); err != nil {
		return nil
	}
	if !telefono.Valid || strings.TrimSpace(telefono.String) == "" {
		return nil
	}
	return &telefono.String
}

func (r *ListaEsperaQxRepository) actualizarTelefonosListaEspera(ctx context.Context, lista []domain.ListaEsperaQxReporte) []domain.ListaEsperaQxReporte {
	if len(lista) == 0 {
		return lista
	}
	idStrs := make([]string, len(lista))
	for i, item := range lista {
		idStrs[i] = fmt.Sprintf("%d", item.Id)
	}
	query := "EXEC dbo.usp_go_ListaEsperaQx_ListarTelefonosPorIds @Ids = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", strings.Join(idStrs, ",")))
	if err != nil {
		return lista
	}
	defer rows.Close()
	telefonos := make(map[int]string, len(lista))
	for rows.Next() {
		var id int
		var telefono sql.NullString
		if err := rows.Scan(&id, &telefono); err != nil {
			continue
		}
		if telefono.Valid && strings.TrimSpace(telefono.String) != "" {
			telefonos[id] = telefono.String
		}
	}
	for i := range lista {
		if t, ok := telefonos[lista[i].Id]; ok {
			lista[i].Telefono = t
		}
	}
	return lista
}

func ptrStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (r *ListaEsperaQxRepository) Modificar(ctx context.Context, item domain.ListaEsperaQxModificar) error {
	query := `EXEC usp_go_ModificarListaEsperaQx
		@id = @p1, @observacion = @p2, @Telefono = @p3, @diagnostico = @p4, @IdEspecialidad = @p5,
		@fechaorden = @p6, @fechalab = @p7, @fechaic_cardio = @p8,
		@fechaic_neumo = @p9, @fechaic_anestesio = @p10`
	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", item.Id),
		sql.Named("p2", item.Observacion),
		sql.Named("p3", item.Telefono),
		sql.Named("p4", item.Diagnostico),
		sql.Named("p5", item.IdEspecialidad),
		sql.Named("p6", item.FechaOrden),
		sql.Named("p7", item.FechaLab),
		sql.Named("p8", item.FechaICCardio),
		sql.Named("p9", item.FechaICNeumo),
		sql.Named("p10", item.FechaICAnestesio),
	)
	return err
}

// ActualizarHospitalizado invoca el procedimiento almacenado
// usp_go_ActualizarListaEsperaQx, que marca el registro de la lista de espera
// quirurgica como HOSPITALIZADO (Tipo = 'HOSPITALIZADO').
func (r *ListaEsperaQxRepository) ActualizarHospitalizado(ctx context.Context, id int) error {
	query := "EXEC usp_go_ActualizarListaEsperaQx @IdListaEspera = @p1"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", id))
	return err
}

func (r *ListaEsperaQxRepository) Crear(ctx context.Context, item domain.ListaEsperaQxCrear, idUsuario int) error {
	query := `EXEC usp_go_GrabarListaEsperaQx
		@FechaOrden = @p1, @IdPaciente = @p2, @Diagnostico = @p3, @IdEspecialidad = @p4,
		@FechaLab = @p5, @FechaIC_Cardio = @p6, @FechaIC_Neumo = @p7,
		@FechaIC_Anestesio = @p8, @IdMedico = @p9, @Observacion = @p10,
		@Telefono = @p11, @IdUsuarioRegistra = @p12`
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
		sql.Named("p11", item.Telefono),
		sql.Named("p12", idUsuario),
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

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Escaneo dinámico: se conserva el mapeo posicional original (índices 0..16)
	// y se busca IdTipoSexo por nombre, para tolerar que el SP agregue la columna
	// en cualquier posición del SELECT.
	asString := func(v any) string {
		switch t := v.(type) {
		case nil:
			return ""
		case string:
			return t
		case []byte:
			return string(t)
		case time.Time:
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02")
		default:
			return fmt.Sprintf("%v", t)
		}
	}
	asInt := func(v any) int {
		switch t := v.(type) {
		case int64:
			return int(t)
		case int32:
			return int(t)
		case int:
			return t
		case float64:
			return int(t)
		case []byte:
			n, _ := strconv.Atoi(strings.TrimSpace(string(t)))
			return n
		case string:
			n, _ := strconv.Atoi(strings.TrimSpace(t))
			return n
		default:
			return 0
		}
	}

	var lista []domain.ListaEsperaQxReporte
	for rows.Next() {
		raw := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range raw {
			ptrs[i] = &raw[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			continue
		}

		at := func(i int) any {
			if i >= 0 && i < len(raw) {
				return raw[i]
			}
			return nil
		}

		var item domain.ListaEsperaQxReporte
		item.Id = asInt(at(0))
		item.NroHistoriaClinica = asInt(at(1))
		item.NroDocumento = asString(at(2))
		item.Paciente = asString(at(3))
		item.Edad = asInt(at(4))
		item.Telefono = asString(at(5))
		item.FechaOrden = asString(at(6))
		item.Especialidad = asString(at(7))
		item.Diagnostico = asString(at(8))
		item.FechaLab = asString(at(9))
		item.FechaICCardio = asString(at(10))
		item.FechaICNeumo = asString(at(11))
		item.FechaICAnestesio = asString(at(12))
		item.Medico = asString(at(13))
		item.Observacion = asString(at(14))
		// El índice 15 corresponde al teléfono alterno (no se usa).
		item.DiasTranscurridos = asInt(at(16))

		for i, c := range cols {
			if strings.EqualFold(c, "IdTipoSexo") {
				item.IdTipoSexo = asInt(at(i))
				break
			}
		}

		lista = append(lista, item)
	}
	return r.actualizarTelefonosListaEspera(ctx, lista), rows.Err()
}
