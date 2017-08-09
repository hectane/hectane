package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB

	ErrRowCount = errors.New("exactly one row must be returned")
)

// Connect establishes a connection to the PostgreSQL database used for all SQL
// queries. This function should be called before using any other types or
// functions in the package.
func Connect(name, user, password, host string, port int) error {
	c, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"dbname=%s user=%s password=%s host=%s port=%d",
			name,
			user,
			password,
			host,
			port,
		),
	)
	if err != nil {
		return err
	}
	db = c
	return nil
}
