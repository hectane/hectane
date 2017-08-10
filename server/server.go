package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server provides the web interface for interacting with Hectane. This
// includes such things as managing users, viewing emails, etc.
type Server struct {
	listener net.Listener
	log      *logrus.Entry
	stopped  chan bool
}

// New creates a new server.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		router = mux.NewRouter()
		server = http.Server{
			Handler: router,
		}
		s = &Server{
			listener: l,
			log:      logrus.WithField("context", "server"),
			stopped:  make(chan bool),
		}
	)
	router.PathPrefix("/").Handler(http.FileServer(HTTP))
	go func() {
		defer close(s.stopped)
		defer s.log.Info("web server has stopped")
		s.log.Info("starting web server...")
		if err := server.Serve(l); err != nil {
			s.log.Error(err)
		}
	}()
	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stopped
}
