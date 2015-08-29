package main

import (
	"github.com/goji/httpauth"
	"github.com/zenazn/goji"

	"flag"
)

func main() {
	var (
		username string
		password string
	)

	flag.StringVar(&username, "username", "", "username for HTTP basic auth")
	flag.StringVar(&password, "password", "", "password for HTTP basic auth")
	flag.Parse()

	// If username and password are provided, enable HTTP basic auth
	if username != "" && password != "" {
		goji.Use(httpauth.SimpleBasicAuth(username, password))
	}

	goji.Get("/v1/version", Version)
	goji.Post("/v1/send", Send)
	goji.Serve()
}
