package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sqlServerDiagnosticoRepository struct {
	db *sql.DB
}

func NewSqlServerDiagnosticoRepository(db *sql.DB) output.DiagnosticoRepository {
	return &sqlServerDiagnosticoRepository{db: db}
}

func (r *sqlServerDiagnosticoRepository) SearchDiagnosticos(ctx context.Context, filtro string, idAtencion, idPaciente int) ([]domain.DiagnosticoBusqueda, error) {
	query := "EXEC usp_go_SelectDiagnosticos @Filtro = @p1, @IdAtencion = @p2, @IdPaciente = @p3"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", filtro), sql.Named("p2", idAtencion), sql.Named("p3", idPaciente))
	if err != nil {
		return nil, fmt.Errorf("error querying diagnosticos: %w", err)
	}
	defer rows.Close()

	var results []domain.DiagnosticoBusqueda
	for rows.Next() {
		var d domain.DiagnosticoBusqueda
		var eMax, eMin, idSexo, yaReg sql.NullInt64
		var intra, activo, cancer sql.NullBool
		var desc, cie10, descL sql.NullString

		if err := rows.Scan(
			&d.IdDiagnostico,
			&intra,
			&desc,
			&cie10,
			&activo,
			&descL,
			&eMax,
			&eMin,
			&idSexo,
			&cancer,
			&yaReg,
		); err != nil {
			return nil, fmt.Errorf("error scanning diagnostico: %w", err)
		}

		if intra.Valid {
			if intra.Bool {
				d.Intrahospitalario = 1
			} else {
				d.Intrahospitalario = 0
			}
		}
		if desc.Valid {
			d.Descripcion = desc.String
		}
		if cie10.Valid {
			d.CodigoCIE10 = cie10.String
		}
		if activo.Valid {
			if activo.Bool {
				d.EsActivo = 1
			} else {
				d.EsActivo = 0
			}
		}
		if cancer.Valid {
			if cancer.Bool {
				d.Cancer = 1
			} else {
				d.Cancer = 0
			}
		}
		if descL.Valid {
			d.DescripcionLarga = descL.String
		}
		if eMax.Valid {
			d.EdadMaxDias = int(eMax.Int64)
		}
		if eMin.Valid {
			d.EdadMinDias = int(eMin.Int64)
		}
		if idSexo.Valid {
			d.IdTipoSexo = int(idSexo.Int64)
		}
		if yaReg.Valid {
			d.YaRegistrado = int(yaReg.Int64)
		}

		results = append(results, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando diagnosticos: %w", err)
	}

	log.Printf("Repository: Returning %d rows for filtro=%q, idAtencion=%d, idPaciente=%d", len(results), filtro, idAtencion, idPaciente)
	return results, nil
}

func (r *sqlServerDiagnosticoRepository) ListarDiagnosticos(ctx context.Context, filtro string) ([]domain.DiagnosticoSimple, error) {
	query := "EXEC usp_go_ListarDiagnosticos @Filtro = @p1"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", filtro))
	if err != nil {
		return nil, fmt.Errorf("error listing diagnosticos: %w", err)
	}
	defer rows.Close()

	var results []domain.DiagnosticoSimple
	for rows.Next() {
		var d domain.DiagnosticoSimple
		if err := rows.Scan(&d.IdDiagnostico, &d.CodigoCIE10, &d.Descripcion); err != nil {
			return nil, fmt.Errorf("error scanning diagnostico simple: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *sqlServerDiagnosticoRepository) ObtenerDiagnosticosAtencion(ctx context.Context, idAtencion int, idPrimeraAtencion *int, idEvolucion *int) ([]domain.DiagnosticoAtencion, error) {
	var idPrim, idEv any = nil, nil
	if idPrimeraAtencion != nil {
		idPrim = *idPrimeraAtencion
	}
	if idEvolucion != nil {
		idEv = *idEvolucion
	}

	query := "EXEC [dbo].[usp_go_ObtenerDiagnosticosAtencion] @IdAtencion = @p1, @IdPrimeraAtencion = @p2, @IdEvolucion = @p3"
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", idAtencion),
		sql.Named("p2", idPrim),
		sql.Named("p3", idEv),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing usp_go_ObtenerDiagnosticosAtencion: %w", err)
	}
	defer rows.Close()

	var results []domain.DiagnosticoAtencion
	for rows.Next() {
		var (
			d                                                  domain.DiagnosticoAtencion
			idAd, idAt, idDx, idSubDx                          sql.NullInt64
			cie, desc, descL, tCod, tDx, his, his1, his2, his3 sql.NullString
		)

		if err := rows.Scan(
			&idAd, &idAt, &idDx, &cie, &desc, &descL,
			&tCod, &tDx, &idSubDx, &his, &his1, &his2, &his3,
		); err != nil {
			return nil, fmt.Errorf("error scanning diagnostico atencion: %w", err)
		}

		if idAd.Valid {
			d.IdAtencionDiagnostico = int(idAd.Int64)
		}
		if idAt.Valid {
			d.IdAtencion = int(idAt.Int64)
		}
		if idDx.Valid {
			d.IdDiagnostico = int(idDx.Int64)
		}
		if cie.Valid {
			d.CodigoCIE10 = cie.String
		}
		if desc.Valid {
			d.Descripcion = desc.String
		}
		if descL.Valid {
			d.DescripcionLarga = descL.String
		}
		if tCod.Valid {
			d.TipoCodigo = tCod.String
		}
		if tDx.Valid {
			d.TipoDx = tDx.String
		}
		if idSubDx.Valid {
			d.IdSubclasificacionDx = int(idSubDx.Int64)
		}
		if his.Valid {
			d.LabConfHIS = his.String
		}
		if his1.Valid {
			d.LabConfHIS1 = his1.String
		}
		if his2.Valid {
			d.LabConfHIS2 = his2.String
		}
		if his3.Valid {
			d.LabConfHIS3 = his3.String
		}

		results = append(results, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando diagnosticos atencion: %w", err)
	}

	return results, nil
}

func (r *sqlServerDiagnosticoRepository) AgregarDiagnosticoAtencion(ctx context.Context, req domain.AgregarDiagnosticoAtencionRequest) (*domain.AgregarDiagnosticoAtencionResponse, error) {
	idDx := req.IdDiagnostico
	if idDx <= 0 && req.CodigoCIE10 != "" {
		_ = r.db.QueryRowContext(ctx, "SELECT TOP 1 IdDiagnostico FROM dbo.Diagnosticos WITH (NOLOCK) WHERE CodigoCIE10 = @p1", sql.Named("p1", strings.TrimSpace(req.CodigoCIE10))).Scan(&idDx)
	}

	idPac := req.IdPaciente
	if idPac <= 0 && req.IdAtencion > 0 {
		_ = r.db.QueryRowContext(ctx, "SELECT ISNULL(IdPaciente, 0) FROM dbo.Atenciones WITH (NOLOCK) WHERE IdAtencion = @p1", sql.Named("p1", req.IdAtencion)).Scan(&idPac)
	}

	idSubclasificacion := req.IdSubclasificacionDiagnostico
	if idSubclasificacion <= 0 {
		switch strings.ToUpper(strings.TrimSpace(req.TipoDiagnostico)) {
		case "D", "DEFINITIVO":
			idSubclasificacion = 2
		case "R", "REPETIDO":
			idSubclasificacion = 3
		default:
			idSubclasificacion = 1 // Presuntivo
		}
	}

	var resultado sql.NullString
	query := `EXEC dbo.usp_go_AtencionesDiagnosticosAgregar
		@IdAtencionDiagnostico = @p1 OUTPUT,
		@IdSubclasificacionDiagnostico = @p2,
		@IdDiagnostico = @p3,
		@IdAtencion = @p4,
		@LabConfHIS = @p5,
		@LabConfHIS_1 = @p6,
		@LabConfHIS_2 = @p7,
		@LabConfHIS_3 = @p8,
		@GrupoHIS = @p9,
		@SubGrupoHIS = @p10,
		@IdEpisodio = @p11,
		@IdEvolucion = @p12,
		@IdPrimeraAtencion = @p13,
		@IdPaciente = @p14`

	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", sql.Out{Dest: &resultado}),
		sql.Named("p2", idSubclasificacion),
		sql.Named("p3", idDx),
		sql.Named("p4", req.IdAtencion),
		sql.Named("p5", req.LabConfHIS),
		sql.Named("p6", req.LabConfHIS1),
		sql.Named("p7", req.LabConfHIS2),
		sql.Named("p8", req.LabConfHIS3),
		sql.Named("p9", req.GrupoHIS),
		sql.Named("p10", req.SubGrupoHIS),
		sql.Named("p11", req.IdEpisodio),
		sql.Named("p12", req.IdEvolucion),
		sql.Named("p13", req.IdPrimeraAtencion),
		sql.Named("p14", idPac),
	)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando usp_go_AtencionesDiagnosticosAgregar: %w", err)
	}

	resStr := strings.TrimSpace(resultado.String)
	if resStr == "-1" {
		return nil, fmt.Errorf("el diagnóstico ya existe en la atención o parámetros inválidos")
	}
	if strings.HasPrefix(resStr, "ERROR;") {
		return nil, fmt.Errorf("error en BD: %s", strings.TrimPrefix(resStr, "ERROR;"))
	}

	partes := strings.Split(resStr, ";")
	idGen, _ := strconv.Atoi(partes[0])
	descSub := ""
	if len(partes) > 1 {
		descSub = partes[1]
	}

	return &domain.AgregarDiagnosticoAtencionResponse{
		IdAtencionDiagnostico:       idGen,
		DescripcionSubclasificacion: descSub,
	}, nil
}
