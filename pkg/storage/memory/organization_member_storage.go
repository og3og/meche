package memory

import (
	"errors"
	"meche/pkg/models"
	"meche/pkg/storage"
	"sync"
)

// MemoryOrganizationMemberStorage implements OrganizationMemberStorage using in-memory storage
type MemoryOrganizationMemberStorage struct {
	members map[string]*models.OrganizationMember
	mu      sync.RWMutex
}

// NewMemoryOrganizationMemberStorage creates a new instance of MemoryOrganizationMemberStorage
func NewMemoryOrganizationMemberStorage() storage.OrganizationMemberStorage {
	return &MemoryOrganizationMemberStorage{
		members: make(map[string]*models.OrganizationMember),
	}
}

// CreateMember creates a new organization membership
func (s *MemoryOrganizationMemberStorage) CreateMember(member *models.OrganizationMember) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.members[member.ID]; exists {
		return errors.New("membership already exists")
	}

	s.members[member.ID] = member
	return nil
}

// GetMember retrieves a membership by ID
func (s *MemoryOrganizationMemberStorage) GetMember(id string) (*models.OrganizationMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	member, exists := s.members[id]
	if !exists || member.DeletedAt != nil {
		return nil, errors.New("membership not found")
	}

	return member, nil
}

// GetMembersByOrganization returns all members of an organization
func (s *MemoryOrganizationMemberStorage) GetMembersByOrganization(orgID string) ([]*models.OrganizationMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	members := make([]*models.OrganizationMember, 0)
	for _, member := range s.members {
		if member.OrganizationID == orgID && member.DeletedAt == nil {
			members = append(members, member)
		}
	}

	return members, nil
}

// GetMembersByUser returns all organizations a user is a member of
func (s *MemoryOrganizationMemberStorage) GetMembersByUser(userID string) ([]*models.OrganizationMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	members := make([]*models.OrganizationMember, 0)
	for _, member := range s.members {
		if member.UserID == userID && member.DeletedAt == nil {
			members = append(members, member)
		}
	}

	return members, nil
}

// UpdateMember updates an existing membership
func (s *MemoryOrganizationMemberStorage) UpdateMember(member *models.OrganizationMember) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.members[member.ID]
	if !exists || existing.DeletedAt != nil {
		return errors.New("membership not found")
	}

	s.members[member.ID] = member
	return nil
}

// DeleteMember soft deletes a membership
func (s *MemoryOrganizationMemberStorage) DeleteMember(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	member, exists := s.members[id]
	if !exists || member.DeletedAt != nil {
		return errors.New("membership not found")
	}

	member.MarkAsDeleted()
	return nil
}
