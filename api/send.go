package api

import (
	"github.com/nathan-osman/go-cannon/email"
	"github.com/nathan-osman/go-cannon/queue"
	"github.com/zenazn/goji/web"

	"encoding/json"
	"net/http"
)

// Parameters for the /send method.
type sendParams struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	Html    string   `json:"html"`
}

// Send an email with the specified parameters.
func Send(c web.C, w http.ResponseWriter, r *http.Request) {
	var p sendParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, "malformed JSON")
	} else {
		if emails, err := email.NewEmails(p.From, p.To, p.Cc, p.Bcc, p.Subject, p.Text, p.Html); err != nil {
			respondWithError(w, err.Error())
		} else {
			for _, e := range emails {
				c.Env["queue"].(*queue.Queue).Deliver(e)
			}
			respondWithJSON(w, struct{}{})
		}
	}
}
