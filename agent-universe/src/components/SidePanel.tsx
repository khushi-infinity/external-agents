import { useStore } from '../store';
import type { FileNode, AgentSession, Checkpoint } from '../data/types';

function riskBadge(score: number) {
  if (score >= 70) return { label: 'HIGH', color: '#ef4444', bg: '#451a1a' };
  if (score >= 50) return { label: 'MEDIUM', color: '#f59e0b', bg: '#452a0a' };
  if (score >= 30) return { label: 'LOW', color: '#eab308', bg: '#453a0a' };
  return { label: 'SAFE', color: '#22c55e', bg: '#0a4520' };
}

const stat = (label: string, value: string) => (
  <div key={label} style={{ background: '#1e293b', padding: '8px 10px', borderRadius: 6 }}>
    <div style={{ fontSize: '10px', color: '#94a3b8' }}>{label}</div>
    <div style={{ fontSize: '15px', fontWeight: 600, color: '#e2e8f0' }}>{value}</div>
  </div>
);

function FileDetails({ file }: { file: FileNode }) {
  const badge = riskBadge(file.risk_score);
  return (
    <div style={{ padding: '16px' }}>
      <div style={{ fontSize: '11px', color: '#94a3b8', textTransform: 'uppercase', letterSpacing: '1px', marginBottom: 8 }}>File</div>
      <div style={{ fontSize: '16px', fontWeight: 600, color: '#e2e8f0', marginBottom: 4 }}>{file.name}</div>
      <div style={{ fontSize: '12px', color: '#64748b', marginBottom: 16, fontFamily: 'monospace' }}>{file.path}</div>
      <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
        <span style={{ padding: '4px 10px', borderRadius: 4, fontSize: '11px', fontWeight: 700, color: badge.color, background: badge.bg }}>
          Risk: {badge.label} ({file.risk_score})
        </span>
        <span style={{ padding: '4px 10px', borderRadius: 4, fontSize: '11px', color: '#94a3b8', background: '#1e293b' }}>
          {file.change_count} changes
        </span>
      </div>
      {file.last_agent && (
        <div style={{ marginBottom: 12 }}>
          <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Last Agent</div>
          <div style={{ fontSize: '13px', color: '#e2e8f0' }}>{file.last_agent.toUpperCase()}</div>
        </div>
      )}
      {file.last_checkpoint && (
        <div style={{ marginBottom: 12 }}>
          <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Checkpoint</div>
          <div style={{ fontSize: '13px', color: '#60a5fa', fontFamily: 'monospace' }}>{file.last_checkpoint}</div>
        </div>
      )}
      {file.last_prompt && (
        <div style={{ marginBottom: 12 }}>
          <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Agent Prompt</div>
          <div style={{ fontSize: '13px', color: '#e2e8f0', fontStyle: 'italic', lineHeight: 1.5 }}>"{file.last_prompt}"</div>
        </div>
      )}
      <div>
        <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Connected Files</div>
        {file.connections.map(c => (
          <div key={c} style={{ fontSize: '12px', color: '#94a3b8', fontFamily: 'monospace', padding: '2px 0' }}>→ {c}</div>
        ))}
      </div>
    </div>
  );
}

function SessionDetails({ session }: { session: AgentSession }) {
  const badge = riskBadge(Math.round((1 - session.success_rate) * 100));
  return (
    <div style={{ padding: '16px' }}>
      <div style={{ fontSize: '11px', color: '#94a3b8', textTransform: 'uppercase', letterSpacing: '1px', marginBottom: 8 }}>Agent Session</div>
      <div style={{ fontSize: '16px', fontWeight: 600, color: '#e2e8f0', marginBottom: 12 }}>{session.agent.toUpperCase()}</div>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginBottom: 16 }}>
        {stat('Prompts', String(session.prompt_count))}
        {stat('File Changes', String(session.file_changes))}
        {stat('Tool Calls', String(session.tool_calls))}
        {stat('Tokens', `${(session.tokens_used / 1000).toFixed(0)}k`)}
      </div>
      <div style={{ marginBottom: 12 }}>
        <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Success Rate</div>
        <div style={{ height: 6, background: '#1e293b', borderRadius: 3, overflow: 'hidden' }}>
          <div style={{ height: '100%', width: `${session.success_rate * 100}%`, background: badge.color, borderRadius: 3 }} />
        </div>
        <div style={{ fontSize: '12px', color: badge.color, marginTop: 4 }}>{(session.success_rate * 100).toFixed(0)}%</div>
      </div>
      <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Files Touched</div>
      {session.files_touched.map(f => (
        <div key={f} style={{ fontSize: '12px', color: '#94a3b8', fontFamily: 'monospace', padding: '2px 0' }}>• {f}</div>
      ))}
      {session.checkpoints.length > 0 && (
        <>
          <div style={{ fontSize: '11px', color: '#94a3b8', marginTop: 12, marginBottom: 4 }}>Checkpoints</div>
          {session.checkpoints.map(c => (
            <div key={c} style={{ fontSize: '12px', color: '#60a5fa', fontFamily: 'monospace', padding: '2px 0' }}>{c}</div>
          ))}
        </>
      )}
    </div>
  );
}

