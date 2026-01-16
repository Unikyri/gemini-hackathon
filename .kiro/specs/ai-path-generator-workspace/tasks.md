# Implementation Plan: AI Path Generator & Workspace

## Overview

Este plan implementa la generación de rutas de aprendizaje con IA y el workspace de código. Se divide en tareas incrementales que construyen sobre las anteriores, comenzando con la infraestructura del backend, luego la integración con Vertex AI, persistencia, y finalmente el frontend.

## Tasks

- [ ] 1. Configurar estructura base del backend
  - [ ] 1.1 Crear estructura de directorios siguiendo Clean Architecture
    - Crear carpetas: `cmd/api`, `internal/domain`, `internal/usecase`, `internal/adapter`, `internal/infrastructure`
    - Configurar `go.mod` con dependencias iniciales (gin, uuid, viper)
    - _Requirements: 1.1, 1.2_
  - [ ] 1.2 Definir entidades de dominio
    - Crear `LearningPath` y `PathNode` en `internal/domain/entity`
    - Definir tipos `PathStatus` y `NodeStatus`
    - _Requirements: 1.4, 2.2, 2.3_
  - [ ] 1.3 Crear interfaces de repositorio
    - Definir `PathRepository` y `NodeRepository` en `internal/domain/repository`
    - _Requirements: 2.2, 2.3, 2.4_

- [ ] 2. Implementar capa de infraestructura de base de datos
  - [ ] 2.1 Configurar conexión a PostgreSQL
    - Crear cliente de base de datos con pool de conexiones
    - Implementar health check de conexión
    - _Requirements: 2.1_
  - [ ] 2.2 Implementar PathRepository con PostgreSQL
    - Implementar métodos `Create`, `GetByID`, `GetByUserID`
    - _Requirements: 2.2, 2.4_
  - [ ] 2.3 Implementar NodeRepository con PostgreSQL
    - Implementar métodos `CreateBatch`, `GetByPathID`, `GetByID`
    - _Requirements: 2.3, 2.4_
  - [ ] 2.4 Write property test for Path persistence round-trip
    - **Property 3: Path Persistence Round-Trip**
    - **Validates: Requirements 2.2, 2.3, 2.4**

- [ ] 3. Implementar cliente de Vertex AI
  - [ ] 3.1 Crear cliente de Vertex AI
    - Implementar interfaz `AIClient` con método `GeneratePath`
    - Configurar autenticación con Google Cloud
    - Definir System Prompt para generación de JSON estructurado
    - _Requirements: 1.2_
  - [ ] 3.2 Implementar parsing de respuesta de IA
    - Parsear JSON de respuesta a estructuras Go
    - Validar que la respuesta contenga exactamente 5 nodos
    - Validar campos requeridos en cada nodo
    - _Requirements: 1.3, 1.4_
  - [ ] 3.3 Write property test for AI response structure
    - **Property 2: AI Response Structure Integrity**
    - **Validates: Requirements 1.3, 1.4, 1.5**

- [ ] 4. Implementar caso de uso GeneratePath
  - [ ] 4.1 Crear GeneratePath use case
    - Implementar validación de prompt (mínimo 3, máximo 500 caracteres)
    - Orquestar llamada a AIClient y persistencia
    - Asignar UUIDs a path y nodos
    - Establecer primer nodo como "unlocked", resto como "locked"
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 2.2, 2.3_
  - [ ] 4.2 Write property test for prompt validation
    - **Property 1: Prompt Validation**
    - **Validates: Requirements 1.1, 1.6**

- [ ] 5. Implementar controladores HTTP
  - [ ] 5.1 Crear PathController
    - Implementar `POST /api/v1/paths/generate`
    - Implementar `GET /api/v1/paths/{id}`
    - Mapear entidades a DTOs de respuesta
    - _Requirements: 1.5, 1.6, 1.7, 2.4, 2.5_
  - [ ] 5.2 Crear NodeController
    - Implementar `GET /api/v1/nodes/{node_id}`
    - Retornar markdown_content, boilerplate_code, documentation
    - _Requirements: 4.1_
  - [ ] 5.3 Configurar router y middleware
    - Configurar Gin con rutas
    - Agregar middleware de CORS, logging y recovery
    - _Requirements: 1.5_

- [ ] 6. Checkpoint - Backend funcional
  - Ensure all tests pass, ask the user if questions arise.
  - Verificar que los endpoints respondan correctamente con Postman/curl

- [ ] 7. Configurar estructura base del frontend
  - [ ] 7.1 Crear proyecto React con Vite
    - Configurar TypeScript en modo estricto
    - Instalar dependencias: tailwindcss, @monaco-editor/react, zustand, axios
    - _Requirements: 3.1_
  - [ ] 7.2 Crear cliente API
    - Implementar funciones `generatePath`, `getPath`, `getNode`
    - Configurar base URL desde variables de entorno
    - _Requirements: 4.1_

- [ ] 8. Implementar componente Workspace
  - [ ] 8.1 Crear layout de 3 columnas
    - Panel izquierdo: Instrucciones (25%)
    - Panel central: Editor (50%)
    - Panel derecho: Info/Salida (25%)
    - Implementar redimensionamiento responsive
    - _Requirements: 3.1, 3.3, 3.4_
  - [ ] 8.2 Integrar Monaco Editor
    - Configurar para lenguaje Go
    - Aplicar tema oscuro
    - Deshabilitar minimap
    - _Requirements: 3.2_
  - [ ] 8.3 Write property test for boilerplate injection
    - **Property 5: Boilerplate Code Injection**
    - **Validates: Requirements 4.3**

- [ ] 9. Implementar carga de ejercicio
  - [ ] 9.1 Crear hook useNodeDetail
    - Fetch de datos del nodo al montar componente
    - Manejar estados de loading y error
    - _Requirements: 4.1, 4.4_
  - [ ] 9.2 Renderizar contenido Markdown
    - Integrar librería de renderizado Markdown (react-markdown)
    - Mostrar enunciado en panel izquierdo
    - _Requirements: 4.2_
  - [ ] 9.3 Write property test for Markdown rendering
    - **Property 4: Markdown Content Rendering**
    - **Validates: Requirements 4.2**
  - [ ] 9.4 Inyectar boilerplate en editor
    - Cargar código inicial en Monaco al recibir respuesta
    - Mostrar título del nodo en cabecera
    - _Requirements: 4.3, 4.5_
  - [ ] 9.5 Write property test for title display
    - **Property 6: Node Title Display Consistency**
    - **Validates: Requirements 4.5**

- [ ] 10. Implementar manejo de errores en frontend
  - [ ] 10.1 Crear componentes de error
    - Toast para errores de red y servidor
    - Pantalla de error para 404
    - Mensajes inline para validación
    - _Requirements: 4.4_
  - [ ] 10.2 Implementar retry para errores 503/504
    - Backoff exponencial con máximo 3 intentos
    - _Requirements: 1.7_

- [ ] 11. Checkpoint final - Integración completa
  - Ensure all tests pass, ask the user if questions arise.
  - Verificar flujo completo: generar path → seleccionar nodo → ver workspace

## Notes

- All tasks including property tests are required for comprehensive coverage
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
- El backend usa Go con Clean Architecture
- El frontend usa React + TypeScript + Tailwind
- Property-based testing: gopter (Go), fast-check (TypeScript)
