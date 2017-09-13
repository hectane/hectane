package sender

import (
	"time"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/sender/queue"
	"github.com/hectane/hectane/storage"
	"github.com/sirupsen/logrus"
)

const pollInterval = 2 * time.Second

// Sender monitors the database for qualified hosts with messages for delivery
// and starts a queue for their delivery.
type Sender struct {
	storage   *storage.Storage
	log       *logrus.Entry
	stopCh    chan bool
	stoppedCh chan bool
}

// run creates queues in response to available hosts. A map is used to track
// active queues so that they can easily be shut down when the sender is being
// closed. A queue indicates that it is finished and has shut down by sending
// its host's name on the queueFinishedCh channel.
func (s *Sender) run() {
	var (
		queues          = make(map[string]*queue.Queue)
		queueFinishedCh = make(chan string)
	)
	defer func() {
		s.log.Info("shutting down queues...")
		for _, q := range queues {
			q.Close()
		}
		s.log.Info("sender shut down")
		close(s.stoppedCh)
	}()
	s.log.Info("sender started")
	for {
		h, err := db.GetAvailableHost()
		if err == nil && h != nil {
			queues[h.Name] = queue.New(&queue.Config{
				Host:            h,
				Storage:         s.storage,
				QueueFinishedCh: queueFinishedCh,
			})
		}
		if err != nil {
			s.log.Error(err.Error())
		}
		select {
		case name := <-queueFinishedCh:
			delete(queues, name)
		case <-time.After(pollInterval):
		case <-s.stopCh:
			return
		}
	}
}

// New creates a new sender.
func New(cfg *Config) *Sender {
	s := &Sender{
		storage:   cfg.Storage,
		log:       logrus.WithField("context", "sender"),
		stopCh:    make(chan bool),
		stoppedCh: make(chan bool),
	}
	go s.run()
	return s
}

// Close shuts down the sender.
func (s *Sender) Close() {
	close(s.stopCh)
	<-s.stoppedCh
}
