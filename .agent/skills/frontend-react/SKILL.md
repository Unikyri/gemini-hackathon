---
name: Frontend React - TypeScript + Vite
description: Genera código React siguiendo los patrones del proyecto Gemini Coding Path
---

# Skill: Frontend React

## Contexto

Este proyecto usa **React 19** + **TypeScript** (strict mode) + **Vite** + **Tailwind CSS**.
El código debe seguir la estructura feature-based y los patrones establecidos.

---

## Stack del Frontend

| Tecnología | Versión | Propósito |
|------------|---------|-----------|
| React | 19.x | UI Library |
| TypeScript | 5.x | Type safety (strict mode) |
| Vite | 7.x | Build tool |
| Tailwind CSS | 4.x | Styling |
| Zustand | 5.x | State management |
| Axios | 1.x | HTTP client |
| Monaco Editor | 4.x | Code editor |
| Lucide React | latest | Icons |

---

## Estructura de Archivos

```
frontend/core/src/
├── features/                  # Módulos por funcionalidad
│   ├── path-generator/        # Feature de generación de paths
│   │   ├── components/        # Componentes del feature
│   │   │   └── PathGenerator.tsx
│   │   ├── hooks/             # Hooks específicos (opcional)
│   │   └── index.ts           # Exports públicos
│   └── node-workspace/        # Feature del workspace
│       ├── components/
│       │   ├── Workspace.tsx
│       │   ├── InstructionsPanel.tsx
│       │   ├── EditorPanel.tsx
│       │   └── OutputPanel.tsx
│       └── index.ts
├── shared/                    # Código compartido
│   ├── api/                   # Cliente API
│   │   └── client.ts
│   ├── components/            # Componentes reutilizables
│   │   └── ErrorMessage.tsx
│   ├── hooks/                 # Custom hooks globales
│   │   └── useApi.ts
│   └── store/                 # Zustand stores
│       └── pathStore.ts
├── App.tsx                    # Componente raíz
├── main.tsx                   # Entry point
└── index.css                  # Estilos globales (Tailwind)
```

---

## Patrones de Código

### Componente Funcional con TypeScript
```tsx
// features/path-generator/components/PathGenerator.tsx
import { useState } from 'react';
import { Sparkles, Loader2, AlertCircle } from 'lucide-react';
import { usePathStore } from '@/shared/store/pathStore';

interface PathGeneratorProps {
  onPathGenerated?: (pathId: string) => void;
}

export function PathGenerator({ onPathGenerated }: PathGeneratorProps) {
  const [prompt, setPrompt] = useState('');
  const { generatePath, isLoading, error } = usePathStore();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!prompt.trim()) return;
    
    try {
      const path = await generatePath(prompt);
      onPathGenerated?.(path.id);
    } catch (err) {
      // El error se maneja en el store
    }
  };

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <div className="flex items-center gap-2 mb-4">
        <Sparkles className="text-yellow-400" size={24} />
        <h1 className="text-2xl font-bold text-white">
          ¿Qué quieres aprender hoy?
        </h1>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          type="text"
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
          placeholder="Ej: Aprender Go desde cero"
          className="w-full px-4 py-3 bg-gray-800 border border-gray-700 rounded-lg 
                     text-white placeholder-gray-500 focus:outline-none 
                     focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
          disabled={isLoading}
        />

        {error && (
          <div className="flex items-center gap-2 text-red-400">
            <AlertCircle size={16} />
            <span>{error}</span>
          </div>
        )}

        <button
          type="submit"
          disabled={isLoading || !prompt.trim()}
          className="w-full px-4 py-3 bg-blue-600 text-white font-semibold rounded-lg
                     hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed
                     flex items-center justify-center gap-2 transition-colors"
        >
          {isLoading ? (
            <>
              <Loader2 className="animate-spin" size={20} />
              Generando...
            </>
          ) : (
            <>
              <Sparkles size={20} />
              Generar mi ruta
            </>
          )}
        </button>
      </form>
    </div>
  );
}
```

### Zustand Store
```tsx
// shared/store/pathStore.ts
import { create } from 'zustand';
import { apiClient } from '@/shared/api/client';

interface PathNode {
  id: string;
  title: string;
  description: string;
  status: 'locked' | 'unlocked' | 'completed';
  position: number;
}

interface Path {
  id: string;
  title: string;
  nodes: PathNode[];
}

interface PathState {
  currentPath: Path | null;
  isLoading: boolean;
  error: string | null;
  generatePath: (prompt: string) => Promise<Path>;
  fetchPath: (pathId: string) => Promise<void>;
  clearError: () => void;
}

export const usePathStore = create<PathState>((set, get) => ({
  currentPath: null,
  isLoading: false,
  error: null,

  generatePath: async (prompt: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.post('/paths/generate', { prompt });
      const path = response.data;
      set({ currentPath: path, isLoading: false });
      return path;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Error desconocido';
      set({ error: message, isLoading: false });
      throw err;
    }
  },

  fetchPath: async (pathId: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.get(`/paths/${pathId}`);
      set({ currentPath: response.data, isLoading: false });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Error desconocido';
      set({ error: message, isLoading: false });
    }
  },

  clearError: () => set({ error: null }),
}));
```

