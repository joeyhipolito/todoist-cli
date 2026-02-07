// Package transform provides data transformation utilities for Todoist CLI output.
package transform

import (
	"fmt"
	"strconv"
)

// Todoist API priority is inverted from the UI:
//   API priority 1 = UI "P4" (normal/no priority)
//   API priority 2 = UI "P3"
//   API priority 3 = UI "P2"
//   API priority 4 = UI "P1" (most urgent)

// FormatPriority converts an API priority value to a human-readable UI label.
func FormatPriority(apiPriority int) string {
	switch apiPriority {
	case 4:
		return "P1"
	case 3:
		return "P2"
	case 2:
		return "P3"
	case 1:
		return "P4"
	default:
		return "P4"
	}
}

// ParsePriority converts a user-facing priority string (1-4) to an API priority value.
// User input "1" means P1 (urgent) = API priority 4.
func ParsePriority(input string) (int, error) {
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > 4 {
		return 0, fmt.Errorf("priority must be 1-4 (1=urgent, 4=normal): %s", input)
	}
	// Invert: user P1 -> API 4, user P2 -> API 3, etc.
	return 5 - n, nil
}
