package queue

import (
	"flag"
	"os"
	"path"
)

// Application configuration.
type Config struct {
	Directory              string `json:"directory"`
	DisableSSLVerification bool   `json:"disable-ssl-verification"`
}

// Initialize the configuration.
func InitConfig() *Config {
	c := &Config{}
	flag.StringVar(&c.Directory, "directory", path.Join(os.TempDir(), "hectane"), "`directory` for persistent storage")
	flag.BoolVar(&c.DisableSSLVerification, "disable-ssl-verification", false, "don't verify SSL certificates")
	return c
}
