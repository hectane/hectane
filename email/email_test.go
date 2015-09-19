package email

import (
	"os"
	"testing"
)

var (
	directory = os.TempDir()
	addr1A    = "1@a.com"
	addr2A    = "2@a.com"
	addr1B    = "1@b.com"
	subject   = "Test"
	text      = "test"
	html      = "<em>test</em>"
)

func TestEmailMessageCount(t *testing.T) {
	e := &Email{
		To: []string{addr1A, addr2A},
		Cc: []string{addr1B},
	}
	if messages, err := e.Messages(directory); err == nil {
		if len(messages) != 2 {
			t.Fatalf("%d != 2", len(messages))
		}
	} else {
		t.Fatal(err)
	}
}
