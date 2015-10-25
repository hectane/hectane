package exec

// Information about an application command.
type Command struct {
	Name        string
	Description string
	Exec        func(c *Config) error
}

// Map of available commands.
var Commands = make(map[string]*Command)

// Initialize the list of commands. The map is first populated with commands
// common to all platforms followed by platform-specific commands.
func InitCommands() {
	Commands["run"] = RunCommand
	initPlatformCommands()
}
