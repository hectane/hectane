package receiver

import (
	"time"

	"github.com/hectane/go-smtpsrv"
	"github.com/hectane/hectane/storage"
	"github.com/sirupsen/logrus"
)

// Receiver manages an incoming mail queue, handling mail delivery as messages
// are received.
type Receiver struct {
	server    *smtpsrv.Server
	storage   *storage.Storage
	log       *logrus.Entry
	stoppedCh chan bool
}

func (r *Receiver) run() {
	defer close(r.stoppedCh)
	defer r.log.Info("mail receiver stopped")
	r.log.Info("mail receiver started")
	for m := range r.server.NewMessage {
		r.deliver(m)
	}
}

// New creates and starts a new receiver.
func New(cfg *Config) (*Receiver, error) {
	s, err := smtpsrv.NewServer(&smtpsrv.Config{
		Addr:        cfg.Addr,
		Banner:      "Hectane",
		ReadTimeout: 2 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	r := &Receiver{
		server:    s,
		storage:   cfg.Storage,
		log:       logrus.WithField("context", "receiver"),
		stoppedCh: make(chan bool),
	}
	go r.run()
	return r, nil
}

// Close shuts down the receiver and waits for it to terminate.
func (r *Receiver) Close() {
	r.server.Close(false)
	<-r.stoppedCh
}
