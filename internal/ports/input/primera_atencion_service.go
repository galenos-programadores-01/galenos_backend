package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// PrimeraAtencionService es el puerto de entrada (casos de uso) de Primera Atención.
type PrimeraAtencionService interface {
	ListarBandeja(ctx context.Context, params shared.PrimeraAtencionParams) ([]map[string]any, error)
	IniciarBorrador(ctx context.Context, idAtencion, idUsuario int) (domain.PrimeraAtencionUpsertResult, error)
	GuardarParcial(ctx context.Context, in domain.PrimeraAtencionUpsert) (domain.PrimeraAtencionUpsertResult, error)
	ObtenerDetalle(ctx context.Context, idAtencion int) (map[string]any, error)
	ListarDocumentosFirmados(ctx context.Context, idAtencion int) (nroVigente int, docs []domain.DocumentoFirmadoPA, err error)
	RegistrarDocumentosFirmados(ctx context.Context, idPrimeraAtencion, nroEdicion int, docs []domain.DocumentoFirmadoPA) error
	CerrarFirma(ctx context.Context, in domain.PrimeraAtencionCerrarFirma) (domain.PrimeraAtencionCerrarFirmaResult, error)
	RegistrarDiagnosticos(ctx context.Context, idAtencion, idPaciente, idPrimeraAtencion, idEpisodio, idClasificacionDx int, items []domain.DiagnosticoRegistrarItem) ([]domain.DiagnosticoAgregarResult, error)
	ListarDiagnosticos(ctx context.Context, idAtencion, idPrimeraAtencion int) ([]domain.DiagnosticoAtencion, error)
	EliminarDiagnostico(ctx context.Context, idAtencion, idPrimeraAtencion, idAtencionDiagnostico int) error
	RegistrarAltaHospitalaria(ctx context.Context, in domain.AltaHospitalariaRegistrar) (domain.AltaHospitalariaResult, error)

	RegistrarPlanTrabajoItem(ctx context.Context, idAtencion, idPrimeraAtencion int, item domain.PlanTrabajoItemRegistrarItem) (domain.PlanTrabajoItemAgregarResult, error)
	ListarPlanTrabajoItems(ctx context.Context, idAtencion, idPrimeraAtencion int, tipo string) ([]domain.PlanTrabajoItem, error)
	EliminarPlanTrabajoItem(ctx context.Context, idAtencion, idPrimeraAtencion, idPlanTrabajoItem int) error

	RegistrarMedicamento(ctx context.Context, in domain.RecetaMedicamentoRegistrar) (domain.RecetaMedicamentoRegistrarResult, error)
	ActualizarMedicamentoCampos(ctx context.Context, in domain.RecetaMedicamentoCamposActualizar) error
	ListarMedicamentos(ctx context.Context, idCuentaAtencion, idPrimeraAtencion, idPuntoCarga int) ([]domain.RecetaMedicamentoItem, error)
	EliminarMedicamento(ctx context.Context, idAtencion, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion int) error

	RegistrarExamen(ctx context.Context, in domain.RecetaExamenRegistrar) (domain.RecetaExamenRegistrarResult, error)
	ListarExamenes(ctx context.Context, idCuentaAtencion, idPrimeraAtencion int) ([]domain.RecetaMedicamentoItem, error)
	EliminarExamen(ctx context.Context, idAtencion, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion int) error

	RegistrarProcedimiento(ctx context.Context, in domain.RecetaProcedimientoRegistrar) (domain.RecetaExamenRegistrarResult, error)
	ListarProcedimientos(ctx context.Context, idCuentaAtencion, idPrimeraAtencion int) ([]domain.RecetaMedicamentoItem, error)
	EliminarProcedimiento(ctx context.Context, idAtencion, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion int) error

	RegistrarSolicitudInterconsulta(ctx context.Context, in domain.InterconsultaRegistrarSolicitud) (idInterconsulta int, err error)
	ListarInterconsultas(ctx context.Context, idAtencion, idEpisodio, idPrimeraAtencion int) ([]domain.InterconsultaListadoItem, error)
	EliminarSolicitudInterconsulta(ctx context.Context, idAtencion, idInterconsulta, idEpisodio int) error
}
