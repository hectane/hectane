package db

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Entry represents a logfile entry.
type Entry struct {
	ID      int       `json:"id"`
	Context string    `json:"context"`
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

func migrateEntriesTable(t *Token) error {
	_, err := t.exec(
		`
CREATE TABLE IF NOT EXISTS Entries (
    ID      SERIAL PRIMARY KEY,
    Context VARCHAR(40) NOT NULL,
    Time    TIMESTAMPTZ NOT NULL,
    Level   VARCHAR(40) NOT NULL,
    Message TEXT
)
        `,
	)
	if err != nil {
		return err
	}
	_, err = t.exec(
		`
CREATE INDEX IF NOT EXISTS entries_context ON Entries (context)
        `,
	)
	if err != nil {
		return err
	}
	_, err = t.exec(
		`
CREATE INDEX IF NOT EXISTS entries_level ON Entries (level)
        `,
	)
	return err
}

// NewEntry creates a new unsaved Entry from a logrus.Entry.
func NewEntry(entry *logrus.Entry) *Entry {
	context, ok := entry.Data["context"].(string)
	if !ok {
		context = "unknown"
	}
	level := "unknown"
	switch entry.Level {
	case logrus.PanicLevel:
		level = "panic"
	case logrus.FatalLevel:
		level = "fatal"
	case logrus.ErrorLevel:
		level = "error"
	case logrus.WarnLevel:
		level = "warn"
	case logrus.InfoLevel:
		level = "info"
	case logrus.DebugLevel:
		level = "debug"
	}
	return &Entry{
		Context: context,
		Time:    entry.Time,
		Level:   level,
		Message: entry.Message,
	}
}

// ClearEntries removes all entries from the database.
func ClearEntries(t *Token) error {
	_, err := t.exec(
		`
TRUNCATE TABLE Entries
        `,
	)
	return err
}

// Save persists the entry in the database.
func (e *Entry) Save(t *Token) error {
	return t.queryRow(
		`
INSERT INTO Entries (Context, Time, Level, Message)
VALUES ($1, $2, $3, $4) RETURNING ID
        `,
		e.Context,
		e.Time,
		e.Level,
		e.Message,
	).Scan(&e.ID)
}
