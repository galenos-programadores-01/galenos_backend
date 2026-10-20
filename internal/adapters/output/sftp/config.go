package sftp

import "time"

// Config credenciales y opciones del cliente SFTP (desde env / composition root).
type Config struct {
	Host                   string
	Port                   int
	User                   string
	Password               string
	Timeout                time.Duration
	// InsecureIgnoreHostKey omite verificación de host key (redes privadas/IIS).
	// En producción preferir KnownHostsPath cuando sea posible.
	InsecureIgnoreHostKey bool
	KnownHostsPath        string
}
