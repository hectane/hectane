package api

import (
	"github.com/hectane/hectane/queue"

	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
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
	stop     chan bool
	serveMux *http.ServeMux
	config   *Config
	queue    *queue.Queue
}

// Log the specified message.
func (a *API) log(msg string, v ...interface{}) {
	log.Printf(fmt.Sprintf("[API] %s", msg), v...)
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
		server: &http.Server{
			Addr: config.Addr,
		},
		stop:     make(chan bool),
		serveMux: http.NewServeMux(),
		config:   config,
		queue:    queue,
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

// Listen for new requests until an error occurs.
func (a *API) Listen() error {
	if l, err := net.Listen("tcp", a.config.Addr); err == nil {
		if a.config.TLSCert != "" && a.config.TLSKey != "" {
			if cert, err := tls.LoadX509KeyPair(a.config.TLSCert, a.config.TLSKey); err == nil {
				l = tls.NewListener(l, &tls.Config{
					Certificates: []tls.Certificate{cert},
				})
			} else {
				return err
			}
		}
		var (
			done = make(chan bool)
			err  error
		)
		go func() {
			a.log("binding to %s", a.config.Addr)
			err = a.server.Serve(l)
			// Supress benign errors - see http://bit.ly/1WUhgDj
			if oe, ok := err.(*net.OpError); ok && oe.Op == "accept" || oe.Op == "AcceptEx" {
				err = nil
			}
			if err == nil {
				a.log("API server shutdown")
			}
			close(done)
		}()
		for {
			select {
			case <-a.stop:
				l.Close()
			case <-done:
				return err
			}
		}
	} else {
		return err
	}
}

// If currently listening, stop and shutdown the server.
func (a *API) Stop() {
	if a.stop != nil {
		a.stop <- true
	}
}
