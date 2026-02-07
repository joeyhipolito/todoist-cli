package transform

import (
	"fmt"
	"strings"

	"github.com/joeyhipolito/todoist-cli/internal/api"
)

// FormatTaskLine renders a single task as a human-readable one-liner.
//
// Format: "  <id> [P1] Task content (due-date) @label1 @label2"
func FormatTaskLine(t *api.Task) string {
	priority := FormatPriority(t.Priority)

	due := ""
	if t.Due != nil {
		d := FormatDueDate(t.Due.Date, t.Due.Datetime)
		if d != "" {
			due = " (" + d + ")"
		}
	}

	labels := FormatLabels(t.Labels)

	return fmt.Sprintf("  %s [%s] %s%s%s", t.ID, priority, t.Content, due, labels)
}

// FormatLabels returns labels as a space-separated "@label" string.
// Returns empty string if no labels.
func FormatLabels(labels []string) string {
	if len(labels) == 0 {
		return ""
	}
	return " @" + strings.Join(labels, " @")
}

// FormatProjectLine renders a single project as a human-readable one-liner.
//
// Format: "  <id> Project Name (*)"  where (*) indicates favorite
func FormatProjectLine(p *api.Project) string {
	fav := ""
	if p.IsFavorite {
		fav = " (*)"
	}
	return fmt.Sprintf("  %s %s%s", p.ID, p.Name, fav)
}

// FormatLabelLine renders a single label as a human-readable one-liner.
//
// Format: "  <id> @label-name (*)"  where (*) indicates favorite
func FormatLabelLine(l *api.Label) string {
	fav := ""
	if l.IsFavorite {
		fav = " (*)"
	}
	return fmt.Sprintf("  %s @%s%s", l.ID, l.Name, fav)
}

// MaskToken returns a masked version of an access token for display.
// Shows first 4 and last 4 characters with "..." in between.
func MaskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
