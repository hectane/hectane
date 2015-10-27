// +build !windows

package util

import (
	"os"
)

// Ensure that the specified path is only accessible to the current user. Note
// that files are directories require different permissions.
func SecurePath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return os.Chmod(path, 0700)
	} else {
		return os.Chmod(path, 0600)
	}
}
