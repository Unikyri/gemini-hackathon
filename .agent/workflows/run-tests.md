---
description: Ejecutar tests del proyecto (backend y frontend)
---

# Workflow: Ejecutar Tests

## Backend (Go)

// turbo
### Ejecutar todos los tests
```bash
cd backend && go test ./... -v
```

// turbo
### Ejecutar tests con cobertura
```bash
cd backend && go test ./... -cover -coverprofile=coverage.out
```

// turbo
### Ver reporte de cobertura en HTML
```bash
cd backend && go tool cover -html=coverage.out -o coverage.html
```

### Ejecutar tests de un paquete específico
```bash
cd backend && go test ./internal/usecase/... -v
```

### Ejecutar un test específico
```bash
cd backend && go test ./internal/usecase/... -run TestNombreDelTest -v
```

---

## Frontend (React/Vite)

// turbo
### Ejecutar linting
```bash
cd frontend/core && npm run lint
```

// turbo
### Ejecutar build (verifica compilación)
```bash
cd frontend/core && npm run build
```

### Ejecutar tests (si están configurados con Vitest)
```bash
cd frontend/core && npm run test
```

---

## Pre-commit Hooks

// turbo
### Ejecutar todos los hooks manualmente
```bash
pre-commit run --all-files
```

### Ejecutar un hook específico
```bash
pre-commit run gofmt --all-files
pre-commit run eslint --all-files
```

---

## Property Tests

### Backend - Con gopter
```go
// Ejemplo de property test
func TestPromptValidation(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("prompts vacíos son rechazados", prop.ForAll(
        func(s string) bool {
            if strings.TrimSpace(s) == "" {
                _, err := validatePrompt(s)
                return err != nil
            }
            return true
        },
        gen.AnyString(),
    ))
    
    properties.TestingRun(t)
}
```

### Frontend - Con fast-check
```typescript
// Ejemplo de property test
import fc from 'fast-check';

test('boilerplate se inyecta correctamente', () => {
  fc.assert(
    fc.property(fc.string(), (code) => {
      const editor = mountEditor(code);
      expect(editor.getValue()).toBe(code);
    })
  );
});
```

---

## CI/CD Tests

Los tests se ejecutan automáticamente en GitHub Actions:

- **backend-ci.yml**: go test, golangci-lint
- **frontend-ci.yml**: npm run lint, npm run build

Para ver el estado: GitHub → Actions tab
