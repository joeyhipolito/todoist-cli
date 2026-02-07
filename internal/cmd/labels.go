package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joeyhipolito/todoist-cli/internal/api"
)

// LabelsCmd lists all personal labels.
func LabelsCmd(token string, args []string, jsonOutput bool) error {
	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	labels, err := client.GetLabels()
	if err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(labels)
	}

	if len(labels) == 0 {
		fmt.Println("No labels found.")
		return nil
	}

	for _, l := range labels {
		marker := " "
		if l.IsFavorite {
			marker = "*"
		}
		fmt.Printf("  %s @%s (%s)\n", marker, l.Name, l.ID)
	}

	return nil
}
