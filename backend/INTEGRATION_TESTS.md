# Integration Test Runner for Gemini Coding Path

## Prerequisites
- Docker installed and running
- Go 1.23+ installed

## Quick Start

### 1. Start the database
```bash
docker compose -f docker-compose.dev.yml up -d
```

### 2. Wait for database to be ready
```bash
docker compose -f docker-compose.dev.yml ps
# Should show "healthy" status
```

### 3. Run integration tests
```bash
cd backend
INTEGRATION_TEST=true go test -v ./internal/adapter/repository/...
```

### 4. Stop the database (optional)
```bash
docker compose -f docker-compose.dev.yml down
```

### 5. Stop and remove data (clean start)
```bash
docker compose -f docker-compose.dev.yml down -v
```

## Environment Variables

The tests use these environment variables (with defaults):

| Variable | Default | Description |
|----------|---------|-------------|
| `INTEGRATION_TEST` | `false` | Set to `true` to run integration tests |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `gemini_hackathon` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |

## Running Specific Tests

```bash
# Run all repository tests
INTEGRATION_TEST=true go test -v ./internal/adapter/repository/...

# Run specific test
INTEGRATION_TEST=true go test -v -run TestPathPersistenceRoundTrip ./internal/adapter/repository/...

# Run with race detection
INTEGRATION_TEST=true go test -v -race ./internal/adapter/repository/...
```
