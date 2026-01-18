---
name: OWASP ASVS 5.0
description: Controles de seguridad seleccionados por tipo de tarea para Gemini Coding Path
---

# Skill: OWASP ASVS 5.0

## Contexto

Este proyecto aplica **OWASP ASVS 5.0** de forma selectiva.
No usamos TODOS los controles - seleccionamos los relevantes según la tarea.

**Niveles ASVS:**
- **L1**: Mínimo - cualquier aplicación
- **L2**: Estándar - aplicaciones con datos sensibles
- **L3**: Avanzado - aplicaciones críticas

> **Gemini Coding Path aplica L1 + controles L2 selectos** (datos de usuarios, código ejecutado)

---

## 📋 Controles por Tipo de Tarea

### 🔐 Autenticación y Sesiones (Sprint 2 - E01)

| Control | Descripción | Nivel |
|---------|-------------|-------|
| **V6.2.1** | Contraseñas mínimo 8 caracteres (recomendado 15) | L1 |
| **V6.2.4** | Verificar contra top 3000 contraseñas comunes | L1 |
| **V6.2.5** | No restringir tipos de caracteres en contraseñas | L1 |
| **V6.3.1** | Rate limiting contra credential stuffing | L1 |
| **V6.3.2** | No usar cuentas default (root, admin) | L1 |
| **V7.2.1** | Verificar tokens de sesión en backend | L1 |
| **V7.2.3** | Tokens de sesión con 128 bits de entropía | L1 |
| **V7.2.4** | Nuevo token al autenticar (regenerar sesión) | L1 |
| **V7.4.1** | Invalidar sesión en logout/expiración | L1 |

```go
// Ejemplo: V6.3.1 - Rate limiting
type RateLimiter struct {
    maxAttempts int
    window      time.Duration
}

func (rl *RateLimiter) Allow(key string) bool {
    // Implementar usando Redis o memoria
    attempts := rl.getAttempts(key)
    if attempts >= rl.maxAttempts {
        return false
    }
    rl.increment(key)
    return true
}
```

---

### 🛡️ Inyección y Sanitización (Backend)

| Control | Descripción | Nivel | Aplica a |
|---------|-------------|-------|----------|
| **V1.2.1** | Output encoding para HTML/XML | L1 | Respuestas API |
| **V1.2.4** | Queries parametrizadas (SQL injection) | L1 | GORM queries |
| **V1.2.5** | Proteger contra OS command injection | L1 | Docker sandbox |
| **V1.3.1** | Sanitizar HTML de WYSIWYG | L1 | Markdown rendering |
| **V1.5.1** | XML parsers con config restrictiva (XXE) | L1 | Si usamos XML |

```go
// Ejemplo: V1.2.4 - GORM ya usa parametrización
// ✅ CORRECTO
db.Where("id = ?", userInput).First(&path)

// ❌ INCORRECTO - vulnerable a SQL injection
db.Raw("SELECT * FROM paths WHERE id = " + userInput)
```

---

### ✅ Validación de Entrada

| Control | Descripción | Nivel | Aplica a |
|---------|-------------|-------|----------|
| **V2.2.1** | Validar input contra reglas de negocio | L1 | Prompt, código |
| **V2.2.2** | Validar en backend, no confiar en frontend | L1 | Todos los endpoints |
| **V2.3.1** | Flujos de negocio en orden secuencial | L1 | Path → Node → Submit |

```go
// Ejemplo: V2.2.1 + V2.2.2 - Validación en backend
func ValidatePrompt(prompt string) error {
    prompt = strings.TrimSpace(prompt)
    
    if len(prompt) < 3 {
        return errors.New("prompt must be at least 3 characters")
    }
    if len(prompt) > 500 {
        return errors.New("prompt must be at most 500 characters")
    }
    
    // Validar contra patrones peligrosos
    if containsInjectionPatterns(prompt) {
        return errors.New("invalid characters in prompt")
    }
    
    return nil
}
```

