// +build windows

package exec

import (
	"golang.org/x/sys/windows/svc/mgr"

	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	serviceName = "Hectane"
	displayName = "Hectane"
	description = "Lightweight SMTP client"
)

// Retrieves the full path to the current executable.
func exePath() (string, error) {
	if l, err := syscall.LoadDLL("kernel32.dll"); err == nil {
		defer l.Release()
		if p, err := l.FindProc("GetModuleFileNameW"); err == nil {
			b := make([]uint16, syscall.MAX_PATH)
			if ret, _, err := p.Call(
				0,
				uintptr(unsafe.Pointer(&b[0])),
				uintptr(len(b)),
			); ret != 0 {
				return string(utf16.Decode(b[:ret])), nil
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

// Installs the Hectane service using the service manager. If the service is
// already registered, the function immediately returns.
func installService(configFile string) error {
	if m, err := mgr.Connect(); err == nil {
		defer m.Disconnect()
		if s, err := m.OpenService(serviceName); err == nil {
			s.Close()
			return nil
		}
		if p, err := exePath(); err == nil {
			if s, err := m.CreateService(serviceName, p, mgr.Config{
				StartType:      mgr.StartAutomatic,
				BinaryPathName: fmt.Sprintf("\"%s\" -f \"%s\"", p, configFile),
				DisplayName:    displayName,
				Description:    description,
			}); err == nil {
				s.Close()
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

// Uninstalls the service using the service manager.
func uninstallService() error {
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
