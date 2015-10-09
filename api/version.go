package api

import (
	"net/http"
)

// Retrieve version information, including the current version of the
// application.
func (a *API) version(w http.ResponseWriter, r *http.Request) {
	a.respondWithJSON(w, map[string]string{
		"version": "0.3.0",
	})
}
