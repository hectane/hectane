package server

import (
	"net/http"

	"github.com/hectane/hectane/db"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(contextParams).(*loginParams)
	u, err := db.FindUser(db.DefaultToken, "Username", p.Username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if err := u.Authenticate(p.Password); err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.Values[sessionUserID] = u.ID
	s.writeJson(w, map[string]interface{}{
		"user": u,
	})
}
