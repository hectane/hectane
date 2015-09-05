package queue

import (
	"github.com/nathan-osman/go-cannon/email"
	"github.com/nathan-osman/go-cannon/util"

	"io/ioutil"
	"path"
	"time"
)

// Mail queue managing the sending of emails to hosts.
type Queue struct {
	newEmail *util.NonBlockingChan
	stop     chan bool
}

// Attempt to load all emails from the specified directory.
func loadEmails(directory string) ([]*email.Email, error) {

	// Create an array for storing the emails
	emails := []*email.Email{}

	// Enumerate the files in the directory
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	// Attempt to load each file and ignore ones that fail
	for _, f := range files {
		if e, err := email.LoadEmail(path.Join(directory, f.Name())); err == nil {
			emails = append(emails, e)
		}
	}

	return emails, nil
}

// Create a new mail queue.
func NewQueue(directory string) *Queue {

	// Create the two channels the queue will need
	q := &Queue{
		newEmail: util.NewNonBlockingChan(),
		stop:     make(chan bool),
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
			case i := <-q.newEmail.Recv:

				// Convert to an Email pointer
				e := i.(*email.Email)

				// Create the specified host if it doesn't exist
				if _, ok := hosts[e.Host]; !ok {
					hosts[e.Host] = NewHost(e.Host)
				}

				// Deliver the message to the host
				hosts[e.Host].Deliver(e)

			case <-ticker.C:

				// Loop through all of the hosts and remove ones that have been
				// idle for longer than 5 minutes and stops them
				for h := range hosts {
					if hosts[h].Idle() > 5*time.Minute {
						hosts[h].Stop()
						delete(hosts, h)
					}
				}

			case <-q.stop:
				break loop
			}
		}

		// Stop all host queues
		for h := range hosts {
			hosts[h].Stop()
		}
	}()

	return q
}

// Deliver the provided email.
func (q *Queue) Deliver(e *email.Email) {
	q.newEmail.Send <- e
}

// Stop all active host queues.
func (q *Queue) Stop() {

	// Send on the channel to stop it and wait for it to be closed
	q.stop <- true
	<-q.stop
}
