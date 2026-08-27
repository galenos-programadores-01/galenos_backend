package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// CatalogRepository expone los catálogos de datos maestros (etnias,
// idiomas, tipos y ubicación geográfica) usados por los formularios del
// frontend. Cada consulta invoca el procedimiento almacenado homólogo del
// proyecto FastAPI.
type CatalogRepository interface {
	ListEtnias(ctx context.Context) ([]domain.Etnia, error)
	ListIdiomas(ctx context.Context) ([]domain.Idioma, error)
	ListTipoSexos(ctx context.Context) ([]domain.TipoSexo, error)
	ListEstadosCivil(ctx context.Context) ([]domain.TipoEstadoCivil, error)
	ListGradosInstruccion(ctx context.Context) ([]domain.TipoGradoInstruccion, error)
	ListOcupaciones(ctx context.Context) ([]domain.TipoOcupacion, error)
	ListTiposDocumentos(ctx context.Context) ([]domain.TipoDocumento, error)
	ListDepartamentos(ctx context.Context) ([]domain.Departamento, error)
	ListProvincias(ctx context.Context, idDepartamento int64) ([]domain.Provincia, error)
	ListDistritos(ctx context.Context, idProvincia int64) ([]domain.Distrito, error)
	ListCentrosPoblados(ctx context.Context, idDistrito int64) ([]domain.CentroPoblado, error)
	ListPaises(ctx context.Context) ([]domain.Pais, error)
	ListEstadosLlegoPaciente(ctx context.Context) ([]domain.EstadoLlegoPaciente, error)
	ListFuentesFinanciamiento(ctx context.Context) ([]domain.FuenteFinanciamiento, error)
	ListServicios(ctx context.Context, idTipoServicio int64) ([]domain.Servicio, error)
	ListEspecialidades(ctx context.Context) ([]domain.Especialidad, error)
	ListarEspecialidadesPorDepartamento(ctx context.Context, idDepartamento int) ([]domain.EspecialidadSimple, error)
	ListarEspecialidadesQx(ctx context.Context) ([]domain.EspecialidadSimple, error)
	GetDatosInstitucion(ctx context.Context) (*domain.DatosInstitucion, error)
	GetParametro(ctx context.Context, idParametro int64) (*domain.Parametro, error)
	ListRecetaFrecuencias(ctx context.Context) ([]domain.CatalogItem, error)
	ListRecetaUnidadesDosis(ctx context.Context) ([]domain.CatalogItem, error)
	ListRecetaViasAdministracion(ctx context.Context) ([]domain.CatalogItem, error)
	BuscarMedicamentosReceta(ctx context.Context, filtro string, idPaciente int) ([]domain.MedicamentoBusqueda, error)
}
