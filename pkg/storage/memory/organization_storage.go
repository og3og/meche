package memory

import (
	"errors"
	"meche/pkg/models"
	"meche/pkg/storage"
	"sync"
)

// MemoryOrganizationStorage implements OrganizationStorage using in-memory storage
type MemoryOrganizationStorage struct {
	organizations map[string]*models.Organization
	mu            sync.RWMutex
}

// NewMemoryOrganizationStorage creates a new instance of MemoryOrganizationStorage
func NewMemoryOrganizationStorage() storage.OrganizationStorage {
	return &MemoryOrganizationStorage{
		organizations: make(map[string]*models.Organization),
	}
}

// CreateOrganization creates a new organization
func (s *MemoryOrganizationStorage) CreateOrganization(org *models.Organization) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.organizations[org.ID]; exists {
		return errors.New("organization already exists")
	}

	s.organizations[org.ID] = org
	return nil
}

// GetOrganization retrieves an organization by ID
func (s *MemoryOrganizationStorage) GetOrganization(id string) (*models.Organization, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	org, exists := s.organizations[id]
	if !exists {
		return nil, errors.New("organization not found")
	}

	// Don't return soft-deleted organizations
	if org.DeletedAt != nil {
		return nil, errors.New("organization not found")
	}

	return org, nil
}

// ListOrganizations returns all non-deleted organizations
func (s *MemoryOrganizationStorage) ListOrganizations() ([]*models.Organization, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orgs := make([]*models.Organization, 0, len(s.organizations))
	for _, org := range s.organizations {
		// Only include non-deleted organizations
		if org.DeletedAt == nil {
			orgs = append(orgs, org)
		}
	}

	return orgs, nil
}

// UpdateOrganization updates an existing organization
func (s *MemoryOrganizationStorage) UpdateOrganization(org *models.Organization) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.organizations[org.ID]
	if !exists {
		return errors.New("organization not found")
	}

	// Don't update soft-deleted organizations
	if existing.DeletedAt != nil {
		return errors.New("cannot update deleted organization")
	}

	s.organizations[org.ID] = org
	return nil
}

// DeleteOrganization soft deletes an organization by ID
func (s *MemoryOrganizationStorage) DeleteOrganization(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	org, exists := s.organizations[id]
	if !exists {
		return errors.New("organization not found")
	}

	// Don't delete already deleted organizations
	if org.DeletedAt != nil {
		return errors.New("organization already deleted")
	}

	org.MarkAsDeleted()
	return nil
}
