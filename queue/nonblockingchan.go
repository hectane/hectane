package queue

import (
	"container/list"
	"reflect"
)

// Channel that does not block when items are sent. To use the struct, simply
// send on the Send channel and receive on the Recv channel. Items will be
// stored internally until they are received. Closing the Send channel will
// cause the Recv channel to be also be closed once all items are received.
type NonBlockingChan struct {
	Send chan<- interface{}
	Recv <-chan interface{}
}

// Create a new non-blocking channel.
func NewNonBlockingChan() *NonBlockingChan {

	// Create the two channels that will be used for sending and receiving
	var (
		send = make(chan interface{})
		recv = make(chan interface{})
	)

	// Assign the channels to the public members of the struct (which limits
	// each of their direction)
	nbChan := &NonBlockingChan{
		Send: send,
		Recv: recv,
	}

	// Start a goroutine to perform the sending and receiving
	go func() {

		// Create the list that will temporarily hold items for receiving and
		// create a boolean to track whether the Send channel has been closed
		var (
			items      = list.New()
			sendClosed = false
		)

		// Create constants for the select case array
		const (
			incomingCase = iota
			outgoingCase
			numCases
		)

		for {

			// Close the Recv channel and quit if the Send channel was closed
			// and there are no more items left in the list to send
			if sendClosed && items.Len() == 0 {
				close(recv)
				break
			}

			// Because the number of cases changes depending on the current
			// state, a dynamic select must be used
			cases := make([]reflect.SelectCase, numCases)
			cases[incomingCase].Dir = reflect.SelectRecv
			cases[outgoingCase].Dir = reflect.SelectSend

			// If the Send channel is still open, then add a select case to
			// receive from it - otherwise, leave it empty
			if !sendClosed {
				cases[incomingCase].Chan = reflect.ValueOf(send)
			}

			// If the list contains at least one item, add a select case to
			// send the first item from the list on the Recv channel
			if items.Len() > 0 {
				cases[outgoingCase].Chan = reflect.ValueOf(recv)
				cases[outgoingCase].Send = reflect.ValueOf(items.Front().Value)
			}

			// Three items will be returned - the index of the select case that
			// was chosen, the item received (if applicable), and the state of
			// the channel (open or closed)
			i, item, ok := reflect.Select(cases)

			// Switch on the chosen channel
			switch i {
			case incomingCase:

				// If an item was received, add it to the list - otherwise,
				// make a note that the Send channel was closed
				if ok {
					items.PushBack(item)
				} else {
					sendClosed = true
				}
			case outgoingCase:

				// The first item was sent, so it can be removed from the list
				items.Remove(items.Front())
			}
		}
	}()
	return nbChan
}
