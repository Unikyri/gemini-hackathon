---
name: DevOps - Docker & CI/CD
description: Configuración de Docker, GitHub Actions y pipelines para Gemini Coding Path
---

# Skill: DevOps

## Contexto

Este proyecto usa **Docker** para containerización y **GitHub Actions** para CI/CD.
Incluye configuración para el sandbox seguro de ejecución de código.

---

## Docker - Backend

### Dockerfile para API
```dockerfile
# backend/Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar go.mod y go.sum primero (cache de dependencias)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# --- Imagen final ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
```

### Docker Compose para Desarrollo
```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: hackathon-postgres
    environment:
      POSTGRES_DB: hackathon_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD:-localdev}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: hackathon-backend
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=hackathon_db
      - DB_USER=postgres
      - DB_PASSWORD=${DB_PASSWORD:-localdev}
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres_data:
```

---

## Docker - Code Runner (Sandbox)

### Dockerfile Seguro para Ejecución de Código
```dockerfile
# sandbox/Dockerfile.runner
FROM golang:1.23-alpine

# Crear usuario no-root
RUN addgroup -S runner && adduser -S runner -G runner

# Directorio de trabajo
WORKDIR /sandbox

# Cambiar a usuario no-root
USER runner

# El código se monta en /sandbox/code
# Los tests se ejecutan con: go test ./...
```

### Docker Compose para Sandbox
```yaml
# docker-compose.sandbox.yml
version: '3.8'

services:
  code-runner:
    build:
      context: ./sandbox
      dockerfile: Dockerfile.runner
    container_name: code-runner
    # SIN acceso a red
    network_mode: none
    # Límites de recursos
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 128M
        reservations:
          cpus: '0.25'
          memory: 64M
    # Solo lectura del sistema de archivos (excepto /tmp y /sandbox)
    read_only: true
    tmpfs:
      - /tmp:size=10M
    # Timeout de 30 segundos
    stop_grace_period: 30s
    # Sin privilegios
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
```

---

## GitHub Actions - CI

### Backend CI
```yaml
# .github/workflows/backend-ci.yml
name: Backend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/**'
  pull_request:
    branches: [main, develop]
    paths:
      - 'backend/**'

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache-dependency-path: backend/go.sum
      
      - name: Install dependencies
        working-directory: ./backend
        run: go mod download
      
      - name: Run go vet
        working-directory: ./backend
        run: go vet ./...
      
      - name: Install golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          working-directory: ./backend
      
      - name: Run tests
        working-directory: ./backend
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./backend/coverage.out
          fail_ci_if_error: false
```

### Frontend CI
```yaml
# .github/workflows/frontend-ci.yml
name: Frontend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'frontend/**'
  pull_request:
    branches: [main, develop]
    paths:
      - 'frontend/**'

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: frontend/core/package-lock.json
      
      - name: Install dependencies
        working-directory: ./frontend/core
        run: npm ci
      
      - name: Run linter
        working-directory: ./frontend/core
        run: npm run lint
      
      - name: Type check
        working-directory: ./frontend/core
        run: npm run type-check
      
      - name: Build
        working-directory: ./frontend/core
        run: npm run build
        env:
          VITE_API_BASE_URL: https://api.example.com
```

---

## Pre-commit Hooks

```yaml
# .pre-commit-config.yaml
repos:
  # Hooks generales
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-added-large-files
        args: ['--maxkb=1000']
      - id: detect-private-key

  # Detección de secretos
  - repo: https://github.com/Yelp/detect-secrets
    rev: v1.4.0
    hooks:
      - id: detect-secrets
        args: ['--baseline', '.secrets.baseline']

  # Go
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
      - id: go-vet
      - id: go-imports
      - id: go-mod-tidy

  # ESLint para JS/TS
  - repo: https://github.com/pre-commit/mirrors-eslint
    rev: v8.56.0
    hooks:
      - id: eslint
        files: \.[jt]sx?$
        types: [file]
        additional_dependencies:
          - eslint
          - typescript
```

---

## Comandos Útiles

```bash
# Levantar base de datos local
docker compose up -d postgres

# Levantar todo el stack
docker compose up -d

# Ver logs
docker compose logs -f backend

# Ejecutar pre-commit
pre-commit run --all-files

# Build de imagen
docker build -t gemini-backend ./backend
```

---

## Checklist

- [ ] Dockerfile multi-stage para imagen pequeña
- [ ] Docker Compose para desarrollo local
- [ ] Sandbox con limits de CPU/memoria
- [ ] Sandbox sin acceso a red
- [ ] Usuario no-root en contenedores
- [ ] CI para backend (go test, golangci-lint)
- [ ] CI para frontend (lint, type-check, build)
- [ ] Pre-commit hooks configurados
