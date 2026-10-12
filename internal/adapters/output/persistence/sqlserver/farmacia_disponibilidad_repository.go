package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type farmaciaDisponibilidadRepository struct {
	db *sql.DB
}

func NewFarmaciaDisponibilidadRepository(db *sql.DB) output.FarmaciaDisponibilidadRepository {
	return &farmaciaDisponibilidadRepository{
		db: db,
	}
}

func (r *farmaciaDisponibilidadRepository) ListarAlmacenes(ctx context.Context) ([]domain.AlmacenFarmaciaItem, error) {
	query := "EXEC dbo.farmAlmacenSeleccionarTodosMenosExternos"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al listar almacenes con SP: %w", err)
	}
	defer rows.Close()

	almacenes := make([]domain.AlmacenFarmaciaItem, 0)
	for rows.Next() {
		var idAlmacen int
		var descripcion, idTipoLocales, idTipoSuministro sql.NullString
		var idEstado sql.NullInt64
		var codigoSISMED, regDias, regHora, regEstado, nombrePC, obs sql.NullString

		if err := rows.Scan(
			&idAlmacen,
			&descripcion,
			&idTipoLocales,
			&idTipoSuministro,
			&idEstado,
			&codigoSISMED,
			&regDias,
			&regHora,
			&regEstado,
			&nombrePC,
			&obs,
		); err != nil {
			return nil, fmt.Errorf("error al escanear almacen: %w", err)
		}

		if idTipoLocales.String == "F" || idTipoLocales.String == "A" {
			almacenes = append(almacenes, domain.AlmacenFarmaciaItem{
				IdAlmacen:   idAlmacen,
				Descripcion: strings.TrimSpace(descripcion.String),
				Estado:      int(idEstado.Int64),
			})
		}
	}

	return almacenes, rows.Err()
}

func (r *farmaciaDisponibilidadRepository) ObtenerKpis(ctx context.Context, idAlmacen int, canastaTipo string) (domain.DisponibilidadKpis, error) {
	almParam := idAlmacen
	if almParam <= 0 {
		almParam = 2
	}
	canastaParam := strings.TrimSpace(canastaTipo)
	if canastaParam == "" {
		canastaParam = "TODOS"
	}

	query := "EXEC dbo.usp_go_FarmDisponibilidad_ObtenerKpis @IdAlmacen = @p1, @CanastaTipo = @p2"
	row := r.db.QueryRowContext(ctx, query, sql.Named("p1", almParam), sql.Named("p2", canastaParam))

	var kpis domain.DisponibilidadKpis
	err := row.Scan(
		&kpis.TotalProductos,
		&kpis.Disponibles,
		&kpis.EnRiesgo,
		&kpis.EnQuiebre,
		&kpis.Normostock,
		&kpis.Substock,
		&kpis.Sobrestock,
		&kpis.TotalEsenciales,
		&kpis.EsencialesDisponibles,
		&kpis.PctDisponibilidadEsenciales,
		&kpis.PorcentajeOptimo,
		&kpis.LotesPorVencer30,
		&kpis.LotesPorVencer60,
	)
	if err != nil {
		return domain.DisponibilidadKpis{}, fmt.Errorf("error al obtener KPIs de disponibilidad mediante SP: %w", err)
	}

	return kpis, nil
}

