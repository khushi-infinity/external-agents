import type { Repository } from './types';
import { sampleRepository } from './sample-data';

/**
 * Data loading priority:
 *   1. /data/entire-export.json  — drop a real Entire export here (public/data/)
 *      Schema: { name, files: [{ path, risk_score, change_count, last_agent,
 *                last_checkpoint, last_prompt, connections }], sessions, checkpoints }
 *   2. /api/repository           — live backend (Databricks/Entire server)
 *   3. sample data               — demo fallback, always works offline
 */
export async function loadRepositoryData(): Promise<Repository> {
  // Route 1: drop-in Entire export (place in public/data/entire-export.json)
  try {
    const resp = await fetch('/data/entire-export.json');
    if (resp.ok) {
      const data = await resp.json();
      if (data?.files?.length) {
        console.info('[adapter] Loaded real Entire export from /data/entire-export.json');
        return importEntireExport(JSON.stringify(data));
      }
    }
  } catch { /* not present */ }

  // Route 2: live backend
  try {
    const resp = await fetch('/api/repository');
    if (resp.ok) {
      const data = await resp.json();
      if (data?.files?.length) return data;
    }
  } catch { /* no server */ }

  // Route 3: sample data
  console.info('[adapter] Using sample dataset for demo');
  return sampleRepository;
}

/** Import Entire data from a JSON export (dashboard export, CLI output, or drop-in file). */
export function importEntireExport(json: string): Repository {
  const raw = JSON.parse(json);
  const files = (raw.files || []).map((f: Record<string, unknown>, i: number) => ({
    id: `f${i}`,
    path: (f.path as string) || `file-${i}`,
    name: ((f.path as string) || '').split('/').pop() || `file-${i}`,
    directory: ((f.path as string) || '').split('/').slice(0, -1).join('/') || 'root',
    extension: ((f.path as string)?.match(/\.[^.]+$/) || [''])[0],
    risk_score: (f.risk_score as number) || 0,
    change_count: (f.change_count as number) || 0,
    last_agent: f.last_agent as string,
    last_checkpoint: f.last_checkpoint as string,
    last_prompt: f.last_prompt as string,
    position: [(Math.random() - 0.5) * 10, (Math.random() - 0.5) * 8, (Math.random() - 0.5) * 6] as [number, number, number],
    connections: (f.connections as string[]) || [],
  }));
  return {
    name: (raw.name as string) || 'imported-repo',
    files,
    sessions: raw.sessions || [],
    checkpoints: raw.checkpoints || [],
  };
}

/**
 * Build a synthetic timeline of events from the repository for Activity Replay.
 * In production this reads real AgentEvents from Entire; here it derives an
 * ordered sequence from sessions, checkpoints and file change counts.
 */
export function buildTimeline(repo: Repository): ReplayStep[] {
  const steps: ReplayStep[] = [];
  const sortedSessions = [...repo.sessions].sort((a, b) =>
    new Date(a.start_time).getTime() - new Date(b.start_time).getTime()
  );

  for (const s of sortedSessions) {
    steps.push({
      time: new Date(s.start_time),
      label: `${s.agent.toUpperCase()} started session`,
      kind: 'session_start',
      agent: s.agent,
    });
    // Interleave file touches (stable order via id sort for determinism)
    const touched = [...s.files_touched].sort();
    touched.forEach((fid, i) => {
      const file = repo.files.find(f => f.id === fid);
      steps.push({
        time: new Date(new Date(s.start_time).getTime() + (i + 1) * 4 * 60000),
        label: file ? `Modified ${file.path}` : `Touched ${fid}`,
        kind: 'file_change',
        agent: s.agent,
        fileId: fid,
      });
    });
    for (const cp of repo.checkpoints.filter(c => c.session_id === s.id)) {
      steps.push({
        time: new Date(cp.timestamp),
        label: `Checkpoint ${cp.id} · "${cp.intent.slice(0, 48)}"`,
        kind: 'checkpoint',
        agent: s.agent,
        fileId: cp.files_changed.length ? repo.files.find(f => cp.files_changed.includes(f.path))?.id : undefined,
      });
    }
  }
  return steps.sort((a, b) => a.time.getTime() - b.time.getTime());
}

export interface ReplayStep {
  time: Date;
  label: string;
  kind: 'session_start' | 'file_change' | 'checkpoint';
  agent?: string;
  fileId?: string;
}