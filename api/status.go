package api

import (
	"github.com/hectane/hectane/queue"
	"github.com/zenazn/goji/web"

	"net/http"
)

// Retrieve status information.
func Status(c web.C, w http.ResponseWriter, r *http.Request) {
	q := c.Env["queue"].(*queue.Queue)
	respondWithJSON(w, map[string]interface{}{
		"hosts": q.Status(),
	})
}
