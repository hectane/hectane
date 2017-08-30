package queue

import (
	"time"

	"github.com/hectane/go-smtpsrv"
	"github.com/sirupsen/logrus"
)

// Queue manages an incoming mail queue, handling mail delivery as messages are
// received.
type Queue struct {
	server    *smtpsrv.Server
	log       *logrus.Entry
	stoppedCh chan bool
}

func (q *Queue) run() {
	defer close(q.stoppedCh)
	defer q.log.Info("incoming mail queue stopped")
	q.log.Info("incoming mail queue started")
	for m := range q.server.NewMessage {
		q.deliver(m)
	}
}

// New creates and starts a new queue.
func New(cfg *Config) (*Queue, error) {
	s, err := smtpsrv.NewServer(&smtpsrv.Config{
		Addr:        cfg.Addr,
		Banner:      "Hectane",
		ReadTimeout: 2 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	q := &Queue{
		server:    s,
		log:       logrus.WithField("context", "queue"),
		stoppedCh: make(chan bool),
	}
	go q.run()
	return q, nil
}

// Close shuts down the mail queue and waits for it to terminate.
func (q *Queue) Close() {
	q.server.Close(false)
	<-q.stoppedCh
}
