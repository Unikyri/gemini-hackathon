import { useState, useCallback } from 'react';
import { apiService, type GeneratePathResponse } from '../api';
import { usePathStore } from '../store';

export const useGeneratePath = () => {
  const [isGenerating, setIsGenerating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { setCurrentPath, setLoading } = usePathStore();

  const generatePath = useCallback(async (prompt: string): Promise<GeneratePathResponse | null> => {
    setIsGenerating(true);
    setLoading(true);
    setError(null);

    try {
      const response = await apiService.generatePath(prompt);
      setCurrentPath(response.path_id, response.nodes);
      return response;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error generando el path';
      setError(errorMessage);
      return null;
    } finally {
      setIsGenerating(false);
      setLoading(false);
    }
  }, [setCurrentPath, setLoading]);

  return { generatePath, isGenerating, error };
};

