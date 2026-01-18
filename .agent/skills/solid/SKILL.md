---
name: SOLID Principles
description: Aplicar principios SOLID en el código de Gemini Coding Path
---

# Skill: Principios SOLID

## Contexto

Los principios SOLID son fundamentales para escribir código mantenible y extensible.
Este skill muestra cómo aplicarlos en el contexto de Gemini Coding Path.

---

## S - Single Responsibility Principle

> Una clase/módulo debe tener una sola razón para cambiar.

### ❌ Violación
```go
// Un controller que hace demasiado
type PathController struct {
    db *sql.DB
}

func (c *PathController) GeneratePath(ctx *gin.Context) {
    // Parsear request
    var req GeneratePathRequest
    ctx.ShouldBindJSON(&req)
    
    // Validar prompt (responsabilidad de validación)
    if len(req.Prompt) < 3 {
        ctx.JSON(400, gin.H{"error": "prompt too short"})
        return
    }
    
    // Llamar a IA (responsabilidad de AI)
    aiResponse, _ := callVertexAI(req.Prompt)
    
    // Guardar en BD (responsabilidad de persistencia)
    c.db.Exec("INSERT INTO paths...")
    
    // Formatear respuesta (responsabilidad de presentación)
    ctx.JSON(200, formatResponse(aiResponse))
}
```

### ✅ Correcto
```go
// Controller solo maneja HTTP
type PathController struct {
    useCase *usecase.GeneratePathUseCase
}

func (c *PathController) GeneratePath(ctx *gin.Context) {
    var req dto.GeneratePathRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    path, err := c.useCase.Execute(ctx, req.Prompt)
    if err != nil {
        ctx.JSON(mapErrorToStatus(err), dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    ctx.JSON(200, dto.ToResponse(path))
}

// Use Case maneja lógica de negocio
type GeneratePathUseCase struct {
    aiClient   AIClient      // Responsabilidad: IA
    pathRepo   PathRepository // Responsabilidad: Persistencia
    validator  Validator      // Responsabilidad: Validación
}
```

---

## O - Open/Closed Principle

> Abierto para extensión, cerrado para modificación.

### ❌ Violación
```go
func (uc *GeneratePathUseCase) Execute(prompt string, provider string) (*Path, error) {
    var response *AIResponse
    
    // Cada nuevo provider requiere modificar este switch
    switch provider {
    case "vertex":
        response, _ = callVertexAI(prompt)
    case "openai":
        response, _ = callOpenAI(prompt)
    case "claude":
        response, _ = callClaude(prompt)
    default:
        return nil, errors.New("unknown provider")
    }
    
    return mapToPath(response), nil
}
```

### ✅ Correcto
```go
// Interfaz que permite extensión
type AIClient interface {
    GeneratePath(ctx context.Context, prompt string) (*AIResponse, error)
}

// Implementaciones específicas (extensión sin modificar el use case)
type VertexAIClient struct { /* ... */ }
func (c *VertexAIClient) GeneratePath(ctx context.Context, prompt string) (*AIResponse, error) { /* ... */ }

type OpenAIClient struct { /* ... */ }
func (c *OpenAIClient) GeneratePath(ctx context.Context, prompt string) (*AIResponse, error) { /* ... */ }

// Use case cerrado para modificación
type GeneratePathUseCase struct {
    aiClient AIClient // Cualquier implementación funciona
}

func (uc *GeneratePathUseCase) Execute(ctx context.Context, prompt string) (*Path, error) {
    response, err := uc.aiClient.GeneratePath(ctx, prompt)
    if err != nil {
        return nil, err
    }
    return mapToPath(response), nil
}
```

---

## L - Liskov Substitution Principle

> Las subclases deben poder sustituir a sus clases base sin romper el programa.

### ❌ Violación
```typescript
interface CodeExecutor {
  execute(code: string): ExecutionResult;
}

class GoExecutor implements CodeExecutor {
  execute(code: string): ExecutionResult {
    return runGoCode(code);
  }
}

class PythonExecutor implements CodeExecutor {
  execute(code: string): ExecutionResult {
    // Rompe el contrato: lanza excepción en lugar de retornar resultado
    throw new Error("Python not supported yet");
  }
}
```

### ✅ Correcto
```typescript
interface CodeExecutor {
  execute(code: string): ExecutionResult;
  isSupported(): boolean;
}

class GoExecutor implements CodeExecutor {
  execute(code: string): ExecutionResult {
    return runGoCode(code);
  }
  
  isSupported(): boolean {
    return true;
  }
}

class PythonExecutor implements CodeExecutor {
  execute(code: string): ExecutionResult {
    return { success: false, error: "Python not supported yet" };
  }
  
  isSupported(): boolean {
    return false; // Maneja el caso honestamente
  }
}
```

---

## I - Interface Segregation Principle

