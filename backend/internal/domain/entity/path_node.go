package entity

import "time"

// NodeStatus representa el estado de un nodo dentro del path
type NodeStatus string

const (
	NodeStatusLocked    NodeStatus = "locked"
	NodeStatusUnlocked  NodeStatus = "unlocked"
	NodeStatusCompleted NodeStatus = "completed"
)

// PathNode representa un ejercicio dentro de un learning path
type PathNode struct {
	ID                   string     // UUID v4
	PathID               string     // FK a LearningPath
	Position             int        // Orden del nodo (1, 2, 3...)
	Title                string     // Título del ejercicio
	Description          string     // Descripción corta
	MarkdownContent      string     // Enunciado completo en Markdown
	BoilerplateCode      string     // Código inicial para el estudiante
	DocumentationSnippet string     // Mini-documentación relevante
	HiddenTests          string     // Tests ocultos en formato JSON
	Status               NodeStatus // Estado actual del nodo
	XPReward             int        // Puntos por completar
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// IsUnlocked verifica si el nodo está disponible para resolver
func (n *PathNode) IsUnlocked() bool {
	return n.Status == NodeStatusUnlocked
}

// IsCompleted verifica si el nodo ya fue completado
func (n *PathNode) IsCompleted() bool {
	return n.Status == NodeStatusCompleted
}

// Unlock desbloquea el nodo para que pueda ser resuelto
func (n *PathNode) Unlock() {
	if n.Status == NodeStatusLocked {
		n.Status = NodeStatusUnlocked
		n.UpdatedAt = time.Now()
	}
}

// Complete marca el nodo como completado
func (n *PathNode) Complete() {
	if n.Status == NodeStatusUnlocked {
		n.Status = NodeStatusCompleted
		n.UpdatedAt = time.Now()
	}
}
