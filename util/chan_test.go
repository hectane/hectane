package util

import (
	"github.com/hectane/hectane/assert"

	"testing"
	"time"
)

func TestNonBlockingChan(t *testing.T) {
	n := NewNonBlockingChan()
	<-time.After(50 * time.Millisecond)
	if err := assert.ChanSend(n.Send, true); err != nil {
		t.Fatal(err)
	}
	<-time.After(50 * time.Millisecond)
	if n.Len() != 1 {
		t.Fatalf("%d != %d", n.Len(), 1)
	}
	if err := assert.ChanRecvVal(n.Recv, true); err != nil {
		t.Fatal(err)
	}
	<-time.After(50 * time.Millisecond)
	if n.Len() != 0 {
		t.Fatalf("%d != %d", n.Len(), 0)
	}
	close(n.Send)
	<-time.After(50 * time.Millisecond)
	if err := assert.ChanClosed(n.Recv); err != nil {
		t.Fatal(err)
	}
}
