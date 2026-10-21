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
	rows, err := r.db.QueryContext(ctx, "EXEC dbo.usp_go_Recetas_ListarCabecerasPorCuenta @IdCuentaAtencion = @p1",
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
		idAtencionScan        sql.NullInt64
		idPacienteScan        sql.NullInt64
		idCuentaAtencion      sql.NullInt64
		idServicioIngresoScan sql.NullInt64
		idTipoServicioScan    sql.NullInt64
	)
	err = tx.QueryRowContext(ctx, "EXEC dbo.usp_go_Atenciones_ObtenerDatosCabecera @IdAtencion = @p1",
		sql.Named("p1", orden.IdRegAtencion),
	).Scan(&idAtencionScan, &idPacienteScan, &idCuentaAtencion, &idServicioIngresoScan, &idTipoServicioScan)
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
	_, err = tx.ExecContext(ctx, `
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
	)
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
		_, err = tx.ExecContext(ctx, `
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
		)
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
	var idAt, idPac, idCuenta, idServ, idTipo int
	err := r.db.QueryRowContext(ctx, "EXEC dbo.usp_go_Atenciones_ObtenerDatosCabecera @IdAtencion = @p1",
		sql.Named("p1", idRegAtencion),
	).Scan(&idAt, &idPac, &idCuenta, &idServ, &idTipo)
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
	err := tx.QueryRowContext(ctx, "EXEC dbo.usp_go_FactCatalogoBienesInsumos_ObtenerPrecioDistribucion @IdProducto = @p1",
		sql.Named("p1", idProducto),
	).Scan(&precio)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("obteniendo precio del producto %d mediante SP: %w", idProducto, err)
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

