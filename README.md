# gemini-hackathon

Plataforma de aprendizaje tipo LeetCode/HackerRank con learning paths generados por IA, pero sin uso de IA para resolver ejercicios. Enfocada en aprender programación desde la lógica y la práctica, con dificultad progresiva y énfasis en buenas prácticas, clean code, patrones de diseño y arquitectura. MVP en Go.

## 📋 Tabla de Contenidos

- [Arquitectura](#arquitectura)
- [Configuración del Entorno](#configuración-del-entorno)
- [Seguridad y Secretos](#seguridad-y-secretos)
- [Pre-commit Hooks](#pre-commit-hooks)
- [GitHub Secrets](#github-secrets)
- [Desarrollo Local](#desarrollo-local)

## 🏗️ Arquitectura

El proyecto está dividido en dos componentes principales:

- **Backend** (`/backend`): API REST en Go
- **Frontend** (`/frontend`): Aplicación React
- **Base de Datos**: PostgreSQL desplegada en Docker

## ⚙️ Configuración del Entorno

### Backend

1. Navega al directorio del backend:
   ```bash
   cd backend
   ```

2. Copia el archivo de ejemplo y configura tus variables de entorno:
   ```bash
   cp .env.example .env
   ```

3. Edita el archivo `.env` con tus credenciales reales:
   - Configura las credenciales de PostgreSQL (Docker)
   - Añade tu `GOOGLE_CLOUD_PROJECT` y ruta a `GOOGLE_APPLICATION_CREDENTIALS`
   - Genera secretos seguros para `JWT_SECRET` y `API_SECRET_KEY`

### Frontend

1. Navega al directorio del frontend:
   ```bash
   cd frontend
   ```

2. Copia el archivo de ejemplo y configura tus variables de entorno:
   ```bash
   cp .env.example .env
   ```

3. Edita el archivo `.env` con la URL de tu API backend:
   - `REACT_APP_API_URL`: URL del backend (por defecto: `http://localhost:8080`)
   - **IMPORTANTE**: Solo variables con prefijo `REACT_APP_` serán accesibles en el cliente
   - **NUNCA** pongas secretos sensibles aquí, ya que serán públicamente visibles

## 🔒 Seguridad y Secretos

### Archivos a NO Commitear

El archivo `.gitignore` está configurado para ignorar automáticamente:
- ✅ Archivos `.env` (excepto `.env.example`)
- ✅ Credenciales de servicios (`service-account*.json`, `*.pem`, `*.key`)
- ✅ `node_modules/`, `dist/`, `build/`, `coverage/`
- ✅ Binarios compilados y archivos temporales
- ✅ Configuraciones de IDEs (`.vscode/`, `.idea/`)

### Variables de Entorno Sensibles

**Backend** requiere:
- Credenciales de base de datos PostgreSQL
- Credenciales de Google Cloud / Vertex AI
- Secretos JWT y API keys
- Configuración de CORS

**Frontend** requiere:
- URL de la API backend
- IDs de proyecto (solo valores públicos)

## 🔐 Pre-commit Hooks

Pre-commit hooks ayudan a prevenir commits con secretos o código roto antes de que lleguen al repositorio.

### Instalación

1. Instala pre-commit:
   ```bash
   pip install pre-commit
   ```

2. Instala los hooks en tu repositorio local:
   ```bash
   pre-commit install
   ```

3. (Opcional) Ejecuta los hooks en todos los archivos:
   ```bash
   pre-commit run --all-files
   ```

### Hooks Configurados

Los siguientes hooks se ejecutarán automáticamente antes de cada commit:

- **Detección de secretos** (`detect-secrets`): Previene commits con API keys, contraseñas, tokens
- **Detección de claves privadas**: Bloquea archivos `.pem`, `.key`
- **Archivos grandes**: Rechaza archivos mayores a 1MB
- **Formato de código**:
  - Go: `gofmt`, `go-vet`, `goimports`, `go-mod-tidy`, `golangci-lint`
  - JavaScript/React: `eslint`
- **Validación de sintaxis**: YAML, JSON
- **Limpieza**: Espacios al final de línea, fin de archivo

### Bypass (Solo en casos de emergencia)

Si necesitas hacer un commit sin ejecutar los hooks (NO RECOMENDADO):
```bash
git commit --no-verify -m "mensaje"
```

## 🔑 GitHub Secrets

Para configurar el CI/CD y despliegues, necesitas configurar los siguientes secretos en GitHub:

### Configuración en GitHub

1. Ve a tu repositorio en GitHub
2. Navega a **Settings** → **Secrets and variables** → **Actions**
3. Click en **New repository secret**
4. Añade los siguientes secretos:

### Secretos Requeridos

| Nombre del Secreto | Descripción | Ejemplo |
|-------------------|-------------|---------|
| `GOOGLE_CLOUD_PROJECT` | ID del proyecto de Google Cloud | `my-hackathon-project` |
| `VERTEX_AI_CREDENTIALS` | JSON de la service account de GCP | `{"type": "service_account", ...}` |
| `DB_HOST` | Host de la base de datos PostgreSQL | `db.example.com` |
| `DB_NAME` | Nombre de la base de datos | `hackathon_db` |
| `DB_USER` | Usuario de PostgreSQL | `postgres` |
| `DB_PASSWORD` | Contraseña de PostgreSQL | `secure_password_here` |
| `JWT_SECRET` | Secreto para firmar tokens JWT (min 32 chars) | `your-random-32-char-secret` |
| `API_SECRET_KEY` | Clave secreta de la API | `your-api-secret-key` |

### Uso en GitHub Actions

Los secretos se pueden usar en workflows de GitHub Actions así:

```yaml
env:
  GOOGLE_CLOUD_PROJECT: ${{ secrets.GOOGLE_CLOUD_PROJECT }}
  DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

### Rotación de Secretos

- 🔄 Rota los secretos de producción cada 90 días
- 🔄 Cambia inmediatamente cualquier secreto que pueda haber sido comprometido
- 📝 Mantén un registro de cuándo fueron cambiados los secretos

## 🚀 Desarrollo Local

### Prerrequisitos

- Go 1.21+ (para backend)
- Node.js 18+ y npm/yarn (para frontend)
- Docker y Docker Compose (para PostgreSQL)
- Python 3.8+ (para pre-commit)

### Inicio Rápido

1. **Clona el repositorio**:
   ```bash
   git clone https://github.com/Unikyri/gemini-hackathon.git
   cd gemini-hackathon
   ```

2. **Configura pre-commit**:
   ```bash
   pip install pre-commit
   pre-commit install
   ```

3. **Configura Backend**:
   ```bash
   cd backend
   cp .env.example .env
   # Edita .env con tus valores
   go mod download
   go run main.go
   ```

4. **Configura Frontend**:
   ```bash
   cd frontend
   cp .env.example .env
   # Edita .env con tus valores
   npm install
   npm start
   ```

5. **Inicia PostgreSQL con Docker**:
   ```bash
   docker run -d \
     --name hackathon-postgres \
     -e POSTGRES_DB=hackathon_db \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=your_password \
     -p 5432:5432 \
     postgres:16-alpine
   ```

### Estructura del Proyecto

```
gemini-hackathon/
├── backend/              # API REST en Go
│   ├── .env.example     # Template de variables de entorno
│   └── ...
├── frontend/            # Aplicación React
│   ├── .env.example     # Template de variables de entorno
│   └── ...
├── .gitignore          # Archivos ignorados por Git
├── .pre-commit-config.yaml  # Configuración de pre-commit hooks
└── README.md           # Este archivo
```

## 📚 Recursos Adicionales

- [Pre-commit Documentation](https://pre-commit.com/)
- [GitHub Secrets Documentation](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [React Documentation](https://react.dev/)
- [Vertex AI Documentation](https://cloud.google.com/vertex-ai/docs)

## 🤝 Contribuir

1. Asegúrate de tener pre-commit instalado y configurado
2. Crea una rama para tu feature: `git checkout -b feature/nueva-funcionalidad`
3. Haz commit de tus cambios (pre-commit se ejecutará automáticamente)
4. Push a la rama: `git push origin feature/nueva-funcionalidad`
5. Abre un Pull Request

## 📝 Licencia

Ver archivo [LICENSE](LICENSE) para más detalles.
