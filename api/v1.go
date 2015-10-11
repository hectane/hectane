package api

import (
	"github.com/hectane/hectane/email"

	"encoding/json"
	"net/http"
)

// Send an email with the specified parameters.
func (a *API) send(r *http.Request) interface{} {
	var e email.Email
	if err := json.NewDecoder(r.Body).Decode(&e); err == nil {
		if messages, err := e.Messages(a.queue.Storage); err == nil {
			for _, m := range messages {
				a.queue.Deliver(m)
			}
			return struct{}{}
		} else {
			return map[string]string{
				"error": err.Error(),
			}
		}
	} else {
		return map[string]string{
			"error": "unable to decode JSON",
		}
	}
}

// Retrieve status information.
func (a *API) status(r *http.Request) interface{} {
	return map[string]interface{}{
		"hosts": a.queue.Status(),
	}
}

// Retrieve version information, including the current version of the
// application.
func (a *API) version(r *http.Request) interface{} {
	return map[string]string{
		"version": "0.3.0",
	}
}
