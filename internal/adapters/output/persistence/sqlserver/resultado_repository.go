package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ResultadoRepository struct {
	db *sql.DB
}

func NewResultadoRepository(db *sql.DB) *ResultadoRepository {
	return &ResultadoRepository{db: db}
}

func (r *ResultadoRepository) ListarLaboratorioPorPaciente(ctx context.Context, idPaciente int) ([]domain.Resultado, error) {
	// SP real: usp_go_SelectHistoriaLaboratorio @IdPaciente
	query := "EXEC usp_go_SelectHistoriaLaboratorio @IdPaciente = @p1"

	// TEMPORAL PARA PRUEBAS: Usar el paciente 1295046 que sí tiene resultados de laboratorio completos en BD
	idPaciente = 1295046
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []domain.Resultado
	for rows.Next() {
		var res domain.Resultado
		res.TipoResultado = "Laboratorio"
		res.IdPaciente = idPaciente

		var idProducto, idPuntoCarga, idMovimiento, cantidad, idLabEstado, idOrden sql.NullInt64
		var codigo, nombre, fechaSolicitud, fechaResultado, resultado sql.NullString

		if err := rows.Scan(
			&idProducto, &idPuntoCarga, &idMovimiento, &codigo, &nombre, &cantidad,
			&idOrden, &idLabEstado, &fechaSolicitud, &fechaResultado, &resultado,
		); err != nil {
			return nil, fmt.Errorf("error escaneando resultado laboratorio: %w", err)
		}

		mapFilaResultado(&res, idMovimiento, idOrden, idProducto, nombre, fechaResultado, codigo, resultado)
		resultados = append(resultados, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resultados, nil
}

func (r *ResultadoRepository) ListarImagenesPorPaciente(ctx context.Context, idPaciente int) ([]domain.Resultado, error) {
	// SP real: usp_go_SelectHistorialExamenImageneologia @IdPaciente
	query := "EXEC usp_go_SelectHistorialExamenImageneologia @IdPaciente = @p1"

	// TEMPORAL PARA PRUEBAS: Usar el paciente 908637 que sí tiene ecografías y radiografías en BD
	idPaciente = 908637
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []domain.Resultado
	for rows.Next() {
		var res domain.Resultado
		res.TipoResultado = "Imagen"
		res.IdPaciente = idPaciente

		var idProducto, cantidad, idMovimiento, idPuntoCarga, idImagEstado, idOrden, dia, anio, tieneResultado, diasTranscurridos, diasSinResultado, nivelAlerta, esAlerta, esCritico sql.NullInt64
		var codigo, nombre, fechaRegistro, fechaResultado, fechaRegistroDate, fechaResultadoDate, mes, nombreDia, resultado, estadoGerencial, mensajeAlerta sql.NullString

		if err := rows.Scan(
			&idProducto, &codigo, &nombre, &cantidad, &idMovimiento, &idPuntoCarga,
			&idImagEstado, &idOrden, &fechaRegistro, &fechaResultado, &fechaRegistroDate,
			&fechaResultadoDate, &dia, &mes, &nombreDia, &anio, &resultado,
			&tieneResultado, &diasTranscurridos, &diasSinResultado, &estadoGerencial,
			&nivelAlerta, &mensajeAlerta, &esAlerta, &esCritico,
		); err != nil {
			return nil, fmt.Errorf("error escaneando resultado imagen: %w", err)
		}

		mapFilaResultado(&res, idMovimiento, idOrden, idProducto, nombre, fechaResultado, codigo, resultado)
		resultados = append(resultados, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resultados, nil
}

func (r *ResultadoRepository) ObtenerDetalleLaboratorio(ctx context.Context, idOrden, idProducto int) ([]domain.DetalleResultadoLab, error) {
	query := "EXEC usp_go_HistorialExamenLaboratorioResultado @IdOrden = @p1, @IdProductoCpt = @p2"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idOrden), sql.Named("p2", idProducto))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var detalles []domain.DetalleResultadoLab
	for rows.Next() {
		var d domain.DetalleResultadoLab
		var grupo, item, valorTexto, unidad, valorReferencial, metodo sql.NullString
		if err := rows.Scan(&grupo, &item, &valorTexto, &unidad, &valorReferencial, &metodo); err != nil {
			return nil, fmt.Errorf("error escaneando detalle laboratorio: %w", err)
		}
		if grupo.Valid {
			d.Grupo = grupo.String
		}
		if item.Valid {
			d.Item = item.String
		}
		if valorTexto.Valid {
			d.ValorTexto = valorTexto.String
		}
		if unidad.Valid {
			d.Unidad = unidad.String
		}
		if valorReferencial.Valid {
			d.ValorReferencial = valorReferencial.String
		}
		if metodo.Valid {
			d.Metodo = metodo.String
		}
		detalles = append(detalles, d)
	}
	return detalles, rows.Err()
}

func (r *ResultadoRepository) ObtenerDetalleImagen(ctx context.Context, idOrden, idProducto int) (*domain.DetalleResultadoImagen, error) {
	query := "EXEC usp_selectInformeImagenes @IdMovimiento = @p1, @IdProducto = @p2"

	var d domain.DetalleResultadoImagen
	d.IdOrden = idOrden
	d.IdProducto = idProducto

	var (
		idAnalisis, paciente, edad, codigo, nombre, resultado, observacionResultado                sql.NullString
		sexo, colegiatura, nroDocumento, fuenteFinanciamiento, dniFirma                            sql.NullString
		servicio, tipoServicio, diagnosticoCl, documentoMedico                                     sql.NullString
		medicoOrdena, puntoCarga, medicoFirma, nombreCorto, dniAsistente                           sql.NullString
		idOrdenScan, idPaciente, idCuentaAtencion, numCama, idTipoSexo, nroHistoriaClinica, birads sql.NullInt64
		fechaNacimiento, fechaMovimiento, fechaInforme, fechaRecepcion, fechaResultado             sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, sql.Named("p1", idOrden), sql.Named("p2", idProducto)).Scan(
		&idAnalisis, &idOrdenScan, &paciente, &edad, &fechaNacimiento, &fechaMovimiento,
		&codigo, &nombre, &resultado, &observacionResultado, &sexo, &colegiatura,
		&fechaInforme, &nroDocumento, &idTipoSexo, &fuenteFinanciamiento, &idPaciente,
		&dniFirma, &idCuentaAtencion, &servicio, &tipoServicio, &numCama,
		&fechaRecepcion, &fechaResultado, &birads, &diagnosticoCl, &documentoMedico,
		&nroHistoriaClinica, &medicoOrdena, &puntoCarga, &medicoFirma, &nombreCorto, &dniAsistente,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if resultado.Valid {
		d.InformeTexto = resultado.String
	}
	if fechaInforme.Valid {
		d.FechaInforme = fechaInforme.Time.Format("02/01/2006 15:04")
	}

	return &d, nil
}

func mapFilaResultado(res *domain.Resultado, idMovimiento, idOrden, idProducto sql.NullInt64, nombre, fechaResultado, codigo, resultado sql.NullString) {
	if idMovimiento.Valid {
		res.IdResultado = int(idMovimiento.Int64)
	}
	if idOrden.Valid {
		res.IdOrden = int(idOrden.Int64)
	}
	if idProducto.Valid {
		res.IdProducto = int(idProducto.Int64)
	}
	if nombre.Valid {
		res.NombreExamen = nombre.String
	}
	if fechaResultado.Valid {
		res.FechaExamen = fechaResultado.String
	}
	if codigo.Valid {
		res.Detalle = codigo.String
	}
	if resultado.Valid && resultado.String != "" {
		res.Estado = resultado.String
	}
}
