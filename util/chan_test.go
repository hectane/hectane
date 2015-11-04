package util

import (
	"github.com/hectane/go-attest"

	"testing"
	"time"
)

func TestNonBlockingChan(t *testing.T) {
	n := NewNonBlockingChan()
	<-time.After(50 * time.Millisecond)
	if err := attest.ChanSend(n.Send, true); err != nil {
		t.Fatal(err)
	}
	<-time.After(50 * time.Millisecond)
	if n.Len() != 1 {
		t.Fatalf("%d != %d", n.Len(), 1)
	}
	v, err := attest.ChanRecv(n.Recv)
	if err != nil {
		t.Fatal(err)
	}
	if v != true {
		t.Fatalf("%v != true", v)
	}
	<-time.After(50 * time.Millisecond)
	if n.Len() != 0 {
		t.Fatalf("%d != %d", n.Len(), 0)
	}
	close(n.Send)
	<-time.After(50 * time.Millisecond)
	if err := attest.ChanClosed(n.Recv); err != nil {
		t.Fatal(err)
	}
}
