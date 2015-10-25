// +build !windows

package exec

// No initialization is required for Unix-like platforms.
func Init(c *Config) error {
	return nil
}

// Only signals are available.
func Exec(c *Config) {
	execSignal()
}
