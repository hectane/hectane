package main

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/queue"

	"encoding/json"
	"flag"
	"os"
	"path"
)

// Global configuration for the application. This data is either read from the
// command-line or from a configuration file.
type Config struct {
	API   api.Config   `json:"api"`
	Queue queue.Config `json:"queue"`
}

// Register command-line flags for each of the options.
func (c *Config) RegisterFlags() {
	flag.StringVar(&c.API.Addr, "bind", ":8025", "address and port to bind to")
	flag.StringVar(&c.API.TLSCert, "tls-cert", "", "certificate for TLS")
	flag.StringVar(&c.API.TLSKey, "tls-key", "", "private key for TLS")
	flag.StringVar(&c.API.Username, "username", "", "username for HTTP basic auth")
	flag.StringVar(&c.API.Password, "password", "", "password for HTTP basic auth")
	flag.StringVar(&c.Queue.Directory, "directory", path.Join(os.TempDir(), "hectane"), "directory for persistent storage")
	flag.BoolVar(&c.Queue.DisableSSLVerification, "disable-ssl-verification", false, "don't verify SSL certificates")
}

// Load the configuration from the specified file.
func (c *Config) LoadFromFile(filename string) error {
	if r, err := os.Open(filename); err == nil {
		defer r.Close()
		return json.NewDecoder(r).Decode(c)
	} else {
		return err
	}
}
