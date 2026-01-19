---
name: Backend Go - Clean Architecture
description: Genera código Go siguiendo Clean Architecture para Gemini Coding Path
---

# Skill: Backend Go

## Contexto

Este proyecto usa **Go 1.23+** con **Clean Architecture** para el backend API REST.
El código debe seguir las convenciones y patrones establecidos.

---

## Estructura de Capas

```
backend/
├── cmd/api/                    # Punto de entrada
│   ├── main.go                 # Inicialización
│   └── routes.go               # Registro de rutas
├── internal/
│   ├── domain/                 # CAPA DOMINIO (más interna)
│   │   ├── entity/             # Entidades de negocio
│   │   └── repository/         # Interfaces de repositorio
│   ├── usecase/                # CAPA APLICACIÓN
│   │   └── {nombre}_usecase.go # Casos de uso
│   ├── adapter/                # CAPA ADAPTADORES
│   │   ├── http/               # Controllers Gin
│   │   │   ├── dto/            # Request/Response DTOs
│   │   │   └── {nombre}_controller.go
│   │   └── repository/         # Implementaciones de repos
│   │       └── postgres_{nombre}_repository.go
│   └── infrastructure/         # CAPA INFRAESTRUCTURA (más externa)
│       ├── database/           # Conexión PostgreSQL
│       │   └── postgres.go
│       └── ai/                 # Cliente Vertex AI
│           └── vertex_client.go
```

---

## Convenciones de Nombrado

| Elemento | Convención | Ejemplo |
|----------|------------|---------|
| Archivos | snake_case | `learning_path.go` |
| Packages | lowercase | `entity`, `usecase` |
| Structs públicos | PascalCase | `LearningPath` |
| Campos públicos | PascalCase | `UserID` |
| Campos privados | camelCase | `createdAt` |
| Interfaces | PascalCase + sufijo | `PathRepository` |
| Funciones públicas | PascalCase | `GetByID` |
| Funciones privadas | camelCase | `validateInput` |
| Constantes | PascalCase | `StatusActive` |

---

## Patrones de Código

### Entidad de Dominio
```go
// internal/domain/entity/learning_path.go
package entity

import "time"

type PathStatus string

const (
    PathStatusActive    PathStatus = "active"
    PathStatusCompleted PathStatus = "completed"
    PathStatusArchived  PathStatus = "archived"
)

type LearningPath struct {
    ID         string
    UserID     string
    Topic      string
    Title      string
    Status     PathStatus
    AIMetadata map[string]interface{}
    CreatedAt  time.Time
    Nodes      []PathNode
}
```

### Interfaz de Repositorio
```go
// internal/domain/repository/path_repository.go
package repository

import (
    "context"
    "gemini-hackathon/internal/domain/entity"
)

type PathRepository interface {
    Create(ctx context.Context, path *entity.LearningPath) error
    GetByID(ctx context.Context, id string) (*entity.LearningPath, error)
    GetByUserID(ctx context.Context, userID string) ([]*entity.LearningPath, error)
}
```

### Caso de Uso
```go
// internal/usecase/generate_path_usecase.go
package usecase

import (
    "context"
    "errors"
    "gemini-hackathon/internal/domain/entity"
    "gemini-hackathon/internal/domain/repository"
)

type GeneratePathUseCase struct {
    pathRepo repository.PathRepository
    aiClient AIClient
}

func NewGeneratePathUseCase(
    pathRepo repository.PathRepository,
    aiClient AIClient,
) *GeneratePathUseCase {
    return &GeneratePathUseCase{
        pathRepo: pathRepo,
        aiClient: aiClient,
    }
}

func (uc *GeneratePathUseCase) Execute(
    ctx context.Context,
    userID string,
    prompt string,
) (*entity.LearningPath, error) {
    // Validación
    if len(prompt) < 3 || len(prompt) > 500 {
        return nil, errors.New("prompt must be between 3 and 500 characters")
    }
    
    // Lógica de negocio
    aiResponse, err := uc.aiClient.GeneratePath(ctx, prompt)
    if err != nil {
        return nil, err
    }
    
    // Crear entidad
    path := &entity.LearningPath{
        ID:     uuid.New().String(),
        UserID: userID,
        Topic:  prompt,
        Title:  aiResponse.Title,
        Status: entity.PathStatusActive,
    }
    
    // Persistir
    if err := uc.pathRepo.Create(ctx, path); err != nil {
        return nil, err
    }
    
    return path, nil
}
```

### Controller HTTP (Gin)
```go
// internal/adapter/http/path_controller.go
package http

import (
    "net/http"
    "gemini-hackathon/internal/usecase"
    "github.com/gin-gonic/gin"
)

type PathController struct {
    generatePathUC *usecase.GeneratePathUseCase
}

func NewPathController(uc *usecase.GeneratePathUseCase) *PathController {
    return &PathController{generatePathUC: uc}
}

// GeneratePath godoc
// @Summary Genera un nuevo path de aprendizaje
// @Tags paths
// @Accept json
// @Produce json
// @Param request body GeneratePathRequest true "Prompt del usuario"
// @Success 200 {object} GeneratePathResponse
// @Failure 400 {object} ErrorResponse
// @Router /paths/generate [post]
func (c *PathController) GeneratePath(ctx *gin.Context) {
    var req GeneratePathRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }
    
    // TODO: Obtener userID del contexto de auth
    userID := "anonymous"
    
    path, err := c.generatePathUC.Execute(ctx, userID, req.Prompt)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, toGeneratePathResponse(path))
}
```

### DTOs
```go
// internal/adapter/http/dto/path_dto.go
package dto

type GeneratePathRequest struct {
    Prompt    string `json:"prompt" binding:"required,min=3,max=500"`
    UserLevel string `json:"user_level,omitempty"`
}

type GeneratePathResponse struct {
    PathID string        `json:"path_id"`
    Title  string        `json:"title"`
    Nodes  []NodeSummary `json:"nodes"`
}

type NodeSummary struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
    Position    int    `json:"position"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
```

---

## Manejo de Errores

```go
// Siempre manejar errores explícitamente
result, err := someFunction()
if err != nil {
    return nil, fmt.Errorf("failed to do X: %w", err)
}

// Errores de dominio personalizados
var (
    ErrPathNotFound = errors.New("path not found")
    ErrInvalidPrompt = errors.New("invalid prompt")
)
```

---

## Inyección de Dependencias

```go
// cmd/api/main.go
func main() {
    // Inicializar infraestructura
    db := database.NewPostgresConnection()
    aiClient := ai.NewVertexClient()
    
    // Inicializar repositorios
    pathRepo := repository.NewPostgresPathRepository(db)
    nodeRepo := repository.NewPostgresNodeRepository(db)
    
    // Inicializar casos de uso
    generatePathUC := usecase.NewGeneratePathUseCase(pathRepo, aiClient)
    
    // Inicializar controllers
    pathController := http.NewPathController(generatePathUC)
    
    // Registrar rutas
    router := gin.Default()
    routes.SetupRoutes(router, pathController)
    
    router.Run(":8080")
}
```

---

## Checklist para Nuevo Código

- [ ] Sigue estructura de capas (Domain → UseCase → Adapter → Infra)
- [ ] No hay imports de capas externas en capas internas
- [ ] Usa interfaces para dependencias
- [ ] Maneja todos los errores explícitamente
- [ ] Nombrado consistente con convenciones
- [ ] DTOs separados de entidades de dominio
- [ ] Tests unitarios para casos de uso
