package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joeyhipolito/todoist-cli/internal/api"
	"github.com/joeyhipolito/todoist-cli/internal/transform"
)

// CompletedCmd lists completed tasks with optional filtering.
func CompletedCmd(token string, args []string, jsonOutput bool) error {
	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	// Parse flags
	var projectName, since string
	limit := 50
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 >= len(args) {
				return fmt.Errorf("--project requires a name")
			}
			projectName = args[i+1]
			i++
		case "--since":
			if i+1 >= len(args) {
				return fmt.Errorf("--since requires a date (YYYY-MM-DD)")
			}
			since = args[i+1]
			i++
		case "--limit":
			if i+1 >= len(args) {
				return fmt.Errorf("--limit requires a number")
			}
			n := 0
			if _, err := fmt.Sscanf(args[i+1], "%d", &n); err != nil || n <= 0 {
				return fmt.Errorf("--limit must be a positive integer")
			}
			limit = n
			i++
		}
	}

	// Resolve project name to ID
	var projectID string
	if projectName != "" {
		projects, err := client.GetProjects()
		if err != nil {
			return fmt.Errorf("failed to resolve project: %w", err)
		}
		for _, p := range projects {
			if strings.EqualFold(p.Name, projectName) {
				projectID = p.ID
				break
			}
		}
		if projectID == "" {
			return fmt.Errorf("project not found: %s", projectName)
		}
	}

	// Convert YYYY-MM-DD to ISO datetime for the API
	sinceParam := since
	if since != "" && len(since) == 10 {
		sinceParam = since + "T00:00:00"
	}

	resp, err := client.GetCompletedTasks(projectID, sinceParam, limit)
	if err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(resp)
	}

	if len(resp.Items) == 0 {
		fmt.Println("No completed tasks found.")
		return nil
	}

	for _, t := range resp.Items {
		fmt.Println(transform.FormatCompletedTaskLine(&t))
	}

	return nil
}
