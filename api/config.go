package api

// Configuration for the HTTP API.
type Config struct {
	Addr     string
	TLSCert  string
	TLSKey   string
	Username string
	Password string
}
