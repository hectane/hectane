package exec

// Add commands for interacting with Windows services.
func initPlatformCommands() {
	Commands[InstallCommand.Name] = InstallCommand
	Commands[StartCommand.Name] = StartCommand
	Commands[StopCommand.Name] = StopCommand
	Commands[RemoveCommand.Name] = RemoveCommand
}
