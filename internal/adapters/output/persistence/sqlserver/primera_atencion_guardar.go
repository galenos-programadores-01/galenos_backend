package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// GuardarUpsert invoca usp_go_InsertarPrimeraAtencion (INSERT o UPDATE por IdAtencion).
func (r *primeraAtencionRepository) GuardarUpsert(ctx context.Context, in domain.PrimeraAtencionUpsert) (domain.PrimeraAtencionUpsertResult, error) {
	const procedure = `dbo.usp_go_InsertarPrimeraAtencion`

	var idPrimeraAtencion int
	var mensaje string

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdAtencion", in.IdAtencion),
		sql.Named("FechaAtencion", in.FechaAtencion),
		sql.Named("TipoAnamnesis", in.TipoAnamnesis),
		sql.Named("TiempoEnfermedad", in.TiempoEnfermedad),
		sql.Named("IdTiempoEnfermedad", in.IdTiempoEnfermedad),
		sql.Named("SignosSintomas", in.SignosSintomas),
		sql.Named("RelatoCronologico", in.RelatoCronologico),
		sql.Named("Antecedentes", in.Antecedentes),
		sql.Named("IdEstadoGeneral", in.IdEstadoGeneral),
		sql.Named("IdEstadoHidratacion", in.IdEstadoHidratacion),
		sql.Named("IdNivelConciencia", in.IdNivelConciencia),
		sql.Named("PielAnexos", in.PielAnexos),
		sql.Named("ExamenClinicoRegional", in.ExamenClinicoRegional),
		sql.Named("ReferenciaOportuna", in.ReferenciaOportuna),
		sql.Named("Evolucion", in.Evolucion),
		sql.Named("EvaluacionEnfermeriaObstetricia", in.EvaluacionEnfermeriaObstetricia),
		sql.Named("IdCondicionEgreso", in.IdCondicionEgreso),
		sql.Named("CodCondicionEgreso", in.CodCondicionEgreso),
		sql.Named("FechaEgreso", in.FechaEgreso),
		sql.Named("IdUsuarioRegistro", in.IdUsuarioRegistro),
		sql.Named("IdUsuarioModifica", in.IdUsuarioModifica),
		sql.Named("FechaModifica", in.FechaModifica),
		sql.Named("EstadoFirma", in.EstadoFirma),
		sql.Named("UbicacionArchivo", in.UbicacionArchivo),
		sql.Named("IdPrimeraAtencion", sql.Out{Dest: &idPrimeraAtencion}),
		sql.Named("Mensaje", sql.Out{Dest: &mensaje}),
	)
	if err != nil {
		return domain.PrimeraAtencionUpsertResult{}, fmt.Errorf("calling usp_go_InsertarPrimeraAtencion: %w", err)
	}

	result := domain.PrimeraAtencionUpsertResult{
		IdPrimeraAtencion: idPrimeraAtencion,
		Mensaje:           mensaje,
	}
	if err := validarMensajeSP(mensaje); err != nil {
		return result, err
	}
	return result, nil
}

// IdsConBorradorActivo devuelve los IdAtencion con registro en Tab_PrimeraAtencion.
func (r *primeraAtencionRepository) IdsConBorradorActivo(ctx context.Context, idsAtencion []int) (map[int]bool, error) {
	return idsAtencionEnTabPrimeraAtencion(ctx, r.db, idsAtencion, false)
}

// IdsFirmados devuelve los IdAtencion con EstadoFirma = 1.
func (r *primeraAtencionRepository) IdsFirmados(ctx context.Context, idsAtencion []int) (map[int]bool, error) {
	return idsAtencionEnTabPrimeraAtencion(ctx, r.db, idsAtencion, true)
}

