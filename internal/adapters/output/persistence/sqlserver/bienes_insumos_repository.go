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

type sqlServerBienesInsumosRepository struct {
	db *sql.DB
}

func NewSqlServerBienesInsumosRepository(db *sql.DB) output.BienesInsumosRepository {
	return &sqlServerBienesInsumosRepository{db: db}
}

func (r *sqlServerBienesInsumosRepository) Listar(ctx context.Context, filtros domain.BienInsumoFiltros) (domain.PaginacionBienesInsumos, error) {
	var resultado domain.PaginacionBienesInsumos
	resultado.Page = filtros.Page
	resultado.PageSize = filtros.PageSize
	resultado.Items = make([]domain.BienInsumoItem, 0)

	filtroLimpio := strings.TrimSpace(filtros.Filtro)
	filtroParam := sql.NullString{String: filtroLimpio, Valid: filtroLimpio != ""}

	var soloPetitorioParam sql.NullBool
	if filtros.SoloPetitorio != nil {
		soloPetitorioParam = sql.NullBool{Bool: *filtros.SoloPetitorio, Valid: true}
	}

	var soloIndicadorParam sql.NullBool
	if filtros.SoloIndicadorGestion != nil {
		soloIndicadorParam = sql.NullBool{Bool: *filtros.SoloIndicadorGestion, Valid: true}
	}

	var estadoParam sql.NullInt32
	if filtros.Estado != nil {
		estadoParam = sql.NullInt32{Int32: int32(*filtros.Estado), Valid: true}
	}

	query := `EXEC dbo.usp_go_FactCatalogoBienesInsumos_ListarPaginado
		@Filtro = @p1,
		@IdGrupoFarmacologico = @p2,
		@TipoProducto = @p3,
		@SoloPetitorio = @p4,
		@SoloIndicadorGestion = @p5,
		@Estado = @p6,
		@PageNumber = @p7,
		@PageSize = @p8`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", filtroParam),
		sql.Named("p2", filtros.IdGrupoFarmacologico),
		sql.Named("p3", filtros.TipoProducto),
		sql.Named("p4", soloPetitorioParam),
		sql.Named("p5", soloIndicadorParam),
		sql.Named("p6", estadoParam),
		sql.Named("p7", filtros.Page),
		sql.Named("p8", filtros.PageSize),
	)
	if err != nil {
		return resultado, fmt.Errorf("error consultando bienes e insumos mediante SP: %w", err)
	}
	defer rows.Close()

	totalRegistros := 0
	for rows.Next() {
		var item domain.BienInsumoItem
		var regTotal int

		err := rows.Scan(
			&item.IdProducto,
			&item.Codigo,
			&item.CodigoSiga,
			&item.CodigoSunat,
			&item.Nombre,
			&item.NombreComercial,
			&item.IdGrupoFarmacologico,
			&item.GrupoFarmacologico,
			&item.PrecioCompra,
			&item.PrecioDistribucion,
			&item.PrecioUltimaCompra,
			&item.StockMinimo,
			&item.TipoProducto,
			&item.Concentracion,
			&item.Presentacion,
			&item.FormaFarmaceutica,
			&item.Petitorio,
			&item.Estado,
			&item.EsVital,
			&item.EsEstrategico,
			&item.EsOncologico,
			&item.CargaEnFua,
			&item.EsIndicadorGestion,
			&regTotal,
		)
		if err != nil {
			return resultado, fmt.Errorf("error escaneando fila de bien/insumo: %w", err)
		}

		totalRegistros = regTotal
		item.Codigo = strings.TrimSpace(item.Codigo)
		item.CodigoSiga = strings.TrimSpace(item.CodigoSiga)
		item.CodigoSunat = strings.TrimSpace(item.CodigoSunat)
		item.Nombre = strings.TrimSpace(item.Nombre)
		item.NombreComercial = strings.TrimSpace(item.NombreComercial)
		item.GrupoFarmacologico = strings.TrimSpace(item.GrupoFarmacologico)
		item.Concentracion = strings.TrimSpace(item.Concentracion)
		item.Presentacion = strings.TrimSpace(item.Presentacion)
		item.FormaFarmaceutica = strings.TrimSpace(item.FormaFarmaceutica)

		resultado.Items = append(resultado.Items, item)
	}

	if err := rows.Err(); err != nil {
		return resultado, fmt.Errorf("error recorriendo filas de bienes: %w", err)
	}

	resultado.TotalItems = totalRegistros
	if filtros.PageSize > 0 {
		resultado.TotalPages = int(math.Ceil(float64(totalRegistros) / float64(filtros.PageSize)))
	}

	return resultado, nil
}

