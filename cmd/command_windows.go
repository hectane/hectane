package cmd

// Initialize the commands available for the current platform.
func Init() {
	commands = []*command{
		installCommand,
		removeCommand,
		startCommand,
		stopCommand,
	}
}
