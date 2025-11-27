package main

import (
	"context"
	"fmt"
	"t-log/internal/config"
	"t-log/internal/note"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
)

// App struct
type App struct {
	ctx    context.Context
	config *config.AppConfig
	hk     *hotkey.Hotkey
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		a.config = config.DefaultConfig()
	} else {
		a.config = cfg
	}

	// Register global hotkey
	// Default is Alt+Space.
	// Note: A robust implementation would parse a.config.Hotkey string.
	// For MVP/POC we hardcode Alt+Space or use a simple parser if needed.
	// Here we assume Alt+Space for simplicity as per MVP spec default.
	// Changing default to Ctrl+Alt+Space to avoid conflicts
	a.hk = hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt}, hotkey.KeySpace)
	if err := a.hk.Register(); err != nil {
		fmt.Printf("Failed to register hotkey: %v\n", err)
		// Consider showing a dialog to user
	}

	// Start listening for hotkey events in a goroutine
	go func() {
		for range a.hk.Keydown() {
			runtime.WindowShow(a.ctx)
			// Force focus - WindowShow typically does this, but AlwaysOnTop ensures it pops over other apps
			runtime.WindowSetAlwaysOnTop(a.ctx, true)
			// We don't immediately disable AlwaysOnTop to avoid focus flicker
		}
	}()
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.hk != nil {
		a.hk.Unregister()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetConfig returns the current configuration
func (a *App) GetConfig() *config.AppConfig {
	return a.config
}

// SaveNote appends a new note to today's markdown file
func (a *App) SaveNote(content string) error {
	return note.SaveNote(a.config.RootPath, content)
}

// GetRecentNotes reads and parses notes from the last N days
func (a *App) GetRecentNotes() []note.NoteEntry {
	entries, err := note.GetRecentNotes(a.config.RootPath, a.config.HistoryDays)
	if err != nil {
		fmt.Printf("Error getting recent notes: %v\n", err)
		return []note.NoteEntry{}
	}
	return entries
}

// OpenDailyNote opens the current day's markdown file in the system default editor
func (a *App) OpenDailyNote() error {
	return note.OpenDailyNote(a.config.RootPath)
}

// HideWindow hides the application window
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
	// Reset AlwaysOnTop when hidden so it doesn't interfere next time or if logic changes
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
}
