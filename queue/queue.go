package queue

import (
	"github.com/Sirupsen/logrus"
	"github.com/hectane/hectane/util"

	"sync"
	"time"
)

// Queue status information.
type QueueStatus struct {
	Uptime int                    `json:"uptime"`
	Hosts  map[string]*HostStatus `json:"hosts"`
}

// Mail queue managing the sending of messages to hosts.
type Queue struct {
	m          sync.Mutex
	config     *Config
	Storage    *Storage
	log        *logrus.Entry
	hosts      map[string]*Host
	newMessage *util.NonBlockingChan
	startTime  time.Time
	stop       chan bool
}

// Deliver the specified message to the appropriate host queue.
func (q *Queue) deliverMessage(m *Message) {
	q.m.Lock()
	if _, ok := q.hosts[m.Host]; !ok {
		q.hosts[m.Host] = NewHost(m.Host, q.Storage, q.config)
	}
	q.hosts[m.Host].Deliver(m)
	q.m.Unlock()
}

// Check for inactive host queues and shut them down.
func (q *Queue) checkForInactiveQueues() {
	q.m.Lock()
	for n, h := range q.hosts {
		if h.Idle() > time.Minute {
			h.Stop()
			delete(q.hosts, n)
		}
	}
	q.m.Unlock()
}

// Receive new messages and deliver them to the specified host queue. Check for
// idle queues every so often and shut them down if they haven't been used.
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
	q.log.Info("stopping host queues")
	q.m.Lock()
	for h := range q.hosts {
		q.hosts[h].Stop()
	}
	q.m.Unlock()
	q.log.Info("shutting down")
}

// Create a new message queue. Any undelivered messages on disk will be added
// to the appropriate queue.
func NewQueue(c *Config) (*Queue, error) {
	q := &Queue{
		config:     c,
		Storage:    NewStorage(c.Directory),
		log:        logrus.WithField("context", "Queue"),
		hosts:      make(map[string]*Host),
		newMessage: util.NewNonBlockingChan(),
		startTime:  time.Now(),
		stop:       make(chan bool),
	}
	messages, err := q.Storage.LoadMessages()
	if err != nil {
		return nil, err
	}
	q.log.Infof("loaded %d message(s) from %s", len(messages), c.Directory)
	for _, m := range messages {
		q.deliverMessage(m)
	}
	go q.run()
	return q, nil
}

// Provide the status of each host queue.
func (q *Queue) Status() *QueueStatus {
	s := &QueueStatus{
		Uptime: int(time.Now().Sub(q.startTime) / time.Second),
		Hosts:  make(map[string]*HostStatus),
	}
	q.m.Lock()
	for n, h := range q.hosts {
		s.Hosts[n] = h.Status()
	}
	q.m.Unlock()
	return s
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
