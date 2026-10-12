package domain

type AlmacenFarmaciaItem struct {
	IdAlmacen   int    `json:"idAlmacen"`
	Descripcion string `json:"descripcion"`
	Estado      int    `json:"estado"`
}

type DisponibilidadKpis struct {
	TotalProductos              int     `json:"totalProductos"`
	Disponibles                 int     `json:"disponibles"`
	EnRiesgo                    int     `json:"enRiesgo"`
	EnQuiebre                   int     `json:"enQuiebre"`
	Normostock                  int     `json:"normostock"`
	Substock                    int     `json:"substock"`
	Sobrestock                  int     `json:"sobrestock"`
	TotalEsenciales             int     `json:"totalEsenciales"`
	EsencialesDisponibles       int     `json:"esencialesDisponibles"`
	PctDisponibilidadEsenciales float64 `json:"pctDisponibilidadEsenciales"`
	PorcentajeOptimo            float64 `json:"porcentajeOptimo"`
	LotesPorVencer30            int     `json:"lotesPorVencer30"`
	LotesPorVencer60            int     `json:"lotesPorVencer60"`
}

type ProductoDisponibilidadItem struct {
	IdProducto        int     `json:"idProducto"`
	Codigo            string  `json:"codigo"`
	Nombre            string  `json:"nombre"`
	Concentracion     string  `json:"concentracion"`
	FormaFarmaceutica string  `json:"formaFarmaceutica"`
	IdAlmacen         int     `json:"idAlmacen"`
	AlmacenNombre     string  `json:"almacenNombre"`
	Cantidad          float64 `json:"cantidad"`
	Precio            float64 `json:"precio"`
	StockMinimo       int     `json:"stockMinimo"`
	EstadoSemaforo    string  `json:"estadoSemaforo"`
	TotalLotes        int     `json:"totalLotes"`
	Petitorio         bool    `json:"petitorio"`
	SalidasPeriodo    float64 `json:"salidasPeriodo"`
}

type LoteItem struct {
	Lote              string  `json:"lote"`
	FechaVencimiento  string  `json:"fechaVencimiento"`
	Cantidad          float64 `json:"cantidad"`
	Precio            float64 `json:"precio"`
	NombreMarca       string  `json:"nombreMarca"`
	DiasParaVencer    int     `json:"diasParaVencer"`
	EstadoVencimiento string  `json:"estadoVencimiento"`
}

type LoteVencimientoDetalleItem struct {
	IdProducto        int     `json:"idProducto"`
	Codigo            string  `json:"codigo"`
	Nombre            string  `json:"nombre"`
	Lote              string  `json:"lote"`
	FechaVencimiento  string  `json:"fechaVencimiento"`
	DiasParaVencer    int     `json:"diasParaVencer"`
	EstadoVencimiento string  `json:"estadoVencimiento"`
	Cantidad          float64 `json:"cantidad"`
	Precio            float64 `json:"precio"`
	NombreMarca       string  `json:"nombreMarca"`
}

type KardexMovimientoItem struct {
	MovNumero        int     `json:"movNumero"`
	MovTipo          string  `json:"movTipo"`
	FechaCreacion    string  `json:"fechaCreacion"`
	AlmacenOrigen    string  `json:"almacenOrigen"`
	AlmacenDestino   string  `json:"almacenDestino"`
	Lote             string  `json:"lote"`
	FechaVencimiento string  `json:"fechaVencimiento"`
	Cantidad         float64 `json:"cantidad"`
	Precio           float64 `json:"precio"`
	Total            float64 `json:"total"`
	DocumentoTipo    string  `json:"documentoTipo"`
	DocumentoNumero  string  `json:"documentoNumero"`
	Concepto         string  `json:"concepto"`
	CuentaAtencion   string  `json:"cuentaAtencion"`
	Paciente         string  `json:"paciente"`
	Usuario          string  `json:"usuario"`
	Observaciones    string  `json:"observaciones"`
}

type FarmaciaDisponibilidadFiltros struct {
	IdAlmacen      int
	Filtro         string
	EstadoSemaforo string
	SoloConStock   *bool
	SoloEsenciales *bool
	CanastaTipo    string
	FechaInicio       string
	FechaFin          string
	SoloConMovimiento *bool
	Page              int
	PageSize       int
}

type PaginacionDisponibilidad struct {
	Items      []ProductoDisponibilidadItem `json:"items"`
	Page       int                          `json:"page"`
	PageSize   int                          `json:"pageSize"`
	TotalItems int                          `json:"totalItems"`
	TotalPages int                          `json:"totalPages"`
}

type SerieTendenciaSalida struct {
	Nombre string    `json:"nombre"`
	Data   []float64 `json:"data"`
}

type TendenciasSalidasFarmacia struct {
	Periodos []string               `json:"periodos"`
	Series   []SerieTendenciaSalida `json:"series"`
}

type TopStockItem struct {
	IdProducto     int     `json:"idProducto"`
	Codigo         string  `json:"codigo"`
	Nombre         string  `json:"nombre"`
	Cantidad       float64 `json:"cantidad"`
	ConsumoMensual float64 `json:"consumoMensual"`
}

type ItemCanastaEvaluada struct {
	IdProducto        int     `json:"idProducto"`
	Codigo            string  `json:"codigo"`
	Nombre            string  `json:"nombre"`
	FormaFarmaceutica string  `json:"formaFarmaceutica"`
	CantidadTotal     float64 `json:"cantidadTotal"`
	ConsumoMensual    float64 `json:"consumoMensual"`
	EstadoSemaforo    string  `json:"estadoSemaforo"`
	EsDigemid         bool    `json:"esDigemid"`
	EsSis             bool    `json:"esSis"`
}

type FormaFarmaceuticaStockItem struct {
	FormaFarmaceutica string `json:"formaFarmaceutica"`
	TotalItems        int    `json:"totalItems"`
	ConStock          int    `json:"conStock"`
	SinStock          int    `json:"sinStock"`
}

type AlmacenComparativaStock struct {
	IdAlmacen   int    `json:"idAlmacen"`
	Descripcion string `json:"descripcion"`
	Normales    int    `json:"normales"`
	EnRiesgo    int    `json:"enRiesgo"`
	EnQuiebre   int    `json:"enQuiebre"`
	TotalItems  int    `json:"totalItems"`
}

type VencimientoTramoItem struct {
	Rango           string  `json:"rango"`
	CantidadLotes   int     `json:"cantidadLotes"`
	UnidadesTotales float64 `json:"unidadesTotales"`
	NivelAlerta     string  `json:"nivelAlerta"`
}

type EstadisticasFarmaciaDashboard struct {
	TopExistencias           []TopStockItem               `json:"topExistencias"`
	ComparativaAlmacenes     []AlmacenComparativaStock    `json:"comparativaAlmacenes"`
	DistribucionVencimientos []VencimientoTramoItem       `json:"distribucionVencimientos"`
	EsencialesCriticos       []ItemCanastaEvaluada        `json:"esencialesCriticos"`
	EsencialesPorForma       []FormaFarmaceuticaStockItem `json:"esencialesPorForma"`
}
