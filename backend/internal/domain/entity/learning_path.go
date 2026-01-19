// Package entity contiene las entidades de dominio del negocio.
// Las entidades NO deben tener dependencias de frameworks externos.
package entity

import "time"

// PathStatus representa el estado de un learning path
type PathStatus string

const (
	PathStatusActive    PathStatus = "active"
	PathStatusCompleted PathStatus = "completed"
	PathStatusArchived  PathStatus = "archived"
)

// LearningPath representa una ruta de aprendizaje generada por IA
type LearningPath struct {
	ID        string     // UUID v4
	UserID    string     // ID del usuario (anonymous en Sprint 1)
	Topic     string     // Prompt original del usuario
	Title     string     // Título generado por IA
	Status    PathStatus // Estado actual del path
	CreatedAt time.Time
	UpdatedAt time.Time
	Nodes     []PathNode // Nodos del path (cargados opcionalmente)
}

// IsActive verifica si el path está activo
func (p *LearningPath) IsActive() bool {
	return p.Status == PathStatusActive
}

// MarkAsCompleted marca el path como completado
func (p *LearningPath) MarkAsCompleted() {
	p.Status = PathStatusCompleted
	p.UpdatedAt = time.Now()
}
