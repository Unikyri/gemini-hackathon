import axios, { type AxiosInstance } from 'axios';

// Configuración del cliente API
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api';

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000,
});

// Tipos para las respuestas de la API
export interface GeneratePathRequest {
  prompt: string;
}

export interface GeneratePathResponse {
  path_id: string;
  nodes: Node[];
  message?: string;
}

export interface GetPathResponse {
  path_id: string;
  prompt: string;
  nodes: Node[];
  created_at: string;
}

export interface Node {
  node_id: string;
  title: string;
  description: string;
  content: string;
  order: number;
  completed: boolean;
  path_id: string;
}

export interface GetNodeResponse extends Node {}

// Funciones del cliente API
export const apiService = {
  /**
   * Genera un nuevo path de aprendizaje basado en un prompt
   */
  generatePath: async (prompt: string): Promise<GeneratePathResponse> => {
    const response = await apiClient.post<GeneratePathResponse>('/generate-path', {
      prompt,
    });
    return response.data;
  },

  /**
   * Obtiene un path de aprendizaje existente por su ID
   */
  getPath: async (pathId: string): Promise<GetPathResponse> => {
    const response = await apiClient.get<GetPathResponse>(`/paths/${pathId}`);
    return response.data;
  },

  /**
   * Obtiene un nodo específico de un path de aprendizaje
   */
  getNode: async (pathId: string, nodeId: string): Promise<GetNodeResponse> => {
    const response = await apiClient.get<GetNodeResponse>(
      `/paths/${pathId}/nodes/${nodeId}`
    );
    return response.data;
  },

  /**
   * Actualiza el estado de completado de un nodo
   */
  updateNodeCompletion: async (
    pathId: string,
    nodeId: string,
    completed: boolean
  ): Promise<GetNodeResponse> => {
    const response = await apiClient.patch<GetNodeResponse>(
      `/paths/${pathId}/nodes/${nodeId}`,
      { completed }
    );
    return response.data;
  },
};

export default apiClient;