func (r *farmaciaDisponibilidadRepository) ListarDisponibilidad(ctx context.Context, filtros domain.FarmaciaDisponibilidadFiltros) (domain.PaginacionDisponibilidad, error) {
	almParam := filtros.IdAlmacen
	if almParam <= 0 {
		almParam = 2
	}

	page := filtros.Page
	if page < 1 {
		page = 1
	}
	pageSize := filtros.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	soloConStockParam := false
	if filtros.SoloConStock != nil {
		soloConStockParam = *filtros.SoloConStock
	}

	soloEsencialesParam := false
	if filtros.SoloEsenciales != nil {
		soloEsencialesParam = *filtros.SoloEsenciales
	}

	var estadoSemaforoParam sql.NullString
	if filtros.EstadoSemaforo != "" {
		estadoSemaforoParam = sql.NullString{String: filtros.EstadoSemaforo, Valid: true}
	}

	var filtroParam sql.NullString
	filtroTexto := strings.TrimSpace(filtros.Filtro)
	if filtroTexto != "" {
		filtroParam = sql.NullString{String: filtroTexto, Valid: true}
	}

	canastaParam := strings.TrimSpace(filtros.CanastaTipo)
	if canastaParam == "" {
		canastaParam = "TODOS"
	}

	fInicioStr := strings.TrimSpace(filtros.FechaInicio)
	fFinStr := strings.TrimSpace(filtros.FechaFin)
	soloConMovimientoParam := false
	if filtros.SoloConMovimiento != nil {
		soloConMovimientoParam = *filtros.SoloConMovimiento
	}

	query := `EXEC dbo.usp_go_FarmDisponibilidad_ListarPaginado 
		@IdAlmacen = @p1, 
		@Filtro = @p2, 
		@SoloConStock = @p3, 
		@EstadoSemaforo = @p4, 
		@SoloEsenciales = @p5,
		@CanastaTipo = @p6,
		@FechaInicio = @p7,
		@FechaFin = @p8,
		@SoloConMovimiento = @p9,
		@PageNumber = @p10, 
		@PageSize = @p11`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", almParam),
		sql.Named("p2", filtroParam),
		sql.Named("p3", soloConStockParam),
		sql.Named("p4", estadoSemaforoParam),
		sql.Named("p5", soloEsencialesParam),
		sql.Named("p6", canastaParam),
		sql.Named("p7", fInicioStr),
		sql.Named("p8", fFinStr),
		sql.Named("p9", soloConMovimientoParam),
		sql.Named("p10", page),
		sql.Named("p11", pageSize),
	)
	if err != nil {
		return domain.PaginacionDisponibilidad{}, fmt.Errorf("error al listar disponibilidad mediante SP: %w", err)
	}
	defer rows.Close()

	paginados := make([]domain.ProductoDisponibilidadItem, 0)
	totalItems := 0

	for rows.Next() {
		var item domain.ProductoDisponibilidadItem
		var regTotal int

		err := rows.Scan(
			&item.IdProducto,
			&item.Codigo,
			&item.Nombre,
			&item.Concentracion,
			&item.FormaFarmaceutica,
			&item.IdAlmacen,
			&item.AlmacenNombre,
			&item.Cantidad,
			&item.Precio,
			&item.StockMinimo,
			&item.EstadoSemaforo,
			&item.TotalLotes,
			&item.Petitorio,
			&item.SalidasPeriodo,
			&regTotal,
		)
		if err != nil {
			return domain.PaginacionDisponibilidad{}, fmt.Errorf("error al escanear producto de disponibilidad: %w", err)
		}

		totalItems = regTotal
		item.Codigo = strings.TrimSpace(item.Codigo)
		item.Nombre = strings.TrimSpace(item.Nombre)
		item.Concentracion = strings.TrimSpace(item.Concentracion)
		item.FormaFarmaceutica = strings.TrimSpace(item.FormaFarmaceutica)
		item.AlmacenNombre = strings.TrimSpace(item.AlmacenNombre)

		paginados = append(paginados, item)
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	return domain.PaginacionDisponibilidad{
		Items:      paginados,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r *farmaciaDisponibilidadRepository) ListarLotesPorProducto(ctx context.Context, idProducto int, idAlmacen int) ([]domain.LoteItem, error) {
	almParam := idAlmacen
	if almParam <= 0 {
		almParam = 2
	}

	query := "EXEC dbo.usp_go_FarmDisponibilidad_ListarLotesPorProducto @IdProducto = @p1, @IdAlmacen = @p2"
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", idProducto),
		sql.Named("p2", almParam),
	)
	if err != nil {
		return nil, fmt.Errorf("error al listar lotes por producto mediante SP: %w", err)
	}
	defer rows.Close()

	lotes := make([]domain.LoteItem, 0)
	for rows.Next() {
		var item domain.LoteItem
		errScan := rows.Scan(
			&item.Lote,
			&item.FechaVencimiento,
			&item.Cantidad,
			&item.Precio,
			&item.NombreMarca,
			&item.DiasParaVencer,
			&item.EstadoVencimiento,
		)
		if errScan != nil {
			return nil, fmt.Errorf("error al escanear lote: %w", errScan)
		}

		item.Lote = strings.TrimSpace(item.Lote)
		item.NombreMarca = strings.TrimSpace(item.NombreMarca)
		item.FechaVencimiento = strings.TrimSpace(item.FechaVencimiento)
		item.EstadoVencimiento = strings.TrimSpace(item.EstadoVencimiento)

		lotes = append(lotes, item)
	}

	return lotes, nil
}

