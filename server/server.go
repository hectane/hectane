package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

// Server provides the web interface for interacting with Hectane. This
// includes such things as managing users, viewing emails, etc.
type Server struct {
	listener net.Listener
	sessions *sessions.CookieStore
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
			sessions: sessions.NewCookieStore([]byte(cfg.SecretKey)),
			stopped:  make(chan bool),
		}
	)
	router.HandleFunc(
		"/api/auth/login",
		s.post(s.json(s.login, loginParams{})),
	)
	router.HandleFunc(
		"/api/auth/logout",
		s.post(s.auth(s.logout)),
	)
	router.HandleFunc(
		"/api/folders",
		s.auth(s.folders),
	)
	router.HandleFunc(
		"/api/folders/new",
		s.post(s.auth(s.json(s.newFolder, newFolderParams{}))),
	)
	router.HandleFunc(
		"/api/folders/{id:[0-9]+}/delete",
		s.post(s.auth(s.deleteFolder)),
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
