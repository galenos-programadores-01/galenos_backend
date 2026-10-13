package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// RegistrarAltaHospitalaria persiste el egreso vía usp_go_RegistrarAltaHospitalaria.
// Domicilio: IdTipoAlta + IdCondicionAlta obligatorios.
// Hosp/Fuga/Fallecido: IdDestinoAtencion (catálogo HOS/FUG/FAL); tipo/condición NULL.
func (uc *primeraAtencionUseCase) RegistrarAltaHospitalaria(
	ctx context.Context,
	in domain.AltaHospitalariaRegistrar,
) (domain.AltaHospitalariaResult, error) {
	if in.IdAtencion <= 0 {
		return domain.AltaHospitalariaResult{}, fmt.Errorf("idAtencion inválido")
	}
	if in.IdMedicoEgreso <= 0 {
		return domain.AltaHospitalariaResult{}, fmt.Errorf("idMedicoEgreso inválido")
	}
	if in.FechaEgreso.IsZero() {
		return domain.AltaHospitalariaResult{}, fmt.Errorf("fechaEgreso inválida")
	}
	if in.IdServicioEgreso <= 0 {
		return domain.AltaHospitalariaResult{}, fmt.Errorf("idServicioEgreso inválido")
	}
	if in.IdEpisodio == nil || *in.IdEpisodio <= 0 {
		return domain.AltaHospitalariaResult{}, fmt.Errorf("idEpisodio inválido")
	}

	// Hosp / fuga: tipo y condición NULL; destino = Id catálogo de la tarjeta.
	sinTipoCondicion := in.IdTipoAlta == nil && in.IdCondicionAlta == nil
	if sinTipoCondicion {
		if in.IdDestinoAtencion == nil || *in.IdDestinoAtencion <= 0 {
			return domain.AltaHospitalariaResult{}, fmt.Errorf("idDestinoAtencion inválido (catálogo de egreso)")
		}
		if in.IdNotificacionFuga != nil {
			n := *in.IdNotificacionFuga
			if n < 1 || n > 3 {
				return domain.AltaHospitalariaResult{}, fmt.Errorf("idNotificacionFuga inválido (espere 1, 2 o 3)")
			}
		}
	} else {
		if in.IdTipoAlta == nil || *in.IdTipoAlta <= 0 {
			return domain.AltaHospitalariaResult{}, fmt.Errorf("idTipoAlta inválido")
		}
		if in.IdCondicionAlta == nil || *in.IdCondicionAlta <= 0 {
			return domain.AltaHospitalariaResult{}, fmt.Errorf("idCondicionAlta inválido")
		}
	}

	if err := uc.asegurarEditable(ctx, in.IdAtencion); err != nil {
		return domain.AltaHospitalariaResult{}, err
	}

	return uc.repo.RegistrarAltaHospitalaria(ctx, in)
}
