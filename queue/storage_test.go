package queue

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestStorage(t *testing.T) {
	var (
		data        = []byte("test")
		numMessages = 5
	)
	if d, err := ioutil.TempDir(os.TempDir(), ""); err == nil {
		defer os.RemoveAll(d)
		s := NewStorage(d)
		if w, body, err := s.NewBody(); err == nil {
			if _, err := w.Write(data); err != nil {
				t.Fatal(err)
			}
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
			for i := 0; i < numMessages; i++ {
				if err := s.SaveMessage(&Message{}, body); err != nil {
					t.Fatal(err)
				}
			}
		} else {
			t.Fatal(err)
		}
		if messages, err := s.LoadMessages(); err == nil {
			for _, m := range messages {
				if r, err := s.GetMessageBody(m); err == nil {
					if b, err := ioutil.ReadAll(r); err == nil {
						if !reflect.DeepEqual(b, data) {
							t.Fatalf("%v != %v", b, data)
						}
					} else {
						t.Fatal(err)
					}
					if err := r.Close(); err != nil {
						t.Fatal(err)
					}
				} else {
					t.Fatal(err)
				}
				if err := s.DeleteMessage(m); err != nil {
					t.Fatal(err)
				}
			}
		} else {
			t.Fatal(err)
		}
		if e, err := ioutil.ReadDir(s.directory); err == nil {
			if len(e) != 0 {
				t.Fatalf("%d != 0", len(e))
			}
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
