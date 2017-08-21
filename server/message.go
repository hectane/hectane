package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/db/models"
	"github.com/hectane/hectane/db/util"
)

func (s *Server) messages(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*models.User)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	i, err := util.SelectItems(db.Token, models.Message{}, util.SelectParams{
		Where: &util.AndClause{
			&util.EqClause{
				Field: "UserID",
				Value: u.ID,
			},
			&util.EqClause{
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