// ObtenerPorIdAtencion devuelve el borrador vigente de Tab_PrimeraAtencion.
func (r *primeraAtencionRepository) ObtenerPorIdAtencion(ctx context.Context, idAtencion int) (map[string]any, error) {
	const query = `
		SELECT TOP 1 *
		FROM dbo.Tab_PrimeraAtencion
		WHERE IdAtencion = @p1
		ORDER BY IdPrimeraAtencion DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idAtencion))
	if err != nil {
		return nil, fmt.Errorf("querying Tab_PrimeraAtencion: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading Tab_PrimeraAtencion: %w", err)
	}
	if len(maps) == 0 {
		return nil, nil
	}
	return maps[0], nil
}

// AgregarDiagnostico invoca usp_go_AtencionesDiagnosticosAgregar.
func (r *primeraAtencionRepository) AgregarDiagnostico(ctx context.Context, in domain.DiagnosticoAgregar) (domain.DiagnosticoAgregarResult, error) {
	const procedure = `dbo.usp_go_AtencionesDiagnosticosAgregar`

	var outID string

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdAtencionDiagnostico", sql.Out{Dest: &outID}),
		sql.Named("IdClasificacionDx", in.IdClasificacionDx),
		sql.Named("IdSubclasificacionDiagnostico", in.IdSubclasificacionDiagnostico),
		sql.Named("IdDiagnostico", in.IdDiagnostico),
		sql.Named("IdAtencion", in.IdAtencion),
		sql.Named("LabConfHIS", in.LabConfHIS),
		sql.Named("LabConfHIS_1", in.LabConfHIS1),
		sql.Named("LabConfHIS_2", in.LabConfHIS2),
		sql.Named("LabConfHIS_3", in.LabConfHIS3),
		sql.Named("GrupoHIS", in.GrupoHIS),
		sql.Named("SubGrupoHIS", in.SubGrupoHIS),
		sql.Named("IdEpisodio", nullIntOrValue(in.IdEpisodio)),
		sql.Named("IdEvolucion", nullIntOrValue(in.IdEvolucion)),
		sql.Named("IdPrimeraAtencion", nullIntOrValue(in.IdPrimeraAtencion)),
		sql.Named("IdPaciente", in.IdPaciente),
	)
	if err != nil {
		return domain.DiagnosticoAgregarResult{}, fmt.Errorf("calling usp_go_AtencionesDiagnosticosAgregar: %w", err)
	}

	result, err := parseOutputAtencionDiagnostico(outID)
	if err != nil {
		return result, err
	}
	return result, nil
}

// RegistrarAltaHospitalaria invoca usp_go_RegistrarAltaHospitalaria.
func (r *primeraAtencionRepository) RegistrarAltaHospitalaria(
	ctx context.Context,
	in domain.AltaHospitalariaRegistrar,
) (domain.AltaHospitalariaResult, error) {
	const procedure = `dbo.usp_go_RegistrarAltaHospitalaria`

	var resultado int
	var mensaje string

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdAtencion", in.IdAtencion),
		sql.Named("Idepisodio", nullIntPtr(in.IdEpisodio)),
		sql.Named("IdDestinoAtencion", nullIntPtr(in.IdDestinoAtencion)),
		sql.Named("IdDestinoServicio", nullIntPtr(in.IdDestinoServicio)),
		sql.Named("IdMedicoEgreso", in.IdMedicoEgreso),
		sql.Named("FechaEgreso", in.FechaEgreso),
		sql.Named("IdCondicionAlta", nullIntPtr(in.IdCondicionAlta)),
		sql.Named("IdTipoCertificado", nullIntPtr(in.IdTipoCertificado)),
		sql.Named("ActDefuncion", in.ActDefuncion),
		sql.Named("IdTipoAlta", nullIntPtr(in.IdTipoAlta)),
		sql.Named("Idpronostico", nullIntPtr(in.IdPronostico)),
		sql.Named("Recomendacion", in.Recomendacion),
		sql.Named("IdServicioEgreso", in.IdServicioEgreso),
		sql.Named("IdCamaEgreso", nullIntPtr(in.IdCamaEgreso)),
		sql.Named("IdUsuario", nullIntPtr(in.IdUsuario)),
		sql.Named("MotivoTraslado", nullBlankString(in.MotivoTraslado)),
		sql.Named("FechaFuga", nullTimePtr(in.FechaFuga)),
		sql.Named("ReportaFuga", nullBlankString(in.ReportaFuga)),
		sql.Named("IdNotificacionFuga", nullIntPtr(in.IdNotificacionFuga)),
		sql.Named("Observaciones", nullBlankString(in.Observaciones)),
		sql.Named("IdtiempoEstanciaDeseso", nullIntPtr(in.IdtiempoEstanciaDeseso)),
		sql.Named("RequiereAutopsia", nullBoolPtr(in.RequiereAutopsia)),
		sql.Named("Resultado", sql.Out{Dest: &resultado}),
		sql.Named("Mensaje", sql.Out{Dest: &mensaje}),
	)
	if err != nil {
		return domain.AltaHospitalariaResult{}, fmt.Errorf("calling usp_go_RegistrarAltaHospitalaria: %w", err)
	}

	out := domain.AltaHospitalariaResult{Resultado: resultado, Mensaje: mensaje}
	if resultado < 1 {
		msg := strings.TrimSpace(mensaje)
		if msg == "" {
			msg = fmt.Sprintf("alta hospitalaria rechazada (Resultado=%d)", resultado)
		}
		return out, fmt.Errorf("%s", msg)
	}
	if err := validarMensajeSP(mensaje); err != nil {
		return out, err
	}
	return out, nil
}

func nullIntPtr(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullBlankString(s string) any {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return s
}

func nullTimePtr(v *time.Time) any {
	if v == nil || v.IsZero() {
		return nil
	}
	return *v
}

func nullBoolPtr(v *bool) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullIntOrValue(v int) any {
	if v <= 0 {
		return nil
	}
	return v
}

// ListarDiagnosticos invoca usp_go_AtencionesDiagnosticosListar.
func (r *primeraAtencionRepository) ListarDiagnosticos(ctx context.Context, idAtencion, idPrimeraAtencion int) ([]domain.DiagnosticoAtencion, error) {
	const procedure = `dbo.usp_go_AtencionesDiagnosticosListar`

	var idPrimera any
	if idPrimeraAtencion > 0 {
		idPrimera = idPrimeraAtencion
	} else {
		idPrimera = nil
	}

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("IdAtencion", idAtencion),
		sql.Named("IdPrimeraAtencion", idPrimera),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_AtencionesDiagnosticosListar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading AtencionesDiagnosticos list: %w", err)
	}

	items := make([]domain.DiagnosticoAtencion, 0, len(maps))
	for _, m := range maps {
		idAD := rowInt64(m, "IdAtencionDiagnostico", "idAtencionDiagnostico")
		idDx := rowInt64(m, "IdDiagnostico", "idDiagnostico")
		if idAD == nil || idDx == nil || *idAD <= 0 || *idDx <= 0 {
			continue
		}
		item := domain.DiagnosticoAtencion{
			IdAtencionDiagnostico: int(*idAD),
			IdDiagnostico:         int(*idDx),
		}
		if v := rowInt64(m, "IdAtencion", "idAtencion"); v != nil {
			item.IdAtencion = int(*v)
		}
		if v := rowInt64(m, "idPrimeraAtencion", "IdPrimeraAtencion"); v != nil {
			item.IdPrimeraAtencion = int(*v)
		}
		if v := rowInt64(m, "IdSubclasificacionDx", "idSubclasificacionDx", "IdSubclasificacionDiagnostico"); v != nil {
			item.IdSubclasificacionDx = int(*v)
		}
		if v := rowInt64(m, "NumOrdenCorelativo", "NumOrdenCorrelativo", "numOrdenCorelativo"); v != nil {
			item.NumOrdenCorrelativo = int(*v)
		}
		if s := rowString(m, "Descripcion", "descripcion"); s != nil {
			item.Descripcion = strings.TrimSpace(*s)
		}
		if s := rowString(m, "CodigoCIE10", "codigoCIE10", "CodigoCie10"); s != nil {
			item.CodigoCIE10 = strings.TrimSpace(*s)
		}
		// Tipo (P|D|R) lo resuelve el FE con el catálogo SP; aquí solo persistimos el id.
		items = append(items, item)
	}
	return items, nil
}

// EliminarDiagnostico invoca usp_go_AtencionesDiagnosticosEliminar.
func (r *primeraAtencionRepository) EliminarDiagnostico(ctx context.Context, idAtencion, idPrimeraAtencion, idAtencionDiagnostico int) error {
	const procedure = `dbo.usp_go_AtencionesDiagnosticosEliminar`

	var mensaje string
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdAtencionDiagnostico", idAtencionDiagnostico),
		sql.Named("IdAtencion", idAtencion),
		sql.Named("IdPrimeraAtencion", idPrimeraAtencion),
		sql.Named("Mensaje", sql.Out{Dest: &mensaje}),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_AtencionesDiagnosticosEliminar: %w", err)
	}
	if err := validarMensajeSP(mensaje); err != nil {
		return err
	}
	return nil
}

// AgregarPlanTrabajoItem invoca usp_go_PA_PlanTrabajoItemAgregar.
func (r *primeraAtencionRepository) AgregarPlanTrabajoItem(ctx context.Context, in domain.PlanTrabajoItemAgregar) (domain.PlanTrabajoItemAgregarResult, error) {
	const procedure = `dbo.usp_go_PA_PlanTrabajoItemAgregar`

	var outID string
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdPlanTrabajoItem", sql.Out{Dest: &outID}),
		sql.Named("IdAtencion", in.IdAtencion),
		sql.Named("IdPrimeraAtencion", in.IdPrimeraAtencion),
		sql.Named("Tipo", in.Tipo),
		sql.Named("IdProducto", in.IdProducto),
		sql.Named("IdEspecialidad", in.IdEspecialidad),
		sql.Named("Codigo", in.Codigo),
		sql.Named("Nombre", in.Nombre),
		sql.Named("NumOrden", in.NumOrden),
		sql.Named("IdUsuarioRegistro", in.IdUsuarioRegistro),
	)
	if err != nil {
		return domain.PlanTrabajoItemAgregarResult{}, fmt.Errorf("calling usp_go_PA_PlanTrabajoItemAgregar: %w", err)
	}
	return parseOutputPlanTrabajoItem(outID)
}

// ListarPlanTrabajoItems invoca usp_go_PA_PlanTrabajoItemListar.
func (r *primeraAtencionRepository) ListarPlanTrabajoItems(ctx context.Context, idAtencion, idPrimeraAtencion int, tipo string) ([]domain.PlanTrabajoItem, error) {
	const procedure = `dbo.usp_go_PA_PlanTrabajoItemListar`

	var idPrimera any
	if idPrimeraAtencion > 0 {
		idPrimera = idPrimeraAtencion
	} else {
		idPrimera = nil
	}
	var tipoArg any
	if strings.TrimSpace(tipo) != "" {
		tipoArg = strings.ToUpper(strings.TrimSpace(tipo))
	} else {
		tipoArg = nil
	}

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("IdAtencion", idAtencion),
		sql.Named("IdPrimeraAtencion", idPrimera),
		sql.Named("Tipo", tipoArg),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_PA_PlanTrabajoItemListar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading plan trabajo items: %w", err)
	}

	items := make([]domain.PlanTrabajoItem, 0, len(maps))
	for _, m := range maps {
		id := rowInt64(m, "IdPlanTrabajoItem", "idPlanTrabajoItem")
		if id == nil || *id <= 0 {
			continue
		}
		item := domain.PlanTrabajoItem{IdPlanTrabajoItem: int(*id)}
		if v := rowInt64(m, "IdAtencion", "idAtencion"); v != nil {
			item.IdAtencion = int(*v)
		}
		if v := rowInt64(m, "IdPrimeraAtencion", "idPrimeraAtencion"); v != nil {
			item.IdPrimeraAtencion = int(*v)
		}
		if s := rowString(m, "Tipo", "tipo"); s != nil {
			item.Tipo = strings.ToUpper(strings.TrimSpace(*s))
		}
		if v := rowInt64(m, "IdProducto", "idProducto"); v != nil && *v > 0 {
			iv := int(*v)
			item.IdProducto = &iv
		}
		if v := rowInt64(m, "IdEspecialidad", "idEspecialidad"); v != nil && *v > 0 {
			iv := int(*v)
			item.IdEspecialidad = &iv
		}
		if s := rowString(m, "Codigo", "codigo"); s != nil {
			item.Codigo = strings.TrimSpace(*s)
		}
		if s := rowString(m, "Nombre", "nombre"); s != nil {
			item.Nombre = strings.TrimSpace(*s)
		}
		if v := rowInt64(m, "NumOrden", "numOrden"); v != nil {
			item.NumOrden = int(*v)
		}
		items = append(items, item)
	}
	return items, nil
}

// EliminarPlanTrabajoItem invoca usp_go_PA_PlanTrabajoItemEliminar.
func (r *primeraAtencionRepository) EliminarPlanTrabajoItem(ctx context.Context, idAtencion, idPrimeraAtencion, idPlanTrabajoItem int) error {
	const procedure = `dbo.usp_go_PA_PlanTrabajoItemEliminar`

	var mensaje string
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdPlanTrabajoItem", idPlanTrabajoItem),
		sql.Named("IdAtencion", idAtencion),
		sql.Named("IdPrimeraAtencion", idPrimeraAtencion),
		sql.Named("Mensaje", sql.Out{Dest: &mensaje}),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_PA_PlanTrabajoItemEliminar: %w", err)
	}
	if err := validarMensajeSP(mensaje); err != nil {
		return err
	}
	return nil
}

func parseOutputPlanTrabajoItem(raw string) (domain.PlanTrabajoItemAgregarResult, error) {
	out := domain.PlanTrabajoItemAgregarResult{Raw: strings.TrimSpace(raw)}
	m := out.Raw
	if m == "" {
		return out, fmt.Errorf("respuesta vacía del SP de plan de trabajo")
	}
	upper := strings.ToUpper(m)
	if strings.HasPrefix(upper, "ERROR") {
		parts := strings.SplitN(m, ";", 2)
		if len(parts) == 2 {
			return out, fmt.Errorf("%s", strings.TrimSpace(parts[1]))
		}
		return out, fmt.Errorf("%s", m)
	}
	// Formato esperado: OK;<id>
	parts := strings.SplitN(m, ";", 2)
	idRaw := strings.TrimSpace(parts[0])
	if strings.EqualFold(idRaw, "OK") && len(parts) == 2 {
		idRaw = strings.TrimSpace(parts[1])
	}
	id, err := strconv.Atoi(idRaw)
	if err != nil || id <= 0 {
		return out, fmt.Errorf("OUTPUT inválido del SP: %s", m)
	}
	out.IdPlanTrabajoItem = id
	return out, nil
}

func idsAtencionEnTabPrimeraAtencion(ctx context.Context, db *sql.DB, idsAtencion []int, soloFirmados bool) (map[int]bool, error) {
	result := make(map[int]bool)
	if len(idsAtencion) == 0 {
		return result, nil
	}

	unique := make([]int, 0, len(idsAtencion))
	seen := make(map[int]struct{}, len(idsAtencion))
	for _, id := range idsAtencion {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	if len(unique) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(unique))
	args := make([]any, len(unique))
	for i, id := range unique {
		pname := fmt.Sprintf("p%d", i+1)
		placeholders[i] = "@" + pname
		args[i] = sql.Named(pname, id)
	}

	filtroFirma := ""
	if soloFirmados {
		filtroFirma = " AND ISNULL(EstadoFirma, 0) = 1"
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT IdAtencion
		FROM dbo.Tab_PrimeraAtencion
		WHERE IdAtencion IN (%s)%s
	`, strings.Join(placeholders, ","), filtroFirma)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying Tab_PrimeraAtencion ids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scanning Tab_PrimeraAtencion id: %w", err)
		}
		result[id] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating Tab_PrimeraAtencion ids: %w", err)
	}

	return result, nil
}

