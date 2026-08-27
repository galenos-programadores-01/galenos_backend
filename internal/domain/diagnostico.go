package domain

type DiagnosticoBusqueda struct {
	IdDiagnostico     int    `json:"idDiagnostico"`
	Intrahospitalario int    `json:"intrahospitalario"`
	Descripcion       string `json:"descripcion"`
	CodigoCIE10       string `json:"codigoCIE10"`
	EsActivo          int    `json:"esActivo"`
	DescripcionLarga  string `json:"descripcionLarga"`
	EdadMaxDias       int    `json:"edadMaxDias"`
	EdadMinDias       int    `json:"edadMinDias"`
	IdTipoSexo        int    `json:"idTipoSexo"`
	Cancer            int    `json:"cancer"`
	YaRegistrado      int    `json:"yaRegistrado"`
}

type DiagnosticoSimple struct {
	IdDiagnostico int    `json:"idDiagnostico"`
	CodigoCIE10   string `json:"codigoCIE10"`
	Descripcion   string `json:"descripcion"`
}
