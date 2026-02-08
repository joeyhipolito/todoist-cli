# todoist-cli

A standalone Go CLI for managing [Todoist](https://todoist.com/) tasks, projects, and labels from the terminal. Zero external dependencies.

## Features

- **Task management** — list, add, close, and delete tasks
- **Filtering** — filter by date, priority, project, or label
- **Project and label listing** — view all projects and labels
- **Priority support** — P1–P4 priority levels
- **Due dates** — natural language dates (`today`, `tomorrow`) or `YYYY-MM-DD`
- **Interactive configuration** — `todoist configure` setup
- **Diagnostics** — built-in `doctor` command for troubleshooting
- **JSON output** — machine-readable format for scripting (`--json`)
- **Cross-platform** — macOS (arm64/amd64) and Linux (amd64/arm64)
- **Pure Go stdlib** — zero external dependencies

## Installation

### Prerequisites

- Go 1.25 or later
- Todoist API token ([get one here](https://todoist.com/app/settings/integrations/developer))

### Build and Install

```bash
make install            # Build and symlink to ~/bin
todoist configure       # Interactive setup
todoist doctor          # Verify everything works
```

## Configuration

Config is stored in `~/.todoist/config` (INI format, `chmod 600`).

```bash
todoist configure           # Interactive setup (recommended)
todoist configure show      # Show current config (token masked)
```

### Config file keys

| Key | Description |
|-----|-------------|
| `access_token` | Todoist API token |

### Environment variables (fallback)

| Variable | Description |
|----------|-------------|
| `TODOIST_ACCESS_TOKEN` | API token (used if no config file) |

## Commands

### Listing tasks

```bash
todoist list                        # Today's and overdue tasks (default)
todoist list --filter "overdue"     # Overdue tasks only
todoist list --filter "#Work"       # Tasks in Work project
todoist list --filter "p1"          # Priority 1 tasks
todoist list --filter "@urgent"     # Tasks with @urgent label
todoist list --json                 # JSON output
```

### Adding tasks

```bash
# Basic task
todoist add "Buy groceries"

# With options
todoist add "Review PR" \
  --date tomorrow \
  --priority 1 \
  --project "Work" \
  --labels "dev,urgent"

# With specific date
todoist add "File taxes" --date 2024-04-15
```

### Completing and deleting

```bash
todoist close <task-id>     # Mark task as complete
todoist delete <task-id>    # Permanently delete task
```

### Projects and labels

```bash
todoist projects            # List all projects
todoist projects --json     # JSON output
todoist labels              # List all labels
todoist labels --json       # JSON output
```

### Diagnostics

```bash
todoist doctor              # Check binary, config, permissions, API connectivity
```

## Architecture

```
cmd/todoist-cli/             # Entry point and command routing
internal/
├── api/                     # Todoist REST API v2 client
│   ├── client.go            # HTTP client with retry and rate limiting
│   ├── methods.go           # GetTasks, CreateTask, CloseTask, etc.
│   ├── types.go             # Task, Project, Label, Due types
│   └── errors.go            # Error types with IsRetryable(), IsAuthError()
├── cmd/                     # Command implementations
│   ├── list.go              # List tasks with filters
│   ├── add.go               # Add new tasks
│   ├── close.go             # Complete tasks
│   ├── delete.go            # Delete tasks
│   ├── projects.go          # List projects
│   ├── labels.go            # List labels
│   ├── configure.go         # Configuration management
│   └── doctor.go            # Diagnostics
├── config/                  # Config file loading/saving
└── transform/               # Display formatting
    ├── priority.go          # Priority conversion (UI ↔ API)
    ├── date.go              # Date formatting and overdue detection
    └── display.go           # Human-readable output
```

### Design decisions

- **Priority inversion** — Todoist API uses inverted priorities (API `4` = UI `P1`). The `transform/priority.go` module handles conversion automatically.
- **Retry with backoff** — exponential backoff (1s, 2s, 4s) with rate-limit (`429`) awareness and 60s retry-after
- **No CLI framework** — simple string-based command dispatch, no external dependencies
- **Secure config** — config directory `700`, config file `600` permissions

## Development

```bash
make build              # Build for current platform
make install            # Build and install to ~/bin
make test               # Run tests
make test-coverage      # Tests with HTML coverage report
make build-all          # Cross-compile for macOS/Linux
make fmt                # Format code
make vet                # Vet for issues
make lint               # fmt + vet
make clean              # Remove build artifacts
```

## License

MIT License. See [LICENSE](LICENSE) for details.

## Links

- [Todoist API Documentation](https://developer.todoist.com/rest/v2/)
- [Get API Token](https://todoist.com/app/settings/integrations/developer)
