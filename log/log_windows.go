package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/hectane/hectane/exec"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"

	"fmt"
)

// Write log messages to the Windows event log.
type eventLogHook struct {
	log *eventlog.Log
}

// Create a new hook for the event log.
func newEventLogHook() (*eventLogHook, error) {
	e, err := eventlog.Open(exec.ServiceName)
	if err != nil {
		return nil, err
	}
	return &eventLogHook{
		log: e,
	}, nil
}

// Indicate which event levels should be logged.
func (e *eventLogHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
	}
}

// Log the specified entry to the event log. Add the context to the message.
func (e *eventLogHook) Fire(entry *logrus.Entry) error {
	msg := entry.Message
	if c, ok := entry.Data["context"]; ok {
		msg = fmt.Sprintf("[%s] %s", c, msg)
	}
	switch entry.Level {
	case logrus.InfoLevel:
		return e.log.Info(1, msg)
	case logrus.WarnLevel:
		return e.log.Warning(1, msg)
	case logrus.ErrorLevel:
		return e.log.Error(1, msg)
	default:
		return nil
	}
}

// Close the event log.
func (e *eventLogHook) close() {
	e.log.Close()
}

var hook *eventLogHook

func platformInit(config *Config) error {
	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		return err
	}
	if !isInteractive {
		h, err := newEventLogHook()
		if err != nil {
			return err
		}
		hook = h
		logrus.AddHook(hook)
	}
	return nil
}

func platformCleanup() {
	if hook != nil {
		hook.close()
	}
}