func (r *farmaciaDisponibilidadRepository) ListarLotesPorVencer(ctx context.Context, idAlmacen int, diasMax int) ([]domain.LoteVencimientoDetalleItem, error) {
	almParam := idAlmacen
	if almParam <= 0 {
		almParam = 2
	}
	if diasMax <= 0 {
		diasMax = 60
	}

	query := "EXEC dbo.usp_go_FarmDisponibilidad_ListarLotesPorVencer @IdAlmacen = @p1, @DiasMax = @p2"
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", almParam),
		sql.Named("p2", diasMax),
	)
	if err != nil {
		return nil, fmt.Errorf("error al listar lotes por vencer mediante SP: %w", err)
	}
	defer rows.Close()

	lotes := make([]domain.LoteVencimientoDetalleItem, 0)
	for rows.Next() {
		var item domain.LoteVencimientoDetalleItem
		errScan := rows.Scan(
			&item.IdProducto,
			&item.Codigo,
			&item.Nombre,
			&item.Lote,
			&item.FechaVencimiento,
			&item.DiasParaVencer,
			&item.EstadoVencimiento,
			&item.Cantidad,
			&item.Precio,
			&item.NombreMarca,
		)
		if errScan != nil {
			return nil, fmt.Errorf("error al escanear lote por vencer: %w", errScan)
		}

		item.Codigo = strings.TrimSpace(item.Codigo)
		item.Nombre = strings.TrimSpace(item.Nombre)
		item.Lote = strings.TrimSpace(item.Lote)
		item.FechaVencimiento = strings.TrimSpace(item.FechaVencimiento)
		item.EstadoVencimiento = strings.TrimSpace(item.EstadoVencimiento)
		item.NombreMarca = strings.TrimSpace(item.NombreMarca)

		lotes = append(lotes, item)
	}

	return lotes, nil
}

func (r *farmaciaDisponibilidadRepository) ConsultarKardex(ctx context.Context, idProducto int, idAlmacen int, fechaInicio, fechaFin time.Time) ([]domain.KardexMovimientoItem, error) {
	almParam := idAlmacen
	if almParam <= 0 {
		almParam = 2
	}

	query := `EXEC dbo.FarmDevuelveMovimientosDeProducto_Kardex 
		@IdProducto = @p1, 
		@FechaInicio = @p2, 
		@FechaFin = @p3, 
		@IdAlmacen = @p4`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", idProducto),
		sql.Named("p2", fechaInicio),
		sql.Named("p3", fechaFin),
		sql.Named("p4", almParam),
	)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar SP FarmDevuelveMovimientosDeProducto_Kardex: %w", err)
	}
	defer rows.Close()

	movimientos := make([]domain.KardexMovimientoItem, 0)
	for rows.Next() {
		var item domain.KardexMovimientoItem
		var fDestino, fOrigen, lote, abrev, docNum, concepto, datPaciente, usuario, observaciones sql.NullString
		var fechaVenc, fechaCreacion sql.NullTime
		var idProd, idEstadoMov, idAlmOrig, idAlmDest int
		var idTipoSalida sql.NullInt64

		errScan := rows.Scan(
			&item.MovNumero,
			&item.MovTipo,
			&fDestino,
			&fOrigen,
			&idProd,
			&lote,
			&fechaVenc,
			&item.Cantidad,
			&item.Precio,
			&item.Total,
			&abrev,
			&docNum,
			&fechaCreacion,
			&idEstadoMov,
			&concepto,
			&idAlmOrig,
			&idAlmDest,
			&idTipoSalida,
			&datPaciente,
			&usuario,
			&observaciones,
		)
		if errScan != nil {
			return nil, fmt.Errorf("error al escanear kardex: %w", errScan)
		}

		item.AlmacenDestino = fDestino.String
		item.AlmacenOrigen = fOrigen.String
		item.Lote = lote.String
		if fechaVenc.Valid {
			item.FechaVencimiento = fechaVenc.Time.Format("2006-01-02")
		}
		item.DocumentoTipo = abrev.String
		item.DocumentoNumero = docNum.String
		if fechaCreacion.Valid {
			item.FechaCreacion = fechaCreacion.Time.Format("2006-01-02 15:04:05")
		}
		item.Concepto = concepto.String
		item.Usuario = usuario.String
		item.Observaciones = observaciones.String

		movimientos = append(movimientos, item)
	}

	return movimientos, nil
}

