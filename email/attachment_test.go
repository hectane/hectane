package email

import (
	"github.com/hectane/go-attest"

	"bytes"
	"mime"
	"mime/multipart"
	"testing"
)

func TestAttachmentWrite(t *testing.T) {
	var (
		filename    = "test.txt"
		contentType = "text/plain"
		content     = "test"
		a           = &Attachment{
			Filename:    filename,
			ContentType: contentType,
			Content:     content,
		}
		buff = &bytes.Buffer{}
		w    = multipart.NewWriter(buff)
		r    = multipart.NewReader(buff, w.Boundary())
	)
	if err := a.Write(w); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	part, err := r.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	mediatype, params, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	if mediatype != contentType {
		t.Fatalf("%s != %s", mediatype, contentType)
	}
	name, ok := params["name"]
	if !ok {
		t.Fatal("\"name\" parameter missing")
	}
	if name != filename {
		t.Fatalf("%s != %s", name, filename)
	}
	if err := attest.Read(part, []byte(content)); err != nil {
		t.Fatal(err)
	}
}
