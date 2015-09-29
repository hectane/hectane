package assert

import (
	"bytes"
	"mime"
	"mime/multipart"
	"net/textproto"
	"testing"
)

func TestMultipart(t *testing.T) {
	var (
		b           = &bytes.Buffer{}
		w           = multipart.NewWriter(b)
		contentType = "multipart/mixed"
		header      = mime.FormatMediaType(contentType, map[string]string{
			"boundary": w.Boundary(),
		})
		data            = []byte("test")
		dataContentType = "text/plain"
		badData         = []byte("test2")
	)
	if p, err := w.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{dataContentType},
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
	if err := Multipart(b, header, &MultipartDesc{
		ContentType: contentType,
		Parts: []*MultipartDesc{
			&MultipartDesc{
				ContentType: dataContentType,
				Content:     data,
			},
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := Multipart(b, header, &MultipartDesc{
		ContentType: contentType,
		Parts: []*MultipartDesc{
			&MultipartDesc{
				ContentType: dataContentType,
				Content:     badData,
			},
		},
	}); err == nil {
		t.Fatal("error expected")
	}
}
