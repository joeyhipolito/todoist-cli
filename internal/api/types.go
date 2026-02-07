package api

// Task represents a Todoist task.
type Task struct {
	ID           string   `json:"id"`
	ProjectID    string   `json:"project_id"`
	SectionID    string   `json:"section_id,omitempty"`
	Content      string   `json:"content"`
	Description  string   `json:"description,omitempty"`
	IsCompleted  bool     `json:"is_completed"`
	Labels       []string `json:"labels"`
	ParentID     string   `json:"parent_id,omitempty"`
	Order        int      `json:"order"`
	Priority     int      `json:"priority"` // 1=normal, 4=urgent (inverted from UI)
	Due          *Due     `json:"due,omitempty"`
	URL          string   `json:"url"`
	CommentCount int      `json:"comment_count"`
	CreatorID    string   `json:"creator_id"`
	CreatedAt    string   `json:"created_at"`
	AssigneeID   string   `json:"assignee_id,omitempty"`
	AssignerID   string   `json:"assigner_id,omitempty"`
	Duration     *Duration `json:"duration,omitempty"`
}

// Due represents a task's due date.
type Due struct {
	String      string `json:"string"`
	Date        string `json:"date"`       // YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime,omitempty"` // RFC3339
	Timezone    string `json:"timezone,omitempty"`
}

// Duration represents a task's duration.
type Duration struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"` // "minute" or "day"
}

// Project represents a Todoist project.
type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CommentCount   int    `json:"comment_count"`
	Order          int    `json:"order"`
	Color          string `json:"color"`
	IsShared       bool   `json:"is_shared"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	ViewStyle      string `json:"view_style"` // "list" or "board"
	URL            string `json:"url"`
	ParentID       string `json:"parent_id,omitempty"`
}

// Label represents a Todoist label.
type Label struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Order      int    `json:"order"`
	IsFavorite bool   `json:"is_favorite"`
}

// CreateTaskRequest represents the payload for creating a task.
type CreateTaskRequest struct {
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"project_id,omitempty"`
	SectionID   string   `json:"section_id,omitempty"`
	ParentID    string   `json:"parent_id,omitempty"`
	Order       int      `json:"order,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	DueString   string   `json:"due_string,omitempty"`
	DueDate     string   `json:"due_date,omitempty"`
	DueDatetime string   `json:"due_datetime,omitempty"`
	DueLang     string   `json:"due_lang,omitempty"`
	AssigneeID  string   `json:"assignee_id,omitempty"`
}
