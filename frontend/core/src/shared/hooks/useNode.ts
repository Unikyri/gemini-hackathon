import { useState, useCallback } from 'react';
import { apiService, type GetNodeResponse } from '../api';
import { usePathStore } from '../store';

export const useNode = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { currentPathId, setCurrentNode, updateNodeInList } = usePathStore();

  const fetchNode = useCallback(async (
    pathId: string,
    nodeId: string
  ): Promise<GetNodeResponse | null> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await apiService.getNode(pathId, nodeId);
      setCurrentNode(nodeId);
      return response;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error cargando el nodo';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
    }
  }, [setCurrentNode]);

  const updateNodeCompletion = useCallback(async (
    nodeId: string,
    completed: boolean
  ): Promise<boolean> => {
    if (!currentPathId) {
      setError('No hay un path activo');
      return false;
    }

    try {
      const response = await apiService.updateNodeCompletion(currentPathId, nodeId, completed);
      updateNodeInList(nodeId, { completed: response.completed });
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error actualizando el nodo';
      setError(errorMessage);
      return false;
    }
  }, [currentPathId, updateNodeInList]);

  return { fetchNode, updateNodeCompletion, isLoading, error };
};

