package cmd

import (
	"fmt"
	"os"
)

// Information about an application command.
type command struct {
	name        string
	description string
	exec        func() error
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
func Exec(name string) error {
	for _, c := range commands {
		if name == c.name {
			return c.exec()
		}
	}
	return fmt.Errorf("invalid command \"%s\"", name)
}
