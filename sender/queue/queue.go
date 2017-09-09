package queue

import (
	"context"

	"github.com/hectane/hectane/storage"
	"github.com/sirupsen/logrus"
)

var queueNum int

// Queue attempts to deliver messages as efficiently as possible. A host is
// "locked" and messages are delivered continuously to it until either no more
// messages remain or an error occurs. The queue will then check for other
// available hosts and shut down if none exist.
type Queue struct {
	storage   *storage.Storage
	log       *logrus.Entry
	stopCh    chan bool
	stoppedCh chan bool
}

func (q *Queue) run() {
	defer close(q.stoppedCh)
	defer q.log.Debug("queue shut down")
	q.log.Debug("queue started")
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		<-q.stopCh
		cancelFn()
	}()
	for {
		if err := q.process(ctx); err != nil {
			q.log.Error(err.Error())
			return
		}
	}
}

// New creates and initializes a queue for the host.
func New(cfg *Config) *Queue {
	queueNum++
	q := &Queue{
		storage: cfg.Storage,
		log: logrus.WithFields(logrus.Fields{
			"context": "queue",
			"number":  queueNum,
		}),
		stopCh:    make(chan bool),
		stoppedCh: make(chan bool),
	}
	go q.run()
	return q
}

// Close shuts down the queue.
func (q *Queue) Close() {
	close(q.stopCh)
	<-q.stoppedCh
}
