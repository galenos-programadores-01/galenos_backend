package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type catalogRepository struct {
	db *sql.DB
}

// NewCatalogRepository construye el adaptador de catálogos de datos
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
// devuelve una única fila con los datos de la institución (RUC, nombre,
// dirección, teléfono, logos, etc.). Devuelve nil si el SP no retorna filas.
func (r *catalogRepository) GetDatosInstitucion(ctx context.Context) (*domain.DatosInstitucion, error) {
	const procedure = `webParametrosDatosInstitucion`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling webParametrosDatosInstitucion: %w", err)
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
	const procedure = `usp_go_ListarEspecialidades`

	rows, err := r.db.QueryContext(ctx, procedure)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarEspecialidades: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading especialidades: %w", err)
	}

	items := make([]domain.Especialidad, 0, len(maps))
	for _, m := range maps {
		items = append(items, domain.Especialidad{
			ID:     *rowInt64(m, "id"),
			Nombre: rowString(m, "descripcion"),
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
