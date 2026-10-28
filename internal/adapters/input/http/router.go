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
	refconhttp "github.com/galenos-pro/appointments-api/internal/refcon/adapters/input/http"
)

type RouterParams struct {
	AppointmentHandler            *AppointmentHandler
	PatientHandler                *PatientHandler
	CatalogHandler                *CatalogHandler
	ReniecHandler                 *ReniecHandler
	SisHandler                    *SisHandler
	TriageHandler                 *TriageHandler
	PrimeraAtencionHandler        *PrimeraAtencionHandler
	AuthHandler                   *AuthHandler
	EvolucionHandler              *EvolucionHandler
	MotivoHandler                 *MotivoHandler
	OrdenHandler                  *OrdenHandler
	ResultadoHandler              *ResultadoHandler
	InterconsultaHandler          *InterconsultaHandler
	SintomaHandler                *SintomaHandler
	FirmaPeruHandler              *FirmaPeruHandler
	DiagnosticoHandler            *DiagnosticoHandler
	ListaEsperaQxHandler          *ListaEsperaQxHandler
	MedicoListaEsperaHandler      *MedicoListaEsperaHandler
	CausaExternaMorbilidadHandler *CausaExternaMorbilidadHandler
	AuditoriaHandler              *AuditoriaHandler
	RefConHandler                 *refconhttp.RefConHandler
	AuthService                   input.AuthService
	AllowedOrigins                []string
}

