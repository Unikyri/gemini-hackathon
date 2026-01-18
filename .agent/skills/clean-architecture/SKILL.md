---
name: Clean Architecture
description: Estructurar código por capas para Gemini Coding Path
---

# Skill: Clean Architecture

## Contexto

Gemini Coding Path sigue **Clean Architecture** en el backend.
La clave es la **regla de dependencia**: las capas internas NO conocen a las externas.

---

## Diagrama de Capas

```
┌─────────────────────────────────────────────────────────────┐
│                    FRAMEWORKS & DRIVERS                      │
│  (Gin HTTP, PostgreSQL, Vertex AI SDK, Docker)              │
│                                                              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │               INTERFACE ADAPTERS                       │  │
│  │    (Controllers, Repositories Impl, Presenters)       │  │
│  │                                                        │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │           APPLICATION LAYER                      │  │  │
│  │  │    (Use Cases / Interactors / Services)         │  │  │
│  │  │                                                  │  │  │
│  │  │  ┌───────────────────────────────────────────┐  │  │  │
│  │  │  │           DOMAIN LAYER                     │  │  │  │
│  │  │  │    (Entities, Repository Interfaces)      │  │  │  │
│  │  │  │                                           │  │  │  │
│  │  │  └───────────────────────────────────────────┘  │  │  │
│  │  │               ▲ NO dependencies                  │  │  │
│  │  └───────────────│─────────────────────────────────┘  │  │
│  │                  │ depends on                          │  │
│  └──────────────────│────────────────────────────────────┘  │
│                     │ depends on                             │
└─────────────────────│───────────────────────────────────────┘
                      │ depends on
                      ▼
               EXTERNAL WORLD
```

---

## Estructura de Directorios

```
backend/
├── cmd/
│   └── api/
│       ├── main.go          # Punto de entrada, DI
│       └── routes.go        # Registro de rutas
│
├── internal/
│   ├── domain/              # CAPA 1: DOMINIO (más interna)
│   │   ├── entity/          # Entidades de negocio
│   │   │   ├── learning_path.go
│   │   │   ├── path_node.go
│   │   │   └── user.go
│   │   └── repository/      # Interfaces de repositorio
│   │       ├── path_repository.go
│   │       └── node_repository.go
│   │
│   ├── usecase/             # CAPA 2: APLICACIÓN
│   │   ├── generate_path.go
│   │   ├── get_path.go
│   │   ├── get_node_detail.go
│   │   └── interfaces.go    # Interfaces de servicios externos
│   │
│   ├── adapter/             # CAPA 3: ADAPTADORES
│   │   ├── http/            # Controllers HTTP
│   │   │   ├── path_controller.go
│   │   │   ├── node_controller.go
│   │   │   └── dto/         # Request/Response DTOs
│   │   │       └── path_dto.go
│   │   └── repository/      # Implementaciones de repo
│   │       ├── postgres_path_repository.go
│   │       └── postgres_node_repository.go
│   │
│   └── infrastructure/      # CAPA 4: INFRAESTRUCTURA (más externa)
│       ├── database/
│       │   └── postgres.go
│       └── ai/
│           └── vertex_client.go
```

---

## Regla de Dependencias

```go
// ✅ CORRECTO: Domain no importa nada externo
// internal/domain/entity/learning_path.go
package entity

import "time"  // Solo librería estándar

type LearningPath struct {
    ID        string
    Title     string
    CreatedAt time.Time
}

// ✅ CORRECTO: Domain define interfaces (contratos)
// internal/domain/repository/path_repository.go
package repository

import (
    "context"
    "gemini-hackathon/internal/domain/entity"  // Solo domain
)

type PathRepository interface {
    Create(ctx context.Context, path *entity.LearningPath) error
    GetByID(ctx context.Context, id string) (*entity.LearningPath, error)
}
```

```go
// ✅ CORRECTO: UseCase depende de Domain
// internal/usecase/generate_path.go
package usecase

import (
    "context"
    "gemini-hackathon/internal/domain/entity"      // ✓ Domain
    "gemini-hackathon/internal/domain/repository"  // ✓ Domain
)

// NO importar: adapter, infrastructure, gin, gorm, etc.

type GeneratePathUseCase struct {
    pathRepo repository.PathRepository  // Interface de domain
    aiClient AIClient                   // Interface definida en usecase
}
```

```go
// ✅ CORRECTO: Adapter depende de Domain y UseCase
// internal/adapter/http/path_controller.go
package http

import (
    "gemini-hackathon/internal/usecase"  // ✓ UseCase
    "github.com/gin-gonic/gin"           // ✓ Framework en adapter
)
```

