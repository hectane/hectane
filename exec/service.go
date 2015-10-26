// +build windows

package exec

import (
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	serviceName = "Hectane"
	displayName = "Hectane"
	description = "Lightweight SMTP client"
)

// Remove the service using the service manager.
func removeService() error {
	if m, err := mgr.Connect(); err == nil {
		defer m.Disconnect()
		if s, err := m.OpenService(serviceName); err == nil {
			defer s.Close()
			return s.Delete()
		} else {
			return err
		}
	} else {
		return err
	}
}

// Run the service.
func execService() {
	// TODO
}
