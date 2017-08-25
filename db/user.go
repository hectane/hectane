package db

import (
	"encoding/base64"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User represents an individual user within the system that can login, send,
// and receive emails.
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(40);not null;unique_index"`
	Password string `gorm:"type:varchar(80);not null"`
	IsAdmin  bool
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