```go
// ✅ CORRECTO: Infrastructure implementa interfaces
// internal/adapter/repository/postgres_path_repository.go
package repository

import (
    "gemini-hackathon/internal/domain/entity"      // ✓
    "gemini-hackathon/internal/domain/repository"  // ✓ Implementa esta interfaz
    "gorm.io/gorm"                                 // ✓ Framework en infra
)

type PostgresPathRepository struct {
    db *gorm.DB
}

// Implementa repository.PathRepository
func (r *PostgresPathRepository) Create(ctx context.Context, path *entity.LearningPath) error {
    // ... implementación con GORM
}
```

---

## Flujo de una Request

```
1. HTTP Request
        │
        ▼
2. [Controller] Parsea request, valida formato
        │
        ▼
3. [Use Case] Ejecuta lógica de negocio
        │
        ├──▶ [AI Client] Llama a Vertex AI
        │
        └──▶ [Repository] Persiste en PostgreSQL
                    │
                    ▼
4. [Controller] Mapea a DTO, responde HTTP
        │
        ▼
5. HTTP Response
```

### Código del Flujo
```go
// 1. Request llega al router (cmd/api/routes.go)
api.POST("/paths/generate", pathController.GeneratePath)

// 2. Controller (adapter/http/path_controller.go)
func (c *PathController) GeneratePath(ctx *gin.Context) {
    var req dto.GeneratePathRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    // 3. Llama al Use Case
    path, err := c.generatePathUC.Execute(ctx, req.Prompt)
    
    // 4. Mapea respuesta
    ctx.JSON(200, dto.ToGeneratePathResponse(path))
}

// 3. Use Case (usecase/generate_path.go)
func (uc *GeneratePathUseCase) Execute(ctx context.Context, prompt string) (*entity.LearningPath, error) {
    // Validación de negocio
    if len(prompt) < 3 {
        return nil, ErrInvalidPrompt
    }
    
    // Llama a IA
    aiResponse, err := uc.aiClient.GeneratePath(ctx, prompt)
    if err != nil {
        return nil, err
    }
    
    // Crea entidad
    path := &entity.LearningPath{
        ID:    uuid.New().String(),
        Title: aiResponse.Title,
    }
    
    // Persiste
    if err := uc.pathRepo.Create(ctx, path); err != nil {
        return nil, err
    }
    
    return path, nil
}
```

---

## Inyección de Dependencias

```go
// cmd/api/main.go
func main() {
    // 1. Inicializar infraestructura
    db, _ := database.NewPostgresConnection()
    aiClient := ai.NewVertexClient(os.Getenv("GOOGLE_CLOUD_PROJECT"))
    
    // 2. Crear implementaciones de repositorios
    pathRepo := repository.NewPostgresPathRepository(db)
    nodeRepo := repository.NewPostgresNodeRepository(db)
    
    // 3. Crear use cases con dependencias inyectadas
    generatePathUC := usecase.NewGeneratePathUseCase(pathRepo, nodeRepo, aiClient)
    getPathUC := usecase.NewGetPathUseCase(pathRepo)
    
    // 4. Crear controllers con use cases
    pathController := http.NewPathController(generatePathUC, getPathUC)
    
    // 5. Configurar router
    router := gin.Default()
    routes.Setup(router, pathController)
    
    router.Run(":8080")
}
```

---

## Beneficios

| Beneficio | Descripción |
|-----------|-------------|
| **Testeable** | Use cases se prueban con mocks de repos/clients |
| **Mantenible** | Cambiar de PostgreSQL a MongoDB solo afecta infra |
| **Escalable** | Agregar nuevas features no afecta el core |
| **Independiente** | Domain no depende de frameworks |

---

## Errores Comunes

### ❌ Use Case que depende de Gin
```go
// MAL: Use case acoplado a HTTP
func (uc *GeneratePathUseCase) Execute(ctx *gin.Context, prompt string) {
    // El use case no debería saber nada de Gin
}
```

### ❌ Entidad con tags de GORM
```go
// MAL: Domain acoplado a infraestructura
type LearningPath struct {
    ID    string `gorm:"primaryKey"`  // NO en domain
    Title string `gorm:"type:varchar(255)"`
}
```

### ✅ Correcto: Modelos de infra separados
```go
// Domain: entidad pura
type LearningPath struct {
    ID    string
    Title string
}

// Infra: modelo con tags de framework
type LearningPathModel struct {
    ID    string `gorm:"type:uuid;primaryKey"`
    Title string `gorm:"type:varchar(255)"`
}

// Mapper entre domain e infra
func toEntity(m *LearningPathModel) *entity.LearningPath { ... }
func toModel(e *entity.LearningPath) *LearningPathModel { ... }
```

---

## Checklist

- [ ] Domain NO importa paquetes externos (solo std lib)
- [ ] Interfaces de repos/clients definidas en domain/usecase
- [ ] Use cases NO conocen Gin, GORM u otros frameworks
- [ ] DTOs separados de entidades de domain
- [ ] Mappers entre modelos de infra y entidades de domain
- [ ] Dependencias inyectadas en main.go
