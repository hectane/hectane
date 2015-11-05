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
	d, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)
	s := NewStorage(d)
	w, body, err := s.NewBody()
	if err != nil {
		t.Fatal(err)
	}
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
	messages, err := s.LoadMessages()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range messages {
		r, err := s.GetMessageBody(m)
		if err != nil {
			t.Fatal(err)
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(b, data) {
			t.Fatalf("%v != %v", b, data)
		}
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
		if err := s.DeleteMessage(m); err != nil {
			t.Fatal(err)
		}
	}
	e, err := ioutil.ReadDir(s.directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(e) != 0 {
		t.Fatalf("%d != 0", len(e))
	}
}
