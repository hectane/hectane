package util

import (
	"code.google.com/p/go-uuid/uuid"

	"os"
	"path"
)

// Temporary directory that is guaranteed to be empty.
type TempDir struct {
	Path string
}

// Create a new temporary directory.
func NewTempDir() (*TempDir, error) {
	t := &TempDir{
		Path: path.Join(os.TempDir(), uuid.New()),
	}
	if err := os.MkdirAll(t.Path, 0700); err == nil {
		return t, nil
	} else {
		return nil, err
	}
}

// Delete the temporary directory and its contents.
func (t *TempDir) Delete() error {
	return os.RemoveAll(t.Path)
}
