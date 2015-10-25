package api

import (
	"flag"
)

// Configuration for the HTTP API.
type Config struct {
	Addr     string `json:"bind"`
	TLSCert  string `json:"tls-cert"`
	TLSKey   string `json:"tls-key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Initialize the configuration.
func InitConfig() *Config {
	c := &Config{}
	flag.StringVar(&c.Addr, "bind", ":8025", "`address` and port to bind to")
	flag.StringVar(&c.TLSCert, "tls-cert", "", "certificate `file` for TLS")
	flag.StringVar(&c.TLSKey, "tls-key", "", "private key `file` for TLS")
	flag.StringVar(&c.Username, "username", "", "`username` for HTTP basic auth")
	flag.StringVar(&c.Password, "password", "", "`password` for HTTP basic auth")
	return c
}
