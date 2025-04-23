package models

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationRole represents the possible roles a user can have in an organization
type OrganizationRole string

const (
	RoleOwner  OrganizationRole = "owner"
	RoleAdmin  OrganizationRole = "admin"
	RoleMember OrganizationRole = "member"
)

// OrganizationMember represents a user's membership in an organization
type OrganizationMember struct {
	ID             string           `json:"id"`
	OrganizationID string           `json:"organization_id"`
	UserID         string           `json:"user_id"`
	Role           OrganizationRole `json:"role"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      *time.Time       `json:"deleted_at,omitempty"`
}

// NewOrganizationMember creates a new organization membership
func NewOrganizationMember(orgID, userID string, role OrganizationRole) *OrganizationMember {
	now := time.Now()
	return &OrganizationMember{
		ID:             uuid.New().String()[:8],
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// MarkAsDeleted sets the DeletedAt timestamp
func (m *OrganizationMember) MarkAsDeleted() {
	now := time.Now()
	m.DeletedAt = &now
}
