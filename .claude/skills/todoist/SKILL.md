---
name: todoist
description: Manages Todoist tasks, projects, and labels via todoist CLI. Use when user asks about tasks, to-dos, deadlines, or wants to add/complete tasks.
allowed-tools: Bash(todoist:*)
---

# Todoist - Task Management

Task management via the standalone `todoist` CLI tool.

## When to Use

- User mentions tasks, to-dos, deadlines, or what's due
- User wants to add a new task
- User wants to check what's due today or overdue
- User says they completed something
- User asks about projects or task organization

## Commands

### View Tasks

```bash
todoist list                              # All tasks
todoist list --filter "today"             # Due today
todoist list --filter "overdue"           # Overdue tasks
todoist list --filter "#Work"             # Tasks in Work project
todoist list --json                       # JSON output
```

### Completed Tasks

```bash
todoist completed                            # Recently completed tasks
todoist completed --since "2026-02-01"       # Completed after date
todoist completed --project "Work"           # Completed in project
todoist completed --limit 10                 # Limit results
todoist completed --json                     # JSON output
```

### Add Tasks

```bash
# Simple task
todoist add "Buy groceries"

# With due date and priority
todoist add "Submit report" --date "tomorrow" --priority 1

# With project and labels
todoist add "Review PR" --project "Work" --labels "urgent"
```

### Complete & Delete

```bash
todoist close <task-id>                   # Mark task complete
todoist delete <task-id>                  # Delete task
```

### Projects

```bash
todoist projects                          # List all projects
todoist projects add "Project Name"       # Create a new project
todoist projects delete <project-id>      # Delete a project
```

### Labels

```bash
todoist labels                            # List all labels
```

### Setup

```bash
todoist configure                         # Interactive setup (API token)
todoist configure show                    # Show current config
todoist doctor                            # Health checks
```

## Configuration

Config file: `~/.todoist/config`

```ini
access_token=your-todoist-api-token
```

Get token from: https://todoist.com/app/settings/integrations/developer

Or use environment variable: `TODOIST_API_TOKEN`

## Examples

**User**: "what's due today"
**Action**: `todoist list --filter "today"`

**User**: "add task: buy milk"
**Action**: `todoist add "Buy milk"`

**User**: "finished the grocery shopping"
**Action**: `todoist list --filter "today"` to find the task, then `todoist close <id>`

**User**: "what's overdue"
**Action**: `todoist list --filter "overdue"`

All commands support `--json` for machine-readable output.
