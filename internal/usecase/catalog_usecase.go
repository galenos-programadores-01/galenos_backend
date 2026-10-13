package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/adapters/output/cache"
	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type catalogUseCase struct {
	repo  output.CatalogRepository
	cache cache.Cache
}

// NewCatalogUseCase construye el caso de uso de catÃ¡logos de datos maestros.
func NewCatalogUseCase(repo output.CatalogRepository, cache cache.Cache) input.CatalogService {
	return &catalogUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *catalogUseCase) ListEtnias(ctx context.Context) ([]domain.Etnia, error) {
	items, err := uc.repo.ListEtnias(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing etnias: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListIdiomas(ctx context.Context) ([]domain.Idioma, error) {
	items, err := uc.repo.ListIdiomas(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing idiomas: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListTipoSexos(ctx context.Context) ([]domain.TipoSexo, error) {
	items, err := uc.repo.ListTipoSexos(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing tipos sexo: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListEstadosCivil(ctx context.Context) ([]domain.TipoEstadoCivil, error) {
	items, err := uc.repo.ListEstadosCivil(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing estados civil: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListGradosInstruccion(ctx context.Context) ([]domain.TipoGradoInstruccion, error) {
	items, err := uc.repo.ListGradosInstruccion(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing grados instruccion: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListOcupaciones(ctx context.Context) ([]domain.TipoOcupacion, error) {
	items, err := uc.repo.ListOcupaciones(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing ocupaciones: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListTiposDocumentos(ctx context.Context) ([]domain.TipoDocumento, error) {
	items, err := uc.repo.ListTiposDocumentos(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing tipos documentos: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListDepartamentos(ctx context.Context) ([]domain.Departamento, error) {
	items, err := uc.repo.ListDepartamentos(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing departamentos: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListProvincias(ctx context.Context, idDepartamento int64) ([]domain.Provincia, error) {
	items, err := uc.repo.ListProvincias(ctx, idDepartamento)
	if err != nil {
		return nil, fmt.Errorf("listing provincias: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListDistritos(ctx context.Context, idProvincia int64) ([]domain.Distrito, error) {
	items, err := uc.repo.ListDistritos(ctx, idProvincia)
	if err != nil {
		return nil, fmt.Errorf("listing distritos: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListCentrosPoblados(ctx context.Context, idDistrito int64) ([]domain.CentroPoblado, error) {
	items, err := uc.repo.ListCentrosPoblados(ctx, idDistrito)
	if err != nil {
		return nil, fmt.Errorf("listing centros poblados: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListPaises(ctx context.Context) ([]domain.Pais, error) {
	items, err := uc.repo.ListPaises(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing paises: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListEstadosLlegoPaciente(ctx context.Context) ([]domain.EstadoLlegoPaciente, error) {
	items, err := uc.repo.ListEstadosLlegoPaciente(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing estados llego paciente: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListFuentesFinanciamiento(ctx context.Context) ([]domain.FuenteFinanciamiento, error) {
	items, err := uc.repo.ListFuentesFinanciamiento(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing fuentes financiamiento: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListServicios(ctx context.Context, idTipoServicio int64) ([]domain.Servicio, error) {
	items, err := uc.repo.ListServicios(ctx, idTipoServicio)
	if err != nil {
		return nil, fmt.Errorf("listing servicios: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) GetDatosInstitucion(ctx context.Context) (*domain.DatosInstitucion, error) {
	item, err := uc.repo.GetDatosInstitucion(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting datos institucion: %w", err)
	}
	return item, nil
}

func (uc *catalogUseCase) ListEspecialidades(ctx context.Context) ([]domain.Especialidad, error) {
	items, err := uc.repo.ListEspecialidades(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing especialidades: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListarEspecialidadesPorDepartamento(ctx context.Context, idDepartamento int) ([]domain.EspecialidadSimple, error) {
	items, err := uc.repo.ListarEspecialidadesPorDepartamento(ctx, idDepartamento)
	if err != nil {
		return nil, fmt.Errorf("listing especialidades por departamento: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListarEspecialidadesQx(ctx context.Context) ([]domain.EspecialidadSimple, error) {
	items, err := uc.repo.ListarEspecialidadesQx(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing especialidades qx: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListEspecialidadesPorTipoServicio(ctx context.Context, idTipoServicio int64) ([]domain.Especialidad, error) {
	items, err := uc.repo.ListEspecialidadesPorTipoServicio(ctx, idTipoServicio)
	if err != nil {
		return nil, fmt.Errorf("listing especialidades por tipo servicio: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListEspecialidadesListar(ctx context.Context) ([]domain.Especialidad, error) {
	const cacheKey = "cache:especialidades:listar"
	if uc.cache != nil {
		var cached []domain.Especialidad
		if uc.cache.Get(ctx, cacheKey, &cached) {
			return cached, nil
		}
	}
	items, err := uc.repo.ListEspecialidadesListar(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing especialidades listar: %w", err)
	}
	if uc.cache != nil && len(items) > 0 {
		uc.cache.Set(ctx, cacheKey, items, 24*time.Hour)
	}
	return items, nil
}

func (uc *catalogUseCase) ListarServiciosPorPrioridad(ctx context.Context, idPrioridad int, fechaNac string) ([]domain.ServicioSimple, error) {
	items, err := uc.repo.ListarServiciosPorPrioridad(ctx, idPrioridad, fechaNac)
	if err != nil {
		return nil, fmt.Errorf("listing servicios por prioridad: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) GetParametro(ctx context.Context, idParametro int64) (*domain.Parametro, error) {
	item, err := uc.repo.GetParametro(ctx, idParametro)
	if err != nil {
		return nil, fmt.Errorf("getting parametro: %w", err)
	}
	return item, nil
}

func (uc *catalogUseCase) ListParametrosClinicos(ctx context.Context, idGrupo int64) ([]domain.ParametroClinico, error) {
	cacheKey := fmt.Sprintf("cache:param_clinicos:grupo:%d", idGrupo)
	if uc.cache != nil {
		var cached []domain.ParametroClinico
		if uc.cache.Get(ctx, cacheKey, &cached) {
			return cached, nil
		}
	}
	items, err := uc.repo.ListParametrosClinicos(ctx, idGrupo)
	if err != nil {
		return nil, fmt.Errorf("listing parametros clinicos: %w", err)
	}
	if uc.cache != nil && len(items) > 0 {
		uc.cache.Set(ctx, cacheKey, items, 24*time.Hour)
	}
	return items, nil
}

func (uc *catalogUseCase) ListPrioridadesPrimeraAtencion(ctx context.Context) ([]domain.PrioridadPrimeraAtencion, error) {
	const cacheKey = "cache:prioridades:primera-atencion"
	if uc.cache != nil {
		var cached []domain.PrioridadPrimeraAtencion
		if uc.cache.Get(ctx, cacheKey, &cached) {
			return cached, nil
		}
	}
	items, err := uc.repo.ListPrioridadesPrimeraAtencion(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing prioridades primera atencion: %w", err)
	}
	if uc.cache != nil && len(items) > 0 {
		uc.cache.Set(ctx, cacheKey, items, 24*time.Hour)
	}
	return items, nil
}

func (uc *catalogUseCase) ListSubclasificacionesDiagnostico(ctx context.Context, idTipoServicio, idClasificacionDx int64) ([]domain.SubclasificacionDiagnostico, error) {
	if idTipoServicio <= 0 || idClasificacionDx <= 0 {
		return nil, fmt.Errorf("idTipoServicio e idClasificacionDx son obligatorios")
	}
	cacheKey := fmt.Sprintf("cache:subclasif-dx:ts:%d:clasif:%d", idTipoServicio, idClasificacionDx)
	if uc.cache != nil {
		var cached []domain.SubclasificacionDiagnostico
		if uc.cache.Get(ctx, cacheKey, &cached) {
			return cached, nil
		}
	}
	items, err := uc.repo.ListSubclasificacionesDiagnostico(ctx, idTipoServicio, idClasificacionDx)
	if err != nil {
		return nil, fmt.Errorf("listing subclasificaciones diagnostico: %w", err)
	}
	if uc.cache != nil && len(items) > 0 {
		uc.cache.Set(ctx, cacheKey, items, 24*time.Hour)
	}
	return items, nil
}

func (uc *catalogUseCase) ListDiagnosticosCatalogo(ctx context.Context, filtro string, idAtencion, idPaciente int64) ([]domain.DiagnosticoCatalogo, error) {
	items, err := uc.repo.ListDiagnosticosCatalogo(ctx, filtro, idAtencion, idPaciente)
	if err != nil {
		return nil, fmt.Errorf("listing diagnosticos catalogo: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListMedicamentosCatalogo(ctx context.Context, filtro string, idPaciente int64) ([]domain.MedicamentoCatalogo, error) {
	items, err := uc.repo.ListMedicamentosCatalogo(ctx, filtro, idPaciente)
	if err != nil {
		return nil, fmt.Errorf("listing medicamentos catalogo: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListExamenesCatalogo(ctx context.Context, idTipoFinanciamiento int64, filtro string) ([]domain.ExamenCatalogo, error) {
	if idTipoFinanciamiento <= 0 {
		return nil, fmt.Errorf("idTipoFinanciamiento inválido")
	}
	cleanFilter := strings.ToLower(strings.TrimSpace(filtro))
	cacheKey := fmt.Sprintf("cache:examenes:fin:%d:%s", idTipoFinanciamiento, cleanFilter)
	if uc.cache != nil {
		var cached []domain.ExamenCatalogo
		if uc.cache.Get(ctx, cacheKey, &cached) {
			return cached, nil
		}
	}
	items, err := uc.repo.ListExamenesCatalogo(ctx, idTipoFinanciamiento, filtro)
	if err != nil {
		return nil, fmt.Errorf("listing examenes catalogo: %w", err)
	}
	if uc.cache != nil && len(items) > 0 {
		uc.cache.Set(ctx, cacheKey, items, 60*time.Minute)
	}
	return items, nil
}

func (uc *catalogUseCase) ListProcedimientosCatalogo(ctx context.Context, idPuntoCarga, idTipoFinanciamiento int64, filtro string) ([]domain.ProcedimientoCatalogo, error) {
	items, err := uc.repo.ListProcedimientosCatalogo(ctx, idPuntoCarga, idTipoFinanciamiento, filtro)
	if err != nil {
		return nil, fmt.Errorf("listing procedimientos catalogo: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListRecetaDosis(ctx context.Context) ([]domain.RecetaDosis, error) {
	items, err := uc.repo.ListRecetaDosis(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing receta dosis: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListRecetaUnidadesDosis(ctx context.Context) ([]domain.RecetaUnidadDosis, error) {
	items, err := uc.repo.ListRecetaUnidadesDosis(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing receta unidades dosis: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListRecetaFrecuencias(ctx context.Context) ([]domain.RecetaFrecuencia, error) {
	items, err := uc.repo.ListRecetaFrecuencias(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing receta frecuencias: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) ListViasAdministracion(ctx context.Context) ([]domain.ViaAdministracion, error) {
	items, err := uc.repo.ListViasAdministracion(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing vias administracion: %w", err)
	}
	return items, nil
}

func (uc *catalogUseCase) BuscarMedicamentosReceta(ctx context.Context, filtro string, idPaciente int) ([]domain.MedicamentoBusqueda, error) {
	cleanFilter := strings.ToLower(strings.TrimSpace(filtro))
	cacheKey := fmt.Sprintf("cache:meds:%s:%d", cleanFilter, idPaciente)
	if uc.cache != nil {
		var cached []domain.MedicamentoBusqueda
		if uc.cache.Get(ctx, cacheKey, &cached) {
			return cached, nil
		}
	}
	items, err := uc.repo.BuscarMedicamentosReceta(ctx, filtro, idPaciente)
	if err != nil {
		return nil, fmt.Errorf("buscar medicamentos receta: %w", err)
	}
	if uc.cache != nil && len(items) > 0 {
		uc.cache.Set(ctx, cacheKey, items, 30*time.Minute)
	}
	return items, nil
}

// ActualizarParametroValorTexto marca el parámetro con ValorTexto = 'N' (SP
// usp_go_ActualizarParametroValorTextoPorId) e invalida la cache de parámetros
// clínicos. La cache está indexada por grupo, y el SP no devuelve a qué grupo
// pertenece el parámetro, así que se purgan los cinco grupos conocidos.
func (uc *catalogUseCase) ActualizarParametroValorTexto(ctx context.Context, idParametro int) error {
	if err := uc.repo.ActualizarParametroValorTexto(ctx, idParametro); err != nil {
		return err
	}

	if uc.cache != nil {
		keys := make([]string, 0, 5)
		for idGrupo := 1; idGrupo <= 5; idGrupo++ {
			keys = append(keys, fmt.Sprintf("cache:param_clinicos:grupo:%d", idGrupo))
		}
		uc.cache.Delete(ctx, keys...)
	}

	return nil
}

func (uc *catalogUseCase) ListEstablecimientos(ctx context.Context, filtro string, tipo string) ([]domain.EstablecimientoBusqueda, error) {
	items, err := uc.repo.ListEstablecimientos(ctx, filtro, tipo)
	if err != nil {
		return nil, fmt.Errorf("listing establecimientos: %w", err)
	}
	return items, nil
}
