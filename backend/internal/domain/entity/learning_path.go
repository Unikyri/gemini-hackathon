// Package entity contains the domain business entities.
// Entities MUST NOT have dependencies on external frameworks.
package entity

import "time"

// PathStatus represents the status of a learning path
type PathStatus string

const (
	PathStatusActive    PathStatus = "active"
	PathStatusCompleted PathStatus = "completed"
	PathStatusArchived  PathStatus = "archived"
)

// LearningPath represents an AI-generated learning path
type LearningPath struct {
	ID        string     // UUID v4
	UserID    string     // User ID (anonymous in Sprint 1)
	Topic     string     // Original user prompt
	Title     string     // AI-generated title
	Status    PathStatus // Current path status
	CreatedAt time.Time
	UpdatedAt time.Time
	Nodes     []PathNode // Path nodes (loaded optionally)
}

// IsActive checks if the path is active
func (p *LearningPath) IsActive() bool {
	return p.Status == PathStatusActive
}

// MarkAsCompleted marks the path as completed
func (p *LearningPath) MarkAsCompleted() {
	p.Status = PathStatusCompleted
	p.UpdatedAt = time.Now()
}
