// +build !windows

package util

// Ensure that the specified path is only accessible to the current user. On
// Unix platforms, this requires nothing since the correct permissions are set
// when the path is created.
func SecurePath(path string) error {
	return nil
}
