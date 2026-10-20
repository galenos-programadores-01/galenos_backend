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

// InsertarInterconsulta invoca usp_go_Interconsulta_Insertar.
// OUTPUT @Resultado: "OK;{IdInterconsulta}" | "ERROR: …"
func (r *primeraAtencionRepository) InsertarInterconsulta(
	ctx context.Context,
	in domain.InterconsultaInsertar,
) (int, error) {
	const procedure = `dbo.usp_go_Interconsulta_Insertar`

	var resultado string
	var fecha any
	if in.FechaSolicitud != nil && strings.TrimSpace(*in.FechaSolicitud) != "" {
		fecha = strings.TrimSpace(*in.FechaSolicitud)
	} else {
		fecha = nil
	}

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdEpisodio", in.IdEpisodio),
		sql.Named("IdRegistroAtencion", in.IdRegistroAtencion),
		sql.Named("IdPrimeraAtencion", in.IdPrimeraAtencion),
		sql.Named("IdServicioSolicitante", in.IdServicioSolicitante),
		sql.Named("IdEspecialidadSolicitada", in.IdEspecialidadSolicitada),
		sql.Named("IdEmpleadoSolicitante", in.IdEmpleadoSolicitante),
		sql.Named("FechaSolicitud", fecha),
		sql.Named("IdPrioridad", in.IdPrioridad),
		sql.Named("OpinionDiagnostica", in.OpinionDiagnostica),
		sql.Named("ManejoPaciente", in.ManejoPaciente),
		sql.Named("TrasferenciaPaciente", in.TrasferenciaPaciente),
		sql.Named("MotivoClinico", in.MotivoClinico),
		sql.Named("ResumenClinico", in.ResumenClinico),
		sql.Named("IdEstado", in.IdEstado),
		sql.Named("IdUsuarioRegistro", in.IdUsuarioRegistro),
		sql.Named("IdCama", nil),
		sql.Named("EsPreQuirurgica", nil),
		sql.Named("Resultado", sql.Out{Dest: &resultado}),
	)
	if err != nil {
		return 0, fmt.Errorf("calling usp_go_Interconsulta_Insertar: %w", err)
	}
	id, err := parseResultadoOKID(resultado)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// InsertarInterconsultaDiagnostico invoca usp_go_InterconsultaDiagnostico_Insertar.
func (r *primeraAtencionRepository) InsertarInterconsultaDiagnostico(
	ctx context.Context,
	in domain.InterconsultaDiagnosticoInsertar,
) error {
	const procedure = `dbo.usp_go_InterconsultaDiagnostico_Insertar`

	var resultado string
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdInterconsulta", in.IdInterconsulta),
		sql.Named("IdDiagnostico", in.IdDiagnostico),
		sql.Named("IdSubclasificacionDx", in.IdSubclasificacionDx),
		sql.Named("IdTipo", in.IdTipo),
		sql.Named("IdEstado", in.IdEstado),
		sql.Named("Resultado", sql.Out{Dest: &resultado}),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_InterconsultaDiagnostico_Insertar: %w", err)
	}
	return validarResultadoOK(resultado)
}

