package config

// AppConfig represents the application configuration
type AppConfig struct {
	RootPath    string `json:"root_path"`    // Root directory for notes
	Hotkey      string `json:"hotkey"`       // Global hotkey to toggle window
	HistoryDays int    `json:"history_days"` // Number of days to show in history
}

// DefaultConfig returns the default configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		RootPath:    "QuickNotes", // Will be relative to user home if not absolute
		Hotkey:      "Ctrl+Alt+Space",
		HistoryDays: 3,
	}
}
