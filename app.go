package main

import (
	"context"
	"fmt"
	"t-log/internal/attachment"
	"t-log/internal/command"
	"t-log/internal/config"
	hk "t-log/internal/hotkey"
	"t-log/internal/note"

	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
)

// App struct
type App struct {
	ctx         context.Context
	config      *config.AppConfig
	hk          *hotkey.Hotkey
	cmdRegistry *command.CommandRegistry
	attachMgr   *attachment.Manager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		cmdRegistry: command.NewRegistry(),
	}
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

	// Initialize managers
	a.attachMgr = attachment.NewManager(a.config)

	// Register global hotkey
	mods, key, err := hk.ParseHotkey(a.config.Hotkey)
	if err != nil {
		fmt.Printf("Error parsing hotkey '%s': %v. Using default Alt+Space.\n", a.config.Hotkey, err)
		mods = []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt}
		key = hotkey.KeySpace
	}

	a.hk = hotkey.New(mods, key)
	if err := a.hk.Register(); err != nil {
		fmt.Printf("Failed to register hotkey: %v\n", err)
	}

	// Start listening for hotkey events in a goroutine
	go func() {
		for range a.hk.Keydown() {
			// On Windows with Acrylic, resizing a hidden window or resizing immediately
			// after show can crash. The safest way is:
			// 1. Show the window (it might be wrong size)
			// 2. Wait a tiny bit (let DWM catch up) - handled by frontend event delay
			// 3. Emit event for frontend to focus input
			runtime.WindowShow(a.ctx)
			// Force restore to ensure it's not minimized
			if runtime.WindowIsMinimised(a.ctx) {
				runtime.WindowUnminimise(a.ctx)
			}
			// Flash Top: Set AlwaysOnTop to bring to front, then disable it
			// This allows the user to Alt-Tab away or click other windows later.
			runtime.WindowSetAlwaysOnTop(a.ctx, true)
			go func() {
				time.Sleep(100 * time.Millisecond) // Short delay to ensure it pops up
				runtime.WindowSetAlwaysOnTop(a.ctx, false)
			}()

			// Delay event emission slightly to ensure window is fully rendered
			// This helps with the "flash crash" on some Windows systems
			go func() {
				// Short sleep (e.g. 50ms) could be done here if needed,
				// but frontend timeout is usually enough.
				// Let's keep it immediate here but rely on frontend delay.
				runtime.EventsEmit(a.ctx, "app:reset")
			}()
		}
	}()

	// Register Commands
	a.registerCommands()
}

// registerCommands registers all available commands
func (a *App) registerCommands() {
	// Open Specific Date
	a.cmdRegistry.Register(command.Command{
		ID:          "cmd:open-date",
		Title:       "Open Date...",
		Description: "Select a specific date to open",
	}, func(args []string) error {
		// This command is handled by frontend to show date picker
		// But we register it so it appears in the list
		return nil
	})

	// Help
	a.cmdRegistry.Register(command.Command{
		ID:          "cmd:help",
		Title:       "Help",
		Description: "Open documentation in browser",
	}, func(args []string) error {
		runtime.BrowserOpenURL(a.ctx, "https://github.com/yourusername/t-log/blob/master/README.md") // Update URL as needed
		return nil
	})

	// Find (placeholder for discovery, logic handled in frontend/backend search)
	a.cmdRegistry.Register(command.Command{
		ID:          "cmd:find",
		Title:       "Find / Search",
		Description: "Search notes by keyword (Type 'find ')",
		Usage:       "find <keyword>",
	}, func(args []string) error {
		// Backend doesn't need to do anything for 'find' as it switches mode in frontend
		// But we register it so it shows up in the list
		return nil
	})

	// Settings
	a.cmdRegistry.Register(command.Command{
		ID:          "cmd:settings",
		Title:       "Settings",
		Description: "Open configuration settings",
	}, func(args []string) error {
		// The frontend CommandPalette will intercept this and handle it,
		// OR we can emit an event here if executed via backend logic.
		// Let's emit the event just in case.
		runtime.EventsEmit(a.ctx, "app:open-settings")
		return nil
	})
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

// UpdateConfig updates the application configuration
func (a *App) UpdateConfig(cfg config.AppConfig) error {
	a.config = &cfg
	// In a real app, we might want to re-init things if config changes (like hotkey)
	// For now, we just save it.
	// Hotkey update would require unregistering old and registering new, which is complex safely.
	// We'll assume restart required for hotkey for now, as per spec US3 test scenario "Restart".
	return config.SaveConfig(config.ConfigFileName, a.config)
}

// SelectRootPath opens a dialog to select the root path
func (a *App) SelectRootPath() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Notes Directory",
	})
	if err != nil {
		return "", err
	}
	return path, nil
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

// OpenDateNote opens the markdown file for a specific date (YYYY-MM-DD)
func (a *App) OpenDateNote(dateStr string) error {
	return note.OpenDateNote(a.config.RootPath, dateStr)
}

// OpenNoteAt opens a specific note file at a specific line number
func (a *App) OpenNoteAt(filePath string, lineNo int) error {
	return note.OpenNoteAt(filePath, lineNo)
}

// HideWindow hides the application window
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
	// Reset AlwaysOnTop when hidden so it doesn't interfere next time or if logic changes
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
}

// GetCommands returns all available commands
func (a *App) GetCommands() []command.Command {
	return a.cmdRegistry.GetCommands()
}

// ExecuteCommand executes a specific command by ID
func (a *App) ExecuteCommand(id string, args []string) error {
	return a.cmdRegistry.Execute(id, args)
}

// SearchNotes performs a text search across all notes
func (a *App) SearchNotes(query string) []note.SearchResult {
	results, err := note.SearchNotes(a.config.RootPath, query)
	if err != nil {
		fmt.Printf("Error searching notes: %v\n", err)
		return []note.SearchResult{}
	}
	return results
}

// UploadAttachment saves the provided content as a file in the attachment directory
func (a *App) UploadAttachment(content []byte, filename string) (string, error) {
	return a.attachMgr.SaveAttachment(content, filename)
}

// ListNoteDates returns a list of all available note dates
func (a *App) ListNoteDates() ([]string, error) {
	return note.ListNoteDates(a.config.RootPath)
}
