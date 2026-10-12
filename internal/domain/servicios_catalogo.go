package domain

type ServicioCatalogoItem struct {
	IdProducto       int     `json:"idProducto"`
	Codigo           string  `json:"codigo"`
	CodigoHis        string  `json:"codigoHis"`
	Nombre           string  `json:"nombre"`
	NombreCorto      string  `json:"nombreCorto"`
	IdServicioGrupo  int     `json:"idServicioGrupo"`
	GrupoNombre      string  `json:"grupoNombre"`
	CodMinsa         string  `json:"codMinsa"`
	NombreMinsa      string  `json:"nombreMinsa"`
	EsCpt            bool    `json:"esCpt"`
	EsQuirurgico     bool    `json:"esQuirurgico"`
	CodigoSis        string  `json:"codigoSis"`
	CodigoSoat       string  `json:"codigoSoat"`
	Cpms             string  `json:"cpms"`
	PrecioParticular float64 `json:"precioParticular"`
	PrecioSis        float64 `json:"precioSis"`
	IdEstado         int     `json:"idEstado"`
	CargaEnFua       bool    `json:"cargaEnFua"`
	CargaEnHis       bool    `json:"cargaEnHis"`
}

type TarifaFinanciamiento struct {
	IdTipoFinanciamiento int     `json:"idTipoFinanciamiento"`
	FinanciamientoNombre string  `json:"financiamientoNombre"`
	PrecioUnitario       float64 `json:"precioUnitario"`
	Activo               bool    `json:"activo"`
}

type ServicioCatalogoDetalle struct {
	IdProducto       int                    `json:"idProducto"`
	Codigo           string                 `json:"codigo"`
	CodigoHis        string                 `json:"codigoHis"`
	CodigoSunat      string                 `json:"codigoSunat"`
	Nombre           string                 `json:"nombre"`
	NombreCorto      string                 `json:"nombreCorto"`
	IdServicioGrupo  int                    `json:"idServicioGrupo"`
	GrupoNombre      string                 `json:"grupoNombre"`
	CodMinsa         string                 `json:"codMinsa"`
	NombreMinsa      string                 `json:"nombreMinsa"`
	EsCpt            bool                   `json:"esCpt"`
	EsQuirurgico     bool                   `json:"esQuirurgico"`
	CodigoSis        string                 `json:"codigoSis"`
	CodigoSoat       string                 `json:"codigoSoat"`
	Cpms             string                 `json:"cpms"`
	PrecioParticular float64                `json:"precioParticular"`
	PrecioSis        float64                `json:"precioSis"`
	IdEstado         int                    `json:"idEstado"`
	CargaEnFua       bool                   `json:"cargaEnFua"`
	CargaEnHis       bool                   `json:"cargaEnHis"`
	Tarifas          []TarifaFinanciamiento `json:"tarifas"`
}

type ServicioGrupo struct {
	IdServicioGrupo int    `json:"idServicioGrupo"`
	Codigo          string `json:"codigo"`
	Descripcion     string `json:"descripcion"`
}

type ServicioCatalogoFiltros struct {
	Filtro          string
	IdServicioGrupo int
	SoloCpt         *bool
	SoloQuirurgico  *bool
	Estado          *int
	Page            int
	PageSize        int
}

type PaginacionServiciosCatalogo struct {
	Items      []ServicioCatalogoItem `json:"items"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"pageSize"`
	TotalItems int                    `json:"totalItems"`
	TotalPages int                    `json:"totalPages"`
}

type ServicioCatalogoGuardarPayload struct {
	Codigo           string  `json:"codigo"`
	CodigoHis        string  `json:"codigoHis"`
	CodigoSunat      string  `json:"codigoSunat"`
	Nombre           string  `json:"nombre"`
	NombreCorto      string  `json:"nombreCorto"`
	IdServicioGrupo  int     `json:"idServicioGrupo"`
	CodMinsa         string  `json:"codMinsa"`
	NombreMinsa      string  `json:"nombreMinsa"`
	EsCpt            bool    `json:"esCpt"`
	EsQuirurgico     bool    `json:"esQuirurgico"`
	CodigoSis        string  `json:"codigoSis"`
	CodigoSoat       string  `json:"codigoSoat"`
	Cpms             string  `json:"cpms"`
	PrecioParticular float64 `json:"precioParticular"`
	PrecioSis        float64 `json:"precioSis"`
	IdEstado         int     `json:"idEstado"`
	CargaEnFua       bool    `json:"cargaEnFua"`
	CargaEnHis       bool    `json:"cargaEnHis"`
}