---

### 🌐 Seguridad HTTP/Headers (Frontend + Backend)

| Control | Descripción | Nivel | Implementar en |
|---------|-------------|-------|----------------|
| **V3.3.1** | Cookies con 'Secure' attribute | L1 | Backend (Gin) |
| **V3.4.1** | HSTS header (min 1 año) | L1 | Nginx/Gin |
| **V3.4.2** | CORS con allowlist de origins | L1 | Gin middleware |
| **V3.5.1** | Anti-CSRF tokens o validar Origin | L1 | Formularios |
| **V4.1.1** | Content-Type header correcto | L1 | Todas las respuestas |

```go
// Ejemplo: Middleware de seguridad en Gin
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // V3.4.1 - HSTS
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // V3.4.4 - No MIME sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // V3.4.6 - Frame protection
        c.Header("Content-Security-Policy", "frame-ancestors 'none'")
        
        // V4.1.1 - Content-Type
        c.Header("Content-Type", "application/json; charset=utf-8")
        
        c.Next()
    }
}
```

---

### 🔒 Criptografía y Secretos

| Control | Descripción | Nivel | Aplica a |
|---------|-------------|-------|----------|
| **V11.3.1** | No usar ECB mode ni padding inseguro | L1 | Si ciframos datos |
| **V11.3.2** | Usar AES-GCM u otros aprobados | L1 | Cifrado de datos |
| **V11.4.1** | No usar MD5/SHA1 para seguridad | L1 | Hashing |
| **V11.4.2** | Passwords con bcrypt/argon2 | L2 | Sprint 2 - Auth |
| **V11.5.1** | Usar CSPRNG para tokens (128 bits) | L2 | Session tokens |
| **V13.3.1** | Secretos en vault, no en código | L2 | API keys, DB creds |

```go
// Ejemplo: V11.4.2 - Password hashing con bcrypt
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

---

### 🐳 Docker Sandbox (E05 - Code Runner)

| Control | Descripción | Nivel | Implementación |
|---------|-------------|-------|----------------|
| **V1.2.5** | Proteger contra OS command injection | L1 | No ejecutar comandos de usuario directamente |
| **V1.3.6** | Proteger contra SSRF | L2 | network_mode: none |
| **V5.2.1** | Limitar tamaño de archivos | L1 | Límite de código: 50KB |
| **V5.3.1** | No ejecutar archivos subidos como código | L1 | Sandbox aislado |
| **V5.3.2** | No usar paths de usuario en file operations | L1 | Paths hardcodeados |
| **V15.2.5** | Sandboxing de componentes peligrosos | L3 | Docker isolado |

```yaml
# docker-compose.sandbox.yml - Aplicando ASVS
services:
  code-runner:
    # V1.2.5 - Sin acceso a shell del host
    network_mode: none  # V1.3.6 - Sin red (anti-SSRF)
    
    # V5.2.1 - Límites de recursos
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 128M
    
    # V15.2.5 - Sandboxing
    read_only: true
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
```

---

### 📝 Logging y Errores

| Control | Descripción | Nivel | Aplica a |
|---------|-------------|-------|----------|
| **V16.2.1** | Logs con metadata (when, where, who, what) | L2 | Todos los logs |
| **V16.2.5** | No loguear credenciales ni tokens | L2 | Auth, sesiones |
| **V16.3.1** | Loguear intentos de auth (éxito y fallo) | L2 | Sprint 2 |
| **V16.5.1** | Mensajes de error genéricos al usuario | L2 | API responses |

```go
// Ejemplo: V16.5.1 - Error handling seguro
func (c *PathController) GeneratePath(ctx *gin.Context) {
    path, err := c.useCase.Execute(ctx, req.Prompt)
    if err != nil {
        // Loguear error real internamente (V16.2.1)
        log.Error().
            Str("user_id", ctx.GetString("user_id")).
            Str("prompt", req.Prompt).
            Err(err).
            Msg("failed to generate path")
        
        // Retornar mensaje genérico (V16.5.1)
        ctx.JSON(500, dto.ErrorResponse{
            Error: "Unable to generate path. Please try again.",
        })
        return
    }
    // ...
}
```

---

### 🔄 TLS/Comunicación Segura

| Control | Descripción | Nivel | Aplica a |
|---------|-------------|-------|----------|
| **V12.1.1** | Solo TLS 1.2+ | L1 | Producción |
| **V12.2.1** | TLS para todas las conexiones externas | L1 | Frontend ↔ Backend |
| **V12.2.2** | Certificados públicos de confianza | L1 | Producción |
| **V12.3.1** | TLS para conexiones internas también | L2 | Backend ↔ DB |

---

### 📦 Dependencias y Terceros

| Control | Descripción | Nivel | Aplica a |
|---------|-------------|-------|----------|
| **V15.1.1** | Remediation timeframes para vulnerabilidades | L1 | npm, go mod |
| **V15.1.2** | Inventario SBOM de dependencias | L2 | CI/CD |
| **V15.2.1** | Componentes sin vulnerabilidades conocidas | L1 | Dependabot |
| **V13.4.1** | No exponer .git en producción | L1 | Dockerfile |

```bash
# Ejemplo: V15.2.1 - Escaneo de vulnerabilidades en CI
# .github/workflows/security.yml
- name: Go vulnerability check
  run: govulncheck ./...

