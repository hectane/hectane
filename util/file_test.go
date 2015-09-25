package util

import (
	"testing"
)

func TestTempDir(t *testing.T) {
	if d, err := NewTempDir(); err == nil {
		if err := AssertFileState(d.Path, true); err != nil {
			t.Fatal(err)
		}
		if err := d.Delete(); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
