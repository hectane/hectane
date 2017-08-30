package imap

import (
	"github.com/emersion/go-imap/server"
	"github.com/sirupsen/logrus"
)

// IMAP allows clients to connect using IMAP and access their mailboxes.
type IMAP struct {
	server    *server.Server
	log       *logrus.Entry
	stoppedCh chan bool
}

// New creates a new IMAP server.
func New(cfg *Config) (*IMAP, error) {
	var (
		s = server.New(&dbBackend{
			storage: cfg.Storage,
		})
		i = &IMAP{
			server:    s,
			log:       logrus.WithField("context", "imap"),
			stoppedCh: make(chan bool),
		}
	)
	s.Addr = cfg.Addr
	s.AllowInsecureAuth = true
	go func() {
		defer close(i.stoppedCh)
		defer i.log.Info("IMAP server stopped")
		i.log.Info("starting IMAP server...")
		if err := s.ListenAndServe(); err != nil {
			i.log.Error(err.Error())
		}
	}()
	return i, nil
}

// Close shuts down the IMAP server.
func (i *IMAP) Close() {
	i.server.Close()
	<-i.stoppedCh
}
