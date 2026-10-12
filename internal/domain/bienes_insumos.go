package domain

type BienInsumoItem struct {
	IdProducto           int     `json:"idProducto"`
	Codigo               string  `json:"codigo"`
	CodigoSiga           string  `json:"codigoSiga"`
	CodigoSunat          string  `json:"codigoSunat"`
	Nombre               string  `json:"nombre"`
	NombreComercial      string  `json:"nombreComercial"`
	IdGrupoFarmacologico int     `json:"idGrupoFarmacologico"`
	GrupoFarmacologico   string  `json:"grupoFarmacologico"`
	PrecioCompra         float64 `json:"precioCompra"`
	PrecioDistribucion   float64 `json:"precioDistribucion"`
	PrecioUltimaCompra   float64 `json:"precioUltimaCompra"`
	StockMinimo          int     `json:"stockMinimo"`
	TipoProducto         int     `json:"tipoProducto"`
	Concentracion        string  `json:"concentracion"`
	Presentacion         string  `json:"presentacion"`
	FormaFarmaceutica    string  `json:"formaFarmaceutica"`
	Petitorio            bool    `json:"petitorio"`
	Estado               int     `json:"estado"`
	EsVital              bool    `json:"esVital"`
	EsEstrategico        bool    `json:"esEstrategico"`
	EsOncologico         bool    `json:"esOncologico"`
	CargaEnFua           bool    `json:"cargaEnFua"`
	EsIndicadorGestion   bool    `json:"esIndicadorGestion"`
}

type BienInsumoPrecioPlan struct {
	IdPlanCatalogo       int     `json:"idPlanCatalogo"`
	IdTipoFinanciamiento int     `json:"idTipoFinanciamiento"`
	Plan                 string  `json:"plan"`
	PrecioUnitario       float64 `json:"precioUnitario"`
	Activo               bool    `json:"activo"`
}

type BienInsumoDetalle struct {
	IdProducto              int                    `json:"idProducto"`
	Codigo                  string                 `json:"codigo"`
	CodigoSiga              string                 `json:"codigoSiga"`
	CodigoSunat             string                 `json:"codigoSunat"`
	Nombre                  string                 `json:"nombre"`
	NombreComercial         string                 `json:"nombreComercial"`
	IdGrupoFarmacologico    int                    `json:"idGrupoFarmacologico"`
	GrupoFarmacologico      string                 `json:"grupoFarmacologico"`
	IdSubGrupoFarmacologico int                    `json:"idSubGrupoFarmacologico"`
	SubGrupoFarmacologico   string                 `json:"subGrupoFarmacologico"`
	IdPartida               int                    `json:"idPartida"`
	Partida                 string                 `json:"partida"`
	IdCentroCosto           int                    `json:"idCentroCosto"`
	CentroCosto             string                 `json:"centroCosto"`
	PrecioCompra            float64                `json:"precioCompra"`
	PrecioDistribucion      float64                `json:"precioDistribucion"`
	PrecioDonacion          float64                `json:"precioDonacion"`
	PrecioUltimaCompra      float64                `json:"precioUltimaCompra"`
	StockMinimo             int                    `json:"stockMinimo"`
	TipoProducto            int                    `json:"tipoProducto"`
	Denominacion            string                 `json:"denominacion"`
	Concentracion           string                 `json:"concentracion"`
	Presentacion            string                 `json:"presentacion"`
	FormaFarmaceutica       string                 `json:"formaFarmaceutica"`
	MaterialEnvase          string                 `json:"materialEnvase"`
	PresentacionEnvase      string                 `json:"presentacionEnvase"`
	Fabricante              string                 `json:"fabricante"`
	IdPaisOrigen            int                    `json:"idPaisOrigen"`
	PaisOrigen              string                 `json:"paisOrigen"`
	Petitorio               bool                   `json:"petitorio"`
	TipoProductoSismed      string                 `json:"tipoProductoSismed"`
	Estado                  int                    `json:"estado"`
	EstadoDigemid           int                    `json:"estadoDigemid"`
	IdTipoPrefijo           int                    `json:"idTipoPrefijo"`
	TipoPrefijo             string                 `json:"tipoPrefijo"`
	PoliticaCanje           int                    `json:"politicaCanje"`
	PetitorioIst            bool                   `json:"petitorioIst"`
	Comunes                 bool                   `json:"comunes"`
	ComiteFarmacologico     bool                   `json:"comiteFarmacologico"`
	Especializados          bool                   `json:"especializados"`
	EsVital                 bool                   `json:"esVital"`
	EsEstrategico           bool                   `json:"esEstrategico"`
	EsOncologico            bool                   `json:"esOncologico"`
	CargaEnFua              bool                   `json:"cargaEnFua"`
	ProductoIndicadorGestion bool                  `json:"productoIndicadorGestion"`
	PreciosPorPlan          []BienInsumoPrecioPlan `json:"preciosPorPlan"`
}

