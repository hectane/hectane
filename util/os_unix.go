// +build !windows

package util

import (
	"os"
)

// Ensure that the specified path is only accessible to the current user.
func SecurePath(path string) error {
	return os.Chmod(path, 0600)
}
