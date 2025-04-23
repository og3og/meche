package storage

import (
	"meche/pkg/models"
)

// OrganizationStorage defines the interface for organization storage operations
type OrganizationStorage interface {
	// CreateOrganization creates a new organization
	CreateOrganization(org *models.Organization) error

	// GetOrganization retrieves an organization by ID
	GetOrganization(id string) (*models.Organization, error)

	// ListOrganizations returns all organizations
	ListOrganizations() ([]*models.Organization, error)

	// UpdateOrganization updates an existing organization
	UpdateOrganization(org *models.Organization) error

	// DeleteOrganization deletes an organization by ID
	DeleteOrganization(id string) error
}
