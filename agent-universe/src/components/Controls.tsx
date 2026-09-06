import { useState } from 'react';
import { useStore } from '../store';

const AGENT_COLORS: Record<string, string> = { claude: '#d97706', codex: '#10b981', copilot: '#6366f1', aider: '#ec4899' };

function AnalyticsModal({ onClose }: { onClose: () => void }) {
  const analytics = useStore(s => s.analytics);
  const dirs = [...new Set(useStore(s => s.repository).files.map(f => f.directory))];

  const maxChanges = Math.max(...analytics.file_hotspots.map(h => h.changes), 1);

  return (
    <div style={{
      position: 'fixed', top: '50%', left: '50%', transform: 'translate(-50%, -50%)',
      width: 720, maxWidth: '92vw', maxHeight: '80vh', overflowY: 'auto',
      background: '#0f172a', border: '1px solid #1e293b', borderRadius: 12, zIndex: 200, padding: 24,
    }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <div>
          <div style={{ fontSize: '18px', fontWeight: 700, color: '#e2e8f0' }}>Databricks Analytics</div>
          <div style={{ fontSize: '12px', color: '#94a3b8' }}>Risk scoring & activity insights over development history</div>
        </div>
        <button onClick={onClose} style={{ background: '#1e293b', border: 'none', color: '#e2e8f0', borderRadius: 6, padding: '6px 12px', cursor: 'pointer' }}>Close</button>
      </div>

      <div style={{ fontSize: '13px', fontWeight: 600, color: '#e2e8f0', margin: '16px 0 8px' }}>Agent Performance</div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))', gap: 8 }}>
        {Object.entries(analytics.agent_performance).map(([name, a]) => (
          <div key={name} style={{ background: '#1e293b', borderRadius: 8, padding: 12 }}>
            <div style={{ color: AGENT_COLORS[name] || '#94a3b8', fontSize: 13, fontWeight: 700 }}>{name.toUpperCase()}</div>
            <div style={{ fontSize: 11, color: '#94a3b8', marginTop: 4 }}>{a.sessions} sessions · {a.avg_duration} min avg</div>
            <div style={{ marginTop: 6, height: 4, background: '#0f172a', borderRadius: 2, overflow: 'hidden' }}>
              <div style={{ height: '100%', width: `${a.success_rate * 100}%`, background: AGENT_COLORS[name] || '#94a3b8' }} />
            </div>
            <div style={{ fontSize: 11, color: '#e2e8f0', marginTop: 2 }}>{(a.success_rate * 100).toFixed(0)}% success</div>
          </div>
        ))}
      </div>

      <div style={{ fontSize: '13px', fontWeight: 600, color: '#e2e8f0', margin: '16px 0 8px' }}>File Hotspots</div>
      {analytics.file_hotspots.slice(0, 8).map(h => (
        <div key={h.path} style={{ marginBottom: 6 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 12 }}>
            <span style={{ color: '#94a3b8', fontFamily: 'monospace' }}>{h.path}</span>
            <span style={{ color: '#e2e8f0' }}>{h.changes} changes · risk {h.risk}</span>
          </div>
          <div style={{ marginTop: 2, height: 4, background: '#1e293b', borderRadius: 2, overflow: 'hidden' }}>
            <div style={{ height: '100%', width: `${(h.changes / maxChanges) * 100}%`, background: h.risk >= 70 ? '#ef4444' : h.risk >= 50 ? '#f59e0b' : '#22c55e' }} />
          </div>
        </div>
      ))}

      <div style={{ fontSize: '13px', fontWeight: 600, color: '#e2e8f0', margin: '16px 0 8px' }}>Failure Patterns</div>
      {analytics.failure_patterns.map(p => (
        <div key={p.module} style={{ display: 'flex', justifyContent: 'space-between', background: '#1e293b', borderRadius: 6, padding: '8px 12px', marginBottom: 4 }}>
          <span style={{ fontSize: 13, color: '#e2e8f0' }}>{p.module}</span>
          <span style={{ fontSize: 12, color: p.risk >= 70 ? '#ef4444' : '#f59e0b' }}>{p.failures} failures · risk {p.risk}</span>
        </div>
      ))}

      <div style={{ fontSize: '13px', fontWeight: 600, color: '#e2e8f0', margin: '16px 0 8px' }}>Development Velocity</div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 8 }}>
        {[
          [`${analytics.velocity.avg_task_minutes} min`, 'avg task time'],
          [`${analytics.velocity.avg_files_per_session}`, 'files / session'],
          [`${analytics.velocity.avg_retries}`, 'avg retries'],
        ].map(([v, l]) => (
          <div key={l} style={{ background: '#1e293b', borderRadius: 8, padding: 12, textAlign: 'center' }}>
            <div style={{ fontSize: 18, fontWeight: 700, color: '#e2e8f0' }}>{v}</div>
            <div style={{ fontSize: 11, color: '#94a3b8' }}>{l}</div>
          </div>
        ))}
      </div>

      <div style={{ fontSize: '13px', fontWeight: 600, color: '#e2e8f0', margin: '16px 0 8px' }}>Modules</div>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
        {dirs.map(d => (
          <span key={d} style={{ background: '#1e293b', borderRadius: 4, padding: '4px 10px', fontSize: 12, color: '#94a3b8', fontFamily: 'monospace' }}>{d}</span>
        ))}
      </div>
    </div>
  );
}

