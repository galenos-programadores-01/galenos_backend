package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sqlServerServiciosCatalogoRepository struct {
	db *sql.DB
}

func NewSqlServerServiciosCatalogoRepository(db *sql.DB) output.ServiciosCatalogoRepository {
	return &sqlServerServiciosCatalogoRepository{db: db}
}

func (r *sqlServerServiciosCatalogoRepository) Listar(ctx context.Context, filtros domain.ServicioCatalogoFiltros) (domain.PaginacionServiciosCatalogo, error) {
	var resultado domain.PaginacionServiciosCatalogo
	resultado.Page = filtros.Page
	resultado.PageSize = filtros.PageSize
	resultado.Items = make([]domain.ServicioCatalogoItem, 0)

	filtroLimpio := strings.TrimSpace(filtros.Filtro)
	filtroParam := sql.NullString{String: filtroLimpio, Valid: filtroLimpio != ""}

	var soloCptParam sql.NullBool
	if filtros.SoloCpt != nil {
		soloCptParam = sql.NullBool{Bool: *filtros.SoloCpt, Valid: true}
	}

	var soloQxParam sql.NullBool
	if filtros.SoloQuirurgico != nil {
		soloQxParam = sql.NullBool{Bool: *filtros.SoloQuirurgico, Valid: true}
	}

	var estadoParam sql.NullInt32
	if filtros.Estado != nil {
		estadoParam = sql.NullInt32{Int32: int32(*filtros.Estado), Valid: true}
	}

	query := `EXEC dbo.usp_go_FactCatalogoServicios_ListarPaginado
		@Filtro = @p1,
		@IdServicioGrupo = @p2,
		@EsCPT = @p3,
		@EsQuirurgico = @p4,
		@IdEstado = @p5,
		@PageNumber = @p6,
		@PageSize = @p7`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", filtroParam),
		sql.Named("p2", filtros.IdServicioGrupo),
		sql.Named("p3", soloCptParam),
		sql.Named("p4", soloQxParam),
		sql.Named("p5", estadoParam),
		sql.Named("p6", filtros.Page),
		sql.Named("p7", filtros.PageSize),
	)
	if err != nil {
		return resultado, fmt.Errorf("error consultando servicios de catalogo mediante SP: %w", err)
	}
	defer rows.Close()

	totalRegistros := 0
	for rows.Next() {
		var item domain.ServicioCatalogoItem
		var regTotal int

		err := rows.Scan(
			&item.IdProducto,
			&item.Codigo,
			&item.CodigoHis,
			&item.Nombre,
			&item.NombreCorto,
			&item.IdServicioGrupo,
			&item.GrupoNombre,
			&item.CodMinsa,
			&item.NombreMinsa,
			&item.EsCpt,
			&item.EsQuirurgico,
			&item.CodigoSis,
			&item.CodigoSoat,
			&item.Cpms,
			&item.PrecioParticular,
			&item.PrecioSis,
			&item.IdEstado,
			&item.CargaEnFua,
			&item.CargaEnHis,
			&regTotal,
		)
		if err != nil {
			return resultado, fmt.Errorf("error escaneando servicio catalogo: %w", err)
		}

		totalRegistros = regTotal
		item.Codigo = strings.TrimSpace(item.Codigo)
		item.CodigoHis = strings.TrimSpace(item.CodigoHis)
		item.Nombre = strings.TrimSpace(item.Nombre)
		item.NombreCorto = strings.TrimSpace(item.NombreCorto)
		item.GrupoNombre = strings.TrimSpace(item.GrupoNombre)
		item.CodMinsa = strings.TrimSpace(item.CodMinsa)
		item.NombreMinsa = strings.TrimSpace(item.NombreMinsa)
		item.CodigoSis = strings.TrimSpace(item.CodigoSis)
		item.CodigoSoat = strings.TrimSpace(item.CodigoSoat)
		item.Cpms = strings.TrimSpace(item.Cpms)

		resultado.Items = append(resultado.Items, item)
	}

	if err := rows.Err(); err != nil {
		return resultado, fmt.Errorf("error iterando servicios: %w", err)
	}

	resultado.TotalItems = totalRegistros
	if filtros.PageSize > 0 {
		resultado.TotalPages = int(math.Ceil(float64(totalRegistros) / float64(filtros.PageSize)))
	}

	return resultado, nil
}

