package util

import (
	"code.google.com/p/go-uuid/uuid"

	"os"
	"path"
	"testing"
)

func TestFileExists(t *testing.T) {
	var (
		directory = os.TempDir()
		filename  = path.Join(directory, uuid.New())
	)
	if f, err := os.Create(filename); err == nil {
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	if exists, err := FileExists(filename); err != nil {
		t.Fatal(err)
	} else {
		if !exists {
			t.Fatal("FileExists() == false")
		}
	}
	if err := os.Remove(filename); err != nil {
		t.Fatal(err)
	}
	if exists, err := FileExists(filename); err != nil {
		t.Fatal(err)
	} else {
		if exists {
			t.Fatal("FileExists() == true")
		}
	}
}
