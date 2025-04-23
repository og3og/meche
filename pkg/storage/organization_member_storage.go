package storage

import (
	"meche/pkg/models"
)

// OrganizationMemberStorage defines the interface for organization member storage operations
type OrganizationMemberStorage interface {
	// CreateMember creates a new organization membership
	CreateMember(member *models.OrganizationMember) error

	// GetMember retrieves a membership by ID
	GetMember(id string) (*models.OrganizationMember, error)

	// GetMembersByOrganization returns all members of an organization
	GetMembersByOrganization(orgID string) ([]*models.OrganizationMember, error)

	// GetMembersByUser returns all organizations a user is a member of
	GetMembersByUser(userID string) ([]*models.OrganizationMember, error)

	// UpdateMember updates an existing membership
	UpdateMember(member *models.OrganizationMember) error

	// DeleteMember soft deletes a membership
	DeleteMember(id string) error
}
