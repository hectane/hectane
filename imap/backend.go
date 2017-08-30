package imap

import (
	"github.com/emersion/go-imap/backend"
	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/storage"
)

// dbBackend authenticates users by using the database.
type dbBackend struct {
	storage *storage.Storage
}

// Login determines if a user is authorized for access. If so, a user instance
// is returned, allowing access to all mailbox operations.
func (d *dbBackend) Login(username, password string) (backend.User, error) {
	u := &db.User{}
	if err := db.C.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, backend.ErrInvalidCredentials
	}
	if err := u.Authenticate(password); err != nil {
		return nil, backend.ErrInvalidCredentials
	}
	return &user{user: u}, nil
}
