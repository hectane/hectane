package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	statusDatabaseError  = "database error"
	statusObjectNotFound = "object not found"
)

// writeJson writes JSON data to the client.
func (s *Server) writeJson(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		s.log.Error(err)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
