package log

import (
	"github.com/Sirupsen/logrus"

	"os"
)

// Initialize the logging backend. Colored output is disabled since it isn't
// supported on all platforms and because it will cause problems when
// redirecting log output to a file.
func Init(config *Config) error {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
	if config.Logfile != "" {
		f, err := os.OpenFile(config.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		logrus.SetOutput(f)
	}
	return platformInit(config)
}

// Shutdown the logging backend.
func Cleanup() {
	platformCleanup()
}
