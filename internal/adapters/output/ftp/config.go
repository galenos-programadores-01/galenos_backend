package ftp

import "time"

// Config credenciales y opciones del cliente FTP clasico (puerto 21 tipico).
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Timeout  time.Duration
}