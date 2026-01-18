---
name: PostgreSQL Database
description: Queries, migraciones y patrones de base de datos para Gemini Coding Path
---

# Skill: PostgreSQL

## Contexto

Este proyecto usa **PostgreSQL 16** como base de datos principal.
El esquema sigue el modelo definido en la wiki (Database.md).

---

## Esquema de la Base de Datos

```sql
-- Habilitar extensión para UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabla de usuarios
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    username VARCHAR(50),
    is_guest BOOLEAN DEFAULT FALSE,
    total_xp INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de rutas de aprendizaje
CREATE TABLE learning_paths (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    topic TEXT NOT NULL,
    title VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active',
    ai_metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índice para buscar paths por usuario
CREATE INDEX idx_learning_paths_user_id ON learning_paths(user_id);

-- Tabla de nodos (ejercicios)
CREATE TABLE path_nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    path_id UUID REFERENCES learning_paths(id) ON DELETE CASCADE,
    position INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255),
    description_md TEXT,
    boilerplate_code TEXT,
    tests_payload JSONB NOT NULL DEFAULT '{}',
    documentation_md TEXT,
    status VARCHAR(20) DEFAULT 'locked',
    xp_reward INT DEFAULT 100,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índice para buscar nodos por path
CREATE INDEX idx_path_nodes_path_id ON path_nodes(path_id);
-- Índice único para orden dentro del path
CREATE UNIQUE INDEX idx_path_nodes_position ON path_nodes(path_id, position);

-- Tabla de submissions
CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    node_id UUID REFERENCES path_nodes(id) ON DELETE CASCADE,
    submitted_code TEXT NOT NULL,
    is_passed BOOLEAN DEFAULT FALSE,
    execution_logs TEXT,
    ai_feedback TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para submissions
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_node_id ON submissions(node_id);
```

---

## Conexión desde Go

### Configuración con GORM
```go
// internal/infrastructure/database/postgres.go
package database

import (
    "fmt"
    "os"
    "time"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewPostgresConnection() (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Configurar pool de conexiones
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db, nil
}

// Health check
func HealthCheck(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Ping()
}
```

### Modelos GORM
```go
// internal/infrastructure/database/models.go
package database

import (
    "time"
    
    "github.com/lib/pq"
    "gorm.io/datatypes"
)

type UserModel struct {
    ID           string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Email        *string   `gorm:"type:varchar(255);uniqueIndex"`
    PasswordHash *string   `gorm:"type:varchar(255)"`
    Username     *string   `gorm:"type:varchar(50)"`
    IsGuest      bool      `gorm:"default:false"`
    TotalXP      int       `gorm:"default:0"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (UserModel) TableName() string {
    return "users"
}

type LearningPathModel struct {
    ID         string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    UserID     string         `gorm:"type:uuid;index"`
    Topic      string         `gorm:"type:text;not null"`
    Title      string         `gorm:"type:varchar(255)"`
    Status     string         `gorm:"type:varchar(20);default:active"`
    AIMetadata datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
    CreatedAt  time.Time      `gorm:"autoCreateTime"`
    
    // Relación
    User  UserModel       `gorm:"foreignKey:UserID"`
    Nodes []PathNodeModel `gorm:"foreignKey:PathID"`
}

func (LearningPathModel) TableName() string {
    return "learning_paths"
}

type PathNodeModel struct {
    ID              string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    PathID          string         `gorm:"type:uuid;index"`
    Position        int            `gorm:"not null"`
    Title           string         `gorm:"type:varchar(255);not null"`
    Slug            string         `gorm:"type:varchar(255)"`
    DescriptionMD   string         `gorm:"type:text;column:description_md"`
    BoilerplateCode string         `gorm:"type:text"`
    TestsPayload    datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
    DocumentationMD string         `gorm:"type:text;column:documentation_md"`
    Status          string         `gorm:"type:varchar(20);default:locked"`
    XPReward        int            `gorm:"default:100"`
    CreatedAt       time.Time      `gorm:"autoCreateTime"`
}

func (PathNodeModel) TableName() string {
    return "path_nodes"
}
```

---

## Implementación de Repositorios

```go
// internal/adapter/repository/postgres_path_repository.go
package repository

