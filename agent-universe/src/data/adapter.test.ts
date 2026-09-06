import { describe, it, expect } from 'vitest';
import { importEntireExport, buildTimeline } from './adapter';
import { sampleRepository } from './sample-data';

describe('importEntireExport', () => {
  it('parses a valid Entire export into a Repository', () => {
    const repo = importEntireExport(JSON.stringify({
      name: 'demo',
      files: [{ path: 'src/auth/service.ts', risk_score: 82, change_count: 17, last_agent: 'claude', connections: ['src/auth/routes.ts'] }],
      sessions: [{ id: 's1', agent: 'claude', start_time: '2026-09-06T09:00:00Z', prompt_count: 2 }],
      checkpoints: [{ id: 'CP-1', commit: 'abc', agent: 'claude', prompt: 'fix', intent: 'fix auth' }],
    }));
    expect(repo.name).toBe('demo');
    expect(repo.files).toHaveLength(1);
    expect(repo.files[0].name).toBe('service.ts');
    expect(repo.files[0].directory).toBe('src/auth');
    expect(repo.files[0].risk_score).toBe(82);
    expect(repo.files[0].connections).toEqual(['src/auth/routes.ts']);
    expect(repo.sessions).toHaveLength(1);
    expect(repo.checkpoints).toHaveLength(1);
  });

  it('handles empty / malformed entries gracefully', () => {
    const repo = importEntireExport(JSON.stringify({ files: [{ path: '' }] }));
    expect(repo.files).toHaveLength(1);
    expect(repo.files[0].name).toBeTruthy();
    expect(repo.files[0].directory).toBe('root');
  });
});

describe('buildTimeline', () => {
  it('produces an ordered timeline of replay steps', () => {
    const steps = buildTimeline(sampleRepository);
    expect(steps.length).toBeGreaterThan(0);
    // chronological order
    for (let i = 1; i < steps.length; i++) {
      expect(steps[i].time.getTime()).toBeGreaterThanOrEqual(steps[i - 1].time.getTime());
    }
    // every session start is represented
    const sessionStarts = steps.filter(s => s.kind === 'session_start');
    expect(sessionStarts.length).toBe(sampleRepository.sessions.length);
    // checkpoints appear with intent text
    const checkpoints = steps.filter(s => s.kind === 'checkpoint');
    expect(checkpoints.length).toBeGreaterThanOrEqual(sampleRepository.checkpoints.length);
  });

  it('sorts sessions by start time regardless of input order', () => {
    const shuffled = {
      ...sampleRepository,
      sessions: [...sampleRepository.sessions].reverse(),
    };
    const steps = buildTimeline(shuffled);
    const starts = steps.filter(s => s.kind === 'session_start');
    expect(starts[0].time.getTime()).toBeLessThanOrEqual(starts[starts.length - 1].time.getTime());
  });
});