### Cliente API con Axios
```tsx
// shared/api/client.ts
import axios from 'axios';

const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

export const apiClient = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30 segundos (IA puede tardar)
});

// Interceptor para errores
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      // Error de servidor
      const message = error.response.data?.error || 'Error del servidor';
      return Promise.reject(new Error(message));
    } else if (error.request) {
      // Error de red
      return Promise.reject(new Error('No se pudo conectar con el servidor'));
    }
    return Promise.reject(error);
  }
);

// Funciones tipadas
export const pathApi = {
  generatePath: (prompt: string, userLevel?: string) =>
    apiClient.post('/paths/generate', { prompt, user_level: userLevel }),
  
  getPath: (pathId: string) =>
    apiClient.get(`/paths/${pathId}`),
  
  getNode: (nodeId: string) =>
    apiClient.get(`/nodes/${nodeId}`),
};
```

### Custom Hook
```tsx
// shared/hooks/useNodeDetail.ts
import { useState, useEffect } from 'react';
import { pathApi } from '@/shared/api/client';

interface NodeDetail {
  id: string;
  title: string;
  markdownContent: string;
  boilerplateCode: string;
  documentation?: string;
}

export function useNodeDetail(nodeId: string) {
  const [node, setNode] = useState<NodeDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!nodeId) return;

    const fetchNode = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await pathApi.getNode(nodeId);
        setNode(response.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar');
      } finally {
        setIsLoading(false);
      }
    };

    fetchNode();
  }, [nodeId]);

  return { node, isLoading, error };
}
```

### Monaco Editor Integration
```tsx
// features/node-workspace/components/EditorPanel.tsx
import Editor from '@monaco-editor/react';

interface EditorPanelProps {
  code: string;
  onChange: (value: string | undefined) => void;
  language?: string;
}

export function EditorPanel({ code, onChange, language = 'go' }: EditorPanelProps) {
  return (
    <div className="h-full w-full">
      <Editor
        height="100%"
        language={language}
        theme="vs-dark"
        value={code}
        onChange={onChange}
        options={{
          minimap: { enabled: false },
          fontSize: 14,
          lineNumbers: 'on',
          scrollBeyondLastLine: false,
          automaticLayout: true,
          tabSize: 4,
          wordWrap: 'on',
        }}
      />
    </div>
  );
}
```

---

## Clases Tailwind Comunes

### Layouts
```css
/* Container centrado */
.container: max-w-2xl mx-auto px-4

/* Flex horizontal con gap */
.flex-row: flex items-center gap-2

/* Grid de 3 columnas */
.grid-3: grid grid-cols-3 gap-4

/* Workspace 3 paneles */
.workspace: flex h-screen
.panel-left: w-1/4 p-4 bg-gray-900
.panel-center: w-1/2 h-full
.panel-right: w-1/4 p-4 bg-gray-900
```

### Componentes
```css
/* Botón primario */
.btn-primary: px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 
              disabled:opacity-50 disabled:cursor-not-allowed transition-colors

/* Input */
.input: w-full px-4 py-3 bg-gray-800 border border-gray-700 rounded-lg 
        text-white placeholder-gray-500 focus:outline-none 
        focus:border-blue-500 focus:ring-1 focus:ring-blue-500

/* Card */
.card: p-4 bg-gray-800 rounded-lg border border-gray-700

/* Error message */
.error: flex items-center gap-2 text-red-400 text-sm
```

### Tipografía
```css
/* Títulos */
.heading-1: text-2xl font-bold text-white
.heading-2: text-xl font-semibold text-white
.heading-3: text-lg font-medium text-gray-200

/* Texto */
.text-body: text-gray-300
.text-muted: text-gray-500
.text-error: text-red-400
.text-success: text-green-400
```

---

## Checklist para Nuevo Código

- [ ] TypeScript estricto (no `any`, interfaces definidas)
- [ ] Props tipadas con `interface`
- [ ] Estados de loading y error manejados
- [ ] Iconos de `lucide-react`
- [ ] Estilos con Tailwind (no CSS personalizado)
- [ ] Componente exportado desde `index.ts` del feature
- [ ] Lógica de API en stores o hooks, no en componentes
- [ ] Compila sin errores (`npm run build`)
