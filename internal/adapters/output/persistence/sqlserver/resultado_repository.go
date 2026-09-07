package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ResultadoRepository struct {
	db *sql.DB
}

func NewResultadoRepository(db *sql.DB) *ResultadoRepository {
	return &ResultadoRepository{db: db}
}

func (r *ResultadoRepository) ListarLaboratorioPorPaciente(ctx context.Context, idPaciente int) ([]domain.Resultado, error) {
	const query = "EXEC dbo.usp_go_SelectHistoriaLaboratorio @IdPaciente = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		log.Printf("[ResultadoRepository] ERROR en usp_go_SelectHistoriaLaboratorio (idPaciente=%d): %v", idPaciente, err)
		return nil, fmt.Errorf("ejecutando usp_go_SelectHistoriaLaboratorio: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		log.Printf("[ResultadoRepository] ERROR mapeando filas laboratorio (idPaciente=%d): %v", idPaciente, err)
		return nil, fmt.Errorf("mapeando resultado laboratorio: %w", err)
	}

	resultados := make([]domain.Resultado, 0, len(maps))
	for _, m := range maps {
		var res domain.Resultado
		res.TipoResultado = "Laboratorio"
		res.IdPaciente = idPaciente

		if v := rowInt64(m, "idMovimiento", "idResultado", "id"); v != nil {
			res.IdResultado = int(*v)
		}
		if v := rowInt64(m, "idOrden"); v != nil {
			res.IdOrden = int(*v)
		}
		if v := rowInt64(m, "idProducto", "idProductoCpt"); v != nil {
			res.IdProducto = int(*v)
		}
		if v := rowString(m, "nombre", "nombreProducto", "nombreExamen", "descripcion"); v != nil {
			res.NombreExamen = *v
		}
		if v := rowString(m, "fechaResultado", "fechaResultadoDate", "fechaExamen", "fechaRegistro", "fechaSolicitud", "fechaMovimiento"); v != nil {
			res.FechaExamen = *v
		}
		if v := rowString(m, "codigo", "codigoCpt", "detalle"); v != nil {
			res.Detalle = *v
		}
		if v := rowString(m, "resultado", "estado", "estadoGerencial"); v != nil {
			res.Estado = *v
		}

		resultados = append(resultados, res)
	}

	return resultados, nil
}

func (r *ResultadoRepository) ListarImagenesPorPaciente(ctx context.Context, idPaciente int) ([]domain.Resultado, error) {
	const query = "EXEC dbo.usp_go_SelectHistorialExamenImageneologia @IdPaciente = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		log.Printf("[ResultadoRepository] ERROR en usp_go_SelectHistorialExamenImageneologia (idPaciente=%d): %v", idPaciente, err)
		return nil, fmt.Errorf("ejecutando usp_go_SelectHistorialExamenImageneologia: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		log.Printf("[ResultadoRepository] ERROR mapeando filas imágenes (idPaciente=%d): %v", idPaciente, err)
		return nil, fmt.Errorf("mapeando resultado imágenes: %w", err)
	}

	resultados := make([]domain.Resultado, 0, len(maps))
	for _, m := range maps {
		var res domain.Resultado
		res.TipoResultado = "Imagen"
		res.IdPaciente = idPaciente

		if v := rowInt64(m, "idMovimiento", "idResultado", "id"); v != nil {
			res.IdResultado = int(*v)
		}
		if v := rowInt64(m, "idOrden"); v != nil {
			res.IdOrden = int(*v)
		}
		if v := rowInt64(m, "idProducto", "idProductoCpt"); v != nil {
			res.IdProducto = int(*v)
		}
		if v := rowString(m, "nombre", "nombreProducto", "nombreExamen", "descripcion"); v != nil {
			res.NombreExamen = *v
		}
		if v := rowString(m, "fechaResultado", "fechaResultadoDate", "fechaExamen", "fechaRegistro", "fechaSolicitud", "fechaMovimiento"); v != nil {
			res.FechaExamen = *v
		}
		if v := rowString(m, "codigo", "codigoCpt", "detalle"); v != nil {
			res.Detalle = *v
		}
		if v := rowString(m, "resultado", "estado", "estadoGerencial"); v != nil {
			res.Estado = *v
		}

		resultados = append(resultados, res)
	}

	return resultados, nil
}

