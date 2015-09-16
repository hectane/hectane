package queue

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/nathan-osman/go-cannon/util"

	"io/ioutil"
	"os"
	"reflect"
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
		if err := b.Add(); err != nil {
			t.Fatal(err)
		}
		if b.m.RefCount != 1 {
			t.Fatalf("%d != 1", b.m.RefCount)
		}
		expectExists(t, b.metadataFilename(), true)
		expectExists(t, b.messageBodyFilename(), true)
		if err := b.Release(); err != nil {
			t.Fatal(err)
		}
		if b.m.RefCount != 0 {
			t.Fatalf("%d != 0", b.m.RefCount)
		}
		expectExists(t, b.metadataFilename(), false)
		expectExists(t, b.messageBodyFilename(), false)
	} else {
		t.Fatal(err)
	}
}

func TestLoadBody(t *testing.T) {
	var (
		directory = os.TempDir()
		id        = uuid.New()
		data      = []byte("test")
	)
	if b, writer, err := NewBody(directory, id); err == nil {
		if _, err := writer.Write(data); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		if err := b.Add(); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	if b, err := LoadBody(directory, id); err == nil {
		if r, err := b.Reader(); err == nil {
			if d, err := ioutil.ReadAll(r); err == nil {
				if !reflect.DeepEqual(d, data) {
					t.Fatal("%v != %v", d, data)
				}
			} else {
				t.Fatal(err)
			}
		} else {
			t.Fatal(err)
		}
		if err := b.Release(); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
