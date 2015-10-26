package util

import (
	"testing"
)

func TestExecutable(t *testing.T) {
	p, err := Executable()
	if err != nil {
		t.Fatal(err)
	}
	if len(p) == 0 {
		t.Fatalf("empty string")
	}
}
