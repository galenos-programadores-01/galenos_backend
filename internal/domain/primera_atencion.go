package domain

import "time"

// PrimeraAtencionUpsert agrupa los campos del SP usp_go_InsertarPrimeraAtencion.
// Punteros nil = no enviar / NULL en actualización parcial.
type PrimeraAtencionUpsert struct {
	IdAtencion                      int
	FechaAtencion                   *time.Time
	TipoAnamnesis                   *string
	TiempoEnfermedad                *int
	IdTiempoEnfermedad              *int
	SignosSintomas                  *string
	RelatoCronologico               *string
	Antecedentes                    *string
	IdEstadoGeneral                 *int
	IdEstadoHidratacion             *int
	IdNivelConciencia               *int
	PielAnexos                      *string
	ExamenClinicoRegional           *string
	ReferenciaOportuna              *string
	Evolucion                       *string
	EvaluacionEnfermeriaObstetricia *string
	IdCondicionEgreso               *int
	CodCondicionEgreso              *string
	FechaEgreso                     *time.Time
	IdUsuarioRegistro               int
	IdUsuarioModifica               *int
	FechaModifica                   *time.Time
	EstadoFirma                     *int
	UbicacionArchivo                *string
}

const (
	// EstadoFirmaPendiente borrador editable.
	EstadoFirmaPendiente = 0
	// EstadoFirmaCerrada PDF firmado con DNIe; ya no se edita.
	EstadoFirmaCerrada = 1
)

// PrimeraAtencionCerrarFirma cierra el registro tras Firma Perú.
type PrimeraAtencionCerrarFirma struct {
	IdAtencion        int
	IdUsuarioModifica int
	UbicacionArchivo  string
}

// PrimeraAtencionCerrarFirmaResult resultado del cierre (mismo SP Insertar).
type PrimeraAtencionCerrarFirmaResult struct {
	IdPrimeraAtencion int
	Mensaje           string
}

// PrimeraAtencionUpsertResult resultado del SP (OUTPUT).
type PrimeraAtencionUpsertResult struct {
	IdPrimeraAtencion int
	Mensaje           string
}

// Clasificaciones DX (usp_go_SubclasificacionDiagnosticos_Select / Agregar).
const (
	ClasificacionDxPA               = 9
	ClasificacionDxInterconsulta    = 7
	ClasificacionDxEgreso           = 3
	ClasificacionDxEgresoFallecido  = 4
)

// DiagnosticoAgregar es la entrada a usp_go_AtencionesDiagnosticosAgregar.
type DiagnosticoAgregar struct {
	IdAtencion                    int
	IdDiagnostico                 int
	IdClasificacionDx             int // 9=PA, 7=IC, 3=egreso, 4=egreso fallecido
	IdSubclasificacionDiagnostico int // viene del SP usp_go_SubclasificacionDiagnosticos_Select
	IdPaciente                    int
	// IdPrimeraAtencion: PA/IC >0; egreso (3/4) debe ir 0 → NULL en el SP.
	IdPrimeraAtencion             int
	IdEpisodio                    int
	IdEvolucion                   int
	LabConfHIS                    *string
	LabConfHIS1                   *string
	LabConfHIS2                   *string
	LabConfHIS3                   *string
	GrupoHIS                      *int
	SubGrupoHIS                   *int
}

// AltaHospitalariaRegistrar — usp_go_RegistrarAltaHospitalaria (domicilio / hosp / fuga / fallecido PA).
type AltaHospitalariaRegistrar struct {
	IdAtencion         int
	IdEpisodio         *int
	IdDestinoAtencion  *int
	IdDestinoServicio  *int // hosp/fuga: NULL; fallecido: destino cadáver (grupo 11)
	IdMedicoEgreso     int
	FechaEgreso        time.Time
	IdCondicionAlta    *int // domicilio: requerido; hosp/fuga/fallecido: NULL
	IdTipoCertificado  *int // fallecido: tipo certificado (grupo 12)
	ActDefuncion       string // fallecido: N° certificado
	IdTipoAlta         *int // domicilio: requerido; hosp/fuga/fallecido: NULL
	IdPronostico       *int
	Recomendacion      string
	IdServicioEgreso   int
	IdCamaEgreso       *int
	IdUsuario          *int
	MotivoTraslado     string     // hospitalización
	FechaFuga           *time.Time // fuga: último contacto (enviar aunque SP aún no persista)
	ReportaFuga         string     // fuga: reportado por
	IdNotificacionFuga  *int       // fuga: 1=seguridad, 2=familiar, 3=incidente
	Observaciones       string     // fuga; fallecido: familiar informado
	IdtiempoEstanciaDeseso *int  // fallecido: tiempo estancia deceso (grupo 13)
	RequiereAutopsia       *bool // fallecido: necropsia si/no
}

