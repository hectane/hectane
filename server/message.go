package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/db/models"
	"github.com/hectane/hectane/db/sql"
)

func (s *Server) messages(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*models.User)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	i, err := sql.SelectItems(db.Token, models.Message{}, sql.SelectParams{
		Where: &sql.AndClause{
			&sql.EqClause{
				Field: "UserID",
				Value: u.ID,
			},
			&sql.EqClause{
				Field: "FolderID",
				Value: id,
			},
		},
		OrderBy:   "Time",
		OrderDesc: true,
	})
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, i)
}
