package util

import (
	"database/sql"
)

// Token abstracts access to the database, allowing models to use both a
// database connection and a transaction in the same way.
type Token struct {
	tx *sql.Tx
	db *sql.DB
}

func (t *Token) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if t.tx != nil {
		return t.tx.Query(query, args...)
	}
	return t.db.Query(query, args...)
}

func (t *Token) QueryRow(query string, args ...interface{}) *sql.Row {
	if t.tx != nil {
		return t.tx.QueryRow(query, args...)
	}
	return t.db.QueryRow(query, args...)
}

func (t *Token) Exec(query string, args ...interface{}) (sql.Result, error) {
	if t.tx != nil {
		return t.tx.Exec(query, args...)
	}
	return t.db.Exec(query, args...)
}

// NewToken creates a new token for the provided database connection.
func NewToken(db *sql.DB) *Token {
	return &Token{
		db: db,
	}
}

// Transaction begins a new transaction and passes it to the provided callback.
// If no error is returned, Commit() is invoked - otherwise, Rollback(). The
// error is returned.
func (t *Token) Transaction(f func(*Token) error) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	tt := &Token{tx: tx}
	if err := f(tt); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
