package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type OrdenRepository struct {
	db *sql.DB
}

func NewOrdenRepository(db *sql.DB) *OrdenRepository {
	return &OrdenRepository{db: db}
}

func (r *OrdenRepository) ListarPorCuenta(ctx context.Context, idRegAtencion int) ([]domain.OrdenMedica, error) {
	idCuenta, err := r.resolverIdCuenta(ctx, idRegAtencion)
	if err != nil {
		return nil, err
	}

	query := "EXEC usp_go_OrdenesListarIdCuentaAtencion @idCuentaAtencion = @p1, @RecetaAdicional = @p2"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idCuenta), sql.Named("p2", -100))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordenes []domain.OrdenMedica
	for rows.Next() {
		var o domain.OrdenMedica

		var idProducto, cantidadPedida, idReceta, idItem, idUnidDosis, idFrecuencia sql.NullInt64
		var codigo, nombre, descripcionPuntoCarga, duracion, tipoProducto, notificacion sql.NullString
		var precio, total sql.NullFloat64

		if err := rows.Scan(
			&idProducto, &codigo, &nombre, &descripcionPuntoCarga, &cantidadPedida, &precio,
			&total, &idReceta, &idItem, &idUnidDosis, &idFrecuencia, &duracion,
			&tipoProducto, &notificacion,
		); err != nil {
			return nil, fmt.Errorf("error escaneando orden: %w", err)
		}

		det := domain.DetalleOrden{
			IdProducto:     int(idProducto.Int64),
			NombreProducto: nombre.String,
			Codigo:         codigo.String,
			Cantidad:       int(cantidadPedida.Int64),
			Precio:         precio.Float64,
			Total:          total.Float64,
		}

		found := false
		for i := range ordenes {
			if ordenes[i].IdOrden == int(idReceta.Int64) {
				ordenes[i].Detalles = append(ordenes[i].Detalles, det)
				found = true
				break
			}
		}
		if !found {
			o.IdOrden = int(idReceta.Int64)
			if tipoProducto.Valid {
				o.Estado = tipoProducto.String
			}
			if notificacion.Valid {
				o.Observacion = notificacion.String
			}
			o.Detalles = []domain.DetalleOrden{det}
			ordenes = append(ordenes, o)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ordenes) > 0 {
		r.completarCabeceras(ctx, idCuenta, ordenes)
	}

	return ordenes, nil
}

