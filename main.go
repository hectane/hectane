package main

import (
	"github.com/goji/httpauth"
	"github.com/mitchellh/go-homedir"
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
	"path"
)

func main() {
	var (
		tlsCert   string
		tlsKey    string
		username  string
		password  string
		directory string
	)
	if s := bind.Sniff(); s == "" {
		flag.Lookup("bind").Value.Set(":8025")
		flag.Lookup("bind").DefValue = ":8025"
	}
	if home, err := homedir.Dir(); err == nil {
		directory = path.Join(home, ".go-cannon")
	} else {
		log.Println(err)
		os.Exit(1)
	}
	flag.StringVar(&tlsCert, "tls-cert", "", "certificate for TLS")
	flag.StringVar(&tlsKey, "tls-key", "", "private key for TLS")
	flag.StringVar(&username, "username", "", "username for HTTP basic auth")
	flag.StringVar(&password, "password", "", "password for HTTP basic auth")
	flag.StringVar(&directory, "directory", directory, "directory for persistent storage")
	flag.Parse()
	if q, err := queue.NewQueue(directory); err == nil {
		defer q.Stop()
		goji.Use(func(c *web.C, h http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				c.Env["queue"] = q
				h.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		})
	} else {
		log.Println(err)
		os.Exit(1)
	}
	goji.Get("/v1/version", api.Version)
	goji.Post("/v1/send", api.Send)
	if username != "" && password != "" {
		goji.Use(httpauth.SimpleBasicAuth(username, password))
	}
	if tlsCert != "" && tlsKey != "" {
		if cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey); err == nil {
			goji.ServeTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
			})
		} else {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		goji.Serve()
	}
}
