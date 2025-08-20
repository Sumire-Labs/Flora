package commands

import "github.com/bwmarrin/discordgo"

// Command defines the structure for a slash command.
type Command struct {
	Def     *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Manager holds all the registered commands.
type Manager struct {
	Commands map[string]*Command
}

// NewManager creates a new command manager and registers all commands.
func NewManager() *Manager {
	m := &Manager{
		Commands: make(map[string]*Command),
	}

	// Register all commands here
	m.Add(PingCommand)
	m.Add(HelpCommand)

	return m
}

// Add adds a new command to the manager.
func (m *Manager) Add(cmd *Command) {
	m.Commands[cmd.Def.Name] = cmd
}

// Get retrieves a command by its name.
func (m *Manager) Get(name string) (*Command, bool) {
	cmd, exists := m.Commands[name]
	return cmd, exists
}

// RegisterAll registers all added commands globally with Discord.
func (m *Manager) RegisterAll(s *discordgo.Session) error {
	for _, cmd := range m.Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd.Def)
		if err != nil {
			return err
		}
	}
	return nil
}
