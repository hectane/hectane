package util

import (
	"container/list"
)

// Channel that does not block when items are sent. To use the type, simply
// send on the Send channel and receive on the Recv channel. Items will be
// stored internally until they are received. Closing the Send channel will
// cause the Recv channel to be also be closed after all items are received.
type NonBlockingChan struct {
	Send chan<- interface{}
	Recv <-chan interface{}
}

// Create a new non-blocking channel.
func NewNonBlockingChan() *NonBlockingChan {
	var (
		in  = make(chan interface{})
		out = make(chan interface{})
		n   = &NonBlockingChan{
			Send: in,
			Recv: out,
		}
	)
	go func() {
		var (
			items    = list.New()
			inClosed = false
		)
		for {
			if inClosed && items.Len() == 0 {
				close(out)
				break
			}
			var (
				inChan  chan interface{}
				outChan chan interface{}
				outVal  interface{}
			)
			if !inClosed {
				inChan = in
			}
			if items.Len() > 0 {
				outChan, outVal = out, items.Front().Value
			}
			select {
			case i, ok := <-inChan:
				if ok {
					items.PushBack(i)
				} else {
					inClosed = true
				}
			case outChan <- outVal:
				items.Remove(items.Front())
			}
		}
	}()
	return n
}
