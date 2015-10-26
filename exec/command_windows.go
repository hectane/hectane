package exec

// Add commands for interacting with Windows services.
func initPlatformCommands() {
	Commands[InstallCommand.Name] = InstallCommand
	Commands[RemoveCommand.Name] = RemoveCommand
}
