# 🔐 Configuración de GitHub Secrets

Esta guía explica cómo configurar los secretos necesarios para el CI/CD del proyecto.

## Secretos Requeridos

| Secreto | Uso | Obligatorio |
|---------|-----|-------------|
| `CODECOV_TOKEN` | Subir reportes de cobertura | Opcional |
| `GOOGLE_CLOUD_PROJECT` | ID del proyecto en GCP | Sprint 2 |
| `VERTEXAI_SERVICE_ACCOUNT` | Credenciales de Vertex AI | Sprint 2 |

---

## Configurar CODECOV_TOKEN (Opcional)

[Codecov](https://codecov.io) es un servicio que visualiza la cobertura de tests.

### Pasos:

1. **Ir a Codecov:**
   - Visita [https://codecov.io](https://codecov.io)
   - Inicia sesión con GitHub

2. **Agregar el repositorio:**
   - Busca `gemini-hackathon` en la lista
   - Click en "Setup repo"

3. **Copiar el token:**
   - Codecov te mostrará un token único
   - Cópialo

4. **Agregar como GitHub Secret:**
   - Ve a tu repo en GitHub
   - `Settings` → `Secrets and variables` → `Actions`
   - Click `New repository secret`
   - **Name:** `CODECOV_TOKEN`
   - **Value:** (pegar el token)
   - Click `Add secret`

### Verificación:
El badge de cobertura aparecerá en los PRs después del primer push a `main` o `develop`.

> **Nota:** Si no configuras CODECOV_TOKEN, el CI NO fallará (está configurado con `fail_ci_if_error: false`), simplemente no se subirán reportes de cobertura.

---

## Configurar Secrets para Vertex AI (Sprint 2)

Cuando implementemos la integración con Vertex AI, necesitarás:

### 1. GOOGLE_CLOUD_PROJECT
```
El ID de tu proyecto en Google Cloud Console
Ejemplo: gemini-hackathon-123456
```

### 2. VERTEXAI_SERVICE_ACCOUNT
```
El JSON de la cuenta de servicio con permisos para Vertex AI.
Contenido del archivo service-account.json (todo en una línea).
```

### Pasos para obtener las credenciales:

1. **Ir a Google Cloud Console:**
   - [https://console.cloud.google.com](https://console.cloud.google.com)

2. **Crear cuenta de servicio:**
   - `IAM & Admin` → `Service Accounts`
   - `Create Service Account`
   - Nombre: `gemini-hackathon-backend`
   - Rol: `Vertex AI User`

3. **Generar clave JSON:**
   - Click en la cuenta creada
   - `Keys` → `Add Key` → `Create new key`
   - Tipo: JSON
   - Descargar y guardar seguro

4. **Agregar como GitHub Secret:**
   - Abrir el JSON descargado
   - Copiar TODO el contenido
   - Pegarlo como valor del secret `VERTEXAI_SERVICE_ACCOUNT`

---

## Cómo Agregar un Secret en GitHub

```
1. Ir a: https://github.com/[tu-usuario]/gemini-hackathon
2. Settings (pestaña)
3. Secrets and variables → Actions
4. New repository secret
5. Completar Name y Value
6. Add secret
```

![GitHub Secrets Location](https://docs.github.com/assets/cb-28263/images/help/repository/repo-actions-secrets.png)

---

## Verificar que los Secrets Funcionan

Los secrets se enmascaran automáticamente en los logs de GitHub Actions.
Si ves `***` en los logs, significa que el secret se está usando correctamente.

---

## Seguridad

⚠️ **NUNCA** hagas lo siguiente:
- Commitear secretos en el código
- Compartir secretos en Discord/Slack
- Guardar secretos en archivos no ignorados por .gitignore

✅ **SIEMPRE** usa GitHub Secrets para credenciales.
