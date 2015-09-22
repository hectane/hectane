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

	// Attempt to decode the parameters as an email.Email instance
	var e email.Email
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		respondWithError(w, err.Error())
	} else {

		// Ensure that if either 'text' or 'html' was not provided, its value
		// is populated by the other field
		if e.Html == "" {
			e.Html = html.EscapeString(e.Text)
		} else if e.Text == "" {
			e.Text = sanitize.HTML(e.Html)
		}

		// Convert the email into an array of messages
		q := c.Env["queue"].(*queue.Queue)
		if messages, err := e.Messages(q); err == nil {

			for _, m := range messages {
				q.Deliver(m)
			}

			// Respond with an empty object
			respondWithJSON(w, struct{}{})

		} else {
			respondWithError(w, err.Error())
		}
	}
}
