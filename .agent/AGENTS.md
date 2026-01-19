# Gemini Coding Path - Agent Context

> Este archivo proporciona contexto completo del proyecto para cualquier agente de IA.
> Léelo antes de generar código para mantener coherencia arquitectónica.

---

## 🎯 Visión del Proyecto

**Gemini Coding Path** es una plataforma de aprendizaje tipo LeetCode/HackerRank donde:

- **Gemini 3 (Vertex AI)** actúa como "Arquitecto Pedagógico" generando rutas de aprendizaje personalizadas
- El **estudiante resuelve los retos SIN asistencia de IA** (código puro, lógica propia)
- Después de resolver, **Gemini actúa como Mentor Senior** dando feedback de Clean Code/SOLID

### Filosofía Core
```
AI for Structure → Human for Logic → AI for Feedback
```

---

## 🛠️ Stack Técnico

| Capa | Tecnología | Versión |
|------|------------|---------|
| **Backend** | Go | 1.23+ |
| **Framework HTTP** | Gin | latest |
| **ORM** | GORM | latest |
| **Frontend** | React | 19.x |
| **Build Tool** | Vite | 7.x |
| **Lenguaje Frontend** | TypeScript | 5.x (strict) |
| **CSS** | Tailwind CSS | 4.x |
| **Estado Global** | Zustand | 5.x |
| **Editor de Código** | Monaco Editor | 4.x |
| **HTTP Client** | Axios | 1.x |
| **Base de Datos** | PostgreSQL | 16 |
| **IA** | Vertex AI (Gemini 3) | - |
| **Contenedores** | Docker | latest |
| **CI/CD** | GitHub Actions | - |

---

## 📁 Estructura del Proyecto

```
gemini-hackathon/
├── .agent/                    # Configuración para agentes IA
│   ├── AGENTS.md              # Este archivo - contexto global
│   ├── workflows/             # Flujos de trabajo paso a paso
│   └── skills/                # Skills específicos por dominio
│
├── backend/                   # API REST en Go
│   ├── cmd/api/               # Punto de entrada (main.go, routes.go)
│   ├── internal/
│   │   ├── domain/            # Entidades y reglas de negocio
│   │   │   ├── entity/        # Structs de dominio
│   │   │   └── repository/    # Interfaces de repositorio
│   │   ├── usecase/           # Casos de uso (lógica de aplicación)
│   │   ├── adapter/           # Adaptadores (HTTP, presenters)
│   │   │   ├── http/          # Controllers Gin
│   │   │   └── repository/    # Implementaciones de repos
│   │   └── infrastructure/    # Infraestructura externa
│   │       ├── database/      # Conexión PostgreSQL
│   │       └── ai/            # Cliente Vertex AI
│   └── .env.example           # Template de variables de entorno
│
├── frontend/core/             # Aplicación React + Vite
│   ├── src/
│   │   ├── features/          # Módulos por funcionalidad
│   │   │   ├── path-generator/
│   │   │   └── node-workspace/
│   │   ├── shared/            # Código compartido
│   │   │   ├── api/           # Cliente API (Axios)
│   │   │   ├── hooks/         # Custom hooks
│   │   │   ├── store/         # Zustand stores
│   │   │   └── components/    # Componentes reutilizables
│   │   ├── App.tsx
│   │   └── main.tsx
│   └── .env.example           # Template de variables de entorno
│
├── docs/                      # Documentación
│   └── gemini-hackathon.wiki/ # Wiki del proyecto
│       ├── Vision.md
│       ├── Architecture.md
│       ├── Database.md
│       ├── API-Docs.md
│       ├── Team.md
│       ├── Workflow.md
│       └── Roadmap.md
│
├── .github/workflows/         # GitHub Actions CI/CD
├── .pre-commit-config.yaml    # Hooks de pre-commit
├── .gitignore
├── LICENSE
└── README.md
```

---

## 🏗️ Arquitectura Backend (Clean Architecture)

```
┌─────────────────────────────────────────┐
│              Frameworks & Drivers        │
│  (Gin HTTP, PostgreSQL, Vertex AI SDK)  │
└─────────────────────────────────────────┘
                    ▲
                    │
┌─────────────────────────────────────────┐
│           Interface Adapters             │
│    (Controllers, Repositories Impl)      │
└─────────────────────────────────────────┘
                    ▲
                    │
┌─────────────────────────────────────────┐
│           Application Layer              │
│    (Use Cases / Interactors)             │
└─────────────────────────────────────────┘
                    ▲
                    │
┌─────────────────────────────────────────┐
│              Domain Layer                │
│    (Entities, Repository Interfaces)     │
└─────────────────────────────────────────┘
```

