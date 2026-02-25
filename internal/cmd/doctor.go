package cmd

import (
	"fmt"
	"os"

	"github.com/joeyhipolito/publishing-shared/doctor"
	"github.com/joeyhipolito/todoist-cli/internal/api"
	"github.com/joeyhipolito/todoist-cli/internal/config"
)

// DoctorCmd validates the Todoist CLI installation and configuration.
func DoctorCmd(jsonOutput bool) error {
	var checks []doctor.Check
	allOK := true

	add := func(name, status, msg string) {
		checks = append(checks, doctor.Check{Name: name, Status: status, Message: msg})
		if status == "fail" {
			allOK = false
		}
	}

	// 1. Binary in PATH
	checks = append(checks, doctor.BinaryInPath("todoist")())

	// 2-3. Config exists + permissions
	configPath := config.Path()
	checks = append(checks, doctor.ConfigExists(configPath, "todoist configure")())
	checks = append(checks, doctor.ConfigPermissions(configPath)())

	if !config.Exists() {
		return doctor.Render(os.Stdout, checks, false, jsonOutput, "Todoist Doctor", true)
	}

	// 4. Config parseable + access token
	cfg, err := config.Load()
	if err != nil {
		add("Config format", "fail", fmt.Sprintf("failed to parse config: %v", err))
		return doctor.Render(os.Stdout, checks, false, jsonOutput, "Todoist Doctor", true)
	}

	token := cfg.AccessToken
	if token == "" {
		token = os.Getenv("TODOIST_ACCESS_TOKEN")
	}
	if token == "" {
		add("Access token", "fail", "not found in config or TODOIST_ACCESS_TOKEN env var")
		return doctor.Render(os.Stdout, checks, false, jsonOutput, "Todoist Doctor", true)
	}
	masked := "****"
	if len(token) > 8 {
		masked = token[:4] + "..." + token[len(token)-4:]
	}
	add("Access token", "ok", fmt.Sprintf("present (%s)", masked))

	// 5. API connection
	client, err := api.NewClient(token)
	if err != nil {
		add("API connection", "fail", fmt.Sprintf("failed to create client: %v", err))
	} else {
		projects, err := client.GetProjects()
		if err != nil {
			add("API connection", "fail", fmt.Sprintf("failed: %v", err))
		} else {
			add("API connection", "ok", fmt.Sprintf("success (%d project(s) found)", len(projects)))
		}
	}

	return doctor.Render(os.Stdout, checks, allOK, jsonOutput, "Todoist Doctor", true)
}