// ListarInterconsultasPorAtencion invoca usp_go_ListarInterconsultasPorAtencion.
// PA: IdRegistroAtencion=0/NULL → engrana por @IdPrimeraAtencion
// (también se envía @IdEpisodio).
func (r *primeraAtencionRepository) ListarInterconsultasPorAtencion(
	ctx context.Context,
	idAtencion, idEpisodio, idPrimeraAtencion int,
) ([]domain.InterconsultaListadoItem, error) {
	const query = `
EXEC dbo.usp_go_ListarInterconsultasPorAtencion
  @IdAtencion = @IdAtencion,
  @IdEpisodio = @IdEpisodio,
  @IdPrimeraAtencion = @IdPrimeraAtencion`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("IdAtencion", idAtencion),
		sql.Named("IdEpisodio", idEpisodio),
		sql.Named("IdPrimeraAtencion", idPrimeraAtencion),
	)
	if err != nil {
		return nil, fmt.Errorf("calling usp_go_ListarInterconsultasPorAtencion: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading interconsultas list: %w", err)
	}

	items := make([]domain.InterconsultaListadoItem, 0, len(maps))
	for _, m := range maps {
		idIC := rowInt64(m, "idInterconsulta", "IdInterconsulta")
		if idIC == nil || *idIC <= 0 {
			continue
		}
		item := domain.InterconsultaListadoItem{
			IdInterconsulta: int(*idIC),
		}
		if v := rowInt64(m, "IdAtencionOrigen", "idAtencionOrigen"); v != nil {
			item.IdAtencionOrigen = int(*v)
		}
		if v := rowInt64(m, "IdEspecialidad", "idEspecialidad"); v != nil {
			item.IdEspecialidad = int(*v)
		}
		if v := rowInt64(m, "IdMedicoDestino", "idMedicoDestino"); v != nil {
			item.IdMedicoDestino = int(*v)
		}
		item.Motivo = derefStr(rowString(m, "Motivo", "motivo"))
		item.Estado = derefStr(rowString(m, "Estado", "estado"))
		item.Especialidad = derefStr(rowString(m, "Especialidad", "especialidad"))
		item.ResumenHC = derefStr(rowString(m, "ResumenHC", "resumenHC"))
		item.FechaSolicitud = rowFechaISO(m, "FechaSolicitud", "fechaSolicitud")
		items = append(items, item)
	}
	return items, nil
}

// EliminarInterconsulta invoca usp_go_Interconsulta_Eliminar (cascada DX → cabecera).
func (r *primeraAtencionRepository) EliminarInterconsulta(
	ctx context.Context,
	idInterconsulta, idEpisodio int,
) error {
	const procedure = `dbo.usp_go_Interconsulta_Eliminar`

	var resultado string
	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdInterconsulta", idInterconsulta),
		sql.Named("IdEpisodio", idEpisodio),
		sql.Named("Resultado", sql.Out{Dest: &resultado}),
	)
	if err != nil {
		return fmt.Errorf("calling usp_go_Interconsulta_Eliminar: %w", err)
	}
	return validarResultadoOK(resultado)
}

func parseResultadoOKID(raw string) (int, error) {
	m := strings.TrimSpace(raw)
	if m == "" {
		return 0, fmt.Errorf("respuesta vacía del SP de interconsulta")
	}
	upper := strings.ToUpper(m)
	if strings.HasPrefix(upper, "ERROR") {
		return 0, fmt.Errorf("%s", mensajeErrorSP(m))
	}
	parts := strings.SplitN(m, ";", 2)
	if len(parts) < 2 || !strings.EqualFold(strings.TrimSpace(parts[0]), "OK") {
		return 0, fmt.Errorf("OUTPUT inválido del SP: %s", m)
	}
	id, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("IdInterconsulta inválido en OUTPUT: %s", m)
	}
	return id, nil
}

func validarResultadoOK(raw string) error {
	m := strings.TrimSpace(raw)
	if m == "" {
		return fmt.Errorf("respuesta vacía del SP")
	}
	upper := strings.ToUpper(m)
	if strings.HasPrefix(upper, "ERROR") {
		return fmt.Errorf("%s", mensajeErrorSP(m))
	}
	if strings.HasPrefix(upper, "OK") {
		return nil
	}
	return fmt.Errorf("OUTPUT inesperado del SP: %s", m)
}

func mensajeErrorSP(raw string) string {
	m := strings.TrimSpace(raw)
	for _, sep := range []string{":", ";"} {
		if i := strings.Index(m, sep); i >= 0 {
			rest := strings.TrimSpace(m[i+1:])
			if rest != "" {
				return rest
			}
		}
	}
	return m
}

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return strings.TrimSpace(*p)
}

func rowFechaISO(m map[string]any, keys ...string) string {
	for _, k := range keys {
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		switch t := v.(type) {
		case time.Time:
			return t.Format("2006-01-02T15:04:05")
		case string:
			return strings.TrimSpace(t)
		}
	}
	// fallback case-insensitive via rowString
	if s := rowString(m, keys...); s != nil {
		return strings.TrimSpace(*s)
	}
	return ""
}
