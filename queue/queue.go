package queue

import (
	"github.com/nathan-osman/go-cannon/email"
	"github.com/nathan-osman/go-cannon/util"

	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

// Mail queue managing the sending of emails to hosts.
type Queue struct {
	directory string
	hosts     map[string]*Host
	newEmail  *util.NonBlockingChan
	stop      chan bool
}

// Load all emails in the storage directory.
func (q *Queue) loadEmails() error {

	// Enumerate the files in the directory
	files, err := ioutil.ReadDir(q.directory)
	if err != nil {
		return err
	}

	// Attempt to load each file and ignore ones that fail
	for _, f := range files {
		if e, err := email.LoadEmail(path.Join(q.directory, f.Name())); err == nil {
			q.newEmail.Send <- e
		}
	}

	return nil
}

// Ensure the storage directory exists and load any emails in the directory.
func (q *Queue) prepareStorage() error {

	// If the directory exists, load the emails contained in it - otherwise,
	// attempt to create the directory
	if _, err := os.Stat(q.directory); err == nil {
		return q.loadEmails()
	} else {
		return os.MkdirAll(q.directory, 0755)
	}
}

// Deliver the specified email to the appropriate host queue.
func (q *Queue) deliverEmail(e *email.Email) {

	log.Printf("delivering email to %s queue", e.Host)

	// Save the email to the storage directory
	e.Save(q.directory)

	// Create the specified host if it doesn't exist
	if _, ok := q.hosts[e.Host]; !ok {
		q.hosts[e.Host] = NewHost(e.Host, q.directory)
	}

	// Deliver the message to the host
	q.hosts[e.Host].Deliver(e)
}

// Check for inactive host queues and shut them down.
func (q *Queue) checkForInactiveQueues() {
	for h := range q.hosts {
		if q.hosts[h].Idle() > 5*time.Minute {
			q.hosts[h].Stop()
			delete(q.hosts, h)
		}
	}
}

// Receive new emails and deliver them to the specified host queue.
func (q *Queue) run() {

	// Close the stop channel when the goroutine exits
	defer close(q.stop)

	// Create a ticker to periodically check for inactive hosts
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Loop to wait for (1) a new email (2) inactive timer (3) stop request
loop:
	for {
		select {
		case i := <-q.newEmail.Recv:
			q.deliverEmail(i.(*email.Email))
		case <-ticker.C:
			q.checkForInactiveQueues()
		case <-q.stop:
			break loop
		}
	}

	log.Println("shutting down host queues")

	// Stop all host queues
	for h := range q.hosts {
		q.hosts[h].Stop()
	}
}

// Create a new mail queue.
func NewQueue(directory string) (*Queue, error) {

	q := &Queue{
		directory: directory,
		hosts:     make(map[string]*Host),
		newEmail:  util.NewNonBlockingChan(),
		stop:      make(chan bool),
	}

	// Prepare the storage directory
	if err := q.prepareStorage(); err != nil {
		return nil, err
	}

	// Start a goroutine to manage the lifecycle of the queue
	go q.run()

	return q, nil
}

// Deliver the specified email to the appropriate host queue.
func (q *Queue) Deliver(e *email.Email) {
	q.newEmail.Send <- e
}

// Stop all active host queues.
func (q *Queue) Stop() {
	q.stop <- true
	<-q.stop
}
