package transform

import (
	"time"
)

// FormatDueDate converts a Todoist due date string to a human-friendly display.
// Handles both date-only (YYYY-MM-DD) and datetime (RFC3339) formats.
//
// Examples:
//
//	FormatDueDate("2024-01-15", "")           // "2024-01-15"
//	FormatDueDate("2024-01-15", "2024-01-15T10:00:00Z")  // "2024-01-15 10:00"
//	FormatDueDate("today's date", "")         // "Today"
//	FormatDueDate("tomorrow's date", "")      // "Tomorrow"
func FormatDueDate(date, datetime string) string {
	if datetime != "" {
		t, err := time.Parse(time.RFC3339, datetime)
		if err == nil {
			return t.Local().Format("2006-01-02 15:04")
		}
	}

	if date == "" {
		return ""
	}

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date // Return as-is if unparseable
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())

	switch {
	case d.Equal(today):
		return "Today"
	case d.Equal(today.AddDate(0, 0, 1)):
		return "Tomorrow"
	case d.Equal(today.AddDate(0, 0, -1)):
		return "Yesterday"
	case d.Before(today):
		return t.Format("2006-01-02") + " (overdue)"
	default:
		return t.Format("2006-01-02")
	}
}

// IsOverdue returns true if the given date string is before today.
func IsOverdue(date string) bool {
	if date == "" {
		return false
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())
	return d.Before(today)
}
