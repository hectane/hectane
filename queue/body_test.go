package queue

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/nathan-osman/go-cannon/util"

	"os"
	"testing"
)

// Ensure that the specified file is in the specified state of existence.
func expectExists(t *testing.T, filename string, exists bool) {
	if e, err := util.FileExists(filename); err == nil {
		if e != exists {
			if e {
				t.Fatal("file exists")
			} else {
				t.Fatal("file does not exist")
			}
		}
	} else {
		t.Fatal(err)
	}
}

func TestBodyRefCount(t *testing.T) {
	directory := os.TempDir()
	if b, writer, err := NewBody(directory, uuid.New()); err == nil {
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		expectExists(t, b.metadataFilename(), true)
		expectExists(t, b.messageBodyFilename(), true)
		if err := b.Add(); err != nil {
			t.Fatal(err)
		}
		if b.m.RefCount != 1 {
			t.Fatalf("%s != %s", b.m.RefCount, 1)
		}
		if err := b.Release(); err != nil {
			t.Fatal(err)
		}
		expectExists(t, b.metadataFilename(), false)
		expectExists(t, b.messageBodyFilename(), false)
	} else {
		t.Fatal(err)
	}
}
