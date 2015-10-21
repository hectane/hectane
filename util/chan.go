package util

import (
	"container/list"
	"sync"
)

// Channel that does not block when items are sent. To use the type, simply
// send on the Send channel and receive on the Recv channel. Items will be
// stored internally until they are received. Closing the Send channel will
// cause the Recv channel to be also be closed after all items are received.
type NonBlockingChan struct {
	sync.Mutex
	Send  chan<- interface{}
	Recv  <-chan interface{}
	items *list.List
}

// Create a new non-blocking channel.
func NewNonBlockingChan() *NonBlockingChan {
	var (
		in  = make(chan interface{})
		out = make(chan interface{})
		n   = &NonBlockingChan{
			Send:  in,
			Recv:  out,
			items: list.New(),
		}
	)
	go func() {
		inClosed := false
		for {
			if inClosed && n.items.Len() == 0 {
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
			if n.items.Len() > 0 {
				outChan, outVal = out, n.items.Front().Value
			}
			select {
			case i, ok := <-inChan:
				if ok {
					n.Lock()
					n.items.PushBack(i)
					n.Unlock()
				} else {
					inClosed = true
				}
			case outChan <- outVal:
				n.Lock()
				n.items.Remove(n.items.Front())
				n.Unlock()
			}
		}
	}()
	return n
}

// Retrieve the number of items waiting to be received.
func (n *NonBlockingChan) Len() int {
	n.Lock()
	defer n.Unlock()
	return n.items.Len()
}
