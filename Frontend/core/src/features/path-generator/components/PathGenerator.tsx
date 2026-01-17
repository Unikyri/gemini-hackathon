import { useState } from 'react';
import { Sparkles, Loader2, AlertCircle, BookOpen } from 'lucide-react';
import { useGeneratePath } from '@/shared/hooks';

export const PathGenerator = () => {
  const [prompt, setPrompt] = useState('');
  const { generatePath, isGenerating, error } = useGeneratePath();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!prompt.trim()) return;

    const result = await generatePath(prompt);
    if (result) {
      console.log('Path generado:', result);
      // Aquí puedes redirigir o mostrar el path generado
    }
  };

  return (
    <div className="w-full max-w-2xl mx-auto p-6">
      <div className="flex items-center gap-3 mb-6">
        <BookOpen className="w-8 h-8 text-blue-600" />
        <h2 className="text-2xl font-bold">Generar Path de Aprendizaje</h2>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label 
            htmlFor="prompt" 
            className="block text-sm font-medium mb-2"
          >
            ¿Qué quieres aprender?
          </label>
          <textarea
            id="prompt"
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            placeholder="Ejemplo: Quiero aprender React desde cero..."
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows={4}
            disabled={isGenerating}
          />
        </div>

        {error && (
          <div className="flex items-center gap-2 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg">
            <AlertCircle className="w-5 h-5 flex-shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <button
          type="submit"
          disabled={isGenerating || !prompt.trim()}
          className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
        >
          {isGenerating ? (
            <>
              <Loader2 className="w-5 h-5 animate-spin" />
              <span>Generando...</span>
            </>
          ) : (
            <>
              <Sparkles className="w-5 h-5" />
              <span>Generar Path</span>
            </>
          )}
        </button>
      </form>
    </div>
  );
};

