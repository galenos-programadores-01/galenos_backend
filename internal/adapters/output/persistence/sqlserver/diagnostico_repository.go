package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
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

func calcularPuntuacionRelevanciaDiagnostico(item domain.DiagnosticoBusqueda, termino string) int {
	terminoNormalizado := strings.ToUpper(strings.TrimSpace(termino))
	if terminoNormalizado == "" {
		return 100
	}
	codigo := strings.ToUpper(strings.TrimSpace(item.CodigoCIE10))
	descripcion := strings.ToUpper(strings.TrimSpace(item.Descripcion))

	if codigo == terminoNormalizado {
		return 1
	}
	if strings.HasPrefix(codigo, terminoNormalizado) {
		return 2
	}
	if strings.HasPrefix(descripcion, terminoNormalizado) {
		return 3
	}
	if strings.Contains(codigo, terminoNormalizado) {
		return 4
	}
	if strings.Contains(descripcion, " "+terminoNormalizado) {
		return 5
	}
	if strings.Contains(descripcion, terminoNormalizado) {
		return 6
	}
	return 10
}

func (r *sqlServerDiagnosticoRepository) SearchDiagnosticos(ctx context.Context, filtro string, idAtencion, idPaciente int) ([]domain.DiagnosticoBusqueda, error) {
	filtroLimpio := strings.TrimSpace(filtro)
	query := "EXEC usp_go_SelectDiagnosticos @Filtro = @p1, @IdAtencion = @p2, @IdPaciente = @p3"
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", filtroLimpio), sql.Named("p2", idAtencion), sql.Named("p3", idPaciente))
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
		if len(results) >= 200 {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando diagnosticos: %w", err)
	}

	if filtroLimpio != "" && len(results) > 1 {
		sort.SliceStable(results, func(i, j int) bool {
			puntuacionI := calcularPuntuacionRelevanciaDiagnostico(results[i], filtroLimpio)
			puntuacionJ := calcularPuntuacionRelevanciaDiagnostico(results[j], filtroLimpio)
			if puntuacionI != puntuacionJ {
				return puntuacionI < puntuacionJ
			}
			return results[i].Descripcion < results[j].Descripcion
		})
	}

	if len(results) > 50 {
		results = results[:50]
	}

	log.Printf("Repository: Returning %d rows for filtro=%q, idAtencion=%d, idPaciente=%d", len(results), filtro, idAtencion, idPaciente)
	return results, nil
}

// ListarDiagnosticos invoca el SP usp_go_SelectDiagnosticos, que recibe el
// filtro de búsqueda y, de forma opcional, el id de atención y el de paciente
// para acotar el catálogo CIE-10. Por defecto ambos se envían en 0 (sin
// acotar), que es el valor que espera el SP cuando no se filtran.
//
// El SP devuelve el set completo de columnas del buscador predictivo, así que
// se reutiliza SearchDiagnosticos y se proyecta el resultado a la vista
// simple que consume este endpoint.
func (r *sqlServerDiagnosticoRepository) ListarDiagnosticos(ctx context.Context, filtro string, idAtencion, idPaciente int) ([]domain.DiagnosticoSimple, error) {
	busqueda, err := r.SearchDiagnosticos(ctx, filtro, idAtencion, idPaciente)
	if err != nil {
		return nil, err
	}

	results := make([]domain.DiagnosticoSimple, 0, len(busqueda))
	for _, d := range busqueda {
		results = append(results, domain.DiagnosticoSimple{
			IdDiagnostico: d.IdDiagnostico,
			CodigoCIE10:   d.CodigoCIE10,
			Descripcion:   d.Descripcion,
		})
	}
	return results, nil
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
	identificadorDiagnostico := req.IdDiagnostico
	if identificadorDiagnostico <= 0 && req.CodigoCIE10 != "" {
		_ = r.db.QueryRowContext(ctx, "EXEC dbo.usp_go_Diagnosticos_ObtenerIdPorCodigoCIE10 @CodigoCIE10 = @p1", sql.Named("p1", strings.TrimSpace(req.CodigoCIE10))).Scan(&identificadorDiagnostico)
	}

	identificadorPaciente := req.IdPaciente
	if identificadorPaciente <= 0 && req.IdAtencion > 0 {
		var idAt, idCta, idServ, idTipo int
		_ = r.db.QueryRowContext(ctx, "EXEC dbo.usp_go_Atenciones_ObtenerDatosCabecera @IdAtencion = @p1", sql.Named("p1", req.IdAtencion)).Scan(&idAt, &identificadorPaciente, &idCta, &idServ, &idTipo)
	}

	identificadorClasificacion := 1
	if req.IdClasificacionDiagnostico > 0 {
		identificadorClasificacion = req.IdClasificacionDiagnostico
	}

	identificadorSubclasificacion := req.IdSubclasificacionDiagnostico
	if identificadorSubclasificacion <= 0 {
		switch strings.ToUpper(strings.TrimSpace(req.TipoDiagnostico)) {
		case "D", "DEFINITIVO":
			identificadorSubclasificacion = 102
		case "R", "REPETIDO":
			identificadorSubclasificacion = 103
		default:
			identificadorSubclasificacion = 101
		}
	} else if identificadorSubclasificacion == 1 {
		identificadorSubclasificacion = 101
	} else if identificadorSubclasificacion == 2 {
		identificadorSubclasificacion = 102
	} else if identificadorSubclasificacion == 3 {
		identificadorSubclasificacion = 103
	}

	// La clasificación del diagnóstico (CIE-10) va fija en 9 por defecto.
	idClasificacionDx := req.IdClasificacionDx
	if idClasificacionDx <= 0 {
		idClasificacionDx = 9
	}

	// La clasificación del diagnóstico (CIE-10) va fija en 9 por defecto.
	idClasificacionDx := req.IdClasificacionDx
	if idClasificacionDx <= 0 {
		idClasificacionDx = 9
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
		@IdPaciente = @p14,
		@IdClasificacionDx = @p15`

	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", sql.Out{Dest: &resultado}),
		sql.Named("p2", idSubclasificacion),
		sql.Named("p3", idDx),
		sql.Named("p4", req.IdAtencion),
		sql.Named("p5", nullIfEmpty(req.LabConfHIS)),
		sql.Named("p6", nullIfEmpty(req.LabConfHIS1)),
		sql.Named("p7", nullIfEmpty(req.LabConfHIS2)),
		sql.Named("p8", nullIfEmpty(req.LabConfHIS3)),
		sql.Named("p9", nullIfZero(req.GrupoHIS)),
		sql.Named("p10", nullIfZero(req.SubGrupoHIS)),
		sql.Named("p11", nullIfZero(req.IdEpisodio)),
		sql.Named("p12", nullIfZero(req.IdEvolucion)),
		sql.Named("p13", nullIfZero(req.IdPrimeraAtencion)),
		sql.Named("p14", idPac),
		sql.Named("p15", idClasificacionDx),
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

// nullIfEmpty devuelve NULL cuando el texto viene vacío (o solo espacios),
// de lo contrario el texto recortado.
func nullIfEmpty(valor string) any {
	recortado := strings.TrimSpace(valor)
	if recortado == "" {
		return nil
	}
	return recortado
}

// nullIfZero devuelve NULL cuando el entero es 0, de lo contrario el valor.
func nullIfZero(valor int) any {
	if valor == 0 {
		return nil
	}
	return valor
}
