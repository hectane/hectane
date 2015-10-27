package util

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSecurePath(t *testing.T) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)
	if err = SecurePath(d); err != nil {
		t.Fatal(err)
	}
}

func TestExecutable(t *testing.T) {
	p, err := Executable()
	if err != nil {
		t.Fatal(err)
	}
	if len(p) == 0 {
		t.Fatalf("empty string")
	}
}
