package server

import (
	"net/http"

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
	s.writeJson(w, f)
}