func (r *sqlServerBienesInsumosRepository) ObtenerPorID(ctx context.Context, idProducto int) (*domain.BienInsumoDetalle, error) {
	query := `EXEC dbo.usp_go_FactCatalogoBienesInsumos_ObtenerPorId @IdProducto = @p1`

	row := r.db.QueryRowContext(ctx, query, sql.Named("p1", idProducto))

	var item domain.BienInsumoDetalle
	err := row.Scan(
		&item.IdProducto,
		&item.Codigo,
		&item.CodigoSiga,
		&item.CodigoSunat,
		&item.Nombre,
		&item.NombreComercial,
		&item.IdGrupoFarmacologico,
		&item.GrupoFarmacologico,
		&item.IdSubGrupoFarmacologico,
		&item.SubGrupoFarmacologico,
		&item.IdPartida,
		&item.Partida,
		&item.IdCentroCosto,
		&item.CentroCosto,
		&item.PrecioCompra,
		&item.PrecioDistribucion,
		&item.PrecioDonacion,
		&item.PrecioUltimaCompra,
		&item.StockMinimo,
		&item.TipoProducto,
		&item.Denominacion,
		&item.Concentracion,
		&item.Presentacion,
		&item.FormaFarmaceutica,
		&item.MaterialEnvase,
		&item.PresentacionEnvase,
		&item.Fabricante,
		&item.IdPaisOrigen,
		&item.PaisOrigen,
		&item.Petitorio,
		&item.TipoProductoSismed,
		&item.Estado,
		&item.EstadoDigemid,
		&item.IdTipoPrefijo,
		&item.TipoPrefijo,
		&item.PoliticaCanje,
		&item.PetitorioIst,
		&item.Comunes,
		&item.ComiteFarmacologico,
		&item.Especializados,
		&item.EsVital,
		&item.EsEstrategico,
		&item.EsOncologico,
		&item.CargaEnFua,
		&item.ProductoIndicadorGestion,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error consultando detalle de producto mediante SP: %w", err)
	}

	item.Codigo = strings.TrimSpace(item.Codigo)
	item.CodigoSiga = strings.TrimSpace(item.CodigoSiga)
	item.CodigoSunat = strings.TrimSpace(item.CodigoSunat)
	item.Nombre = strings.TrimSpace(item.Nombre)
	item.NombreComercial = strings.TrimSpace(item.NombreComercial)
	item.GrupoFarmacologico = strings.TrimSpace(item.GrupoFarmacologico)
	item.SubGrupoFarmacologico = strings.TrimSpace(item.SubGrupoFarmacologico)
	item.Partida = strings.TrimSpace(item.Partida)
	item.CentroCosto = strings.TrimSpace(item.CentroCosto)
	item.Denominacion = strings.TrimSpace(item.Denominacion)
	item.Concentracion = strings.TrimSpace(item.Concentracion)
	item.Presentacion = strings.TrimSpace(item.Presentacion)
	item.FormaFarmaceutica = strings.TrimSpace(item.FormaFarmaceutica)
	item.MaterialEnvase = strings.TrimSpace(item.MaterialEnvase)
	item.PresentacionEnvase = strings.TrimSpace(item.PresentacionEnvase)
	item.Fabricante = strings.TrimSpace(item.Fabricante)
	item.PaisOrigen = strings.TrimSpace(item.PaisOrigen)
	item.TipoProductoSismed = strings.TrimSpace(item.TipoProductoSismed)
	item.TipoPrefijo = strings.TrimSpace(item.TipoPrefijo)

	precios, err := r.ListarPreciosPorPlan(ctx, idProducto)
	if err == nil {
		item.PreciosPorPlan = precios
	} else {
		item.PreciosPorPlan = make([]domain.BienInsumoPrecioPlan, 0)
	}

	return &item, nil
}

