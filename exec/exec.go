package exec

import (
	"github.com/Sirupsen/logrus"
	"github.com/hectane/hectane/cfg"

	"os"
)

// Setup log redirection (if requested) and initialize the execution
// environment for the current platform.
func Init(config *cfg.Config) error {
	if config.Log != "" {
		f, err := os.OpenFile(config.Log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		logrus.SetOutput(f)
	}
	return platformInit()
}
