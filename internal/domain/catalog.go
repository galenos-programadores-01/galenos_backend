package domain

// Etnia es un registro del catálogo de etnias (SP ups_go_ListarEtnias).
type Etnia struct {
	Codigo      string
	Descripcion *string
}

// Idioma es un registro del catálogo de lenguas (SP ups_go_ListarIdiomas).
type Idioma struct {
	ID     int64
	Lengua *string
}

// TipoSexo es un registro del catálogo de sexos (SP usp_go_ListarTiposSexos).
type TipoSexo struct {
	ID          int64
	Descripcion *string
}

// TipoEstadoCivil es un registro del catálogo de estados civiles
// (SP usp_go_ListarEstadosCivil).
type TipoEstadoCivil struct {
	ID          int64
	Descripcion *string
}

// TipoGradoInstruccion es un registro del catálogo de grados de instrucción
// (SP usp_go_ListarGradoInstruccion).
type TipoGradoInstruccion struct {
	ID          int64
	Descripcion *string
}

// TipoOcupacion es un registro del catálogo de ocupaciones
// (SP usp_go_ListarOcupaciones).
type TipoOcupacion struct {
	ID          int64
	Descripcion *string
}

// TipoDocumento es un registro del catálogo de tipos de documento de
// identidad (SP usp_go_ListarTiposDocumentos).
type TipoDocumento struct {
	ID          int64
	Descripcion *string
}

// Departamento es un registro de la tabla Departamentos
// (SP usp_go_ListarDepartamentos).
type Departamento struct {
	ID     int64
	Nombre *string
}

// Provincia es un registro de la tabla Provincias
// (SP usp_go_ListarProvincias @IdDepartamento).
type Provincia struct {
	ID     int64
	Nombre *string
}

// Distrito es un registro de la tabla Distritos
// (SP usp_go_ListarDistritos @IdProvincia).
type Distrito struct {
	ID     int64
	Nombre *string
}

// CentroPoblado es un registro del catálogo de centros poblados
// (SP usp_go_ListarCentrosPoblados @IdDistrito).
type CentroPoblado struct {
	ID     int64
	Nombre *string
}

// Pais es un registro del catálogo de países (SP usp_go_ListarPaises).
type Pais struct {
	ID     int64
	Nombre *string
}

// EstadoLlegoPaciente es un registro del catálogo de estados de llegada del
// paciente (SP usp_go_listarEstadosLlegoPaciente).
type EstadoLlegoPaciente struct {
	ID          int64
	Descripcion *string
}

// FuenteFinanciamiento es un registro del catálogo de fuentes de
// financiamiento (SP usp_go_ListarFuentesFinanciamiento).
type FuenteFinanciamiento struct {
	ID                   int64
	Descripcion          *string
	IdTipoFinanciamiento int64
}

// Servicio es un registro del catálogo de servicios por tipo
// (SP usp_go_ListarServicios @IdTipoServicio).
type Servicio struct {
	ID     int64
	Nombre *string
}

// DatosInstitucion contiene los datos de la institución (EE.SS.) que el
// SP webParametrosDatosInstitucion devuelve en una única fila.
type DatosInstitucion struct {
	RucEESS    *string
	Nombre     *string
	Direccion  *string
	Telefono   *string
	Codigo     *string
	CodRenaes  *string
	LogoHospi  *string
	LogoMinsa  *string
	UbigeoHosp *string
}

// Especialidad es un registro del catálogo de especialidades
// (SP usp_go_ListarEspecialidades) o filtrado por tipo de servicio
// (SP usp_go_EspecialidadesXidTipoServicio @IdTipoServicio).
type Especialidad struct {
	ID     int64
	Nombre *string
}

// EspecialidadSimple contiene el id y nombre de una especialidad
// (SP usp_go_ListarEspecialidadXDepartamento / usp_go_ListarEspecialidadesQx).
type EspecialidadSimple struct {
	IdEspecialidad int    `json:"idEspecialidad"`
	Nombre         string `json:"nombre"`
}

// ServicioSimple contiene el id y nombre de un servicio
// (SP usp_go_ListarServiciosXPrioridad).
type ServicioSimple struct {
	IdServicio int    `json:"idServicio"`
	Nombre     string `json:"nombre"`
}

