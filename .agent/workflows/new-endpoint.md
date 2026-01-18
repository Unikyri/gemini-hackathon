---
description: Crear un nuevo endpoint REST siguiendo Clean Architecture en Go
---

# Workflow: Nuevo Endpoint

## Prerrequisitos
- Backend configurado con estructura Clean Architecture
- Base de datos PostgreSQL corriendo
- Go 1.21+ instalado

## Pasos

### 1. Definir el contrato en la documentación
```bash
# Editar el archivo de API docs
docs/gemini-hackathon.wiki/API-Docs.md
```
- Agregar método HTTP, ruta, request body y response esperado
- Definir códigos de error posibles

### 2. Crear/actualizar entidad de dominio
```bash
# Crear archivo en:
backend/internal/domain/entity/nombre_entidad.go
```
```go
package entity

type NombreEntidad struct {
    ID        string    `json:"id"`
    // ... otros campos
    CreatedAt time.Time `json:"created_at"`
}
```

### 3. Definir interfaz de repositorio
```bash
# Crear/editar en:
backend/internal/domain/repository/nombre_repository.go
```
```go
package repository

type NombreRepository interface {
    Create(ctx context.Context, entity *entity.NombreEntidad) error
    GetByID(ctx context.Context, id string) (*entity.NombreEntidad, error)
    // ... otros métodos
}
```

### 4. Implementar el caso de uso
```bash
# Crear archivo en:
backend/internal/usecase/nombre_usecase.go
```
```go
package usecase

type NombreUseCase struct {
    repo repository.NombreRepository
}

func NewNombreUseCase(repo repository.NombreRepository) *NombreUseCase {
    return &NombreUseCase{repo: repo}
}

func (uc *NombreUseCase) Execute(ctx context.Context, input InputDTO) (*OutputDTO, error) {
    // Validaciones
    // Lógica de negocio
    // Llamada al repositorio
    return output, nil
}
```

### 5. Implementar el repositorio con PostgreSQL
```bash
# Crear archivo en:
backend/internal/adapter/repository/postgres_nombre_repository.go
```

### 6. Crear el controlador HTTP
```bash
# Crear archivo en:
backend/internal/adapter/http/nombre_controller.go
```
```go
package http

type NombreController struct {
    useCase *usecase.NombreUseCase
}

func (c *NombreController) Handler(ctx *gin.Context) {
    // Parse request
    // Call use case
    // Return response
}
```

### 7. Registrar la ruta
```bash
# Editar:
backend/cmd/api/routes.go
```
```go
api := router.Group("/api/v1")
{
    api.POST("/ruta", controller.Handler)
}
```

### 8. Escribir tests
```bash
# Crear test en:
backend/internal/usecase/nombre_usecase_test.go
```

// turbo
### 9. Verificar que compila
```bash
cd backend && go build ./...
```

// turbo
### 10. Ejecutar tests
```bash
cd backend && go test ./...
```

## Checklist Final
- [ ] Contrato documentado en API-Docs.md
- [ ] Entidad creada en domain/entity
- [ ] Interfaz de repositorio definida
- [ ] Caso de uso implementado con validaciones
- [ ] Repositorio PostgreSQL implementado
- [ ] Controller HTTP creado
- [ ] Ruta registrada
- [ ] Tests escritos y pasando
