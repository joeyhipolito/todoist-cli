package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// paginatedResponse wraps the v1 API list responses which return
// {"results": [...], "next_cursor": ...} instead of a raw array.
type paginatedResponse[T any] struct {
	Results    []T     `json:"results"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

// GetTasks returns all active tasks, optionally filtered by filter query and/or project ID.
func (c *Client) GetTasks(filter, projectID string) ([]*Task, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	if projectID != "" {
		params.Set("project_id", projectID)
	}

	endpoint := "/tasks"
	if encoded := params.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	body, _, err := c.request(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var resp paginatedResponse[*Task]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return resp.Results, nil
}

// CreateTask creates a new task.
func (c *Client) CreateTask(req *CreateTaskRequest) (*Task, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	body, _, err := c.request(http.MethodPost, "/tasks", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}

	return &task, nil
}

// CloseTask marks a task as complete. Returns 204 No Content on success.
func (c *Client) CloseTask(taskID string) error {
	_, statusCode, err := c.request(http.MethodPost, "/tasks/"+taskID+"/close", nil)
	if err != nil {
		return fmt.Errorf("failed to close task: %w", err)
	}
	if statusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}
	return nil
}

// DeleteTask permanently deletes a task. Returns 204 No Content on success.
func (c *Client) DeleteTask(taskID string) error {
	_, statusCode, err := c.request(http.MethodDelete, "/tasks/"+taskID, nil)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	if statusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}
	return nil
}

// GetProjects returns all projects.
func (c *Client) GetProjects() ([]*Project, error) {
	body, _, err := c.request(http.MethodGet, "/projects", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	var resp paginatedResponse[*Project]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return resp.Results, nil
}

// CreateProject creates a new project.
func (c *Client) CreateProject(req *CreateProjectRequest) (*Project, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	body, _, err := c.request(http.MethodPost, "/projects", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}

	return &project, nil
}

// DeleteProject permanently deletes a project. Returns 204 No Content on success.
func (c *Client) DeleteProject(projectID string) error {
	_, statusCode, err := c.request(http.MethodDelete, "/projects/"+projectID, nil)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	if statusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}
	return nil
}

// GetLabels returns all personal labels.
func (c *Client) GetLabels() ([]*Label, error) {
	body, _, err := c.request(http.MethodGet, "/labels", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get labels: %w", err)
	}

	var resp paginatedResponse[*Label]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %w", err)
	}

	return resp.Results, nil
}
