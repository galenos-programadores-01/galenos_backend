package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// RegistrarSolicitudInterconsulta: Insertar cabecera → N× DX.
// Si falla un DX a mitad, compensa eliminando la IC (todo o nada).
func (uc *primeraAtencionUseCase) RegistrarSolicitudInterconsulta(
	ctx context.Context,
	in domain.InterconsultaRegistrarSolicitud,
) (int, error) {
	if in.IdAtencion <= 0 {
		return 0, fmt.Errorf("idAtencion inválido")
	}
	if in.IdEpisodio <= 0 {
		return 0, fmt.Errorf("idEpisodio es obligatorio")
	}
	if in.IdPrimeraAtencion <= 0 {
		return 0, fmt.Errorf("idPrimeraAtencion es obligatorio")
	}
	if in.IdServicio <= 0 {
		return 0, fmt.Errorf("idServicio es obligatorio")
	}
	if in.IdEspecialidad <= 0 {
		return 0, fmt.Errorf("idEspecialidad es obligatorio")
	}
	if in.IdPrioridad <= 0 {
		return 0, fmt.Errorf("idPrioridad es obligatorio")
	}
	if in.IdEmpleado <= 0 {
		return 0, fmt.Errorf("idEmpleado es obligatorio")
	}
	motivo := strings.TrimSpace(in.MotivoClinico)
	if motivo == "" {
		return 0, fmt.Errorf("motivoClinico es obligatorio")
	}
	opinion, manejo, transf, okBits := domain.BitsMotivoPorOrdenCatalogo(in.MotivoOrden)
	if !okBits {
		return 0, fmt.Errorf("motivoOrden inválido (espere 0, 1 o 2 según orden del catálogo IdGrupo=%d)", domain.IdGrupoMotivoInterconsulta)
	}
	if in.IdMotivoParametro <= 0 {
		return 0, fmt.Errorf("idMotivoParametro es obligatorio")
	}
	if len(in.Diagnosticos) < 1 {
		return 0, fmt.Errorf("debe incluir al menos un diagnóstico")
	}

	seen := make(map[int]struct{}, len(in.Diagnosticos))
	for _, d := range in.Diagnosticos {
		if d.IdDiagnostico <= 0 {
			return 0, fmt.Errorf("idDiagnostico inválido")
		}
		if d.IdSubclasificacionDx <= 0 {
			return 0, fmt.Errorf("idSubclasificacionDx es obligatorio")
		}
		if _, dup := seen[d.IdDiagnostico]; dup {
			return 0, fmt.Errorf("diagnóstico duplicado en la solicitud: %d", d.IdDiagnostico)
		}
		seen[d.IdDiagnostico] = struct{}{}
	}

	if err := uc.asegurarEditable(ctx, in.IdAtencion); err != nil {
		return 0, err
	}

	idReg := domain.ResolverIdRegistroAtencionInterconsulta(in.IdAtencion)

	idIC, err := uc.repo.InsertarInterconsulta(ctx, domain.InterconsultaInsertar{
		IdEpisodio:               in.IdEpisodio,
		IdRegistroAtencion:       idReg,
		IdPrimeraAtencion:        in.IdPrimeraAtencion,
		IdServicioSolicitante:    in.IdServicio,
		IdEspecialidadSolicitada: in.IdEspecialidad,
		IdEmpleadoSolicitante:    in.IdEmpleado,
		IdPrioridad:              in.IdPrioridad,
		OpinionDiagnostica:       opinion,
		ManejoPaciente:           manejo,
		TrasferenciaPaciente:     transf,
		MotivoClinico:            motivo,
		ResumenClinico:           "",
		IdEstado:                 domain.IdEstadoInterconsultaSolicitado,
		IdUsuarioRegistro:        in.IdEmpleado,
	})
	if err != nil {
		return 0, fmt.Errorf("insertar interconsulta: %w", err)
	}

	for _, d := range in.Diagnosticos {
		errDx := uc.repo.InsertarInterconsultaDiagnostico(ctx, domain.InterconsultaDiagnosticoInsertar{
			IdInterconsulta:      idIC,
			IdDiagnostico:        d.IdDiagnostico,
			IdSubclasificacionDx: d.IdSubclasificacionDx,
			IdTipo:               domain.IdTipoSolicitudDxInterconsulta,
			IdEstado:             domain.IdEstadoInterconsultaSolicitado,
		})
		if errDx != nil {
			_ = uc.repo.EliminarInterconsulta(ctx, idIC, in.IdEpisodio)
			return 0, fmt.Errorf("insertar diagnóstico %d: %w", d.IdDiagnostico, errDx)
		}
	}

	return idIC, nil
}

// ListarInterconsultas rehidrata la tabla inferior del paso 5.
// PA: IdRegistroAtencion=0/NULL → el List SP engrana por IdPrimeraAtencion
// (IdEpisodio se sigue enviando).
func (uc *primeraAtencionUseCase) ListarInterconsultas(
	ctx context.Context,
	idAtencion, idEpisodio, idPrimeraAtencion int,
) ([]domain.InterconsultaListadoItem, error) {
	if idAtencion <= 0 {
		return nil, fmt.Errorf("idAtencion inválido")
	}
	if idEpisodio <= 0 {
		return nil, fmt.Errorf("idEpisodio es obligatorio para listar interconsultas PA")
	}
	if idPrimeraAtencion <= 0 {
		return nil, fmt.Errorf("idPrimeraAtencion es obligatorio para listar interconsultas PA")
	}
	items, err := uc.repo.ListarInterconsultasPorAtencion(ctx, idAtencion, idEpisodio, idPrimeraAtencion)
	if err != nil {
		return nil, fmt.Errorf("listar interconsultas: %w", err)
	}
	return items, nil
}

// EliminarSolicitudInterconsulta: cascada DX → cabecera (solo SOLICITUD en SP).
func (uc *primeraAtencionUseCase) EliminarSolicitudInterconsulta(
	ctx context.Context,
	idAtencion, idInterconsulta, idEpisodio int,
) error {
	if idAtencion <= 0 {
		return fmt.Errorf("idAtencion inválido")
	}
	if idInterconsulta <= 0 {
		return fmt.Errorf("idInterconsulta inválido")
	}
	if idEpisodio <= 0 {
		return fmt.Errorf("idEpisodio es obligatorio")
	}
	if err := uc.asegurarEditable(ctx, idAtencion); err != nil {
		return err
	}
	if err := uc.repo.EliminarInterconsulta(ctx, idInterconsulta, idEpisodio); err != nil {
		return fmt.Errorf("eliminar interconsulta: %w", err)
	}
	return nil
}
