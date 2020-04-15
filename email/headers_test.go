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
	if buff.String() != value {
		t.Fatalf("%s != %s", buff.String(), value)
	}
}

func TestHeadersUTF8(t *testing.T) {
	var (
		headers = Headers{
			"Test": "日本語",
		}
		buff  = &bytes.Buffer{}
		value = "Test: =?utf-8?q?=E6=97=A5=E6=9C=AC=E8=AA=9E?=\r\n\r\n"
	)
	if err := headers.Write(buff); err != nil {
		t.Fatal(err)
	}
	if buff.String() != value {
		t.Fatalf("%s != %s", buff.String(), value)
	}
}
