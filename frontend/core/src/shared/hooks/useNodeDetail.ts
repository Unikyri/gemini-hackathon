import { useState, useEffect, useCallback } from 'react';
import { apiService, type GetNodeResponse } from '../api';

interface UseNodeDetailOptions {
  pathId: string | null;
  nodeId: string | null;
}

interface UseNodeDetailReturn {
  node: GetNodeResponse | null;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

/**
 * Hook para cargar los detalles de un nodo específico
 * Fetch de datos del nodo al montar componente
 * Manejar estados de loading y error
 */
export const useNodeDetail = ({ pathId, nodeId }: UseNodeDetailOptions): UseNodeDetailReturn => {
  const [node, setNode] = useState<GetNodeResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchNodeDetail = useCallback(async () => {
    if (!pathId || !nodeId) {
      setNode(null);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const response = await apiService.getNode(pathId, nodeId);
      setNode(response);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error cargando los detalles del nodo';
      setError(errorMessage);
      setNode(null);
    } finally {
      setIsLoading(false);
    }
  }, [pathId, nodeId]);

  useEffect(() => {
    fetchNodeDetail();
  }, [fetchNodeDetail]);

  return {
    node,
    isLoading,
    error,
    refetch: fetchNodeDetail,
  };
};

