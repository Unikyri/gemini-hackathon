---
description: Crear un nuevo componente React siguiendo los patrones del proyecto
---

# Workflow: Nuevo Componente React

## Prerrequisitos
- Frontend configurado con Vite + React + TypeScript
- Tailwind CSS configurado
- Feature folder existente o a crear

## Pasos

### 1. Determinar ubicación del componente

**Componente de feature específica:**
```bash
frontend/core/src/features/{nombre-feature}/components/
```

**Componente compartido:**
```bash
frontend/core/src/shared/components/
```

### 2. Crear el archivo del componente
```bash
# Crear archivo:
NombreComponente.tsx
```

### 3. Estructura base del componente
```tsx
import { useState } from 'react';
// Imports de lucide-react para iconos
import { IconName } from 'lucide-react';

interface NombreComponenteProps {
  // Props tipadas
  title: string;
  onAction?: () => void;
}

export function NombreComponente({ title, onAction }: NombreComponenteProps) {
  // Estados locales
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Handlers
  const handleClick = () => {
    if (onAction) {
      onAction();
    }
  };

  // Render
  return (
    <div className="p-4 rounded-lg bg-gray-800">
      <h2 className="text-xl font-bold text-white">{title}</h2>
      {error && (
        <p className="text-red-500">{error}</p>
      )}
      <button
        onClick={handleClick}
        disabled={isLoading}
        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
      >
        {isLoading ? 'Cargando...' : 'Acción'}
      </button>
    </div>
  );
}
```

### 4. Si necesita estado global, crear/usar store Zustand
```bash
# Crear/editar:
frontend/core/src/shared/store/nombreStore.ts
```
```tsx
import { create } from 'zustand';

interface NombreState {
  data: DataType | null;
  isLoading: boolean;
  error: string | null;
  fetchData: () => Promise<void>;
}

export const useNombreStore = create<NombreState>((set) => ({
  data: null,
  isLoading: false,
  error: null,
  fetchData: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await api.getData();
      set({ data: response, isLoading: false });
    } catch (err) {
      set({ error: 'Error al cargar', isLoading: false });
    }
  },
}));
```

### 5. Si necesita llamadas API, crear función en api client
```bash
# Editar:
frontend/core/src/shared/api/client.ts
```

### 6. Exportar desde index del feature
```bash
# Editar:
frontend/core/src/features/{feature}/index.ts
```
```tsx
export { NombreComponente } from './components/NombreComponente';
```

// turbo
### 7. Verificar que compila
```bash
cd frontend/core && npm run build
```

## Patrones del Proyecto

### Iconos
Usar `lucide-react`:
```tsx
import { Sparkles, Loader2, AlertCircle } from 'lucide-react';
```

### Estados de carga
```tsx
{isLoading ? (
  <Loader2 className="animate-spin" />
) : (
  <span>Contenido</span>
)}
```

### Manejo de errores
```tsx
{error && (
  <div className="flex items-center gap-2 text-red-500">
    <AlertCircle size={16} />
    <span>{error}</span>
  </div>
)}
```

### Clases Tailwind comunes
- Container: `p-4 rounded-lg bg-gray-800`
- Botón primario: `px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700`
- Botón disabled: `disabled:opacity-50 disabled:cursor-not-allowed`
- Texto título: `text-xl font-bold text-white`
- Texto error: `text-red-500`

## Checklist Final
- [ ] Componente creado con TypeScript estricto
- [ ] Props tipadas con interface
- [ ] Estados de loading y error manejados
- [ ] Iconos de lucide-react usados
- [ ] Exportado desde index del feature
- [ ] Compila sin errores
