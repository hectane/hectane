package log

import (
	"github.com/Sirupsen/logrus"

	"os"
)

// Disable colored output since it doesn't work universally and leads to
// problems when writing to files.
func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
}

// If the user has supplied a filename for logging, open the file and redirect
// all log output there.
func Init(config *Config) error {
	if config.Logfile != "" {
		f, err := os.OpenFile(config.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		logrus.SetOutput(f)
	}
	return nil
}
