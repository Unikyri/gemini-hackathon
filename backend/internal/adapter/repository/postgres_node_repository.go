package repository

import (
	"context"
	"errors"

	"gemini-hackathon/internal/domain/entity"
	"gemini-hackathon/internal/domain/repository"
	"gemini-hackathon/internal/infrastructure/database"

	"gorm.io/gorm"
)

// Ensure PostgresNodeRepository implements repository.NodeRepository
var _ repository.NodeRepository = (*PostgresNodeRepository)(nil)

// PostgresNodeRepository is a PostgreSQL implementation of NodeRepository
type PostgresNodeRepository struct {
	db *gorm.DB
}

// NewPostgresNodeRepository creates a new PostgresNodeRepository
func NewPostgresNodeRepository(db *gorm.DB) *PostgresNodeRepository {
	return &PostgresNodeRepository{db: db}
}

// CreateBatch saves multiple nodes at once (transactional)
func (r *PostgresNodeRepository) CreateBatch(ctx context.Context, nodes []*entity.PathNode) error {
	if len(nodes) == 0 {
		return nil
	}

	models := make([]database.PathNodeModel, len(nodes))
	for i, node := range nodes {
		models[i] = *database.FromNodeEntity(node)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&models).Error
	})
}

// GetByID retrieves a node by its ID
func (r *PostgresNodeRepository) GetByID(ctx context.Context, id string) (*entity.PathNode, error) {
	var model database.PathNodeModel
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return model.ToEntity(), nil
}

// GetByPathID retrieves all nodes for a path ordered by position
func (r *PostgresNodeRepository) GetByPathID(ctx context.Context, pathID string) ([]*entity.PathNode, error) {
	var models []database.PathNodeModel
	result := r.db.WithContext(ctx).
		Where("path_id = ?", pathID).
		Order("position ASC").
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	nodes := make([]*entity.PathNode, len(models))
	for i, model := range models {
		nodes[i] = model.ToEntity()
	}
	return nodes, nil
}

// UpdateStatus updates the status of a node
func (r *PostgresNodeRepository) UpdateStatus(ctx context.Context, id string, status entity.NodeStatus) error {
	result := r.db.WithContext(ctx).
		Model(&database.PathNodeModel{}).
		Where("id = ?", id).
		Update("status", string(status))
	return result.Error
}

// UnlockNext unlocks the next node after completing one
func (r *PostgresNodeRepository) UnlockNext(ctx context.Context, pathID string, currentPosition int) error {
	nextPosition := currentPosition + 1
	result := r.db.WithContext(ctx).
		Model(&database.PathNodeModel{}).
		Where("path_id = ? AND position = ? AND status = ?", pathID, nextPosition, string(entity.NodeStatusLocked)).
		Update("status", string(entity.NodeStatusUnlocked))
	return result.Error
}
