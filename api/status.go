package api

import (
	"net/http"
)

// Retrieve status information.
func (a *API) status(w http.ResponseWriter, r *http.Request) {
	if a.validRequest(w, r, get) {
		a.respondWithJSON(w, map[string]interface{}{
			"hosts": a.queue.Status(),
		})
	}
}
