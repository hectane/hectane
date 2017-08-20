package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/db/util"
)

const (
	statusInvalidFolderName = "invalid folder name"
)

func (s *Server) folders(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*db.User)
	i, err := util.SelectItems(db.Token, db.Folder{}, util.SelectParams{
		Where: &util.EqClause{
			Field: "UserID",
			Value: u.ID,
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
		u = r.Context().Value(contextUser).(*db.User)
		p = r.Context().Value(contextParams).(*newFolderParams)
	)
	if len(p.Name) == 0 || len(p.Name) > 40 {
		http.Error(w, statusInvalidFolderName, http.StatusBadRequest)
		return
	}
	f := &db.Folder{
		Name:   p.Name,
		UserID: u.ID,
	}
	if err := util.InsertItem(db.Token, f); err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, f)
}

func (s *Server) deleteFolder(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*db.User)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := db.Token.Transaction(func(t *util.Token) error {
		i, err := util.SelectItem(t, db.Folder{}, util.SelectParams{
			Where: &util.AndClause{
				&util.EqClause{
					Field: "ID",
					Value: id,
				},
				&util.EqClause{
					Field: "UserID",
					Value: u.ID,
				},
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
