// Command api arranca la REST API del sistema hospitalario. Es el único
// punto donde se conectan las implementaciones concretas de los adaptadores
// con los puertos que el dominio define.
package main

// @title Galenos Pro Appointments API
// @version 1.0
// @description REST API del sistema hospitalario Galenos Pro. Arquitectura hexagonal, persiste en SQL Server 2022.
// @host localhost:8080
// @BasePath /api/v1

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/galenos-pro/appointments-api/docs"
	httpadapter "github.com/galenos-pro/appointments-api/internal/adapters/input/http"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/cache"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/firmaperu"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/persistence/sqlserver"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/reniec"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/sis"
	ftpadapter "github.com/galenos-pro/appointments-api/internal/adapters/output/ftp"
	sftpadapter "github.com/galenos-pro/appointments-api/internal/adapters/output/sftp"
	"github.com/galenos-pro/appointments-api/internal/config"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}

func run() error {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// --- Adaptador de salida: persistencia SQL Server ---
	db, err := sqlserver.NewConnection(sqlserver.Config{
		DSN:             cfg.SQLServerDSN,
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	// --- Adaptador de salida: Caché Redis ---
	redisCache := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.RedisTTL)

	appointmentRepo := sqlserver.NewAppointmentRepository(db)
	patientRepo := sqlserver.NewPatientRepository(db)
	catalogRepo := sqlserver.NewCatalogRepository(db)
	triageRepo := sqlserver.NewTriageRepository(db)
	sisRepo := sqlserver.NewSisRepository(db)
	evolucionRepo := sqlserver.NewSqlServerEvolucionRepository(db)
	motivoRepo := sqlserver.NewMotivoRepository(db)
	ordenRepo := sqlserver.NewOrdenRepository(db)
	resultadoRepo := sqlserver.NewResultadoRepository(db)
	interconsultaRepo := sqlserver.NewInterconsultaRepository(db)
	sintomaRepo := sqlserver.NewSintomaRepository(db)
	diagnosticoRepo := sqlserver.NewSqlServerDiagnosticoRepository(db)
	listaEsperaQxRepo := sqlserver.NewListaEsperaQxRepository(db)
	medicoListaEsperaRepo := sqlserver.NewMedicoListaEsperaRepository(db)
	refConRepo := sqlserver.NewRefConRepository(db)
	primeraAtencionRepo := sqlserver.NewPrimeraAtencionRepository(db)

	// --- Adaptador de salida: servicio externo RENIEC ---
	reniecClient := reniec.New(reniec.Config{
		App:     cfg.ReniecApp,
		Usuario: cfg.ReniecUsuario,
		Clave:   cfg.ReniecClave,
		URL:     cfg.ReniecURL,
		Timeout: cfg.ReniecTimeout,
	})

	// --- Adaptador de salida: servicio externo SIS ---
	sisClient := sis.New(sis.Config{
		Usuario:       cfg.SISUsuario,
		Clave:         cfg.SISClave,
		URL:           cfg.SISURL,
		DNIAutorizado: cfg.SISDNIAutorizado,
		Timeout:       cfg.SISTimeout,
	})

	// --- Adaptador de salida: Firmador Web de Firma Perú ---
	firmaPeruClient := firmaperu.New(firmaperu.Config{
		TokenURL:     cfg.FirmaPeruTokenURL,
		ClientID:     cfg.FirmaPeruClientID,
		ClientSecret: cfg.FirmaPeruClientSecret,
		Timeout:      cfg.FirmaPeruTimeout,
	})
	firmaPeruStore := firmaperu.NewMemoryStore(cfg.FirmaPeruSignedDir)
	firmaPeruArchive := firmaperu.NewSevenZipTool(cfg.SevenZipPath)

	// --- Adaptador de salida: almacenamiento remoto (SFTP o FTP; noop si no hay host) ---
	almacenamientoRemoto, err := newAlmacenamientoRemoto(cfg)
	if err != nil {
		return err
	}
	rutasArchivoRemoto := domain.NuevaRutasArchivoRemoto(cfg.SFTPBasePath, map[string]string{
		domain.PatronArchivoPrimeraAtencion: cfg.SFTPPatternPrimeraAtencion,
		domain.PatronArchivoEvolucion:       cfg.SFTPPatternEvolucion,
	})
	archivoRemotoService := usecase.NewArchivoRemotoService(almacenamientoRemoto, rutasArchivoRemoto)
	if archivoRemotoService.Enabled() {
		log.Printf("almacenamiento remoto habilitado protocol=%s host=%s:%d base=%s",
			cfg.SFTPProtocol, cfg.SFTPHost, cfg.SFTPPort, cfg.SFTPBasePath)
	} else {
		log.Printf("almacenamiento remoto deshabilitado (SFTP_HOST/SFTP_USER vacíos)")
	}

	// --- Núcleo de dominio: casos de uso implementando los puertos de entrada ---
	appointmentService := usecase.NewAppointmentUseCase(appointmentRepo)
	patientService := usecase.NewPatientUseCase(patientRepo)
	catalogService := usecase.NewCatalogUseCase(catalogRepo, redisCache)
	reniecService := usecase.NewReniecUseCase(reniecClient)
	sisService := usecase.NewSisUseCase(sisClient, sisRepo)
	firmaPeruService := usecase.NewFirmaPeruUseCase(firmaPeruClient, firmaPeruStore, firmaPeruArchive)
	triageService := usecase.NewTriageUseCase(triageRepo)
	evolucionService := usecase.NewEvolucionUseCase(evolucionRepo)
	motivoService := usecase.NewMotivoService(motivoRepo)
	ordenService := usecase.NewOrdenService(ordenRepo)
	resultadoService := usecase.NewResultadoService(resultadoRepo)
	interconsultaService := usecase.NewInterconsultaService(interconsultaRepo)
	sintomaService := usecase.NewSintomaService(sintomaRepo)
	diagnosticoUseCase := usecase.NewDiagnosticoUseCase(diagnosticoRepo)
	listaEsperaQxService := usecase.NewListaEsperaQxService(listaEsperaQxRepo)
	medicoListaEsperaService := usecase.NewMedicoListaEsperaService(medicoListaEsperaRepo)
	refConService := usecase.NewRefConService(refConRepo)
	primeraAtencionService := usecase.NewPrimeraAtencionUseCase(primeraAtencionRepo)

	authRepo := sqlserver.NewAuthRepository(db)
	authService := usecase.NewAuthUseCase(authRepo, cfg.AuthSecret, cfg.AuthTTL)

	// El host de Swagger se ajusta en runtime al IP detectado de la
	// máquina (o al configurado en SERVER_HOST) para que otros equipos de
	// la red consuman la API por IP en lugar de localhost.
	docs.SwaggerInfo.Host = cfg.ServerHost + ":" + cfg.ServerPort

	// --- Adaptador de entrada: HTTP/REST para Angular ---
	appointmentHandler := httpadapter.NewAppointmentHandler(appointmentService)
	patientHandler := httpadapter.NewPatientHandler(patientService)
	catalogHandler := httpadapter.NewCatalogHandler(catalogService)
	reniecHandler := httpadapter.NewReniecHandler(reniecService)
	sisHandler := httpadapter.NewSisHandler(sisService)
	firmaPeruHandler := httpadapter.NewFirmaPeruHandler(firmaPeruService, cfg.FirmaPeruPublicURL)
	triageHandler := httpadapter.NewTriageHandler(triageService)
	authHandler := httpadapter.NewAuthHandler(authService)
	evolucionHandler := httpadapter.NewEvolucionHandler(evolucionService)
	motivoHandler := httpadapter.NewMotivoHandler(motivoService)
	ordenHandler := httpadapter.NewOrdenHandler(ordenService)
	resultadoHandler := httpadapter.NewResultadoHandler(resultadoService)
	interconsultaHandler := httpadapter.NewInterconsultaHandler(interconsultaService)
	sintomaHandler := httpadapter.NewSintomaHandler(sintomaService)
	diagnosticoHandler := httpadapter.NewDiagnosticoHandler(diagnosticoUseCase)
	listaEsperaQxHandler := httpadapter.NewListaEsperaQxHandler(listaEsperaQxService)
	medicoListaEsperaHandler := httpadapter.NewMedicoListaEsperaHandler(medicoListaEsperaService)
	refConHandler := httpadapter.NewRefConHandler(refConService)
	primeraAtencionHandler := httpadapter.NewPrimeraAtencionHandler(
		primeraAtencionService,
		archivoRemotoService,
		firmaPeruStore,
	)

	router := httpadapter.NewRouter(httpadapter.RouterParams{
		AppointmentHandler:       appointmentHandler,
		PatientHandler:           patientHandler,
		CatalogHandler:           catalogHandler,
		ReniecHandler:            reniecHandler,
		SisHandler:               sisHandler,
		FirmaPeruHandler:         firmaPeruHandler,
		TriageHandler:            triageHandler,
		PrimeraAtencionHandler:   primeraAtencionHandler,
		AuthHandler:              authHandler,
		EvolucionHandler:         evolucionHandler,
		MotivoHandler:            motivoHandler,
		OrdenHandler:             ordenHandler,
		ResultadoHandler:         resultadoHandler,
		InterconsultaHandler:     interconsultaHandler,
		SintomaHandler:           sintomaHandler,
		DiagnosticoHandler:       diagnosticoHandler,
		ListaEsperaQxHandler:     listaEsperaQxHandler,
		MedicoListaEsperaHandler: medicoListaEsperaHandler,
		RefConHandler:            refConHandler,
		AuthService:              authService,
		AllowedOrigins:           cfg.AllowedOrigins,
	})

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on http://%s:%s (swagger: http://%s:%s/swagger/index.html)", cfg.ServerHost, cfg.ServerPort, cfg.ServerHost, cfg.ServerPort)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on :%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case <-stop:
		log.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}

// newAlmacenamientoRemoto elige el adaptador según SFTP_PROTOCOL (sftp|ftp).
// Mismas variables SFTP_* de credenciales; solo cambia el transporte.
func newAlmacenamientoRemoto(cfg *config.Config) (output.AlmacenamientoRemoto, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.SFTPProtocol)) {
	case "ftp":
		return ftpadapter.New(ftpadapter.Config{
			Host:     cfg.SFTPHost,
			Port:     cfg.SFTPPort,
			User:     cfg.SFTPUser,
			Password: cfg.SFTPPassword,
			Timeout:  cfg.SFTPTimeout,
		})
	case "sftp", "":
		return sftpadapter.New(sftpadapter.Config{
			Host:                  cfg.SFTPHost,
			Port:                  cfg.SFTPPort,
			User:                  cfg.SFTPUser,
			Password:              cfg.SFTPPassword,
			Timeout:               cfg.SFTPTimeout,
			InsecureIgnoreHostKey: cfg.SFTPInsecureIgnoreHostKey,
			KnownHostsPath:        cfg.SFTPKnownHostsPath,
		})
	default:
		return nil, fmt.Errorf("SFTP_PROTOCOL inválido %q (use sftp o ftp)", cfg.SFTPProtocol)
	}
}
