package main

import (
	"github.com/goji/httpauth"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"flag"
	"net/http"
)

func send(c web.C, w http.ResponseWriter, r *http.Request) {
	// TODO: do stuff here
}

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

	goji.Post("/v1/send", send)
	goji.Serve()
}
