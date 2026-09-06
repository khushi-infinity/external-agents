import type { Repository, FileNode, AgentSession, Checkpoint, DatabricksAnalytics } from './types';

const files: FileNode[] = [
  { id: 'f1', path: 'src/auth/service.ts', name: 'service.ts', directory: 'src/auth', extension: '.ts', risk_score: 82, change_count: 17, last_agent: 'claude', last_checkpoint: 'CP-1842', last_prompt: 'Implement OAuth authentication flow with refresh tokens', position: [-4, 2, 0], connections: ['f2', 'f3', 'f7'] },
  { id: 'f2', path: 'src/auth/middleware.ts', name: 'middleware.ts', directory: 'src/auth', extension: '.ts', risk_score: 45, change_count: 8, last_agent: 'claude', position: [-5, 0.5, 1], connections: ['f1', 'f3'] },
  { id: 'f3', path: 'src/auth/routes.ts', name: 'routes.ts', directory: 'src/auth', extension: '.ts', risk_score: 30, change_count: 12, last_agent: 'codex', position: [-3, 0.5, -1], connections: ['f1', 'f2', 'f4'] },
  { id: 'f4', path: 'src/payments/service.ts', name: 'service.ts', directory: 'src/payments', extension: '.ts', risk_score: 71, change_count: 14, last_agent: 'codex', last_checkpoint: 'CP-1835', last_prompt: 'Add Stripe webhook handling for payment intents', position: [3, 2, 0], connections: ['f5', 'f7'] },
  { id: 'f5', path: 'src/payments/api.ts', name: 'api.ts', directory: 'src/payments', extension: '.ts', risk_score: 55, change_count: 9, last_agent: 'codex', position: [4, 0.5, 1], connections: ['f4', 'f6'] },
  { id: 'f6', path: 'src/database/client.ts', name: 'client.ts', directory: 'src/database', extension: '.ts', risk_score: 60, change_count: 9, last_agent: 'aider', last_checkpoint: 'CP-1830', position: [0, -2, 0], connections: ['f7'] },
  { id: 'f7', path: 'src/database/models.ts', name: 'models.ts', directory: 'src/database', extension: '.ts', risk_score: 35, change_count: 6, last_agent: 'claude', position: [1, -3, 1], connections: ['f6'] },
  { id: 'f8', path: 'src/api/routes.ts', name: 'routes.ts', directory: 'src/api', extension: '.ts', risk_score: 40, change_count: 11, last_agent: 'copilot', position: [0, 3, -1], connections: ['f1', 'f4', 'f9'] },
  { id: 'f9', path: 'src/api/validators.ts', name: 'validators.ts', directory: 'src/api', extension: '.ts', risk_score: 20, change_count: 4, last_agent: 'copilot', position: [-1, 4, 0], connections: ['f8'] },
  { id: 'f10', path: 'tests/auth.test.ts', name: 'auth.test.ts', directory: 'tests', extension: '.ts', risk_score: 15, change_count: 7, last_agent: 'claude', position: [-6, 3, 2], connections: ['f1', 'f2'] },
  { id: 'f11', path: 'tests/payments.test.ts', name: 'payments.test.ts', directory: 'tests', extension: '.ts', risk_score: 25, change_count: 5, last_agent: 'codex', position: [6, 3, 2], connections: ['f4', 'f5'] },
  { id: 'f12', path: 'src/config.ts', name: 'config.ts', directory: 'src', extension: '.ts', risk_score: 10, change_count: 3, last_agent: 'copilot', position: [0, 5, 0], connections: ['f8', 'f6'] },
];

