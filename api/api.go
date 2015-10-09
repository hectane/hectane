package api

import (
	"github.com/hectane/hectane/queue"

	"net/http"
)

// HTTP API for managing a mail queue.
type API struct {
	server  *http.Server
	config  *Config
	queue   *queue.Queue
	storage *queue.Storage
}

// Create a new API instance for the specified queue.
func New(config *Config, queue *queue.Queue, storage *queue.Storage) *API {
	var (
		s = http.NewServeMux()
		a = &API{
			server: &http.Server{
				Addr:    config.Addr,
				Handler: s,
			},
			config:  config,
			queue:   queue,
			storage: storage,
		}
	)
	s.HandleFunc("/v1/send", a.send)
	s.HandleFunc("/v1/status", a.status)
	s.HandleFunc("/v1/version", a.version)
	return a
}

// Listen for new requests.
func (a *API) Listen() error {
	if a.config.TLSCert != "" && a.config.TLSKey != "" {
		return a.server.ListenAndServeTLS(a.config.TLSKey, a.config.TLSKey)
	} else {
		return a.server.ListenAndServe()
	}
}

func (a *API) send(w http.ResponseWriter, r *http.Request)    {}
func (a *API) status(w http.ResponseWriter, r *http.Request)  {}
func (a *API) version(w http.ResponseWriter, r *http.Request) {}