export function Controls() {
  const store = useStore();
  const [showAnalytics, setShowAnalytics] = useState(false);
  const dirs = [...new Set(store.repository.files.map(f => f.directory))];

  const toggleBtn = (active: boolean, color?: string) => ({
    background: active ? (color || '#334155') : '#0f172a',
    color: active ? '#ffffff' : '#94a3b8',
    border: active ? 'none' : '1px solid #334155',
    padding: '6px 12px',
    borderRadius: 6,
    fontSize: 12,
    cursor: 'pointer',
    fontWeight: 600,
  });

  return (
    <>
      <div style={{
        position: 'fixed', top: 0, left: 0, right: 0, zIndex: 50,
        padding: '12px 16px', display: 'flex', alignItems: 'center', gap: 12,
        background: 'linear-gradient(to bottom, rgba(2,6,23,0.9), transparent)',
      }}>
        <div style={{ fontWeight: 800, fontSize: 15, color: '#e2e8f0', letterSpacing: -0.2 }}>
          ✦ Agent Activity Universe
        </div>
        <div style={{ fontSize: 11, color: '#64748b' }}>
          {store.repository.name} · Entire Checkpoints + Graph · Databricks Analytics
        </div>

        <div style={{ flex: 1 }} />

        {/* Risk view toggle */}
        <button
          onClick={store.toggleRiskView}
          style={toggleBtn(store.riskViewEnabled, store.riskViewEnabled ? '#ef4444' : undefined)}
        >
          {store.riskViewEnabled ? '⦿ Risk View ON' : '○ Risk View'}
        </button>

        <button
          onClick={() => store.setReplayPlaying(!store.replayPlaying)}
          style={toggleBtn(store.replayPlaying, store.replayPlaying ? '#22c55e' : undefined)}
        >
          {store.replayPlaying ? '⏸ Activity Replay' : '▶ Activity Replay'}
        </button>

        <button onClick={() => setShowAnalytics(true)} style={toggleBtn(false)}>📊 Analytics</button>
      </div>

      <div style={{
        position: 'fixed', bottom: 0, left: 0, right: 0, zIndex: 50,
        padding: '10px 16px', display: 'flex', alignItems: 'center', gap: 10, flexWrap: 'wrap',
        background: 'linear-gradient(to top, rgba(2,6,23,0.9), transparent)',
      }}>
        <span style={{ fontSize: 11, color: '#64748b', textTransform: 'uppercase', letterSpacing: 1 }}>Agents</span>
        {Object.keys(AGENT_COLORS).map(a => (
          <button key={a} onClick={() => store.toggleAgent(a)} style={toggleBtn(store.activeAgents.includes(a), AGENT_COLORS[a])}>
            {a}
          </button>
        ))}
        <div style={{ width: 1, height: 20, background: '#334155', margin: '0 4px' }} />
        <span style={{ fontSize: 11, color: '#64748b', textTransform: 'uppercase', letterSpacing: 1 }}>Modules</span>
        {dirs.map(d => (
          <button key={d} onClick={() => store.toggleDirectory(d)} style={toggleBtn(store.activeDirectories.includes(d))}>
            {d}
          </button>
        ))}
      </div>

      {showAnalytics && <AnalyticsModal onClose={() => setShowAnalytics(false)} />}
    </>
  );
}