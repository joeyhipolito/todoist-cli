package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joeyhipolito/todoist-cli/internal/api"
)

// ProjectsCmd handles the "projects" command and its subcommands.
func ProjectsCmd(token string, args []string, jsonOutput bool) error {
	// No subcommand → list projects (default)
	if len(args) == 0 {
		return projectsListCmd(token, jsonOutput)
	}

	// Dispatch subcommand
	switch args[0] {
	case "add":
		return projectsAddCmd(token, args[1:], jsonOutput)
	case "delete":
		return projectsDeleteCmd(token, args[1:], jsonOutput)
	default:
		return fmt.Errorf("unknown projects subcommand: %s\n\nUsage:\n  todoist projects              List all projects\n  todoist projects add <name>   Create a project\n  todoist projects delete <id>  Delete a project", args[0])
	}
}

// projectsListCmd lists all projects.
func projectsListCmd(token string, jsonOutput bool) error {
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

// projectsAddCmd creates a new project.
func projectsAddCmd(token string, args []string, jsonOutput bool) error {
	if len(args) < 1 {
		return fmt.Errorf("projects add requires a project name\n\nUsage: todoist projects add <name>")
	}

	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	req := &api.CreateProjectRequest{
		Name: args[0],
	}

	project, err := client.CreateProject(req)
	if err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(project)
	}

	fmt.Printf("Created project: %s (ID: %s)\n", project.Name, project.ID)
	return nil
}

// projectsDeleteCmd permanently deletes a project.
func projectsDeleteCmd(token string, args []string, jsonOutput bool) error {
	if len(args) < 1 {
		return fmt.Errorf("projects delete requires a project ID\n\nUsage: todoist projects delete <project-id>")
	}

	client, err := api.NewClient(token)
	if err != nil {
		return err
	}

	projectID := args[0]
	if err := client.DeleteProject(projectID); err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(map[string]string{
			"status":     "ok",
			"project_id": projectID,
			"message":    "Project deleted",
		})
	}

	fmt.Printf("Project %s deleted.\n", projectID)
	return nil
}
