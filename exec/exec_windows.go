package exec

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"

	"fmt"
)

const (
	ServiceName = "Hectane"
)

// Write log messages to the Windows event log.
type eventLogHook struct {
	log *eventlog.Log
}

// Create a new hook for the event log.
func newEventLogHook() (*eventLogHook, error) {
	e, err := eventlog.Open(ServiceName)
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

// A service must implement the svc.Handler interface.
type service struct{}

// Run the service, responding to control commands as necessary.
func (s *service) Execute(args []string, chChan <-chan svc.ChangeRequest, stChan chan<- svc.Status) (bool, uint32) {
	logrus.Info("service started")
	stChan <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}
loop:
	for {
		c := <-chChan
		switch c.Cmd {
		case svc.Interrogate:
			stChan <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			stChan <- svc.Status{State: svc.StopPending}
			break loop
		}
	}
	logrus.Info("service stopped")
	return false, 0
}

var (
	isInteractive bool
	hook          *eventLogHook
)

// Determine if the application is running interactively or as a service. If
// running as a service, add the event log hook.
func platformInit() error {
	i, err := svc.IsAnInteractiveSession()
	if err != nil {
		return err
	}
	isInteractive = i
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

// If the application is running in an interactive session, run until
// terminated. Otherwise, run the application as a Windows service.
func Exec() error {
	if isInteractive {
		execSignal()
		return nil
	} else {
		return svc.Run(ServiceName, &service{})
	}
}

// If the event log was opened, close it.
func Cleanup() {
	if hook != nil {
		hook.close()
	}
}
