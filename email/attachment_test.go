package email

import (
	"bytes"
	"io/ioutil"
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
	if mediatype, params, err := mime.ParseMediaType(part.Header.Get("Content-Type")); err == nil {
		if mediatype != contentType {
			t.Fatalf("%s != %s", mediatype, contentType)
		}
		if name, ok := params["name"]; ok {
			if name != filename {
				t.Fatalf("%s != %s", name, filename)
			}
		} else {
			t.Fatal("\"name\" parameter missing")
		}
	} else {
		t.Fatal(err)
	}
	if data, err := ioutil.ReadAll(part); err == nil {
		if string(data) != content {
			t.Fatalf("%s != %s", string(data), content)
		}
	} else {
		t.Fatal(err)
	}
}
