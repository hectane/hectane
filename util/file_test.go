package util

import (
	"os"
	"testing"
)

func TestTempDir(t *testing.T) {
	if d, err := NewTempDir(); err == nil {
		if err := AssertFileState(d.Path, true); err != nil {
			t.Fatal(err)
		}
		if err := os.RemoveAll(d.Path); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
