package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

func (uc *primeraAtencionUseCase) asegurarEditable(ctx context.Context, idAtencion int) error {
	row, err := uc.repo.ObtenerPorIdAtencion(ctx, idAtencion)
	if err != nil {
		return fmt.Errorf("verificar estado primera atencion: %w", err)
	}
	if primeraAtencionEstaFirmada(row) {
		return domain.ErrPrimeraAtencionYaFirmada
	}
	return nil
}

func primeraAtencionEstaFirmada(row map[string]any) bool {
	if row == nil {
		return false
	}
	claves := []string{"EstadoFirma", "estadoFirma"}
	for _, key := range claves {
		if enteroDesdeValor(row[key]) == domain.EstadoFirmaCerrada {
			return true
		}
	}
	for key, value := range row {
		if strings.EqualFold(key, "EstadoFirma") && enteroDesdeValor(value) == domain.EstadoFirmaCerrada {
			return true
		}
	}
	return false
}

// IniciarBorrador crea el registro mínimo en Tab_PrimeraAtencion.
func (uc *primeraAtencionUseCase) IniciarBorrador(ctx context.Context, idAtencion, idUsuario int) (domain.PrimeraAtencionUpsertResult, error) {
	if err := uc.asegurarEditable(ctx, idAtencion); err != nil {
		return domain.PrimeraAtencionUpsertResult{}, err
	}
	now := time.Now()
	tipo := "D"
	in := domain.PrimeraAtencionUpsert{
		IdAtencion:        idAtencion,
		FechaAtencion:     &now,
		TipoAnamnesis:     &tipo,
		IdUsuarioRegistro: idUsuario,
	}
	result, err := uc.repo.GuardarUpsert(ctx, in)
	if err != nil {
		return domain.PrimeraAtencionUpsertResult{}, fmt.Errorf("iniciar borrador primera atencion: %w", err)
	}
	return result, nil
}

// GuardarParcial persiste cambios del formulario (upsert parcial vía SP).
func (uc *primeraAtencionUseCase) GuardarParcial(ctx context.Context, in domain.PrimeraAtencionUpsert) (domain.PrimeraAtencionUpsertResult, error) {
	if err := uc.asegurarEditable(ctx, in.IdAtencion); err != nil {
		return domain.PrimeraAtencionUpsertResult{}, err
	}
	result, err := uc.repo.GuardarUpsert(ctx, in)
	if err != nil {
		return domain.PrimeraAtencionUpsertResult{}, fmt.Errorf("guardar parcial primera atencion: %w", err)
	}
	return result, nil
}

// ObtenerDetalle rehidrata el borrador por IdAtencion.
func (uc *primeraAtencionUseCase) ObtenerDetalle(ctx context.Context, idAtencion int) (map[string]any, error) {
	row, err := uc.repo.ObtenerPorIdAtencion(ctx, idAtencion)
	if err != nil {
		return nil, fmt.Errorf("obtener detalle primera atencion: %w", err)
	}
	return row, nil
}

// CerrarFirma persiste UUID/ruta del PDF firmado y bloquea nuevas ediciones.
func (uc *primeraAtencionUseCase) CerrarFirma(ctx context.Context, in domain.PrimeraAtencionCerrarFirma) (domain.PrimeraAtencionCerrarFirmaResult, error) {
	if in.IdAtencion <= 0 {
		return domain.PrimeraAtencionCerrarFirmaResult{}, fmt.Errorf("idAtencion inválido")
	}
	if in.IdUsuarioModifica <= 0 {
		return domain.PrimeraAtencionCerrarFirmaResult{}, fmt.Errorf("idUsuario inválido")
	}
	in.UbicacionArchivo = strings.TrimSpace(in.UbicacionArchivo)
	if in.UbicacionArchivo == "" {
		return domain.PrimeraAtencionCerrarFirmaResult{}, fmt.Errorf("ubicacionArchivo es obligatorio")
	}

	row, err := uc.repo.ObtenerPorIdAtencion(ctx, in.IdAtencion)
	if err != nil {
		return domain.PrimeraAtencionCerrarFirmaResult{}, fmt.Errorf("verificar primera atencion: %w", err)
	}
	if row == nil {
		return domain.PrimeraAtencionCerrarFirmaResult{}, domain.ErrPrimeraAtencionSinBorrador
	}

	estadoFirma := domain.EstadoFirmaCerrada
	now := time.Now()
	result, err := uc.repo.GuardarUpsert(ctx, domain.PrimeraAtencionUpsert{
		IdAtencion:        in.IdAtencion,
		IdUsuarioRegistro: in.IdUsuarioModifica,
		IdUsuarioModifica: &in.IdUsuarioModifica,
		FechaModifica:     &now,
		EstadoFirma:       &estadoFirma,
		UbicacionArchivo:  &in.UbicacionArchivo,
	})
	if err != nil {
		if errors.Is(err, domain.ErrPrimeraAtencionYaFirmada) || strings.Contains(strings.ToLower(err.Error()), "ya fue firmada") {
			return domain.PrimeraAtencionCerrarFirmaResult{}, domain.ErrPrimeraAtencionYaFirmada
		}
		return domain.PrimeraAtencionCerrarFirmaResult{}, fmt.Errorf("cerrar firma primera atencion: %w", err)
	}
	return domain.PrimeraAtencionCerrarFirmaResult{
		IdPrimeraAtencion: result.IdPrimeraAtencion,
		Mensaje:           result.Mensaje,
	}, nil
}
