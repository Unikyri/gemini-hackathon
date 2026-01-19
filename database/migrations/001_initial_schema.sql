-- Initial database schema for Gemini Coding Path
-- This file runs automatically on first PostgreSQL container startup

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (for Sprint 2)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    xp_total INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Learning paths table
CREATE TABLE IF NOT EXISTS learning_paths (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,  -- 'anonymous' in Sprint 1
    topic TEXT NOT NULL,
    title VARCHAR(500) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for faster user lookups
CREATE INDEX IF NOT EXISTS idx_learning_paths_user_id ON learning_paths(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_paths_status ON learning_paths(status);

-- Path nodes table
CREATE TABLE IF NOT EXISTS path_nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    position INT NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    markdown_content TEXT,
    boilerplate_code TEXT,
    documentation_snippet TEXT,
    hidden_tests TEXT,  -- JSON format
    status VARCHAR(50) DEFAULT 'locked',
    xp_reward INT DEFAULT 100,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for path_nodes
CREATE INDEX IF NOT EXISTS idx_path_nodes_path_id ON path_nodes(path_id);
CREATE INDEX IF NOT EXISTS idx_path_nodes_position ON path_nodes(path_id, position);

-- Submissions table (for Sprint 2)
CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    node_id UUID NOT NULL REFERENCES path_nodes(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    code TEXT NOT NULL,
    passed BOOLEAN DEFAULT FALSE,
    execution_time_ms INT,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for submissions
CREATE INDEX IF NOT EXISTS idx_submissions_node_id ON submissions(node_id);
CREATE INDEX IF NOT EXISTS idx_submissions_user_id ON submissions(user_id);

-- Insert anonymous user for Sprint 1
-- (This way we can use a real user_id instead of just 'anonymous' string)
-- INSERT INTO users (id, email, password_hash, display_name)
-- VALUES ('00000000-0000-0000-0000-000000000000', 'anonymous@local', '', 'Anonymous')
-- ON CONFLICT DO NOTHING;

COMMENT ON TABLE learning_paths IS 'AI-generated learning paths for users';
COMMENT ON TABLE path_nodes IS 'Individual exercises within a learning path';
COMMENT ON TABLE submissions IS 'User code submissions for exercises';
