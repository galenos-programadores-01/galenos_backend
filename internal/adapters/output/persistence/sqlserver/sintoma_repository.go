package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type SintomaRepository struct {
	db *sql.DB
}

func NewSintomaRepository(db *sql.DB) *SintomaRepository {
	return &SintomaRepository{db: db}
}

func (r *SintomaRepository) ListarCatalogo(ctx context.Context) ([]domain.SintomaCatalogo, error) {
	query := "EXEC dbo.usp_go_SelectSintomaCatalogo"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando catálogo de síntomas: %w", err)
	}
	defer rows.Close()

	var sintomas []domain.SintomaCatalogo
	for rows.Next() {
		var s domain.SintomaCatalogo
		var orden sql.NullInt64
		if err := rows.Scan(&s.Sistema, &s.Sintoma, &s.IdSintoma, &orden); err != nil {
			return nil, fmt.Errorf("error escaneando síntoma del catálogo: %w", err)
		}
		if orden.Valid {
			s.Orden = int(orden.Int64)
		}
		sintomas = append(sintomas, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sintomas, nil
}

func (r *SintomaRepository) AgregarCatalogo(ctx context.Context, sistema, sintoma string, idUsuario int) error {
	query := "EXEC usp_go_AgregarSintomaCatalogo @Sistema = @p1, @Sintoma = @p2, @IdUsuario = @p3"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", sistema), sql.Named("p2", sintoma), sql.Named("p3", idUsuario))
	if err != nil {
		return fmt.Errorf("error agregando síntoma al catálogo: %w", err)
	}
	return nil
}

func (r *SintomaRepository) GuardarEvolucionSintomas(ctx context.Context, idRegAtencion int, sintomas []domain.SintomaSeleccionado, idUsuario int) error {
	usr := idUsuario
	if usr <= 0 {
		usr = 1
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error iniciando transacción de síntomas: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM dbo.Tab_Atencion_Sintomas WHERE IdRegAtencion = @p1", sql.Named("p1", idRegAtencion))
	if err != nil {
		return fmt.Errorf("error limpiando síntomas previos: %w", err)
	}

	if len(sintomas) == 0 {
		return tx.Commit()
	}

	insertQuery := `
		INSERT INTO dbo.Tab_Atencion_Sintomas
			(IdRegAtencion, Sistema, Sintoma, EsPredefinido, SG_LOG_CRE_FECH, SG_LOG_CRE_USUARIO)
		VALUES
			(@p1, @p2, @p3, 1, GETDATE(), @p4)
	`
	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("error preparando inserción de síntomas: %w", err)
	}
	defer stmt.Close()

	for _, s := range sintomas {
		sintTexto := strings.TrimSpace(s.Sintoma)
		if sintTexto == "" {
			continue
		}
		sistemaTexto := strings.TrimSpace(s.Sistema)
		if sistemaTexto == "" {
			sistemaTexto = "General"
		}
		_, err = stmt.ExecContext(ctx,
			sql.Named("p1", idRegAtencion),
			sql.Named("p2", sistemaTexto),
			sql.Named("p3", sintTexto),
			sql.Named("p4", usr),
		)
		if err != nil {
			return fmt.Errorf("error insertando síntoma '%s': %w", sintTexto, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error confirmando transacción de síntomas: %w", err)
	}

	return nil
}

func (r *SintomaRepository) ObtenerAtencionSintomas(ctx context.Context, idRegAtencion int) ([]domain.SintomaSeleccionado, error) {
	query := "EXEC usp_go_ObtenerAtencionSintomas @IdRegAtencion = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idRegAtencion))
	if err != nil {
		return nil, fmt.Errorf("error obteniendo síntomas de la atención: %w", err)
	}
	defer rows.Close()

	var sintomas []domain.SintomaSeleccionado
	for rows.Next() {
		var (
			idAtencionSintoma, idRegAtencionScan sql.NullInt64
			sistema, sintoma                     sql.NullString
			esPredefinido                        interface{}
		)
		if err := rows.Scan(&idAtencionSintoma, &idRegAtencionScan, &sistema, &sintoma, &esPredefinido); err != nil {
			return nil, fmt.Errorf("error escaneando síntoma de la atención: %w", err)
		}
		sintomas = append(sintomas, domain.SintomaSeleccionado{
			IdSintoma: int(idAtencionSintoma.Int64),
			Sistema:   sistema.String,
			Sintoma:   sintoma.String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sintomas, nil
}
