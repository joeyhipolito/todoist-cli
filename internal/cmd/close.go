package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joeyhipolito/todoist-cli/internal/api"
)

// CloseCmd marks a task as complete.
func CloseCmd(token string, args []string, jsonOutput bool) error {
	if len(args) < 1 {
		return fmt.Errorf("close requires a task ID\n\nUsage: todoist close <task-id>")
	}

	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	taskID := args[0]
	if err := client.CloseTask(taskID); err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(map[string]string{
			"status":  "ok",
			"task_id": taskID,
			"message": "Task completed",
		})
	}

	fmt.Printf("Task %s completed.\n", taskID)
	return nil
}
