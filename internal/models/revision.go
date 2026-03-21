package models

import "time"

// Revision represents a project revision (document/spec).
// Stored as YAML files in .watchfire/revisions/ directory.
type Revision struct {
	Version        int       `yaml:"version"`
	RevisionNumber int       `yaml:"revision_number"` // Sequential, user-facing
	Title          string    `yaml:"title"`
	Content        string    `yaml:"content"`  // Markdown
	Complete       bool      `yaml:"complete"` // All tasks done
	Position       int       `yaml:"position"`
	CreatedAt      time.Time `yaml:"created_at"`
	UpdatedAt      time.Time `yaml:"updated_at"`
}

// NewRevision creates a new revision with default values.
func NewRevision(revisionNumber int, title, content string) *Revision {
	now := time.Now().UTC()
	return &Revision{
		Version:        1,
		RevisionNumber: revisionNumber,
		Title:          title,
		Content:        content,
		Complete:       false,
		Position:       revisionNumber,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
