package queue

import (
	"github.com/nathan-osman/go-cannon/util"

	"testing"
)

func TestStorage(t *testing.T) {
	var (
		data        = []byte("test")
		numMessages = 5
		messageIDs  = make([]string, numMessages)
		m           = &Message{}
	)
	if d, err := util.NewTempDir(); err == nil {
		defer d.Delete()
		if s, err := NewStorage(d.Path); err == nil {
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
			if messages, err := s.LoadMessages(); err == nil {
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
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
