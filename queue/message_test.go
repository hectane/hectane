package queue

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/nathan-osman/go-cannon/util"

	"os"
	"testing"
)

func TestMessage(t *testing.T) {
	var (
		directory = os.TempDir()
		host      = "example.com"
		from      = "me@example.com"
		to        = []string{"you@example.com"}
		id        = uuid.New()
	)
	if _, w, err := NewBody(directory, id); err == nil {
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		if m, err := NewMessage(directory, host, from, to, id); err == nil {
			if err := util.AssertFileState(m.metadataFilename(), true); err != nil {
				t.Fatal(err)
			}
			//...
			if err := m.Delete(); err != nil {
				t.Fatal(err)
			}
			if err := util.AssertFileState(m.metadataFilename(), false); err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
