package api

// Task represents a Todoist task (API v1).
type Task struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	SectionID   string    `json:"section_id,omitempty"`
	Content     string    `json:"content"`
	Description string    `json:"description,omitempty"`
	IsCompleted bool      `json:"checked"`
	Labels      []string  `json:"labels"`
	ParentID    string    `json:"parent_id,omitempty"`
	Order       int       `json:"child_order"`
	Priority    int       `json:"priority"` // 1=normal, 4=urgent (inverted from UI)
	Due         *Due      `json:"due,omitempty"`
	Deadline    *Deadline `json:"deadline,omitempty"`
	NoteCount   int       `json:"note_count"`
	CreatorID   string    `json:"added_by_uid"`
	CreatedAt   string    `json:"added_at"`
	CompletedAt string    `json:"completed_at,omitempty"`
	AssigneeID  string    `json:"responsible_uid,omitempty"`
	AssignerID  string    `json:"assigned_by_uid,omitempty"`
	Duration    *Duration `json:"duration,omitempty"`
}

// Due represents a task's due date.
type Due struct {
	String      string `json:"string"`
	Date        string `json:"date"` // YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime,omitempty"` // RFC3339
	Timezone    string `json:"timezone,omitempty"`
	Lang        string `json:"lang,omitempty"` // IETF language tag
}

// Deadline represents a task's hard deadline (separate from due date).
type Deadline struct {
	Date string `json:"date"`           // YYYY-MM-DD
	Lang string `json:"lang,omitempty"` // IETF language tag
}

// Duration represents a task's duration.
type Duration struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"` // "minute" or "day"
}

// Project represents a Todoist project (API v1).
type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Order          int    `json:"child_order"`
	Color          string `json:"color"`
	IsShared       bool   `json:"is_shared"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"inbox_project"`
	ViewStyle      string `json:"view_style"` // "list", "board", or "calendar"
	ParentID       string `json:"parent_id,omitempty"`
	CreatorID      string `json:"creator_uid,omitempty"`
}

// CompletedTask represents a completed task from the Todoist API v1.
// This is a different shape than active Task — returned by GET /tasks/completed.
type CompletedTask struct {
	ID          string  `json:"id"`
	TaskID      string  `json:"task_id"`
	ProjectID   string  `json:"project_id"`
	SectionID   *string `json:"section_id"`
	Content     string  `json:"content"`
	CompletedAt string  `json:"completed_at"`
	NoteCount   int     `json:"note_count"`
}

// CompletedResponse wraps the response from GET /tasks/completed.
type CompletedResponse struct {
	Items    []CompletedTask    `json:"items"`
	Projects map[string]Project `json:"projects"`
}

// Label represents a Todoist personal label.
type Label struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Order      int    `json:"order"`
	IsFavorite bool   `json:"is_favorite"`
}

// CreateProjectRequest represents the payload for creating a project.
type CreateProjectRequest struct {
	Name string `json:"name"`
}

// CreateTaskRequest represents the payload for creating a task.
type CreateTaskRequest struct {
	Content      string   `json:"content"`
	Description  string   `json:"description,omitempty"`
	ProjectID    string   `json:"project_id,omitempty"`
	SectionID    string   `json:"section_id,omitempty"`
	ParentID     string   `json:"parent_id,omitempty"`
	Order        int      `json:"order,omitempty"`
	Labels       []string `json:"labels,omitempty"`
	Priority     int      `json:"priority,omitempty"`
	DueString    string   `json:"due_string,omitempty"`
	DueDate      string   `json:"due_date,omitempty"`
	DueDatetime  string   `json:"due_datetime,omitempty"`
	DueLang      string   `json:"due_lang,omitempty"`
	AssigneeID   string   `json:"assignee_id,omitempty"`
	DeadlineDate string   `json:"deadline_date,omitempty"`
}
