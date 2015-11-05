package email

import (
	"github.com/hectane/hectane/assert"
	"github.com/hectane/hectane/queue"

	"bytes"
	"errors"
	"io/ioutil"
	"net/mail"
	"os"
	"testing"
)

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
	if b, err := ioutil.ReadAll(r); err == nil {
		return m, b, nil
	} else {
		return nil, nil, err
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

func TestEmailContent(t *testing.T) {
	var (
		from        = "me@example.com"
		to          = "you@example.com"
		text        = "<test>\r\n"
		description = &assert.MultipartDesc{
			ContentType: "multipart/mixed",
			Parts: []*assert.MultipartDesc{
				&assert.MultipartDesc{
					ContentType: "multipart/alternative",
					Parts: []*assert.MultipartDesc{
						&assert.MultipartDesc{
							ContentType: "text/plain",
							Content:     []byte("<test>\r\n"),
						},
						&assert.MultipartDesc{
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
	if err := assert.Multipart(m.Body, m.Header.Get("Content-Type"), description); err != nil {
		t.Fatal(err)
	}
}
