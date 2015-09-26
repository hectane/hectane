package email

import (
	"github.com/nathan-osman/go-cannon/queue"
	"github.com/nathan-osman/go-cannon/util"

	"testing"
)

func emailToMessages(e *Email) (*queue.Storage, []*queue.Message, error) {
	if d, err := util.NewTempDir(); err == nil {
		defer d.Delete()
		if s, err := queue.NewStorage(d.Path); err == nil {
			if messages, err := e.Messages(s); err == nil {
				return s, messages, nil
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
	if _, messages, err := emailToMessages(&Email{
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
