package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = "config.json"

// LoadConfig loads the configuration from config.json or creates a default one
func LoadConfig() (*AppConfig, error) {
	configPath := ConfigFileName // Currently relative to CWD, can be improved to AppData

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultCfg := DefaultConfig()

		// Try to resolve default root path to user home
		home, err := os.UserHomeDir()
		if err == nil {
			defaultCfg.RootPath = filepath.Join(home, "QuickNotes")
		}

		if err := SaveConfig(configPath, defaultCfg); err != nil {
			return nil, err
		}
		return defaultCfg, nil
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to disk
func SaveConfig(path string, cfg *AppConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
