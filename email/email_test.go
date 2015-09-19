package email

import (
	"os"
	"testing"
)

var (
	directory = os.TempDir()
	from      = "a@example.com"
	to        = "b@example.com"
	cc        = "c@example.com"
	bcc       = "d@example.org"
	subject   = "Test"
	text      = "test"
	html      = "<em>test</em>"
	e         = &Email{
		From:    from,
		To:      []string{to},
		Cc:      []string{cc},
		Bcc:     []string{bcc},
		Subject: subject,
		Text:    text,
		Html:    html,
	}
)

func TestEmailMessageCount(t *testing.T) {
	if messages, err := e.Messages(directory); err == nil {
		if len(messages) != 2 {
			t.Fatalf("%d != 2", len(messages))
		}
	} else {
		t.Fatal(err)
	}
}
