package cmd

import (
	"github.com/hectane/hectane/cfg"
	"github.com/hectane/hectane/util"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"fmt"
	"path"
)

const (
	serviceName = "Hectane"
	displayName = "Hectane"
	description = "Lightweight SMTP client"
)

// Run the specified command on the service.
func serviceCommand(name string) error {
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
	switch name {
	case "remove":
		return s.Delete()
	case "start":
		return s.Start()
	case "stop":
		_, err := s.Control(svc.Stop)
		return err
	}
	return nil
}

// Connect to the service manager and install the service.
var installCommand = &command{
	name:        "install",
	description: "install the service (Windows only)",
	exec: func(config *cfg.Config) error {
		m, err := mgr.Connect()
		if err != nil {
			return err
		}
		defer m.Disconnect()
		exePath, err := util.Executable()
		if err != nil {
			return err
		}
		cfgPath := path.Join(path.Dir(exePath), "config.json")
		if err := config.Save(cfgPath); err != nil {
			return err
		}
		s, err := m.CreateService(serviceName, exePath, mgr.Config{
			StartType:      mgr.StartAutomatic,
			BinaryPathName: fmt.Sprintf("\"%s\" -f \"%s\"", exePath, cfgPath),
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
var startCommand = &command{
	name:        "start",
	description: "start the service (Windows only)",
	exec: func(config *cfg.Config) error {
		return serviceCommand("start")
	},
}

// Stop the service.
var stopCommand = &command{
	name:        "stop",
	description: "stop the service (Windows only)",
	exec: func(config *cfg.Config) error {
		return serviceCommand("stop")
	},
}

// Remove the service.
var removeCommand = &command{
	name:        "remove",
	description: "remove the service (Windows only)",
	exec: func(config *cfg.Config) error {
		return serviceCommand("remove")
	},
}
