package api

import (
	"github.com/hectane/hectane/email"
	"github.com/hectane/hectane/version"

	"encoding/json"
	"net/http"
)

// Send a raw MIME message.
func (a *API) raw(r *http.Request) interface{} {
	var raw email.Raw
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return err
	}
	if err := raw.DeliverToQueue(a.queue); err != nil {
		return err
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
		"version": version.Version,
	}
}
