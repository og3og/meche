package storage

import "meche/pkg/models"

type ProjectStorage interface {
	CreateProject(project *models.Project) error
	GetProject(id string) (*models.Project, error)
	ListProjects(orgID string) ([]*models.Project, error)
	UpdateProject(project *models.Project) error
	DeleteProject(id string) error
}
