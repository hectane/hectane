package log

// Configuration for logging.
type Config struct {
	Debug   bool   `json:"debug"`
	Logfile string `json:"logfile"`
}
