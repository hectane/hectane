package auth

import (
	"net/http"
)

// Logout destroys the current session.
func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := a.store.Get(r, sessionName)
	session.Values[sessionUserID] = ""
	session.Save(r, w)
}
