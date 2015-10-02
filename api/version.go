package api

import (
	"github.com/zenazn/goji/web"

	"net/http"
)

// Retrieve version information.
func Version(c web.C, w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, map[string]string{
		"version": "0.2.1",
	})
}
