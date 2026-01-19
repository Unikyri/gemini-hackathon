// Package repository provides PostgreSQL implementations of domain repositories.
package repository

import (
	"context"
	"errors"

	"gemini-hackathon/internal/domain/entity"
	"gemini-hackathon/internal/domain/repository"
	"gemini-hackathon/internal/infrastructure/database"

	"gorm.io/gorm"
)

// Ensure PostgresPathRepository implements repository.PathRepository
var _ repository.PathRepository = (*PostgresPathRepository)(nil)

// PostgresPathRepository is a PostgreSQL implementation of PathRepository
type PostgresPathRepository struct {
	db *gorm.DB
}

// NewPostgresPathRepository creates a new PostgresPathRepository
func NewPostgresPathRepository(db *gorm.DB) *PostgresPathRepository {
	return &PostgresPathRepository{db: db}
}

// Create saves a new path to the database
func (r *PostgresPathRepository) Create(ctx context.Context, path *entity.LearningPath) error {
	model := database.FromPathEntity(path)
	result := r.db.WithContext(ctx).Create(model)
	return result.Error
}

// GetByID retrieves a path by its ID (without nodes)
func (r *PostgresPathRepository) GetByID(ctx context.Context, id string) (*entity.LearningPath, error) {
	var model database.LearningPathModel
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return model.ToEntity(), nil
}

// GetByIDWithNodes retrieves a path with all its nodes
func (r *PostgresPathRepository) GetByIDWithNodes(ctx context.Context, id string) (*entity.LearningPath, error) {
	var model database.LearningPathModel
	result := r.db.WithContext(ctx).
		Preload("Nodes", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("id = ?", id).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return model.ToEntity(), nil
}

// GetByUserID retrieves all paths for a user
func (r *PostgresPathRepository) GetByUserID(ctx context.Context, userID string) ([]*entity.LearningPath, error) {
	var models []database.LearningPathModel
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	paths := make([]*entity.LearningPath, len(models))
	for i, model := range models {
		paths[i] = model.ToEntity()
	}
	return paths, nil
}

// UpdateStatus updates the status of a path
func (r *PostgresPathRepository) UpdateStatus(ctx context.Context, id string, status entity.PathStatus) error {
	result := r.db.WithContext(ctx).
		Model(&database.LearningPathModel{}).
		Where("id = ?", id).
		Update("status", string(status))
	return result.Error
}
