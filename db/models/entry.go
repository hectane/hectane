package models

import (
	"time"

	"github.com/hectane/hectane/db/util"
	"github.com/sirupsen/logrus"
)

const Entries = "entries"

// Entry represents a logfile entry.
type Entry struct {
	ID      int       `json:"id"`
	Context string    `json:"context"`
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

func MigrateEntries(t *util.Token) error {
	_, err := t.Exec(
		`
CREATE TABLE IF NOT EXISTS Entry (
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
	_, err = t.Exec(
		`
CREATE INDEX IF NOT EXISTS entries_context ON Entries (context)
		`,
	)
	if err != nil {
		return err
	}
	_, err = t.Exec(
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
