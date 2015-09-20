package util

import (
	"code.google.com/p/go-uuid/uuid"

	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestFilestore(t *testing.T) {
	var (
		directory = path.Join(os.TempDir(), uuid.New())
		id        = "a"
		data      = []byte("test")
		refCount  = 5
	)
	if err := os.MkdirAll(directory, 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(directory)
	if f, err := NewFilestore(directory); err == nil {
		if w, err := f.New(id); err == nil {
			if _, err := w.Write(data); err != nil {
				t.Fatal(err)
			}
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatal(err)
		}
		for i := 0; i < refCount; i++ {
			if err := f.Add(id); err != nil {
				t.Fatal(err)
			}
		}
		if r, err := f.Get(id); err == nil {
			if d, err := ioutil.ReadAll(r); err == nil {
				if !reflect.DeepEqual(d, data) {
					t.Fatalf("%v != %v", d, data)
				}
			} else {
				t.Fatal(err)
			}
		} else {
			t.Fatal(err)
		}
		for i := 0; i < refCount; i++ {
			if err := f.Release(id); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := f.Get(id); err != InvalidItem {
			t.Fatalf("%v != %v", err, InvalidItem)
		}
	} else {
		t.Fatal(err)
	}
}