func parseOutputAtencionDiagnostico(raw string) (domain.DiagnosticoAgregarResult, error) {
	out := domain.DiagnosticoAgregarResult{Raw: strings.TrimSpace(raw)}
	m := out.Raw
	if m == "" {
		return out, fmt.Errorf("respuesta vacía del SP de diagnósticos")
	}
	upper := strings.ToUpper(m)
	if strings.HasPrefix(upper, "ERROR") {
		parts := strings.SplitN(m, ";", 2)
		if len(parts) == 2 {
			return out, fmt.Errorf("%s", strings.TrimSpace(parts[1]))
		}
		return out, fmt.Errorf("%s", m)
	}
	if m == "-1" {
		return out, fmt.Errorf("diagnóstico duplicado o parámetros inválidos")
	}
	parts := strings.SplitN(m, ";", 2)
	id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || id <= 0 {
		return out, fmt.Errorf("OUTPUT inválido del SP: %s", m)
	}
	out.IdAtencionDiagnostico = id
	if len(parts) == 2 {
		out.SubclasificacionAnterior = strings.TrimSpace(parts[1])
	}
	return out, nil
}

func validarMensajeSP(mensaje string) error {
	m := strings.TrimSpace(mensaje)
	if m == "" {
		return nil
	}
	upper := strings.ToUpper(m)
	if strings.HasPrefix(upper, "ERROR") {
		parts := strings.SplitN(m, ";", 2)
		if len(parts) == 2 {
			return fmt.Errorf("%s", strings.TrimSpace(parts[1]))
		}
		return fmt.Errorf("%s", m)
	}
	return nil
}

// ParseFechaHoraLocal interpreta ISO/RFC3339 del front como time local.
func ParseFechaHoraLocal(raw string) (*time.Time, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, nil
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	var parsed time.Time
	var err error
	for _, layout := range layouts {
		parsed, err = time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			return &parsed, nil
		}
	}
	return nil, fmt.Errorf("fecha inválida: %s", raw)
}
