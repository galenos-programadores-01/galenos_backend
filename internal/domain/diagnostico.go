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

type DiagnosticoAtencion struct {
	IdAtencionDiagnostico int    `json:"idAtencionDiagnostico"`
	IdAtencion            int    `json:"idAtencion"`
	IdDiagnostico         int    `json:"idDiagnostico"`
	CodigoCIE10           string `json:"codigoCIE10"`
	Descripcion           string `json:"descripcion"`
	DescripcionLarga      string `json:"descripcionLarga"`
	TipoCodigo            string `json:"tipoCodigo"`
	TipoDx                string `json:"tipoDx"`
	IdSubclasificacionDx  int    `json:"idSubclasificacionDx"`
	LabConfHIS            string `json:"labConfHIS"`
	LabConfHIS1           string `json:"labConfHIS1"`
	LabConfHIS2           string `json:"labConfHIS2"`
	LabConfHIS3           string `json:"labConfHIS3"`
}

type AgregarDiagnosticoAtencionRequest struct {
	IdSubclasificacionDiagnostico int    `json:"idSubclasificacionDiagnostico"`
	IdDiagnostico                 int    `json:"idDiagnostico"`
	IdAtencion                    int    `json:"idAtencion"`
	CodigoCIE10                   string `json:"codigoCIE10"`
	TipoDiagnostico               string `json:"tipoDiagnostico"`
	LabConfHIS                    string `json:"labConfHIS"`
	LabConfHIS1                   string `json:"labConfHIS1"`
	LabConfHIS2                   string `json:"labConfHIS2"`
	LabConfHIS3                   string `json:"labConfHIS3"`
	GrupoHIS                      int    `json:"grupoHIS"`
	SubGrupoHIS                   int    `json:"subGrupoHIS"`
	IdEpisodio                    int    `json:"idEpisodio"`
	IdEvolucion                   int    `json:"idEvolucion"`
	IdPrimeraAtencion             int    `json:"idPrimeraAtencion"`
	IdPaciente                    int    `json:"idPaciente"`
}

type AgregarDiagnosticoAtencionResponse struct {
	IdAtencionDiagnostico       int    `json:"idAtencionDiagnostico"`
	DescripcionSubclasificacion string `json:"descripcionSubclasificacion"`
}
