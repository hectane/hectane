package exec

import (
	"golang.org/x/sys/windows/svc"
)

const (
	ServiceName = "Hectane"
)

// A service must implement the svc.Handler interface.
type service struct{}

// Run the service, responding to control commands as necessary.
func (s *service) Execute(args []string, chChan <-chan svc.ChangeRequest, stChan chan<- svc.Status) (bool, uint32) {
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
	return false, 0
}

// Operate as a Windows service.
func execService() error {
	return svc.Run(ServiceName, &service{})
}
