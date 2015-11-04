package exec

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

const (
	ServiceName = "Hectane"
)

// A service must implement the svc.Handler interface.
type service struct{}

// Run the service, responding to control commands as necessary.
func (s *service) Execute(args []string, chChan <-chan svc.ChangeRequest, stChan chan<- svc.Status) (bool, uint32) {
	logrus.Debug("service started")
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
	logrus.Debug("service stopped")
	return false, 0
}

// If the application is running in an interactive session, run until
// terminated. Otherwise, run the application as a Windows service.
func Exec() error {
	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		return err
	}
	if !isInteractive {
		return svc.Run(ServiceName, &service{})
	} else {
		execSignal()
		return nil
	}
}
