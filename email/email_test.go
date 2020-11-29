package email

import (
	"github.com/hectane/go-attest"
	"github.com/hectane/hectane/queue"

	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"testing"
)

// Description of a multipart MIME message.
type multipartDesc struct {
	ContentType string
	Content     []byte
	Parts       []*multipartDesc
}

// Ensure that a multipart message conforms to the specified description.
func checkMultipart(r io.Reader, contentType string, d *multipartDesc) error {
	c, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	if c != d.ContentType {
		return fmt.Errorf("%s != %s", c, d.ContentType)
	}
	if len(d.Parts) == 0 {
		return attest.Read(r, d.Content)
	}
	boundary, ok := params["boundary"]
	if !ok {
		return errors.New("\"boundary\" parameter missing")
	}
	reader := multipart.NewReader(r, boundary)
	for _, part := range d.Parts {
		p, err := reader.NextPart()
		if err != nil {
			return err
		}
		if err := checkMultipart(p, p.Header.Get("Content-Type"), part); err != nil {
			return err
		}
	}
	return nil
}

func emailToMessages(e *Email) ([]*queue.Message, []byte, error) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(d)
	s := queue.NewStorage(d)
	m, err := e.Messages(s)
	if err != nil {
		return nil, nil, err
	}
	if len(m) < 1 {
		return nil, nil, errors.New("no messages")
	}
	r, err := s.GetMessageBody(m[0])
	if err != nil {
		return nil, nil, err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	return m, b, nil
}

func TestEmailErrors(t *testing.T) {
	badEmails := []*Email{
		{},
		{From: "me@example.com"},
		{To: []string{"you@example.com"}},
	}
	for _, e := range badEmails {
		_, _, err := emailToMessages(e)
		if err == nil {
			t.Fatal("error expected")
		}
	}
}

func TestEmailCount(t *testing.T) {
	m, _, err := emailToMessages(&Email{
		From: "me@example.com",
		To:   []string{"1@a.com", "1@b.com"},
		Cc:   []string{"2@a.com", "2@b.com"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatalf("%d != 2", len(m))
	}
}

func TestEmailHeaders(t *testing.T) {
	var (
		from    = "me@example.com"
		to      = "you@example.com"
		bcc     = "hidden@example.com"
		subject = "Test"
	)
	_, body, err := emailToMessages(&Email{
		From:    from,
		To:      []string{to},
		Bcc:     []string{bcc},
		Subject: subject,
	})
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewBuffer(body)
	m, err := mail.ReadMessage(r)
	if err != nil {
		t.Fatal(err)
	}
	if v := m.Header.Get("From"); v != from {
		t.Fatalf("%s != %s", v, from)
	}
	if v := m.Header.Get("To"); v != to {
		t.Fatalf("%s != %s", v, to)
	}
	if v := m.Header.Get("Bcc"); v != "" {
		t.Fatalf("%s != \"\"", v)
	}
	if v := m.Header.Get("Subject"); v != subject {
		t.Fatalf("%s != %s", v, subject)
	}
}

func TestEmailHeaders_importantHeadersAreNotClobbered(t *testing.T) {
	var (
		from    = "me@example.com"
		to      = "you@example.com"
		bcc     = "hidden@example.com"
		subject = "Test"
	)
	_, body, err := emailToMessages(&Email{
		From:    from,
		To:      []string{to},
		Bcc:     []string{bcc},
		Subject: subject,
		Headers: Headers{
			"From": "someone-else@example.com",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewBuffer(body)
	m, err := mail.ReadMessage(r)
	if err != nil {
		t.Fatal(err)
	}
	if v := m.Header.Get("From"); v != from {
		t.Fatalf("%s != %s", v, from)
	}
}

func TestEmailContent(t *testing.T) {
	var (
		from        = "me@example.com"
		to          = "you@example.com"
		text        = "<test>\r\n"
		description = &multipartDesc{
			ContentType: "multipart/mixed",
			Parts: []*multipartDesc{
				&multipartDesc{
					ContentType: "multipart/alternative",
					Parts: []*multipartDesc{
						&multipartDesc{
							ContentType: "text/plain",
							Content:     []byte("<test>\r\n"),
						},
						&multipartDesc{
							ContentType: "text/html",
							Content:     []byte("&lt;test&gt;<br>"),
						},
					},
				},
			},
		}
	)
	_, body, err := emailToMessages(&Email{
		From: from,
		Text: text,
		To:   []string{to},
	})
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewBuffer(body)
	m, err := mail.ReadMessage(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkMultipart(m.Body, m.Header.Get("Content-Type"), description); err != nil {
		t.Fatal(err)
	}
}
