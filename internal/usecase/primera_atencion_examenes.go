package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

const cantidadExamenFija = 1

// RegistrarExamen: 1) cabecera PA (crear/reusar por punto + IdPrimeraAtencion)
// 2) detalle con Cantidad=1, Precio=precioUnitario, Total=Cantidad*Precio + DX.
func (uc *primeraAtencionUseCase) RegistrarExamen(
	ctx context.Context,
	in domain.RecetaExamenRegistrar,
) (domain.RecetaExamenRegistrarResult, error) {
	if in.IdAtencion <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idAtencion inválido")
	}
	if in.IdCuentaAtencion <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idCuentaAtencion inválido")
	}
	if in.IdPrimeraAtencion <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idPrimeraAtencion inválido")
	}
	if in.IdPuntoCarga <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idPuntoCarga inválido")
	}
	if in.IdServicioReceta <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idServicioReceta inválido")
	}
	if in.IdMedicoReceta <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idMedicoReceta inválido")
	}
	if in.IdPaciente <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idPaciente inválido")
	}
	if in.IdProducto <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idProducto inválido")
	}
	if in.IdDiagnostico <= 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("idDiagnostico inválido: el examen debe asociarse a un CIE")
	}
	codigoDx := strings.TrimSpace(in.CodigoDiagnostico)
	if codigoDx == "" {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("codigoDiagnostico inválido: se requiere el CIE-10 del diagnóstico")
	}
	if in.Precio < 0 {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("precio inválido")
	}
	if err := uc.asegurarEditable(ctx, in.IdAtencion); err != nil {
		return domain.RecetaExamenRegistrarResult{}, err
	}

	cantidad := cantidadExamenFija
	total := float64(cantidad) * in.Precio

	idReceta, err := uc.repo.AgregarRecetaCabecera(ctx, domain.RecetaCabeceraAgregar{
		IdPuntoCarga:       in.IdPuntoCarga,
		IdCuentaAtencion:   in.IdCuentaAtencion,
		IdServicioReceta:   in.IdServicioReceta,
		IdMedicoReceta:     in.IdMedicoReceta,
		IdProducto:         0,
		IdPaciente:         in.IdPaciente,
		IdUsuarioAuditoria: in.IdUsuarioAuditoria,
		IdEvolucion:        0,
		IdPrimeraAtencion:  in.IdPrimeraAtencion,
	})
	if err != nil {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("cabecera receta examen: %w", err)
	}
	if err := uc.rechazarCabeceraAnulada(ctx, idReceta); err != nil {
		return domain.RecetaExamenRegistrarResult{}, err
	}

	idDx := in.IdDiagnostico
	err = uc.repo.AgregarRecetaDetalle(ctx, domain.RecetaDetalleAgregar{
		IdReceta:            idReceta,
		IdProducto:          in.IdProducto,
		Cantidad:            cantidad,
		Precio:              in.Precio,
		Total:               total,
		SaldoEnRegistro:     0,
		IdDosisRecetada:     nil,
		Observaciones:       "",
		IdViaAdministracion: nil,
		IdUNIDDosisReceta:   nil,
		IdFrecuencia:        nil,
		Duracion:            0,
		IdDiagnostico:       &idDx,
		CodigoDiagnostico:   &codigoDx,
	})
	if err != nil {
		return domain.RecetaExamenRegistrarResult{}, fmt.Errorf("detalle receta examen: %w", err)
	}

	return domain.RecetaExamenRegistrarResult{
		IdReceta:   idReceta,
		IdProducto: in.IdProducto,
		Cantidad:   cantidad,
		Precio:     in.Precio,
		Total:      total,
		Mensaje:    "Examen registrado",
	}, nil
}

// ListarExamenes: SP excluye farmacia (=5). Luego se omiten procedimientos (=26)
// para no mezclarlos en la grilla de exámenes.
func (uc *primeraAtencionUseCase) ListarExamenes(
	ctx context.Context,
	idCuentaAtencion, idPrimeraAtencion int,
) ([]domain.RecetaMedicamentoItem, error) {
	if idCuentaAtencion <= 0 {
		return nil, fmt.Errorf("idCuentaAtencion inválido")
	}
	items, err := uc.repo.ListarRecetaDetalle(ctx, domain.RecetaDetalleListarFiltro{
		IdCuentaAtencion:    idCuentaAtencion,
		IdPrimeraAtencion:   idPrimeraAtencion,
		ExcluirIdPuntoCarga: domain.IdPuntoCargaFarmacia,
	})
	if err != nil {
		return nil, fmt.Errorf("listar examenes: %w", err)
	}
	out := make([]domain.RecetaMedicamentoItem, 0, len(items))
	for _, it := range items {
		if it.IdPuntoCarga == domain.IdPuntoCargaProcedimientos {
			continue
		}
		out = append(out, it)
	}
	return out, nil
}

func (uc *primeraAtencionUseCase) EliminarExamen(
	ctx context.Context,
	idAtencion, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion int,
) error {
	if idAtencion <= 0 {
		return fmt.Errorf("idAtencion inválido")
	}
	if idReceta <= 0 {
		return fmt.Errorf("idReceta inválido")
	}
	if idProducto <= 0 {
		return fmt.Errorf("idProducto inválido")
	}
	if idCuentaAtencion <= 0 {
		return fmt.Errorf("idCuentaAtencion inválido")
	}
	if idPrimeraAtencion <= 0 {
		return fmt.Errorf("idPrimeraAtencion inválido")
	}
	if err := uc.asegurarEditable(ctx, idAtencion); err != nil {
		return err
	}
	if err := uc.eliminarDetalleReceta(ctx, idAtencion, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion); err != nil {
		return fmt.Errorf("eliminar examen: %w", err)
	}
	return nil
}
