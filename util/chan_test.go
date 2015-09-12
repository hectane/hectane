package util

import (
	"testing"
	"time"
)

// Ensure that a value is immediately sent on the specified channel.
func expectSend(sendChan chan<- interface{}) bool {
	select {
	case sendChan <- true:
		return true
	case <-time.After(50 * time.Millisecond):
		return false
	}
}

// Ensure that a value is immediately received on the specified channel. The
// return value indicates whether the channel is open and whether a value was
// received.
func expectRecv(recvChan <-chan interface{}) (bool, bool) {
	select {
	case _, ok := <-recvChan:
		return ok, ok
	case <-time.After(50 * time.Millisecond):
		return true, false
	}
}

func TestNonBlockingChan(t *testing.T) {
	n := NewNonBlockingChan()

	// Send and receive a value on the channel
	if !expectSend(n.Send) {
		t.Fatal("sending on channel blocked")
	}
	if open, recvd := expectRecv(n.Recv); !open || !recvd {
		t.Fatal("unable to receive item")
	}

	// Close the sending channel and ensure the receiving channel is closed
	close(n.Send)
	if open, recvd := expectRecv(n.Recv); open || recvd {
		t.Fatal("receiving channel not closed")
	}
}
