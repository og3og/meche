package models

import "time"

type Project struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OrgID       string     `json:"org_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
	PinnedAt    *time.Time `json:"pinned_at,omitempty"`
}

func NewProject(name, description, orgID string) *Project {
	now := time.Now()
	return &Project{
		ID:          generateID(),
		Name:        name,
		Description: description,
		OrgID:       orgID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Project) Update(name, description string) {
	p.Name = name
	p.Description = description
	p.UpdatedAt = time.Now()
}

func (p *Project) Delete() {
	now := time.Now()
	p.DeletedAt = &now
}

func (p *Project) Archive() {
	now := time.Now()
	p.ArchivedAt = &now
	p.UpdatedAt = now
}

func (p *Project) Unarchive() {
	p.ArchivedAt = nil
	p.UpdatedAt = time.Now()
}

func (p *Project) Pin() {
	now := time.Now()
	p.PinnedAt = &now
	p.UpdatedAt = now
}

func (p *Project) Unpin() {
	p.PinnedAt = nil
	p.UpdatedAt = time.Now()
}
