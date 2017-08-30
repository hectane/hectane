package auth

import (
	"net/http"

	"github.com/hectane/hectane/db"
)

// Login attempts to login the specified user.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request, username, password string) (*db.User, error) {
	u := &db.User{}
	if err := db.C.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	if err := u.Authenticate(password); err != nil {
		return nil, err
	}
	session, _ := a.store.Get(r, sessionName)
	session.Values[sessionUserID] = u.ID
	session.Save(r, w)
	return u, nil
}
