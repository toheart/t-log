package note

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// SaveNote appends a note to the file at RootPath/YYYY/MM/YYYY-MM-DD.md
func SaveNote(rootPath, content string) error {
	if content == "" {
		return nil
	}

	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("2006-01-02")

	// Directory structure: RootPath/YYYY/MM
	dirPath := filepath.Join(rootPath, year, month)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filename := fmt.Sprintf("%s.md", day)
	filePath := filepath.Join(dirPath, filename)

	// Format: - [HH:MM] Content
	timeStr := now.Format("15:04")
	line := fmt.Sprintf("- [%s] %s\n", timeStr, content)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		return fmt.Errorf("failed to write note: %w", err)
	}

	return nil
}

// GetRecentNotes reads notes from the last n days
func GetRecentNotes(rootPath string, n int) ([]NoteEntry, error) {
	var entries []NoteEntry
	now := time.Now()

	// Check last N days
	for i := 0; i < n; i++ {
		date := now.AddDate(0, 0, -i)
		year := date.Format("2006")
		month := date.Format("01")
		day := date.Format("2006-01-02")

		filename := fmt.Sprintf("%s.md", day)
		filePath := filepath.Join(rootPath, year, month, filename)

		dayEntries, err := parseNoteFile(filePath, day)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Skip if file doesn't exist
			}
			// Log error but continue? Or return?
			// For MVP let's skip unreadable files
			continue
		}

		// Prepend to keep newest overall first, but parseNoteFile returns file order (oldest to newest in file).
		// We want the Timeline to usually show newest at top.
		// Let's append all and sort later or handle display order in frontend.
		// Usually daily logs are chronological.
		// If we iterate days backwards (Today, Yesterday...), and parse file (09:00, 10:00...),
		// We get: [Today 09:00, Today 10:00, Yesterday 09:00...]
		// If we want strictly reverse chronological:
		// Reverse dayEntries then append.
		for j := len(dayEntries) - 1; j >= 0; j-- {
			entries = append(entries, dayEntries[j])
		}
	}

	return entries, nil
}

var noteLineRegex = regexp.MustCompile(`^- \[(\d{2}:\d{2})\] (.*)$`)

func parseNoteFile(path, date string) ([]NoteEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []NoteEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		matches := noteLineRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			entries = append(entries, NoteEntry{
				Timestamp: matches[1],
				Content:   matches[2],
				Date:      date,
				RawLine:   line,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// OpenDailyNote opens the daily note file in the system default editor
// Implements US3 logic
func OpenDailyNote(rootPath string) error {
	// We need to find today's file. If it doesn't exist, we should probably create it first?
	// Or just open the directory if file missing?
	// Spec says "launch current file". Let's ensure it exists.

	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("2006-01-02")

	dirPath := filepath.Join(rootPath, year, month)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.md", day)
	filePath := filepath.Join(dirPath, filename)

	// Ensure file exists
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	f.Close()

	// Use 'open' command equivalent
	return openFileInOS(filePath)
}
