package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// AgregarRecetaCabecera invoca usp_go_PA_RecetaCabeceraAgregar.
// Reusa solo cabeceras de la misma IdPrimeraAtencion (no el SP compartido de Órdenes/HIS).
// Script: scripts/usp_go_PA_RecetaCabeceraAgregar.sql
func (r *primeraAtencionRepository) AgregarRecetaCabecera(ctx context.Context, in domain.RecetaCabeceraAgregar) (int, error) {
	const procedure = `dbo.usp_go_PA_RecetaCabeceraAgregar`

	if in.IdPrimeraAtencion <= 0 {
		return 0, fmt.Errorf("IdPrimeraAtencion inválido para cabecera PA")
	}

	var respuesta string
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("Respuesta", sql.Out{Dest: &respuesta}),
		sql.Named("IdPuntoCarga", in.IdPuntoCarga),
		sql.Named("idCuentaAtencion", in.IdCuentaAtencion),
		sql.Named("idServicioReceta", in.IdServicioReceta),
		sql.Named("idMedicoReceta", in.IdMedicoReceta),
		sql.Named("IdProducto", in.IdProducto),
		sql.Named("Idpaciente", in.IdPaciente),
		sql.Named("IdUsuarioAuditoria", in.IdUsuarioAuditoria),
		sql.Named("IdEvolucion", in.IdEvolucion),
		sql.Named("IdPrimeraAtencion", in.IdPrimeraAtencion),
	)
	if err != nil {
		return 0, fmt.Errorf("calling usp_go_PA_RecetaCabeceraAgregar: %w", err)
	}
	id, err := parseOKID(respuesta)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// AgregarRecetaDetalle invoca usp_go_RecetaDetalleAgregar.
// Luego actualiza Precio/Total y, si vienen, IdDiagnostico + CodigoDiagnostico.
func (r *primeraAtencionRepository) AgregarRecetaDetalle(ctx context.Context, in domain.RecetaDetalleAgregar) error {
	const procedure = `dbo.usp_go_RecetaDetalleAgregar`

	var mensaje string
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("Mensaje", sql.Out{Dest: &mensaje}),
		sql.Named("idReceta", in.IdReceta),
		sql.Named("IdProducto", in.IdProducto),
		sql.Named("Cantidad", in.Cantidad),
		sql.Named("Precio", in.Precio),
		sql.Named("SaldoEnRegistroReceta", in.SaldoEnRegistro),
		sql.Named("idDosisRecetada", in.IdDosisRecetada),
		sql.Named("observaciones", in.Observaciones),
		sql.Named("IdViaAdministracion", in.IdViaAdministracion),
		// DX no via SP: el SP histórico puede tipar mal CodigoDiagnostico.
		sql.Named("CodigoDiagnostico", nil),
		sql.Named("Justificacion", nil),
		sql.Named("idUNIDDosisReceta", in.IdUNIDDosisReceta),
		sql.Named("idFrecuencia", in.IdFrecuencia),
		sql.Named("DescripcionadicionalReceta", nil),
		sql.Named("Duracion", in.Duracion),
		sql.Named("idCuentaAtencionProxCita", nil),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_RecetaDetalleAgregar: %w", err)
	}
	if err := validarMensajeOK(mensaje); err != nil {
		return err
	}

	// 1) Precio/Total/Cantidad (+ IdDiagnostico si viene). Agrupación UI depende de esto.
	if in.IdDiagnostico != nil && *in.IdDiagnostico > 0 {
		_, errUpd := r.db.ExecContext(ctx, `
			UPDATE dbo.RecetaDetalle
			SET Precio = @Precio,
			    Total = @Total,
			    CantidadPedida = @Cantidad,
			    IdDiagnostico = @IdDiagnostico
			WHERE idReceta = @idReceta
			  AND idItem = @IdProducto`,
			sql.Named("Precio", in.Precio),
			sql.Named("Total", in.Total),
			sql.Named("Cantidad", in.Cantidad),
			sql.Named("IdDiagnostico", *in.IdDiagnostico),
			sql.Named("idReceta", in.IdReceta),
			sql.Named("IdProducto", in.IdProducto),
		)
		if errUpd != nil {
			return fmt.Errorf("actualizando Precio/Total/IdDiagnostico en RecetaDetalle: %w", errUpd)
		}
	} else {
		_, errUpd := r.db.ExecContext(ctx, `
			UPDATE dbo.RecetaDetalle
			SET Precio = @Precio,
			    Total = @Total,
			    CantidadPedida = @Cantidad
			WHERE idReceta = @idReceta
			  AND idItem = @IdProducto`,
			sql.Named("Precio", in.Precio),
			sql.Named("Total", in.Total),
			sql.Named("Cantidad", in.Cantidad),
			sql.Named("idReceta", in.IdReceta),
			sql.Named("IdProducto", in.IdProducto),
		)
		if errUpd != nil {
			return fmt.Errorf("actualizando Precio/Total en RecetaDetalle: %w", errUpd)
		}
	}

	// 2) CodigoDiagnostico (CIE-10) best-effort: si la columna es INT en algún HIS, no tumba el alta.
	if in.CodigoDiagnostico != nil {
		c := strings.TrimSpace(*in.CodigoDiagnostico)
		if c != "" {
			_, errCie := r.db.ExecContext(ctx, `
				UPDATE dbo.RecetaDetalle
				SET CodigoDiagnostico = @CodigoDiagnostico
				WHERE idReceta = @idReceta
				  AND idItem = @IdProducto`,
				sql.Named("CodigoDiagnostico", c),
				sql.Named("idReceta", in.IdReceta),
				sql.Named("IdProducto", in.IdProducto),
			)
			if errCie != nil {
				// No borrar el detalle: IdDiagnostico ya alcanza para agrupar por DX.
				fmt.Printf("warning: CodigoDiagnostico no actualizado (idReceta=%d idProducto=%d): %v\n",
					in.IdReceta, in.IdProducto, errCie)
			}
		}
	}
	return nil
}

