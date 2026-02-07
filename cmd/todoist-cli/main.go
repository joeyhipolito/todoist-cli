// Package main implements the todoist binary.
package main

import (
	"fmt"
	"os"

	"github.com/joeyhipolito/todoist-cli/internal/cmd"
	"github.com/joeyhipolito/todoist-cli/internal/config"
)

const version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]

	// Handle help and version flags
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return nil
	}

	if args[0] == "--version" || args[0] == "-v" {
		fmt.Printf("todoist version %s\n", version)
		return nil
	}

	// Parse subcommand
	subcommand := args[0]
	remainingArgs := args[1:]

	// Check for global --json flag
	jsonOutput := false
	var filteredArgs []string
	for _, arg := range remainingArgs {
		if arg == "--json" {
			jsonOutput = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	// Commands that don't require authentication
	switch subcommand {
	case "configure":
		if len(filteredArgs) > 0 && filteredArgs[0] == "show" {
			return cmd.ConfigureShowCmd(jsonOutput)
		}
		return cmd.ConfigureCmd()
	case "doctor":
		return cmd.DoctorCmd(jsonOutput)
	}

	// Resolve access token: config file > environment variable
	token := config.ResolveToken()
	if token == "" {
		return fmt.Errorf("no access token found\n\nRun 'todoist configure' to set up, or set TODOIST_ACCESS_TOKEN")
	}

	// Dispatch to authenticated commands
	switch subcommand {
	case "list":
		return cmd.ListCmd(token, filteredArgs, jsonOutput)
	case "add":
		return cmd.AddCmd(token, filteredArgs, jsonOutput)
	case "close":
		return cmd.CloseCmd(token, filteredArgs, jsonOutput)
	case "delete":
		return cmd.DeleteCmd(token, filteredArgs, jsonOutput)
	case "projects":
		return cmd.ProjectsCmd(token, filteredArgs, jsonOutput)
	case "labels":
		return cmd.LabelsCmd(token, filteredArgs, jsonOutput)
	default:
		return fmt.Errorf("unknown command: %s\n\nRun 'todoist --help' for usage", subcommand)
	}
}

func printUsage() {
	fmt.Printf(`todoist - Todoist command-line interface (v%s)

USAGE:
    todoist <command> [options]

COMMANDS:
    list                    List tasks (default: today & overdue)
    add                     Add a new task
    close                   Complete a task
    delete                  Delete a task
    projects                List all projects
    labels                  List all labels
    configure               Set up Todoist access token
    configure show          Show current configuration
    doctor                  Validate installation and configuration

LIST TASKS:
    todoist list [options]
        --filter <query>        Filter (today, overdue, p1, @label, #project)

ADD TASK:
    todoist add <task name> [options]
        --date <date>           Due date (today, tomorrow, YYYY-MM-DD)
        --priority <1-4>        Priority (1=urgent, 4=normal)
        --project <name>        Target project
        --labels <l1,l2>        Comma-separated labels

CLOSE/DELETE TASK:
    todoist close <task-id>
    todoist delete <task-id>

GLOBAL OPTIONS:
    --json              Output in JSON format
    --help, -h          Show this help
    --version, -v       Show version

CONFIGURATION:
    todoist configure           Interactive setup
    todoist configure show      Show current config (token masked)
    todoist doctor              Validate setup and troubleshoot
    Config file: ~/.todoist/config

EXAMPLES:
    todoist configure                           # First-time setup
    todoist list                                # Today's tasks
    todoist list --filter "overdue"             # Overdue tasks
    todoist list --filter "#Work"               # Tasks in Work project
    todoist add "Buy groceries" --date tomorrow --priority 2
    todoist add "Review PR" --project "Work" --labels "dev,urgent"
    todoist close 1234567890                    # Complete a task
    todoist projects --json                     # List projects as JSON
    todoist doctor                              # Check setup

For more information, visit: https://developer.todoist.com/rest/v2/
`, version)
}
