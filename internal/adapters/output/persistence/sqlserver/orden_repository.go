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
	// El frontend envía idRegAtencion = Atenciones.IdAtencion; el SP exige IdCuentaAtencion.
	idCuenta, err := r.resolverCuentaAtencion(ctx, idRegAtencion)
	if err != nil {
		return nil, err
	}

	query := "EXEC webOrdenesListarIdCuentaAtencion @IdCuentaAtencion = @p1, @RecetaAdicional = @p2"

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

	// Completar fecha y médico desde la cabecera real de la receta
	if len(ordenes) > 0 {
		r.completarCabeceras(ctx, idCuenta, ordenes)
	}

	return ordenes, nil
}

// completarCabeceras rellena FechaOrden y Medico de cada receta desde RecetaCabecera.
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

	// 1. Resolver la atención real (el frontend envía IdAtencion como idRegAtencion)
	var idCuentaAtencion, idPaciente, idServicioIngreso int
	var (
		apellidoPaterno, apellidoMaterno, primerNombre, segundoNombre sql.NullString
		nroHistoriaClinica, idAtencionScan, idPacienteScan, edad      sql.NullInt64
		fechaIngreso                                                  sql.NullTime
		horaIngreso, idDestinoAtencion, idDestinoServicio             sql.NullString
		idTipoCondicionAlServicio, idTipoCondicionALEstab             sql.NullInt64
		idServicioIngresoScan, idMedicoIngreso, idEspecialidadMedico  sql.NullInt64
		idMedicoEgreso                                                sql.NullInt64
		fechaEgreso                                                   sql.NullTime
		horaEgreso, idOrigenAtencion                                  sql.NullString
		fechaEgresoAdministrativo                                     sql.NullTime
		horaEgresoAdministrativo, idCondicionAlta, idTipoAlta         sql.NullString
		idServicioEgreso, idCamaIngreso, idCamaEgreso                 sql.NullInt64
		idTipoGravedad, idTipoGravedadEgreso, idTipoEdad              sql.NullInt64
		idTipoServicioScan, idFormaPago, idFuenteFinanciamiento       sql.NullInt64
		idEstadoAtencion                                              sql.NullInt64
		esPacienteExterno                                             sql.NullBool
		idSunasaPacienteHistorico, idCartaGarantia                    sql.NullInt64
		numCartaGarantia, cuentaAtencionRelacion, numIntegracion      sql.NullString
		idMotivoAnulacion                                             sql.NullInt64
		observacion, tieneSolicitudSOP, horaInicioAtencion            sql.NullString
		rutaBase                                                      sql.NullString
		fechaFirma                                                    sql.NullTime
		estadoFirma                                                   sql.NullString
		fechaRecuperacion                                             sql.NullTime
		idUsuarioRecupera                                             sql.NullInt64
		fechaRecepcion                                                sql.NullTime
		idUsuarioRecepciona                                           sql.NullInt64
		estadoFicha                                                   sql.NullString
		idInterconsulta, cirugiaDeDia, idTriaje                       sql.NullInt64
		numeradorReracion                                             sql.NullInt64
		emailEnvioAtencion                                            sql.NullString
		envioFua                                                      sql.NullInt64
		idRecetaGeneraCita                                            sql.NullInt64
		descripcion, nroDocumento                                     sql.NullString
		idDocIdentidad                                                sql.NullInt64
		email                                                         sql.NullString
		fechaNacimiento                                               sql.NullTime
		idTipoSexo                                                    sql.NullInt64
		telefono, idDistritoDomicilio                                 sql.NullString
	)
	err = tx.QueryRowContext(ctx, "EXEC AtencionesSeleccionarPorIdAtencion @idAtencion = @p1",
		sql.Named("p1", orden.IdRegAtencion),
	).Scan(
		&apellidoPaterno, &apellidoMaterno, &primerNombre, &segundoNombre,
		&nroHistoriaClinica, &idAtencionScan, &idPacienteScan, &edad,
		&fechaIngreso, &horaIngreso, &idDestinoAtencion, &idDestinoServicio,
		&idTipoCondicionAlServicio, &idTipoCondicionALEstab,
		&idServicioIngresoScan, &idMedicoIngreso, &idEspecialidadMedico,
		&idMedicoEgreso, &fechaEgreso, &horaEgreso, &idOrigenAtencion,
		&fechaEgresoAdministrativo, &horaEgresoAdministrativo,
		&idCondicionAlta, &idTipoAlta, &idServicioEgreso,
		&idCamaIngreso, &idCamaEgreso, &idTipoGravedad, &idTipoGravedadEgreso,
		&idTipoEdad, &idCuentaAtencion, &idTipoServicioScan, &idFormaPago,
		&idFuenteFinanciamiento, &idEstadoAtencion, &esPacienteExterno,
		&idSunasaPacienteHistorico, &idCartaGarantia, &numCartaGarantia,
		&cuentaAtencionRelacion, &numIntegracion, &idMotivoAnulacion,
		&observacion, &tieneSolicitudSOP, &horaInicioAtencion,
		&rutaBase, &fechaFirma, &estadoFirma, &fechaRecuperacion,
		&idUsuarioRecupera, &fechaRecepcion, &idUsuarioRecepciona,
		&estadoFicha, &idInterconsulta, &cirugiaDeDia, &idTriaje,
		&numeradorReracion, &emailEnvioAtencion, &envioFua,
		&idRecetaGeneraCita, &descripcion, &nroDocumento,
		&idDocIdentidad, &email, &fechaNacimiento, &idTipoSexo,
		&telefono, &idDistritoDomicilio,
	)
	if err != nil {
		return fmt.Errorf("resolviendo atención %d: %w", orden.IdRegAtencion, err)
	}
	idPaciente = int(idPacienteScan.Int64)
	idServicioIngreso = int(idServicioIngresoScan.Int64)

	// 2. Resolver el médico real desde el empleado autenticado (JWT)
	var idMedico int
	var idEmpOut int
	var colegiatura, idColegioHIS sql.NullString
	err = tx.QueryRowContext(ctx, "EXEC MedicosXidEmpleado @IdEmpleado = @p1",
		sql.Named("p1", idEmpleado),
	).Scan(&idEmpOut, &colegiatura, &idMedico, &idColegioHIS)
	if err != nil {
		return fmt.Errorf("el empleado %d no tiene médico asociado: %w", idEmpleado, err)
	}

	// 3. Crear cabecera de receta (SP real)
	// Params: @Respuesta OUTPUT, @IdPuntoCarga, @idCuentaAtencion, @idServicioReceta,
	//         @idMedicoReceta, @IdProducto, @Idpaciente, @IdUsuarioAuditoria,
	//         @IdEvolucion, @IdPrimeraAtencion
	// Respuesta: "OK;<IdReceta>" o mensaje de error
	const idPuntoCargaFarmacia = 5 // Punto de carga Farmacia (el más usado en recetas reales)
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
		sql.Named("p3", idCuentaAtencion),
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

	// 4. Agregar cada detalle (SP real)
	// Params: @Mensaje OUTPUT, @idReceta, @IdProducto, @Cantidad, @Precio,
	//         @SaldoEnRegistroReceta, @idDosisRecetada, @observaciones, @IdViaAdministracion,
	//         @CodigoDiagnostico, @Justificacion, @idUNIDDosisReceta, @idFrecuencia,
	//         @DescripcionadicionalReceta, @Duracion, @idCuentaAtencionProxCita
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
			return fmt.Errorf("agregando detalle %d a la receta %d: %w", det.IdProducto, idReceta, err)
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

