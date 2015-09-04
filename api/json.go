package api

import (
	"encoding/json"
	"net/http"
)

// Write the specified JSON object to the client. No error checking is done
// since little could be done if the method were to fail anyway.
func respondWithJSON(w http.ResponseWriter, o interface{}) {
	d, _ := json.Marshal(o)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(d)
}

// Write the specified error message to the client.
func respondWithStatus(w http.ResponseWriter, m string) {
	respondWithJSON(w, map[string]string{
		"error": m,
	})
}