func (r *sqlServerServiciosCatalogoRepository) ObtenerPorID(ctx context.Context, idProducto int) (*domain.ServicioCatalogoDetalle, error) {
	query := `EXEC dbo.usp_go_FactCatalogoServicios_ObtenerPorId @IdProducto = @p1`
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idProducto))
	if err != nil {
		return nil, fmt.Errorf("error consultando detalle servicio mediante SP: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var item domain.ServicioCatalogoDetalle
	err = rows.Scan(
		&item.IdProducto,
		&item.Codigo,
		&item.CodigoHis,
		&item.CodigoSunat,
		&item.Nombre,
		&item.NombreCorto,
		&item.IdServicioGrupo,
		&item.GrupoNombre,
		&item.CodMinsa,
		&item.NombreMinsa,
		&item.EsCpt,
		&item.EsQuirurgico,
		&item.CodigoSis,
		&item.CodigoSoat,
		&item.Cpms,
		&item.IdEstado,
		&item.CargaEnFua,
		&item.CargaEnHis,
	)
	if err != nil {
		return nil, fmt.Errorf("error escaneando detalle servicio: %w", err)
	}

	item.Codigo = strings.TrimSpace(item.Codigo)
	item.CodigoHis = strings.TrimSpace(item.CodigoHis)
	item.CodigoSunat = strings.TrimSpace(item.CodigoSunat)
	item.Nombre = strings.TrimSpace(item.Nombre)
	item.NombreCorto = strings.TrimSpace(item.NombreCorto)
	item.GrupoNombre = strings.TrimSpace(item.GrupoNombre)
	item.CodMinsa = strings.TrimSpace(item.CodMinsa)
	item.NombreMinsa = strings.TrimSpace(item.NombreMinsa)
	item.CodigoSis = strings.TrimSpace(item.CodigoSis)
	item.CodigoSoat = strings.TrimSpace(item.CodigoSoat)
	item.Cpms = strings.TrimSpace(item.Cpms)

	item.Tarifas = make([]domain.TarifaFinanciamiento, 0)
	for rows.Next() {
	}
	if rows.NextResultSet() {
		for rows.Next() {
			var t domain.TarifaFinanciamiento
			var activoBit bool
			if errScan := rows.Scan(&t.IdTipoFinanciamiento, &t.FinanciamientoNombre, &t.PrecioUnitario, &activoBit); errScan == nil {
				t.FinanciamientoNombre = strings.TrimSpace(t.FinanciamientoNombre)
				t.Activo = activoBit
				item.Tarifas = append(item.Tarifas, t)
				if t.IdTipoFinanciamiento == 1 {
					item.PrecioParticular = t.PrecioUnitario
				} else if t.IdTipoFinanciamiento == 2 {
					item.PrecioSis = t.PrecioUnitario
				}
			}
		}
	}

	return &item, nil
}

func (r *sqlServerServiciosCatalogoRepository) Crear(ctx context.Context, payload domain.ServicioCatalogoGuardarPayload) (int, error) {
	esCptInt := 0
	if payload.EsCpt {
		esCptInt = 1
	}
	esQxInt := 0
	if payload.EsQuirurgico {
		esQxInt = 1
	}

	var nuevoID int
	querySP := `
		EXEC [dbo].[FactCatalogoServiciosAgregar]
			@IdServicioSubGrupo = 0,
			@IdPartida = 0,
			@IdCentroCosto = 0,
			@IdServicioSubSeccion = 0,
			@IdServicioSeccion = 0,
			@IdServicioGrupo = @p1,
			@Nombre = @p2,
			@Codigo = @p3,
			@CodigoSunat = @p4,
			@IdProducto = @p5 OUTPUT,
			@EsCPT = @p6,
			@NombreMINSA = @p7,
			@idEstado = @p8,
			@codigoSIS = @p9,
			@codigoSOAT = @p10,
			@EsPerfil = 0,
			@Tip_prov_atiende = 'M',
			@EsQuirurgico = @p11,
			@IdUsuarioAuditoria = 1
	`
	_, err := r.db.ExecContext(ctx, querySP,
		sql.Named("p1", payload.IdServicioGrupo),
		sql.Named("p2", payload.Nombre),
		sql.Named("p3", payload.Codigo),
		sql.Named("p4", payload.CodigoSunat),
		sql.Named("p5", sql.Out{Dest: &nuevoID}),
		sql.Named("p6", esCptInt),
		sql.Named("p7", payload.NombreMinsa),
		sql.Named("p8", payload.IdEstado),
		sql.Named("p9", payload.CodigoSis),
		sql.Named("p10", payload.CodigoSoat),
		sql.Named("p11", esQxInt),
	)

	if err != nil {
		return 0, fmt.Errorf("error al ejecutar procedimiento FactCatalogoServiciosAgregar: %w", err)
	}

	guardarTarifaSP := `EXEC dbo.usp_go_FactCatalogoServiciosHosp_GuardarTarifa 
		@IdProducto = @p1, 
		@IdTipoFinanciamiento = @p2, 
		@PrecioUnitario = @p3, 
		@Activo = 1`
	if _, err := r.db.ExecContext(ctx, guardarTarifaSP, sql.Named("p1", nuevoID), sql.Named("p2", 1), sql.Named("p3", payload.PrecioParticular)); err != nil {
		return 0, fmt.Errorf("error guardando tarifa particular: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, guardarTarifaSP, sql.Named("p1", nuevoID), sql.Named("p2", 2), sql.Named("p3", payload.PrecioSis)); err != nil {
		return 0, fmt.Errorf("error guardando tarifa SIS: %w", err)
	}

	if err := r.actualizarNormativa(ctx, nuevoID, payload); err != nil {
		return 0, fmt.Errorf("error actualizando normativa del servicio: %w", err)
	}

	return nuevoID, nil
}

func (r *sqlServerServiciosCatalogoRepository) Actualizar(ctx context.Context, idProducto int, payload domain.ServicioCatalogoGuardarPayload) error {
	esCptInt := 0
	if payload.EsCpt {
		esCptInt = 1
	}
	esQxInt := 0
	if payload.EsQuirurgico {
		esQxInt = 1
	}

	querySP := `
		EXEC [dbo].[FactCatalogoServiciosModificar]
			@IdServicioSubGrupo = 0,
			@IdPartida = 0,
			@IdCentroCosto = 0,
			@IdServicioSubSeccion = 0,
			@IdServicioSeccion = 0,
			@IdServicioGrupo = @p1,
			@Nombre = @p2,
			@Codigo = @p3,
			@CodigoSunat = @p4,
			@IdProducto = @p5,
			@EsCPT = @p6,
			@NombreMINSA = @p7,
			@idEstado = @p8,
			@codigoSIS = @p9,
			@codigoSOAT = @p10,
			@EsPerfil = 0,
			@Tip_prov_atiende = 'M',
			@EsQuirurgico = @p11,
			@IdUsuarioAuditoria = 1
	`

	_, err := r.db.ExecContext(ctx, querySP,
		sql.Named("p1", payload.IdServicioGrupo),
		sql.Named("p2", payload.Nombre),
		sql.Named("p3", payload.Codigo),
		sql.Named("p4", payload.CodigoSunat),
		sql.Named("p5", idProducto),
		sql.Named("p6", esCptInt),
		sql.Named("p7", payload.NombreMinsa),
		sql.Named("p8", payload.IdEstado),
		sql.Named("p9", payload.CodigoSis),
		sql.Named("p10", payload.CodigoSoat),
		sql.Named("p11", esQxInt),
	)

	if err != nil {
		return fmt.Errorf("error al ejecutar procedimiento FactCatalogoServiciosModificar: %w", err)
	}

	guardarTarifaSP := `EXEC dbo.usp_go_FactCatalogoServiciosHosp_GuardarTarifa 
		@IdProducto = @p1, 
		@IdTipoFinanciamiento = @p2, 
		@PrecioUnitario = @p3, 
		@Activo = 1`
	if _, err := r.db.ExecContext(ctx, guardarTarifaSP, sql.Named("p1", idProducto), sql.Named("p2", 1), sql.Named("p3", payload.PrecioParticular)); err != nil {
		return fmt.Errorf("error actualizando tarifa particular: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, guardarTarifaSP, sql.Named("p1", idProducto), sql.Named("p2", 2), sql.Named("p3", payload.PrecioSis)); err != nil {
		return fmt.Errorf("error actualizando tarifa SIS: %w", err)
	}

	if err := r.actualizarNormativa(ctx, idProducto, payload); err != nil {
		return fmt.Errorf("error actualizando normativa del servicio: %w", err)
	}

	return nil
}

func (r *sqlServerServiciosCatalogoRepository) actualizarNormativa(ctx context.Context, idProducto int, payload domain.ServicioCatalogoGuardarPayload) error {
	cptFuaInt := 0
	if payload.CargaEnFua {
		cptFuaInt = 1
	}
	cptHisInt := 0
	if payload.CargaEnHis {
		cptHisInt = 1
	}
	sp := `EXEC dbo.usp_go_FactCatalogoServicios_ActualizarNormativa
		@IdProducto = @p1,
		@NombreCorto = @p2,
		@CodMinsa = @p3,
		@CodigoHis = @p4,
		@Cpms = @p5,
		@CargaEnFua = @p6,
		@CargaEnHis = @p7`
	_, err := r.db.ExecContext(ctx, sp,
		sql.Named("p1", idProducto),
		sql.Named("p2", payload.NombreCorto),
		sql.Named("p3", payload.CodMinsa),
		sql.Named("p4", payload.CodigoHis),
		sql.Named("p5", payload.Cpms),
		sql.Named("p6", cptFuaInt),
		sql.Named("p7", cptHisInt),
	)
	return err
}

func (r *sqlServerServiciosCatalogoRepository) ListarGrupos(ctx context.Context) ([]domain.ServicioGrupo, error) {
	querySP := `EXEC [dbo].[FactCatalogoServiciosGrupoSeleccionarTodos]`
	rows, err := r.db.QueryContext(ctx, querySP)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar FactCatalogoServiciosGrupoSeleccionarTodos: %w", err)
	}
	defer rows.Close()

	grupos := make([]domain.ServicioGrupo, 0)
	for rows.Next() {
		var g domain.ServicioGrupo
		if err := rows.Scan(&g.IdServicioGrupo, &g.Codigo, &g.Descripcion); err != nil {
			return nil, fmt.Errorf("error escaneando grupo servicio: %w", err)
		}
		g.Codigo = strings.TrimSpace(g.Codigo)
		g.Descripcion = strings.TrimSpace(g.Descripcion)
		grupos = append(grupos, g)
	}

	return grupos, nil
}

func (r *sqlServerServiciosCatalogoRepository) Eliminar(ctx context.Context, idProducto int) error {
	query := `EXEC dbo.usp_go_FactCatalogoServicios_Inactivar @IdProducto = @p1, @IdUsuarioAuditoria = 1`
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", idProducto))
	if err != nil {
		return fmt.Errorf("error inactivando servicio %d mediante SP: %w", idProducto, err)
	}
	return nil
}
