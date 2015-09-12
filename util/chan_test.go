package util

import (
	"testing"
	"time"
)

func TestNonBlockingChan(t *testing.T) {
	n := NewNonBlockingChan()
	time.Sleep(50 * time.Millisecond)
	select {
	case n.Send <- true:
	default:
		t.Fatal("sending on channel blocked")
	}
}