func (r *OrdenRepository) completarCabeceras(ctx context.Context, idCuentaAtencion int, ordenes []domain.OrdenMedica) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT rc.idReceta, CONVERT(varchar(19), rc.FechaReceta, 120),
		       RTRIM(ISNULL(em.ApellidoPaterno,'')) + ' ' + RTRIM(ISNULL(em.ApellidoMaterno,'')) + ' ' + RTRIM(ISNULL(em.Nombres,''))
		FROM RecetaCabecera rc
		LEFT JOIN Medicos m ON m.IdMedico = rc.idMedicoReceta
		LEFT JOIN Empleados em ON em.IdEmpleado = m.IdEmpleado
		WHERE rc.idCuentaAtencion = @p1`,
		sql.Named("p1", idCuentaAtencion),
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var idReceta int
		var fecha, medico string
		if err := rows.Scan(&idReceta, &fecha, &medico); err != nil {
			log.Printf("escaneando cabecera de receta: %v", err)
			continue
		}
		for i := range ordenes {
			if ordenes[i].IdOrden == idReceta {
				ordenes[i].FechaOrden = fecha
				ordenes[i].Medico = medico
			}
		}
	}
}

func (r *OrdenRepository) CrearOrden(ctx context.Context, orden domain.OrdenMedica, detalles []domain.DetalleOrden, idEmpleado int) error {
	if len(detalles) == 0 {
		return fmt.Errorf("la orden debe tener al menos un detalle")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Resolver la atención real
	var (
		idPacienteScan        sql.NullInt64
		idServicioIngresoScan sql.NullInt64
		idCuentaAtencion      sql.NullInt64
		idTipoServicioScan    sql.NullInt64
	)
	err = tx.QueryRowContext(ctx, `
		SELECT 
			ISNULL(IdCuentaAtencion, 0),
			ISNULL(IdPaciente, 0),
			ISNULL(IdServicioIngreso, 0),
			ISNULL(IdTipoServicio, 0)
		FROM dbo.Atenciones WITH (NOLOCK)
		WHERE IdAtencion = @p1`,
		sql.Named("p1", orden.IdRegAtencion),
	).Scan(&idCuentaAtencion, &idPacienteScan, &idServicioIngresoScan, &idTipoServicioScan)
	if err != nil {
		return fmt.Errorf("resolviendo atención %d: %w", orden.IdRegAtencion, err)
	}
	idPaciente := int(idPacienteScan.Int64)
	idServicioIngreso := int(idServicioIngresoScan.Int64)

	// 2. Resolver el médico real desde el empleado autenticado (JWT)
	var idMedico int
	var idEmpOut int
	var colegiatura, idColegioHIS sql.NullString
	err = tx.QueryRowContext(ctx, "EXEC usp_go_MedicosXidEmpleado @IdEmpleado = @p1",
		sql.Named("p1", idEmpleado),
	).Scan(&idEmpOut, &colegiatura, &idMedico, &idColegioHIS)
	if err != nil {
		return fmt.Errorf("el empleado %d no tiene médico asociado: %w", idEmpleado, err)
	}

	// 3. Insertar la cabecera de la receta
	const idPuntoCargaFarmacia = 5
	var respuesta string
	err = tx.QueryRowContext(ctx, `
		EXEC usp_go_RecetaCabeceraAgregar
			@Respuesta = @p1 OUTPUT,
			@IdPuntoCarga = @p2,
			@idCuentaAtencion = @p3,
			@idServicioReceta = @p4,
			@idMedicoReceta = @p5,
			@IdProducto = @p6,
			@Idpaciente = @p7,
			@IdUsuarioAuditoria = @p8,
			@IdEvolucion = @p9,
			@IdPrimeraAtencion = @p10`,
		sql.Named("p1", sql.Out{Dest: &respuesta}),
		sql.Named("p2", idPuntoCargaFarmacia),
		sql.Named("p3", int(idCuentaAtencion.Int64)),
		sql.Named("p4", idServicioIngreso),
		sql.Named("p5", idMedico),
		sql.Named("p6", detalles[0].IdProducto),
		sql.Named("p7", idPaciente),
		sql.Named("p8", idEmpleado),
		sql.Named("p9", 0),
		sql.Named("p10", 0),
	).Scan(&respuesta)
	if err != nil {
		return fmt.Errorf("creando cabecera de receta: %w", err)
	}

	idReceta, ok := parsearRespuestaReceta(respuesta)
	if !ok {
		return fmt.Errorf("el sistema rechazó la receta: %s", respuesta)
	}

	// 4. Agregar cada detalle
	for _, det := range detalles {
		precio, err := r.precioProducto(ctx, tx, det.IdProducto)
		if err != nil {
			return err
		}

		var mensaje string
		err = tx.QueryRowContext(ctx, `
			EXEC usp_go_RecetaDetalleAgregar
				@Mensaje = @p1 OUTPUT,
				@idReceta = @p2,
				@IdProducto = @p3,
				@Cantidad = @p4,
				@Precio = @p5,
				@SaldoEnRegistroReceta = @p6,
				@idDosisRecetada = @p7,
				@observaciones = @p8,
				@IdViaAdministracion = @p9,
				@CodigoDiagnostico = @p10,
				@Justificacion = @p11,
				@idUNIDDosisReceta = @p12,
				@idFrecuencia = @p13,
				@DescripcionadicionalReceta = @p14,
				@Duracion = @p15,
				@idCuentaAtencionProxCita = @p16`,
			sql.Named("p1", sql.Out{Dest: &mensaje}),
			sql.Named("p2", idReceta),
			sql.Named("p3", det.IdProducto),
			sql.Named("p4", det.Cantidad),
			sql.Named("p5", precio),
			sql.Named("p6", nil),
			sql.Named("p7", nil),
			sql.Named("p8", det.Indicaciones),
			sql.Named("p9", nil),
			sql.Named("p10", nil),
			sql.Named("p11", nil),
			sql.Named("p12", nil),
			sql.Named("p13", nil),
			sql.Named("p14", nil),
			sql.Named("p15", nil),
			sql.Named("p16", nil),
		).Scan(&mensaje)
		if err != nil {
			return fmt.Errorf("agregando producto %d a receta: %w", det.IdProducto, err)
		}
		if !strings.HasPrefix(strings.TrimSpace(mensaje), "OK") {
			return fmt.Errorf("el sistema rechazó el producto %d: %s", det.IdProducto, mensaje)
		}
	}

	return tx.Commit()
}

func (r *OrdenRepository) BuscarProductos(ctx context.Context, filtro string, limite int) ([]domain.ProductoCatalogo, error) {
	if limite <= 0 || limite > 100 {
		limite = 50
	}

	query := "EXEC usp_go_SelectMedicamentosFiltro @Filtro = @p1, @IdPaciente = @p2"

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", filtro),
		sql.Named("p2", 0),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []domain.ProductoCatalogo
	for rows.Next() {
		var p domain.ProductoCatalogo
		var (
			nombreLargo, formaFarm                                               sql.NullString
			stock, tipoProducto, ultimaCantidad, idDosisRecetada                 sql.NullInt64
			idUNIDDosisReceta, idFrecuencia, idViaAdministracion, ultimaDuracion sql.NullInt64
			tieneRecetaAnterior, cargaEnFUA                                      sql.NullInt64
			ultimaFechaReceta                                                    sql.NullTime
		)
		if err := rows.Scan(
			&p.IdProducto, &p.Codigo, &p.Nombre, &nombreLargo,
			&formaFarm, &stock, &p.PrecioVenta, &tipoProducto,
			&ultimaFechaReceta, &ultimaCantidad, &idDosisRecetada,
			&idUNIDDosisReceta, &idFrecuencia, &idViaAdministracion,
			&ultimaDuracion, &tieneRecetaAnterior, &cargaEnFUA,
		); err != nil {
			return nil, fmt.Errorf("error escaneando producto: %w", err)
		}
		if formaFarm.Valid {
			p.FormaFarmaceutica = formaFarm.String
		}
		if len(productos) >= limite {
			break
		}
		productos = append(productos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productos, nil
}

func (r *OrdenRepository) resolverIdCuenta(ctx context.Context, idRegAtencion int) (int, error) {
	var idCuenta int
	err := r.db.QueryRowContext(ctx, `
		SELECT ISNULL(IdCuentaAtencion, 0)
		FROM dbo.Atenciones WITH (NOLOCK)
		WHERE IdAtencion = @p1`,
		sql.Named("p1", idRegAtencion),
	).Scan(&idCuenta)
	if err == nil && idCuenta > 0 {
		return idCuenta, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return idRegAtencion, nil
}

func (r *OrdenRepository) precioProducto(ctx context.Context, tx *sql.Tx, idProducto int) (float64, error) {
	var precio float64
	err := tx.QueryRowContext(ctx, `
		SELECT ISNULL(PrecioDistribucion, 0)
		FROM dbo.FactCatalogoBienesInsumos WITH (NOLOCK)
		WHERE IdProducto = @p1`,
		sql.Named("p1", idProducto),
	).Scan(&precio)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("obteniendo precio del producto %d: %w", idProducto, err)
	}
	return precio, nil
}

func parsearRespuestaReceta(respuesta string) (int, bool) {
	respuesta = strings.TrimSpace(respuesta)
	if !strings.HasPrefix(respuesta, "OK;") {
		return 0, false
	}
	var idReceta int
	if _, err := fmt.Sscanf(strings.TrimPrefix(respuesta, "OK;"), "%d", &idReceta); err != nil {
		return 0, false
	}
	return idReceta, true
}
