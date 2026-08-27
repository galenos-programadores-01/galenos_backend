package httpadapter

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/galenos-pro/appointments-api/docs"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

type RouterParams struct {
	AppointmentHandler   *AppointmentHandler
	PatientHandler       *PatientHandler
	CatalogHandler       *CatalogHandler
	ReniecHandler        *ReniecHandler
	SisHandler           *SisHandler
	TriageHandler        *TriageHandler
	AuthHandler          *AuthHandler
	EvolucionHandler     *EvolucionHandler
	MotivoHandler        *MotivoHandler
	OrdenHandler         *OrdenHandler
	ResultadoHandler     *ResultadoHandler
	InterconsultaHandler *InterconsultaHandler
	SintomaHandler       *SintomaHandler
	FirmaPeruHandler     *FirmaPeruHandler
	DiagnosticoHandler   *DiagnosticoHandler
	ListaEsperaQxHandler *ListaEsperaQxHandler
	MedicoListaEsperaHandler *MedicoListaEsperaHandler
	AuthService          input.AuthService
	AllowedOrigins       []string
}

func NewRouter(p RouterParams) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     p.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", p.AuthHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(RequireBearerToken(p.AuthService))
		{
			authProtected := protected.Group("/auth")
			authProtected.GET("/menus", p.AuthHandler.GetMenus)

			pacientes := protected.Group("/pacientes")
			pacientes.GET("", p.PatientHandler.List)
			pacientes.POST("", p.PatientHandler.Create)
			pacientes.GET("/buscar", p.PatientHandler.Search)
			pacientes.GET("/por-documento", p.PatientHandler.GetByDocumentAndType)
			pacientes.GET("/:idOrDoc", p.PatientHandler.Get)
			pacientes.PUT("/:id", p.PatientHandler.Update)
			pacientes.DELETE("/:id", p.PatientHandler.Delete)

			appointments := protected.Group("/appointments")
			{
				appointments.POST("", p.AppointmentHandler.Create)
				appointments.GET("/:id", p.AppointmentHandler.GetByID)
			}
		}

		etnias := v1.Group("/etnias")
		{
			etnias.GET("", p.CatalogHandler.ListEtnias)
		}

		idiomas := v1.Group("/idiomas")
		{
			idiomas.GET("", p.CatalogHandler.ListIdiomas)
		}

		v1.GET("/tipos-sexo", p.CatalogHandler.ListTipoSexos)
		v1.GET("/estados-civil", p.CatalogHandler.ListEstadosCivil)
		v1.GET("/grados-instruccion", p.CatalogHandler.ListGradosInstruccion)
		v1.GET("/ocupaciones", p.CatalogHandler.ListOcupaciones)
		v1.GET("/tipos-documentos", p.CatalogHandler.ListTiposDocumentos)

		v1.GET("/departamentos", p.CatalogHandler.ListDepartamentos)
		v1.GET("/provincias/:id", p.CatalogHandler.ListProvincias)
		v1.GET("/distritos/:id", p.CatalogHandler.ListDistritos)
		v1.GET("/centros-poblados/:id", p.CatalogHandler.ListCentrosPoblados)
		v1.GET("/paises", p.CatalogHandler.ListPaises)
		v1.GET("/estados-llego-paciente", p.CatalogHandler.ListEstadosLlegoPaciente)
		v1.GET("/fuentes-financiamiento", p.CatalogHandler.ListFuentesFinanciamiento)
		v1.GET("/servicios/:idTipoServicio", p.CatalogHandler.ListServicios)
		v1.GET("/datos-institucion", p.CatalogHandler.GetDatosInstitucion)
		v1.GET("/especialidades", p.CatalogHandler.ListEspecialidades)
		v1.GET("/especialidades-qx", p.CatalogHandler.HandleListarEspecialidadesQx)
		v1.GET("/especialidades-departamento/:idDepartamento", p.CatalogHandler.HandleListarEspecialidadesPorDepartamento)
		v1.GET("/parametros/:idParametro", p.CatalogHandler.GetParametro)

		reniec := v1.Group("/reniec")
		{
			reniec.GET("/:nrodoc", p.ReniecHandler.Consultar)
		}

		sis := v1.Group("/sis")
		{
			sis.GET("/afiliado/:nrodoc", p.SisHandler.ConsultarAfiliado)
			sis.GET("/filiaciones", p.SisHandler.BuscarAfiliacion)
			sis.POST("/filiaciones", p.SisHandler.GestionarAfiliacion)
			sis.POST("/fua", p.SisHandler.ForzarGuardadoFua)
			sis.POST("/fua/agregar", p.SisHandler.AgregarFua)
			sis.GET("/fua/imprimir", p.SisHandler.GetFuaImprimir)
			sis.GET("/diagnosticos", p.SisHandler.ListDiagnosticos)
			sis.GET("/medicamentos", p.SisHandler.ListMedicamentos)
			sis.GET("/procedimientos", p.SisHandler.ListProcedimientos)
			sis.GET("/consumo", p.SisHandler.ListConsumo)
		}

		firmaperu := v1.Group("/firmaperu")
		{
			firmaperu.POST("/firmar", p.FirmaPeruHandler.Firmar)
			firmaperu.POST("/lote", p.FirmaPeruHandler.FirmarLote)
			firmaperu.POST("/params/:uuid", p.FirmaPeruHandler.ParametrosFirma)
			firmaperu.GET("/estampado.png", p.FirmaPeruHandler.Estampado)
			firmaperu.GET("/estampado/:uuid", p.FirmaPeruHandler.EstampadoUuid)
			firmaperu.GET("/documentos/:uuid", p.FirmaPeruHandler.DescargarDocumento)
			firmaperu.POST("/documentos/:uuid", p.FirmaPeruHandler.RecibirDocumentoFirmado)
			firmaperu.GET("/documentos/:uuid/firmado", p.FirmaPeruHandler.DescargarDocumentoFirmado)
			firmaperu.GET("/documentos/:uuid/lote", p.FirmaPeruHandler.DescargarLoteDocumentos)
			firmaperu.POST("/documentos/:uuid/lote", p.FirmaPeruHandler.RecibirLoteFirmado)
			firmaperu.GET("/documentos/:uuid/lote/firmado", p.FirmaPeruHandler.DescargarLoteFirmado)
		}

		triaje := v1.Group("/triaje")
		{
			triaje.GET("", p.TriageHandler.List)
			triaje.GET("/pendientes-admision", p.TriageHandler.ListPendingAdmission)
			triaje.GET("/reporte", p.TriageHandler.GetReporte)
			triaje.GET("/ficha-admision", p.TriageHandler.GetFichaAdmision)
			triaje.GET("/medicos/:IdEspecialidad", p.TriageHandler.ListMedicosPorEspecialidad)
			triaje.GET("/consulta", p.TriageHandler.ListTriajeConsulta)
			triaje.GET("/consulta/atencion/:idAtencion", p.TriageHandler.GetTriajeConsultaPorAtencion)
			triaje.PUT("/consulta/:id/estado", p.TriageHandler.UpdateEstadoTriajeConsulta)
			triaje.POST("", p.TriageHandler.Create)
			triaje.POST("/consulta", p.TriageHandler.CreateTriajeConsulta)
			triaje.POST("/admision", p.TriageHandler.CreateAdmission)
		}

		evoluciones := protected.Group("/evoluciones")
		{
			evoluciones.POST("", p.EvolucionHandler.HandleCreateEvolucion)
			evoluciones.GET("/pacientes", p.EvolucionHandler.HandleListPatients)
			evoluciones.GET("/paciente/:pacienteId", p.EvolucionHandler.HandleListEvoluciones)
			evoluciones.GET("/:idRegAtencion/motivos", p.MotivoHandler.HandleListMotivos)
			evoluciones.POST("/:idRegAtencion/motivos", p.MotivoHandler.HandleCreateMotivo)
			evoluciones.POST("/:idRegAtencion/sintomas", p.SintomaHandler.HandleGuardarSintomas)
		}

		ordenes := protected.Group("/ordenes")
		{
			ordenes.GET("/cuenta/:idCuentaAtencion", p.OrdenHandler.HandleListOrdenes)
			ordenes.POST("", p.OrdenHandler.HandleCreateOrden)
			ordenes.GET("/productos", p.OrdenHandler.HandleBuscarProductos)
		}

		resultados := protected.Group("/resultados")
		{
			resultados.GET("/laboratorio/paciente/:idPaciente", p.ResultadoHandler.HandleListResultadosLaboratorio)
			resultados.GET("/imagenes/paciente/:idPaciente", p.ResultadoHandler.HandleListResultadosImagenes)
		}

		interconsultas := protected.Group("/interconsultas")
		{
			interconsultas.GET("/especialidades", p.InterconsultaHandler.HandleListarEspecialidades)
			interconsultas.GET("/medicos/:IdEspecialidad", p.InterconsultaHandler.HandleListarMedicosPorEspecialidad)
			interconsultas.GET("/:id", p.InterconsultaHandler.HandleObtenerPorId)
			interconsultas.GET("/servicio/:tipoServicio", p.InterconsultaHandler.HandleListarPorServicio)
			interconsultas.GET("/atencion/:idAtencion", p.InterconsultaHandler.HandleListarPorAtencion)
			interconsultas.POST("", p.InterconsultaHandler.HandleCrear)
			interconsultas.PUT("/:id/estado", p.InterconsultaHandler.HandleActualizarEstado)
			interconsultas.POST("/:id/firma", p.InterconsultaHandler.HandleGuardarFirma)
		}

		sintomas := protected.Group("/sintomas")
		{
			sintomas.GET("/catalogo", p.SintomaHandler.HandleListarCatalogo)
			sintomas.POST("/catalogo", p.SintomaHandler.HandleAgregarCatalogo)
		}

		diagnosticos := protected.Group("/diagnosticos")
		{
			diagnosticos.GET("/search", p.DiagnosticoHandler.SearchDiagnosticos)
			diagnosticos.GET("/listar", p.DiagnosticoHandler.HandleListarDiagnosticos)
		}

		listaEsperaQx := protected.Group("/lista-espera-qx")
		{
			listaEsperaQx.GET("", p.ListaEsperaQxHandler.HandleListar)
			listaEsperaQx.GET("/reporte", p.ListaEsperaQxHandler.HandleReporte)
			listaEsperaQx.GET("/:id", p.ListaEsperaQxHandler.HandleObtenerPorId)
			listaEsperaQx.POST("", p.ListaEsperaQxHandler.HandleCrear)
			listaEsperaQx.PUT("/:id", p.ListaEsperaQxHandler.HandleModificar)
		}

		v1.GET("/medicos-lista-espera", p.MedicoListaEsperaHandler.HandleListar)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