import (
    "context"
    "errors"
    
    "gemini-hackathon/internal/domain/entity"
    "gemini-hackathon/internal/infrastructure/database"
    "gorm.io/gorm"
)

type PostgresPathRepository struct {
    db *gorm.DB
}

func NewPostgresPathRepository(db *gorm.DB) *PostgresPathRepository {
    return &PostgresPathRepository{db: db}
}

func (r *PostgresPathRepository) Create(ctx context.Context, path *entity.LearningPath) error {
    model := toPathModel(path)
    
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Insertar path
        if err := tx.Create(&model).Error; err != nil {
            return err
        }
        
        // Insertar nodos
        nodeModels := toNodeModels(path.Nodes, model.ID)
        if len(nodeModels) > 0 {
            if err := tx.Create(&nodeModels).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}

func (r *PostgresPathRepository) GetByID(ctx context.Context, id string) (*entity.LearningPath, error) {
    var model database.LearningPathModel
    
    err := r.db.WithContext(ctx).
        Preload("Nodes", func(db *gorm.DB) *gorm.DB {
            return db.Order("position ASC")
        }).
        First(&model, "id = ?", id).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrPathNotFound
        }
        return nil, err
    }
    
    return toPathEntity(&model), nil
}

func (r *PostgresPathRepository) GetByUserID(ctx context.Context, userID string) ([]*entity.LearningPath, error) {
    var models []database.LearningPathModel
    
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&models).Error
    
    if err != nil {
        return nil, err
    }
    
    paths := make([]*entity.LearningPath, len(models))
    for i, m := range models {
        paths[i] = toPathEntity(&m)
    }
    
    return paths, nil
}

// Mappers
func toPathModel(e *entity.LearningPath) database.LearningPathModel {
    return database.LearningPathModel{
        ID:     e.ID,
        UserID: e.UserID,
        Topic:  e.Topic,
        Title:  e.Title,
        Status: string(e.Status),
    }
}

func toPathEntity(m *database.LearningPathModel) *entity.LearningPath {
    nodes := make([]entity.PathNode, len(m.Nodes))
    for i, n := range m.Nodes {
        nodes[i] = *toNodeEntity(&n)
    }
    
    return &entity.LearningPath{
        ID:        m.ID,
        UserID:    m.UserID,
        Topic:     m.Topic,
        Title:     m.Title,
        Status:    entity.PathStatus(m.Status),
        CreatedAt: m.CreatedAt,
        Nodes:     nodes,
    }
}
```

---

## Queries Útiles

### Obtener progreso de un usuario
```sql
SELECT 
    lp.id AS path_id,
    lp.title,
    COUNT(pn.id) AS total_nodes,
    COUNT(CASE WHEN pn.status = 'completed' THEN 1 END) AS completed_nodes,
    ROUND(COUNT(CASE WHEN pn.status = 'completed' THEN 1 END)::numeric / 
          COUNT(pn.id) * 100, 2) AS progress_percent
FROM learning_paths lp
JOIN path_nodes pn ON pn.path_id = lp.id
WHERE lp.user_id = $1
GROUP BY lp.id, lp.title;
```

### Obtener último intento de un usuario en un nodo
```sql
SELECT * FROM submissions
WHERE user_id = $1 AND node_id = $2
ORDER BY created_at DESC
LIMIT 1;
```

### Desbloquear siguiente nodo
```sql
UPDATE path_nodes
SET status = 'unlocked'
WHERE path_id = $1 
  AND position = (
      SELECT position + 1 
      FROM path_nodes 
      WHERE id = $2
  );
```

---

## Migraciones

### Con golang-migrate
```bash
# Instalar
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Crear migración
migrate create -ext sql -dir migrations -seq create_users_table

# Aplicar migraciones
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up

# Rollback
migrate -path migrations -database "..." down 1
```

### Archivo de migración ejemplo
```sql
-- migrations/000001_create_users_table.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    username VARCHAR(50),
    is_guest BOOLEAN DEFAULT FALSE,
    total_xp INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

---

## Checklist

- [ ] Esquema DDL creado con todos los índices
- [ ] Pool de conexiones configurado
- [ ] Health check implementado
- [ ] Modelos GORM con tags correctos
- [ ] Repositorios usan transacciones donde corresponde
- [ ] Migraciones versionadas
