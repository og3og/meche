package memory

import (
	"errors"
	"meche/pkg/models"
	"sync"
)

type MemoryProjectStorage struct {
	projects map[string]*models.Project
	mu       sync.RWMutex
}

func NewMemoryProjectStorage() *MemoryProjectStorage {
	return &MemoryProjectStorage{
		projects: make(map[string]*models.Project),
	}
}

func (s *MemoryProjectStorage) CreateProject(project *models.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.projects[project.ID] = project
	return nil
}

func (s *MemoryProjectStorage) GetProject(id string) (*models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, exists := s.projects[id]
	if !exists {
		return nil, errors.New("project not found")
	}
	return project, nil
}

func (s *MemoryProjectStorage) ListProjects(orgID string) ([]*models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var projects []*models.Project
	for _, project := range s.projects {
		if project.OrgID == orgID && project.DeletedAt == nil {
			projects = append(projects, project)
		}
	}
	return projects, nil
}

func (s *MemoryProjectStorage) UpdateProject(project *models.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.projects[project.ID]; !exists {
		return errors.New("project not found")
	}

	s.projects[project.ID] = project
	return nil
}

func (s *MemoryProjectStorage) DeleteProject(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	project, exists := s.projects[id]
	if !exists {
		return errors.New("project not found")
	}

	project.Delete()
	return nil
}
