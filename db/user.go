package db

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// User represents an individual user within the system that can login, send,
// and receive emails.
type User struct {
	ID       int64
	Username string `gorm:"type:varchar(40);not null;unique_index"`
	Password string `gorm:"type:varchar(80);not null"`
	IsAdmin  bool
}

// MarshalJSON writes the content of the struct to JSON data.
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Username string `json:"username"`
		IsAdmin  bool   `json:"is-admin"`
	}{
		Username: u.Username,
		IsAdmin:  u.IsAdmin,
	})
}

// UnmarshalJSON reads JSON data into the struct.
func (u *User) UnmarshalJSON(data []byte) error {
	d := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is-admin"`
	}{}
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}
	u.Username = d.Username
	u.IsAdmin = d.IsAdmin
	if len(d.Password) != 0 {
		return u.SetPassword(d.Password)
	}
	return nil
}

// GetID retrieves the ID of the user.
func (u *User) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}

// SetID sets the ID for the user.
func (u *User) SetID(id string) error {
	u.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
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
