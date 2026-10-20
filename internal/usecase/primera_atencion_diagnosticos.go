package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// RegistrarDiagnosticos inserta ≥1 diagnóstico vía usp_go_AtencionesDiagnosticosAgregar.
// IdSubclasificacionDx debe venir del catálogo SP (no se inventan ids 401–406 en código).
// IdEpisodio e IdClasificacionDx son obligatorios (PA=9, egreso=3, fallecido=4, IC=7).
//
// IdPrimeraAtencion:
//   - PA (9) / IC (7): obligatorio (>0) — amarre a Tab_PrimeraAtencion.
//   - Egreso (3) / egreso fallecido (4): siempre NULL (0) — no pertenecen a la PA.
// Unicidad de Principal/Secundario en egreso se valida en UI (catálogo clasif. 3/4).
func (uc *primeraAtencionUseCase) RegistrarDiagnosticos(
	ctx context.Context,
	idAtencion, idPaciente, idPrimeraAtencion, idEpisodio, idClasificacionDx int,
	items []domain.DiagnosticoRegistrarItem,
) ([]domain.DiagnosticoAgregarResult, error) {
	if idAtencion <= 0 {
		return nil, fmt.Errorf("idAtencion inválido")
	}
	if idPaciente <= 0 {
		return nil, fmt.Errorf("idPaciente inválido")
	}
	if idEpisodio <= 0 {
		return nil, fmt.Errorf("idEpisodio inválido")
	}
	if idClasificacionDx <= 0 {
		return nil, fmt.Errorf("idClasificacionDx inválido")
	}
	if len(items) < 1 {
		return nil, fmt.Errorf("debe registrar al menos un diagnóstico")
	}

	esEgreso := idClasificacionDx == domain.ClasificacionDxEgreso ||
		idClasificacionDx == domain.ClasificacionDxEgresoFallecido
	idPrimeraPersistir := idPrimeraAtencion
	if esEgreso {
		// DX de egreso no se amarran a Primera Atención.
		idPrimeraPersistir = 0
	} else if idPrimeraAtencion <= 0 {
		return nil, fmt.Errorf("idPrimeraAtencion es obligatorio para diagnósticos de primera atención / interconsulta")
	}

	if err := uc.asegurarEditable(ctx, idAtencion); err != nil {
		return nil, err
	}

	seen := make(map[int]struct{}, len(items))
	resultados := make([]domain.DiagnosticoAgregarResult, 0, len(items))

	for _, item := range items {
		if item.IdDiagnostico <= 0 {
			return nil, fmt.Errorf("idDiagnostico inválido")
		}
		if item.IdSubclasificacionDx <= 0 {
			return nil, fmt.Errorf("idSubclasificacionDx es obligatorio (catálogo SubclasificacionDiagnosticos)")
		}
		if _, dup := seen[item.IdDiagnostico]; dup {
			return nil, fmt.Errorf("diagnóstico duplicado en la solicitud: %d", item.IdDiagnostico)
		}
		seen[item.IdDiagnostico] = struct{}{}

		in := domain.DiagnosticoAgregar{
			IdAtencion:                    idAtencion,
			IdDiagnostico:                 item.IdDiagnostico,
			IdClasificacionDx:             idClasificacionDx,
			IdSubclasificacionDiagnostico: item.IdSubclasificacionDx,
			IdPaciente:                    idPaciente,
			IdPrimeraAtencion:             idPrimeraPersistir,
			IdEpisodio:                    idEpisodio,
		}
		res, err := uc.repo.AgregarDiagnostico(ctx, in)
		if err != nil {
			return resultados, fmt.Errorf("agregar diagnóstico %d: %w", item.IdDiagnostico, err)
		}
		resultados = append(resultados, res)
	}

	return resultados, nil
}

// ListarDiagnosticos rehidrata DX desde AtencionesDiagnosticos (otra máquina / Continuar).
func (uc *primeraAtencionUseCase) ListarDiagnosticos(
	ctx context.Context,
	idAtencion, idPrimeraAtencion int,
) ([]domain.DiagnosticoAtencion, error) {
	if idAtencion <= 0 {
		return nil, fmt.Errorf("idAtencion inválido")
	}
	items, err := uc.repo.ListarDiagnosticos(ctx, idAtencion, idPrimeraAtencion)
	if err != nil {
		return nil, fmt.Errorf("listar diagnosticos primera atencion: %w", err)
	}
	return items, nil
}

// EliminarDiagnostico quita una fila de AtencionesDiagnosticos (X en UI).
// Exige idPrimeraAtencion para no borrar DX de otros módulos / otras PA.
func (uc *primeraAtencionUseCase) EliminarDiagnostico(
	ctx context.Context,
	idAtencion, idPrimeraAtencion, idAtencionDiagnostico int,
) error {
	if idAtencion <= 0 {
		return fmt.Errorf("idAtencion inválido")
	}
	if idPrimeraAtencion <= 0 {
		return fmt.Errorf("idPrimeraAtencion inválido")
	}
	if idAtencionDiagnostico <= 0 {
		return fmt.Errorf("idAtencionDiagnostico inválido")
	}
	if err := uc.asegurarEditable(ctx, idAtencion); err != nil {
		return err
	}
	if err := uc.repo.EliminarDiagnostico(ctx, idAtencion, idPrimeraAtencion, idAtencionDiagnostico); err != nil {
		return fmt.Errorf("eliminar diagnostico primera atencion: %w", err)
	}
	return nil
}
