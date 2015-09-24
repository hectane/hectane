package util

import (
	"code.google.com/p/go-uuid/uuid"

	"os"
	"path"
	"testing"
)

func TestAssertChanSend(t *testing.T) {
	c := make(chan interface{}, 1)
	if err := AssertChanSend(c, true); err != nil {
		t.Fatal(err)
	}
	if err := AssertChanSend(c, true); err == nil {
		t.Fatal("error expected")
	}
}

func TestAssertChanRecv(t *testing.T) {
	c := make(chan interface{}, 1)
	c <- true
	if _, err := AssertChanRecv(c); err != nil {
		t.Fatal(err)
	}
	if _, err := AssertChanRecv(c); err == nil {
		t.Fatal("error expected")
	}
}

func TestAssertChanRecvVal(t *testing.T) {
	c := make(chan interface{}, 2)
	c <- true
	c <- false
	if err := AssertChanRecvVal(c, true); err != nil {
		t.Fatal(err)
	}
	if err := AssertChanRecvVal(c, true); err == nil {
		t.Fatal("error expected")
	}
}

func TestAssertChanClosed(t *testing.T) {
	c := make(chan interface{})
	if err := AssertChanClosed(c); err == nil {
		t.Fatal("error expected")
	}
	close(c)
	if err := AssertChanClosed(c); err != nil {
		t.Fatal(err)
	}
}

func TestAssertFileState(t *testing.T) {
	var (
		directory = os.TempDir()
		filename  = path.Join(directory, uuid.New())
	)
	if f, err := os.Create(filename); err == nil {
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	if err := AssertFileState(filename, true); err != nil {
		t.Fatal("unexpected error")
	}
	if err := AssertFileState(filename, false); err == nil {
		t.Fatal("error expected")
	}
	if err := os.Remove(filename); err != nil {
		t.Fatal(err)
	}
	if err := AssertFileState(filename, true); err == nil {
		t.Fatal("error expected")
	}
	if err := AssertFileState(filename, false); err != nil {
		t.Fatal("unexpected error")
	}
}
