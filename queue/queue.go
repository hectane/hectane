package queue

import (
	"github.com/nathan-osman/go-cannon/util"

	"log"
	"os"
	"time"
)

// Mail queue managing the sending of messages to hosts.
type Queue struct {
	directory  string
	hosts      map[string]*Host
	newMessage *util.NonBlockingChan
	stop       chan bool
}

// Ensure the storage directory exists and load any messages in the directory.
// If the directory does not exist, it will be created.
func (q *Queue) prepareStorage() error {
	if _, err := os.Stat(q.directory); err == nil {
		if messages, err := LoadMessages(q.directory); err == nil {
			for _, m := range messages {
				q.newMessage.Send <- m
			}
			return nil
		} else {
			return err
		}
	} else {
		return os.MkdirAll(q.directory, 0700)
	}
}

// Deliver the specified message to the appropriate host queue.
func (q *Queue) deliverMessage(m *Message) {
	log.Printf("delivering message to %s queue", m.m.Host)
	if _, ok := q.hosts[m.m.Host]; !ok {
		q.hosts[m.m.Host] = NewHost(m.m.Host)
	}
	q.hosts[m.m.Host].Deliver(m)
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

// Receive new messages and deliver them to the specified host queue.
func (q *Queue) run() {
	defer close(q.stop)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
loop:
	for {
		select {
		case i := <-q.newMessage.Recv:
			q.deliverMessage(i.(*Message))
		case <-ticker.C:
			q.checkForInactiveQueues()
		case <-q.stop:
			break loop
		}
	}
	log.Println("shutting down host queues")
	for h := range q.hosts {
		q.hosts[h].Stop()
	}
}

// Create a new message queue.
func NewQueue(directory string) (*Queue, error) {
	q := &Queue{
		directory:  directory,
		hosts:      make(map[string]*Host),
		newMessage: util.NewNonBlockingChan(),
		stop:       make(chan bool),
	}
	if err := q.prepareStorage(); err != nil {
		return nil, err
	}
	go q.run()
	return q, nil
}

// Deliver the specified message to the appropriate host queue.
func (q *Queue) Deliver(m *Message) {
	q.newMessage.Send <- m
}

// Stop all active host queues.
func (q *Queue) Stop() {
	q.stop <- true
	<-q.stop
}
