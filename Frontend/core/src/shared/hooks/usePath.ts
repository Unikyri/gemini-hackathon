import { useState, useCallback } from 'react';
import { apiService, type GetPathResponse } from '../api';
import { usePathStore } from '../store';

export const usePath = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { setCurrentPath, setLoading } = usePathStore();

  const fetchPath = useCallback(async (pathId: string): Promise<GetPathResponse | null> => {
    setIsLoading(true);
    setLoading(true);
    setError(null);

    try {
      const response = await apiService.getPath(pathId);
      setCurrentPath(response.path_id, response.nodes);
      return response;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error cargando el path';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
      setLoading(false);
    }
  }, [setCurrentPath, setLoading]);

  return { fetchPath, isLoading, error };
};