func (r *sqlServerBienesInsumosRepository) ListarPreciosPorPlan(ctx context.Context, idProducto int) ([]domain.BienInsumoPrecioPlan, error) {
	planesQuery := `EXEC dbo.TiposFinanciamientoSeleccionarTodos`
	rowsPlanes, err := r.db.QueryContext(ctx, planesQuery)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo planes de financiamiento: %w", err)
	}
	defer rowsPlanes.Close()

	nombresPlanes := make(map[int]string)
	idsOrdenados := make([]int, 0)
	for rowsPlanes.Next() {
		var idFinan int
		var desc string
		var descLarga string
		if err := rowsPlanes.Scan(&idFinan, &desc, &descLarga); err == nil {
			descLimpia := strings.TrimSpace(desc)
			nombresPlanes[idFinan] = descLimpia
			idsOrdenados = append(idsOrdenados, idFinan)
		}
	}

	mapaPrecios := make(map[int]domain.BienInsumoPrecioPlan)
	if idProducto > 0 {
		preciosQuery := `EXEC dbo.CatalogoBienesInsumosHospSeleccionarXIdProducto @IdProducto = @p1`
		rowsPrecios, err := r.db.QueryContext(ctx, preciosQuery, sql.Named("p1", idProducto))
		if err == nil {
			defer rowsPrecios.Close()
			for rowsPrecios.Next() {
				var idPlanCat int
				var precioUnit float64
				var idProd int
				var idTipoFinan int
				var activo bool
				if err := rowsPrecios.Scan(&idPlanCat, &precioUnit, &idProd, &idTipoFinan, &activo); err == nil {
					mapaPrecios[idTipoFinan] = domain.BienInsumoPrecioPlan{
						IdPlanCatalogo:       idPlanCat,
						IdTipoFinanciamiento: idTipoFinan,
						PrecioUnitario:       precioUnit,
						Activo:               activo,
					}
				}
			}
		}
	}

	resultado := make([]domain.BienInsumoPrecioPlan, 0, len(idsOrdenados))
	for _, idFinan := range idsOrdenados {
		nombre := nombresPlanes[idFinan]
		if precioExistente, ok := mapaPrecios[idFinan]; ok {
			precioExistente.Plan = nombre
			resultado = append(resultado, precioExistente)
		} else {
			resultado = append(resultado, domain.BienInsumoPrecioPlan{
				IdPlanCatalogo:       0,
				IdTipoFinanciamiento: idFinan,
				Plan:                 nombre,
				PrecioUnitario:       0,
				Activo:               true,
			})
		}
	}

	return resultado, nil
}

func (r *sqlServerBienesInsumosRepository) GuardarPreciosPorPlan(ctx context.Context, idProducto int, precios []domain.BienInsumoPrecioPlan) error {
	if idProducto <= 0 || len(precios) == 0 {
		return nil
	}

	for _, item := range precios {
		activoBit := 0
		if item.Activo {
			activoBit = 1
		}

		if item.IdPlanCatalogo > 0 {
			updateSP := `EXEC dbo.FactCatalogoBienesInsumosHospModificar
				@IdPlanCatalogo = @p1,
				@PrecioUnitario = @p2,
				@IdProducto = @p3,
				@IdTipoFinanciamiento = @p4,
				@Activo = @p5,
				@IdUsuarioAuditoria = 1`

			_, err := r.db.ExecContext(ctx, updateSP,
				sql.Named("p1", item.IdPlanCatalogo),
				sql.Named("p2", item.PrecioUnitario),
				sql.Named("p3", idProducto),
				sql.Named("p4", item.IdTipoFinanciamiento),
				sql.Named("p5", activoBit),
			)
			if err != nil {
				return fmt.Errorf("error al modificar precio de plan %d: %w", item.IdTipoFinanciamiento, err)
			}
		} else {
			var nuevoIdPlan int
			insertSP := `EXEC dbo.FactCatalogoBienesInsumosHospAgregar
				@IdPlanCatalogo = @p1 OUTPUT,
				@PrecioUnitario = @p2,
				@IdProducto = @p3,
				@IdTipoFinanciamiento = @p4,
				@Activo = @p5,
				@IdUsuarioAuditoria = 1`

			_, err := r.db.ExecContext(ctx, insertSP,
				sql.Named("p1", sql.Out{Dest: &nuevoIdPlan}),
				sql.Named("p2", item.PrecioUnitario),
				sql.Named("p3", idProducto),
				sql.Named("p4", item.IdTipoFinanciamiento),
				sql.Named("p5", activoBit),
			)
			if err != nil {
				return fmt.Errorf("error al registrar precio de plan %d: %w", item.IdTipoFinanciamiento, err)
			}
		}
	}

	return nil
}