func (r *farmaciaDisponibilidadRepository) ObtenerTendenciasSalidas(ctx context.Context, idAlmacen int) (domain.TendenciasSalidasFarmacia, error) {
	almParam := idAlmacen
	if almParam <= 0 {
		almParam = 2
	}

	query := "EXEC dbo.usp_go_FarmDisponibilidad_ObtenerTendenciasSalidas @MesesAtras = 6, @IdAlmacen = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", almParam))
	if err != nil {
		return domain.TendenciasSalidasFarmacia{}, fmt.Errorf("error al obtener tendencias de salidas mediante SP: %w", err)
	}
	defer rows.Close()

	periodosSet := make([]string, 0)
	periodosMap := make(map[string]bool)
	modalidadesMap := make(map[string]map[string]float64)

	for rows.Next() {
		var periodo, modalidad string
		var unidades float64
		if errScan := rows.Scan(&periodo, &modalidad, &unidades); errScan == nil {
			periodo = strings.TrimSpace(periodo)
			modalidad = strings.TrimSpace(modalidad)
			if !periodosMap[periodo] {
				periodosMap[periodo] = true
				periodosSet = append(periodosSet, periodo)
			}
			if _, ok := modalidadesMap[modalidad]; !ok {
				modalidadesMap[modalidad] = make(map[string]float64)
			}
			modalidadesMap[modalidad][periodo] = unidades
		}
	}

	seriesNombres := []string{
		"Dispensación SIS (Asegurados)",
		"Venta Contado (Particulares)",
		"Programas Estratégicos (MINSA)",
		"Distribución Interna a Servicios",
	}

	series := make([]domain.SerieTendenciaSalida, 0, len(seriesNombres))
	for _, nom := range seriesNombres {
		data := make([]float64, len(periodosSet))
		for i, p := range periodosSet {
			if modalidadesMap[nom] != nil {
				data[i] = modalidadesMap[nom][p]
			}
		}
		series = append(series, domain.SerieTendenciaSalida{
			Nombre: nom,
			Data:   data,
		})
	}

	return domain.TendenciasSalidasFarmacia{
		Periodos: periodosSet,
		Series:   series,
	}, nil
}

