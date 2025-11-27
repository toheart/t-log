package command

import (
	"fmt"
)

// Command represents an executable action in the palette
type Command struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

// Handler is the function that executes the command
type Handler func(args []string) error

// CommandRegistry manages available commands
type CommandRegistry struct {
	commands map[string]Command
	handlers map[string]Handler
}

// NewRegistry creates a new command registry
func NewRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
		handlers: make(map[string]Handler),
	}
}

// Register adds a command to the registry
func (r *CommandRegistry) Register(cmd Command, handler Handler) {
	r.commands[cmd.ID] = cmd
	r.handlers[cmd.ID] = handler
}

// GetCommands returns all registered commands
func (r *CommandRegistry) GetCommands() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// Execute runs a command by ID
func (r *CommandRegistry) Execute(id string, args []string) error {
	handler, ok := r.handlers[id]
	if !ok {
		return fmt.Errorf("command not found: %s", id)
	}
	return handler(args)
}
