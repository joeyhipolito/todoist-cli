package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joeyhipolito/todoist-cli/internal/api"
	"github.com/joeyhipolito/todoist-cli/internal/transform"
)

// ListCmd lists active tasks with optional filtering.
func ListCmd(token string, args []string, jsonOutput bool) error {
	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	// Parse --filter flag
	filter := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--filter" {
			if i+1 >= len(args) {
				return fmt.Errorf("--filter requires a query string")
			}
			filter = args[i+1]
			i++
		}
	}

	tasks, err := client.GetTasks(filter, "")
	if err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(tasks)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	for _, t := range tasks {
		priority := transform.FormatPriority(t.Priority)
		due := ""
		if t.Due != nil {
			due = " (" + t.Due.Date + ")"
		}
		labels := ""
		if len(t.Labels) > 0 {
			labels = " @" + strings.Join(t.Labels, " @")
		}
		fmt.Printf("  %s [%s] %s%s%s\n", t.ID, priority, t.Content, due, labels)
	}

	return nil
}
