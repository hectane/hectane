package assert

import (
	"bytes"
	"testing"
)

func TestRead(t *testing.T) {
	var (
		data = []byte("test")
		b    = bytes.NewBuffer(data)
	)
	if err := Read(b, data); err != nil {
		t.Fatal(err)
	}
}
