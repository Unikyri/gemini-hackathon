# Learning Path Generator - Frontend

Frontend del proyecto Learning Path Generator construido con React, TypeScript y Vite siguiendo los principios de **Screaming Architecture**.

## 🚀 Stack Tecnológico

- **React 19.2.0** - Framework UI
- **TypeScript 5.9.3** - Tipado estático (modo estricto)
- **Vite 7.2.4** - Build tool y dev server
- **Tailwind CSS 4.1.18** - Framework CSS utility-first
- **Zustand 5.0.10** - Estado global ligero
- **Axios 1.13.2** - Cliente HTTP
- **Monaco Editor React 4.7.0** - Editor de código

## 📦 Instalación

```bash
# Instalar dependencias
npm install

# Copiar variables de entorno
cp .env.example .env

# Configurar la URL del API en .env
VITE_API_BASE_URL=http://localhost:8000/api
```

## 🏃 Scripts Disponibles

```bash
# Iniciar servidor de desarrollo
npm run dev

# Build para producción
npm run build

# Preview del build de producción
npm run preview

# Ejecutar linter
npm run lint
```

## 📁 Arquitectura

Este proyecto usa **Screaming Architecture**, donde la estructura refleja las características del negocio:

```
src/
├── features/              # Features del negocio
│   └── path-generator/    # Generación de paths
├── shared/                # Código compartido
│   ├── api/              # Cliente API
│   ├── hooks/            # Custom hooks
│   └── store/            # Estado global
└── ...
```

Ver [ARCHITECTURE.md](./ARCHITECTURE.md) para más detalles.

## 🔧 Configuración

### Variables de Entorno

Crea un archivo `.env` con:

```env
VITE_API_BASE_URL=http://localhost:8000/api
```

### TypeScript

TypeScript está configurado en modo estricto con las siguientes opciones:
- `strict: true`
- `noUnusedLocals: true`
- `noUnusedParameters: true`

### Path Aliases

El proyecto usa path aliases para imports limpios:

```typescript
import { apiService } from '@/shared/api';
import { usePathStore } from '@/shared/store';
import { PathGenerator } from '@/features/path-generator';
```

## 🔌 Cliente API

El cliente API proporciona métodos para interactuar con el backend:

```typescript
import { apiService } from '@/shared/api';

// Generar path
const response = await apiService.generatePath(prompt);

// Obtener path
const path = await apiService.getPath(pathId);

// Obtener nodo
const node = await apiService.getNode(pathId, nodeId);

// Actualizar nodo
await apiService.updateNodeCompletion(pathId, nodeId, completed);
```

## 🎣 Hooks Personalizados

### useGeneratePath
```typescript
const { generatePath, isGenerating, error } = useGeneratePath();
await generatePath("Quiero aprender React");
```

### usePath
```typescript
const { fetchPath, isLoading, error } = usePath();
await fetchPath(pathId);
```

### useNode
```typescript
const { fetchNode, updateNodeCompletion, isLoading, error } = useNode();
await fetchNode(pathId, nodeId);
```

## 📚 Documentación Adicional

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Arquitectura detallada del proyecto
- [Vite Documentation](https://vite.dev/)
- [React Documentation](https://react.dev/)
- [TypeScript Documentation](https://www.typescriptlang.org/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)

## 🤝 Contribuir

1. Seguir la estructura de screaming architecture
2. Usar TypeScript en modo estricto
3. Usar path aliases para imports
4. Crear hooks reutilizables en `shared/hooks/`
5. Mantener componentes pequeños y enfocados

## 📝 Licencia

Este proyecto es parte del Gemini Hackathon.

