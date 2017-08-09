package db

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User represents an individual user within the system that can login, send,
// and receive emails.
type User struct {
	ID       int
	Username string
	Password string
	IsAdmin  bool
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

// FindUser attempts to find a user where the specified field matches the
// specified value.
func FindUser(t *Token, field, value string) (*User, error) {
	u := &User{}
	err := t.queryRow(
		fmt.Sprintf(
			`
SELECT ID, Username, Password, IsAdmin
FROM Users WHERE %s = $1
            `,
			field,
		),
		value,
	).Scan(&u.ID, &u.Username, &u.Password, &u.IsAdmin)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Save persists changes to the user. If ID is set to zero, a new user is
// created and its ID updated.
func (u *User) Save(t *Token) error {
	if u.ID == 0 {
		var id int
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
		u.ID = id
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
