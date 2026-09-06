import { create } from 'zustand';
import type { Repository, SelectedNode, AgentName, DatabricksAnalytics } from './data/types';
import { sampleRepository, sampleAnalytics } from './data/sample-data';

interface AppState {
  repository: Repository;
  analytics: DatabricksAnalytics;
  selectedNode: SelectedNode | null;
  sidePanelOpen: boolean;
  riskViewEnabled: boolean;
  activeAgents: AgentName[];
  activeDirectories: string[];
  // replay
  replayPlaying: boolean;
  replayStep: number;       // index into the timeline
  replayActiveFile: string | null;
  setSelectedNode: (node: SelectedNode | null) => void;
  toggleRiskView: () => void;
  toggleAgent: (agent: AgentName) => void;
  toggleDirectory: (dir: string) => void;
  setReplayPlaying: (v: boolean) => void;
  setReplayStep: (i: number | ((prev: number) => number)) => void;
  setReplayActiveFile: (id: string | null) => void;
}

const allDirs = [...new Set(sampleRepository.files.map(f => f.directory))];

export const useStore = create<AppState>((set) => ({
  repository: sampleRepository,
  analytics: sampleAnalytics,
  selectedNode: null,
  sidePanelOpen: false,
  riskViewEnabled: false,
  activeAgents: ['claude', 'codex', 'copilot', 'aider'],
  activeDirectories: [...allDirs],
  replayPlaying: false,
  replayStep: 0,
  replayActiveFile: null,
  setSelectedNode: (node) => set({ selectedNode: node, sidePanelOpen: node !== null }),
  toggleRiskView: () => set((s) => ({ riskViewEnabled: !s.riskViewEnabled })),
  toggleAgent: (agent) => set((s) => ({
    activeAgents: s.activeAgents.includes(agent)
      ? s.activeAgents.filter(a => a !== agent)
      : [...s.activeAgents, agent],
  })),
  toggleDirectory: (dir) => set((s) => ({
    activeDirectories: s.activeDirectories.includes(dir)
      ? s.activeDirectories.filter(d => d !== dir)
      : [...s.activeDirectories, dir],
  })),
  setReplayPlaying: (v) => set({ replayPlaying: v }),
  setReplayStep: (i) => set((s) => ({ replayStep: typeof i === 'function' ? (i as (p: number) => number)(s.replayStep) : i })),
  setReplayActiveFile: (id) => set({ replayActiveFile: id }),
}));