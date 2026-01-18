---
name: Testing - Property & Unit Tests
description: Escribir tests unitarios y property-based tests para Gemini Coding Path
---

# Skill: Testing

## Contexto

Este proyecto usa:
- **Go**: `testing` estándar + `gopter` para property tests
- **TypeScript**: Vitest (opcional) + `fast-check` para property tests

---

## Backend (Go)

### Unit Test Básico
```go
// internal/usecase/generate_path_usecase_test.go
package usecase_test

import (
    "context"
    "testing"
    
    "gemini-hackathon/internal/domain/entity"
    "gemini-hackathon/internal/usecase"
)

// Mock del repositorio
type mockPathRepo struct {
    paths map[string]*entity.LearningPath
}

func (m *mockPathRepo) Create(ctx context.Context, path *entity.LearningPath) error {
    m.paths[path.ID] = path
    return nil
}

func (m *mockPathRepo) GetByID(ctx context.Context, id string) (*entity.LearningPath, error) {
    if path, ok := m.paths[id]; ok {
        return path, nil
    }
    return nil, usecase.ErrPathNotFound
}

// Mock del cliente AI
type mockAIClient struct {
    response *usecase.AIPathResponse
    err      error
}

func (m *mockAIClient) GeneratePath(ctx context.Context, prompt string) (*usecase.AIPathResponse, error) {
    return m.response, m.err
}

func TestGeneratePathUseCase_Execute_Success(t *testing.T) {
    // Arrange
    repo := &mockPathRepo{paths: make(map[string]*entity.LearningPath)}
    aiClient := &mockAIClient{
        response: &usecase.AIPathResponse{
            Title: "Learn Go",
            Nodes: []usecase.AINode{
                {Title: "Node 1", Description: "Desc 1"},
                {Title: "Node 2", Description: "Desc 2"},
                {Title: "Node 3", Description: "Desc 3"},
                {Title: "Node 4", Description: "Desc 4"},
                {Title: "Node 5", Description: "Desc 5"},
            },
        },
    }
    
    uc := usecase.NewGeneratePathUseCase(repo, aiClient)
    
    // Act
    path, err := uc.Execute(context.Background(), "user-1", "Learn Go")
    
    // Assert
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if path.Title != "Learn Go" {
        t.Errorf("expected title 'Learn Go', got %s", path.Title)
    }
    if len(path.Nodes) != 5 {
        t.Errorf("expected 5 nodes, got %d", len(path.Nodes))
    }
}

func TestGeneratePathUseCase_Execute_EmptyPrompt(t *testing.T) {
    // Arrange
    uc := usecase.NewGeneratePathUseCase(nil, nil)
    
    // Act
    _, err := uc.Execute(context.Background(), "user-1", "")
    
    // Assert
    if err == nil {
        t.Fatal("expected error for empty prompt")
    }
}
```

### Property Test con gopter
```go
// internal/usecase/generate_path_usecase_property_test.go
package usecase_test

import (
    "strings"
    "testing"
    
    "github.com/leanovate/gopter"
    "github.com/leanovate/gopter/gen"
    "github.com/leanovate/gopter/prop"
    
    "gemini-hackathon/internal/usecase"
)

// Property 1: Prompt Validation
// For any string input as prompt, if the string is empty or contains 
// only whitespace characters, the Backend_API should reject it with 
// a 400 error; otherwise, it should accept the prompt for processing.
func TestProperty_PromptValidation(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("empty or whitespace prompts are rejected", prop.ForAll(
        func(s string) bool {
            err := usecase.ValidatePrompt(s)
            
            isEmptyOrWhitespace := strings.TrimSpace(s) == ""
            
            if isEmptyOrWhitespace {
                // Debe retornar error
                return err != nil
            }
            // Prompts válidos solo rechazan por longitud
            if len(s) < 3 || len(s) > 500 {
                return err != nil
            }
            return err == nil
        },
        gen.AnyString(),
    ))
    
    properties.TestingRun(t)
}

// Property 2: AI Response Structure Integrity
// For any successful response from Vertex_AI, the parsed response 
// should contain exactly 5 nodes with valid fields.
func TestProperty_AIResponseStructure(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    // Generador de respuestas de IA
    nodeGen := gen.Struct(reflect.TypeOf(usecase.AINode{}), map[string]gopter.Gen{
        "Title":       gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
        "Description": gen.AlphaString(),
    })
    
    responseGen := gen.Struct(reflect.TypeOf(usecase.AIPathResponse{}), map[string]gopter.Gen{
        "Title": gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
        "Nodes": gen.SliceOfN(5, nodeGen),
    })
    
    properties.Property("parsed response has exactly 5 nodes", prop.ForAll(
        func(resp usecase.AIPathResponse) bool {
            err := usecase.ValidateAIResponse(&resp)
            if err != nil {
                return false
            }
            return len(resp.Nodes) == 5
        },
        responseGen,
    ))
    
    properties.TestingRun(t)
}

// Property 3: Path Persistence Round-Trip
func TestProperty_PathPersistenceRoundTrip(t *testing.T) {
    // Este test requiere una BD de test
    // Se ejecuta con: go test -tags=integration ./...
}
```

### Table-Driven Tests
```go
func TestValidatePrompt(t *testing.T) {
    tests := []struct {
        name      string
        prompt    string
        wantError bool
    }{
        {"empty string", "", true},
        {"whitespace only", "   ", true},
        {"too short", "ab", true},
        {"valid minimum", "abc", false},
        {"valid prompt", "Learn Go programming", false},
        {"too long", strings.Repeat("a", 501), true},
        {"exactly 500 chars", strings.Repeat("a", 500), false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := usecase.ValidatePrompt(tt.prompt)
            if (err != nil) != tt.wantError {
                t.Errorf("ValidatePrompt(%q) error = %v, wantError %v", 
                    tt.prompt, err, tt.wantError)
            }
        })
    }
}
```

