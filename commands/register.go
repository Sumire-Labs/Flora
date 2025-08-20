package commands

// RegisterAllCommands adds all command definitions to the manager.
func RegisterAllCommands(m *Manager) {
	m.Add(PingCommand)
	m.Add(HelpCommand)
	// Add new commands here in the future
}
