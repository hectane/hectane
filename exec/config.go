package exec

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/queue"

	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

// Global configuration for the application. This data is either read from the
// command-line or from a configuration file.
type Config struct {
	API   *api.Config   `json:"api"`
	Queue *queue.Config `json:"queue"`
}

// Display usage information for the application.
func usage() {
	fmt.Fprintf(os.Stderr, "USAGE\n\thectane [options] [command]\n\n")
	fmt.Fprintf(os.Stderr, "COMMANDS\n")
	for _, c := range Commands {
		fmt.Fprintf(os.Stderr, "  %s\n\t%s\n", c.Name, c.Description)
	}
	fmt.Fprintf(os.Stderr, "\nOPTIONS\n")
	flag.PrintDefaults()
}

// Initialize the global application configuration and parse the command line.
func InitConfig() (string, *Config, error) {
	var (
		cmd string
		c   = &Config{
			API:   api.InitConfig(),
			Queue: queue.InitConfig(),
		}
		filename = flag.String("config", "", "file containing configuration")
	)
	flag.Usage = usage
	flag.Parse()
	switch {
	case flag.NArg() == 0:
		cmd = "run"
	case flag.NArg() == 1:
		cmd = flag.Args()[0]
	default:
		return "", nil, errors.New("single command expected")
	}
	if *filename != "" {
		r, err := os.Open(*filename)
		if err != nil {
			return "", nil, err
		}
		defer r.Close()
		if err = json.NewDecoder(r).Decode(c); err != nil {
			return "", nil, err
		}
	}
	return cmd, c, nil
}
