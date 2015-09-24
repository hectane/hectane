package util

import (
	"testing"
	"time"
)

func TestNonBlockingChan(t *testing.T) {
	n := NewNonBlockingChan()
	<-time.After(50 * time.Millisecond)
	if err := AssertChanSend(n.Send, true); err != nil {
		t.Fatal(err)
	}
	<-time.After(50 * time.Millisecond)
	if err := AssertChanRecvVal(n.Recv, true); err != nil {
		t.Fatal(err)
	}
	close(n.Send)
	<-time.After(50 * time.Millisecond)
	if err := AssertChanClosed(n.Recv); err != nil {
		t.Fatal(err)
	}
}
