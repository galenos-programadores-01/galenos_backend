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
// (SP usp_go_ListarEspecialidades).
type Especialidad struct {
	ID     int64
	Nombre *string
}

// EspecialidadSimple contiene el id y nombre de una especialidad
// (SP usp_go_ListarEspecialidadXDepartamento).
type EspecialidadSimple struct {
	IdEspecialidad int    `json:"idEspecialidad"`
	Nombre         string `json:"nombre"`
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

type CatalogItem struct {
	ID          int64  `json:"id"`
	Descripcion string `json:"descripcion"`
}

type MedicamentoBusqueda struct {
	IdProducto          int     `json:"idProducto"`
	Codigo              string  `json:"codigo"`
	Nombre              string  `json:"nombre"`
	Stock               int     `json:"stock"`
	Precio              float64 `json:"precio"`
	IdDosisRecetada     int     `json:"idDosisRecetada"`
	IdUNIDDosisReceta   int     `json:"idUNIDDosisReceta"`
	IdFrecuencia        int     `json:"idFrecuencia"`
	IdViaAdministracion int     `json:"idViaAdministracion"`
	TieneRecetaAnterior int     `json:"tieneRecetaAnterior"`
}
