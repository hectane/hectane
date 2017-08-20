package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/db/util"
)

const (
	statusInvalidUsername = "invalid username"
)

func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	i, err := util.SelectItems(db.Token, db.User{}, util.SelectParams{
		OrderBy: "Username",
	})
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, i)
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
	if err := util.InsertItem(db.Token, u); err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, u)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := db.Token.Transaction(func(t *util.Token) error {
		i, err := util.SelectItem(t, db.User{}, util.SelectParams{
			Where: &util.EqClause{
				Field: "ID",
				Value: id,
			},
		})
		if err != nil {
			return err
		}
		return util.DeleteItem(t, i)
	})
	if err == util.ErrRowCount {
		http.Error(w, statusObjectNotFound, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, nil)
}
