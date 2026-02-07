package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/joeyhipolito/todoist-cli/internal/api"
	"github.com/joeyhipolito/todoist-cli/internal/config"
)

// DoctorCheck represents a single doctor check result.
type DoctorCheck struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "ok", "warn", "fail"
	Message string `json:"message"`
}

// DoctorOutput represents the JSON output of the doctor command.
type DoctorOutput struct {
	Checks  []DoctorCheck `json:"checks"`
	Summary string        `json:"summary"`
	AllOK   bool          `json:"all_ok"`
}

// DoctorCmd validates the Todoist CLI installation and configuration.
func DoctorCmd(jsonOutput bool) error {
	var checks []DoctorCheck
	allOK := true

	// 1. Check binary location
	binaryPath, err := exec.LookPath("todoist")
	if err != nil {
		checks = append(checks, DoctorCheck{
			Name:    "Binary",
			Status:  "warn",
			Message: "todoist not found in PATH (running from local build?)",
		})
	} else {
		checks = append(checks, DoctorCheck{
			Name:    "Binary",
			Status:  "ok",
			Message: binaryPath,
		})
	}

	// 2. Check config file exists
	configPath := config.Path()
	if !config.Exists() {
		checks = append(checks, DoctorCheck{
			Name:    "Config file",
			Status:  "fail",
			Message: fmt.Sprintf("%s not found. Run 'todoist configure'", configPath),
		})
		allOK = false
	} else {
		checks = append(checks, DoctorCheck{
			Name:    "Config file",
			Status:  "ok",
			Message: configPath,
		})

		// 3. Check config permissions
		perms, err := config.Permissions()
		if err != nil {
			checks = append(checks, DoctorCheck{
				Name:    "Config permissions",
				Status:  "fail",
				Message: fmt.Sprintf("Cannot read permissions: %v", err),
			})
			allOK = false
		} else if perms != 0600 {
			checks = append(checks, DoctorCheck{
				Name:    "Config permissions",
				Status:  "warn",
				Message: fmt.Sprintf("%o (should be 600). Fix: chmod 600 %s", perms, configPath),
			})
		} else {
			checks = append(checks, DoctorCheck{
				Name:    "Config permissions",
				Status:  "ok",
				Message: "600 (secure)",
			})
		}
	}

	// 4. Check access token
	cfg, err := config.Load()
	if err != nil {
		checks = append(checks, DoctorCheck{
			Name:    "Config format",
			Status:  "fail",
			Message: fmt.Sprintf("Failed to parse config: %v", err),
		})
		allOK = false
	} else {
		token := cfg.AccessToken
		if token == "" {
			token = os.Getenv("TODOIST_ACCESS_TOKEN")
		}

		if token == "" {
			checks = append(checks, DoctorCheck{
				Name:    "Access token",
				Status:  "fail",
				Message: "Not found in config or TODOIST_ACCESS_TOKEN env var",
			})
			allOK = false
		} else {
			// Mask token
			masked := "****"
			if len(token) > 8 {
				masked = token[:4] + "..." + token[len(token)-4:]
			}
			checks = append(checks, DoctorCheck{
				Name:    "Access token",
				Status:  "ok",
				Message: fmt.Sprintf("Present (%s)", masked),
			})

			// 5. Test API connection
			client, err := api.NewClient(token)
			if err != nil {
				checks = append(checks, DoctorCheck{
					Name:    "API connection",
					Status:  "fail",
					Message: fmt.Sprintf("Failed to create client: %v", err),
				})
				allOK = false
			} else {
				projects, err := client.GetProjects()
				if err != nil {
					checks = append(checks, DoctorCheck{
						Name:    "API connection",
						Status:  "fail",
						Message: fmt.Sprintf("Failed: %v", err),
					})
					allOK = false
				} else {
					checks = append(checks, DoctorCheck{
						Name:    "API connection",
						Status:  "ok",
						Message: fmt.Sprintf("Success (%d project(s) found)", len(projects)),
					})
				}
			}
		}
	}

	// Determine summary
	summary := "All checks passed!"
	if !allOK {
		failCount := 0
		for _, c := range checks {
			if c.Status == "fail" {
				failCount++
			}
		}
		summary = fmt.Sprintf("%d check(s) failed. Run 'todoist configure' to fix.", failCount)
	}

	// JSON output
	if jsonOutput {
		output := DoctorOutput{
			Checks:  checks,
			Summary: summary,
			AllOK:   allOK,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(output)
	}

	// Human-readable output
	fmt.Println("Todoist CLI Doctor")
	fmt.Println("==================")
	fmt.Println()

	for _, c := range checks {
		var icon string
		switch c.Status {
		case "ok":
			icon = "OK"
		case "warn":
			icon = "WARN"
		case "fail":
			icon = "FAIL"
		}
		fmt.Printf("  [%4s] %-20s %s\n", icon, c.Name+":", c.Message)
	}

	fmt.Println()
	if allOK {
		fmt.Println(summary)
	} else {
		fmt.Println(summary)
		return fmt.Errorf("doctor checks failed")
	}

	return nil
}
