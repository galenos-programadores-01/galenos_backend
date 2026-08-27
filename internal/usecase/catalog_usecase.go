package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type catalogUseCase struct {
	repo output.CatalogRepository
}

// NewCatalogUseCase construye el caso de uso de catálogos de datos maestros.
func NewCatalogUseCase(repo output.CatalogRepository) input.CatalogService {
	return &catalogUseCase{repo: repo}
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

func (uc *catalogUseCase) GetParametro(ctx context.Context, idParametro int64) (*domain.Parametro, error) {
	item, err := uc.repo.GetParametro(ctx, idParametro)
	if err != nil {
		return nil, fmt.Errorf("getting parametro: %w", err)
	}
	return item, nil
}

func (uc *catalogUseCase) ListRecetaFrecuencias(ctx context.Context) ([]domain.CatalogItem, error) {
	return uc.repo.ListRecetaFrecuencias(ctx)
}

func (uc *catalogUseCase) ListRecetaUnidadesDosis(ctx context.Context) ([]domain.CatalogItem, error) {
	return uc.repo.ListRecetaUnidadesDosis(ctx)
}

func (uc *catalogUseCase) ListRecetaViasAdministracion(ctx context.Context) ([]domain.CatalogItem, error) {
	return uc.repo.ListRecetaViasAdministracion(ctx)
}

func (uc *catalogUseCase) BuscarMedicamentosReceta(ctx context.Context, filtro string, idPaciente int) ([]domain.MedicamentoBusqueda, error) {
	return uc.repo.BuscarMedicamentosReceta(ctx, filtro, idPaciente)
}