type BienInsumoFiltros struct {
	Filtro               string
	IdGrupoFarmacologico int
	TipoProducto         int
	SoloPetitorio        *bool
	SoloIndicadorGestion *bool
	Estado               *int
	Page                 int
	PageSize             int
}

type GrupoFarmacologico struct {
	IdGrupoFarmacologico int    `json:"idGrupoFarmacologico"`
	Descripcion          string `json:"descripcion"`
}

type SubGrupoFarmacologico struct {
	IdSubGrupoFarmacologico int    `json:"idSubGrupoFarmacologico"`
	Descripcion             string `json:"descripcion"`
	IdGrupoFarmacologico    int    `json:"idGrupoFarmacologico"`
}

type CentroCostoItem struct {
	IdCentroCosto int    `json:"idCentroCosto"`
	Codigo        string `json:"codigo"`
	Descripcion   string `json:"descripcion"`
}

type PartidaPresupuestalItem struct {
	IdPartidaPresupuestal int    `json:"idPartidaPresupuestal"`
	Codigo                string `json:"codigo"`
	Descripcion           string `json:"descripcion"`
}

type PaisItem struct {
	IdPais int    `json:"idPais"`
	Nombre string `json:"nombre"`
}

type TipoPrefijoItem struct {
	IdTipoSufijos int    `json:"idTipoSufijos"`
	Descripcion   string `json:"descripcion"`
}

type TipoFinanciamientoItem struct {
	IdTipoFinanciamiento int    `json:"idTipoFinanciamiento"`
	Descripcion          string `json:"descripcion"`
}

type CatalogoBienesInsumosAuxiliares struct {
	CentrosCosto    []CentroCostoItem         `json:"centrosCosto"`
	Partidas        []PartidaPresupuestalItem `json:"partidas"`
	Paises          []PaisItem                `json:"paises"`
	Prefijos        []TipoPrefijoItem         `json:"prefijos"`
	Financiamientos []TipoFinanciamientoItem  `json:"financiamientos"`
	Grupos          []GrupoFarmacologico      `json:"grupos"`
}

type BienInsumoGuardarPayload struct {
	Codigo                  string                 `json:"codigo"`
	CodigoSiga              string                 `json:"codigoSiga"`
	CodigoSunat             string                 `json:"codigoSunat"`
	Nombre                  string                 `json:"nombre"`
	NombreComercial         string                 `json:"nombreComercial"`
	IdGrupoFarmacologico    int                    `json:"idGrupoFarmacologico"`
	IdSubGrupoFarmacologico int                    `json:"idSubGrupoFarmacologico"`
	IdPartida               int                    `json:"idPartida"`
	IdCentroCosto           int                    `json:"idCentroCosto"`
	PrecioCompra            float64                `json:"precioCompra"`
	PrecioDistribucion      float64                `json:"precioDistribucion"`
	PrecioDonacion          float64                `json:"precioDonacion"`
	PrecioUltimaCompra      float64                `json:"precioUltimaCompra"`
	StockMinimo             int                    `json:"stockMinimo"`
	TipoProducto            int                    `json:"tipoProducto"`
	Denominacion            string                 `json:"denominacion"`
	Concentracion           string                 `json:"concentracion"`
	Presentacion            string                 `json:"presentacion"`
	FormaFarmaceutica       string                 `json:"formaFarmaceutica"`
	MaterialEnvase          string                 `json:"materialEnvase"`
	PresentacionEnvase      string                 `json:"presentacionEnvase"`
	Fabricante              string                 `json:"fabricante"`
	IdPaisOrigen            int                    `json:"idPaisOrigen"`
	Petitorio               bool                   `json:"petitorio"`
	TipoProductoSismed      string                 `json:"tipoProductoSismed"`
	Estado                  int                    `json:"estado"`
	EstadoDigemid           int                    `json:"estadoDigemid"`
	IdTipoPrefijo           int                    `json:"idTipoPrefijo"`
	PoliticaCanje           int                    `json:"politicaCanje"`
	PetitorioIst            bool                   `json:"petitorioIst"`
	Comunes                 bool                   `json:"comunes"`
	ComiteFarmacologico     bool                   `json:"comiteFarmacologico"`
	Especializados          bool                   `json:"especializados"`
	EsVital                 bool                   `json:"esVital"`
	EsEstrategico           bool                   `json:"esEstrategico"`
	EsOncologico            bool                   `json:"esOncologico"`
	CargaEnFua              bool                   `json:"cargaEnFua"`
	ProductoIndicadorGestion bool                  `json:"productoIndicadorGestion"`
	PreciosPorPlan          []BienInsumoPrecioPlan `json:"preciosPorPlan"`
}

type PaginacionBienesInsumos struct {
	Items      []BienInsumoItem `json:"items"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalItems int              `json:"totalItems"`
	TotalPages int              `json:"totalPages"`
}
