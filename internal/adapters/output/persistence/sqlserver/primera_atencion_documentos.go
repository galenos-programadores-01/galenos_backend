package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

func (r *primeraAtencionRepository) ListarDocumentosFirmados(
	ctx context.Context,
	idPrimeraAtencion int,
) ([]domain.DocumentoFirmadoPA, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT NroEdicion, TipoDocumento, IdPuntoCarga, IdReceta, NombreArchivo, UbicacionArchivo, EstadoFirma
		FROM dbo.Tab_PA_DocumentoFirmado
		WHERE IdPrimeraAtencion = @id
		ORDER BY NroEdicion DESC, IdDocumentoFirmado ASC`,
		sql.Named("id", idPrimeraAtencion),
	)
	if err != nil {
		return nil, fmt.Errorf("listar documentos firmados PA: %w", err)
	}
	defer rows.Close()

	var out []domain.DocumentoFirmadoPA
	for rows.Next() {
		var item domain.DocumentoFirmadoPA
		var punto, receta sql.NullInt64
		item.IdPrimeraAtencion = idPrimeraAtencion
		if err := rows.Scan(
			&item.NroEdicion,
			&item.TipoDocumento,
			&punto,
			&receta,
			&item.NombreArchivo,
			&item.UbicacionArchivo,
			&item.EstadoFirma,
		); err != nil {
			return nil, fmt.Errorf("leyendo documento firmado PA: %w", err)
		}
		if punto.Valid {
			v := int(punto.Int64)
			item.IdPuntoCarga = &v
		}
		if receta.Valid {
			v := int(receta.Int64)
			item.IdReceta = &v
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *primeraAtencionRepository) RegistrarDocumentosFirmados(
	ctx context.Context,
	idPrimeraAtencion, nroEdicion int,
	docs []domain.DocumentoFirmadoPA,
) error {
	if idPrimeraAtencion <= 0 || len(docs) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for _, doc := range docs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO dbo.Tab_PA_DocumentoFirmado
				(IdPrimeraAtencion, NroEdicion, TipoDocumento, IdPuntoCarga, IdReceta, NombreArchivo, UbicacionArchivo, EstadoFirma, FechaRegistro)
			VALUES
				(@idPA, @nro, @tipo, @punto, @receta, @nombre, @ubi, @estado, GETDATE())`,
			sql.Named("idPA", idPrimeraAtencion),
			sql.Named("nro", nroEdicion),
			sql.Named("tipo", doc.TipoDocumento),
			sql.Named("punto", nullableInt(doc.IdPuntoCarga)),
			sql.Named("receta", nullableInt(doc.IdReceta)),
			sql.Named("nombre", doc.NombreArchivo),
			sql.Named("ubi", doc.UbicacionArchivo),
			sql.Named("estado", domain.EstadoFirmaCerrada),
		); err != nil {
			return fmt.Errorf("insertar documento firmado PA: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE dbo.Tab_PrimeraAtencion
		SET NroEdicionVigente = @nro
		WHERE IdPrimeraAtencion = @idPA`,
		sql.Named("nro", nroEdicion),
		sql.Named("idPA", idPrimeraAtencion),
	); err != nil {
		return fmt.Errorf("actualizar NroEdicionVigente: %w", err)
	}
	return tx.Commit()
}

func nullableInt(v *int) any {
	if v == nil || *v <= 0 {
		return nil
	}
	return *v
}

func (r *primeraAtencionRepository) ActualizarEdicionVigente(ctx context.Context, idPrimeraAtencion, nroEdicion int, ubicacion string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE dbo.Tab_PrimeraAtencion
		SET NroEdicionVigente = @nro,
		    UbicacionArchivo = @ubi,
		    EstadoFirma = @estado
		WHERE IdPrimeraAtencion = @id`,
		sql.Named("nro", nroEdicion),
		sql.Named("ubi", ubicacion),
		sql.Named("estado", domain.EstadoFirmaCerrada),
		sql.Named("id", idPrimeraAtencion),
	)
	if err != nil {
		return fmt.Errorf("actualizar edición vigente: %w", err)
	}
	return nil
}

func (r *primeraAtencionRepository) EliminarDetalleRectificacion(ctx context.Context, idReceta, idProducto int) error {
	if _, err := r.db.ExecContext(ctx, `
		DELETE FROM dbo.RecetaDetalle
		WHERE idReceta = @idReceta AND idItem = @idProducto`,
		sql.Named("idReceta", idReceta),
		sql.Named("idProducto", idProducto),
	); err != nil {
		return fmt.Errorf("eliminar detalle rectificación: %w", err)
	}
	var quedan int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(1) FROM dbo.RecetaDetalle WHERE idReceta = @idReceta`,
		sql.Named("idReceta", idReceta),
	).Scan(&quedan); err != nil {
		return fmt.Errorf("contar detalle receta: %w", err)
	}
	if quedan == 0 {
		if _, err := r.db.ExecContext(ctx, `
			UPDATE dbo.RecetaCabecera
			SET IdEstado = 0
			WHERE IdReceta = @idReceta`,
			sql.Named("idReceta", idReceta),
		); err != nil {
			return fmt.Errorf("anular cabecera receta: %w", err)
		}
	}
	return nil
}

func (r *primeraAtencionRepository) CabeceraRecetaAnulada(ctx context.Context, idReceta int) (bool, error) {
	var estado sql.NullInt64
	err := r.db.QueryRowContext(ctx, `
		SELECT IdEstado FROM dbo.RecetaCabecera WHERE IdReceta = @id`,
		sql.Named("id", idReceta),
	).Scan(&estado)
	if err != nil {
		return false, err
	}
	return estado.Valid && estado.Int64 == 0, nil
}