func (r *ResultadoRepository) ObtenerDetalleLaboratorio(ctx context.Context, idOrden, idProducto int) ([]domain.DetalleResultadoLab, error) {
	const query = "EXEC dbo.usp_go_HistorialExamenLaboratorioResultado @IdOrden = @p1, @IdProductoCpt = @p2"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idOrden), sql.Named("p2", idProducto))
	if err != nil {
		log.Printf("[ResultadoRepository] ERROR en usp_go_HistorialExamenLaboratorioResultado (idOrden=%d, idProducto=%d): %v", idOrden, idProducto, err)
		return nil, fmt.Errorf("ejecutando usp_go_HistorialExamenLaboratorioResultado: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		log.Printf("[ResultadoRepository] ERROR leyendo detalle lab (idOrden=%d, idProducto=%d): %v", idOrden, idProducto, err)
		return nil, fmt.Errorf("mapeando detalle lab: %w", err)
	}

	detalles := make([]domain.DetalleResultadoLab, 0, len(maps))
	for _, m := range maps {
		var d domain.DetalleResultadoLab
		if v := rowString(m, "grupo", "gruponombre", "nombregrupo", "seccion"); v != nil {
			d.Grupo = *v
		}
		if v := rowString(m, "item", "analisis", "prueba", "nombreitem", "parametro", "nombreanalisis", "descripcion"); v != nil {
			d.Item = *v
		}
		if v := rowString(m, "valortexto", "resultado", "valor", "resultadoanalisis", "valornumerico"); v != nil {
			d.ValorTexto = *v
		}
		if v := rowString(m, "unidad", "unidadmedida", "unidaddosis", "undmedida"); v != nil {
			d.Unidad = *v
		}
		if v := rowString(m, "valorreferencial", "rangoreferencial", "referencia", "valoresreferenciales", "rangonormal", "valorreferencia"); v != nil {
			d.ValorReferencial = *v
		}
		if v := rowString(m, "metodo", "metodologia", "observacion", "metodoanalisis"); v != nil {
			d.Metodo = *v
		}
		detalles = append(detalles, d)
	}

	return detalles, nil
}

func (r *ResultadoRepository) ejecutarConsultaDetalle(ctx context.Context, query string, p1, p2 int) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error
	if p2 > 0 {
		rows, err = r.db.QueryContext(ctx, query, sql.Named("p1", p1), sql.Named("p2", p2))
	} else {
		rows, err = r.db.QueryContext(ctx, query, sql.Named("p1", p1))
	}
	if err != nil || rows == nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToMaps(rows)
}

func poblarDetalleImagen(m map[string]interface{}, d *domain.DetalleResultadoImagen) {
	if v := rowString(m, "nombre", "nombrecorto", "nombreexamen", "codigo", "descripcion", "servicio"); v != nil {
		d.NombreExamen = *v
	}
	if v := rowString(m, "resultado", "informetexto", "informe", "conclusion", "observacionresultado", "observaciones", "descripcion", "texto", "detalle", "hallazgos", "impresiondiagnostica"); v != nil {
		d.InformeTexto = *v
	}
	if v := rowString(m, "fechainforme", "fecharesultado", "fechamovimiento", "fecharegistro", "fecha"); v != nil {
		d.FechaInforme = *v
	}
}

func (r *ResultadoRepository) ObtenerDetalleImagen(ctx context.Context, idOrden, idProducto int) (*domain.DetalleResultadoImagen, error) {
	d := &domain.DetalleResultadoImagen{
		IdOrden:    idOrden,
		IdProducto: idProducto,
	}

	queries := []struct {
		sql string
		p1  int
		p2  int
	}{
		{"EXEC dbo.usp_go_SelectInformeImagenes @IdMovimiento = @p1, @IdProducto = @p2", idOrden, idProducto},
		{"EXEC dbo.usp_go_SelectInformeImagenes @IdOrden = @p1, @IdProducto = @p2", idOrden, idProducto},
		{"EXEC dbo.usp_go_SelectInformeImagenes @p1, @p2", idOrden, idProducto},
		{"EXEC dbo.usp_go_SelectInformeImagenes @IdMovimiento = @p1", idOrden, 0},
		{"EXEC dbo.usp_go_SelectInformeImagenes @IdOrden = @p1", idOrden, 0},
		{"EXEC dbo.usp_go_SelectInformeImagenes @p1", idOrden, 0},
	}

	for _, q := range queries {
		maps, err := r.ejecutarConsultaDetalle(ctx, q.sql, q.p1, q.p2)
		if err != nil || len(maps) == 0 {
			continue
		}

		poblarDetalleImagen(maps[0], d)
		if d.InformeTexto != "" {
			return d, nil
		}
	}

	return d, nil
}
