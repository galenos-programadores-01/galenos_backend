package sqlserver

import (
	"context"
	"database/sql"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type MotivoRepository struct {
	db *sql.DB
}

func NewMotivoRepository(db *sql.DB) *MotivoRepository {
	return &MotivoRepository{db: db}
}

func (r *MotivoRepository) ListarPorAtencion(ctx context.Context, idRegAtencion int) ([]domain.MotivoAtencion, error) {
	// 1. Buscar en las evoluciones médicas modernas (Tab_EvolucionesMedicas_H_E)
	queryEvol := `
		SELECT TOP 5 
			IdEvolucion AS IdMotivo, 
			NroAtencion AS IdRegAtencion, 
			ISNULL(Motivo, '') AS Motivo, 
			FechaAtencion AS FechaRegistro
		FROM dbo.Tab_EvolucionesMedicas_H_E WITH (NOLOCK)
		WHERE NroAtencion = @p1 AND ISNULL(Motivo, '') <> ''
		ORDER BY FechaAtencion DESC, IdEvolucion DESC
	`
	rows, err := r.db.QueryContext(ctx, queryEvol, sql.Named("p1", idRegAtencion))
	if err == nil {
		defer rows.Close()
		var motivos []domain.MotivoAtencion
		for rows.Next() {
			var m domain.MotivoAtencion
			var fechaRegistro sql.NullTime
			if err := rows.Scan(&m.IdMotivo, &m.IdRegAtencion, &m.Motivo, &fechaRegistro); err == nil {
				if fechaRegistro.Valid {
					m.FechaRegistro = fechaRegistro.Time.Format(time.RFC3339)
				}
				motivos = append(motivos, m)
			}
		}
		if len(motivos) > 0 {
			return motivos, nil
		}
	}

	// 2. Si no hay en evoluciones modernas, consultar registro de ingreso
	queryAntiguo := "EXEC webEvolucionMotivoListar @IdRegAtencion = @p1"
	rows2, err := r.db.QueryContext(ctx, queryAntiguo, sql.Named("p1", idRegAtencion))
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	var motivos []domain.MotivoAtencion
	for rows2.Next() {
		var m domain.MotivoAtencion
		var fechaRegistro sql.NullTime
		if err := rows2.Scan(&m.IdMotivo, &m.IdRegAtencion, &m.Motivo, &fechaRegistro); err != nil {
			return nil, err
		}
		if fechaRegistro.Valid {
			m.FechaRegistro = fechaRegistro.Time.Format(time.RFC3339)
		}
		motivos = append(motivos, m)
	}
	return motivos, rows2.Err()
}

func (r *MotivoRepository) Guardar(ctx context.Context, idRegAtencion int, motivo string) error {
	queryMod := `
		EXEC dbo.usp_go_EvolucionesMedicas_Insertar
			@NroAtencion = @p1,
			@Motivo = @p2,
			@IdEvolucion = @p3 OUTPUT,
			@Mensaje = @p4 OUTPUT
	`
	var idEvolucion sql.NullInt64
	var mensaje string
	_, err := r.db.ExecContext(ctx, queryMod,
		sql.Named("p1", idRegAtencion),
		sql.Named("p2", motivo),
		sql.Named("p3", sql.Out{Dest: &idEvolucion}),
		sql.Named("p4", sql.Out{Dest: &mensaje}),
	)
	return err
}
