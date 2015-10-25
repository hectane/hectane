package exec

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/queue"

	"encoding/json"
	"flag"
	"os"
)

// Global configuration for the application. This data is either read from the
// command-line or from a configuration file.
type Config struct {
	API   *api.Config   `json:"api"`
	Queue *queue.Config `json:"queue"`
}

// Initialize the global application configuration.
func InitConfig() (*Config, error) {
	var (
		c = &Config{
			API:   api.InitConfig(),
			Queue: queue.InitConfig(),
		}
		filename = flag.String("config", "", "file containing configuration")
	)
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
