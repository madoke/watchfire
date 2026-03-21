// Package revision handles revision management for the daemon.
package revision

import (
	"fmt"
	"sort"
	"time"

	"github.com/watchfire-io/watchfire/internal/config"
	"github.com/watchfire-io/watchfire/internal/models"
)

// Manager handles revision operations.
type Manager struct{}

// NewManager creates a new revision manager.
func NewManager() *Manager {
	return &Manager{}
}

// CreateOptions contains options for creating a revision.
type CreateOptions struct {
	Title   string
	Content string
}

// UpdateOptions contains options for updating a revision.
type UpdateOptions struct {
	RevisionNumber int
	Title          *string
	Content        *string
	Complete       *bool
}

// ListRevisions returns all revisions for a project, sorted by position.
func (m *Manager) ListRevisions(projectPath string) ([]*models.Revision, error) {
	revisions, err := config.LoadAllRevisions(projectPath)
	if err != nil {
		return nil, err
	}

	sort.Slice(revisions, func(i, j int) bool {
		if revisions[i].Position != revisions[j].Position {
			return revisions[i].Position < revisions[j].Position
		}
		return revisions[i].RevisionNumber < revisions[j].RevisionNumber
	})

	return revisions, nil
}

// GetRevision retrieves a revision by number.
func (m *Manager) GetRevision(projectPath string, revisionNumber int) (*models.Revision, error) {
	rev, err := config.LoadRevision(projectPath, revisionNumber)
	if err != nil {
		return nil, err
	}
	if rev == nil {
		return nil, fmt.Errorf("revision not found: %d", revisionNumber)
	}
	return rev, nil
}

// CreateRevision creates a new revision.
func (m *Manager) CreateRevision(projectPath string, opts CreateOptions) (*models.Revision, error) {
	_ = config.SyncNextRevisionNumber(projectPath)

	project, err := config.LoadProject(projectPath)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found: %s", projectPath)
	}

	revisionNumber := project.NextRevisionNumber
	rev := models.NewRevision(revisionNumber, opts.Title, opts.Content)

	if err := config.SaveRevision(projectPath, rev); err != nil {
		return nil, err
	}

	project.NextRevisionNumber++
	project.UpdatedAt = time.Now().UTC()
	if err := config.SaveProject(projectPath, project); err != nil {
		return nil, err
	}

	return rev, nil
}

// UpdateRevision updates an existing revision.
func (m *Manager) UpdateRevision(projectPath string, opts UpdateOptions) (*models.Revision, error) {
	rev, err := config.LoadRevision(projectPath, opts.RevisionNumber)
	if err != nil {
		return nil, err
	}
	if rev == nil {
		return nil, fmt.Errorf("revision not found: %d", opts.RevisionNumber)
	}

	if opts.Title != nil {
		rev.Title = *opts.Title
	}
	if opts.Content != nil {
		rev.Content = *opts.Content
	}
	if opts.Complete != nil {
		rev.Complete = *opts.Complete
	}

	rev.UpdatedAt = time.Now().UTC()

	if err := config.SaveRevision(projectPath, rev); err != nil {
		return nil, err
	}

	return rev, nil
}

// DeleteRevision permanently deletes a revision.
func (m *Manager) DeleteRevision(projectPath string, revisionNumber int) error {
	return config.DeleteRevisionFile(projectPath, revisionNumber)
}

