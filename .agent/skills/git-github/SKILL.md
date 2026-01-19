---
name: Git & GitHub Workflow
description: Flujos de trabajo Git, branches, PRs y colaboración para Gemini Coding Path
---

# Skill: Git & GitHub

## Contexto

Este proyecto usa **Git** para control de versiones y **GitHub** para colaboración.
Seguimos un flujo de trabajo basado en Pull Requests con code review obligatorio.

---

## Estrategia de Branches

```
main (protegida)
  │
  ├── feat/hu-02-1-generate-path
  │
  ├── feat/hu-04-1-workspace-layout
  │
  ├── fix/validation-error
  │
  └── docs/update-api-docs
```

### Nomenclatura de Branches

| Prefijo | Uso | Ejemplo |
|---------|-----|---------|
| `feat/` | Nueva funcionalidad | `feat/hu-02-1-generate-path` |
| `fix/` | Corrección de bugs | `fix/login-error` |
| `docs/` | Documentación | `docs/update-readme` |
| `refactor/` | Mejoras de código sin cambio funcional | `refactor/path-controller` |
| `test/` | Agregar o modificar tests | `test/path-usecase` |
| `chore/` | Tareas de mantenimiento | `chore/update-dependencies` |

---

## Conventional Commits

```
<tipo>(<alcance>): <descripción>

[cuerpo opcional]

[footer opcional]
```

### Tipos de Commit

| Tipo | Descripción | Ejemplo |
|------|-------------|---------|
| `feat` | Nueva funcionalidad | `feat(api): agregar endpoint de generación de path` |
| `fix` | Corrección de bug | `fix(ui): corregir validación de prompt vacío` |
| `docs` | Solo documentación | `docs: actualizar README con instrucciones` |
| `style` | Formato (no afecta lógica) | `style: aplicar gofmt` |
| `refactor` | Cambio de código sin cambiar funcionalidad | `refactor(usecase): extraer validación a función` |
| `test` | Agregar o corregir tests | `test(path): agregar test de round-trip` |
| `chore` | Tareas de mantenimiento | `chore: actualizar dependencias` |

### Ejemplos Buenos vs Malos

```bash
# ❌ MAL
git commit -m "fix"
git commit -m "cambios"
git commit -m "WIP"

# ✅ BIEN
git commit -m "feat(api): implementar POST /paths/generate"
git commit -m "fix(ui): corregir estado de loading en PathGenerator"
git commit -m "docs(wiki): agregar diagrama de arquitectura"
```

---

## Flujo de Trabajo

### 1. Crear Branch desde main
```bash
# Asegurarse de estar en main actualizado
git checkout main
git pull origin main

# Crear nueva branch
git checkout -b feat/hu-02-1-generate-path
```

### 2. Hacer Commits Pequeños y Frecuentes
```bash
# Agregar cambios
git add backend/internal/usecase/generate_path.go

# Commit con mensaje descriptivo
git commit -m "feat(usecase): crear estructura base de GeneratePathUseCase"

# Continuar trabajando...
git add .
git commit -m "feat(usecase): implementar validación de prompt"
```

### 3. Push a Remote
```bash
# Primera vez
git push -u origin feat/hu-02-1-generate-path

# Siguientes pushes
git push
```

### 4. Crear Pull Request
```markdown
## Descripción
Implementa el caso de uso de generación de path con IA.

## Cambios
- Crear `GeneratePathUseCase` con validación de prompt
- Agregar interfaz `AIClient` 
- Implementar tests unitarios

## Relacionado
Closes #7 (HU-02.1)

## Checklist
- [x] Tests pasan
- [x] Código sigue convenciones
- [x] Documentación actualizada (si aplica)
```

### 5. Code Review
El otro miembro del equipo revisa:
- ¿El código sigue las convenciones?
- ¿Hay tests suficientes?
- ¿El PR está relacionado con el issue correcto?

### 6. Merge a Main
Después de aprobación, el autor hace merge (Squash and merge recomendado).

---

## Comandos Git Útiles

### Ver estado
```bash
git status                  # Ver archivos modificados
git log --oneline -10       # Ver últimos 10 commits
git branch -a               # Ver todas las branches
```

### Sincronizar
```bash
git fetch origin            # Obtener cambios remotos
git pull origin main        # Actualizar main local
git merge main              # Traer cambios de main a tu branch
```

### Resolver conflictos
```bash
# Si hay conflictos al hacer merge/pull
git status                  # Ver archivos en conflicto
# Editar archivos manualmente
git add <archivos>
git commit -m "fix: resolver conflictos de merge"
```

### Deshacer cambios
```bash
git checkout -- <archivo>   # Descartar cambios no staged
git reset HEAD <archivo>    # Unstage un archivo
git reset --soft HEAD~1     # Deshacer último commit (mantener cambios)
git reset --hard HEAD~1     # Deshacer último commit (perder cambios)
```

### Stash
```bash
git stash                   # Guardar cambios temporalmente
git stash list              # Ver stashes
git stash pop               # Recuperar último stash
git stash drop              # Eliminar último stash
```

---

## GitHub Features

### Issues
```markdown
## Issue Template

### Descripción
[Descripción clara del problema o feature]

### Criterios de Aceptación
- [ ] Criterio 1
- [ ] Criterio 2

### Contexto Adicional
[Screenshots, logs, etc.]
```

### Labels Recomendados
| Label | Color | Uso |
|-------|-------|-----|
| `bug` | #d73a4a | Bug confirmado |
| `feature` | #a2eeef | Nueva funcionalidad |
| `documentation` | #0075ca | Documentación |
| `good first issue` | #7057ff | Para principiantes |
| `priority: high` | #b60205 | Urgente |
| `priority: low` | #c5def5 | Puede esperar |

### Milestones
```
Sprint 1 - MVP Core
├── HU-02.1: Generación de Path
├── HU-02.2: Persistencia de Path
├── HU-04.1: Layout del Editor
└── HU-04.2: Carga de Ejercicio
```

---

## Protección de Branch Main

### Reglas Recomendadas
1. ✅ Require pull request before merging
2. ✅ Require approvals (1)
3. ✅ Require status checks to pass (CI)
4. ✅ Require conversation resolution
5. ❌ Do not allow bypassing

---

## .gitignore

```gitignore
# Binarios
*.exe
*.dll
*.so
*.dylib
/backend/main

# Dependencias
node_modules/
vendor/

# Archivos de configuración local
.env
*.local

# Credenciales
*.pem
*.key
service-account*.json

# IDE
.vscode/
.idea/
*.swp
*.swo

# Build
dist/
build/
coverage/

# OS
.DS_Store
Thumbs.db

# Logs
*.log
npm-debug.log*
```

---

## Checklist

- [ ] Branch creada desde main actualizado
- [ ] Nombre de branch sigue nomenclatura
- [ ] Commits siguen Conventional Commits
- [ ] PR tiene descripción y referencia al issue
- [ ] CI pasa antes de solicitar review
- [ ] Code review aprobado
- [ ] Merge con Squash
