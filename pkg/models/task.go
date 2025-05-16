package models

import "time"

type Task struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OrgID       string     `json:"org_id"`
	ProjectID   string     `json:"project_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
	PinnedAt    *time.Time `json:"pinned_at,omitempty"`
	DueAt       *time.Time `json:"due_at,omitempty"`
}

func NewTask(name, description, orgID, projectID string) *Task {
	now := time.Now()
	return &Task{
		ID:          generateID(),
		Name:        name,
		Description: description,
		OrgID:       orgID,
		ProjectID:   projectID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) Update(name, description string) {
	t.Name = name
	t.Description = description
	t.UpdatedAt = time.Now()
}

func (t *Task) SetDueDate(dueAt *time.Time) {
	t.DueAt = dueAt
	t.UpdatedAt = time.Now()
}

func (t *Task) Delete() {
	now := time.Now()
	t.DeletedAt = &now
}

func (t *Task) Archive() {
	now := time.Now()
	t.ArchivedAt = &now
	t.UpdatedAt = now
}

func (t *Task) Unarchive() {
	t.ArchivedAt = nil
	t.UpdatedAt = time.Now()
}

func (t *Task) Pin() {
	now := time.Now()
	t.PinnedAt = &now
	t.UpdatedAt = now
}

func (t *Task) Unpin() {
	t.PinnedAt = nil
	t.UpdatedAt = time.Now()
}
