package api

import (
	"net/http"
)

// Retrieve status information.
func (a *API) status(r *http.Request) interface{} {
	return map[string]interface{}{
		"hosts": a.queue.Status(),
	}
}
