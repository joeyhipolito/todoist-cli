package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joeyhipolito/todoist-cli/internal/api"
)

// DeleteCmd permanently deletes a task.
func DeleteCmd(token string, args []string, jsonOutput bool) error {
	if len(args) < 1 {
		return fmt.Errorf("delete requires a task ID\n\nUsage: todoist delete <task-id>")
	}

	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	taskID := args[0]
	if err := client.DeleteTask(taskID); err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(map[string]string{
			"status":  "ok",
			"task_id": taskID,
			"message": "Task deleted",
		})
	}

	fmt.Printf("Task %s deleted.\n", taskID)
	return nil
}
