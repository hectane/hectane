package email

import (
	"os"
	"testing"
)

func TestEmailMessages(t *testing.T) {
	var (
		directory = os.TempDir()
		from      = "a@hotmail.com"
		to        = "b@hotmail.com"
		cc        = "c@hotmail.com"
		bcc       = "d@yahoo.com"
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
	if messages, err := e.Messages(directory); err == nil {
		hosts := make([]string, 0, 2)
		for _, m := range messages {
			_ = hosts
			_ = m
			//...
		}
	} else {
		t.Fatal(err)
	}
}