---

## Frontend (TypeScript)

### Vitest Test
```typescript
// src/shared/store/pathStore.test.ts
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { usePathStore } from './pathStore';
import { apiClient } from '../api/client';

// Mock de axios
vi.mock('../api/client', () => ({
  apiClient: {
    post: vi.fn(),
    get: vi.fn(),
  },
}));

describe('pathStore', () => {
  beforeEach(() => {
    // Reset store entre tests
    usePathStore.setState({
      currentPath: null,
      isLoading: false,
      error: null,
    });
    vi.clearAllMocks();
  });

  it('should generate path successfully', async () => {
    // Arrange
    const mockPath = {
      id: 'path-123',
      title: 'Learn Go',
      nodes: [],
    };
    vi.mocked(apiClient.post).mockResolvedValue({ data: mockPath });

    // Act
    const result = await usePathStore.getState().generatePath('Learn Go');

    // Assert
    expect(result).toEqual(mockPath);
    expect(usePathStore.getState().currentPath).toEqual(mockPath);
    expect(usePathStore.getState().isLoading).toBe(false);
    expect(usePathStore.getState().error).toBeNull();
  });

  it('should handle error when generating path', async () => {
    // Arrange
    vi.mocked(apiClient.post).mockRejectedValue(new Error('Network error'));

    // Act & Assert
    await expect(
      usePathStore.getState().generatePath('Learn Go')
    ).rejects.toThrow();
    
    expect(usePathStore.getState().error).toBe('Network error');
  });
});
```

### Property Test con fast-check
```typescript
// src/shared/utils/validation.test.ts
import fc from 'fast-check';
import { describe, it, expect } from 'vitest';
import { validatePrompt } from './validation';

describe('validatePrompt - Property Tests', () => {
  // Property 1: Prompt Validation
  it('should reject empty or whitespace-only prompts', () => {
    fc.assert(
      fc.property(fc.string(), (input) => {
        const result = validatePrompt(input);
        const isEmptyOrWhitespace = input.trim() === '';
        
        if (isEmptyOrWhitespace) {
          expect(result.valid).toBe(false);
        }
      })
    );
  });

  it('should accept prompts with valid length', () => {
    fc.assert(
      fc.property(
        fc.string({ minLength: 3, maxLength: 500 }).filter(s => s.trim().length >= 3),
        (input) => {
          const result = validatePrompt(input);
          expect(result.valid).toBe(true);
        }
      )
    );
  });

  // Property 5: Boilerplate Code Injection
  it('should preserve boilerplate code exactly', () => {
    fc.assert(
      fc.property(fc.string(), (code) => {
        // Simular inyección en editor
        const editorValue = injectBoilerplate(code);
        expect(editorValue).toBe(code);
      })
    );
  });
});

// Property 6: Node Title Display Consistency
describe('NodeTitle - Property Tests', () => {
  it('should display title exactly as received from API', () => {
    fc.assert(
      fc.property(
        fc.string().filter(s => s.length > 0 && s.length <= 50),
        (title) => {
          const rendered = renderNodeTitle(title);
          expect(rendered).toBe(title);
        }
      )
    );
  });
});
```

### Testing de Componentes React
```typescript
// src/features/path-generator/components/PathGenerator.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { PathGenerator } from './PathGenerator';
import { usePathStore } from '@shared/store/pathStore';

// Mock del store
vi.mock('@shared/store/pathStore', () => ({
  usePathStore: vi.fn(),
}));

describe('PathGenerator', () => {
  it('should show loading state while generating', async () => {
    vi.mocked(usePathStore).mockReturnValue({
      generatePath: vi.fn(() => new Promise(() => {})), // Never resolves
      isLoading: true,
      error: null,
    });

    render(<PathGenerator />);
    
    expect(screen.getByText(/Generando/)).toBeInTheDocument();
  });

  it('should show error message on failure', () => {
    vi.mocked(usePathStore).mockReturnValue({
      generatePath: vi.fn(),
      isLoading: false,
      error: 'Network error',
    });

    render(<PathGenerator />);
    
    expect(screen.getByText(/Network error/)).toBeInTheDocument();
  });

  it('should call generatePath on submit', async () => {
    const mockGenerate = vi.fn().mockResolvedValue({ id: '123' });
    vi.mocked(usePathStore).mockReturnValue({
      generatePath: mockGenerate,
      isLoading: false,
      error: null,
    });

    render(<PathGenerator />);
    
    const input = screen.getByPlaceholderText(/Aprender Go/);
    const button = screen.getByRole('button');

    fireEvent.change(input, { target: { value: 'Learn Go' } });
    fireEvent.click(button);

    await waitFor(() => {
      expect(mockGenerate).toHaveBeenCalledWith('Learn Go');
    });
  });
});
```

---

## Comandos

```bash
# Backend - Ejecutar todos los tests
cd backend && go test ./... -v

# Backend - Con cobertura
cd backend && go test ./... -cover

# Backend - Solo property tests
cd backend && go test ./... -run "Property" -v

# Frontend - Ejecutar tests
cd frontend/core && npm test

# Frontend - Con cobertura
cd frontend/core && npm test -- --coverage
```

---

## Checklist

- [ ] Unit tests para cada caso de uso
- [ ] Mocks para dependencias externas (DB, AI)
- [ ] Property tests para las 6 propiedades del design doc
- [ ] Table-driven tests para casos borde
- [ ] Tests de componentes React
- [ ] Cobertura mínima 80%
