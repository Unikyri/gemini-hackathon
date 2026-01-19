package repository

import (
	"context"

	"gemini-hackathon/internal/domain/entity"
)

// NodeRepository defines persistence operations for PathNode
type NodeRepository interface {
	// CreateBatch saves multiple nodes at once (transactional)
	CreateBatch(ctx context.Context, nodes []*entity.PathNode) error

	// GetByID retrieves a node by its ID
	GetByID(ctx context.Context, id string) (*entity.PathNode, error)

	// GetByPathID retrieves all nodes for a path ordered by position
	GetByPathID(ctx context.Context, pathID string) ([]*entity.PathNode, error)

	// UpdateStatus updates the status of a node
	UpdateStatus(ctx context.Context, id string, status entity.NodeStatus) error

	// UnlockNext unlocks the next node after completing one
	UnlockNext(ctx context.Context, pathID string, currentPosition int) error
}
