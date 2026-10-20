package ftp

import "time"

// Config credenciales y opciones del cliente FTP/FTPS (puerto 21 tipico).
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Timeout  time.Duration
	// TLSMode: "plain" | "explicit" | "implicit".
	// explicit = AUTH TLS en :21 (IIS / WinSCP "TLS/SSL explicito").
	TLSMode string
	// InsecureSkipVerify omite validacion del certificado TLS (red privada).
	InsecureSkipVerify bool
}