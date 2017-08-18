package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hectane/hectane/db"
)

const (
	statusInvalidFolderName = "invalid folder name"
)

func (s *Server) folders(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*db.User)
	f, err := db.Folders(db.DefaultToken, u.ID)
	if err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, map[string]interface{}{
		"folders": f,
	})
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
	if err := f.Save(db.DefaultToken); err != nil {
		http.Error(w, statusDatabaseError, http.StatusInternalServerError)
		return
	}
	s.writeJson(w, nil)
}

func (s *Server) deleteFolder(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(contextUser).(*db.User)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := db.Transaction(func(t *db.Token) error {
		f, err := db.FindFolder(t, id, u.ID)
		if err != nil {
			return err
		}
		return f.Delete(t)
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
