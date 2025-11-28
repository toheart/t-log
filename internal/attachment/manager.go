package attachment

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"t-log/internal/config"
)

type Manager struct {
	config *config.AppConfig
}

func NewManager(cfg *config.AppConfig) *Manager {
	return &Manager{
		config: cfg,
	}
}

// EnsureDir creates the attachment directory for the current month if it doesn't exist
func (m *Manager) EnsureDir() (string, error) {
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")

	// {RootPath}/{YYYY}/{MM}/Attachment/
	path := filepath.Join(m.config.RootPath, year, month, "Attachment")

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("failed to create attachment directory: %w", err)
	}

	return path, nil
}

func sanitizeFilename(name string) string {
	// Replace illegal characters with _
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	return re.ReplaceAllString(name, "_")
}

// SaveAttachment saves the content to a file and returns the web-accessible path
func (m *Manager) SaveAttachment(content []byte, filename string) (string, error) {
	dir, err := m.EnsureDir()
	if err != nil {
		return "", err
	}

	// Generate unique filename: {Timestamp}_{OriginalName}
	timestamp := time.Now().UnixMilli()
	sanitized := sanitizeFilename(filename)
	newFilename := fmt.Sprintf("%d_%s", timestamp, sanitized)
	fullPath := filepath.Join(dir, newFilename)

	// Write file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write attachment: %w", err)
	}

	// Return path relative to RootPath/.. -> Web Path
	// The web handler will handle /attachments/ prefix and map it to RootPath
	// We return: /attachments/{YYYY}/{MM}/Attachment/{Filename}

	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")

	// URL Encode the filename part to handle spaces and special chars in URL
	encodedFilename := url.PathEscape(newFilename)

	// Using forward slashes for web URL
	webPath := fmt.Sprintf("/attachments/%s/%s/Attachment/%s", year, month, encodedFilename)

	return webPath, nil
}
