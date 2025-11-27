package command

// InitRegistry initializes the registry with default commands
func InitRegistry() *CommandRegistry {
	registry := NewRegistry()

	// Register default commands here
	// Note: Handlers will be injected or bound in main/app setup,
	// but for purely static commands or those needing basic setup we define here.
	// For commands needing App context (like OpenDailyNote), we might need a better injection strategy.
	// For now, we'll just return the registry and let App register handlers.

	return registry
}
