# Component API Reference - NodeWorkspace

## NodeWorkspace Component

Main container component that implements "The Dojo" - the integrated development environment in the browser.

### Props

```typescript
interface NodeWorkspaceProps {
  pathId: string;  // UUID of the learning path
  nodeId: string;  // UUID of the current exercise node
}
```

### Usage

```tsx
import { NodeWorkspace } from '@/features/node-workspace';

<NodeWorkspace 
  pathId="550e8400-e29b-41d4-a716-446655440000"
  nodeId="660e8400-e29b-41d4-a716-446655440000" 
/>
```

### Features

- Fetches node data automatically via `useNodeDetail` hook
- Displays loading state while fetching
- Shows error state if fetch fails
- Renders 3-column layout when data loads successfully
- Injects boilerplate code into Monaco Editor

---

## CodeEditor Component

Monaco Editor wrapper configured for Go programming with dark theme.

### Props

```typescript
interface CodeEditorProps {
  initialCode?: string;      // Initial code to display (boilerplate)
  language?: string;         // Programming language (default: 'go')
  onChange?: (value: string | undefined) => void; // Callback when code changes
}
```

### Usage

```tsx
import { CodeEditor } from '@/features/node-workspace';

<CodeEditor 
  initialCode="package main\n\nfunc main() {\n\t// Your code here\n}"
  language="go"
  onChange={(code) => console.log('Code changed:', code)}
/>
```

### Configuration

- **Language**: Go (configurable)
- **Theme**: VS Dark
- **Minimap**: Disabled
- **Font Size**: 14px
- **Word Wrap**: Enabled
- **Line Numbers**: Enabled
- **Auto Layout**: Enabled (responsive)
- **Padding**: 16px top and bottom

### Behavior

- Updates automatically when `initialCode` prop changes (via useEffect)
- Calls `onChange` when user edits the code
- Preserves exact boilerplate content (no trimming, no modifications)
- Handles special characters, unicode, and various line endings

---

## InstructionsPanel Component

Panel that displays exercise instructions in Markdown format.

### Props

```typescript
interface InstructionsPanelProps {
  markdownContent: string;  // Markdown text to render
  title: string;            // Panel title
}
```

### Usage

```tsx
import { InstructionsPanel } from '@/features/node-workspace';

<InstructionsPanel 
  markdownContent="# Exercise 1\n\nImplement a function..."
  title="Enunciado"
/>
```

### Features

- Full GFM (GitHub Flavored Markdown) support via `remark-gfm`
- Styled with Tailwind Typography (`prose`)
- Scrollable content area
- Headers, lists, code blocks, emphasis, links, etc.

---

## Component Tree

```
<NodeWorkspace>
  │
  ├── Loading State (if isLoading)
  │   └── Spinner + "Cargando ejercicio..."
  │
  ├── Error State (if error)
  │   └── Error icon + message
  │
  └── Main UI (if node loaded)
      │
      ├── <Header>
      │   ├── Node title + status
      │   └── Action buttons
      │       ├── "Ejecutar Tests" button
      │       └── "Guardar Borrador" button
      │
      └── <Main Layout> (flex row)
          │
          ├── <InstructionsPanel> (w-1/4)
          │   ├── Panel header with title
          │   └── Markdown content (scrollable)
          │
          ├── <CodeEditor> (w-1/2)
          │   └── <MonacoEditor>
          │       └── Go code with dark theme
          │
          └── <InfoPanel> (w-1/4)
              ├── "Info / Salida" header
              └── Sections
                  ├── Documentación
                  ├── Misiones Secundarias
                  └── Resultados
```

---

## Data Flow

```
User navigates to exercise
    ↓
<NodeWorkspace pathId={...} nodeId={...}>
    ↓
useNodeDetail({ pathId, nodeId })
    ↓
API: GET /api/paths/{pathId}/nodes/{nodeId}
    ↓
Response: {
  node_id: string
  title: string
  description: string
  content: string
  markdown_content: string
  boilerplate: string  ← Goes to CodeEditor
  order: number
  completed: boolean
  path_id: string
}
    ↓
NodeWorkspace renders:
  - InstructionsPanel (markdown_content)
  - CodeEditor (boilerplate)
  - InfoPanel (placeholder)
    ↓
User edits code in Monaco Editor
    ↓
onChange callback fires
    ↓
Code state updates
```

---

## Styling

All components use Tailwind CSS for styling:

### NodeWorkspace
- `h-screen flex flex-col` - Full viewport height with flex column
- `bg-gray-50` - Light gray background

### Layout
- `flex-1 flex` - Main content area with flex row
- `w-1/4`, `w-1/2`, `w-1/4` - Column proportions (25-50-25)
- `min-w-[250px]` - Minimum widths to prevent collapse
- `overflow-hidden` - Prevent overflow, enable scroll in children

### Panels
- `border-r border-gray-200` - Right borders for separation
- `overflow-y-auto` - Vertical scroll when needed
- `p-4`, `p-6` - Padding for content

### Header
- `bg-white border-b` - White background with bottom border
- `px-6 py-4` - Padding for header content

---

## Accessibility

- Semantic HTML elements (`<header>`, `<main>`, `<section>`)
- Clear visual hierarchy
- Readable font sizes
- High contrast (dark theme on editor)
- Keyboard navigation support (Monaco Editor built-in)

---

## Performance Considerations

1. **Monaco Editor**: Loaded lazily, only mounts when needed
2. **Markdown Rendering**: Incremental, no blocking
3. **State Management**: Local state with hooks, no unnecessary re-renders
4. **Code Splitting**: Monaco Editor is code-split automatically
5. **Build Size**: Optimized bundle (~236 KB, gzipped 77.65 KB)

---

## Testing

All components have comprehensive property-based tests:

- `InstructionsPanel.property.test.tsx` - 4 tests
- `NodeWorkspace.property.test.tsx` - 5 tests
- `CodeEditor.property.test.tsx` - 9 tests

Run tests: `npm test`

---

## Future Enhancements

1. **Code Execution**
   - Connect "Ejecutar Tests" button to backend
   - Display results in InfoPanel

2. **Tab Navigation**
   - Multiple tabs in right panel
   - Documentation, Missions, Results tabs

3. **Resizable Panels**
   - Drag dividers to resize panels
   - Save preferences

4. **Code Persistence**
   - Auto-save to localStorage
   - "Guardar Borrador" saves to backend

5. **Keyboard Shortcuts**
   - Run tests: Ctrl+Enter
   - Save: Ctrl+S
   - Focus panels: Ctrl+1/2/3

