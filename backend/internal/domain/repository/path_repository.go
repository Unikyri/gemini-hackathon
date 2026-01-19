// Package repository define las interfaces de persistencia.
// Las implementaciones concretas están en adapter/repository.
package repository

import (
	"context"

	"gemini-hackathon/internal/domain/entity"
)

// PathRepository define las operaciones de persistencia para LearningPath
type PathRepository interface {
	// Create guarda un nuevo path en la base de datos
	Create(ctx context.Context, path *entity.LearningPath) error

	// GetByID obtiene un path por su ID (sin nodos)
	GetByID(ctx context.Context, id string) (*entity.LearningPath, error)

	// GetByIDWithNodes obtiene un path con todos sus nodos
	GetByIDWithNodes(ctx context.Context, id string) (*entity.LearningPath, error)

	// GetByUserID obtiene todos los paths de un usuario
	GetByUserID(ctx context.Context, userID string) ([]*entity.LearningPath, error)

	// UpdateStatus actualiza el estado de un path
	UpdateStatus(ctx context.Context, id string, status entity.PathStatus) error
}
