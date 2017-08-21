package server

import (
	"net/http"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/db/models"
	"github.com/hectane/hectane/db/util"
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
	i, err := util.SelectItem(db.Token, models.User{}, util.SelectParams{
		Where: &util.EqClause{
			Field: "Username",
			Value: p.Username,
		},
	})
	if err != nil {
		http.Error(w, statusInvalidCredentials, http.StatusForbidden)
		return
	}
	u := i.(*models.User)
	if err := u.Authenticate(p.Password); err != nil {
		http.Error(w, statusInvalidCredentials, http.StatusForbidden)
		return
	}
	session, _ := s.sessions.Get(r, sessionName)
	session.Values[sessionUserID] = u.ID
	session.Save(r, w)
	s.writeJson(w, u)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.Values[sessionUserID] = ""
	s.writeJson(w, nil)
}
