package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/watchfire-io/watchfire/internal/models"
)

// LoadRevision loads a revision from its YAML file.
func LoadRevision(projectPath string, revisionNumber int) (*models.Revision, error) {
	path := RevisionFile(projectPath, revisionNumber)

	if !FileExists(path) {
		return nil, nil
	}

	var rev models.Revision
	if err := LoadYAML(path, &rev); err != nil {
		return nil, err
	}
	return &rev, nil
}

// SaveRevision saves a revision to its YAML file.
func SaveRevision(projectPath string, rev *models.Revision) error {
	if err := os.MkdirAll(ProjectRevisionsDir(projectPath), 0o755); err != nil {
		return err
	}
	return SaveYAML(RevisionFile(projectPath, rev.RevisionNumber), rev)
}

// DeleteRevisionFile permanently deletes a revision file.
func DeleteRevisionFile(projectPath string, revisionNumber int) error {
	path := RevisionFile(projectPath, revisionNumber)
	if !FileExists(path) {
		return nil
	}
	return os.Remove(path)
}

// LoadAllRevisions loads all revisions from a project's revisions directory.
func LoadAllRevisions(projectPath string) ([]*models.Revision, error) {
	revisionsDir := ProjectRevisionsDir(projectPath)

	if !FileExists(revisionsDir) {
		return []*models.Revision{}, nil
	}

	entries, err := os.ReadDir(revisionsDir)
	if err != nil {
		return nil, err
	}

	var revisions []*models.Revision
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") {
			continue
		}

		numStr := strings.TrimSuffix(name, ".yaml")
		revNum, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}

		rev, err := LoadRevision(projectPath, revNum)
		if err != nil {
			return nil, err
		}
		if rev != nil {
			revisions = append(revisions, rev)
		}
	}

	return revisions, nil
}

// LoadActiveRevisions loads all non-complete revisions from a project.
func LoadActiveRevisions(projectPath string) ([]*models.Revision, error) {
	revisions, err := LoadAllRevisions(projectPath)
	if err != nil {
		return nil, err
	}

	var active []*models.Revision
	for _, r := range revisions {
		if !r.Complete {
			active = append(active, r)
		}
	}
	return active, nil
}

// SyncNextRevisionNumber scans the revisions directory and updates next_revision_number
// in project.yaml if it's behind the highest existing revision file.
func SyncNextRevisionNumber(projectPath string) error {
	project, err := LoadProject(projectPath)
	if err != nil || project == nil {
		return err
	}

	revisionsDir := ProjectRevisionsDir(projectPath)
	if !FileExists(revisionsDir) {
		return nil
	}

	entries, err := os.ReadDir(revisionsDir)
	if err != nil {
		return err
	}

	highest := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") {
			continue
		}
		numStr := strings.TrimSuffix(name, ".yaml")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}
		if num > highest {
			highest = num
		}
	}

	needed := highest + 1
	if needed > project.NextRevisionNumber {
		project.NextRevisionNumber = needed
		return SaveProject(projectPath, project)
	}
	return nil
}
