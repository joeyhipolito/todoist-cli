package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joeyhipolito/todoist-cli/internal/api"
	"github.com/joeyhipolito/todoist-cli/internal/transform"
)

// AddCmd creates a new task.
func AddCmd(token string, args []string, jsonOutput bool) error {
	if len(args) < 1 {
		return fmt.Errorf("add requires a task name\n\nUsage: todoist add <task name> [--date <date>] [--priority <1-4>] [--project <name>] [--labels <l1,l2>]")
	}

	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	// First non-flag argument is the task content
	content := args[0]
	args = args[1:]

	req := &api.CreateTaskRequest{
		Content: content,
	}

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--date":
			if i+1 >= len(args) {
				return fmt.Errorf("--date requires an argument")
			}
			req.DueString = args[i+1]
			i++
		case "--priority":
			if i+1 >= len(args) {
				return fmt.Errorf("--priority requires an argument (1-4)")
			}
			p, err := transform.ParsePriority(args[i+1])
			if err != nil {
				return err
			}
			req.Priority = p
			i++
		case "--project":
			if i+1 >= len(args) {
				return fmt.Errorf("--project requires an argument")
			}
			// Resolve project name to ID
			projects, err := client.GetProjects()
			if err != nil {
				return fmt.Errorf("failed to fetch projects: %w", err)
			}
			projectName := args[i+1]
			found := false
			for _, p := range projects {
				if strings.EqualFold(p.Name, projectName) {
					req.ProjectID = p.ID
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("project not found: %s", projectName)
			}
			i++
		case "--labels":
			if i+1 >= len(args) {
				return fmt.Errorf("--labels requires an argument")
			}
			req.Labels = strings.Split(args[i+1], ",")
			i++
		default:
			return fmt.Errorf("unknown flag: %s", args[i])
		}
	}

	task, err := client.CreateTask(req)
	if err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(task)
	}

	fmt.Printf("Created task: %s (ID: %s)\n", task.Content, task.ID)
	return nil
}