func (r *sqlServerBienesInsumosRepository) Crear(ctx context.Context, payload domain.BienInsumoGuardarPayload) (int, error) {
	petitorioBit := 0
	if payload.Petitorio {
		petitorioBit = 1
	}

	var nuevoID int
	idGrupo := payload.IdGrupoFarmacologico
	if idGrupo <= 0 {
		idGrupo = 999
	}
	idSubGrupo := payload.IdSubGrupoFarmacologico
	if idSubGrupo <= 0 {
		idSubGrupo = 999
	}

	querySP := `
		EXEC [dbo].[FactCatalogoBienesInsumosAgregar]
			@IdCentroCosto = @p1,
			@IdPartida = @p2,
			@IdSubGrupoFarmacologico = @p3,
			@IdGrupoFarmacologico = @p4,
			@NombreComercial = @p5,
			@Nombre = @p6,
			@Codigo = @p7,
			@CodigoSiga = @p8,
			@CodigoSunat = @p9,
			@IdProducto = @p10 OUTPUT,
			@IdUsuarioAuditoria = 1,
			@PrecioCompra = @p11,
			@PrecioDistribucion = @p12,
			@PrecioDonacion = @p13,
			@PrecioUltCompra = @p14,
			@idTipoSalidaBienInsumo = @p15,
			@StockMinimo = @p16,
			@TipoProducto = @p17,
			@Denominacion = @p18,
			@Concentracion = @p19,
			@Presentacion = @p20,
			@FormaFarmaceutica = @p21,
			@MaterialEnvase = @p22,
			@PresentacionEnvase = @p23,
			@Fabricante = @p24,
			@IdPaisOrigen = @p25,
			@Petitorio = @p26,
			@TipoProductoSismed = @p27,
			@Estado = @p28,
			@estadoDigemid = @p29,
			@IdsubTipoProducto = 0
	`
	_, err := r.db.ExecContext(ctx, querySP,
		sql.Named("p1", payload.IdCentroCosto),
		sql.Named("p2", payload.IdPartida),
		sql.Named("p3", idSubGrupo),
		sql.Named("p4", idGrupo),
		sql.Named("p5", payload.NombreComercial),
		sql.Named("p6", payload.Nombre),
		sql.Named("p7", payload.Codigo),
		sql.Named("p8", payload.CodigoSiga),
		sql.Named("p9", payload.CodigoSunat),
		sql.Named("p10", sql.Out{Dest: &nuevoID}),
		sql.Named("p11", payload.PrecioCompra),
		sql.Named("p12", payload.PrecioDistribucion),
		sql.Named("p13", payload.PrecioDonacion),
		sql.Named("p14", payload.PrecioUltimaCompra),
		sql.Named("p15", 1),
		sql.Named("p16", payload.StockMinimo),
		sql.Named("p17", payload.TipoProducto),
		sql.Named("p18", payload.Denominacion),
		sql.Named("p19", payload.Concentracion),
		sql.Named("p20", payload.Presentacion),
		sql.Named("p21", payload.FormaFarmaceutica),
		sql.Named("p22", payload.MaterialEnvase),
		sql.Named("p23", payload.PresentacionEnvase),
		sql.Named("p24", payload.Fabricante),
		sql.Named("p25", payload.IdPaisOrigen),
		sql.Named("p26", petitorioBit),
		sql.Named("p27", payload.TipoProductoSismed),
		sql.Named("p28", payload.Estado),
		sql.Named("p29", payload.EstadoDigemid),
	)

	if err != nil {
		return 0, fmt.Errorf("error al ejecutar procedimiento FactCatalogoBienesInsumosAgregar: %w", err)
	}

	_ = r.actualizarExtensiones(ctx, nuevoID, payload)

	if len(payload.PreciosPorPlan) > 0 {
		_ = r.GuardarPreciosPorPlan(ctx, nuevoID, payload.PreciosPorPlan)
	}

	return nuevoID, nil
}

