package exec

import (
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"encoding/json"
	"fmt"
	"os"
	"path"
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
	l, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return "", err
	}
	defer l.Release()
	p, err := l.FindProc("GetModuleFileNameW")
	if err != nil {
		return "", err
	}
	b := make([]uint16, syscall.MAX_PATH)
	ret, _, err := p.Call(
		0,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(b)),
	)
	if ret != 0 {
		return string(utf16.Decode(b[:ret])), nil
	} else {
		return "", err
	}
}

// Determine the path to the configuration file and if it does not exist,
// create one with the default configuration.
func saveConfig(exePath string, cfg *Config) (string, error) {
	cfgPath := path.Join(path.Dir(exePath), "config.json")
	w, err := os.OpenFile(cfgPath, os.O_WRONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			return cfgPath, nil
		} else {
			return "", err
		}
	}
	defer w.Close()
	if err := json.NewEncoder(w).Encode(cfg); err != nil {
		return "", err
	}
	return cfgPath, nil
}

// Run the specified command on the service.
func serviceCommand(cmd string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()
	switch cmd {
	case "start":
		return s.Start()
	case "stop":
		_, err := s.Control(svc.Stop)
		return err
	case "remove":
		return s.Delete()
	}
	return nil
}

// Connect to the service manager and install the service.
var InstallCommand = &Command{
	Name:        "install",
	Description: "install the service (Windows only)",
	Exec: func(cfg *Config) error {
		m, err := mgr.Connect()
		if err != nil {
			return err
		}
		defer m.Disconnect()
		p, err := exePath()
		if err != nil {
			return err
		}
		c, err := saveConfig(p, cfg)
		if err != nil {
			return err
		}
		s, err := m.CreateService(serviceName, p, mgr.Config{
			StartType:      mgr.StartAutomatic,
			BinaryPathName: fmt.Sprintf("\"%s\" -f \"%s\"", p, c),
			DisplayName:    displayName,
			Description:    description,
		})
		if err != nil {
			return err
		}
		s.Close()
		return nil
	},
}

// Start the service.
var StartCommand = &Command{
	Name:        "start",
	Description: "start the service (Windows only)",
	Exec: func(cfg *Config) error {
		return serviceCommand("start")
	},
}

var StopCommand = &Command{
	Name:        "stop",
	Description: "stop the service (Windows only)",
	Exec: func(cfg *Config) error {
		return serviceCommand("stop")
	},
}

// Remove the service.
var RemoveCommand = &Command{
	Name:        "remove",
	Description: "remove the service (Windows only)",
	Exec: func(cfg *Config) error {
		return serviceCommand("remove")
	},
}
