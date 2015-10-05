package email

import (
	"bytes"
	"testing"
)

func TestHeadersWrite(t *testing.T) {
	var (
		headers = Headers{
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