// resolverCuentaAtencion obtiene IdCuentaAtencion desde el IdAtencion enviado
// por el frontend (idRegAtencion). Si el valor ya es una cuenta válida, lo usa directo.
func (r *OrdenRepository) resolverCuentaAtencion(ctx context.Context, idRegAtencion int) (int, error) {
	var idCuenta int
	var (
		apellidoPaterno, apellidoMaterno, primerNombre, segundoNombre sql.NullString
		nroHistoriaClinica, idAtencion, idPaciente, edad              sql.NullInt64
		fechaIngreso                                                  sql.NullTime
		horaIngreso, idDestinoAtencion, idDestinoServicio             sql.NullString
		idTipoCondicionAlServicio, idTipoCondicionALEstab             sql.NullInt64
		idServicioIngreso, idMedicoIngreso, idEspecialidadMedico      sql.NullInt64
		idMedicoEgreso                                                sql.NullInt64
		fechaEgreso                                                   sql.NullTime
		horaEgreso, idOrigenAtencion                                  sql.NullString
		fechaEgresoAdministrativo                                     sql.NullTime
		horaEgresoAdministrativo, idCondicionAlta, idTipoAlta         sql.NullString
		idServicioEgreso, idCamaIngreso, idCamaEgreso                 sql.NullInt64
		idTipoGravedad, idTipoGravedadEgreso, idTipoEdad              sql.NullInt64
		idTipoServicio, idFormaPago                                   sql.NullInt64
		idFuenteFinanciamiento, idEstadoAtencion                      sql.NullInt64
		esPacienteExterno                                             sql.NullBool
		idSunasaPacienteHistorico, idCartaGarantia                    sql.NullInt64
		numCartaGarantia, cuentaAtencionRelacion, numIntegracion      sql.NullString
		idMotivoAnulacion                                             sql.NullInt64
		observacion, tieneSolicitudSOP, horaInicioAtencion            sql.NullString
		rutaBase                                                      sql.NullString
		fechaFirma                                                    sql.NullTime
		estadoFirma                                                   sql.NullString
		fechaRecuperacion                                             sql.NullTime
		idUsuarioRecupera                                             sql.NullInt64
		fechaRecepcion                                                sql.NullTime
		idUsuarioRecepciona                                           sql.NullInt64
		estadoFicha                                                   sql.NullString
		idInterconsulta, cirugiaDeDia, idTriaje                       sql.NullInt64
		numeradorReracion                                             sql.NullInt64
		emailEnvioAtencion                                            sql.NullString
		envioFua                                                      sql.NullInt64
		idRecetaGeneraCita                                            sql.NullInt64
		descripcion, nroDocumento                                     sql.NullString
		idDocIdentidad                                                sql.NullInt64
		email                                                         sql.NullString
		fechaNacimiento                                               sql.NullTime
		idTipoSexo                                                    sql.NullInt64
		telefono, idDistritoDomicilio                                 sql.NullString
	)
	err := r.db.QueryRowContext(ctx,
		"EXEC AtencionesSeleccionarPorIdAtencion @idAtencion = @p1",
		sql.Named("p1", idRegAtencion),
	).Scan(
		&apellidoPaterno, &apellidoMaterno, &primerNombre, &segundoNombre,
		&nroHistoriaClinica, &idAtencion, &idPaciente, &edad,
		&fechaIngreso, &horaIngreso, &idDestinoAtencion, &idDestinoServicio,
		&idTipoCondicionAlServicio, &idTipoCondicionALEstab,
		&idServicioIngreso, &idMedicoIngreso, &idEspecialidadMedico,
		&idMedicoEgreso, &fechaEgreso, &horaEgreso, &idOrigenAtencion,
		&fechaEgresoAdministrativo, &horaEgresoAdministrativo,
		&idCondicionAlta, &idTipoAlta, &idServicioEgreso,
		&idCamaIngreso, &idCamaEgreso, &idTipoGravedad, &idTipoGravedadEgreso,
		&idTipoEdad, &idCuenta, &idTipoServicio, &idFormaPago,
		&idFuenteFinanciamiento, &idEstadoAtencion, &esPacienteExterno,
		&idSunasaPacienteHistorico, &idCartaGarantia, &numCartaGarantia,
		&cuentaAtencionRelacion, &numIntegracion, &idMotivoAnulacion,
		&observacion, &tieneSolicitudSOP, &horaInicioAtencion,
		&rutaBase, &fechaFirma, &estadoFirma, &fechaRecuperacion,
		&idUsuarioRecupera, &fechaRecepcion, &idUsuarioRecepciona,
		&estadoFicha, &idInterconsulta, &cirugiaDeDia, &idTriaje,
		&numeradorReracion, &emailEnvioAtencion, &envioFua,
		&idRecetaGeneraCita, &descripcion, &nroDocumento,
		&idDocIdentidad, &email, &fechaNacimiento, &idTipoSexo,
		&telefono, &idDistritoDomicilio,
	)
	if err == nil {
		return int(idCuenta), nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	return idRegAtencion, nil
}

// precioProducto obtiene el precio unitario vigente del catálogo hospitalario
// (por producto y tipo de financiamiento de la atención, con prioridad al tipo de la atención).
func (r *OrdenRepository) precioProducto(ctx context.Context, tx *sql.Tx, idProducto int) (float64, error) {
	var idPlanCatalogo, idTipoFinanciamiento int
	var precio float64
	var activo bool
	err := tx.QueryRowContext(ctx, `
		EXEC CatalogoBienesInsumosHospSeleccionarXIdProducto @IdProducto = @p1`,
		sql.Named("p1", idProducto),
	).Scan(&idPlanCatalogo, &precio, &idProducto, &idTipoFinanciamiento, &activo)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("obteniendo precio del producto %d: %w", idProducto, err)
	}
	if !activo {
		return 0, nil
	}
	return precio, nil
}

// parsearRespuestaReceta extrae el IdReceta de una respuesta "OK;<IdReceta>".
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
