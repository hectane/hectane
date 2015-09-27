package email

import (
	"github.com/nathan-osman/go-cannon/queue"

	"bytes"
	"io/ioutil"
	"net/mail"
	"os"
	"testing"
)

func emailToMessages(e *Email) ([]*queue.Message, []byte, error) {
	if d, err := ioutil.TempDir(os.TempDir(), ""); err == nil {
		defer os.RemoveAll(d)
		s := queue.NewStorage(d)
		if messages, err := e.Messages(s); err == nil {
			if len(messages) > 0 {
				if r, err := s.GetMessageBody(messages[0]); err == nil {
					if b, err := ioutil.ReadAll(r); err == nil {
						return messages, b, nil
					} else {
						return nil, nil, err
					}
				} else {
					return nil, nil, err
				}
			} else {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	} else {
		return nil, nil, err
	}
}

func TestEmailCount(t *testing.T) {
	if messages, _, err := emailToMessages(&Email{
		To: []string{"1@a.com", "1@b.com"},
		Cc: []string{"2@a.com", "2@b.com"},
	}); err == nil {
		if len(messages) != 2 {
			t.Fatal("%d != 2", len(messages))
		}
	} else {
		t.Fatal(err)
	}
}

func TestEmailHeaders(t *testing.T) {
	var (
		to      = "test@example.com"
		bcc     = "hidden@example.com"
		subject = "Test"
	)
	if _, body, err := emailToMessages(&Email{
		To:      []string{to},
		Bcc:     []string{bcc},
		Subject: subject,
	}); err == nil {
		r := bytes.NewBuffer(body)
		if m, err := mail.ReadMessage(r); err == nil {
			if v := m.Header.Get("To"); v != to {
				t.Fatalf("%s != %s", v, to)
			}
			if v := m.Header.Get("Bcc"); v != "" {
				t.Fatalf("%s != \"\"", v)
			}
			if v := m.Header.Get("Subject"); v != subject {
				t.Fatalf("%s != %s", v, subject)
			}
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
