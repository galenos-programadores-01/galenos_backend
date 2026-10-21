package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// CatalogService es el puerto de entrada para consultar los catálogos de
// datos maestros que alimentan los formularios del frontend.
type CatalogService interface {
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
	// Develop / QX
	ListarEspecialidadesPorDepartamento(ctx context.Context, idDepartamento int) ([]domain.EspecialidadSimple, error)
	ListarEspecialidadesQx(ctx context.Context) ([]domain.EspecialidadSimple, error)
	// PA interconsulta (legado por tipo servicio)
	ListEspecialidadesPorTipoServicio(ctx context.Context, idTipoServicio int64) ([]domain.Especialidad, error)
	// PA interconsulta: usp_go_Especialidades_Listar
	ListEspecialidadesListar(ctx context.Context) ([]domain.Especialidad, error)
	ListarServiciosPorPrioridad(ctx context.Context, idPrioridad int, fechaNac string) ([]domain.ServicioSimple, error)
	GetDatosInstitucion(ctx context.Context) (*domain.DatosInstitucion, error)
	GetParametro(ctx context.Context, idParametro int64) (*domain.Parametro, error)
	// Compartido PA + develop (mismo SP; JSON tipado en HTTP)
	ListParametrosClinicos(ctx context.Context, idGrupo int64) ([]domain.ParametroClinico, error)
	// PA bandeja + interconsulta
	ListPrioridadesPrimeraAtencion(ctx context.Context) ([]domain.PrioridadPrimeraAtencion, error)
	ListSubclasificacionesDiagnostico(ctx context.Context, idTipoServicio, idClasificacionDx int64) ([]domain.SubclasificacionDiagnostico, error)
	// PA
	ListDiagnosticosCatalogo(ctx context.Context, filtro string, idAtencion, idPaciente int64) ([]domain.DiagnosticoCatalogo, error)
	ListMedicamentosCatalogo(ctx context.Context, filtro string, idPaciente int64) ([]domain.MedicamentoCatalogo, error)
	ListExamenesCatalogo(ctx context.Context, idTipoFinanciamiento int64, filtro string) ([]domain.ExamenCatalogo, error)
	ListProcedimientosCatalogo(ctx context.Context, idPuntoCarga, idTipoFinanciamiento int64, filtro string) ([]domain.ProcedimientoCatalogo, error)
	ListRecetaDosis(ctx context.Context) ([]domain.RecetaDosis, error)
	ListRecetaUnidadesDosis(ctx context.Context) ([]domain.RecetaUnidadDosis, error)
	ListRecetaFrecuencias(ctx context.Context) ([]domain.RecetaFrecuencia, error)
	ListViasAdministracion(ctx context.Context) ([]domain.ViaAdministracion, error)
	// Develop: misma SP medicamentos, proyección enriquecida
	BuscarMedicamentosReceta(ctx context.Context, filtro string, idPaciente int) ([]domain.MedicamentoBusqueda, error)
	ActualizarParametroValorTexto(ctx context.Context, idParametro int) error
	ListEstablecimientos(ctx context.Context, filtro string, tipo string) ([]domain.EstablecimientoBusqueda, error)
	ListSubclasificacionesDiagnosticosItems(ctx context.Context, idTipoServicio, idClasificacionDx int) ([]domain.SubclasificacionDxItem, error)
}
