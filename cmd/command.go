package cmd

import (
	"github.com/hectane/hectane/cfg"

	"fmt"
	"os"
)

// Information about an application command.
type command struct {
	name        string
	description string
	exec        func(config *cfg.Config) error
}

// List of valid commands.
var commands []*command

// Display a list of valid commands.
func Print() {
	for _, c := range commands {
		fmt.Fprintf(os.Stderr, "  %s\n\t%s\n", c.name, c.description)
	}
}

// Execute the specified command if available.
func Exec(name string, config *cfg.Config) error {
	for _, c := range commands {
		if name == c.name {
			return c.exec(config)
		}
	}
	return fmt.Errorf("invalid command \"%s\"", name)
}
