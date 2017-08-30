package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/server/auth"
	"github.com/hectane/hectane/server/resources"
	"github.com/manyminds/api2go"
	"github.com/sirupsen/logrus"
)

// Server provides the web interface for interacting with Hectane. This
// includes such things as managing users, viewing emails, etc.
type Server struct {
	listener net.Listener
	auth     *auth.Auth
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
		api = api2go.NewAPI("api")
		s   = &Server{
			listener: l,
			auth:     auth.New(api.Handler(), cfg.SecretKey),
			log:      logrus.WithField("context", "server"),
			stopped:  make(chan bool),
		}
	)
	api.AddResource(&db.Domain{}, resources.DomainResource)
	api.AddResource(&db.Folder{}, resources.FolderResource)
	api.AddResource(&db.Message{}, resources.MessageResource)
	api.AddResource(&db.User{}, resources.UserResource)
	router.HandleFunc("/api/login", s.login)
	router.HandleFunc("/api/logout", s.logout)
	router.PathPrefix("/api/").Handler(s.auth)
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
