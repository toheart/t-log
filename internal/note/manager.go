package note

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
func GetRecentNotes(rootPath string, n int) ([]DailyNote, error) {
	now := time.Now()
	// Start date is Today - (n-1) days
	startDate := now.AddDate(0, 0, -(n - 1))
	return GetDailyNotes(rootPath, startDate.Format("2006-01-02"), now.Format("2006-01-02"))
}

// GetNotesByDateRange reads notes within a start and end date range (inclusive)
// start, end format: YYYY-MM-DD
func GetNotesByDateRange(rootPath, start, end string) ([]NoteEntry, error) {
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	var entries []NoteEntry

	// Iterate from end date down to start date to keep reverse chronological order by day
	// Or iterate start to end?
	// Existing GetRecentNotes does newest first (reverse chronological).
	// Let's stick to reverse chronological (newest notes first).

	current := endDate
	for !current.Before(startDate) {
		year := current.Format("2006")
		month := current.Format("01")
		day := current.Format("2006-01-02")

		filename := fmt.Sprintf("%s.md", day)
		filePath := filepath.Join(rootPath, year, month, filename)

		dayEntries, err := parseNoteFile(filePath, day)
		if err == nil {
			// Reverse dayEntries to have newest time first
			for j := len(dayEntries) - 1; j >= 0; j-- {
				entries = append(entries, dayEntries[j])
			}
		} else if !os.IsNotExist(err) {
			// Log error?
		}

		current = current.AddDate(0, 0, -1)
	}

	return entries, nil
}

// GetDailyNotes reads full file contents within a start and end date range (inclusive)
// Returns parsed DailyNote structs
func GetDailyNotes(rootPath, start, end string) ([]DailyNote, error) {
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	var notes []DailyNote

	// Iterate from end date down to start date (newest first)
	current := endDate
	for !current.Before(startDate) {
		year := current.Format("2006")
		month := current.Format("01")
		day := current.Format("2006-01-02")

		filename := fmt.Sprintf("%s.md", day)
		filePath := filepath.Join(rootPath, year, month, filename)

		contentBytes, err := os.ReadFile(filePath)
		if err == nil {
			notes = append(notes, DailyNote{
				Date:    day,
				Content: string(contentBytes),
			})
		} else if !os.IsNotExist(err) {
			// Log error?
		}

		current = current.AddDate(0, 0, -1)
	}

	return notes, nil
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

// ListNoteDates returns a list of all available note dates (YYYY-MM-DD)
func ListNoteDates(rootPath string) ([]string, error) {
	var dates []string

	// Walk the root path
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip hidden directories if any, or Attachment directory
			if d.Name() == "Attachment" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file matches YYYY-MM-DD.md
		name := d.Name()
		if strings.HasSuffix(name, ".md") {
			// Simple validation: length is 10+3=13 chars?
			// Better: regex check
			datePart := strings.TrimSuffix(name, ".md")
			if _, err := time.Parse("2006-01-02", datePart); err == nil {
				dates = append(dates, datePart)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort reverse chronological
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	return dates, nil
}