func (r *sqlServerBienesInsumosRepository) Actualizar(ctx context.Context, idProducto int, payload domain.BienInsumoGuardarPayload) error {
	petitorioBit := 0
	if payload.Petitorio {
		petitorioBit = 1
	}

	idGrupo := payload.IdGrupoFarmacologico
	if idGrupo <= 0 {
		idGrupo = 999
	}
	idSubGrupo := payload.IdSubGrupoFarmacologico
	if idSubGrupo <= 0 {
		idSubGrupo = 999
	}

	querySP := `
		EXEC [dbo].[FactCatalogoBienesInsumosModificar]
			@IdCentroCosto = @p1,
			@IdPartida = @p2,
			@IdSubGrupoFarmacologico = @p3,
			@IdGrupoFarmacologico = @p4,
			@NombreComercial = @p5,
			@Nombre = @p6,
			@Codigo = @p7,
			@CodigoSiga = @p8,
			@CodigoSunat = @p9,
			@IdProducto = @p10,
			@IdUsuarioAuditoria = 1,
			@PrecioCompra = @p11,
			@PrecioDistribucion = @p12,
			@PrecioDonacion = @p13,
			@PrecioUltCompra = @p14,
			@idTipoSalidaBienInsumo = @p15,
			@StockMinimo = @p16,
			@TipoProducto = @p17,
			@Denominacion = @p18,
			@Concentracion = @p19,
			@Presentacion = @p20,
			@FormaFarmaceutica = @p21,
			@MaterialEnvase = @p22,
			@PresentacionEnvase = @p23,
			@Fabricante = @p24,
			@IdPaisOrigen = @p25,
			@Petitorio = @p26,
			@TipoProductoSismed = @p27,
			@Estado = @p28,
			@estadoDigemid = @p29,
			@IdsubTipoProducto = 0
	`

	_, err := r.db.ExecContext(ctx, querySP,
		sql.Named("p1", payload.IdCentroCosto),
		sql.Named("p2", payload.IdPartida),
		sql.Named("p3", idSubGrupo),
		sql.Named("p4", idGrupo),
		sql.Named("p5", payload.NombreComercial),
		sql.Named("p6", payload.Nombre),
		sql.Named("p7", payload.Codigo),
		sql.Named("p8", payload.CodigoSiga),
		sql.Named("p9", payload.CodigoSunat),
		sql.Named("p10", idProducto),
		sql.Named("p11", payload.PrecioCompra),
		sql.Named("p12", payload.PrecioDistribucion),
		sql.Named("p13", payload.PrecioDonacion),
		sql.Named("p14", payload.PrecioUltimaCompra),
		sql.Named("p15", 1),
		sql.Named("p16", payload.StockMinimo),
		sql.Named("p17", payload.TipoProducto),
		sql.Named("p18", payload.Denominacion),
		sql.Named("p19", payload.Concentracion),
		sql.Named("p20", payload.Presentacion),
		sql.Named("p21", payload.FormaFarmaceutica),
		sql.Named("p22", payload.MaterialEnvase),
		sql.Named("p23", payload.PresentacionEnvase),
		sql.Named("p24", payload.Fabricante),
		sql.Named("p25", payload.IdPaisOrigen),
		sql.Named("p26", petitorioBit),
		sql.Named("p27", payload.TipoProductoSismed),
		sql.Named("p28", payload.Estado),
		sql.Named("p29", payload.EstadoDigemid),
	)

	if err != nil {
		return fmt.Errorf("error al ejecutar procedimiento FactCatalogoBienesInsumosModificar: %w", err)
	}

	_ = r.actualizarExtensiones(ctx, idProducto, payload)

	if len(payload.PreciosPorPlan) > 0 {
		_ = r.GuardarPreciosPorPlan(ctx, idProducto, payload.PreciosPorPlan)
	}

	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (r *sqlServerBienesInsumosRepository) actualizarExtensiones(ctx context.Context, idProducto int, payload domain.BienInsumoGuardarPayload) error {
	updateExt := `EXEC dbo.usp_go_FactCatalogoBienesInsumos_ActualizarClasificaciones
		@IdProducto = @p1,
		@Vitales = @p2,
		@Estrategicos = @p3,
		@Oncologicos = @p4,
		@CargaEnFUA = @p5,
		@PetitorioIst = @p6,
		@Comunes = @p7,
		@ComiteFarmacologico = @p8,
		@Especializados = @p9,
		@IdTipoPrefijo = @p10,
		@PoliticaCanje = @p11,
		@ProductoIndicadorGestion = @p12`

	_, err := r.db.ExecContext(ctx, updateExt,
		sql.Named("p1", idProducto),
		sql.Named("p2", boolToInt(payload.EsVital)),
		sql.Named("p3", boolToInt(payload.EsEstrategico)),
		sql.Named("p4", boolToInt(payload.EsOncologico)),
		sql.Named("p5", boolToInt(payload.CargaEnFua)),
		sql.Named("p6", boolToInt(payload.PetitorioIst)),
		sql.Named("p7", boolToInt(payload.Comunes)),
		sql.Named("p8", boolToInt(payload.ComiteFarmacologico)),
		sql.Named("p9", boolToInt(payload.Especializados)),
		sql.Named("p10", payload.IdTipoPrefijo),
		sql.Named("p11", payload.PoliticaCanje),
		sql.Named("p12", boolToInt(payload.ProductoIndicadorGestion)),
	)
	return err
}

func (r *sqlServerBienesInsumosRepository) Eliminar(ctx context.Context, idProducto int) error {
	querySP := `EXEC [dbo].[FactCatalogoBienesInsumosEliminar] @IdProducto = @p1, @IdUsuarioAuditoria = 1`
	_, err := r.db.ExecContext(ctx, querySP, sql.Named("p1", idProducto))
	if err != nil {
		return fmt.Errorf("error al ejecutar procedimiento FactCatalogoBienesInsumosEliminar: %w", err)
	}
	return nil
}

func (r *sqlServerBienesInsumosRepository) ListarGruposFarmacologicos(ctx context.Context) ([]domain.GrupoFarmacologico, error) {
	query := `EXEC dbo.usp_go_FactInsumosGrupoFarmacologico_ListarTodos`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando grupos farmacológicos mediante SP: %w", err)
	}
	defer rows.Close()

	var lista []domain.GrupoFarmacologico
	for rows.Next() {
		var item domain.GrupoFarmacologico
		if err := rows.Scan(&item.IdGrupoFarmacologico, &item.Descripcion); err != nil {
			return nil, fmt.Errorf("error escaneando grupo farmacológico: %w", err)
		}
		item.Descripcion = strings.TrimSpace(item.Descripcion)
		lista = append(lista, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error recorriendo filas de grupos farmacológicos: %w", err)
	}

	return lista, nil
}

func (r *sqlServerBienesInsumosRepository) ListarSubGruposFarmacologicos(ctx context.Context, idGrupo int) ([]domain.SubGrupoFarmacologico, error) {
	query := `EXEC dbo.FactInsumosSubGrupoFarmacologicoXIdGrupoFarmacologico @lIdGrupoFarmacologico = @p1`
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idGrupo))
	if err != nil {
		return nil, fmt.Errorf("error listando subgrupos farmacológicos mediante SP: %w", err)
	}
	defer rows.Close()

	var lista []domain.SubGrupoFarmacologico
	for rows.Next() {
		var item domain.SubGrupoFarmacologico
		if err := rows.Scan(&item.IdSubGrupoFarmacologico, &item.Descripcion, &item.IdGrupoFarmacologico); err != nil {
			return nil, fmt.Errorf("error escaneando subgrupo farmacológico: %w", err)
		}
		item.Descripcion = strings.TrimSpace(item.Descripcion)
		lista = append(lista, item)
	}

	return lista, nil
}

