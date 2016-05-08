package api

// Configuration for the HTTP API.
type Config struct {
	Addr       string `json:"bind"`
	CORSOrigin string `json:"cors-origin"`
	TLSCert    string `json:"tls-cert"`
	TLSKey     string `json:"tls-key"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
