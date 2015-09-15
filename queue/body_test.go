package queue

import (
	"code.google.com/p/go-uuid/uuid"

	"os"
	"testing"
)

func TestBodyRefCount(t *testing.T) {
	directory := os.TempDir()
	if b, writer, err := NewBody(directory, uuid.New()); err == nil {
		writer.Close()
		if err := b.Add(); err != nil {
			t.Fatal(err)
		}
		if b.m.RefCount != 1 {
			t.Fatalf("%s != %s", b.m.RefCount, 1)
		}
		if err := b.Release(); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(b.metadataFilename()); !os.IsNotExist(err) {
			t.Fatal("metadata file not removed")
		}
		if _, err := os.Stat(b.messageBodyFilename()); !os.IsNotExist(err) {
			t.Fatal("message body file not removed")
		}
	} else {
		t.Fatal(err)
	}
}
