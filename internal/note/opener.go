package note

import (
	"fmt"
	"os/exec"
	"runtime"
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