### Regla de Dependencias
- Las capas internas NO conocen a las externas
- Solo se depende de abstracciones (interfaces)
- El dominio es 100% independiente de frameworks

---

## 📊 Modelo de Datos

```
┌─────────────┐       ┌─────────────────┐       ┌─────────────┐
│   USERS     │──1:N──│ LEARNING_PATHS  │──1:N──│ PATH_NODES  │
└─────────────┘       └─────────────────┘       └─────────────┘
      │                                               │
      │                      ┌────────────────────────┘
      │                      │
      └──────────1:N─────────┼──────────┐
                             │          │
                      ┌──────┴──────────┴───┐
                      │     SUBMISSIONS     │
                      └─────────────────────┘
```

### Tablas Principales
| Tabla | Propósito |
|-------|-----------|
| `users` | Usuarios (guest o registrados) + XP |
| `learning_paths` | Rutas de aprendizaje generadas por IA |
| `path_nodes` | Ejercicios individuales con tests |
| `submissions` | Intentos de solución + feedback IA |

---

## 🔗 Endpoints API

**Base URL:** `/api/v1`

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/paths/generate` | Generar nuevo path con IA |
| GET | `/paths/{id}` | Obtener path con sus nodos |
| GET | `/nodes/{node_id}` | Obtener detalle de ejercicio |
| POST | `/submissions` | Enviar código para evaluación |

---

## 📋 Convenciones de Código

### Git
- **Commits:** Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`)
- **Branches:** `feat/hu-XX-descripcion`, `fix/descripcion`, `docs/descripcion`
- **PR:** Requiere code review del otro miembro

### Go (Backend)
- **Privado:** `camelCase` → `getUserByID`
- **Público:** `PascalCase` → `GetUserByID`
- **Interfaces:** Sufijo `-er` → `PathRepository`, `AIClient`
- **Errors:** Siempre manejar con `if err != nil`

### TypeScript/React (Frontend)
- **Componentes:** `PascalCase` → `PathGenerator.tsx`
- **Hooks:** `camelCase` con prefijo `use` → `useNodeDetail.ts`
- **Stores:** `camelCase` con sufijo `Store` → `pathStore.ts`
- **Utilities:** `camelCase` → `formatDate.ts`

---

## 🔐 Variables de Entorno

### Backend (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=hackathon_db
DB_USER=postgres
DB_PASSWORD=xxx
GOOGLE_CLOUD_PROJECT=xxx
GOOGLE_APPLICATION_CREDENTIALS=./service-account.json
JWT_SECRET=xxx
API_SECRET_KEY=xxx
```

### Frontend (.env)
```env
VITE_API_BASE_URL=http://localhost:8080/api
```

---

## 👥 Equipo

| Rol | Miembro | Responsabilidades |
|-----|---------|-------------------|
| **PO / Backend** | Daikyri | Arquitectura, Go, Vertex AI, DevOps |
| **SM / Frontend** | Hoxanfox | React, UI/UX, Gestión de sprints |

---

## 📚 Skills Disponibles

Para tareas específicas, consulta los skills en `.agent/skills/`:

| Skill | Uso |
|-------|-----|
| `backend-go` | Generar código Go siguiendo Clean Architecture |
| `frontend-react` | Crear componentes React con patrones del proyecto |
| `vertex-ai` | Diseñar prompts y parsear respuestas de Gemini |
| `gin-go` | Crear handlers HTTP con Gin |
| `vite` | Configuración y optimización de Vite |
| `postgres` | Queries, migraciones, optimización |
| `testing` | Property tests, unit tests, mocks |
| `devops` | Docker, CI/CD, GitHub Actions |
| `solid` | Aplicar principios SOLID |
| `clean-architecture` | Estructurar código por capas |
| `git-github` | Flujos de trabajo Git, PRs, branches |
| `scrum` | Gestión ágil, sprints, HUs |

---

## 🎯 Épicas del Proyecto

| # | Épica | Sprint | Estado |
|---|-------|--------|--------|
| E01 | Gestión de Identidad | 2 | 🔜 Pendiente |
| E02 | AI Path Generator | 1 | 🚧 En progreso |
| E03 | Mapa de Aprendizaje | 3 | 🔜 Pendiente |
| E04 | Workspace de Código | 1 | 🚧 En progreso |
| E05 | Code Runner (Juez) | 2 | 🔜 Pendiente |
| E06 | Gamificación | 3 | 🔜 Pendiente |

---

*Última actualización: 2026-01-18*
