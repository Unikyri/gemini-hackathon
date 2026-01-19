---
name: Vertex AI - Gemini 3 Integration
description: Diseñar prompts y parsear respuestas de Gemini 3 para Gemini Coding Path
---

# Skill: Vertex AI (Gemini 3)

## Contexto

Este proyecto usa **Vertex AI** con el modelo **Gemini 3** para:
1. Generar rutas de aprendizaje personalizadas
2. Generar ejercicios con enunciados, código y tests
3. Proporcionar feedback estilo "Senior Developer"

**Importante:** La IA NO resuelve ejercicios, solo estructura, genera contenido y evalúa.

---

## Configuración del Cliente Go

### Instalación de dependencias
```bash
go get cloud.google.com/go/aiplatform/apiv1
go get google.golang.org/api/option
```

### Cliente Vertex AI
```go
// internal/infrastructure/ai/vertex_client.go
package ai

import (
    "context"
    "encoding/json"
    "fmt"
    
    aiplatform "cloud.google.com/go/aiplatform/apiv1"
    "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
    "google.golang.org/api/option"
    "google.golang.org/protobuf/types/known/structpb"
)

type VertexClient struct {
    client    *aiplatform.PredictionClient
    projectID string
    location  string
    model     string
}

func NewVertexClient(projectID, location string) (*VertexClient, error) {
    ctx := context.Background()
    
    endpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
    client, err := aiplatform.NewPredictionClient(ctx, 
        option.WithEndpoint(endpoint),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }
    
    return &VertexClient{
        client:    client,
        projectID: projectID,
        location:  location,
        model:     "gemini-1.5-pro", // o "gemini-1.5-flash" para más velocidad
    }, nil
}

func (c *VertexClient) GeneratePath(ctx context.Context, prompt string) (*AIPathResponse, error) {
    // Construir el endpoint del modelo
    endpoint := fmt.Sprintf(
        "projects/%s/locations/%s/publishers/google/models/%s",
        c.projectID, c.location, c.model,
    )
    
    // Construir el prompt completo
    fullPrompt := buildPathGenerationPrompt(prompt)
    
    // Crear la instancia del request
    instance, _ := structpb.NewValue(map[string]interface{}{
        "prompt": fullPrompt,
    })
    
    // Parámetros de generación
    params, _ := structpb.NewValue(map[string]interface{}{
        "temperature":     0.7,
        "maxOutputTokens": 4096,
        "topP":            0.95,
        "topK":            40,
    })
    
    req := &aiplatformpb.PredictRequest{
        Endpoint:   endpoint,
        Instances:  []*structpb.Value{instance},
        Parameters: params,
    }
    
    resp, err := c.client.Predict(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("prediction failed: %w", err)
    }
    
    // Parsear respuesta
    return parsePathResponse(resp)
}
```

---

## System Prompts

### Para Generar Path (Path Generator)
```go
const systemPromptPathGenerator = `
Eres un arquitecto pedagógico experto en programación. Tu tarea es generar rutas de aprendizaje estructuradas y personalizadas.

REGLAS ESTRICTAS:
1. Responde ÚNICAMENTE con JSON válido, sin texto adicional ni markdown
2. Genera exactamente 5 nodos de aprendizaje
3. El primer nodo siempre debe estar "unlocked", los demás "locked"
4. Cada nodo debe ser progresivamente más complejo
5. Los títulos deben ser concisos (máximo 50 caracteres)
6. Las descripciones deben ser claras (máximo 150 caracteres)
7. El código boilerplate debe ser funcional y tener TODO markers
8. Los tests deben validar la funcionalidad esperada

FORMATO DE RESPUESTA:
{
  "title": "Título del path (máximo 100 caracteres)",
  "nodes": [
    {
      "title": "Título del nodo",
      "description": "Descripción corta del ejercicio",
      "description_md": "# Enunciado Completo\n\n## Objetivo\nDescripción detallada...\n\n## Requisitos\n- Req 1\n- Req 2\n\n## Ejemplo\n...",
      "boilerplate_code": "package main\n\n// TODO: Implementa la función...\nfunc main() {\n\t// Tu código aquí\n}",
      "tests_payload": {
        "required_tests": [
          {"name": "TestBasicFunction", "code": "..."}
        ],
        "bonus_tests": [
          {"name": "TestEdgeCase", "code": "...", "xp_bonus": 50}
        ]
      },
      "documentation_md": "## Documentación Relevante\n\n### El paquete fmt\n..."
    }
  ]
}
`

func buildPathGenerationPrompt(userPrompt string) string {
    return fmt.Sprintf(`%s

TEMA SOLICITADO POR EL USUARIO:
"%s"

Genera la ruta de aprendizaje ahora:`, systemPromptPathGenerator, userPrompt)
}
```

### Para Feedback de Senior (Post-Solución)
```go
const systemPromptSeniorFeedback = `
Eres un desarrollador senior con 15 años de experiencia revisando código de juniors.
Tu rol es dar feedback constructivo enfocado en:

1. CLEAN CODE
   - Nombres descriptivos
   - Funciones pequeñas y enfocadas
   - Evitar comentarios innecesarios
   - Código autoexplicativo

