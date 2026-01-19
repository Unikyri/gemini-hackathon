# Sandbox Runner - Gemini Coding Path

## Propósito

Este sandbox ejecuta código de usuarios en un entorno **completamente aislado** para evitar:
- Ataques a otros servicios (SSRF)
- Consumo excesivo de recursos (DoS)
- Ejecución de código malicioso
- Exfiltración de datos

## Controles de Seguridad Implementados

| Control | Implementación | OWASP ASVS |
|---------|----------------|------------|
| Sin red | `network_mode: none` | V1.3.6 |
| Usuario no-root | `USER runner` (UID 1001) | V15.2.5 |
| Límite CPU | `cpus: 0.5` | V5.2.1 |
| Límite RAM | `memory: 128M` | V5.2.1 |
| Solo lectura | `read_only: true` | V5.3.2 |
| Sin privilegios | `no-new-privileges: true` | V15.2.5 |
| Límite PIDs | `pids_limit: 50` | V5.2.1 |
| Timeout | 30 segundos | V5.2.1 |

## Archivos

```
sandbox/
├── Dockerfile.runner     # Imagen del sandbox
├── README.md             # Esta documentación
├── submissions/          # Código del usuario (se monta aquí)
└── tests/                # Tests generados por IA
```

## Uso

### Build de la imagen
```bash
docker compose -f docker-compose.sandbox.yml build
```

### Ejecutar tests
```bash
# Montar código en sandbox/submissions/ y luego:
docker compose -f docker-compose.sandbox.yml run --rm code-runner
```

### Desde Go (backend)
```go
func ExecuteUserCode(codeDir string) (*ExecutionResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
    defer cancel()
    
    cmd := exec.CommandContext(ctx, "docker", "compose",
        "-f", "docker-compose.sandbox.yml",
        "run", "--rm", "code-runner")
    
    output, err := cmd.CombinedOutput()
    
    return &ExecutionResult{
        Output:   string(output),
        ExitCode: cmd.ProcessState.ExitCode(),
        Passed:   err == nil,
    }, nil
}
```

## Verificación de Seguridad

### Probar aislamiento de red
```bash
# Dentro del contenedor, esto DEBE fallar:
docker compose -f docker-compose.sandbox.yml run --rm code-runner \
    sh -c "ping -c 1 google.com"
# Error: network is unreachable ✓
```

### Probar límites de memoria
```bash
# Código que intenta usar más de 128MB será matado
# Exit code: 137 (OOMKilled)
```

### Probar límite de tiempo
```bash
# Código que excede 30 segundos será terminado
# Exit code: 124 (timeout)
```

## Flujo de Ejecución

```
┌─────────────────┐
│   Usuario       │
│   envía código  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Backend       │
│   valida código │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│   Sandbox Container                      │
│   ┌──────────────────────────────────┐  │
│   │ Usuario: runner (no root)         │  │
│   │ Red: deshabilitada                │  │
│   │ CPU: máx 0.5                      │  │
│   │ RAM: máx 128MB                    │  │
│   │ Timeout: 30s                      │  │
│   │ Filesystem: read-only             │  │
│   └──────────────────────────────────┘  │
│                                          │
│   go test -v -timeout 30s ./...          │
│                                          │
└────────┬────────────────────────────────┘
         │
         ▼
┌─────────────────┐
│   Backend       │
│   captura output│
│   y exit code   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Frontend      │
│   muestra       │
│   resultados    │
└─────────────────┘
```
