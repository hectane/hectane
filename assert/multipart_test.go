package assert

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"testing"
)

func TestMultipart(t *testing.T) {
	var (
		b           = &bytes.Buffer{}
		w           = multipart.NewWriter(b)
		data        = []byte("test")
		description = &MultipartDesc{
			Parts: map[string]*MultipartDesc{
				"text/plain": &MultipartDesc{Content: data},
			},
		}
	)
	if p, err := w.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain"},
	}); err == nil {
		if _, err := p.Write(data); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := Multipart(b, w.Boundary(), description); err != nil {
		t.Fatal(err)
	}
}