// ActualizarRecetaDetalleCamposClinicos actualiza dosis/unidad/frecuencia/vía
// de un ítem ya insertado (no recrea cabecera ni detalle).
func (r *primeraAtencionRepository) ActualizarRecetaDetalleCamposClinicos(
	ctx context.Context,
	in domain.RecetaMedicamentoCamposActualizar,
) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE dbo.RecetaDetalle
		SET idDosisRecetada = CASE
		        WHEN @IdDosis IS NULL THEN idDosisRecetada
		        ELSE @IdDosis
		    END,
		    idUNIDDosisReceta = CASE
		        WHEN @IdUNIDDosis IS NULL THEN idUNIDDosisReceta
		        ELSE @IdUNIDDosis
		    END,
		    idFrecuencia = CASE
		        WHEN @IdFrecuencia IS NULL THEN idFrecuencia
		        ELSE @IdFrecuencia
		    END,
		    IdViaAdministracion = CASE
		        WHEN @IdViaAdministracion IS NULL THEN IdViaAdministracion
		        ELSE @IdViaAdministracion
		    END
		WHERE idReceta = @idReceta
		  AND idItem = @IdProducto`,
		sql.Named("IdDosis", nullInt(in.IdDosis)),
		sql.Named("IdUNIDDosis", nullInt(in.IdUNIDDosis)),
		sql.Named("IdFrecuencia", nullInt(in.IdFrecuencia)),
		sql.Named("IdViaAdministracion", nullInt(in.IdViaAdministracion)),
		sql.Named("idReceta", in.IdReceta),
		sql.Named("IdProducto", in.IdProducto),
	)
	if err != nil {
		return fmt.Errorf("actualizando campos clínicos RecetaDetalle: %w", err)
	}
	return nil
}

func nullInt(v *int) any {
	if v == nil || *v <= 0 {
		return nil
	}
	return *v
}

// ListarRecetaDetalle invoca usp_go_PA_RecetaDetalleListar.
// IdPuntoCarga vive en RecetaCabecera; el SP une cabecera + detalle.
func (r *primeraAtencionRepository) ListarRecetaDetalle(ctx context.Context, filtro domain.RecetaDetalleListarFiltro) ([]domain.RecetaMedicamentoItem, error) {
	const procedure = `dbo.usp_go_PA_RecetaDetalleListar`

	var idPrimera any
	if filtro.IdPrimeraAtencion > 0 {
		idPrimera = filtro.IdPrimeraAtencion
	} else {
		idPrimera = nil
	}
	var idPunto any
	if filtro.IdPuntoCarga > 0 {
		idPunto = filtro.IdPuntoCarga
	} else {
		idPunto = nil
	}
	var excluir any
	if filtro.ExcluirIdPuntoCarga > 0 {
		excluir = filtro.ExcluirIdPuntoCarga
	} else {
		excluir = nil
	}

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("IdCuentaAtencion", filtro.IdCuentaAtencion),
		sql.Named("IdPrimeraAtencion", idPrimera),
		sql.Named("IdPuntoCarga", idPunto),
		sql.Named("ExcluirIdPuntoCarga", excluir),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_PA_RecetaDetalleListar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading receta detalle list: %w", err)
	}

	items := make([]domain.RecetaMedicamentoItem, 0, len(maps))
	for _, m := range maps {
		idReceta := rowInt64(m, "idReceta", "IdReceta")
		idProd := rowInt64(m, "IdProducto", "idProducto", "idItem")
		if idReceta == nil || idProd == nil || *idReceta <= 0 || *idProd <= 0 {
			continue
		}
		item := domain.RecetaMedicamentoItem{
			IdReceta:   int(*idReceta),
			IdProducto: int(*idProd),
		}
		if v := rowInt64(m, "IdPuntoCarga", "idPuntoCarga"); v != nil {
			item.IdPuntoCarga = int(*v)
		}
		if s := rowString(m, "PuntoCarga", "puntoCarga"); s != nil {
			item.PuntoCarga = strings.TrimSpace(*s)
		}
		if v := rowInt64(m, "idCuentaAtencion", "IdCuentaAtencion"); v != nil {
			item.IdCuentaAtencion = int(*v)
		}
		if v := rowInt64(m, "idServicioReceta", "IdServicioReceta"); v != nil {
			item.IdServicioReceta = int(*v)
		}
		if v := rowInt64(m, "IdPrimeraAtencion", "idPrimeraAtencion"); v != nil {
			item.IdPrimeraAtencion = int(*v)
		}
		if v := rowInt64(m, "Cantidad", "CantidadPedida", "cantidad"); v != nil {
			item.Cantidad = int(*v)
		}
		if v := rowFloat64(m, "Precio", "precio"); v != nil {
			item.Precio = *v
		}
		if v := rowFloat64(m, "Total", "total"); v != nil {
			item.Total = *v
		} else {
			item.Total = float64(item.Cantidad) * item.Precio
		}
		if v := rowInt64(m, "IdDosis", "idDosis", "idDosisRecetada"); v != nil && *v > 0 {
			iv := int(*v)
			item.IdDosis = &iv
		}
		if v := rowInt64(m, "IdUNIDDosis", "idUNIDDosis", "idUNIDDosisReceta"); v != nil && *v > 0 {
			iv := int(*v)
			item.IdUNIDDosis = &iv
		}
		if v := rowInt64(m, "IdFrecuencia", "idFrecuencia"); v != nil && *v > 0 {
			iv := int(*v)
			item.IdFrecuencia = &iv
		}
		if v := rowInt64(m, "IdViaAdministracion", "idViaAdministracion"); v != nil && *v > 0 {
			iv := int(*v)
			item.IdViaAdministracion = &iv
		}
		if v := rowInt64(m, "Duracion", "duracion"); v != nil && *v > 0 {
			iv := int(*v)
			item.Duracion = &iv
		}
		if s := rowString(m, "Codigo", "codigo"); s != nil {
			item.Codigo = strings.TrimSpace(*s)
		}
		if s := rowString(m, "Nombre", "nombre"); s != nil {
			item.Nombre = strings.TrimSpace(*s)
		}
		if v := rowInt64(m, "IdDiagnostico", "idDiagnostico"); v != nil && *v > 0 {
			iv := int(*v)
			item.IdDiagnostico = &iv
		}
		if s := rowString(m, "CodigoDiagnostico", "codigoDiagnostico"); s != nil {
			item.CodigoDiagnostico = strings.TrimSpace(*s)
		}
		items = append(items, item)
	}
	return items, nil
}

// EliminarRecetaDetalle invoca usp_go_PA_RecetaDetalleEliminar.
// El SP borra el detalle y, si la idReceta queda sin ítems, también la cabecera.
func (r *primeraAtencionRepository) EliminarRecetaDetalle(ctx context.Context, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion int) error {
	const procedure = `dbo.usp_go_PA_RecetaDetalleEliminar`

	var mensaje string
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("idReceta", idReceta),
		sql.Named("IdProducto", idProducto),
		sql.Named("IdCuentaAtencion", idCuentaAtencion),
		sql.Named("IdPrimeraAtencion", idPrimeraAtencion),
		sql.Named("Mensaje", sql.Out{Dest: &mensaje}),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_PA_RecetaDetalleEliminar: %w", err)
	}
	if err := validarMensajeSP(mensaje); err != nil {
		return err
	}
	return nil
}

func parseOKID(raw string) (int, error) {
	m := strings.TrimSpace(raw)
	if m == "" {
		return 0, fmt.Errorf("respuesta vacía del SP")
	}
	upper := strings.ToUpper(m)
	if strings.HasPrefix(upper, "ERROR") || (!strings.HasPrefix(upper, "OK") && !isDigitsStart(m)) {
		// Mensajes de negocio sin prefijo ERROR (cabecera puede devolver texto libre)
		if !strings.HasPrefix(upper, "OK") {
			return 0, fmt.Errorf("%s", m)
		}
	}
	parts := strings.SplitN(m, ";", 2)
	idRaw := strings.TrimSpace(parts[0])
	if strings.EqualFold(idRaw, "OK") && len(parts) == 2 {
		idRaw = strings.TrimSpace(parts[1])
	}
	id, err := strconv.Atoi(idRaw)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("OUTPUT inválido del SP: %s", m)
	}
	return id, nil
}

func validarMensajeOK(mensaje string) error {
	m := strings.TrimSpace(mensaje)
	if m == "" {
		return fmt.Errorf("respuesta vacía del SP de detalle")
	}
	upper := strings.ToUpper(m)
	if strings.HasPrefix(upper, "OK") {
		return nil
	}
	return fmt.Errorf("%s", m)
}

func isDigitsStart(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0 && s[0] >= '0' && s[0] <= '9'
}
