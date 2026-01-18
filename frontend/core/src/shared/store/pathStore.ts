import { create } from 'zustand';
import { type Node } from '../api';

interface PathState {
  currentPathId: string | null;
  currentNodes: Node[];
  currentNodeId: string | null;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  setCurrentPath: (pathId: string, nodes: Node[]) => void;
  setCurrentNode: (nodeId: string) => void;
  updateNodeInList: (nodeId: string, updates: Partial<Node>) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  reset: () => void;
}

export const usePathStore = create<PathState>((set) => ({
  currentPathId: null,
  currentNodes: [],
  currentNodeId: null,
  isLoading: false,
  error: null,

  setCurrentPath: (pathId, nodes) =>
    set({ currentPathId: pathId, currentNodes: nodes, error: null }),

  setCurrentNode: (nodeId) =>
    set({ currentNodeId: nodeId }),

  updateNodeInList: (nodeId, updates) =>
    set((state) => ({
      currentNodes: state.currentNodes.map((node) =>
        node.node_id === nodeId ? { ...node, ...updates } : node
      ),
    })),

  setLoading: (loading) =>
    set({ isLoading: loading }),

  setError: (error) =>
    set({ error, isLoading: false }),

  reset: () =>
    set({
      currentPathId: null,
      currentNodes: [],
      currentNodeId: null,
      isLoading: false,
      error: null,
    }),
}));

