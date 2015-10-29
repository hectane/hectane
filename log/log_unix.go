// +build !windows

package log

// No initialization required on Unix platforms.
func platformInit(config *Config) error {
	return nil
}

// No cleanup required on Unix platforms.
func platformCleanup() {
}
