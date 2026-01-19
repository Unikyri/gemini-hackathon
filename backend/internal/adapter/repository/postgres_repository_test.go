// Package repository_test provides integration tests for repository implementations.
package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"gemini-hackathon/internal/adapter/repository"
	"gemini-hackathon/internal/domain/entity"
	"gemini-hackathon/internal/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var testDB *gorm.DB

// TestMain sets up the test database connection
func TestMain(m *testing.M) {
	// Skip if not running integration tests
	if os.Getenv("INTEGRATION_TEST") != "true" {
		os.Exit(0)
	}

	// Setup database connection
	cfg := database.NewConfigFromEnv()
	db, err := database.NewConnection(cfg)
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}
	testDB = db

	// Run tests
	code := m.Run()

	// Cleanup
	sqlDB, _ := db.DB()
	sqlDB.Close()

	os.Exit(code)
}

// cleanupTestData removes test data from the database
func cleanupTestData(t *testing.T) {
	t.Helper()
	testDB.Exec("DELETE FROM path_nodes")
	testDB.Exec("DELETE FROM learning_paths")
}

// TestPathPersistenceRoundTrip verifies path persistence round-trip.
// Property: For any successfully created LearningPath, retrieving it via GetByID
// should return an equivalent structure with the same path_id, title, topic, and status.
func TestPathPersistenceRoundTrip(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test - no database connection")
	}

	cleanupTestData(t)
	pathRepo := repository.NewPostgresPathRepository(testDB)

	testCases := []struct {
		name string
		path *entity.LearningPath
	}{
		{
			name: "basic path",
			path: &entity.LearningPath{
				ID:        uuid.New().String(),
				UserID:    "anonymous",
				Topic:     "Learn Go basics",
				Title:     "Go Programming Fundamentals",
				Status:    entity.PathStatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "path with special characters",
			path: &entity.LearningPath{
				ID:        uuid.New().String(),
				UserID:    "user-123",
				Topic:     "Learn C++ templates & generics",
				Title:     "C++ Advanced: Templates, Generics & Metaprogramming",
				Status:    entity.PathStatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "completed path",
			path: &entity.LearningPath{
				ID:        uuid.New().String(),
				UserID:    "test-user",
				Topic:     "Python for beginners",
				Title:     "Python 101",
				Status:    entity.PathStatusCompleted,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// Create path
			err := pathRepo.Create(ctx, tc.path)
			if err != nil {
				t.Fatalf("Failed to create path: %v", err)
			}

			// Retrieve path
			retrieved, err := pathRepo.GetByID(ctx, tc.path.ID)
			if err != nil {
				t.Fatalf("Failed to retrieve path: %v", err)
			}
			if retrieved == nil {
				t.Fatal("Retrieved path is nil")
			}

			// Verify round-trip properties
			if retrieved.ID != tc.path.ID {
				t.Errorf("ID mismatch: got %s, want %s", retrieved.ID, tc.path.ID)
			}
			if retrieved.UserID != tc.path.UserID {
				t.Errorf("UserID mismatch: got %s, want %s", retrieved.UserID, tc.path.UserID)
			}
			if retrieved.Topic != tc.path.Topic {
				t.Errorf("Topic mismatch: got %s, want %s", retrieved.Topic, tc.path.Topic)
			}
			if retrieved.Title != tc.path.Title {
				t.Errorf("Title mismatch: got %s, want %s", retrieved.Title, tc.path.Title)
			}
			if retrieved.Status != tc.path.Status {
				t.Errorf("Status mismatch: got %s, want %s", retrieved.Status, tc.path.Status)
			}
		})
	}
}

// TestNodePersistenceRoundTrip tests that nodes can be created and retrieved correctly.
func TestNodePersistenceRoundTrip(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test - no database connection")
	}

	cleanupTestData(t)
	pathRepo := repository.NewPostgresPathRepository(testDB)
	nodeRepo := repository.NewPostgresNodeRepository(testDB)

	ctx := context.Background()

	// Create a path first
	pathID := uuid.New().String()
	path := &entity.LearningPath{
		ID:        pathID,
		UserID:    "test-user",
		Topic:     "Test topic",
		Title:     "Test Path",
		Status:    entity.PathStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := pathRepo.Create(ctx, path); err != nil {
		t.Fatalf("Failed to create path: %v", err)
	}

	// Create nodes
	nodes := []*entity.PathNode{
		{
			ID:              uuid.New().String(),
			PathID:          pathID,
			Position:        1,
			Title:           "First Exercise",
			Description:     "Introduction to basics",
			MarkdownContent: "# Exercise 1\nWrite a hello world program",
			BoilerplateCode: "func main() {\n\t// Your code here\n}",
			HiddenTests:     `[{"input": "", "expected": "Hello, World!"}]`,
			Status:          entity.NodeStatusUnlocked,
			XPReward:        100,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			PathID:          pathID,
			Position:        2,
			Title:           "Second Exercise",
			Description:     "Variables and types",
			MarkdownContent: "# Exercise 2\nDeclare variables",
			BoilerplateCode: "func main() {}",
			HiddenTests:     `[]`,
			Status:          entity.NodeStatusLocked,
			XPReward:        150,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	// Save nodes
	if err := nodeRepo.CreateBatch(ctx, nodes); err != nil {
		t.Fatalf("Failed to create nodes: %v", err)
	}

	// Retrieve nodes by path ID
	retrieved, err := nodeRepo.GetByPathID(ctx, pathID)
	if err != nil {
		t.Fatalf("Failed to retrieve nodes: %v", err)
	}

	// Verify count
	if len(retrieved) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(retrieved))
	}

	// Verify order (should be sorted by position)
	if retrieved[0].Position != 1 {
		t.Errorf("First node should have position 1, got %d", retrieved[0].Position)
	}
	if retrieved[1].Position != 2 {
		t.Errorf("Second node should have position 2, got %d", retrieved[1].Position)
	}
}

// TestPathWithNodesRoundTrip verifies GetByIDWithNodes returns correct data.
// Property: For any valid path with N nodes, GetByIDWithNodes should return
// exactly N nodes in the correct position order.
func TestPathWithNodesRoundTrip(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test - no database connection")
	}

	cleanupTestData(t)
	pathRepo := repository.NewPostgresPathRepository(testDB)
	nodeRepo := repository.NewPostgresNodeRepository(testDB)

	ctx := context.Background()

	// Create a path
	pathID := uuid.New().String()
	path := &entity.LearningPath{
		ID:        pathID,
		UserID:    "test-user",
		Topic:     "Full path test",
		Title:     "Complete Learning Path",
		Status:    entity.PathStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := pathRepo.Create(ctx, path); err != nil {
		t.Fatalf("Failed to create path: %v", err)
	}

	// Create 5 nodes
	nodes := make([]*entity.PathNode, 5)
	for i := 0; i < 5; i++ {
		status := entity.NodeStatusLocked
		if i == 0 {
			status = entity.NodeStatusUnlocked
		}
		nodes[i] = &entity.PathNode{
			ID:              uuid.New().String(),
			PathID:          pathID,
			Position:        i + 1,
			Title:           "Exercise " + string(rune('A'+i)),
			Description:     "Description for exercise",
			MarkdownContent: "# Content",
			BoilerplateCode: "// code",
			Status:          status,
			XPReward:        100 * (i + 1),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
	}

	if err := nodeRepo.CreateBatch(ctx, nodes); err != nil {
		t.Fatalf("Failed to create nodes: %v", err)
	}

	// Retrieve path with nodes
	retrieved, err := pathRepo.GetByIDWithNodes(ctx, pathID)
	if err != nil {
		t.Fatalf("Failed to retrieve path with nodes: %v", err)
	}

	// Verify node count
	if len(retrieved.Nodes) != 5 {
		t.Errorf("Expected 5 nodes, got %d", len(retrieved.Nodes))
	}

	// Verify nodes are in position order
	for i, node := range retrieved.Nodes {
		expectedPos := i + 1
		if node.Position != expectedPos {
			t.Errorf("Node at index %d should have position %d, got %d", i, expectedPos, node.Position)
		}
	}

	// Verify first node is unlocked
	if retrieved.Nodes[0].Status != entity.NodeStatusUnlocked {
		t.Errorf("First node should be unlocked, got %s", retrieved.Nodes[0].Status)
	}
}

// TestUnlockNextNode verifies the unlock logic
func TestUnlockNextNode(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test - no database connection")
	}

	cleanupTestData(t)
	pathRepo := repository.NewPostgresPathRepository(testDB)
	nodeRepo := repository.NewPostgresNodeRepository(testDB)

	ctx := context.Background()

	// Create path
	pathID := uuid.New().String()
	path := &entity.LearningPath{
		ID:        pathID,
		UserID:    "test-user",
		Topic:     "Unlock test",
		Title:     "Unlock Test Path",
		Status:    entity.PathStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := pathRepo.Create(ctx, path); err != nil {
		t.Fatalf("Failed to create path: %v", err)
	}

	// Create 3 nodes
	nodes := []*entity.PathNode{
		{ID: uuid.New().String(), PathID: pathID, Position: 1, Title: "Node 1", Status: entity.NodeStatusCompleted, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), PathID: pathID, Position: 2, Title: "Node 2", Status: entity.NodeStatusLocked, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), PathID: pathID, Position: 3, Title: "Node 3", Status: entity.NodeStatusLocked, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if err := nodeRepo.CreateBatch(ctx, nodes); err != nil {
		t.Fatalf("Failed to create nodes: %v", err)
	}

	// Unlock next after position 1
	if err := nodeRepo.UnlockNext(ctx, pathID, 1); err != nil {
		t.Fatalf("Failed to unlock next: %v", err)
	}

	// Verify node 2 is now unlocked
	node2, err := nodeRepo.GetByID(ctx, nodes[1].ID)
	if err != nil {
		t.Fatalf("Failed to get node 2: %v", err)
	}
	if node2.Status != entity.NodeStatusUnlocked {
		t.Errorf("Node 2 should be unlocked after UnlockNext, got %s", node2.Status)
	}

	// Verify node 3 is still locked
	node3, err := nodeRepo.GetByID(ctx, nodes[2].ID)
	if err != nil {
		t.Fatalf("Failed to get node 3: %v", err)
	}
	if node3.Status != entity.NodeStatusLocked {
		t.Errorf("Node 3 should still be locked, got %s", node3.Status)
	}
}