func (r *OrdenRepository) ObtenerFormatoImpresionReceta(
	ctx context.Context,
	identificadorReceta int,
	identificadorCuentaAtencion int,
	identificadorAtencion int,
) (*domain.FormatoImpresionReceta, error) {
	if identificadorCuentaAtencion <= 0 || identificadorAtencion <= 0 {
		var cuentaObtenida sql.NullInt64
		var atencionObtenida sql.NullInt64
		consultaResolucion := `
			SELECT TOP 1 idCuentaAtencion, lnIdRegistroAtencion 
			FROM dbo.RecetaCabecera WITH (NOLOCK) 
			WHERE idReceta = @p1`
		errorResolucion := r.db.QueryRowContext(ctx, consultaResolucion, sql.Named("p1", identificadorReceta)).Scan(
			&cuentaObtenida,
			&atencionObtenida,
		)
		if errorResolucion == nil {
			if identificadorCuentaAtencion <= 0 && cuentaObtenida.Valid {
				identificadorCuentaAtencion = int(cuentaObtenida.Int64)
			}
			if identificadorAtencion <= 0 && atencionObtenida.Valid {
				identificadorAtencion = int(atencionObtenida.Int64)
			}
		}
	}

	if identificadorAtencion <= 0 && identificadorCuentaAtencion > 0 {
		var atencionDesdeCuenta sql.NullInt64
		consultaAtencion := `
			SELECT TOP 1 IdAtencion 
			FROM dbo.Atenciones WITH (NOLOCK) 
			WHERE IdCuentaAtencion = @p1`
		errorAtencion := r.db.QueryRowContext(ctx, consultaAtencion, sql.Named("p1", identificadorCuentaAtencion)).Scan(&atencionDesdeCuenta)
		if errorAtencion == nil && atencionDesdeCuenta.Valid {
			identificadorAtencion = int(atencionDesdeCuenta.Int64)
		}
	}

	formato := &domain.FormatoImpresionReceta{
		IdentificadorReceta:         identificadorReceta,
		NumeroReceta:                fmt.Sprintf("%d", identificadorReceta),
		IdentificadorCuentaAtencion: identificadorCuentaAtencion,
		IdentificadorAtencion:       identificadorAtencion,
		Diagnosticos:                make([]domain.DiagnosticoImpresionReceta, 0),
		Medicamentos:                make([]domain.MedicamentoImpresionReceta, 0),
	}

	if identificadorCuentaAtencion > 0 {
		consultaCabecera := "EXEC dbo.usp_go_AtencionesSelecionarPorCuenta @IdCuentaAtencion = @p1"
		filasCabecera, errorCabecera := r.db.QueryContext(ctx, consultaCabecera, sql.Named("p1", identificadorCuentaAtencion))
		if errorCabecera == nil {
			mapasCabecera, errorMapeo := rowsToMaps(filasCabecera)
			filasCabecera.Close()
			if errorMapeo == nil && len(mapasCabecera) > 0 {
				fila := mapasCabecera[0]
				if valor := rowString(fila, "Paciente", "NombresPaciente", "NombreCompleto", "PacienteNombre", "ApellidosNombres"); valor != nil {
					formato.PacienteNombreCompleto = *valor
				}
				if valor := rowString(fila, "Direccion", "Domicilio", "DireccionPaciente"); valor != nil {
					formato.DireccionDomicilio = *valor
				}
				if valor := rowString(fila, "NumeroFua", "FUA", "NroFua", "NumFua"); valor != nil {
					formato.NumeroFormatoUnicoAtencion = *valor
				}
				if valor := rowString(fila, "TipoAtencion", "DesTipoAtencion", "TipoServicio"); valor != nil {
					formato.TipoAtencion = *valor
				}
				if valor := rowString(fila, "TipoFinanciamiento", "FuenteFinanciamiento", "DesFuenteFinanciamiento", "Financiador", "TipoFinanc"); valor != nil {
					formato.TipoFinanciamiento = *valor
				}
				if valor := rowString(fila, "Telefono", "TelefonoPaciente", "Celular"); valor != nil {
					formato.NumeroTelefono = *valor
				}
				if valor := rowString(fila, "Especialidad", "DesEspecialidad", "Servicio", "DesServicio", "EspMedico"); valor != nil {
					formato.EspecialidadMedica = *valor
				}
				if valor := rowString(fila, "Consultorio", "DesConsultorio", "NombreConsultorio"); valor != nil {
					formato.Consultorio = *valor
				}
				if valor := rowString(fila, "NroHistoriaClinica", "HistoriaClinica", "HC", "NumHistoria"); valor != nil {
					formato.NumeroHistoriaClinica = *valor
				}
				if valor := rowString(fila, "Edad", "EdadPaciente", "EdadTexto"); valor != nil {
					formato.EdadPaciente = *valor
				}
				if valor := rowString(fila, "IdCuentaAtencion", "Cuenta", "CTA", "CuentaAtencion"); valor != nil {
					formato.CuentaAtencion = *valor
				}
				if valor := rowString(fila, "FechaAtencion", "FechaIngreso", "Fecha", "FechaRegistro"); valor != nil {
					formato.FechaAtencion = *valor
				}
				if valor := rowString(fila, "NroDocumento", "Documento", "DNI", "DocIdentidad"); valor != nil {
					formato.DocumentoNacionalIdentidad = *valor
				}
				if valor := rowString(fila, "Peso", "PesoKg", "SignoPeso"); valor != nil {
					formato.PesoKilogramos = *valor
				}
				if valor := rowString(fila, "Talla", "TallaCm", "SignoTalla"); valor != nil {
					formato.TallaCentimetros = *valor
				}
				if valor := rowString(fila, "Medico", "MedicoTratante", "NombresMedico", "NombreMedico"); valor != nil {
					formato.MedicoNombreCompleto = *valor
				}
				if valor := rowString(fila, "Colegiatura", "CMP", "NroColegiatura"); valor != nil {
					formato.MedicoColegiatura = *valor
				}
				if valor := rowString(fila, "RNE", "NroRNE", "RegistroEspecialista"); valor != nil {
					formato.MedicoRegistroEspecialidad = *valor
				}
				if valor := rowString(fila, "Usuario", "UsuarioAuditoria", "NomUsuario"); valor != nil {
					formato.UsuarioResponsable = *valor
				}
			}
		}
	}

	if identificadorAtencion > 0 {
		consultaDiagnosticos := "EXEC dbo.usp_go_AtencionesDiagnosticosIdAtencion @IdAtencion = @p1"
		filasDiagnosticos, errorDiagnosticos := r.db.QueryContext(ctx, consultaDiagnosticos, sql.Named("p1", identificadorAtencion))
		if errorDiagnosticos != nil {
			consultaDiagnosticosAlternativa := "EXEC dbo.usp_go_ObtenerDiagnosticosAtencion @IdAtencion = @p1, @IdPrimeraAtencion = 0, @IdTipoServicio = 0"
			filasDiagnosticos, errorDiagnosticos = r.db.QueryContext(ctx, consultaDiagnosticosAlternativa, sql.Named("p1", identificadorAtencion))
		}

		if errorDiagnosticos == nil {
			mapasDiagnosticos, errorMapeo := rowsToMaps(filasDiagnosticos)
			filasDiagnosticos.Close()
			if errorMapeo == nil {
				textosCodigos := make([]string, 0)
				for _, fila := range mapasDiagnosticos {
					var itemDiagnostico domain.DiagnosticoImpresionReceta
					if valor := rowString(fila, "CodigoCIE10", "Codigo", "CIE10", "CIE"); valor != nil {
						itemDiagnostico.CodigoCIE10 = *valor
					}
					if valor := rowString(fila, "Descripcion", "Nombre", "Diagnostico"); valor != nil {
						itemDiagnostico.Descripcion = *valor
					}
					if valor := rowString(fila, "Tipo", "Clasificacion", "Subclasificacion"); valor != nil {
						itemDiagnostico.TipoDiagnostico = *valor
					}
					if itemDiagnostico.CodigoCIE10 != "" {
						textosCodigos = append(textosCodigos, itemDiagnostico.CodigoCIE10)
					}
					formato.Diagnosticos = append(formato.Diagnosticos, itemDiagnostico)
				}
				formato.DiagnosticoTextoResumido = strings.Join(textosCodigos, ", ")
			}
		}
	}

	if identificadorReceta > 0 {
		consultaMedicamentos := "EXEC dbo.usp_go_wMedicamentosListarIdReceta @IdReceta = @p1"
		filasMedicamentos, errorMedicamentos := r.db.QueryContext(ctx, consultaMedicamentos, sql.Named("p1", identificadorReceta))
		if errorMedicamentos != nil {
			consultaMedicamentosAlternativa := "EXEC dbo.usp_go_Recetas_ListarDetallesPorEvolucion @IdCuentaAtencion = @p1"
			filasMedicamentos, errorMedicamentos = r.db.QueryContext(ctx, consultaMedicamentosAlternativa, sql.Named("p1", identificadorCuentaAtencion))
		}

		if errorMedicamentos == nil {
			mapasMedicamentos, errorMapeo := rowsToMaps(filasMedicamentos)
			filasMedicamentos.Close()
			if errorMapeo == nil {
				for indice, fila := range mapasMedicamentos {
					var itemMedicamento domain.MedicamentoImpresionReceta
					itemMedicamento.NumeroItem = indice + 1
					if valor := rowInt64(fila, "Item", "NroItem", "NumeroItem"); valor != nil && *valor > 0 {
						itemMedicamento.NumeroItem = int(*valor)
					}
					if valor := rowString(fila, "Codigo", "CodigoProducto", "CodProducto", "IdProducto"); valor != nil {
						itemMedicamento.CodigoProducto = *valor
					}
					if valor := rowString(fila, "Medicamento", "Nombre", "NombreProducto", "Descripcion", "NombreMedicamento"); valor != nil {
						itemMedicamento.NombreProductoCompleto = *valor
					}
					if valor := rowString(fila, "FormaFarmaceutica", "FormaFarm", "Presentacion"); valor != nil {
						itemMedicamento.FormaFarmaceutica = *valor
					}
					if valor := rowString(fila, "Dosis", "Posologia", "Indicaciones", "Indicacion", "DescripcionAdicional", "DescripcionadicionalReceta"); valor != nil {
						itemMedicamento.DosisPosologiaCompleta = *valor
					}
					if valor := rowInt64(fila, "Cantidad", "CantidadPedida", "Cant"); valor != nil {
						itemMedicamento.CantidadTotal = int(*valor)
					}
					formato.Medicamentos = append(formato.Medicamentos, itemMedicamento)
				}
			}
		}
	}

	return formato, nil
}

