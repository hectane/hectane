package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/hectane/go-asyncserver"
	"github.com/hectane/hectane/queue"

	"crypto/tls"
	"encoding/json"
	"net/http"
	"strconv"
)

// Request methods.
const (
	get  = "GET"
	post = "POST"
)

// HTTP API for managing a mail queue.
type API struct {
	config   *Config
	log      *logrus.Entry
	server   *server.AsyncServer
	serveMux *http.ServeMux
	queue    *queue.Queue
	stopped  chan bool
}

// Create a handler that logs and validates requests as they come in. The
// return value of the handler is assumed to be either an error or a map.
func (a *API) method(method string, handler func(r *http.Request) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			v := handler(r)
			if err, ok := v.(error); ok {
				v = map[string]string{
					"error": err.Error(),
				}
			}
			if data, err := json.Marshal(v); err == nil {
				w.Header().Set("Content-Length", strconv.Itoa(len(data)))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(data)
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
		config:   config,
		log:      logrus.WithField("context", "API"),
		server:   server.New(config.Addr),
		serveMux: http.NewServeMux(),
		queue:    queue,
		stopped:  make(chan bool),
	}
	a.server.Handler = a
	a.serveMux.HandleFunc("/v1/raw", a.method(post, a.raw))
	a.serveMux.HandleFunc("/v1/send", a.method(post, a.send))
	a.serveMux.HandleFunc("/v1/status", a.method(get, a.status))
	a.serveMux.HandleFunc("/v1/version", a.method(get, a.version))
	return a
}

// Process an incoming request. This method logs the request and checks to
// ensure that HTTP basic auth credentials were supplied if required.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.log.Debugf("%s - %s %s", r.RemoteAddr, r.Method, r.RequestURI)
	if a.config.Username != "" && a.config.Password != "" {
		username, password, ok := r.BasicAuth()
		if !ok || username != a.config.Username || password != a.config.Password {
			w.Header().Set("WWW-Authenticate", "Basic realm=Hectane")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}
	a.serveMux.ServeHTTP(w, r)
}

// Start listening for new requests.
func (a *API) Start() error {
	if a.config.TLSCert != "" && a.config.TLSKey != "" {
		c, err := tls.LoadX509KeyPair(a.config.TLSCert, a.config.TLSKey)
		if err != nil {
			return err
		}
		a.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{c},
		}
	}
	return a.server.Start()
}

// Stop listening for new requests.
func (a *API) Stop() {
	a.server.Stop()
}
