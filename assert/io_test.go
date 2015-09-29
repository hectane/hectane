package assert

import (
	"bytes"
	"testing"
)

func TestRead(t *testing.T) {
	var (
		data    = []byte("test")
		badData = []byte("test2")
	)
	if err := Read(bytes.NewBuffer(data), data); err != nil {
		t.Fatal(err)
	}
	if err := Read(bytes.NewBuffer(data), badData); err == nil {
		t.Fatal("error expected")
	}
}
