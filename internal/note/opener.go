package note

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func openFileInOS(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

// OpenNoteAt opens a file at a specific line number using VS Code if available,
// otherwise falls back to system default opener (ignoring line number).
func OpenNoteAt(path string, lineNo int) error {
	// Check if VS Code ("code") is in PATH
	// This is a heuristic; if users use other editors, we might need configuration.
	// But VS Code is the most common "goto line" capable editor users likely use.
	codePath, err := exec.LookPath("code")
	if err == nil {
		// VS Code found, use "code -g file:line"
		arg := fmt.Sprintf("%s:%d", path, lineNo)
		cmd := exec.Command(codePath, "-g", arg)
		// On Windows, using just "code" might spawn a console window if not careful,
		// but LookPath + Command usually handles it fine for executables.
		// However, to be safe and non-blocking:
		if runtime.GOOS == "windows" {
			// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // requires syscall import
		}
		return cmd.Start()
	}

	// Fallback: just open the file
	return openFileInOS(path)
}

// OpenDateNote opens the markdown file for a specific date
// dateStr should be in "YYYY-MM-DD" format
func OpenDateNote(rootPath, dateStr string) error {
	// Validate date format
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
	}

	year := t.Format("2006")
	month := t.Format("01")
	filename := t.Format("2006-01-02") + ".md"

	// Path: root/2006/01/2006-01-02.md
	filePath := filepath.Join(rootPath, year, month, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		// Create empty file if not exists
		if f, err := os.Create(filePath); err == nil {
			f.Close()
		}
	}

	return openFileInOS(filePath)
}
