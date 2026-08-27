// Package config centraliza la lectura de configuración desde variables de
// entorno. Es la única capa, junto a main, que conoce de dónde vienen los
// valores externos.
package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config agrupa toda la configuración externa de la aplicación.
type Config struct {
	ServerPort            string
	ServerHost            string
	SQLServerDSN          string
	AllowedOrigins        []string
	DBMaxOpenConns        int
	DBMaxIdleConns        int
	DBConnMaxLifetime     time.Duration
	ReniecApp             string
	ReniecUsuario         string
	ReniecClave           string
	ReniecURL             string
	ReniecTimeout         time.Duration
	SISUsuario            string
	SISClave              string
	SISURL                string
	SISDNIAutorizado      string
	SISTimeout            time.Duration
	FirmaPeruTokenURL     string
	FirmaPeruClientID     string
	FirmaPeruClientSecret string
	FirmaPeruPublicURL    string
	FirmaPeruTimeout      time.Duration
	FirmaPeruSignedDir    string
	SevenZipPath          string
	AuthUsername          string
	AuthPassword          string
	AuthSecret            string
	AuthTTL               time.Duration
}

// Load lee la configuración desde el entorno aplicando valores por defecto
// razonables para desarrollo.
func Load() (*Config, error) {
	dsn := os.Getenv("SQLSERVER_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("SQLSERVER_DSN environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	apiUsername := os.Getenv("API_USERNAME")
	if apiUsername == "" {
		return nil, fmt.Errorf("API_USERNAME environment variable is required")
	}

	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		return nil, fmt.Errorf("API_PASSWORD environment variable is required")
	}

	return &Config{
		ServerPort:            envOrDefault("SERVER_PORT", "8080"),
		ServerHost:            envOrDefault("SERVER_HOST", defaultServerHost()),
		SQLServerDSN:          dsn,
		AllowedOrigins:        allowedOriginsWithServer(envOrDefault("ALLOWED_ORIGINS", "http://localhost:4200"), envOrDefault("SERVER_PORT", "8080"), envOrDefault("SERVER_HOST", defaultServerHost())),
		DBMaxOpenConns:        envIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:        envIntOrDefault("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime:     envDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ReniecApp:             envOrDefault("RENIEC_APP", "HNSEB"),
		ReniecUsuario:         envOrDefault("RENIEC_USUARIO", "44602631"),
		ReniecClave:           os.Getenv("RENIEC_CLAVE"),
		ReniecURL:             envOrDefault("RENIEC_URL", "https://wsvmin.minsa.gob.pe/wsreniecmq/serviciomq.asmx"),
		ReniecTimeout:         envDurationOrDefault("RENIEC_TIMEOUT", 30*time.Second),
		SISUsuario:            envOrDefault("SIS_USUARIO", "HNAL"),
		SISClave:              os.Getenv("SIS_CLAVE"),
		SISURL:                envOrDefault("SIS_URL", "http://app.sis.gob.pe/sisWSAFI/Service.asmx"),
		SISDNIAutorizado:      envOrDefault("SIS_DNI_AUTORIZADO", ""),
		SISTimeout:            envDurationOrDefault("SIS_TIMEOUT", 30*time.Second),
		FirmaPeruTokenURL:     envOrDefault("FIRMAPERU_TOKEN_URL", "https://apps.firmaperu.gob.pe/admin/api/security/generate-token"),
		FirmaPeruClientID:     os.Getenv("FIRMAPERU_CLIENT_ID"),
		FirmaPeruClientSecret: os.Getenv("FIRMAPERU_CLIENT_SECRET"),
		FirmaPeruPublicURL:    os.Getenv("FIRMAPERU_PUBLIC_URL"),
		FirmaPeruTimeout:      envDurationOrDefault("FIRMAPERU_TIMEOUT", 60*time.Second),
		FirmaPeruSignedDir:    os.Getenv("FIRMAPERU_SIGNED_DIR"),
		SevenZipPath:          envOrDefault("SEVENZIP_PATH", `C:\Program Files\7-Zip\7z.exe`),
		AuthUsername:          apiUsername,
		AuthPassword:          apiPassword,
		AuthSecret:            jwtSecret,
		AuthTTL:               envDurationOrDefault("JWT_TTL", 24*time.Hour),
	}, nil
}

// defaultServerHost devuelve el IP de la interfaz de red principal para
// que otros usuarios de la red local puedan consumir la API (en lugar de
// localhost). Si no se puede detectar, retorna localhost.
func defaultServerHost() string {
	host := localIPv4()
	if host == "" {
		return "localhost"
	}
	return host
}

// localIPv4 resuelve el IP IPv4 saliente consultando la interfaz por
// defecto (conexión UDP sin tráfico real). Retorna vacío si falla.
func localIPv4() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return addr.IP.String()
	}
	return ""
}

// allowedOriginsWithServer agrega a la lista de orígenes permitidos el
// origen de la propia API (IP configurado y localhost) y, cuando el host es
// una IP de red, el origen del frontend de desarrollo (http://<IP>:4200),
// para que los usuarios de la red local no fallen por CORS.
func allowedOriginsWithServer(csv, port, host string) []string {
	origins := strings.Split(csv, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}

	appendIfMissing := func(origin string) {
		if origin == "" {
			return
		}
		for _, o := range origins {
			if o == origin {
				return
			}
		}
		origins = append(origins, origin)
	}

	appendIfMissing("http://" + host + ":" + port)
	appendIfMissing("http://localhost:" + port)
	if host != "localhost" {
		appendIfMissing("http://" + host + ":4200")
	}
	return origins
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return parsed
}
