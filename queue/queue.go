package queue

import (
	"github.com/nathan-osman/go-cannon/util"

	"time"
)

// Mail queue managing the sending of messages to hosts.
type Queue struct {
	newMessage *util.NonBlockingChan
	stop       chan bool
}

// Create a new mail queue.
func NewQueue() *Queue {

	// Create the two channels the queue will need
	q := &Queue{
		newMessage: util.NewNonBlockingChan(),
		stop:       make(chan bool),
	}

	// Start a goroutine to manage the lifecycle of the queue
	go func() {

		// Close the stop channel when the goroutine exits
		defer close(q.stop)

		// Create a map of hosts and a ticker for freeing up unused hosts
		var (
			hosts  = make(map[string]*Host)
			ticker = time.NewTicker(5 * time.Minute)
		)

		// Stop the ticker when the goroutine exits
		defer ticker.Stop()

		// Main "loop" of the queue
	loop:
		for {
			select {
			case i := <-q.newMessage.Recv:

				// Convert to a Message pointer
				msg := i.(*Message)

				// Create the specified host if it doesn't exist
				if _, ok := hosts[msg.Host]; !ok {
					hosts[msg.Host] = NewHost(msg.Host)
				}

				// Deliver the message to the host
				hosts[msg.Host].Deliver(msg)

			case <-ticker.C:

				// Loop through all of the hosts and remove ones that have been
				// idle for longer than 5 minutes and stops them
				for host := range hosts {
					if hosts[host].Idle() > 5*time.Minute {
						hosts[host].Stop()
						delete(hosts, host)
					}
				}

			case <-q.stop:
				break loop
			}
		}

		// Stop all host queues
		for host := range hosts {
			hosts[host].Stop()
		}
	}()

	return q
}

// Deliver the provided message
func (q *Queue) Deliver(msg *Message) {
	q.newMessage.Send <- msg
}

// Stop all active host queues.
func (q *Queue) Stop() {

	// Send on the channel to stop it and wait for it to be closed
	q.stop <- true
	<-q.stop
}