> Muchas interfaces específicas son mejores que una interfaz general.

### ❌ Violación
```go
// Interfaz demasiado grande
type Repository interface {
    CreatePath(path *Path) error
    GetPath(id string) (*Path, error)
    UpdatePath(path *Path) error
    DeletePath(id string) error
    CreateNode(node *Node) error
    GetNode(id string) (*Node, error)
    UpdateNode(node *Node) error
    DeleteNode(id string) error
    CreateSubmission(sub *Submission) error
    GetSubmission(id string) (*Submission, error)
    // ... más métodos
}

// Los use cases se ven forzados a depender de métodos que no usan
type GeneratePathUseCase struct {
    repo Repository // Solo necesita CreatePath, pero tiene acceso a todo
}
```

### ✅ Correcto
```go
// Interfaces segregadas por responsabilidad
type PathRepository interface {
    Create(ctx context.Context, path *Path) error
    GetByID(ctx context.Context, id string) (*Path, error)
    GetByUserID(ctx context.Context, userID string) ([]*Path, error)
}

type NodeRepository interface {
    CreateBatch(ctx context.Context, nodes []*Node) error
    GetByID(ctx context.Context, id string) (*Node, error)
    GetByPathID(ctx context.Context, pathID string) ([]*Node, error)
}

type SubmissionRepository interface {
    Create(ctx context.Context, sub *Submission) error
    GetByNodeID(ctx context.Context, nodeID string) ([]*Submission, error)
}

// Cada use case depende solo de lo que necesita
type GeneratePathUseCase struct {
    pathRepo PathRepository
    nodeRepo NodeRepository
}
```

---

## D - Dependency Inversion Principle

> Depender de abstracciones, no de implementaciones concretas.

### ❌ Violación
```go
type GeneratePathUseCase struct {
    // Dependencias concretas
    db     *sql.DB           // Implementación específica de BD
    vertex *VertexAIClient   // Implementación específica de IA
}

func (uc *GeneratePathUseCase) Execute(prompt string) (*Path, error) {
    // Acoplado a PostgreSQL
    rows, _ := uc.db.Query("SELECT * FROM paths")
    
    // Acoplado a Vertex AI
    response, _ := uc.vertex.CallAPI(prompt)
    
    return mapToPath(response), nil
}
```

### ✅ Correcto
```go
// Interfaces definidas en el dominio (capa interna)
type PathRepository interface {
    Create(ctx context.Context, path *Path) error
    GetByID(ctx context.Context, id string) (*Path, error)
}

type AIClient interface {
    GeneratePath(ctx context.Context, prompt string) (*AIResponse, error)
}

// Use case depende de abstracciones
type GeneratePathUseCase struct {
    pathRepo PathRepository  // Abstracción
    aiClient AIClient        // Abstracción
}

func NewGeneratePathUseCase(repo PathRepository, ai AIClient) *GeneratePathUseCase {
    return &GeneratePathUseCase{
        pathRepo: repo,
        aiClient: ai,
    }
}

// La inversión ocurre en la capa de infraestructura (inyección)
func main() {
    // Implementaciones concretas se inyectan desde afuera
    db := database.NewPostgresConnection()
    pathRepo := repository.NewPostgresPathRepository(db)
    aiClient := ai.NewVertexClient()
    
    useCase := usecase.NewGeneratePathUseCase(pathRepo, aiClient)
}
```

---

## Aplicación en Frontend (React/TypeScript)

### Single Responsibility
```typescript
// ❌ Componente que hace todo
function PathGenerator() {
  const [prompt, setPrompt] = useState('');
  const [path, setPath] = useState(null);
  const [loading, setLoading] = useState(false);
  
  const handleSubmit = async () => {
    setLoading(true);
    const response = await fetch('/api/paths/generate', {
      method: 'POST',
      body: JSON.stringify({ prompt }),
    });
    const data = await response.json();
    setPath(data);
    setLoading(false);
  };
  
  return (/* UI + lógica mezcladas */);
}

// ✅ Separación de responsabilidades
// Hook para lógica de datos
function usePathGenerator() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const generatePath = async (prompt: string) => {
    setIsLoading(true);
    try {
      const path = await pathApi.generate(prompt);
      return path;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };
  
  return { generatePath, isLoading, error };
}

// Componente solo para UI
function PathGenerator() {
  const [prompt, setPrompt] = useState('');
  const { generatePath, isLoading, error } = usePathGenerator();
  
  return (/* Solo UI */);
}
```

---

## Checklist SOLID

- [ ] ¿Cada clase/función tiene una sola responsabilidad?
- [ ] ¿Puedo agregar nuevas funcionalidades sin modificar código existente?
- [ ] ¿Las implementaciones pueden sustituirse sin romper el sistema?
- [ ] ¿Las interfaces son pequeñas y específicas?
- [ ] ¿Las dependencias son abstracciones inyectadas?
