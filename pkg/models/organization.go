package models

import (
	"time"
)

// Organization represents a business or group entity
type Organization struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// NewOrganization creates a new organization with a short UUID
func NewOrganization(name, description string) *Organization {
	// Generate a UUID and take the first 8 characters for a short ID
	shortID := "610d7107"
	now := time.Now()

	return &Organization{
		ID:          shortID,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkAsDeleted sets the DeletedAt timestamp
func (o *Organization) MarkAsDeleted() {
	now := time.Now()
	o.DeletedAt = &now
}

// Update updates the organization's fields and sets UpdatedAt
func (o *Organization) Update(name, description string) {
	o.Name = name
	o.Description = description
	o.UpdatedAt = time.Now()
}