2. PRINCIPIOS SOLID
   - Single Responsibility
   - Open/Closed
   - Liskov Substitution
   - Interface Segregation
   - Dependency Inversion

3. BUENAS PRÁCTICAS DE GO
   - Manejo de errores
   - Nombrado idiomático
   - Uso de interfaces
   - Concurrencia (si aplica)

FORMATO DE RESPUESTA:
{
  "overall_score": 85,
  "summary": "Resumen de 1-2 líneas",
  "strengths": ["Punto fuerte 1", "Punto fuerte 2"],
  "improvements": [
    {
      "category": "CLEAN_CODE|SOLID|GO_IDIOMS",
      "issue": "Descripción del problema",
      "suggestion": "Cómo mejorarlo",
      "code_example": "// Ejemplo de código mejorado"
    }
  ],
  "learning_tips": ["Tip 1", "Tip 2"]
}
`
```

---

## Estructuras de Respuesta

```go
// internal/infrastructure/ai/types.go
package ai

// Respuesta de generación de path
type AIPathResponse struct {
    Title string   `json:"title"`
    Nodes []AINode `json:"nodes"`
}

type AINode struct {
    Title           string      `json:"title"`
    Description     string      `json:"description"`
    DescriptionMD   string      `json:"description_md"`
    BoilerplateCode string      `json:"boilerplate_code"`
    TestsPayload    TestPayload `json:"tests_payload"`
    DocumentationMD string      `json:"documentation_md"`
}

type TestPayload struct {
    RequiredTests []Test `json:"required_tests"`
    BonusTests    []Test `json:"bonus_tests"`
}

type Test struct {
    Name    string `json:"name"`
    Code    string `json:"code"`
    XPBonus int    `json:"xp_bonus,omitempty"`
}

// Respuesta de feedback de senior
type SeniorFeedbackResponse struct {
    OverallScore int               `json:"overall_score"`
    Summary      string            `json:"summary"`
    Strengths    []string          `json:"strengths"`
    Improvements []Improvement     `json:"improvements"`
    LearningTips []string          `json:"learning_tips"`
}

type Improvement struct {
    Category    string `json:"category"`
    Issue       string `json:"issue"`
    Suggestion  string `json:"suggestion"`
    CodeExample string `json:"code_example,omitempty"`
}
```

---

## Parsing de Respuestas

```go
func parsePathResponse(resp *aiplatformpb.PredictResponse) (*AIPathResponse, error) {
    if len(resp.Predictions) == 0 {
        return nil, errors.New("no predictions in response")
    }
    
    // Extraer el texto de la respuesta
    prediction := resp.Predictions[0]
    text := extractTextFromPrediction(prediction)
    
    // Limpiar posibles artefactos de markdown
    text = cleanJSONResponse(text)
    
    // Parsear JSON
    var result AIPathResponse
    if err := json.Unmarshal([]byte(text), &result); err != nil {
        return nil, fmt.Errorf("failed to parse AI response: %w", err)
    }
    
    // Validar estructura
    if err := validatePathResponse(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

func cleanJSONResponse(text string) string {
    // Remover markdown code blocks si existen
    text = strings.TrimPrefix(text, "```json")
    text = strings.TrimPrefix(text, "```")
    text = strings.TrimSuffix(text, "```")
    return strings.TrimSpace(text)
}

func validatePathResponse(resp *AIPathResponse) error {
    if resp.Title == "" {
        return errors.New("path title is empty")
    }
    if len(resp.Nodes) != 5 {
        return fmt.Errorf("expected 5 nodes, got %d", len(resp.Nodes))
    }
    
    for i, node := range resp.Nodes {
        if node.Title == "" {
            return fmt.Errorf("node %d has empty title", i)
        }
        if node.BoilerplateCode == "" {
            return fmt.Errorf("node %d has no boilerplate code", i)
        }
    }
    
    return nil
}
```

---

## Manejo de Errores y Retry

```go
func (c *VertexClient) GeneratePathWithRetry(ctx context.Context, prompt string) (*AIPathResponse, error) {
    var lastErr error
    
    for attempt := 0; attempt < 3; attempt++ {
        result, err := c.GeneratePath(ctx, prompt)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // Solo reintentar en errores transitorios
        if !isRetriableError(err) {
            return nil, err
        }
        
        // Backoff exponencial
        waitTime := time.Duration(math.Pow(2, float64(attempt))) * time.Second
        time.Sleep(waitTime)
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func isRetriableError(err error) bool {
    // Errores de rate limiting o timeout son retriables
    errStr := err.Error()
    return strings.Contains(errStr, "429") || 
           strings.Contains(errStr, "timeout") ||
           strings.Contains(errStr, "unavailable")
}
```

---

## Checklist para Prompts

- [ ] JSON Schema definido claramente
- [ ] Reglas estrictas numeradas
- [ ] Ejemplos de formato esperado
- [ ] Límites de caracteres especificados
- [ ] Validación de respuesta implementada
- [ ] Retry con backoff para errores transitorios
- [ ] Logging de prompts/respuestas para debugging
