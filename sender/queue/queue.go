package queue

import (
	"fmt"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/storage"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	host            *db.Host
	storage         *storage.Storage
	log             *logrus.Entry
	queueFinishedCh chan string
	stopCh          chan bool
	stoppedCh       chan bool
}

func (q *Queue) run() {
	defer func() {
		select {
		case q.queueFinishedCh <- q.host.Name:
		default:
		}
		q.log.Debug("queue shut down")
		close(q.stoppedCh)
	}()
	q.log.Debug("queue started")
	for {
		select {
		case <-q.stopCh:
			return
		}
	}
}

// New creates a new queue for the specified host.
func New(cfg *Config) *Queue {
	q := &Queue{
		host:            cfg.Host,
		storage:         cfg.Storage,
		log:             logrus.WithField("context", fmt.Sprintf("queue[%s]", cfg.Host.Name)),
		queueFinishedCh: cfg.QueueFinishedCh,
		stopCh:          make(chan bool),
		stoppedCh:       make(chan bool),
	}
	go q.run()
	return q
}

// Close shuts down the queue.
func (q *Queue) Close() {
	close(q.stopCh)
	<-q.stoppedCh
}
