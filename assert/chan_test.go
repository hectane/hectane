package assert

import (
	"testing"
)

func TestChanSend(t *testing.T) {
	c := make(chan interface{}, 1)
	if err := ChanSend(c, true); err != nil {
		t.Fatal("unexpected error")
	}
	if err := ChanSend(c, true); err == nil {
		t.Fatal("error expected")
	}
}

func TestChanRecv(t *testing.T) {
	c := make(chan interface{}, 1)
	c <- true
	if _, err := ChanRecv(c); err != nil {
		t.Fatal("unexpected error")
	}
	if _, err := ChanRecv(c); err == nil {
		t.Fatal("error expected")
	}
}

func TestChanRecvVal(t *testing.T) {
	c := make(chan interface{}, 2)
	c <- true
	c <- false
	if err := ChanRecvVal(c, true); err != nil {
		t.Fatal("unexpected error")
	}
	if err := ChanRecvVal(c, true); err == nil {
		t.Fatal("error expected")
	}
}

func TestChanClosed(t *testing.T) {
	c := make(chan interface{})
	if err := ChanClosed(c); err == nil {
		t.Fatal("error expected")
	}
	close(c)
	if err := ChanClosed(c); err != nil {
		t.Fatal("unexpected error")
	}
}
