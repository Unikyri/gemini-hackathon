package entity

import "time"

// NodeStatus represents the status of a node within the path
type NodeStatus string

const (
	NodeStatusLocked    NodeStatus = "locked"
	NodeStatusUnlocked  NodeStatus = "unlocked"
	NodeStatusCompleted NodeStatus = "completed"
)

// PathNode represents an exercise within a learning path
type PathNode struct {
	ID                   string     // UUID v4
	PathID               string     // FK to LearningPath
	Position             int        // Node order (1, 2, 3...)
	Title                string     // Exercise title
	Description          string     // Short description
	MarkdownContent      string     // Full exercise statement in Markdown
	BoilerplateCode      string     // Initial code for the student
	DocumentationSnippet string     // Relevant mini-documentation
	HiddenTests          string     // Hidden tests in JSON format
	Status               NodeStatus // Current node status
	XPReward             int        // Points earned on completion
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// IsUnlocked checks if the node is available to solve
func (n *PathNode) IsUnlocked() bool {
	return n.Status == NodeStatusUnlocked
}

// IsCompleted checks if the node has been completed
func (n *PathNode) IsCompleted() bool {
	return n.Status == NodeStatusCompleted
}

// Unlock unlocks the node so it can be solved
func (n *PathNode) Unlock() {
	if n.Status == NodeStatusLocked {
		n.Status = NodeStatusUnlocked
		n.UpdatedAt = time.Now()
	}
}

// Complete marks the node as completed
func (n *PathNode) Complete() {
	if n.Status == NodeStatusUnlocked {
		n.Status = NodeStatusCompleted
		n.UpdatedAt = time.Now()
	}
}