func NewRouter(p RouterParams) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "Cache-Control", "X-Requested-With", "X-Cache-Engine", "X-Cache-Provider"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Cache-Engine", "X-Cache-Provider"},
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
			authProtected.GET("/perfil", p.AuthHandler.GetPerfil)

			pacientes := protected.Group("/pacientes")
			pacientes.GET("", p.PatientHandler.List)
			pacientes.POST("", p.PatientHandler.Create)
			pacientes.POST("/historia-clinica", p.PatientHandler.CreateWithHistory)
			pacientes.GET("/buscar", p.PatientHandler.Search)
			pacientes.GET("/por-documento", p.PatientHandler.GetByDocumentAndType)
			pacientes.GET("/:idOrDoc/datos-adicionales", p.PatientHandler.GetDatosAdicionales)
			pacientes.PUT("/:idOrDoc/datos-adicionales", p.PatientHandler.UpdateDatosAdicionales)
			pacientes.GET("/:idOrDoc", p.PatientHandler.Get)
			pacientes.PUT("/:idOrDoc", p.PatientHandler.Update)
			pacientes.DELETE("/:idOrDoc", p.PatientHandler.Delete)

			appointments := protected.Group("/appointments")
			{
				appointments.POST("", p.AppointmentHandler.Create)
				appointments.GET("/:id", p.AppointmentHandler.GetByID)
			}

			// Bandeja Primera Atención (emergencia) — SP web_go_GestionPrimeraAtencionFiltrar
			primeraAtencion := protected.Group("/primera-atencion")
			{
				primeraAtencion.GET("", p.PrimeraAtencionHandler.ListarBandeja)
				primeraAtencion.POST("/iniciar", p.PrimeraAtencionHandler.IniciarBorrador)
				primeraAtencion.GET("/detalle/:idAtencion", p.PrimeraAtencionHandler.ObtenerDetalle)
				primeraAtencion.PATCH("/:idAtencion", p.PrimeraAtencionHandler.GuardarParcial)
				primeraAtencion.POST("/:idAtencion/cerrar-firma", p.PrimeraAtencionHandler.CerrarFirma)
				primeraAtencion.GET("/:idAtencion/documento", p.PrimeraAtencionHandler.DescargarDocumento)
				primeraAtencion.GET("/:idAtencion/documentos-firmados", p.PrimeraAtencionHandler.ListarDocumentosFirmados)
				primeraAtencion.GET("/:idAtencion/diagnosticos", p.PrimeraAtencionHandler.ListarDiagnosticos)
				primeraAtencion.POST("/:idAtencion/diagnosticos", p.PrimeraAtencionHandler.RegistrarDiagnosticos)
				primeraAtencion.DELETE("/:idAtencion/diagnosticos/:idAtencionDiagnostico", p.PrimeraAtencionHandler.EliminarDiagnostico)
				primeraAtencion.POST("/:idAtencion/egreso/alta-domicilio", p.PrimeraAtencionHandler.RegistrarAltaDomicilio)
				primeraAtencion.POST("/:idAtencion/egreso/alta-hospitalizacion", p.PrimeraAtencionHandler.RegistrarAltaHospitalizacion)
				primeraAtencion.POST("/:idAtencion/egreso/alta-fuga", p.PrimeraAtencionHandler.RegistrarAltaFuga)
				primeraAtencion.POST("/:idAtencion/egreso/alta-fallecido", p.PrimeraAtencionHandler.RegistrarAltaFallecido)
				primeraAtencion.GET("/:idAtencion/plan-trabajo", p.PrimeraAtencionHandler.ListarPlanTrabajoItems)
				primeraAtencion.POST("/:idAtencion/plan-trabajo", p.PrimeraAtencionHandler.RegistrarPlanTrabajoItem)
				primeraAtencion.DELETE("/:idAtencion/plan-trabajo/:idPlanTrabajoItem", p.PrimeraAtencionHandler.EliminarPlanTrabajoItem)
				primeraAtencion.GET("/:idAtencion/medicamentos", p.PrimeraAtencionHandler.ListarMedicamentos)
				primeraAtencion.POST("/:idAtencion/medicamentos", p.PrimeraAtencionHandler.RegistrarMedicamento)
				primeraAtencion.PATCH("/:idAtencion/medicamentos/:idReceta/:idProducto", p.PrimeraAtencionHandler.ActualizarMedicamentoCampos)
				primeraAtencion.DELETE("/:idAtencion/medicamentos/:idReceta/:idProducto", p.PrimeraAtencionHandler.EliminarMedicamento)
				primeraAtencion.GET("/:idAtencion/examenes", p.PrimeraAtencionHandler.ListarExamenes)
				primeraAtencion.POST("/:idAtencion/examenes", p.PrimeraAtencionHandler.RegistrarExamen)
				primeraAtencion.DELETE("/:idAtencion/examenes/:idReceta/:idProducto", p.PrimeraAtencionHandler.EliminarExamen)
				primeraAtencion.GET("/:idAtencion/procedimientos", p.PrimeraAtencionHandler.ListarProcedimientos)
				primeraAtencion.POST("/:idAtencion/procedimientos", p.PrimeraAtencionHandler.RegistrarProcedimiento)
				primeraAtencion.DELETE("/:idAtencion/procedimientos/:idReceta/:idProducto", p.PrimeraAtencionHandler.EliminarProcedimiento)
				primeraAtencion.GET("/:idAtencion/interconsultas", p.PrimeraAtencionHandler.ListarInterconsultas)
				primeraAtencion.POST("/:idAtencion/interconsultas", p.PrimeraAtencionHandler.RegistrarSolicitudInterconsulta)
				primeraAtencion.DELETE("/:idAtencion/interconsultas/:idInterconsulta", p.PrimeraAtencionHandler.EliminarSolicitudInterconsulta)
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
		v1.GET("/establecimientos", p.CatalogHandler.ListEstablecimientos)
		v1.GET("/especialidades-qx", p.CatalogHandler.HandleListarEspecialidadesQx)
		v1.GET("/especialidades-departamento/:idDepartamento", p.CatalogHandler.HandleListarEspecialidadesPorDepartamento)
		v1.GET("/especialidades-por-tipo-servicio/:idTipoServicio", p.CatalogHandler.ListEspecialidadesPorTipoServicio)
		v1.GET("/especialidades-listar", p.CatalogHandler.ListEspecialidadesListar)
		v1.GET("/servicios-prioridad/:idPrioridad", p.CatalogHandler.HandleListarServiciosPorPrioridad)
		v1.GET("/servicios-prioridad/:idPrioridad/:fechaNac", p.CatalogHandler.HandleListarServiciosPorPrioridad)
		v1.GET("/parametros/:idParametro", p.CatalogHandler.GetParametro)
		// --- Contrato PA ---
		v1.GET("/parametros-clinicos/:idGrupo", p.CatalogHandler.ListParametrosClinicos)
		v1.GET("/prioridades-primera-atencion", p.CatalogHandler.ListPrioridadesPrimeraAtencion)
		v1.GET("/subclasificaciones-diagnostico", p.CatalogHandler.ListSubclasificacionesDiagnostico)
		v1.GET("/diagnosticos-catalogo", p.CatalogHandler.ListDiagnosticosCatalogo)
		v1.GET("/medicamentos-catalogo", p.CatalogHandler.ListMedicamentosCatalogo)
		v1.GET("/examenes-catalogo", p.CatalogHandler.ListExamenesCatalogo)
		v1.GET("/procedimientos-catalogo", p.CatalogHandler.ListProcedimientosCatalogo)
		v1.GET("/receta-dosis", p.CatalogHandler.ListRecetaDosis)
		v1.GET("/receta-unidades-dosis", p.CatalogHandler.ListRecetaUnidadesDosis)
		v1.GET("/receta-frecuencias", p.CatalogHandler.ListRecetaFrecuencias)
		v1.GET("/vias-administracion", p.CatalogHandler.ListViasAdministracion)

		v1.GET("/receta/frecuencias", p.CatalogHandler.HandleListRecetaFrecuencias)
		v1.GET("/receta/unidades-dosis", p.CatalogHandler.HandleListRecetaUnidadesDosis)
		v1.GET("/receta/vias-administracion", p.CatalogHandler.HandleListRecetaViasAdministracion)
		v1.GET("/receta/medicamentos", p.CatalogHandler.BuscarMedicamentosReceta)
		v1.GET("/catalogos/examenes", p.CatalogHandler.BuscarExamenesCatalogo)
		v1.GET("/catalogos/parametros-clinicos/:idGrupo", p.CatalogHandler.HandleListParametrosClinicos)
		v1.PATCH("/catalogos/parametros/:idParametro/valor-texto", p.CatalogHandler.ActualizarParametroValorTexto)

		reniec := v1.Group("/reniec")
		{
			reniec.GET("/:nrodoc", p.ReniecHandler.Consultar)
		}

		sis := v1.Group("/sis")
		{
			sis.GET("/afiliado/:nrodoc", p.SisHandler.ConsultarAfiliado)
			sis.POST("/filiaciones", p.SisHandler.GestionarAfiliacion)
			sis.GET("/filiaciones", p.SisHandler.ListFiliaciones)
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
			triaje.GET("/:id", p.TriageHandler.GetTriajePorId)
			triaje.GET("/medicos/:IdEspecialidad", p.TriageHandler.ListMedicosPorEspecialidad)
			triaje.GET("/consulta", p.TriageHandler.ListTriajeConsulta)
			triaje.GET("/consulta/atencion/:idAtencion", p.TriageHandler.GetTriajeConsultaPorAtencion)
			triaje.PUT("/consulta/:id/estado", p.TriageHandler.UpdateEstadoTriajeConsulta)
			triaje.PUT("/:id/iafa", p.TriageHandler.UpdateIafa)
			triaje.PUT("/:id", p.TriageHandler.UpdateTriaje)
			triaje.DELETE("/:id", p.TriageHandler.Delete)
			triaje.POST("", p.TriageHandler.Create)
			triaje.POST("/consulta", p.TriageHandler.CreateTriajeConsulta)
			triaje.POST("/admision", p.TriageHandler.CreateAdmission)
			triaje.GET("/reporte-por-empleado", p.TriageHandler.GetReporteTriaje)
			triaje.GET("/referencias", p.TriageHandler.ListReferenciasConsultaExterna)
			triaje.GET("/referencias/:idAtencion", p.TriageHandler.ListarDatosReferencia)
			triaje.GET("/referencias-cabecera/:idCuentaAtencion", p.TriageHandler.WebReferenciaJSONCab)
			triaje.GET("/referencias-detalle/:idCuentaAtencion", p.TriageHandler.WebReferenciaJSONDet)
			triaje.GET("/referencias-tratamiento/:idCuentaAtencion", p.TriageHandler.WebTratamientoReferenciaJSON)
			triaje.GET("/referencias-cpt/:idCuentaAtencion", p.TriageHandler.WebCPTReferenciaJSON)
			triaje.GET("/causas-externas-morbilidad", p.CausaExternaMorbilidadHandler.HandleListar)
		}

		evoluciones := protected.Group("/evoluciones")
		{
			evoluciones.POST("", p.EvolucionHandler.HandleCreateEvolucion)
			evoluciones.GET("/pacientes", p.EvolucionHandler.HandleListPatients)
			evoluciones.GET("/paciente/:pacienteId", p.EvolucionHandler.HandleListEvoluciones)
			evoluciones.GET("/bandeja", p.EvolucionHandler.HandleBandeja)
			evoluciones.POST("/registro", p.EvolucionHandler.HandleInsertEvolucionMedica)
			evoluciones.GET("/:idRegAtencion/motivos", p.MotivoHandler.HandleListMotivos)
			evoluciones.POST("/:idRegAtencion/motivos", p.MotivoHandler.HandleCreateMotivo)
			evoluciones.GET("/:idRegAtencion/sintomas", p.SintomaHandler.HandleObtenerAtencionSintomas)
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
			resultados.GET("/laboratorio/detalle", p.ResultadoHandler.HandleObtenerDetalleLaboratorio)
			resultados.GET("/imagenes/paciente/:idPaciente", p.ResultadoHandler.HandleListResultadosImagenes)
			resultados.GET("/imagenes/detalle", p.ResultadoHandler.HandleObtenerDetalleImagen)
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

		diagnosticos := v1.Group("/diagnosticos")
		{
			diagnosticos.GET("/search", p.DiagnosticoHandler.SearchDiagnosticos)
			diagnosticos.GET("/listar", p.DiagnosticoHandler.HandleListarDiagnosticos)
			diagnosticos.GET("/atencion/:idAtencion", p.DiagnosticoHandler.HandleObtenerDiagnosticosAtencion)
		}

		diagnosticosProtected := protected.Group("/diagnosticos")
		{
			diagnosticosProtected.POST("/atencion", p.DiagnosticoHandler.HandleAgregarDiagnosticoAtencion)
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

		dashrefcon := protected.Group("/dashrefcon")
		{
			dashrefcon.GET("/referencias", p.RefConHandler.HandleListarReferencias)
			dashrefcon.GET("/ups", p.RefConHandler.HandleListarUps)
			dashrefcon.GET("/upss/:codigoRenipress", p.RefConHandler.HandleListadoUpss)
			dashrefcon.GET("/especialidades-minsa", p.RefConHandler.HandleListadoEspecialidades)
			dashrefcon.GET("/consulta-referencia-detalle", p.RefConHandler.HandleConsultarReferenciaDetalle)
			dashrefcon.POST("/enviar-referencia", p.RefConHandler.HandleSaveReferencia)
		}

		refcon := v1.Group("/refcon")
		{
			refcon.GET("/ups", p.RefConHandler.HandleListarTodasLasUps)
			refcon.GET("/establecimientos", p.RefConHandler.HandleListarEstablecimientos)
			refcon.GET("/distritos-reniec/:idReniec", p.RefConHandler.HandleListarDistritosPorIdReniec)
			refcon.POST("/hoja-referencia", p.RefConHandler.HandleGenerarHojaReferencia)
		}

		auditoria := protected.Group("/auditoria")
		{
			auditoria.POST("", p.AuditoriaHandler.HandleRegistrar)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