func (r *farmaciaDisponibilidadRepository) ObtenerEstadisticasDashboard(ctx context.Context, idAlmacen int, canastaTipo string, fechaInicio, fechaFin string) (domain.EstadisticasFarmaciaDashboard, error) {
	almParam := idAlmacen
	if almParam <= 0 {
		almParam = 2
	}
	canastaParam := strings.TrimSpace(canastaTipo)
	if canastaParam == "" {
		canastaParam = "TODOS"
	}

	fInicioStr := strings.TrimSpace(fechaInicio)
	fFinStr := strings.TrimSpace(fechaFin)

	estadisticas := domain.EstadisticasFarmaciaDashboard{
		TopExistencias:           make([]domain.TopStockItem, 0),
		ComparativaAlmacenes:     make([]domain.AlmacenComparativaStock, 0),
		DistribucionVencimientos: make([]domain.VencimientoTramoItem, 0),
		EsencialesCriticos:       make([]domain.ItemCanastaEvaluada, 0),
		EsencialesPorForma:       make([]domain.FormaFarmaceuticaStockItem, 0),
	}

	query := `EXEC dbo.usp_go_FarmDisponibilidad_ObtenerEstadisticasDashboard 
		@IdAlmacen = @p1, 
		@CanastaTipo = @p2, 
		@FechaInicio = @p3, 
		@FechaFin = @p4`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", almParam),
		sql.Named("p2", canastaParam),
		sql.Named("p3", fInicioStr),
		sql.Named("p4", fFinStr),
	)
	if err != nil {
		return estadisticas, fmt.Errorf("error al obtener estadisticas de dashboard con SP: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.TopStockItem
		if errScan := rows.Scan(&item.IdProducto, &item.Codigo, &item.Nombre, &item.Cantidad); errScan == nil {
			item.Codigo = strings.TrimSpace(item.Codigo)
			item.Nombre = strings.TrimSpace(item.Nombre)
			estadisticas.TopExistencias = append(estadisticas.TopExistencias, item)
		}
	}

	if rows.NextResultSet() {
		for rows.Next() {
			var comp domain.AlmacenComparativaStock
			if errScan := rows.Scan(&comp.IdAlmacen, &comp.Descripcion, &comp.Normales, &comp.EnRiesgo, &comp.EnQuiebre, &comp.TotalItems); errScan == nil {
				comp.Descripcion = strings.TrimSpace(comp.Descripcion)
				estadisticas.ComparativaAlmacenes = append(estadisticas.ComparativaAlmacenes, comp)
			}
		}
	}

	if rows.NextResultSet() {
		for rows.Next() {
			var v domain.VencimientoTramoItem
			if errScan := rows.Scan(&v.Rango, &v.CantidadLotes, &v.UnidadesTotales, &v.NivelAlerta); errScan == nil {
				v.Rango = strings.TrimSpace(v.Rango)
				v.NivelAlerta = strings.TrimSpace(v.NivelAlerta)
				estadisticas.DistribucionVencimientos = append(estadisticas.DistribucionVencimientos, v)
			}
		}
	}

	if rows.NextResultSet() {
		for rows.Next() {
			var crit domain.ItemCanastaEvaluada
			if errScan := rows.Scan(
				&crit.IdProducto,
				&crit.Codigo,
				&crit.Nombre,
				&crit.FormaFarmaceutica,
				&crit.CantidadTotal,
				&crit.ConsumoMensual,
				&crit.EstadoSemaforo,
				&crit.EsDigemid,
				&crit.EsSis,
			); errScan == nil {
				crit.Codigo = strings.TrimSpace(crit.Codigo)
				crit.Nombre = strings.TrimSpace(crit.Nombre)
				crit.FormaFarmaceutica = strings.TrimSpace(crit.FormaFarmaceutica)
				crit.EstadoSemaforo = strings.TrimSpace(crit.EstadoSemaforo)
				estadisticas.EsencialesCriticos = append(estadisticas.EsencialesCriticos, crit)
			}
		}
	}

	if rows.NextResultSet() {
		for rows.Next() {
			var forma domain.FormaFarmaceuticaStockItem
			if errScan := rows.Scan(&forma.FormaFarmaceutica, &forma.TotalItems, &forma.ConStock, &forma.SinStock); errScan == nil {
				forma.FormaFarmaceutica = strings.TrimSpace(forma.FormaFarmaceutica)
				estadisticas.EsencialesPorForma = append(estadisticas.EsencialesPorForma, forma)
			}
		}
	}

	return estadisticas, nil
}
