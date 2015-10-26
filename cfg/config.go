package cfg

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/queue"

	"encoding/json"
	"flag"
	"os"
	"path"
)

// Global configuration for the application.
type Config struct {
	API   api.Config   `json:"api"`
	Queue queue.Config `json:"queue"`
}

// Parse the flags passed to the application
func Parse() (*Config, error) {
	var (
		cfg      = &Config{}
		filename = flag.String("config", "", "file containing configuration")
	)
	flag.StringVar(&cfg.API.Addr, "bind", ":8025", "`address` and port to bind to")
	flag.StringVar(&cfg.API.TLSCert, "tls-cert", "", "certificate `file` for TLS")
	flag.StringVar(&cfg.API.TLSKey, "tls-key", "", "private key `file` for TLS")
	flag.StringVar(&cfg.API.Username, "username", "", "`username` for HTTP basic auth")
	flag.StringVar(&cfg.API.Password, "password", "", "`password` for HTTP basic auth")
	flag.StringVar(&cfg.Queue.Directory, "directory", path.Join(os.TempDir(), "hectane"), "`directory` for persistent storage")
	flag.BoolVar(&cfg.Queue.DisableSSLVerification, "disable-ssl-verification", false, "don't verify SSL certificates")
	flag.Parse()
	if *filename != "" {
		r, err := os.Open(*filename)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		if err = json.NewDecoder(r).Decode(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