func (r *sqlServerBienesInsumosRepository) ObtenerAuxiliares(ctx context.Context) (*domain.CatalogoBienesInsumosAuxiliares, error) {
	aux := &domain.CatalogoBienesInsumosAuxiliares{
		CentrosCosto:    make([]domain.CentroCostoItem, 0),
		Partidas:        make([]domain.PartidaPresupuestalItem, 0),
		Paises:          make([]domain.PaisItem, 0),
		Prefijos:        make([]domain.TipoPrefijoItem, 0),
		Financiamientos: make([]domain.TipoFinanciamientoItem, 0),
		Grupos:          make([]domain.GrupoFarmacologico, 0),
	}

	rowsCC, err := r.db.QueryContext(ctx, `EXEC dbo.CentrosCostoSeleccionarTodos`)
	if err == nil {
		defer rowsCC.Close()
		for rowsCC.Next() {
			var cc domain.CentroCostoItem
			if err := rowsCC.Scan(&cc.IdCentroCosto, &cc.Codigo, &cc.Descripcion); err == nil {
				cc.Codigo = strings.TrimSpace(cc.Codigo)
				cc.Descripcion = strings.TrimSpace(cc.Descripcion)
				aux.CentrosCosto = append(aux.CentrosCosto, cc)
			}
		}
	}

	rowsP, err := r.db.QueryContext(ctx, `EXEC dbo.FactPartidasPresupuestalesSeleccionarTodas`)
	if err == nil {
		defer rowsP.Close()
		for rowsP.Next() {
			var p domain.PartidaPresupuestalItem
			var seccion sql.NullInt64
			var fondo sql.NullInt64
			var orden sql.NullInt64
			if err := rowsP.Scan(&p.IdPartidaPresupuestal, &p.Codigo, &p.Descripcion, &seccion, &fondo, &orden); err == nil {
				p.Codigo = strings.TrimSpace(p.Codigo)
				p.Descripcion = strings.TrimSpace(p.Descripcion)
				aux.Partidas = append(aux.Partidas, p)
			}
		}
	}

	rowsPais, err := r.db.QueryContext(ctx, `EXEC dbo.PaisesSeleccionarTodos`)
	if err == nil {
		defer rowsPais.Close()
		for rowsPais.Next() {
			var pais domain.PaisItem
			var descLarga string
			if err := rowsPais.Scan(&pais.IdPais, &pais.Nombre, &descLarga); err == nil {
				pais.Nombre = strings.TrimSpace(pais.Nombre)
				aux.Paises = append(aux.Paises, pais)
			}
		}
	}

	rowsPref, err := r.db.QueryContext(ctx, `EXEC dbo.Usp_farmTiposSufijo_06042016`)
	if err == nil {
		defer rowsPref.Close()
		for rowsPref.Next() {
			var pref domain.TipoPrefijoItem
			if err := rowsPref.Scan(&pref.IdTipoSufijos, &pref.Descripcion); err == nil {
				pref.Descripcion = strings.TrimSpace(pref.Descripcion)
				aux.Prefijos = append(aux.Prefijos, pref)
			}
		}
	}

	rowsFin, err := r.db.QueryContext(ctx, `EXEC dbo.TiposFinanciamientoSeleccionarTodos`)
	if err == nil {
		defer rowsFin.Close()
		for rowsFin.Next() {
			var fin domain.TipoFinanciamientoItem
			var descLarga string
			if err := rowsFin.Scan(&fin.IdTipoFinanciamiento, &fin.Descripcion, &descLarga); err == nil {
				fin.Descripcion = strings.TrimSpace(fin.Descripcion)
				aux.Financiamientos = append(aux.Financiamientos, fin)
			}
		}
	}

	grupos, err := r.ListarGruposFarmacologicos(ctx)
	if err == nil {
		aux.Grupos = grupos
	}

	return aux, nil
}
