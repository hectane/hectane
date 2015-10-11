package queue

import (
	"github.com/hectane/hectane/util"

	"fmt"
	"log"
	"time"
)

// Mail queue managing the sending of messages to hosts.
type Queue struct {
	Storage    *Storage
	config     *Config
	hosts      map[string]*Host
	newMessage *util.NonBlockingChan
	stop       chan bool
}

// Log the specified message.
func (q *Queue) log(msg string, v ...interface{}) {
	log.Printf(fmt.Sprintf("[Queue] %s", msg), v...)
}

// Deliver the specified message to the appropriate host queue.
func (q *Queue) deliverMessage(m *Message) {
	if _, ok := q.hosts[m.Host]; !ok {
		q.hosts[m.Host] = NewHost(m.Host, q.Storage, q.config)
	}
	q.hosts[m.Host].Deliver(m)
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
			q.deliverMessage(i.(*Message))
		case <-ticker.C:
			q.checkForInactiveQueues()
		case <-q.stop:
			break loop
		}
	}
	q.log("shutting down host queues")
	for h := range q.hosts {
		q.hosts[h].Stop()
	}
	q.log("mail queue shutdown")
}

// Create a new message queue. Any undelivered messages on disk will be added
// to the appropriate queue.
func NewQueue(c *Config) (*Queue, error) {
	q := &Queue{
		Storage:    NewStorage(c.Directory),
		config:     c,
		hosts:      make(map[string]*Host),
		newMessage: util.NewNonBlockingChan(),
		stop:       make(chan bool),
	}
	if messages, err := q.Storage.LoadMessages(); err == nil {
		q.log("loaded %d message(s) from storage", len(messages))
		for _, m := range messages {
			q.deliverMessage(m)
		}
	} else {
		return nil, err
	}
	go q.run()
	return q, nil
}

// Provide the status of each host queue.
func (q *Queue) Status() map[string]interface{} {
	m := make(map[string]interface{})
	for n, h := range q.hosts {
		m[n] = map[string]interface{}{
			"active": h.Idle() == 0,
			"idle":   h.Idle() / time.Second,
		}
	}
	return m
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
