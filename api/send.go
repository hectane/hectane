package api

import (
	"github.com/hectane/hectane/email"
	"github.com/hectane/hectane/queue"
	"github.com/zenazn/goji/web"

	"encoding/json"
	"net/http"
)

// Send an email with the specified parameters.
func Send(c web.C, w http.ResponseWriter, r *http.Request) {
	var e email.Email
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		respondWithError(w, err.Error())
	} else {
		s := c.Env["storage"].(*queue.Storage)
		q := c.Env["queue"].(*queue.Queue)
		if messages, err := e.Messages(s); err == nil {
			for _, m := range messages {
				q.Deliver(m)
			}
			respondWithJSON(w, struct{}{})
		} else {
			respondWithError(w, err.Error())
		}
	}
}
