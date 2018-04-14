package smtp

import (
	"github.com/sirupsen/logrus"
	"github.com/hectane/go-smtpsrv"
	"github.com/hectane/hectane/email"
	"github.com/hectane/hectane/queue"
)

// Server awaits incoming connections and delivers them to the mail queue.
type Server struct {
	server *smtpsrv.Server
	queue  *queue.Queue
	log    *logrus.Entry
}

// run continuously delivers
func (s *Server) run() {
	for m := range s.server.NewMessage {
		s.log.Info("email received via SMTP")
		raw := email.Raw{
			From: m.From,
			To:   m.To,
			Body: m.Body,
		}
		if err := raw.DeliverToQueue(s.queue); err != nil {
			s.log.Error(err.Error())
		}
	}
}

// New creates a new SMTP server with the specified configuration.
func New(c *Config, q *queue.Queue) (*Server, error) {
	server, err := smtpsrv.NewServer(c.smtpsrvConfig())
	if err != nil {
		return nil, err
	}
	s := &Server{
		server: server,
		queue:  q,
		log:    logrus.WithField("context", "SMTP"),
	}
	go s.run()
	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Close(false)
}
