package models

import (
	"encoding/base64"

	"github.com/hectane/hectane/db/util"
	"golang.org/x/crypto/bcrypt"
)

const Users = "users"

// User represents an individual user within the system that can login, send,
// and receive emails.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}

func MigrateUsers(t *util.Token) error {
	_, err := t.Exec(
		`
CREATE TABLE IF NOT EXISTS User_ (
	ID       SERIAL PRIMARY KEY,
	Username VARCHAR(40) NOT NULL UNIQUE,
	Password VARCHAR(80) NOT NULL,
	IsAdmin  BOOLEAN
)
		`,
	)
	return err
}

// Authenticate hashes the provided password and compares it to the value
// stored in the database.
func (u *User) Authenticate(password string) error {
	h, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(h, []byte(password))
}

// SetPassword salts and hashes the user's password. It does not update the
// user's row in the database.
func (u *User) SetPassword(password string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	u.Password = base64.StdEncoding.EncodeToString(h)
	return nil
}