// AltaHospitalariaResult — OUTPUT del SP de alta.
type AltaHospitalariaResult struct {
	Resultado int
	Mensaje   string
}

// DiagnosticoAgregarResult parsea el OUTPUT @IdAtencionDiagnostico.
type DiagnosticoAgregarResult struct {
	IdAtencionDiagnostico    int
	SubclasificacionAnterior string
	Raw                      string
}

// TipoDiagnosticoUI es la letra corta en UI: P|D|R (derivada del catálogo en FE).
type TipoDiagnosticoUI string

const (
	TipoDiagnosticoPresuntivo TipoDiagnosticoUI = "P"
	TipoDiagnosticoDefinitivo TipoDiagnosticoUI = "D"
	TipoDiagnosticoRepetido   TipoDiagnosticoUI = "R"
)

// DiagnosticoRegistrarItem ítem del body de registro.
// IdSubclasificacionDx es obligatorio (catálogo SP; no hardcodear 401–406).
type DiagnosticoRegistrarItem struct {
	IdDiagnostico        int
	IdSubclasificacionDx int
	Tipo                 TipoDiagnosticoUI // solo informativo / UI; no se usa para mapear ids
}

// Tipos de ítem en Tab_PA_PlanTrabajoItem (temporal hasta SPs definitivos).
const (
	PlanTrabajoTipoExamen        = "EXA"
	PlanTrabajoTipoInterconsulta = "INT"
	PlanTrabajoTipoProcedimiento = "PRO"
	PlanTrabajoTipoMedicamento   = "MED" // reservado
)

// IdPuntoCargaFarmacia — RecetaCabecera de medicamentos (PA / HIS).
// Exámenes usan otros puntos (lab, patología, rayos, …) del catálogo.
const IdPuntoCargaFarmacia = 5

// IdPuntoCargaProcedimientos — RecetaCabecera fija al guardar procedimientos (PA).
// El catálogo usp_go_ProcedimentosFiltarNombrePuntoCarga filtra con IdPuntoCarga de bandeja.
const IdPuntoCargaProcedimientos = 26

// RecetaDetalleListarFiltro criterios de usp_go_PA_RecetaDetalleListar.
// IdPuntoCarga > 0 → solo ese punto; ExcluirIdPuntoCarga > 0 → todos excepto ese.
// Si ambos > 0, el SP prioriza el filtro exacto.
type RecetaDetalleListarFiltro struct {
	IdCuentaAtencion    int
	IdPrimeraAtencion   int
	IdPuntoCarga        int
	ExcluirIdPuntoCarga int
}

// PlanTrabajoItemAgregar entrada a usp_go_PA_PlanTrabajoItemAgregar.
type PlanTrabajoItemAgregar struct {
	IdAtencion        int
	IdPrimeraAtencion int
	Tipo              string
	IdProducto        *int
	IdEspecialidad    *int
	Codigo            *string
	Nombre            *string
	NumOrden          *int
	IdUsuarioRegistro *int
}

// PlanTrabajoItemAgregarResult OUTPUT @IdPlanTrabajoItem ("OK;<id>").
type PlanTrabajoItemAgregarResult struct {
	IdPlanTrabajoItem int
	Raw               string
}

// PlanTrabajoItemRegistrarItem body de POST (UI).
type PlanTrabajoItemRegistrarItem struct {
	Tipo           string
	IdProducto     *int
	IdEspecialidad *int
	Codigo         string
	Nombre         string
}

// PlanTrabajoItem fila listada desde Tab_PA_PlanTrabajoItem.
type PlanTrabajoItem struct {
	IdPlanTrabajoItem int
	IdAtencion        int
	IdPrimeraAtencion int
	Tipo              string
	IdProducto        *int
	IdEspecialidad    *int
	Codigo            string
	Nombre            string
	NumOrden          int
}

// RecetaCabeceraAgregar entrada a usp_go_PA_RecetaCabeceraAgregar.
// IdPrimeraAtencion es obligatorio: el SP PA aísla cabeceras por PA
// (el SP compartido usp_go_RecetaCabeceraAgregar no se usa desde PA).
type RecetaCabeceraAgregar struct {
	IdPuntoCarga       int
	IdCuentaAtencion   int
	IdServicioReceta   int
	IdMedicoReceta     int // 0 por ahora
	IdProducto         int // 0 en cabecera (exámenes); farmacia puede enviar producto
	IdPaciente         int
	IdUsuarioAuditoria int
	IdEvolucion        int // 0
	IdPrimeraAtencion  int // requerido (>0)
}

