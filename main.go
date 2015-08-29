package main

import (
	"github.com/goji/httpauth"
	"github.com/zenazn/goji"

	"crypto/tls"
	"flag"
	"log"
	"os"
)

func main() {
	var (
		tlsCert  string
		tlsKey   string
		username string
		password string
	)

	// Set the default port
	flag.Lookup("bind").DefValue = ":8025"

	flag.StringVar(&tlsCert, "tls-cert", "", "certificate for TLS")
	flag.StringVar(&tlsKey, "tls-key", "", "private key for TLS")
	flag.StringVar(&username, "username", "", "username for HTTP basic auth")
	flag.StringVar(&password, "password", "", "password for HTTP basic auth")
	flag.Parse()

	goji.Get("/v1/version", Version)
	goji.Post("/v1/send", Send)

	// If username and password are provided, enable HTTP basic auth
	if username != "" && password != "" {
		goji.Use(httpauth.SimpleBasicAuth(username, password))
	}

	// If a TLS certificate and key were provided, enable TLS
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
