package cfg

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/log"
	"github.com/hectane/hectane/queue"
	"github.com/hectane/hectane/smtp"

	"encoding/json"
	"flag"
	"os"
	"path"
)

// Global configuration for the application.
type Config struct {
	API   api.Config   `json:"api"`
	Log   log.Config   `json:"log"`
	Queue queue.Config `json:"queue"`
	SMTP  smtp.Config  `json:"smtp"`
}

// Parse the flags passed to the application
func Parse() (*Config, error) {
	var (
		c        = &Config{}
		filename = flag.String("config", "", "file containing configuration")
	)
	flag.StringVar(&c.API.Addr, "bind", ":8025", "`address` and port to bind to")
	flag.StringVar(&c.API.CORSOrigin, "cors-origin", "", "`origin` to use for CORS headers")
	flag.StringVar(&c.API.TLSCert, "tls-cert", "", "certificate `file` for TLS")
	flag.StringVar(&c.API.TLSKey, "tls-key", "", "private key `file` for TLS")
	flag.StringVar(&c.API.Username, "username", "", "`username` for HTTP basic auth")
	flag.StringVar(&c.API.Password, "password", "", "`password` for HTTP basic auth")
	flag.BoolVar(&c.Log.Debug, "debug", false, "show debug log messages")
	flag.StringVar(&c.Log.Logfile, "logfile", "", "`file` to write log output to")
	flag.StringVar(&c.Queue.Directory, "directory", path.Join(os.TempDir(), "hectane"), "`directory` for persistent storage")
	flag.BoolVar(&c.Queue.DisableSSLVerification, "disable-ssl-verification", false, "don't verify SSL certificates")
	flag.StringVar(&c.SMTP.Addr, "smtp-addr", ":smtp", "`address` and port for SMTP server")
	flag.IntVar(&c.SMTP.ReadTimeout, "read-timeout", 900, "`seconds` before client timeout")
	flag.Parse()
	if *filename != "" {
		r, err := os.Open(*filename)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		if err = json.NewDecoder(r).Decode(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Save the configuration to the specified path.
func (c *Config) Save(path string) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := json.NewEncoder(w).Encode(c); err != nil {
		return err
	}
	return nil
}
