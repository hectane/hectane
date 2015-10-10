package api

import (
	"github.com/hectane/hectane/queue"

	"encoding/json"
	"log"
	"net/http"
)

// Request methods.
const (
	get  = "GET"
	post = "POST"
)

// HTTP API for managing a mail queue.
type API struct {
	server   *http.Server
	serveMux *http.ServeMux
	config   *Config
	queue    *queue.Queue
}

// Create a handler that logs and validates requests as they come in.
func (a *API) method(method string, handler func(r *http.Request) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			if v, err := json.Marshal(handler(r)); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(v)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

// Create a new API instance for the specified queue.
func New(config *Config, queue *queue.Queue) *API {
	a := &API{
		server: &http.Server{
			Addr: config.Addr,
		},
		serveMux: http.NewServeMux(),
		config:   config,
		queue:    queue,
	}
	a.server.Handler = a
	a.serveMux.HandleFunc("/v1/send", a.method(post, a.send))
	a.serveMux.HandleFunc("/v1/status", a.method(get, a.status))
	a.serveMux.HandleFunc("/v1/version", a.method(get, a.version))
	return a
}

// Process an incoming request. This method logs the request and checks to
// ensure that HTTP basic auth credentials were supplied if required.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[API] %s - %s %s", r.RemoteAddr, r.Method, r.RequestURI)
	if a.config.Username != "" && a.config.Password != "" {
		if username, password, ok := r.BasicAuth(); ok {
			if username != a.config.Username || password != a.config.Password {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Hectane")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}
	a.serveMux.ServeHTTP(w, r)
}

// Listen for new requests.
func (a *API) Listen() error {
	log.Printf("[API] Listening on %s...", a.config.Addr)
	if a.config.TLSCert != "" && a.config.TLSKey != "" {
		return a.server.ListenAndServeTLS(a.config.TLSCert, a.config.TLSKey)
	} else {
		return a.server.ListenAndServe()
	}
}
