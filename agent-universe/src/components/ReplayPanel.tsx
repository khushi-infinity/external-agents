import { useMemo, useEffect, useRef } from 'react';
import { useStore } from '../store';
import { buildTimeline, type ReplayStep } from '../data/adapter';

const KIND_COLOR: Record<ReplayStep['kind'], string> = {
  session_start: '#818cf8',
  file_change: '#f59e0b',
  checkpoint: '#22c55e',
};

const KIND_ICON: Record<ReplayStep['kind'], string> = {
  session_start: '▶',
  file_change: '✎',
  checkpoint: '⛁',
};

export function ReplayPanel() {
  const repo = useStore(s => s.repository);
  const playing = useStore(s => s.replayPlaying);
  const step = useStore(s => s.replayStep);
  const setStep = useStore(s => s.setReplayStep);
  const setPlaying = useStore(s => s.setReplayPlaying);
  const setActiveFile = useStore(s => s.setReplayActiveFile);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const steps = useMemo(() => buildTimeline(repo), [repo]);

  // auto-advance while playing
  useEffect(() => {
    if (playing) {
      timerRef.current = setInterval(() => {
        setStep((s: number) => {
          if (s >= steps.length - 1) {
            setPlaying(false);
            return s;
          }
          return s + 1;
        });
      }, 900);
    }
    return () => { if (timerRef.current) clearInterval(timerRef.current); };
  }, [playing, steps.length, setStep, setPlaying]);

  // highlight active file in 3D
  useEffect(() => {
    const current = steps[step];
    setActiveFile(current?.kind === 'file_change' ? current.fileId || null : null);
  }, [step, steps, setActiveFile]);

  const current = steps[step];

  return (
    <div style={{
      position: 'fixed', bottom: 56, left: '50%', transform: 'translateX(-50%)',
      zIndex: 60, width: 560, maxWidth: '92vw',
      background: 'rgba(15, 23, 42, 0.92)', border: '1px solid #334155', borderRadius: 10,
      padding: '10px 14px', backdropFilter: 'blur(10px)',
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
        <button
          onClick={() => setPlaying(!playing)}
          style={{
            width: 34, height: 34, borderRadius: '50%', border: 'none', cursor: 'pointer',
            background: '#6366f1', color: '#fff', fontSize: 14, fontWeight: 700, flexShrink: 0,
          }}
        >
          {playing ? '❚❚' : '▶'}
        </button>
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 11, color: '#64748b', textTransform: 'uppercase', letterSpacing: 1, marginBottom: 4 }}>
            Activity Replay {current && <span style={{ color: '#94a3b8', textTransform: 'none' }}>· {current.time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>}
          </div>
          <div style={{ fontSize: 13, color: '#e2e8f0', minHeight: 18, marginBottom: 6 }}>
            {current ? (
              <span>
                <span style={{ color: KIND_COLOR[current.kind] }}>{KIND_ICON[current.kind]} {current.agent?.toUpperCase()}</span>{' '}
                {current.label}
              </span>
            ) : 'Select a step to begin'}
          </div>
          <input
            type="range" min={0} max={Math.max(steps.length - 1, 1)} value={step}
            onChange={(e) => setStep(Number(e.target.value))}
            style={{ width: '100%', accentColor: '#6366f1', cursor: 'pointer' }}
          />
        </div>
        <button
          onClick={() => { setPlaying(false); setStep(0); }}
          style={{ background: 'none', border: 'none', color: '#94a3b8', cursor: 'pointer', fontSize: 12 }}
        >
          ⟲
        </button>
      </div>
      <div style={{ fontSize: 10, color: '#475569', marginTop: 4 }}>
        {steps.length} events · Prompt → Agent → Tool calls → File changes → Checkpoint → Commit
      </div>
    </div>
  );
}