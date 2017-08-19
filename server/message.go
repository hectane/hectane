package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
)

func (s *Server) messages(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*db.User)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	m, err := db.Messages(db.DefaultToken, u.ID, id)
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, m)
}
