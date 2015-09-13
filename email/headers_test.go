package email

import (
	"bytes"
	"testing"
)

func TestEmailHeadersWrite(t *testing.T) {
	var (
		headers = EmailHeaders{
			"Test": "test",
		}
		buff  = &bytes.Buffer{}
		value = "Test: test\r\n\r\n"
	)
	if err := headers.Write(buff); err != nil {
		t.Fatal(err)
	}
	if string(buff.Bytes()) != value {
		t.Fatalf("%s != %s", string(buff.Bytes()), value)
	}
}
