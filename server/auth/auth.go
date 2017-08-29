package auth

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/hectane/hectane/db"
)

const (
	sessionName   = "session"
	sessionUserID = "user_id"

	User = "user"
)

// Auth confirms that a user is logged in and stores the User instance in the
// session.
type Auth struct {
	handler http.Handler
	store   *sessions.CookieStore
}

// New creates a new auth handler that wraps the provided handler and manages
// session data.
func New(handler http.Handler, secretKey string) *Auth {
	return &Auth{
		handler: handler,
		store:   sessions.NewCookieStore([]byte(secretKey)),
	}
}

func (a *Auth) loadUser(r *http.Request) *db.User {
	session, _ := a.store.Get(r, sessionName)
	v, _ := session.Values[sessionUserID]
	if v == "" {
		return nil
	}
	u := &db.User{}
	if err := db.C.First(u, v).Error; err != nil {
		return nil
	}
	return u
}

// ServeHTTP ensures the user is logged in and stores their information in the
// request context.
func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := a.loadUser(r)
	if u == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	ctx := context.WithValue(r.Context(), User, u)
	a.handler.ServeHTTP(w, r.WithContext(ctx))
}
