package db

import (
	"database/sql"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User represents an individual user within the system that can login, send,
// and receive emails.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}

func migrateUsersTable(t *Token) error {
	_, err := t.exec(
		`
CREATE TABLE IF NOT EXISTS Users (
    ID       SERIAL PRIMARY KEY,
    Username VARCHAR(40) NOT NULL UNIQUE,
    Password VARCHAR(80) NOT NULL,
    IsAdmin  BOOLEAN
)
        `,
	)
	return err
}

func rowsToUsers(r *sql.Rows) ([]*User, error) {
	users := make([]*User, 0, 1)
	for r.Next() {
		u := &User{}
		if err := r.Scan(&u.ID, &u.Username, &u.Password, &u.IsAdmin); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Users retrieves a list all users in the database.
func Users(t *Token) ([]*User, error) {
	r, err := t.query(
		`
SELECT ID, Username, Password, IsAdmin
FROM Users ORDER BY Username
        `,
	)
	if err != nil {
		return nil, err
	}
	return rowsToUsers(r)
}

// FindUser attempts to find a user where the specified field matches the
// specified value. Exactly one row must be returned.
func FindUser(t *Token, field string, value interface{}) (*User, error) {
	r, err := FindUsers(t, field, value)
	if err != nil {
		return nil, err
	}
	if len(r) != 1 {
		return nil, ErrRowCount
	}
	return r[0], nil
}

// FindUsers attempts to retrieve all users where the specified field matches
// the specified value.
func FindUsers(t *Token, field string, value interface{}) ([]*User, error) {
	r, err := t.query(
		fmt.Sprintf(
			`
SELECT ID, Username, Password, IsAdmin
FROM Users WHERE %s = $1 ORDER BY Username
            `,
			field,
		),
		value,
	)
	if err != nil {
		return nil, err
	}
	return rowsToUsers(r)
}

// Save persists changes to the user. If ID is set to zero, a new user is
// created and its ID updated.
func (u *User) Save(t *Token) error {
	if u.ID == 0 {
		err := t.queryRow(
			`
INSERT INTO Users (Username, Password, IsAdmin)
VALUES ($1, $2, $3) RETURNING ID
            `,
			u.Username,
			u.Password,
			u.IsAdmin,
		).Scan(&u.ID)
		if err != nil {
			return err
		}
		return nil
	} else {
		_, err := t.exec(
			`
UPDATE Users SET Username=$1, Password=$2, IsAdmin=$3
WHERE ID = $4
            `,
			u.Username,
			u.Password,
			u.IsAdmin,
			u.ID,
		)
		return err
	}
}

// Delete the user from the database.
func (u *User) Delete(t *Token) error {
	_, err := t.exec(
		`
DELETE FROM Users WHERE ID = $1
        `,
		u.ID,
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
