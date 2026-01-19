// Package database provides PostgreSQL connection and GORM models.
package database

import (
	"time"

	"gemini-hackathon/internal/domain/entity"
)

// LearningPathModel is the GORM model for learning_paths table
type LearningPathModel struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:varchar(255);index"`
	Topic     string    `gorm:"type:text"`
	Title     string    `gorm:"type:varchar(500)"`
	Status    string    `gorm:"type:varchar(50);default:'active'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Nodes     []PathNodeModel `gorm:"foreignKey:PathID;references:ID"`
}

// TableName returns the table name for GORM
func (LearningPathModel) TableName() string {
	return "learning_paths"
}

// ToEntity converts the GORM model to a domain entity
func (m *LearningPathModel) ToEntity() *entity.LearningPath {
	path := &entity.LearningPath{
		ID:        m.ID,
		UserID:    m.UserID,
		Topic:     m.Topic,
		Title:     m.Title,
		Status:    entity.PathStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	// Convert nodes if loaded
	if len(m.Nodes) > 0 {
		path.Nodes = make([]entity.PathNode, len(m.Nodes))
		for i, node := range m.Nodes {
			path.Nodes[i] = *node.ToEntity()
		}
	}

	return path
}

// FromPathEntity converts a domain entity to a GORM model
func FromPathEntity(e *entity.LearningPath) *LearningPathModel {
	return &LearningPathModel{
		ID:        e.ID,
		UserID:    e.UserID,
		Topic:     e.Topic,
		Title:     e.Title,
		Status:    string(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// PathNodeModel is the GORM model for path_nodes table
type PathNodeModel struct {
	ID                   string    `gorm:"type:uuid;primaryKey"`
	PathID               string    `gorm:"type:uuid;index"`
	Position             int       `gorm:"type:int"`
	Title                string    `gorm:"type:varchar(500)"`
	Description          string    `gorm:"type:text"`
	MarkdownContent      string    `gorm:"type:text"`
	BoilerplateCode      string    `gorm:"type:text"`
	DocumentationSnippet string    `gorm:"type:text"`
	HiddenTests          string    `gorm:"type:text"`
	Status               string    `gorm:"type:varchar(50);default:'locked'"`
	XPReward             int       `gorm:"type:int;default:100"`
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM
func (PathNodeModel) TableName() string {
	return "path_nodes"
}

// ToEntity converts the GORM model to a domain entity
func (m *PathNodeModel) ToEntity() *entity.PathNode {
	return &entity.PathNode{
		ID:                   m.ID,
		PathID:               m.PathID,
		Position:             m.Position,
		Title:                m.Title,
		Description:          m.Description,
		MarkdownContent:      m.MarkdownContent,
		BoilerplateCode:      m.BoilerplateCode,
		DocumentationSnippet: m.DocumentationSnippet,
		HiddenTests:          m.HiddenTests,
		Status:               entity.NodeStatus(m.Status),
		XPReward:             m.XPReward,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

// FromNodeEntity converts a domain entity to a GORM model
func FromNodeEntity(e *entity.PathNode) *PathNodeModel {
	return &PathNodeModel{
		ID:                   e.ID,
		PathID:               e.PathID,
		Position:             e.Position,
		Title:                e.Title,
		Description:          e.Description,
		MarkdownContent:      e.MarkdownContent,
		BoilerplateCode:      e.BoilerplateCode,
		DocumentationSnippet: e.DocumentationSnippet,
		HiddenTests:          e.HiddenTests,
		Status:               string(e.Status),
		XPReward:             e.XPReward,
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
	}
}
