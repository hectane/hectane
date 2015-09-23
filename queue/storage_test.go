package queue

import (
	"code.google.com/p/go-uuid/uuid"

	"os"
	"path"
	"testing"
)

func TestStorage(t *testing.T) {
	var (
		directory   = path.Join(os.TempDir(), uuid.New())
		data        = []byte("test")
		numMessages = 5
		messageIDs  = make([]string, numMessages)
		m           = &Message{}
	)
	if err := os.MkdirAll(directory, 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(directory)
	if s, messages, err := NewStorage(directory); err == nil {
		if len(messages) != 0 {
			t.Fatalf("%d != 0", len(messages))
		}
		if w, id, err := s.NewBody(); err == nil {
			if _, err := w.Write(data); err != nil {
				t.Fatal(err)
			}
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
			m.Body = id
		} else {
			t.Fatal(err)
		}
		for i := 0; i < numMessages; i++ {
			if id, err := s.NewMessage(m); err == nil {
				messageIDs[i] = id
			} else {
				t.Fatal(err)
			}
		}
	} else {
		t.Fatal(err)
	}
	if s, messages, err := NewStorage(directory); err == nil {
		if len(messages) != numMessages {
			t.Fatalf("%d != %d", len(messages), numMessages)
		}
		for _, id := range messageIDs {
			if err := s.DeleteMessage(id); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := s.GetBody(m.Body); err != InvalidID {
			t.Fatalf("%v != %v", err, InvalidID)
		}
	} else {
		t.Fatal(err)
	}
}