// Parametro contiene los valores de un parámetro configurable
// (SP usp_go_webParametroSeleccionarPorId @IdParametro).
type Parametro struct {
	Tipo       *string
	Codigo     *string
	ValorTexto *string
	ValorInt   *int64
	ValorFloat *float64
}

// CatalogItem proyección genérica id+descripción (uso interno / mapper develop).
// El JSON HTTP se define en adapters/http, no aquí.
type CatalogItem struct {
	ID          int64
	Descripcion string
}

// MedicamentoBusqueda proyección enriquecida de receta (develop).
// Mismo SP que MedicamentoCatalogo; más columnas. JSON vía DTO HTTP.
type MedicamentoBusqueda struct {
	IdProducto          int
	Codigo              string
	Nombre              string
	Stock               int
	Precio              float64
	IdDosisRecetada     int
	IdUNIDDosisReceta   int
	IdFrecuencia        int
	IdViaAdministracion int
	TieneRecetaAnterior int
}

// ParametroClinico — usp_go_Cat_ParametroClinico_Listar @IdGrupo.
// ID = columna SQL IdParametro. IdGrupo: 1 estado, 2 hidratación, 3 conciencia, 4 egreso.
type ParametroClinico struct {
	ID          int64
	Descripcion *string
	Codigo      *string
}

// PrioridadPrimeraAtencion — usp_go_Cat_Prioridad_PrimeraAtencion.
// ID = IdTipoGravedad del SP (misma clave que bandeja / IdTipoPrioridad del paciente).
type PrioridadPrimeraAtencion struct {
	ID          int64
	Descripcion *string
}

// SubclasificacionDiagnostico — usp_go_SubclasificacionDiagnosticos_Select
// (@IdTipoServicio, @IdClasificacionDx).
// Clasificación 9 = PA general; 7 = interconsulta.
type SubclasificacionDiagnostico struct {
	ID               int64
	DescripcionCorta *string
	DescripcionLarga *string
}

// DiagnosticoCatalogo — usp_go_SelectDiagnosticos (@filtro, @IdAtencion, @IdPaciente).
type DiagnosticoCatalogo struct {
	ID          int64
	Descripcion string
	CodigoCIE10 string
}

// MedicamentoCatalogo — usp_go_SelectMedicamentosFiltro (proyección slim PA).
type MedicamentoCatalogo struct {
	IDProducto     int64
	Codigo         string
	Nombre         string
	PrecioUnitario float64 // columna Precio del SP → RecetaDetalle.Precio
}

// ExamenCatalogo — webProcedimentosRecetasFiltarNombrePuntoCarga
// (@IdTipoFinanciamiento, @Filtro). IdProducto = IdProductoSolo del SP
// (el IdProducto compuesto "54707|2" no se usa como PK).
type ExamenCatalogo struct {
	IDProducto     int64
	Codigo         string
	Nombre         string
	NombreLargo    string
	IDPuntoCarga   int64
	PuntoCarga     string
	Servicio       string
	PrecioUnitario float64
}

// ProcedimientoCatalogo — usp_go_ProcedimentosFiltarNombrePuntoCarga.
type ProcedimientoCatalogo struct {
	IDProducto     int64
	Codigo         string
	Nombre         string
	NombreLargo    string
	IDPuntoCarga   int64
	PuntoCarga     string
	PrecioUnitario float64
	EsCPT          bool
}

// RecetaDosis — usp_go_RecetaDosisSelecionarTodos.
type RecetaDosis struct {
	ID          int64
	NumeroDosis string
}

// RecetaUnidadDosis — Usp_go_SelectRecetaUndDosisSelecionarTodos.
type RecetaUnidadDosis struct {
	ID              int64
	NumeroUNIDDosis string
}

// RecetaFrecuencia — usp_go_SelectRecetaFrecuenciaSelecionarTodos.
type RecetaFrecuencia struct {
	ID          int64
	Descripcion string
}

// ViaAdministracion — usp_go_RecetasListadoViasAdministracion.
type ViaAdministracion struct {
	ID          int64
	Descripcion string
}

type EstablecimientoBusqueda struct {
	IdEstablecimiento int    `json:"idEstablecimiento"`
	Codigo            string `json:"codigo"`
	Nombre            string `json:"nombre"`
	Distrito          string `json:"distrito"`
	Provincia         string `json:"provincia"`
	Departamento      string `json:"departamento"`
	NombreLargo       string `json:"nombreLargo"`
}
