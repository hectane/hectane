package exec

import (
	"golang.org/x/sys/windows/svc"
)

// If the application is running in an interactive session, run until
// terminated. Otherwise, run the application as a Windows service.
func Exec() error {
	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		return err
	}
	if isIntSess {
		execSignal()
		return nil
	} else {
		return nil
	}
}
