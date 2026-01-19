---
name: Gin Framework
description: Crear handlers HTTP con Gin para el backend de Gemini Coding Path
---

# Skill: Gin Framework (Go)

## Contexto

Este proyecto usa **Gin** como framework HTTP para el backend.
Los handlers deben seguir patrones consistentes y trabajar con Clean Architecture.

---

## Estructura de un Controller

```go
// internal/adapter/http/path_controller.go
package http

import (
    "net/http"
    
    "gemini-hackathon/internal/adapter/http/dto"
    "gemini-hackathon/internal/usecase"
    "github.com/gin-gonic/gin"
)

type PathController struct {
    generatePathUC *usecase.GeneratePathUseCase
    getPathUC      *usecase.GetPathUseCase
}

func NewPathController(
    generatePathUC *usecase.GeneratePathUseCase,
    getPathUC *usecase.GetPathUseCase,
) *PathController {
    return &PathController{
        generatePathUC: generatePathUC,
        getPathUC:      getPathUC,
    }
}
```

---

## Patrones de Handlers

### POST con JSON Body
```go
func (c *PathController) GeneratePath(ctx *gin.Context) {
    // 1. Bind y validar request
    var req dto.GeneratePathRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error: "Invalid request body: " + err.Error(),
        })
        return
    }
    
    // 2. Obtener datos del contexto (ej: userID de auth)
    userID := ctx.GetString("user_id")
    if userID == "" {
        userID = "anonymous" // Para MVP sin auth
    }
    
    // 3. Ejecutar caso de uso
    path, err := c.generatePathUC.Execute(ctx.Request.Context(), userID, req.Prompt)
    if err != nil {
        // Mapear error a código HTTP apropiado
        statusCode := mapErrorToStatus(err)
        ctx.JSON(statusCode, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    // 4. Mapear respuesta
    response := dto.ToGeneratePathResponse(path)
    ctx.JSON(http.StatusOK, response)
}
```

### GET con Path Parameter
```go
func (c *PathController) GetPath(ctx *gin.Context) {
    // 1. Obtener path parameter
    pathID := ctx.Param("id")
    if pathID == "" {
        ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error: "path_id is required",
        })
        return
    }
    
    // 2. Ejecutar caso de uso
    path, err := c.getPathUC.Execute(ctx.Request.Context(), pathID)
    if err != nil {
        if errors.Is(err, usecase.ErrPathNotFound) {
            ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
                Error: "Path not found",
            })
            return
        }
        ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error: "Internal server error",
        })
        return
    }
    
    // 3. Responder
    ctx.JSON(http.StatusOK, dto.ToGetPathResponse(path))
}
```

### GET con Query Parameters
```go
func (c *PathController) ListPaths(ctx *gin.Context) {
    // Query params con valores por defecto
    userID := ctx.Query("user_id")
    status := ctx.DefaultQuery("status", "active")
    limit := ctx.DefaultQuery("limit", "10")
    
    // Parsear limit a int
    limitInt, err := strconv.Atoi(limit)
    if err != nil {
        limitInt = 10
    }
    
    // ... resto del handler
}
```

---

## DTOs (Data Transfer Objects)

```go
// internal/adapter/http/dto/path_dto.go
package dto

import (
    "gemini-hackathon/internal/domain/entity"
)

// Request DTOs
type GeneratePathRequest struct {
    Prompt    string `json:"prompt" binding:"required,min=3,max=500"`
    UserLevel string `json:"user_level,omitempty" binding:"omitempty,oneof=beginner intermediate advanced"`
}

// Response DTOs
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

// Mappers
func ToGeneratePathResponse(path *entity.LearningPath) GeneratePathResponse {
    nodes := make([]NodeSummary, len(path.Nodes))
    for i, n := range path.Nodes {
        nodes[i] = NodeSummary{
            ID:          n.ID,
            Title:       n.Title,
            Description: n.Description,
            Status:      string(n.Status),
            Position:    n.Position,
        }
    }
    
    return GeneratePathResponse{
        PathID: path.ID,
        Title:  path.Title,
        Nodes:  nodes,
    }
}
```

---

## Configuración del Router

```go
// cmd/api/routes.go
package main

import (
    "gemini-hackathon/internal/adapter/http"
    "github.com/gin-gonic/gin"
)

func SetupRouter(
    pathController *http.PathController,
    nodeController *http.NodeController,
) *gin.Engine {
    router := gin.Default()
    
    // Middleware global
    router.Use(corsMiddleware())
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API v1
    api := router.Group("/api/v1")
    {
        // Paths
        paths := api.Group("/paths")
        {
            paths.POST("/generate", pathController.GeneratePath)
            paths.GET("/:id", pathController.GetPath)
        }
        
        // Nodes
        nodes := api.Group("/nodes")
        {
            nodes.GET("/:node_id", nodeController.GetNodeDetail)
        }
    }
    
    return router
}
```

---

## Middleware

### CORS
```go
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}
```

### Auth (para futuro)
```go
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            // Por ahora, permitir sin auth (MVP)
            c.Set("user_id", "anonymous")
            c.Next()
            return
        }
        
        // Validar token y extraer user_id
        userID, err := validateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
                Error: "Invalid token",
            })
            c.Abort()
            return
        }
        
        c.Set("user_id", userID)
        c.Next()
    }
}
```

---

## Validación con Binding Tags

```go
type CreateSubmissionRequest struct {
    NodeID string `json:"node_id" binding:"required,uuid"`
    Code   string `json:"code" binding:"required,min=1,max=50000"`
}

// Tags disponibles:
// - required: campo obligatorio
// - min=N, max=N: longitud mínima/máxima
// - oneof=a b c: valor debe ser uno de los listados
// - uuid: debe ser UUID válido
// - email: debe ser email válido
// - url: debe ser URL válida
```

---

## Manejo de Errores

```go
// internal/adapter/http/error_handler.go
package http

import (
    "errors"
    "net/http"
    
    "gemini-hackathon/internal/usecase"
)

var errorStatusMap = map[error]int{
    usecase.ErrPathNotFound:    http.StatusNotFound,
    usecase.ErrNodeNotFound:    http.StatusNotFound,
    usecase.ErrInvalidPrompt:   http.StatusBadRequest,
    usecase.ErrAIUnavailable:   http.StatusServiceUnavailable,
}

func mapErrorToStatus(err error) int {
    for knownErr, status := range errorStatusMap {
        if errors.Is(err, knownErr) {
            return status
        }
    }
    return http.StatusInternalServerError
}
```

---

## Checklist para Nuevo Handler

- [ ] Controller inyecta use case, no repositorio directo
- [ ] Request DTO con binding tags para validación
- [ ] Response DTO separado de entidades de dominio
- [ ] Manejo de errores con códigos HTTP apropiados
- [ ] Context propagado al use case
- [ ] Mapper de entidad a DTO
- [ ] Ruta registrada en routes.go