const sessions: AgentSession[] = [
  { id: 'sess_182', agent: 'claude', start_time: '2026-09-06T09:15:00Z', end_time: '2026-09-06T10:42:00Z', prompt_count: 14, file_changes: 8, tool_calls: 32, checkpoints: ['CP-1842', 'CP-1841'], success_rate: 0.91, tokens_used: 142000, files_touched: ['f1', 'f2', 'f3', 'f10'] },
  { id: 'sess_179', agent: 'codex', start_time: '2026-09-06T09:30:00Z', end_time: '2026-09-06T10:15:00Z', prompt_count: 8, file_changes: 5, tool_calls: 18, checkpoints: ['CP-1839', 'CP-1835'], success_rate: 0.84, tokens_used: 89000, files_touched: ['f4', 'f5', 'f3', 'f11'] },
  { id: 'sess_175', agent: 'copilot', start_time: '2026-09-06T09:45:00Z', end_time: '2026-09-06T10:05:00Z', prompt_count: 5, file_changes: 3, tool_calls: 10, checkpoints: [], success_rate: 0.95, tokens_used: 34000, files_touched: ['f8', 'f9', 'f12'] },
  { id: 'sess_171', agent: 'aider', start_time: '2026-09-06T10:00:00Z', end_time: '2026-09-06T10:30:00Z', prompt_count: 6, file_changes: 2, tool_calls: 12, checkpoints: ['CP-1830'], success_rate: 0.76, tokens_used: 52000, files_touched: ['f6'] },
];

const checkpoints: Checkpoint[] = [
  { id: 'CP-1842', commit: 'a83f91c', session_id: 'sess_182', agent: 'claude', timestamp: '2026-09-06T10:42:00Z', prompt: 'Implement OAuth authentication flow with refresh tokens', files_changed: ['src/auth/service.ts', 'src/auth/middleware.ts', 'src/auth/routes.ts', 'tests/auth.test.ts'], tool_calls: 14, tokens: 42000, intent: 'Add secure OAuth2 authentication with JWT refresh token rotation', risk_score: 82 },
  { id: 'CP-1841', commit: 'b2e7d4a', session_id: 'sess_182', agent: 'claude', timestamp: '2026-09-06T10:15:00Z', prompt: 'Fix auth middleware token validation', files_changed: ['src/auth/middleware.ts'], tool_calls: 6, tokens: 18000, intent: 'Correct JWT expiry check logic', risk_score: 45 },
  { id: 'CP-1839', commit: 'c91f3b8', session_id: 'sess_179', agent: 'codex', timestamp: '2026-09-06T10:00:00Z', prompt: 'Add Stripe webhook handling for payment intents', files_changed: ['src/payments/service.ts', 'src/payments/api.ts'], tool_calls: 10, tokens: 35000, intent: 'Handle Stripe payment_intent.succeeded and payment_intent.failed webhooks', risk_score: 71 },
  { id: 'CP-1835', commit: 'd45a1e2', session_id: 'sess_179', agent: 'codex', timestamp: '2026-09-06T09:45:00Z', prompt: 'Create payment service interface', files_changed: ['src/payments/service.ts'], tool_calls: 5, tokens: 12000, intent: 'Define PaymentService type and Stripe implementation', risk_score: 55 },
  { id: 'CP-1830', commit: 'e78c2f1', session_id: 'sess_171', agent: 'aider', timestamp: '2026-09-06T10:30:00Z', prompt: 'Update database client connection pooling', files_changed: ['src/database/client.ts'], tool_calls: 8, tokens: 22000, intent: 'Switch from single connection to pool with max 10 connections', risk_score: 60 },
];

export const sampleAnalytics: DatabricksAnalytics = {
  agent_performance: {
    claude: { sessions: 42, success_rate: 0.91, avg_duration: 23 },
    codex: { sessions: 27, success_rate: 0.84, avg_duration: 19 },
    aider: { sessions: 12, success_rate: 0.76, avg_duration: 31 },
    copilot: { sessions: 18, success_rate: 0.95, avg_duration: 12 },
  },
  file_hotspots: [...files].sort((a, b) => b.change_count - a.change_count).map(f => ({ path: f.path, changes: f.change_count, risk: f.risk_score })),
  failure_patterns: [
    { module: 'payments', failures: 7, risk: 71 },
    { module: 'auth', failures: 4, risk: 82 },
    { module: 'database', failures: 3, risk: 60 },
  ],
  velocity: { avg_task_minutes: 23, avg_files_per_session: 8.4, avg_retries: 1.8 },
  risk_map: Object.fromEntries(files.map(f => [f.path, f.risk_score])),
};

export const sampleRepository: Repository = {
  name: 'my-app',
  files,
  sessions,
  checkpoints,
};