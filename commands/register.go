package commands

// RegisterAllCommands adds all command definitions to the manager and returns it.
func RegisterAllCommands(m *Manager) *Manager {
	m.Add(PingCommand)
	m.Add(HelpCommand)
	// Add new commands here in the future
	return m
}
