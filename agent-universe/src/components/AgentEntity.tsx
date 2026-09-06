import { useRef, useMemo } from 'react';
import { useFrame } from '@react-three/fiber';
import type * as THREE from 'three';
import { Text, Billboard, Trail } from '@react-three/drei';
import type { AgentSession, FileNode } from '../data/types';
import { useStore } from '../store';

const COLORS: Record<string, string> = { claude: '#d97706', codex: '#10b981', copilot: '#6366f1', aider: '#ec4899' };
const ICONS: Record<string, string> = { claude: '◆', codex: '●', copilot: '▲', aider: '■' };

export function AgentEntity({ session, files }: { session: AgentSession; files: FileNode[] }) {
  const ref = useRef<THREE.Mesh>(null);
  const setSelectedNode = useStore(s => s.setSelectedNode);
  const color = COLORS[session.agent] || '#94a3b8';

  const pos = useMemo(() => {
    const touched = files.filter(f => session.files_touched.includes(f.id));
    if (!touched.length) return [0, 5, 0] as [number, number, number];
    return [
      touched.reduce((s, f) => s + f.position[0], 0) / touched.length,
      touched.reduce((s, f) => s + f.position[1], 0) / touched.length + 2,
      touched.reduce((s, f) => s + f.position[2], 0) / touched.length,
    ] as [number, number, number];
  }, [session, files]);

  useFrame(() => {
    if (!ref.current) return;
    const t = Date.now() * 0.0005;
    ref.current.position.set(
      pos[0] + Math.sin(t) * 0.3,
      pos[1] + Math.sin(t * 1.5) * 0.15,
      pos[2] + Math.cos(t) * 0.3
    );
    ref.current.rotation.y += 0.02;
  });

  return (
    <Trail width={0.6} length={6} color={color} attenuation={(w) => w * w}>
      <group>
        <mesh
          ref={ref}
          position={pos}
          onClick={(e) => { e.stopPropagation(); setSelectedNode({ type: 'session', data: session }); }}
          onPointerOver={() => { document.body.style.cursor = 'pointer'; }}
          onPointerOut={() => { document.body.style.cursor = 'default'; }}
        >
          <octahedronGeometry args={[0.25, 0]} />
          <meshStandardMaterial color={color} emissive={color} emissiveIntensity={0.6} roughness={0.2} metalness={0.8} />
        </mesh>
        <Billboard position={[pos[0], pos[1] + 0.5, pos[2]]}>
          <Text fontSize={0.18} color={color} anchorX="center" anchorY="bottom" outlineWidth={0.015} outlineColor="#000" fontWeight="bold">
            {`${ICONS[session.agent] || ''} ${session.agent.toUpperCase()}`}
          </Text>
        </Billboard>
        <Billboard position={[pos[0], pos[1] - 0.4, pos[2]]}>
          <Text fontSize={0.1} color="#94a3b8" anchorX="center" anchorY="top">
            {`${session.file_changes} files · ${session.tool_calls} calls`}
          </Text>
        </Billboard>
      </group>
    </Trail>
  );
}