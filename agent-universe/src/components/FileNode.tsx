import { useRef, useState } from 'react';
import { useFrame } from '@react-three/fiber';
import { Text, Billboard } from '@react-three/drei';
import type { FileNode as FileNodeType } from '../data/types';
import { useStore } from '../store';
import * as THREE from 'three';

const EXT_COLORS: Record<string, string> = {
  '.ts': '#3178c6', '.tsx': '#61dafb', '.js': '#f7df1e', '.py': '#3776ab',
  '.go': '#00add8', '.css': '#264de4', '.json': '#94a3b8', '.md': '#e2e8f0',
};

function riskColor(r: number) {
  if (r >= 70) return '#ef4444';
  if (r >= 50) return '#f59e0b';
  if (r >= 30) return '#eab308';
  return '#22c55e';
}

export function FileNode3D({ file }: { file: FileNodeType }) {
  const meshRef = useRef<THREE.Mesh>(null);
  const [hovered, setHovered] = useState(false);
  const setSelectedNode = useStore(s => s.setSelectedNode);
  const riskView = useStore(s => s.riskViewEnabled);
  const replayActiveFile = useStore(s => s.replayActiveFile);

  const isReplayActive = replayActiveFile === file.id;
  const color = isReplayActive
    ? '#ffffff'
    : riskView
      ? riskColor(file.risk_score)
      : (EXT_COLORS[file.extension] || '#64748b');
  const baseScale = 0.12 + (file.change_count / 20) * 0.18;
  const target = (hovered || isReplayActive) ? baseScale * 1.5 : baseScale;

  useFrame((_, dt) => {
    if (!meshRef.current) return;
    const s = meshRef.current.scale.x;
    meshRef.current.scale.setScalar(THREE.MathUtils.lerp(s, target, dt * 8));
    meshRef.current.position.y = file.position[1] + Math.sin(Date.now() * 0.001 + file.position[0]) * 0.05;
  });

  return (
    <group position={file.position}>
      <mesh
        ref={meshRef}
        scale={baseScale}
        onPointerOver={(e) => { e.stopPropagation(); setHovered(true); document.body.style.cursor = 'pointer'; }}
        onPointerOut={() => { setHovered(false); document.body.style.cursor = 'default'; }}
        onClick={(e) => { e.stopPropagation(); setSelectedNode({ type: 'file', data: file }); }}
      >
        <dodecahedronGeometry args={[1, 0]} />
        <meshStandardMaterial color={color} emissive={color} emissiveIntensity={isReplayActive ? 1.5 : hovered ? 0.8 : 0.3} roughness={0.3} metalness={0.6} transparent opacity={0.9} />
      </mesh>
      {isReplayActive && (
        <pointLight color="#ffffff" intensity={3} distance={3} decay={2} />
      )}
      <Billboard position={[0, baseScale + 0.3, 0]}>
        <Text fontSize={0.14} color={hovered || isReplayActive ? '#ffffff' : '#a3b3cc'} anchorX="center" anchorY="bottom" outlineWidth={0.01} outlineColor="#000">
          {file.name}
        </Text>
      </Billboard>
      {riskView && !isReplayActive && file.risk_score >= 70 && (
        <pointLight color={riskColor(file.risk_score)} intensity={2} distance={2} decay={2} />
      )}
    </group>
  );
}