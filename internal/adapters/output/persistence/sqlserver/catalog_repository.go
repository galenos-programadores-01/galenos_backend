package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type catalogRepository struct {
	db *sql.DB
}

// NewCatalogRepository construye el adaptador de catÃ¡logos de datos
// maestros, invocando los procedimientos almacenados que el frontend
// consume para poblar los formularios.
func NewCatalogRepository(db *sql.DB) output.CatalogRepository {
	return &catalogRepository{db: db}
}

func (r *catalogRepository) ListEtnias(ctx context.Context) ([]domain.Etnia, error) {
	const procedure = `ups_go_ListarEtnias`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling ups_go_ListarEtnias: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading etnias: %w", err)
	}

	items := make([]domain.Etnia, 0, len(maps))
	for _, m := range maps {
		descripcion := rowString(m, "desetni", "Descripcion")
		items = append(items, domain.Etnia{
			Codigo:      *rowString(m, "codetni", "Codigo"),
			Descripcion: descripcion,
		})
	}

	return items, nil
}

func (r *catalogRepository) ListIdiomas(ctx context.Context) ([]domain.Idioma, error) {
	const procedure = `ups_go_ListarIdiomas`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling ups_go_ListarIdiomas: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading idiomas: %w", err)
	}

	items := make([]domain.Idioma, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.Idioma{
			ID:     *rowInt64(m, "IdIdioma"),
			Lengua: rowString(m, "Lengua", "Nombre"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListTipoSexos(ctx context.Context) ([]domain.TipoSexo, error) {
	const procedure = `usp_go_ListarTiposSexos`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarTiposSexos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading tipos sexos: %w", err)
	}

	items := make([]domain.TipoSexo, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.TipoSexo{
			ID:          *rowInt64(m, "IdTipoSexo"),
			Descripcion: rowString(m, "Descripcion"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListEstadosCivil(ctx context.Context) ([]domain.TipoEstadoCivil, error) {
	const procedure = `usp_go_ListarEstadosCivil`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarEstadosCivil: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading estados civil: %w", err)
	}

	items := make([]domain.TipoEstadoCivil, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.TipoEstadoCivil{
			ID:          *rowInt64(m, "IdEstadoCivil"),
			Descripcion: rowString(m, "Descripcion"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListGradosInstruccion(ctx context.Context) ([]domain.TipoGradoInstruccion, error) {
	const procedure = `usp_go_ListarGradoInstruccion`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarGradoInstruccion: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading grados instruccion: %w", err)
	}

	items := make([]domain.TipoGradoInstruccion, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.TipoGradoInstruccion{
			ID:          *rowInt64(m, "IdGradoInstruccion"),
			Descripcion: rowString(m, "Descripcion"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListOcupaciones(ctx context.Context) ([]domain.TipoOcupacion, error) {
	const procedure = `usp_go_ListarOcupaciones`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarOcupaciones: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading ocupaciones: %w", err)
	}

	items := make([]domain.TipoOcupacion, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.TipoOcupacion{
			ID:          *rowInt64(m, "IdTipoOcupacion"),
			Descripcion: rowString(m, "descripcion", "Descripcion"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListTiposDocumentos(ctx context.Context) ([]domain.TipoDocumento, error) {
	const procedure = `usp_go_ListarTiposDocumentos`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarTiposDocumentos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading tipos documentos: %w", err)
	}

	items := make([]domain.TipoDocumento, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.TipoDocumento{
			ID:          *rowInt64(m, "IdDocIdentidad"),
			Descripcion: rowString(m, "Descripcion"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListDepartamentos(ctx context.Context) ([]domain.Departamento, error) {
	const procedure = `usp_go_ListarDepartamentos`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarDepartamentos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading departamentos: %w", err)
	}

	items := make([]domain.Departamento, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.Departamento{
			ID:     *rowInt64(m, "IdDepartamento"),
			Nombre: rowString(m, "Nombre"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListProvincias(ctx context.Context, idDepartamento int64) ([]domain.Provincia, error) {
	const procedure = `usp_go_ListarProvincias`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdDepartamento", idDepartamento))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarProvincias: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading provincias: %w", err)
	}

	items := make([]domain.Provincia, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.Provincia{
			ID:     *rowInt64(m, "IdProvincia"),
			Nombre: rowString(m, "Nombre"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListDistritos(ctx context.Context, idProvincia int64) ([]domain.Distrito, error) {
	const procedure = `usp_go_ListarDistritos`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdProvincia", idProvincia))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarDistritos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading distritos: %w", err)
	}

	items := make([]domain.Distrito, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.Distrito{
			ID:     *rowInt64(m, "IdDistrito"),
			Nombre: rowString(m, "Nombre"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListCentrosPoblados(ctx context.Context, idDistrito int64) ([]domain.CentroPoblado, error) {
	const procedure = `usp_go_ListarCentrosPoblados`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdDistrito", idDistrito))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarCentrosPoblados: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading centros poblados: %w", err)
	}

	items := make([]domain.CentroPoblado, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.CentroPoblado{
			ID:     *rowInt64(m, "IdCentroPoblado"),
			Nombre: rowString(m, "Nombre"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListPaises(ctx context.Context) ([]domain.Pais, error) {
	const procedure = `usp_go_ListarPaises`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarPaises: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading paises: %w", err)
	}

	items := make([]domain.Pais, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.Pais{
			ID:     *rowInt64(m, "IdPais"),
			Nombre: rowString(m, "nombre", "Nombre"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListEstadosLlegoPaciente(ctx context.Context) ([]domain.EstadoLlegoPaciente, error) {
	const procedure = `usp_go_listarEstadosLlegoPaciente`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_listarEstadosLlegoPaciente: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading estados llego paciente: %w", err)
	}

	items := make([]domain.EstadoLlegoPaciente, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.EstadoLlegoPaciente{
			ID:          *rowInt64(m, "id"),
			Descripcion: rowString(m, "Descripcion"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListFuentesFinanciamiento(ctx context.Context) ([]domain.FuenteFinanciamiento, error) {
	const procedure = `usp_go_ListarFuentesFinanciamiento`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarFuentesFinanciamiento: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading fuentes financiamiento: %w", err)
	}

	items := make([]domain.FuenteFinanciamiento, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.FuenteFinanciamiento{
			ID:                   *rowInt64(m, "IdFuenteFinanciamiento"),
			Descripcion:          rowString(m, "Descripcion"),
			IdTipoFinanciamiento: *rowInt64(m, "IdTipoFinanciamiento"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListServicios(ctx context.Context, idTipoServicio int64) ([]domain.Servicio, error) {
	const procedure = `usp_go_ListarServicios`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdTipoServicio", idTipoServicio))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarServicios: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading servicios: %w", err)
	}

	items := make([]domain.Servicio, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.Servicio{
			ID:     *rowInt64(m, "IdServicio"),
			Nombre: rowString(m, "Nombre"),
		})
	}

	return items, nil
}

// GetDatosInstitucion invoca el SP webParametrosDatosInstitucion, que
// devuelve una Ãºnica fila con los datos de la instituciÃ³n (RUC, nombre,
// direcciÃ³n, telÃ©fono, logos, etc.). Devuelve nil si el SP no retorna filas.
func (r *catalogRepository) GetDatosInstitucion(ctx context.Context) (*domain.DatosInstitucion, error) {
	const procedure = `usp_go_ParametrosDatosInstitucion`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ParametrosDatosInstitucion: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading datos institucion: %w", err)
	}
	if len(maps) == 0 {
		return nil, nil
	}

	m := maps[0]
	return &domain.DatosInstitucion{
		RucEESS:    rowString(m, "RUC_EESS"),
		Nombre:     rowString(m, "NOMBRE"),
		Direccion:  rowString(m, "DIRECCION"),
		Telefono:   rowString(m, "TELEFONO"),
		Codigo:     rowString(m, "CODIGO"),
		CodRenaes:  rowString(m, "COD_RENAES"),
		LogoHospi:  rowString(m, "LOGO_HOSPI"),
		LogoMinsa:  rowString(m, "LOGO_MINSA"),
		UbigeoHosp: rowString(m, "UBIGEO_HOSP"),
	}, nil
}

func (r *catalogRepository) ListEspecialidades(ctx context.Context) ([]domain.Especialidad, error) {
	const procedure = `usp_go_Especialidades_Listar`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_Especialidades_Listar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading especialidades: %w", err)
	}

	items := make([]domain.Especialidad, 0, len(maps))
	for _, m := range maps {
		idVal := rowInt64(m, "IdEspecialidad", "id", "ID")
		var id int64
		if idVal != nil {
			id = *idVal
		}
		items = append(items, domain.Especialidad{
			ID:     id,
			Nombre: rowString(m, "Nombre", "descripcion", "Descripcion"),
		})
	}

	return items, nil
}

func (r *catalogRepository) ListarEspecialidadesPorDepartamento(ctx context.Context, idDepartamento int) ([]domain.EspecialidadSimple, error) {
	const query = `EXEC usp_go_ListarEspecialidadXDepartamento @IdDepartamento = @p1`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idDepartamento))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarEspecialidadXDepartamento: %w", err)
	}
	defer rows.Close()

	var items []domain.EspecialidadSimple
	for rows.Next() {
		var e domain.EspecialidadSimple
		if err := rows.Scan(&e.IdEspecialidad, &e.Nombre); err != nil {
			return nil, fmt.Errorf("scanning especialidad simple: %w", err)
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *catalogRepository) ListarEspecialidadesQx(ctx context.Context) ([]domain.EspecialidadSimple, error) {
	const query = `EXEC usp_go_ListarEspecialidadesQx`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarEspecialidadesQx: %w", err)
	}
	defer rows.Close()

	var items []domain.EspecialidadSimple
	for rows.Next() {
		var e domain.EspecialidadSimple
		var descLarga sql.NullString
		if err := rows.Scan(&e.IdEspecialidad, &descLarga, &e.Nombre); err != nil {
			return nil, fmt.Errorf("scanning especialidad qx: %w", err)
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

// ListEspecialidadesPorTipoServicio invoca usp_go_EspecialidadesXidTipoServicio.
func (r *catalogRepository) ListEspecialidadesPorTipoServicio(ctx context.Context, idTipoServicio int64) ([]domain.Especialidad, error) {
	const procedure = `usp_go_EspecialidadesXidTipoServicio`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdTipoServicio", idTipoServicio))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_EspecialidadesXidTipoServicio: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading especialidades por tipo servicio: %w", err)
	}

	items := make([]domain.Especialidad, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m,
			"IdEspecialidad", "idEspecialidad",
			"Id", "id",
		)
		if id == nil {
			continue
		}
		nombre := rowString(m,
			"Nombre", "nombre",
			"Descripcion", "descripcion",
			"Especialidad", "especialidad",
		)
		items = append(items, domain.Especialidad{
			ID:     *id,
			Nombre: nombre,
		})
	}

	return items, nil
}

// ListEspecialidadesListar invoca usp_go_Especialidades_Listar (sin parámetros).
// Usado por el combo de interconsulta en Primera Atención.
func (r *catalogRepository) ListEspecialidadesListar(ctx context.Context) ([]domain.Especialidad, error) {
	const procedure = `usp_go_Especialidades_Listar`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_Especialidades_Listar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading especialidades listar: %w", err)
	}

	items := make([]domain.Especialidad, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m,
			"IdEspecialidad", "idEspecialidad",
			"Id", "id",
		)
		if id == nil {
			continue
		}
		nombre := rowString(m,
			"Nombre", "nombre",
			"Descripcion", "descripcion",
			"Especialidad", "especialidad",
		)
		items = append(items, domain.Especialidad{
			ID:     *id,
			Nombre: nombre,
		})
	}

	if len(maps) > 0 && len(items) == 0 {
		cols := make([]string, 0, len(maps[0]))
		for k := range maps[0] {
			cols = append(cols, k)
		}
		return nil, fmt.Errorf(
			"usp_go_Especialidades_Listar devolvió %d fila(s) pero no se pudo mapear IdEspecialidad; columnas=%v",
			len(maps), cols,
		)
	}

	return items, nil
}

func (r *catalogRepository) ListarServiciosPorPrioridad(ctx context.Context, idPrioridad int, fechaNac string) ([]domain.ServicioSimple, error) {
	const query = `EXEC usp_go_ListarServiciosXPrioridad @IdPrioridad = @p1, @FechaNac = @p2`

	// Sin fecha de nacimiento se envía NULL para que el SP aplique su
	// comportamiento por defecto (lista general según la prioridad).
	var fechaNacParam any
	if fechaNac != "" {
		fechaNacParam = fechaNac
	}

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPrioridad), sql.Named("p2", fechaNacParam))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarServiciosXPrioridad: %w", err)
	}
	defer rows.Close()

	var items []domain.ServicioSimple
	for rows.Next() {
		var s domain.ServicioSimple
		if err := rows.Scan(&s.IdServicio, &s.Nombre); err != nil {
			return nil, fmt.Errorf("scanning servicio por prioridad: %w", err)
		}
		items = append(items, s)
	}
	return items, rows.Err()
}

// GetParametro invoca el SP usp_go_webParametroSeleccionarPorId, que
// devuelve en una única fila los valores (Tipo, Codigo, ValorTexto,
// ValorInt, ValorFloat) del parámetro solicitado. Devuelve nil si el SP
// no retorna filas.
func (r *catalogRepository) GetParametro(ctx context.Context, idParametro int64) (*domain.Parametro, error) {
	const procedure = `usp_go_webParametroSeleccionarPorId`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdParametro", idParametro))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_webParametroSeleccionarPorId: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading parametro: %w", err)
	}
	if len(maps) == 0 {
		return nil, nil
	}

	m := maps[0]
	return &domain.Parametro{
		Tipo:       rowString(m, "Tipo"),
		Codigo:     rowString(m, "Codigo"),
		ValorTexto: rowString(m, "ValorTexto"),
		ValorInt:   rowInt64(m, "ValorInt"),
		ValorFloat: rowFloat64(m, "ValorFloat"),
	}, nil
}

// ListParametrosClinicos invoca usp_go_Cat_ParametroClinico_Listar @IdGrupo.
// Identidad de fila SQL: IdParametro -> domain.ParametroClinico.ID.
func (r *catalogRepository) ListParametrosClinicos(ctx context.Context, idGrupo int64) ([]domain.ParametroClinico, error) {
	const procedure = `usp_go_Cat_ParametroClinico_Listar`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("IdGrupo", idGrupo))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_Cat_ParametroClinico_Listar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading parametros clinicos: %w", err)
	}

	items := make([]domain.ParametroClinico, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m, "IdParametro", "idParametro", "Id", "id")
		if id == nil {
			continue
		}
		desc := rowString(m,
			"Descripcion", "descripcion",
			"Nombre", "nombre",
			"Valor", "valor",
		)
		codigo := rowString(m, "Codigo", "codigo", "CODIGO")
		items = append(items, domain.ParametroClinico{
			ID:          *id,
			Descripcion: desc,
			Codigo:      codigo,
		})
	}

	return items, nil
}

// ListPrioridadesPrimeraAtencion invoca usp_go_Cat_Prioridad_PrimeraAtencion.
// Identidad de fila SQL: IdTipoGravedad (alias histórico IdTipoPrioridad en bandeja).
func (r *catalogRepository) ListPrioridadesPrimeraAtencion(ctx context.Context) ([]domain.PrioridadPrimeraAtencion, error) {
	const procedure = `usp_go_Cat_Prioridad_PrimeraAtencion`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_Cat_Prioridad_PrimeraAtencion: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading prioridades primera atencion: %w", err)
	}

	items := make([]domain.PrioridadPrimeraAtencion, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m,
			"IdTipoGravedad", "idTipoGravedad",
			"IdTipoPrioridad", "idTipoPrioridad",
			"IdTiposPrioridad", "idTiposPrioridad",
			"IdPrioridad", "idPrioridad",
			"Codigo", "codigo",
			"Id", "id",
		)
		if id == nil {
			continue
		}
		desc := rowString(m,
			"Descripcion", "descripcion",
			"TipoPrioridad", "tipoPrioridad",
			"DescripcionPrioridad", "descripcionPrioridad",
			"Prioridad", "prioridad",
			"Nombre", "nombre",
			"Label", "label",
		)
		items = append(items, domain.PrioridadPrimeraAtencion{
			ID:          *id,
			Descripcion: desc,
		})
	}

	if len(maps) > 0 && len(items) == 0 {
		cols := make([]string, 0, len(maps[0]))
		for k := range maps[0] {
			cols = append(cols, k)
		}
		return nil, fmt.Errorf(
			"usp_go_Cat_Prioridad_PrimeraAtencion devolvió %d fila(s) pero no se pudo mapear el id; columnas=%v",
			len(maps), cols,
		)
	}

	return items, nil
}

// ListSubclasificacionesDiagnostico invoca usp_go_SubclasificacionDiagnosticos_Select
// (@IdTipoServicio, @IdClasificacionDx). 9=PA, 7=IC, 3=egreso, 4=fallecido.
func (r *catalogRepository) ListSubclasificacionesDiagnostico(ctx context.Context, idTipoServicio, idClasificacionDx int64) ([]domain.SubclasificacionDiagnostico, error) {
	const procedure = `usp_go_SubclasificacionDiagnosticos_Select`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("IdTipoServicio", idTipoServicio),
		sql.Named("IdClasificacionDx", idClasificacionDx),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_SubclasificacionDiagnosticos_Select: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading subclasificaciones diagnostico: %w", err)
	}

	items := make([]domain.SubclasificacionDiagnostico, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m,
			"IdSubclasificacionDx", "idSubclasificacionDx",
			"IdSubclasificacionDiagnostico", "idSubclasificacionDiagnostico",
			"Id", "id",
		)
		if id == nil {
			continue
		}
		corta := rowString(m, "DescripcionCorta", "descripcionCorta", "Descripcion", "descripcion")
		larga := rowString(m, "DescripcionLarga", "descripcionLarga")
		items = append(items, domain.SubclasificacionDiagnostico{
			ID:               *id,
			DescripcionCorta: corta,
			DescripcionLarga: larga,
		})
	}

	if len(maps) > 0 && len(items) == 0 {
		cols := make([]string, 0, len(maps[0]))
		for k := range maps[0] {
			cols = append(cols, k)
		}
		return nil, fmt.Errorf(
			"usp_go_SubclasificacionDiagnosticos_Select devolvió %d fila(s) sin IdSubclasificacionDx; columnas=%v",
			len(maps), cols,
		)
	}

	return items, nil
}

// ListDiagnosticosCatalogo invoca usp_go_SelectDiagnosticos.
func (r *catalogRepository) ListDiagnosticosCatalogo(ctx context.Context, filtro string, idAtencion, idPaciente int64) ([]domain.DiagnosticoCatalogo, error) {
	const procedure = `usp_go_SelectDiagnosticos`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("filtro", filtro),
		sql.Named("IdAtencion", idAtencion),
		sql.Named("IdPaciente", idPaciente),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_SelectDiagnosticos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading diagnosticos catalogo: %w", err)
	}

	items := make([]domain.DiagnosticoCatalogo, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m, "IdDiagnostico", "idDiagnostico", "Id", "id")
		if id == nil {
			continue
		}
		desc := rowString(m, "Descripcion", "descripcion")
		cie := rowString(m, "CodigoCIE10", "codigoCIE10", "CodigoCie10")
		item := domain.DiagnosticoCatalogo{ID: *id}
		if desc != nil {
			item.Descripcion = *desc
		}
		if cie != nil {
			item.CodigoCIE10 = strings.TrimSpace(*cie)
		}
		items = append(items, item)
	}

	return items, nil
}

// ListMedicamentosCatalogo invoca usp_go_SelectMedicamentosFiltro.
// Misma forma de llamada/lectura que BuscarMedicamentosReceta y orden_repository:
// EXEC explícito + Scan tipado. No usar Query(nombreProc)+rowsToMaps aquí:
// con go-mssqldb eso duplicaba filas y corrompía money (Precio → 0 / 2× ítem).
func (r *catalogRepository) ListMedicamentosCatalogo(ctx context.Context, filtro string, idPaciente int64) ([]domain.MedicamentoCatalogo, error) {
	const query = `EXEC usp_go_SelectMedicamentosFiltro @Filtro = @Filtro, @IdPaciente = @IdPaciente`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("Filtro", filtro),
		sql.Named("IdPaciente", idPaciente),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_SelectMedicamentosFiltro: %w", err)
	}
	defer rows.Close()

	items := make([]domain.MedicamentoCatalogo, 0)
	for rows.Next() {
		var (
			idProd                                                 sql.NullInt64
			cod, nom, nomLargo, formaFarm                          sql.NullString
			stock, tipoProd, ultCant, ultDosis, ultUnid, ultFrec   sql.NullInt64
			ultVia, ultDur, tieneAnt, cargaFua                     sql.NullInt64
			precio                                                 sql.NullFloat64
			ultFecha                                               sql.NullTime
		)
		// Orden de columnas = SELECT final del SP.
		if err := rows.Scan(
			&idProd, &cod, &nom, &nomLargo,
			&formaFarm, &stock, &precio, &tipoProd,
			&ultFecha, &ultCant, &ultDosis,
			&ultUnid, &ultFrec, &ultVia,
			&ultDur, &tieneAnt, &cargaFua,
		); err != nil {
			return nil, fmt.Errorf("escaneando medicamento catalogo: %w", err)
		}
		if !idProd.Valid || idProd.Int64 <= 0 {
			continue
		}
		item := domain.MedicamentoCatalogo{IDProducto: idProd.Int64}
		if cod.Valid {
			item.Codigo = strings.TrimSpace(cod.String)
		}
		if nom.Valid {
			item.Nombre = strings.TrimSpace(nom.String)
		}
		if precio.Valid {
			item.PrecioUnitario = precio.Float64
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterando medicamentos catalogo: %w", err)
	}
	return items, nil
}

func (r *catalogRepository) BuscarMedicamentosReceta(ctx context.Context, filtro string, idPaciente int) ([]domain.MedicamentoBusqueda, error) {
	rows, err := r.db.QueryContext(ctx, "EXEC usp_go_SelectMedicamentosFiltro @Filtro = @p1, @IdPaciente = @p2", sql.Named("p1", filtro), sql.Named("p2", idPaciente))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_SelectMedicamentosFiltro: %w", err)
	}
	defer rows.Close()

	var items []domain.MedicamentoBusqueda
	for rows.Next() {
		var m domain.MedicamentoBusqueda
		var (
			cod, nom, nomLargo, formaFarm                                sql.NullString
			stock, tipoProd, ultCant, ultDosis, ultUnid, ultFrec, ultVia sql.NullInt64
			ultDur, tieneAnt, cargaFua, idProd                           sql.NullInt64
			precio                                                       sql.NullFloat64
			ultFecha                                                     sql.NullTime
		)

		if err := rows.Scan(
			&idProd, &cod, &nom, &nomLargo,
			&formaFarm, &stock, &precio, &tipoProd,
			&ultFecha, &ultCant, &ultDosis,
			&ultUnid, &ultFrec, &ultVia,
			&ultDur, &tieneAnt, &cargaFua,
		); err != nil {
			return nil, fmt.Errorf("escaneando medicamento: %w", err)
		}
		if idProd.Valid {
			m.IdProducto = int(idProd.Int64)
		}
		if cod.Valid {
			m.Codigo = cod.String
		}
		if nom.Valid {
			m.Nombre = nom.String
		}
		if stock.Valid {
			m.Stock = int(stock.Int64)
		}
		if precio.Valid {
			m.Precio = precio.Float64
		}
		if ultDosis.Valid {
			m.IdDosisRecetada = int(ultDosis.Int64)
		}
		if ultUnid.Valid {
			m.IdUNIDDosisReceta = int(ultUnid.Int64)
		}
		if ultFrec.Valid {
			m.IdFrecuencia = int(ultFrec.Int64)
		}
		if ultVia.Valid {
			m.IdViaAdministracion = int(ultVia.Int64)
		}
		if tieneAnt.Valid {
			m.TieneRecetaAnterior = int(tieneAnt.Int64)
		}

		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// ListExamenesCatalogo invoca webProcedimentosRecetasFiltarNombrePuntoCarga
// (@IdTipoFinanciamiento ← IdFormaPago bandeja, @Filtro).
func (r *catalogRepository) ListExamenesCatalogo(ctx context.Context, idTipoFinanciamiento int64, filtro string) ([]domain.ExamenCatalogo, error) {
	const procedure = `webProcedimentosRecetasFiltarNombrePuntoCarga`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("IdTipoFinanciamiento", idTipoFinanciamiento),
		sql.Named("Filtro", filtro),
	)
	if err != nil {
		return nil, fmt.Errorf("calling webProcedimentosRecetasFiltarNombrePuntoCarga: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading examenes catalogo: %w", err)
	}

	items := make([]domain.ExamenCatalogo, 0, len(maps))
	for _, m := range maps {
		// IdProducto del SP es compuesto ("54707|2"); la PK útil es IdProductoSolo.
		id := rowInt64(m, "IdProductoSolo", "idProductoSolo")
		if id == nil {
			continue
		}
		item := domain.ExamenCatalogo{IDProducto: *id}
		if c := rowString(m, "Codigo", "codigo"); c != nil {
			item.Codigo = strings.TrimSpace(*c)
		}
		if n := rowString(m, "Nombre", "nombre"); n != nil {
			item.Nombre = strings.TrimSpace(*n)
		}
		if nl := rowString(m, "NombreLargo", "nombreLargo"); nl != nil {
			item.NombreLargo = strings.TrimSpace(*nl)
		}
		if pc := rowInt64(m, "idPuntoCarga", "IdPuntoCarga"); pc != nil {
			item.IDPuntoCarga = *pc
		}
		if pn := rowString(m, "PuntoCarga", "puntoCarga"); pn != nil {
			item.PuntoCarga = strings.TrimSpace(*pn)
		}
		if s := rowString(m, "Servicio", "servicio"); s != nil {
			item.Servicio = strings.TrimSpace(*s)
		}
		if p := rowFloat64(m, "PrecioUnitario", "precioUnitario"); p != nil {
			item.PrecioUnitario = *p
		}
		items = append(items, item)
	}
	return items, nil
}

// ListProcedimientosCatalogo invoca usp_go_ProcedimentosFiltarNombrePuntoCarga.
func (r *catalogRepository) ListProcedimientosCatalogo(ctx context.Context, idPuntoCarga, idTipoFinanciamiento int64, filtro string) ([]domain.ProcedimientoCatalogo, error) {
	const procedure = `usp_go_ProcedimentosFiltarNombrePuntoCarga`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("IdPuntoCarga", idPuntoCarga),
		sql.Named("IdTipoFinanciamiento", idTipoFinanciamiento),
		sql.Named("Filtro", filtro),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ProcedimentosFiltarNombrePuntoCarga: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading procedimientos catalogo: %w", err)
	}

	items := make([]domain.ProcedimientoCatalogo, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m, "IdProducto", "idProducto")
		if id == nil {
			continue
		}
		item := domain.ProcedimientoCatalogo{IDProducto: *id}
		if c := rowString(m, "codigo", "Codigo"); c != nil {
			item.Codigo = strings.TrimSpace(*c)
		}
		if n := rowString(m, "Nombre", "nombre"); n != nil {
			item.Nombre = strings.TrimSpace(*n)
		}
		if nl := rowString(m, "NombreLargo", "nombreLargo"); nl != nil {
			item.NombreLargo = strings.TrimSpace(*nl)
		}
		if pc := rowInt64(m, "IdPuntoCarga", "idPuntoCarga"); pc != nil {
			item.IDPuntoCarga = *pc
		}
		if pn := rowString(m, "PuntoCarga", "puntoCarga"); pn != nil {
			item.PuntoCarga = strings.TrimSpace(*pn)
		}
		if p := rowFloat64(m, "PrecioUnitario", "precioUnitario"); p != nil {
			item.PrecioUnitario = *p
		}
		if b := rowBool(m, "EsCPT", "esCPT"); b != nil {
			item.EsCPT = *b
		}
		items = append(items, item)
	}
	return items, nil
}

// ActualizarParametroValorTexto invoca el SP usp_go_ActualizarParametroValorTextoPorId,
// que deja ValorTexto = 'N' para el parámetro recibido.
func (r *catalogRepository) ActualizarParametroValorTexto(ctx context.Context, idParametro int) error {
	const procedure = `EXEC dbo.usp_go_ActualizarParametroValorTextoPorId @IdParametro = @p1`
	if _, err := r.db.ExecContext(ctx, procedure, sql.Named("p1", idParametro)); err != nil {
		return fmt.Errorf("calling usp_go_ActualizarParametroValorTextoPorId: %w", err)
	}
	return nil
}

// ListEstablecimientos invoca el SP usp_go_webEstablecimientosFiltrar para
// buscar establecimientos por filtro y tipo (1=referencia, 0=contrareferencia).
func (r *catalogRepository) ListEstablecimientos(ctx context.Context, filtro string, tipo string) ([]domain.EstablecimientoBusqueda, error) {
	rows, err := r.db.QueryContext(ctx, "EXEC [dbo].[usp_go_webEstablecimientosFiltrar] @filtro = @p1, @Tipo = @p2", sql.Named("p1", filtro), sql.Named("p2", tipo))
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_webEstablecimientosFiltrar: %w", err)
	}
	defer rows.Close()

	var items []domain.EstablecimientoBusqueda
	for rows.Next() {
		var (
			item            domain.EstablecimientoBusqueda
			dist, prov, dep sql.NullString
		)
		if err := rows.Scan(&item.IdEstablecimiento, &item.Codigo, &item.Nombre, &dist, &prov, &dep, &item.NombreLargo); err != nil {
			return nil, fmt.Errorf("escaneando establecimiento: %w", err)
		}
		item.Distrito = dist.String
		item.Provincia = prov.String
		item.Departamento = dep.String
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// ListRecetaDosis invoca usp_go_RecetaDosisSelecionarTodos.
func (r *catalogRepository) ListRecetaDosis(ctx context.Context) ([]domain.RecetaDosis, error) {
	const procedure = `usp_go_RecetaDosisSelecionarTodos`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_RecetaDosisSelecionarTodos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading receta dosis: %w", err)
	}

	items := make([]domain.RecetaDosis, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m, "idDosis", "IdDosis")
		if id == nil {
			continue
		}
		item := domain.RecetaDosis{ID: *id}
		if s := rowString(m, "NumeroDosis", "numeroDosis"); s != nil {
			item.NumeroDosis = strings.TrimSpace(*s)
		} else if f := rowFloat64(m, "NumeroDosis", "numeroDosis"); f != nil {
			item.NumeroDosis = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", *f), "0"), ".")
		}
		items = append(items, item)
	}
	return items, nil
}

// ListRecetaUnidadesDosis invoca Usp_go_SelectRecetaUndDosisSelecionarTodos.
func (r *catalogRepository) ListRecetaUnidadesDosis(ctx context.Context) ([]domain.RecetaUnidadDosis, error) {
	const procedure = `Usp_go_SelectRecetaUndDosisSelecionarTodos`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling Usp_go_SelectRecetaUndDosisSelecionarTodos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading receta unidades dosis: %w", err)
	}

	items := make([]domain.RecetaUnidadDosis, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m, "idUNIDDosis", "IdUNIDDosis", "idUnidDosis")
		if id == nil {
			continue
		}
		item := domain.RecetaUnidadDosis{ID: *id}
		if s := rowString(m, "NumeroUNIDDosis", "numeroUNIDDosis", "NumeroUnidDosis"); s != nil {
			item.NumeroUNIDDosis = strings.TrimSpace(*s)
		}
		items = append(items, item)
	}
	return items, nil
}

// ListRecetaFrecuencias invoca usp_go_SelectRecetaFrecuenciaSelecionarTodos.
func (r *catalogRepository) ListRecetaFrecuencias(ctx context.Context) ([]domain.RecetaFrecuencia, error) {
	const procedure = `usp_go_SelectRecetaFrecuenciaSelecionarTodos`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_SelectRecetaFrecuenciaSelecionarTodos: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading receta frecuencias: %w", err)
	}

	items := make([]domain.RecetaFrecuencia, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m, "IdfRecuencia", "IdFrecuencia", "idFrecuencia")
		if id == nil {
			continue
		}
		item := domain.RecetaFrecuencia{ID: *id}
		if s := rowString(m, "Descripcion", "descripcion"); s != nil {
			item.Descripcion = strings.TrimSpace(*s)
		}
		items = append(items, item)
	}
	return items, nil
}

// ListViasAdministracion invoca usp_go_RecetasListadoViasAdministracion.
func (r *catalogRepository) ListViasAdministracion(ctx context.Context) ([]domain.ViaAdministracion, error) {
	const procedure = `usp_go_RecetasListadoViasAdministracion`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_RecetasListadoViasAdministracion: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading vias administracion: %w", err)
	}

	items := make([]domain.ViaAdministracion, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m, "IdViaAdministracion", "idViaAdministracion")
		if id == nil {
			continue
		}
		item := domain.ViaAdministracion{ID: *id}
		if s := rowString(m, "Descripcion", "descripcion"); s != nil {
			item.Descripcion = strings.TrimSpace(*s)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *catalogRepository) ListSubclasificacionesDiagnosticosItems(ctx context.Context, idTipoServicio, idClasificacionDx int) ([]domain.SubclasificacionDxItem, error) {
	rows, err := r.db.QueryContext(ctx, "EXEC dbo.usp_go_SubclasificacionDiagnosticos_Select @IdTipoServicio = @p1, @IdClasificacionDx = @p2",
		sql.Named("p1", idTipoServicio),
		sql.Named("p2", idClasificacionDx),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_SubclasificacionDiagnosticos_Select: %w", err)
	}
	defer rows.Close()

	var items []domain.SubclasificacionDxItem
	for rows.Next() {
		var (
			idSubDx int
			corta   string
			larga   string
		)
		if err := rows.Scan(&idSubDx, &corta, &larga); err != nil {
			return nil, fmt.Errorf("escaneando subclasificacion dx: %w", err)
		}
		items = append(items, domain.SubclasificacionDxItem{
			IdSubclasificacionDx: idSubDx,
			DescripcionCorta:     corta,
			DescripcionLarga:     larga,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

