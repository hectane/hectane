package main

import (
	"github.com/goji/httpauth"
	"github.com/nathan-osman/go-cannon/api"
	"github.com/nathan-osman/go-cannon/queue"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/web"

	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
)

// Global mail queue exposed to the API methods.
var q *queue.Queue

// Goji middleware to expose the mail queue to the API methods.
func queueMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["queue"] = q
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func main() {

	// Storage for values supplied via command-line parameters
	var (
		tlsCert   string
		tlsKey    string
		username  string
		password  string
		directory string
	)

	// Set the default port, while still allowing for the usual overrides
	if s := bind.Sniff(); s == "" {
		flag.Lookup("bind").Value.Set(":8025")
		flag.Lookup("bind").DefValue = ":8025"
	}

	// Add command-line flags for each of the options and then parse them
	flag.StringVar(&tlsCert, "tls-cert", "", "certificate for TLS")
	flag.StringVar(&tlsKey, "tls-key", "", "private key for TLS")
	flag.StringVar(&username, "username", "", "username for HTTP basic auth")
	flag.StringVar(&password, "password", "", "password for HTTP basic auth")
	flag.StringVar(&directory, "directory", "", "directory for the mail queue")
	flag.Parse()

	// Create the mail queue
	var err error
	q, err := queue.NewQueue(directory)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer q.Stop()

	// Add the two current API methods
	goji.Get("/v1/version", api.Version)
	goji.Post("/v1/send", api.Send)

	// Add the queue middleware
	goji.Use(queueMiddleware)

	// If username and password were provided, enable HTTP basic auth
	if username != "" && password != "" {
		goji.Use(httpauth.SimpleBasicAuth(username, password))
	}

	// If a TLS certificate and key were provided, enable TLS and serve
	if tlsCert != "" && tlsKey != "" {
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		goji.ServeTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
		})
	} else {
		goji.Serve()
	}
}
