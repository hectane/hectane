package cmd

import (
	"testing"
)

func TestBadCommand(t *testing.T) {
	if err := Exec("", nil); err == nil {
		t.Fatal("error expected")
	}
}
