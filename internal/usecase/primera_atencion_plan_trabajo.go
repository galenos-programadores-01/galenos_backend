package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// RegistrarPlanTrabajoItem inserta 1 ítem (EXA|INT|PRO) vía SP temporal.
// Misma UX que diagnósticos: al seleccionar en UI.
func (uc *primeraAtencionUseCase) RegistrarPlanTrabajoItem(
	ctx context.Context,
	idAtencion, idPrimeraAtencion int,
	item domain.PlanTrabajoItemRegistrarItem,
) (domain.PlanTrabajoItemAgregarResult, error) {
	if idAtencion <= 0 {
		return domain.PlanTrabajoItemAgregarResult{}, fmt.Errorf("idAtencion inválido")
	}
	if idPrimeraAtencion <= 0 {
		return domain.PlanTrabajoItemAgregarResult{}, fmt.Errorf("idPrimeraAtencion inválido")
	}

	tipo := strings.ToUpper(strings.TrimSpace(item.Tipo))
	switch tipo {
	case domain.PlanTrabajoTipoExamen, domain.PlanTrabajoTipoProcedimiento, domain.PlanTrabajoTipoMedicamento:
		if item.IdProducto == nil || *item.IdProducto <= 0 {
			return domain.PlanTrabajoItemAgregarResult{}, fmt.Errorf("idProducto obligatorio para tipo %s", tipo)
		}
	case domain.PlanTrabajoTipoInterconsulta:
		if item.IdEspecialidad == nil || *item.IdEspecialidad <= 0 {
			return domain.PlanTrabajoItemAgregarResult{}, fmt.Errorf("idEspecialidad obligatorio para tipo INT")
		}
	default:
		return domain.PlanTrabajoItemAgregarResult{}, fmt.Errorf("tipo inválido (use EXA|INT|PRO)")
	}

	if err := uc.asegurarEditable(ctx, idAtencion); err != nil {
		return domain.PlanTrabajoItemAgregarResult{}, err
	}

	in := domain.PlanTrabajoItemAgregar{
		IdAtencion:        idAtencion,
		IdPrimeraAtencion: idPrimeraAtencion,
		Tipo:              tipo,
		IdProducto:        item.IdProducto,
		IdEspecialidad:    item.IdEspecialidad,
	}
	if c := strings.TrimSpace(item.Codigo); c != "" {
		in.Codigo = &c
	}
	if n := strings.TrimSpace(item.Nombre); n != "" {
		in.Nombre = &n
	}

	res, err := uc.repo.AgregarPlanTrabajoItem(ctx, in)
	if err != nil {
		return domain.PlanTrabajoItemAgregarResult{}, fmt.Errorf("agregar plan trabajo: %w", err)
	}
	return res, nil
}

// ListarPlanTrabajoItems rehidrata ítems (otra máquina / Continuar).
func (uc *primeraAtencionUseCase) ListarPlanTrabajoItems(
	ctx context.Context,
	idAtencion, idPrimeraAtencion int,
	tipo string,
) ([]domain.PlanTrabajoItem, error) {
	if idAtencion <= 0 {
		return nil, fmt.Errorf("idAtencion inválido")
	}
	items, err := uc.repo.ListarPlanTrabajoItems(ctx, idAtencion, idPrimeraAtencion, tipo)
	if err != nil {
		return nil, fmt.Errorf("listar plan trabajo: %w", err)
	}
	return items, nil
}

// EliminarPlanTrabajoItem quita un ítem (X en UI / cambio de interconsulta).
func (uc *primeraAtencionUseCase) EliminarPlanTrabajoItem(
	ctx context.Context,
	idAtencion, idPrimeraAtencion, idPlanTrabajoItem int,
) error {
	if idAtencion <= 0 {
		return fmt.Errorf("idAtencion inválido")
	}
	if idPrimeraAtencion <= 0 {
		return fmt.Errorf("idPrimeraAtencion inválido")
	}
	if idPlanTrabajoItem <= 0 {
		return fmt.Errorf("idPlanTrabajoItem inválido")
	}
	if err := uc.asegurarEditable(ctx, idAtencion); err != nil {
		return err
	}
	if err := uc.repo.EliminarPlanTrabajoItem(ctx, idAtencion, idPrimeraAtencion, idPlanTrabajoItem); err != nil {
		return fmt.Errorf("eliminar plan trabajo: %w", err)
	}
	return nil
}
