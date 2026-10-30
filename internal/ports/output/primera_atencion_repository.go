package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// PrimeraAtencionRepository es el puerto de salida hacia la persistencia
// de Primera Atención.
type PrimeraAtencionRepository interface {
	ListarBandeja(ctx context.Context, params shared.PrimeraAtencionParams) ([]map[string]any, error)
	// IdsConBorradorActivo indica qué IdAtencion tienen fila en Tab_PrimeraAtencion (borrador).
	IdsConBorradorActivo(ctx context.Context, idsAtencion []int) (map[int]bool, error)
	// IdsFirmados indica qué IdAtencion ya tienen EstadoFirma = 1.
	IdsFirmados(ctx context.Context, idsAtencion []int) (map[int]bool, error)
	GuardarUpsert(ctx context.Context, in domain.PrimeraAtencionUpsert) (domain.PrimeraAtencionUpsertResult, error)
	ObtenerPorIdAtencion(ctx context.Context, idAtencion int) (map[string]any, error)
	ListarDocumentosFirmados(ctx context.Context, idPrimeraAtencion int) ([]domain.DocumentoFirmadoPA, error)
	RegistrarDocumentosFirmados(ctx context.Context, idPrimeraAtencion, nroEdicion int, docs []domain.DocumentoFirmadoPA) error
	ActualizarEdicionVigente(ctx context.Context, idPrimeraAtencion, nroEdicion int, ubicacion string) error
	EliminarDetalleRectificacion(ctx context.Context, idReceta, idProducto int) error
	CabeceraRecetaAnulada(ctx context.Context, idReceta int) (bool, error)
	ValorConfiguracionEntero(ctx context.Context, idListItem int, clave string) (int, error)
	AgregarDiagnostico(ctx context.Context, in domain.DiagnosticoAgregar) (domain.DiagnosticoAgregarResult, error)
	ListarDiagnosticos(ctx context.Context, idAtencion, idPrimeraAtencion int) ([]domain.DiagnosticoAtencion, error)
	EliminarDiagnostico(ctx context.Context, idAtencion, idPrimeraAtencion, idAtencionDiagnostico int) error
	RegistrarAltaHospitalaria(ctx context.Context, in domain.AltaHospitalariaRegistrar) (domain.AltaHospitalariaResult, error)

	AgregarPlanTrabajoItem(ctx context.Context, in domain.PlanTrabajoItemAgregar) (domain.PlanTrabajoItemAgregarResult, error)
	ListarPlanTrabajoItems(ctx context.Context, idAtencion, idPrimeraAtencion int, tipo string) ([]domain.PlanTrabajoItem, error)
	EliminarPlanTrabajoItem(ctx context.Context, idAtencion, idPrimeraAtencion, idPlanTrabajoItem int) error

	// AgregarRecetaCabecera → usp_go_PA_RecetaCabeceraAgregar (aislado por IdPrimeraAtencion).
	AgregarRecetaCabecera(ctx context.Context, in domain.RecetaCabeceraAgregar) (int, error)
	AgregarRecetaDetalle(ctx context.Context, in domain.RecetaDetalleAgregar) error
	ActualizarRecetaDetalleCamposClinicos(ctx context.Context, in domain.RecetaMedicamentoCamposActualizar) error
	ListarRecetaDetalle(ctx context.Context, filtro domain.RecetaDetalleListarFiltro) ([]domain.RecetaMedicamentoItem, error)
	EliminarRecetaDetalle(ctx context.Context, idReceta, idProducto, idCuentaAtencion, idPrimeraAtencion int) error

	InsertarInterconsulta(ctx context.Context, in domain.InterconsultaInsertar) (int, error)
	InsertarInterconsultaDiagnostico(ctx context.Context, in domain.InterconsultaDiagnosticoInsertar) error
	ListarInterconsultasPorAtencion(ctx context.Context, idAtencion, idEpisodio, idPrimeraAtencion int) ([]domain.InterconsultaListadoItem, error)
	EliminarInterconsulta(ctx context.Context, idInterconsulta, idEpisodio int) error
}
