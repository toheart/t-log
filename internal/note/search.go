package note

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SearchResult represents a match found in the notes
type SearchResult struct {
	Content  string `json:"content"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	FilePath string `json:"filePath"`
	LineNo   int    `json:"lineNo"`
}

// SearchNotes scans all markdown files in rootPath for the query string (case-insensitive)
func SearchNotes(rootPath, query string) ([]SearchResult, error) {
	var results []SearchResult
	query = strings.ToLower(query)

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // or continue?
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		// Parse date from filename/path?
		// Structure: Root/YYYY/MM/YYYY-MM-DD.md
		// Simpler: Just read file and parse lines

		file, err := os.Open(path)
		if err != nil {
			return nil // Skip unreadable
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			lowerLine := strings.ToLower(line)

			if strings.Contains(lowerLine, query) {
				// Parse timestamp if available: - [HH:MM] Content
				// Reuse regex from manager.go?
				// For search results, raw content is fine, but metadata helps.
				// Let's extract basic info.

				// Basic date extraction from filename
				filename := filepath.Base(path)
				dateStr := strings.TrimSuffix(filename, ".md")

				matches := noteLineRegex.FindStringSubmatch(line)
				timeStr := ""
				content := line
				if len(matches) == 3 {
					timeStr = matches[1]
					content = matches[2]
				}

				results = append(results, SearchResult{
					Content:  content,
					Date:     dateStr,
					Time:     timeStr,
					FilePath: path,
					LineNo:   lineNo,
				})

				// Limit results?
				if len(results) >= 100 {
					return fs.SkipAll
				}
			}
		}
		return nil
	})

	return results, err
}
