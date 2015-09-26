package email

import (
	"github.com/nathan-osman/go-cannon/queue"
	"github.com/nathan-osman/go-cannon/util"

	"testing"
)

var (
	addr1A  = "1@a.com"
	addr2A  = "2@a.com"
	addr1B  = "1@b.com"
	subject = "Test"
	text    = "test"
	html    = "<em>test</em>"
)

func emailToMessages(e *Email) ([]*queue.Message, error) {
	if d, err := util.NewTempDir(); err == nil {
		defer d.Delete()
		if s, err := queue.NewStorage(d.Path); err == nil {
			if messages, err := e.Messages(s); err == nil {
				return messages, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func TestEmailCount(t *testing.T) {
	if messages, err := emailToMessages(&Email{
		To:      []string{addr1A, addr1B},
		Cc:      []string{addr2A},
		Subject: subject,
		Text:    text,
		Html:    html,
	}); err == nil {
		if len(messages) != 2 {
			t.Fatal("%d != 2", len(messages))
		}
	} else {
		t.Fatal(err)
	}
}
