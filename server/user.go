package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
)

const (
	statusInvalidUsername = "invalid username"
)

func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	u, err := db.Users(db.DefaultToken)
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, u)
}

type newUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func (s *Server) newUser(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(contextParams).(*newUserParams)
	if len(p.Username) == 0 || len(p.Username) > 40 {
		http.Error(w, statusInvalidUsername, http.StatusBadRequest)
		return
	}
	u := &db.User{
		Username: p.Username,
		IsAdmin:  p.IsAdmin,
	}
	if err := u.SetPassword(p.Password); err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	if err := u.Save(db.DefaultToken); err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, u)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := db.Transaction(func(t *db.Token) error {
		u, err := db.FindUser(t, "ID", id)
		if err != nil {
			return err
		}
		return u.Delete(t)
	})
	if err == db.ErrRowCount {
		http.Error(w, statusObjectNotFound, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, nil)
}