// RecetaDetalleAgregar entrada a usp_go_RecetaDetalleAgregar.
type RecetaDetalleAgregar struct {
	IdReceta            int
	IdProducto          int
	Cantidad            int
	Precio              float64
	Total               float64 // Cantidad * Precio (exámenes: 1 * Precio)
	SaldoEnRegistro     int     // 0
	IdDosisRecetada     *int
	Observaciones       string // ''
	IdViaAdministracion *int
	IdUNIDDosisReceta   *int
	IdFrecuencia        *int
	Duracion            int
	// IdDiagnostico → columna RecetaDetalle.IdDiagnostico (FK Diagnosticos).
	IdDiagnostico *int
	// CodigoDiagnostico → columna RecetaDetalle.CodigoDiagnostico (CIE-10, ej. "J18.9").
	CodigoDiagnostico *string
}

// RecetaMedicamentoItem fila listada (cabecera + detalle). Sirve tambien para examenes.
type RecetaMedicamentoItem struct {
	IdReceta            int
	IdPuntoCarga        int
	IdCuentaAtencion    int
	IdServicioReceta    int
	IdPrimeraAtencion   int
	IdProducto          int
	Cantidad            int
	Precio              float64
	Total               float64
	IdDosis             *int
	IdUNIDDosis         *int
	IdFrecuencia        *int
	IdViaAdministracion *int
	Duracion            *int
	Codigo              string
	Nombre              string
	IdDiagnostico       *int   // RecetaDetalle.IdDiagnostico
	CodigoDiagnostico   string // RecetaDetalle.CodigoDiagnostico (CIE-10)
	PuntoCarga          string // FactPuntosCarga.Descripcion (ej. PATOLOGIA CLINICA)
}

// RecetaMedicamentoRegistrar entrada de use case (cabecera ctx + detalle).
type RecetaMedicamentoRegistrar struct {
	IdAtencion          int // Tab_PrimeraAtencion / editable (NO es IdCuentaAtencion)
	IdPuntoCarga        int
	IdCuentaAtencion    int // RecetaCabecera (puede diferir de IdAtencion)
	IdServicioReceta    int
	IdPaciente          int
	IdUsuarioAuditoria  int
	IdPrimeraAtencion   int
	IdProducto          int
	Cantidad            int
	Precio              float64 // precioUnitario del catálogo → RecetaDetalle.Precio
	IdDosis             *int
	IdUNIDDosis         *int
	IdFrecuencia        *int
	IdViaAdministracion *int
	Duracion            int
	IdDiagnostico       int    // → RecetaDetalle.IdDiagnostico
	CodigoDiagnostico   string // → RecetaDetalle.CodigoDiagnostico (CIE-10)
}

// RecetaMedicamentoRegistrarResult resultado POST medicamento.
type RecetaMedicamentoRegistrarResult struct {
	IdReceta   int
	IdProducto int
	Cantidad   int
	Precio     float64
	Total      float64
	Mensaje    string
}

// RecetaMedicamentoCamposActualizar actualiza dosis/unidad/frecuencia/vía
// de un detalle ya insertado (idReceta + idProducto).
type RecetaMedicamentoCamposActualizar struct {
	IdAtencion          int // editable PA
	IdReceta            int
	IdProducto          int
	IdCuentaAtencion    int
	IdPrimeraAtencion   int
	IdDosis             *int
	IdUNIDDosis         *int
	IdFrecuencia        *int
	IdViaAdministracion *int
}

// RecetaExamenRegistrar entrada de use case (cabecera por IdPuntoCarga del examen + detalle).
// Cantidad fija 1; Total = Cantidad * Precio (calculado en backend).
type RecetaExamenRegistrar struct {
	IdAtencion         int // editable PA
	IdPuntoCarga       int // del catálogo webProcedimentosRecetas… (lab/rayos/…)
	IdCuentaAtencion   int // RecetaCabecera
	IdServicioReceta   int
	IdPaciente         int
	IdUsuarioAuditoria int
	IdPrimeraAtencion  int
	IdProducto         int // IdProductoSolo
	Precio             float64
	IdDiagnostico      int    // → RecetaDetalle.IdDiagnostico
	CodigoDiagnostico  string // → RecetaDetalle.CodigoDiagnostico (CIE-10)
}

// RecetaExamenRegistrarResult resultado POST examen / procedimiento.
type RecetaExamenRegistrarResult struct {
	IdReceta   int
	IdProducto int
	Cantidad   int
	Precio     float64
	Total      float64
	Mensaje    string
}

