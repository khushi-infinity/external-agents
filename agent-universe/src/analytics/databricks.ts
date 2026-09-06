import type { Repository, DatabricksAnalytics } from '../data/types';

/**
 * Databricks analytics layer.
 *
 * In a full deployment this module:
 *   1. Exports normalized AgentEvents as JSON/NDJSON (Entire checkpoints + git data)
 *   2. Sends them to a Databricks workspace (SQL Warehouse / notebook job)
 *   3. Runs the analysis: agent performance, file hotspots, failure patterns,
 *      development velocity, risk scoring
 *   4. Returns the results as a DatabricksAnalytics object for the 3D view
 *
 * The MVP ships with precomputed sample analytics so the demo never blocks.
 * The risk_map drives the "Risk View" overlay on the 3D universe.
 */
export async function fetchDatabricksAnalytics(_repo: Repository): Promise<DatabricksAnalytics> {
  // Try live endpoint if the backend is running
  try {
    const resp = await fetch('/api/analytics');
    if (resp.ok) {
      const data = await resp.json();
      if (data?.agent_performance) return data;
    }
  } catch {
    /* fall through to sample analytics */
  }
  const { sampleAnalytics } = await import('../data/sample-data');
  return sampleAnalytics;
}

/**
 * Compute risk map locally from events (offline fallback when Databricks is unavailable).
 * Weighted: change count + failed test runs + checkpoint risk.
 */
export function computeRiskMap(repo: Repository): Record<string, number> {
  const map: Record<string, number> = {};
  for (const f of repo.files) {
    let score = f.risk_score || 0;
    if (f.change_count > 10) score = Math.min(100, score + 15);
    if (f.change_count > 5) score = Math.min(100, score + 8);
    map[f.path] = score;
  }
  return map;
}

/** Export repository activity as NDJSON for Databricks ingestion. */
export function exportEventsNDJSON(repo: Repository): string {
  const lines = repo.files.map(f => JSON.stringify({
    event_type: 'file_change',
    file: f.path,
    change_count: f.change_count,
    risk_score: f.risk_score,
    last_agent: f.last_agent,
    last_checkpoint: f.last_checkpoint,
  }));
  return lines.join('\n');
}