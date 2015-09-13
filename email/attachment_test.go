package email

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"reflect"
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
	if part, err := r.NextPart(); err == nil {
		var (
			header          = part.Header.Get("Content-Type")
			contentTypeLine = fmt.Sprintf("%s; name=%s", contentType, filename)
		)
		if header != contentTypeLine {
			t.Fatalf("%s != %s", header, contentTypeLine)
		}
		if data, err := ioutil.ReadAll(part); err == nil {
			if !reflect.DeepEqual(string(data), content) {
				t.Fatalf("%s != %s", string(data), content)
			}
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
