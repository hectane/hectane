package db

import (
	"database/sql"
)

// Token abstracts access to the database, allowing models to use both a
// database connection and a transaction in the same way.
type Token struct {
	tx *sql.Tx
}

func (t *Token) query(query string, args ...interface{}) (*sql.Rows, error) {
	if t.tx != nil {
		return t.tx.Query(query, args...)
	}
	return db.Query(query, args...)
}

func (t *Token) queryRow(query string, args ...interface{}) *sql.Row {
	if t.tx != nil {
		return t.tx.QueryRow(query, args...)
	}
	return db.QueryRow(query, args...)
}

func (t *Token) exec(query string, args ...interface{}) (sql.Result, error) {
	if t.tx != nil {
		return t.tx.Exec(query, args...)
	}
	return db.Exec(query, args...)
}

// DefaultToken may be passed to all database functions and methods when a
// transaction is not required.
var DefaultToken = &Token{}

// Transaction begins a new transaction and passes it to the provided callback.
// If no error is returned, Commit() is invoked - otherwise, Rollback(). The
// error is returned.
func Transaction(f func(*Token) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	t := &Token{tx: tx}
	if err := f(t); err != nil {
		t.tx.Rollback()
		return err
	}
	t.tx.Commit()
	return nil
}