function CheckpointDetails({ cp }: { cp: Checkpoint }) {
  return (
    <div style={{ padding: '16px' }}>
      <div style={{ fontSize: '11px', color: '#94a3b8', textTransform: 'uppercase', letterSpacing: '1px', marginBottom: 8 }}>Checkpoint</div>
      <div style={{ fontSize: '18px', fontWeight: 700, color: '#60a5fa', fontFamily: 'monospace', marginBottom: 4 }}>{cp.id}</div>
      <div style={{ fontSize: '12px', color: '#64748b', marginBottom: 16, fontFamily: 'monospace' }}>commit {cp.commit}</div>
      <div style={{ marginBottom: 12 }}>
        <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Agent</div>
        <div style={{ fontSize: '14px', color: '#e2e8f0' }}>{cp.agent.toUpperCase()}</div>
      </div>
      <div style={{ marginBottom: 12 }}>
        <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Intent</div>
        <div style={{ fontSize: '13px', color: '#e2e8f0', fontStyle: 'italic', lineHeight: 1.5 }}>{cp.intent}</div>
      </div>
      <div style={{ marginBottom: 12 }}>
        <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Prompt</div>
        <div style={{ fontSize: '13px', color: '#e2e8f0', lineHeight: 1.5 }}>"{cp.prompt}"</div>
      </div>
      <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
        <span style={{ padding: '4px 10px', borderRadius: 4, fontSize: '11px', color: '#94a3b8', background: '#1e293b' }}>{cp.tool_calls} tool calls</span>
        <span style={{ padding: '4px 10px', borderRadius: 4, fontSize: '11px', color: '#94a3b8', background: '#1e293b' }}>{(cp.tokens / 1000).toFixed(0)}k tokens</span>
      </div>
      <div style={{ fontSize: '11px', color: '#94a3b8', marginBottom: 4 }}>Files Changed</div>
      {cp.files_changed.map(f => (
        <div key={f} style={{ fontSize: '12px', color: '#94a3b8', fontFamily: 'monospace', padding: '2px 0' }}>• {f}</div>
      ))}
    </div>
  );
}

export function SidePanel() {
  const selected = useStore(s => s.selectedNode);
  const open = useStore(s => s.sidePanelOpen);
  const setSelected = useStore(s => s.setSelectedNode);

  if (!open || !selected) return null;

  return (
    <div style={{
      position: 'fixed', top: 0, right: 0, width: 340, height: '100vh',
      background: 'rgba(15, 23, 42, 0.96)', borderLeft: '1px solid #1e293b',
      backdropFilter: 'blur(12px)', zIndex: 100, overflowY: 'auto',
    }}>
      <div style={{ padding: '12px 16px', borderBottom: '1px solid #1e293b', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <span style={{ fontSize: '12px', color: '#94a3b8', textTransform: 'uppercase', letterSpacing: '1px' }}>Inspector</span>
        <button onClick={() => setSelected(null)} style={{ background: 'none', border: 'none', color: '#94a3b8', cursor: 'pointer', fontSize: '18px' }}>✕</button>
      </div>
      {selected.type === 'file' && <FileDetails file={selected.data as FileNode} />}
      {selected.type === 'session' && <SessionDetails session={selected.data as AgentSession} />}
      {selected.type === 'checkpoint' && <CheckpointDetails cp={selected.data as Checkpoint} />}
    </div>
  );
}