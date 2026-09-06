import { describe, it, expect } from 'vitest';
import { computeRiskMap, exportEventsNDJSON } from './databricks';
import { sampleRepository } from '../data/sample-data';

describe('computeRiskMap', () => {
  it('computes a risk score for every file, clamped to 100', () => {
    const map = computeRiskMap(sampleRepository);
    expect(Object.keys(map).length).toBe(sampleRepository.files.length);
    for (const score of Object.values(map)) {
      expect(score).toBeGreaterThanOrEqual(0);
      expect(score).toBeLessThanOrEqual(100);
    }
  });

  it('boosts risk for files changed frequently', () => {
    const heavy = { ...sampleRepository.files[0], path: 'src/a-heavy.ts', change_count: 50, risk_score: 10 };
    const light = { ...sampleRepository.files[0], path: 'src/b-light.ts', change_count: 1, risk_score: 10 };
    const repo = { ...sampleRepository, files: [heavy, light] };
    const map = computeRiskMap(repo);
    expect(map[heavy.path]).toBeGreaterThan(map[light.path]);
  });
});

describe('exportEventsNDJSON', () => {
  it('exports one JSON line per file (Databricks ingestion format)', () => {
    const ndjson = exportEventsNDJSON(sampleRepository);
    const lines = ndjson.trim().split('\n');
    expect(lines.length).toBe(sampleRepository.files.length);
    // every line is valid JSON with required fields
    for (const line of lines) {
      const event = JSON.parse(line);
      expect(event.event_type).toBe('file_change');
      expect(event.file).toBeTruthy();
      expect(typeof event.risk_score).toBe('number');
    }
  });
});