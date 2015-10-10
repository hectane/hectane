package api

import (
	"net/http"
)

// Retrieve version information, including the current version of the
// application.
func (a *API) version(r *http.Request) interface{} {
	return map[string]string{
		"version": "0.3.0",
	}
}
