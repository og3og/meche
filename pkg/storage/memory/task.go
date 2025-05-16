package memory

import (
	"errors"
	"meche/pkg/models"
	"meche/pkg/storage"
	"sync"
)

// MemoryTaskStorage implements TaskStorage using in-memory storage
type MemoryTaskStorage struct {
	tasks map[string]*models.Task
	mu    sync.RWMutex
}

// NewMemoryTaskStorage creates a new instance of MemoryTaskStorage
func NewMemoryTaskStorage() storage.TaskStorage {
	return &MemoryTaskStorage{
		tasks: make(map[string]*models.Task),
	}
}

// CreateTask creates a new task
func (s *MemoryTaskStorage) CreateTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; exists {
		return errors.New("task already exists")
	}

	if task.OrgID == "" {
		return errors.New("task must belong to an organization")
	}

	if task.ProjectID == "" {
		return errors.New("task must belong to a project")
	}

	s.tasks[task.ID] = task
	return nil
}

// GetTask retrieves a task by ID
func (s *MemoryTaskStorage) GetTask(id string, orgID string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists || task.DeletedAt != nil {
		return nil, errors.New("task not found")
	}

	if task.OrgID != orgID {
		return nil, errors.New("unauthorized: task does not belong to the specified organization")
	}

	return task, nil
}

// ListTasks returns all tasks for a project
func (s *MemoryTaskStorage) ListTasks(projectID string, orgID string) ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0)
	for _, task := range s.tasks {
		if task.ProjectID == projectID && task.OrgID == orgID && task.DeletedAt == nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// ListTasksByOrganization returns all tasks for an organization
func (s *MemoryTaskStorage) ListTasksByOrganization(orgID string) ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0)
	for _, task := range s.tasks {
		if task.OrgID == orgID && task.DeletedAt == nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (s *MemoryTaskStorage) UpdateTask(task *models.Task, orgID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.tasks[task.ID]
	if !exists || existing.DeletedAt != nil {
		return errors.New("task not found")
	}

	if existing.OrgID != orgID {
		return errors.New("unauthorized: task does not belong to the specified organization")
	}

	// Ensure organization and project IDs don't change
	if existing.OrgID != task.OrgID {
		return errors.New("cannot change task's organization")
	}
	if existing.ProjectID != task.ProjectID {
		return errors.New("cannot change task's project")
	}

	s.tasks[task.ID] = task
	return nil
}

// DeleteTask soft deletes a task
func (s *MemoryTaskStorage) DeleteTask(id string, orgID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists || task.DeletedAt != nil {
		return errors.New("task not found")
	}

	if task.OrgID != orgID {
		return errors.New("unauthorized: task does not belong to the specified organization")
	}

	task.Delete()
	return nil
}
