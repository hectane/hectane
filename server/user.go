package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/db/models"
	"github.com/hectane/hectane/db/sql"
)

const (
	statusInvalidUsername = "invalid username"
)

func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	i, err := sql.SelectItems(db.Token, models.User{}, sql.SelectParams{
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
	u := &models.User{
		Username: p.Username,
		IsAdmin:  p.IsAdmin,
	}
	if err := u.SetPassword(p.Password); err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	if err := sql.InsertItem(db.Token, u); err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, u)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := db.Token.Transaction(func(t *sql.Token) error {
		i, err := sql.SelectItem(t, models.User{}, sql.SelectParams{
			Where: &sql.EqClause{
				Field: "ID",
				Value: id,
			},
		})
		if err != nil {
			return err
		}
		return sql.DeleteItem(t, i)
	})
	if err == sql.ErrRowCount {
		http.Error(w, statusObjectNotFound, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, nil)
}
