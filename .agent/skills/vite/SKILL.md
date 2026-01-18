---
name: Vite Configuration
description: Configuración y optimización de Vite para Gemini Coding Path
---

# Skill: Vite

## Contexto

Este proyecto usa **Vite 7.x** como build tool para el frontend React.
Configuración optimizada para desarrollo rápido y builds de producción eficientes.

---

## Configuración Base

```typescript
// frontend/core/vite.config.ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  
  // Path aliases
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@features': path.resolve(__dirname, './src/features'),
      '@shared': path.resolve(__dirname, './src/shared'),
    },
  },
  
  // Server de desarrollo
  server: {
    port: 5173,
    host: true, // Permite acceso desde red local
    proxy: {
      // Proxy para evitar CORS en desarrollo
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  
  // Build de producción
  build: {
    outDir: 'dist',
    sourcemap: false, // true para debugging
    minify: 'esbuild',
    rollupOptions: {
      output: {
        manualChunks: {
          // Separar vendors para mejor caching
          vendor: ['react', 'react-dom'],
          editor: ['@monaco-editor/react'],
          ui: ['lucide-react'],
        },
      },
    },
  },
  
  // Optimización de dependencias
  optimizeDeps: {
    include: ['react', 'react-dom', 'zustand', 'axios'],
    exclude: ['@monaco-editor/react'], // Monaco se carga async
  },
})
```

---

## TypeScript con Path Aliases

```json
// frontend/core/tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    
    /* Bundler mode */
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    
    /* Linting - STRICT MODE */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    
    /* Path aliases */
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@features/*": ["src/features/*"],
      "@shared/*": ["src/shared/*"]
    }
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

---

## Variables de Entorno

```bash
# frontend/core/.env.example
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_NAME=Gemini Coding Path
VITE_ENABLE_DEVTOOLS=true
```

### Uso en código
```typescript
// Las variables deben tener prefijo VITE_
const apiUrl = import.meta.env.VITE_API_BASE_URL;
const appName = import.meta.env.VITE_APP_NAME;

// Tipado para IntelliSense
// src/vite-env.d.ts
/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_APP_NAME: string
  readonly VITE_ENABLE_DEVTOOLS: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
```

---

## Scripts de package.json

```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0",
    "lint:fix": "eslint . --ext ts,tsx --fix",
    "type-check": "tsc --noEmit"
  }
}
```

---

## Comandos Comunes

```bash
# Desarrollo con HMR
npm run dev

# Build de producción
npm run build

# Preview del build
npm run preview

# Solo type-check (sin build)
npm run type-check

# Linting
npm run lint
npm run lint:fix
```

---

## Optimizaciones de Producción

### Code Splitting por Rutas
```typescript
// Con React.lazy para rutas
import { lazy, Suspense } from 'react';

const Workspace = lazy(() => import('@features/node-workspace'));
const PathGenerator = lazy(() => import('@features/path-generator'));

function App() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      {/* Rutas */}
    </Suspense>
  );
}
```

### Loading de Monaco Editor
```typescript
// Monaco se carga de forma asíncrona automáticamente
import Editor from '@monaco-editor/react';

// Componente de loading personalizado
<Editor
  loading={<div>Cargando editor...</div>}
  // ... otras props
/>
```

---

## Checklist

- [ ] Path aliases configurados en vite.config.ts Y tsconfig.json
- [ ] Variables de entorno con prefijo VITE_
- [ ] Proxy configurado para desarrollo
- [ ] Build optimizado con chunks separados
- [ ] TypeScript en modo estricto
