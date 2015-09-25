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
		path.Join(os.TempDir(), uuid.New()),
	}
	if err := os.MkdirAll(t.Path, 0700); err == nil {
		return t, nil
	} else {
		return nil, err
	}
}
