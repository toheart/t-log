package note

// NoteEntry represents a single note item for the frontend
type NoteEntry struct {
	Content   string `json:"content"`   // Note content
	Timestamp string `json:"timestamp"` // Display time (HH:MM)
	Date      string `json:"date"`      // Date (YYYY-MM-DD)
	RawLine   string `json:"-"`         // Original raw line from file (internal use)
}

// DailyNote represents a full day's note content
type DailyNote struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}