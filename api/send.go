package api

import (
	"github.com/kennygrant/sanitize"
	"github.com/nathan-osman/go-cannon/email"
	"github.com/nathan-osman/go-cannon/queue"
	"github.com/zenazn/goji/web"

	"encoding/json"
	"html"
	"net/http"
)

// Send an email with the specified parameters.
func Send(c web.C, w http.ResponseWriter, r *http.Request) {
	var e email.Email
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		respondWithError(w, err.Error())
	} else {
		if e.Html == "" {
			e.Html = html.EscapeString(e.Text)
		} else if e.Text == "" {
			e.Text = sanitize.HTML(e.Html)
		}
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
