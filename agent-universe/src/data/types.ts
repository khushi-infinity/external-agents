export type AgentName = 'claude' | 'codex' | 'copilot' | 'aider' | string;

export interface FileNode {
  id: string;
  path: string;
  name: string;
  directory: string;
  extension: string;
  risk_score: number;
  change_count: number;
  last_agent?: AgentName;
  last_checkpoint?: string;
  last_prompt?: string;
  position: [number, number, number];
  connections: string[];
}

export interface AgentSession {
  id: string;
  agent: AgentName;
  start_time: string;
  end_time?: string;
  prompt_count: number;
  file_changes: number;
  tool_calls: number;
  checkpoints: string[];
  success_rate: number;
  tokens_used: number;
  files_touched: string[];
}

export interface Checkpoint {
  id: string;
  commit: string;
  session_id: string;
  agent: AgentName;
  timestamp: string;
  prompt: string;
  files_changed: string[];
  tool_calls: number;
  tokens: number;
  intent: string;
  risk_score: number;
}

export interface Repository {
  name: string;
  files: FileNode[];
  sessions: AgentSession[];
  checkpoints: Checkpoint[];
}

export interface DatabricksAnalytics {
  agent_performance: Record<string, { sessions: number; success_rate: number; avg_duration: number }>;
  file_hotspots: { path: string; changes: number; risk: number }[];
  failure_patterns: { module: string; failures: number; risk: number }[];
  velocity: { avg_task_minutes: number; avg_files_per_session: number; avg_retries: number };
  risk_map: Record<string, number>;
}

export interface SelectedNode {
  type: 'file' | 'session' | 'checkpoint';
  data: FileNode | AgentSession | Checkpoint;
}