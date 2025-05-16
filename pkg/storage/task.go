package storage

import "meche/pkg/models"

// TaskStorage defines the interface for task storage operations
type TaskStorage interface {
	// CreateTask creates a new task
	CreateTask(task *models.Task) error

	// GetTask retrieves a task by ID
	GetTask(id string, orgID string) (*models.Task, error)

	// ListTasks returns all tasks for a project
	ListTasks(projectID string, orgID string) ([]*models.Task, error)

	// ListTasksByOrganization returns all tasks for an organization
	ListTasksByOrganization(orgID string) ([]*models.Task, error)

	// UpdateTask updates an existing task
	UpdateTask(task *models.Task, orgID string) error

	// DeleteTask soft deletes a task
	DeleteTask(id string, orgID string) error
}
