package api

import (
	"github.com/hectane/hectane/email"
	"github.com/hectane/hectane/queue"

	"encoding/json"
	"net/http"
)

type rawParams struct {
	From string   `json:"from"`
	To   []string `json:"to"`
	Body string   `json:"body"`
}

// Send a raw MIME message.
func (a *API) raw(r *http.Request) interface{} {
	var p rawParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return err
	}
	w, body, err := a.queue.Storage.NewBody()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(p.Body)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	hostMap, err := email.GroupAddressesByHost(p.To)
	if err == nil {
		return err
	}
	for h, to := range hostMap {
		m := &queue.Message{
			Host: h,
			From: p.From,
			To:   to,
		}
		if err := a.queue.Storage.SaveMessage(m, body); err != nil {
			return err
		}
		a.queue.Deliver(m)
	}
	return struct{}{}
}

// Send an email with the specified parameters.
func (a *API) send(r *http.Request) interface{} {
	var e email.Email
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return err
	}
	messages, err := e.Messages(a.queue.Storage)
	if err != nil {
		return map[string]string{
			"error": err.Error(),
		}
	}
	for _, m := range messages {
		a.queue.Deliver(m)
	}
	return struct{}{}
}

// Retrieve status information.
func (a *API) status(r *http.Request) interface{} {
	return a.queue.Status()
}

// Retrieve version information, including the current version of the
// application.
func (a *API) version(r *http.Request) interface{} {
	return map[string]string{
		"version": "0.3.1",
	}
}
