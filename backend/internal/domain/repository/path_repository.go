// Package repository defines the persistence interfaces.
// Concrete implementations are in adapter/repository.
package repository

import (
	"context"

	"gemini-hackathon/internal/domain/entity"
)

// PathRepository defines persistence operations for LearningPath
type PathRepository interface {
	// Create saves a new path to the database
	Create(ctx context.Context, path *entity.LearningPath) error

	// GetByID retrieves a path by its ID (without nodes)
	GetByID(ctx context.Context, id string) (*entity.LearningPath, error)

	// GetByIDWithNodes retrieves a path with all its nodes
	GetByIDWithNodes(ctx context.Context, id string) (*entity.LearningPath, error)

	// GetByUserID retrieves all paths for a user
	GetByUserID(ctx context.Context, userID string) ([]*entity.LearningPath, error)

	// UpdateStatus updates the status of a path
	UpdateStatus(ctx context.Context, id string, status entity.PathStatus) error
}
