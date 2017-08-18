package server

import (
	"net/http"

	"github.com/hectane/hectane/db"
)

const (
	statusInvalidCredentials = "invalid username or password"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(contextParams).(*loginParams)
	u, err := db.FindUser(db.DefaultToken, "Username", p.Username)
	if err != nil {
		http.Error(w, statusInvalidCredentials, http.StatusForbidden)
		return
	}
	if err := u.Authenticate(p.Password); err != nil {
		http.Error(w, statusInvalidCredentials, http.StatusForbidden)
		return
	}
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.Values[sessionUserID] = u.ID
	s.writeJson(w, u)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.Values[sessionUserID] = ""
	s.writeJson(w, nil)
}
