package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// RegistrarMedicamento: 1) cabecera (crear o reusar) 2) detalle.
// Mínimo: IdProducto + Cantidad + Duracion. Resto opcional.
func (uc *primeraAtencionUseCase) RegistrarMedicamento(
	ctx context.Context,
	in domain.RecetaMedicamentoRegistrar,
) (domain.RecetaMedicamentoRegistrarResult, error) {
	// IdPuntoCarga de entrada se ignora: siempre farmacia (HIS = 5).
	if in.IdAtencion <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("idAtencion inválido")
	}
	if in.IdCuentaAtencion <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("idCuentaAtencion inválido")
	}
	if in.IdPrimeraAtencion <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("idPrimeraAtencion inválido")
	}
	if in.IdServicioReceta <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("idServicioReceta inválido")
	}
	if in.IdMedicoReceta <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("idMedicoReceta inválido")
	}
	if in.IdPaciente <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("idPaciente inválido")
	}
	if in.IdProducto <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("idProducto inválido")
	}
	if in.Cantidad <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("cantidad inválida")
	}
	if in.Duracion <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("duracion (días) inválida")
	}
	if in.IdDiagnostico <= 0 {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("idDiagnostico inválido: el medicamento debe asociarse a un CIE")
	}
	codigoDx := strings.TrimSpace(in.CodigoDiagnostico)
	if codigoDx == "" {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("codigoDiagnostico inválido: se requiere el CIE-10 del diagnóstico")
	}
	if err := uc.asegurarEditable(ctx, in.IdAtencion); err != nil {
		return domain.RecetaMedicamentoRegistrarResult{}, err
	}

	idReceta, err := uc.repo.AgregarRecetaCabecera(ctx, domain.RecetaCabeceraAgregar{
		// Farmacia fija (HIS). No usar IdPuntoCarga de bandeja (puede ser otro servicio).
		IdPuntoCarga:       domain.IdPuntoCargaFarmacia,
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
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("cabecera receta: %w", err)
	}
	if err := uc.rechazarCabeceraAnulada(ctx, idReceta); err != nil {
		return domain.RecetaMedicamentoRegistrarResult{}, err
	}

	precio := in.Precio
	if precio < 0 {
		precio = 0
	}
	total := float64(in.Cantidad) * precio

	idDx := in.IdDiagnostico
	err = uc.repo.AgregarRecetaDetalle(ctx, domain.RecetaDetalleAgregar{
		IdReceta:            idReceta,
		IdProducto:          in.IdProducto,
		Cantidad:            in.Cantidad,
		Precio:              precio,
		Total:               total,
		SaldoEnRegistro:     0,
		IdDosisRecetada:     in.IdDosis,
		Observaciones:       "",
		IdViaAdministracion: in.IdViaAdministracion,
		IdUNIDDosisReceta:   in.IdUNIDDosis,
		IdFrecuencia:        in.IdFrecuencia,
		Duracion:            in.Duracion,
		IdDiagnostico:       &idDx,
		CodigoDiagnostico:   &codigoDx,
	})
	if err != nil {
		return domain.RecetaMedicamentoRegistrarResult{}, fmt.Errorf("detalle receta: %w", err)
	}

	return domain.RecetaMedicamentoRegistrarResult{
		IdReceta:   idReceta,
		IdProducto: in.IdProducto,
		Cantidad:   in.Cantidad,
		Precio:     precio,
		Total:      total,
		Mensaje:    "Medicamento registrado",
	}, nil
}

// ActualizarMedicamentoCampos: dosis/unidad/frecuencia/vía de un detalle ya guardado.
func (uc *primeraAtencionUseCase) ActualizarMedicamentoCampos(
	ctx context.Context,
	in domain.RecetaMedicamentoCamposActualizar,
) error {
	if in.IdAtencion <= 0 {
		return fmt.Errorf("idAtencion inválido")
	}
	if in.IdReceta <= 0 {
		return fmt.Errorf("idReceta inválido")
	}
	if in.IdProducto <= 0 {
		return fmt.Errorf("idProducto inválido")
	}
	if in.IdCuentaAtencion <= 0 {
		return fmt.Errorf("idCuentaAtencion inválido")
	}
	if in.IdPrimeraAtencion <= 0 {
		return fmt.Errorf("idPrimeraAtencion inválido")
	}
	if err := uc.asegurarEditable(ctx, in.IdAtencion); err != nil {
		return err
	}
	if err := uc.repo.ActualizarRecetaDetalleCamposClinicos(ctx, in); err != nil {
		return fmt.Errorf("actualizar campos medicamento: %w", err)
	}
	return nil
}

func (uc *primeraAtencionUseCase) ListarMedicamentos(
	ctx context.Context,
	idCuentaAtencion, idPrimeraAtencion, idPuntoCarga int,
) ([]domain.RecetaMedicamentoItem, error) {
	if idCuentaAtencion <= 0 {
		return nil, fmt.Errorf("idCuentaAtencion inválido")
	}
	// Siempre farmacia (=5). El idPuntoCarga del query/bandeja se ignora:
	// la bandeja suele traer otro punto (consultorio) y dejaba el listado vacío.
	_ = idPuntoCarga
	items, err := uc.repo.ListarRecetaDetalle(ctx, domain.RecetaDetalleListarFiltro{
		IdCuentaAtencion:  idCuentaAtencion,
		IdPrimeraAtencion: idPrimeraAtencion,
		IdPuntoCarga:      domain.IdPuntoCargaFarmacia,
	})
	if err != nil {
		return nil, fmt.Errorf("listar medicamentos: %w", err)
	}
	return items, nil
}

func (uc *primeraAtencionUseCase) EliminarMedicamento(
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
		return fmt.Errorf("eliminar medicamento: %w", err)
	}
	return nil
}

func (uc *primeraAtencionUseCase) rechazarCabeceraAnulada(ctx context.Context, idReceta int) error {
	anulada, err := uc.repo.CabeceraRecetaAnulada(ctx, idReceta)
	if err != nil {
		return fmt.Errorf("verificar cabecera receta: %w", err)
	}
	if anulada {
		return fmt.Errorf("la cabecera de receta anulada no se reutiliza")
	}
	return nil
}

func (uc *primeraAtencionUseCase) eliminarDetalleReceta(ctx context.Context, idAtencion, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion int) error {
	row, err := uc.repo.ObtenerPorIdAtencion(ctx, idAtencion)
	if err != nil {
		return err
	}
	if primeraAtencionEstaFirmada(row) {
		return uc.repo.EliminarDetalleRectificacion(ctx, idReceta, idProducto)
	}
	return uc.repo.EliminarRecetaDetalle(ctx, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion)
}