package queue

import (
	"github.com/nathan-osman/go-cannon/util"

	"log"
	"time"
)

// Mail queue managing the sending of messages to hosts.
type Queue struct {
	storage    *Storage
	hosts      map[string]*Host
	newMessage *util.NonBlockingChan
	stop       chan bool
}

// Deliver the specified message to the appropriate host queue.
func (q *Queue) deliverMessage(host, id string) {
	if _, ok := q.hosts[host]; !ok {
		q.hosts[host] = NewHost(host, q.storage)
	}
	q.hosts[host].Deliver(id)
}

// Check for inactive host queues and shut them down.
func (q *Queue) checkForInactiveQueues() {
	for n, h := range q.hosts {
		if h.Idle() > time.Minute {
			h.Stop()
			delete(q.hosts, n)
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
			m := i.(*Message)
			if id, err := q.storage.NewMessage(m); err == nil {
				q.deliverMessage(m.Host, id)
			} else {
				log.Print(err)
			}
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

// Create a new message queue. Any undelivered messages on disk will be added
// to the appropriate queue.
func NewQueue(s *Storage) (*Queue, error) {
	q := &Queue{
		storage:    s,
		hosts:      make(map[string]*Host),
		newMessage: util.NewNonBlockingChan(),
		stop:       make(chan bool),
	}
	if messages, err := s.LoadMessages(); err == nil {
		for _, id := range messages {
			if m, err := q.storage.GetMessage(id); err == nil {
				q.deliverMessage(m.Host, id)
			}
		}
	} else {
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