- name: NPM audit
  run: npm audit --audit-level=high
```

---

## 🎯 Controles por Sprint

### Sprint 1 (E02 + E04)
| Prioridad | Control | Tarea relacionada |
|-----------|---------|-------------------|
| Alta | V1.2.4 | Task 2 - Queries parametrizadas |
| Alta | V2.2.1, V2.2.2 | Task 4 - Validación de prompt |
| Alta | V3.4.2 | Task 5 - CORS middleware |
| Alta | V4.1.1 | Task 5 - Content-Type headers |
| Alta | V1.3.1 | Task 9 - Sanitizar Markdown |
| Media | V13.4.1 | Task 12 - No .git en Docker |
| Media | V15.2.1 | Task 13 - Dependabot |

### Sprint 2 (E01 + E05)
| Prioridad | Control | Tarea relacionada |
|-----------|---------|-------------------|
| Alta | V6.2.1-V6.2.5 | Auth - Password policy |
| Alta | V6.3.1 | Auth - Rate limiting |
| Alta | V7.2.1-V7.4.1 | Auth - Session management |
| Alta | V1.2.5 | Code Runner - Command injection |
| Alta | V1.3.6 | Code Runner - SSRF protection |
| Alta | V11.4.2 | Auth - Password hashing |

### Sprint 3 (E03 + E06)
| Prioridad | Control | Tarea relacionada |
|-----------|---------|-------------------|
| Media | V16.3.1 | Logging de auth |
| Media | V13.3.1 | Secrets management |
| Baja | V3.7.4 | HSTS preload list |

---

## ✅ Checklist por Pull Request

Antes de mergear, verificar:

```markdown
## Security Checklist (ASVS)

### Input/Output
- [ ] V1.2.4: Queries parametrizadas (no concatenación)
- [ ] V2.2.2: Validación en backend

### HTTP
- [ ] V3.4.2: CORS configurado correctamente
- [ ] V4.1.1: Content-Type correcto

### Secrets
- [ ] V13.3.1: Sin secretos hardcodeados
- [ ] V13.4.1: Sin archivos .git expuestos

### Errors
- [ ] V16.5.1: Mensajes de error genéricos
```

---

## 📚 Referencia Rápida

Para ver el control completo, consultar:
- [OWASP ASVS 5.0 CSV](docs/OWASP_Application_Security_Verification_Standard_5.0.0_en.csv)
- [OWASP ASVS Project](https://owasp.org/www-project-application-security-verification-standard/)
