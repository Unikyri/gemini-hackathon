import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

interface InstructionsPanelProps {
  markdownContent: string;
  title: string;
}

/**
 * Panel izquierdo que muestra el enunciado del problema en formato Markdown
 */
export const InstructionsPanel = ({ markdownContent, title }: InstructionsPanelProps) => {
  return (
    <div className="h-full flex flex-col bg-white border-r border-gray-200">
      <div className="p-4 border-b border-gray-200">
        <h2 className="text-xl font-semibold text-gray-900">{title}</h2>
      </div>
      <div className="flex-1 overflow-y-auto p-6">
        <div className="prose prose-slate max-w-none">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>
            {markdownContent}
          </ReactMarkdown>
        </div>
      </div>
    </div>
  );
};

