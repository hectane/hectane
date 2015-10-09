package api

import (
	"github.com/hectane/hectane/email"

	"encoding/json"
	"net/http"
)

// Send an email with the specified parameters.
func (a *API) send(w http.ResponseWriter, r *http.Request) {
	if a.validRequest(w, r, post) {
		var e email.Email
		if err := json.NewDecoder(r.Body).Decode(&e); err == nil {
			if messages, err := e.Messages(a.storage); err == nil {
				for _, m := range messages {
					a.queue.Deliver(m)
				}
				a.respondWithJSON(w, struct{}{})
			} else {
				a.respondWithJSON(w, map[string]string{
					"error": err.Error(),
				})
			}
		} else {
			a.respondWithJSON(w, map[string]string{
				"error": "unable to decode JSON",
			})
		}
	}
}
