package email

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestWrite(t *testing.T) {
	var (
		filename        = "test.txt"
		contentType     = "text/plain"
		contentTypeLine = fmt.Sprintf("%s; name=%s", contentType, filename)
		content         = "test"
		a               = &Attachment{
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
	if part, err := r.NextPart(); err == nil {
		hdr := part.Header.Get("Content-Type")
		if hdr != contentTypeLine {
			t.Fatalf("%s != %s", hdr, contentTypeLine)
		}
		if data, err := ioutil.ReadAll(part); err == nil {
			if !reflect.DeepEqual(data, []byte(content)) {
				t.Fatalf("%s != %s", data, []byte(content))
			}
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