// RecetaProcedimientoRegistrar entrada de use case.
// Cabecera siempre IdPuntoCargaProcedimientos (=26); catálogo usa IdPuntoCarga de bandeja.
type RecetaProcedimientoRegistrar struct {
	IdAtencion         int // editable PA
	IdCuentaAtencion   int // RecetaCabecera
	IdServicioReceta   int
	IdPaciente         int
	IdUsuarioAuditoria int
	IdPrimeraAtencion  int
	IdProducto         int
	Precio             float64
	IdDiagnostico      int
	CodigoDiagnostico  string
}

// ─── Interconsulta PA (usp_go_Interconsulta_*) ───────────────────────────────

// IdEstadoInterconsultaSolicitado — @IdEstado al registrar desde PA.
const IdEstadoInterconsultaSolicitado = 1

// IdTipoSolicitudDxInterconsulta — @IdTipo CHAR(1) → InterconsultaDiagnosticos.IdTipoSolicitud.
// HIS: 'S' = Solicitado (no usar "1"; el resto de interconsultas guardan 'S').
const IdTipoSolicitudDxInterconsulta = "S"

// IdGrupoMotivoInterconsulta — Cat_ParametroClinico (radios paso 5).
const IdGrupoMotivoInterconsulta = 6

// BitsMotivoPorOrdenCatalogo mapea la posición del motivo en el catálogo
// (ordenado por Codigo, luego Id) a los 3 BIT del SP Insertar.
// No depende de IdParametro fijos (37/38/39 pueden cambiar en BD).
// Orden esperado: 0=OpinionDiagnostica, 1=ManejoPaciente, 2=TrasferenciaPaciente.
func BitsMotivoPorOrdenCatalogo(orden int) (opinion, manejo, transferencia bool, ok bool) {
	switch orden {
	case 0:
		return true, false, false, true
	case 1:
		return false, true, false, true
	case 2:
		return false, false, true, true
	default:
		return false, false, false, false
	}
}

// ResolverIdRegistroAtencionInterconsulta — PA alineado a Evolución:
// IdRegistroAtencion = 0 e IdRegAtencion queda NULL (no se inserta).
// Listado PA engrana por IdPrimeraAtencion (+ IdEpisodio se sigue enviando).
// Marcador: PA-INTERC-REGATENCION
func ResolverIdRegistroAtencionInterconsulta(_idAtencion int) int {
	return 0
}

// InterconsultaInsertar entrada a usp_go_Interconsulta_Insertar.
type InterconsultaInsertar struct {
	IdEpisodio               int
	IdRegistroAtencion       int
	IdPrimeraAtencion        int
	IdServicioSolicitante    int
	IdEspecialidadSolicitada int
	IdEmpleadoSolicitante    int
	FechaSolicitud           *string // nil → SP usa GETDATE()
	IdPrioridad              int
	OpinionDiagnostica       bool
	ManejoPaciente           bool
	TrasferenciaPaciente     bool // typo del SP (doble r)
	MotivoClinico            string
	ResumenClinico           string
	IdEstado                 int
	IdUsuarioRegistro        int
}

// InterconsultaDiagnosticoInsertar entrada a usp_go_InterconsultaDiagnostico_Insertar.
type InterconsultaDiagnosticoInsertar struct {
	IdInterconsulta      int
	IdDiagnostico        int
	IdSubclasificacionDx int
	IdTipo               string
	IdEstado             int
}

// InterconsultaRegistrarSolicitud body del caso de uso (1 cabecera + N DX).
type InterconsultaRegistrarSolicitud struct {
	IdAtencion         int // bandeja
	IdEpisodio         int
	IdPrimeraAtencion  int // amarre PA; lo usa el List SP
	IdServicio         int
	IdEspecialidad     int
	IdPrioridad        int
	IdMotivoParametro  int // IdParametro del catálogo (trazabilidad)
	MotivoOrden       int // índice 0..2 tras ordenar catálogo por Codigo
	MotivoClinico      string
	IdEmpleado         int
	Diagnosticos       []InterconsultaDiagnosticoItem
}

// InterconsultaDiagnosticoItem DX de la solicitud (pila local FE).
type InterconsultaDiagnosticoItem struct {
	IdDiagnostico        int
	IdSubclasificacionDx int
}

// InterconsultaListadoItem fila de usp_go_ListarInterconsultasPorAtencion.
type InterconsultaListadoItem struct {
	IdInterconsulta  int
	IdAtencionOrigen int
	IdEspecialidad   int
	IdMedicoDestino  int
	Motivo           string
	FechaSolicitud   string
	Estado           string
	Especialidad     string
	ResumenHC        string
}
