import { useMemo } from 'react';
import { Line } from '@react-three/drei';
import type { FileNode } from '../data/types';
import { useStore } from '../store';

export function Connection({ from, to }: { from: FileNode; to: FileNode }) {
  const riskView = useStore(s => s.riskViewEnabled);

  const color = useMemo(() => {
    if (riskView) {
      const avg = (from.risk_score + to.risk_score) / 2;
      if (avg >= 70) return '#ef4444';
      if (avg >= 50) return '#f59e0b';
      if (avg >= 30) return '#eab308';
    }
    return '#334155';
  }, [from, to, riskView]);

  const opacity = useMemo(() => {
    if (riskView) return 0.3 + ((from.risk_score + to.risk_score) / 200) * 0.5;
    return 0.25;
  }, [from, to, riskView]);

  const points = useMemo(() => [
    from.position,
    [
      (from.position[0] + to.position[0]) / 2,
      (from.position[1] + to.position[1]) / 2 + 0.3,
      (from.position[2] + to.position[2]) / 2,
    ],
    to.position,
  ] as [number, number, number][], [from, to]);

  return (
    <Line
      points={points}
      color={color}
      lineWidth={1.5}
      transparent
      opacity={opacity}
      dashed
      dashScale={5}
      dashSize={0.3}
      gapSize={0.2}
    />
  );
}