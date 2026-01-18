import { useNodeDetail } from '@/shared/hooks';
import { InstructionsPanel } from './InstructionsPanel';
import { CodeEditor } from './CodeEditor';

interface NodeWorkspaceProps {
  pathId: string;
  nodeId: string;
}

/**
 * El Workspace de Código (The Dojo)
 * Componente principal que muestra el entorno de desarrollo integrado
 * con panel de instrucciones, editor de código y acciones
 */
export const NodeWorkspace = ({ pathId, nodeId }: NodeWorkspaceProps) => {
  const { node, isLoading, error } = useNodeDetail({ pathId, nodeId });

  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Cargando ejercicio...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center max-w-md">
          <div className="text-red-600 text-5xl mb-4">⚠️</div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Error</h2>
          <p className="text-gray-600">{error}</p>
        </div>
      </div>
    );
  }

  if (!node) {
    return (
      <div className="h-screen flex items-center justify-center bg-gray-50">
        <p className="text-gray-600">No se encontró el ejercicio</p>
      </div>
    );
  }

  const markdownContent = node.markdown_content || node.content || node.description;
  const boilerplateCode = node.boilerplate || '// Escribe tu código aquí\n';

  return (
    <div className="h-screen flex flex-col bg-gray-50">
      {/* Header con título del nodo */}
      <header className="bg-white border-b border-gray-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{node.title}</h1>
            <p className="text-sm text-gray-500 mt-1">
              Ejercicio {node.order} • {node.completed ? '✓ Completado' : 'En progreso'}
            </p>
          </div>
          <div className="flex gap-3">
            <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium">
              Ejecutar Tests
            </button>
            <button className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors font-medium">
              Guardar Borrador
            </button>
          </div>
        </div>
      </header>

      {/* Main workspace con 3 paneles */}
      <div className="flex-1 flex overflow-hidden">
        {/* Panel Izquierdo: Instrucciones */}
        <div className="w-1/3 min-w-[300px] max-w-[500px] overflow-hidden">
          <InstructionsPanel 
            markdownContent={markdownContent} 
            title="Enunciado"
          />
        </div>

        {/* Panel Central: Editor de código */}
        <div className="flex-1 overflow-hidden">
          <CodeEditor 
            initialCode={boilerplateCode}
            language="go"
          />
        </div>

        {/* Panel Derecho: Documentación y Misiones (futuro) */}
        <div className="w-1/4 min-w-[250px] max-w-[400px] bg-white border-l border-gray-200">
          <div className="p-4">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Recursos</h3>
            <p className="text-sm text-gray-500">Próximamente: Documentación y misiones secundarias</p>
          </div>
        </div>
      </div>
    </div>
  );
};

