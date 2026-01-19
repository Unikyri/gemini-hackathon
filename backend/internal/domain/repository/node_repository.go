package repository

import (
	"context"

	"gemini-hackathon/internal/domain/entity"
)

// NodeRepository define las operaciones de persistencia para PathNode
type NodeRepository interface {
	// CreateBatch guarda múltiples nodos de una vez (transaccional)
	CreateBatch(ctx context.Context, nodes []*entity.PathNode) error

	// GetByID obtiene un nodo por su ID
	GetByID(ctx context.Context, id string) (*entity.PathNode, error)

	// GetByPathID obtiene todos los nodos de un path ordenados por posición
	GetByPathID(ctx context.Context, pathID string) ([]*entity.PathNode, error)

	// UpdateStatus actualiza el estado de un nodo
	UpdateStatus(ctx context.Context, id string, status entity.NodeStatus) error

	// UnlockNext desbloquea el siguiente nodo después de completar uno
	UnlockNext(ctx context.Context, pathID string, currentPosition int) error
}
