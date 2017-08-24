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
	statusInvalidFolderName = "invalid folder name"
)

func (s *Server) folders(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*models.User)
	i, err := sql.SelectItems(db.Token, models.Folder{}, sql.SelectParams{
		Where: &sql.ComparisonClause{
			Field:    "UserID",
			Operator: sql.OpEq,
			Value:    u.ID,
		},
		OrderBy: "Name",
	})
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, i)
}

type newFolderParams struct {
	Name string `json:"name"`
}

func (s *Server) newFolder(w http.ResponseWriter, r *http.Request) {
	var (
		u = r.Context().Value(contextUser).(*models.User)
		p = r.Context().Value(contextParams).(*newFolderParams)
	)
	if len(p.Name) == 0 || len(p.Name) > 40 {
		http.Error(w, statusInvalidFolderName, http.StatusBadRequest)
		return
	}
	f := &models.Folder{
		Name:   p.Name,
		UserID: u.ID,
	}
	if err := sql.InsertItem(db.Token, f); err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, f)
}

func (s *Server) deleteFolder(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*models.User)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := db.Token.Transaction(func(t *sql.Token) error {
		i, err := sql.SelectItem(t, models.Folder{}, sql.SelectParams{
			Where: &sql.AndClause{
				&sql.ComparisonClause{
					Field:    "ID",
					Operator: sql.OpEq,
					Value:    id,
				},
				&sql.ComparisonClause{
					Field:    "UserID",
					Operator: sql.OpEq,
					Value:    u.ID,
				},
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
