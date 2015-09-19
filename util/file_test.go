package util

import (
	"code.google.com/p/go-uuid/uuid"

	"os"
	"path"
	"testing"
)

func newFilestore() (*Filestore, error) {
	directory := path.Join(os.TempDir(), uuid.New())
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, err
	}
	return NewFilestore(directory)
}

func TestFilestore(t *testing.T) {
	if f, err := newFilestore(); err == nil {
		_ = f
	} else {
		t.Fatal(err)
	}
}
