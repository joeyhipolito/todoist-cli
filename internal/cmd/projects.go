package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joeyhipolito/todoist-cli/internal/api"
)

// ProjectsCmd lists all projects.
func ProjectsCmd(token string, args []string, jsonOutput bool) error {
	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	projects, err := client.GetProjects()
	if err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(projects)
	}

	if len(projects) == 0 {
		fmt.Println("No projects found.")
		return nil
	}

	for _, p := range projects {
		marker := " "
		if p.IsFavorite {
			marker = "*"
		}
		fmt.Printf("  %s %s (%s)\n", marker, p.Name, p.ID)
	}

	return nil
}
