package domain

// CausaExternaMorbilidad representa un registro del catálogo de causas
// externas de morbilidad (tabla EmergenciaCausaExternaMorbilidad).
type CausaExternaMorbilidad struct {
	IdCausaExternaMorbilidad int    `json:"idCausaExternaMorbilidad"`
	Descripcion              string `json:"descripcion"`
}
