import { Canvas } from '@react-three/fiber';
import { OrbitControls, Stars, Sparkles } from '@react-three/drei';
import { useStore } from '../store';
import { FileNode3D } from './FileNode';
import { AgentEntity } from './AgentEntity';
import { Connection } from './Connection';

function UniverseScene() {
  const repository = useStore(s => s.repository);
  const activeAgents = useStore(s => s.activeAgents);
  const activeDirectories = useStore(s => s.activeDirectories);

  const visibleFiles = repository.files.filter(f => activeDirectories.includes(f.directory));
  const visibleSessions = repository.sessions.filter(s => activeAgents.includes(s.agent));
  const fileMap = Object.fromEntries(repository.files.map(f => [f.id, f]));
  const visibleIds = new Set(visibleFiles.map(f => f.id));

  return (
    <>
      <ambientLight intensity={0.4} />
      <pointLight position={[10, 10, 10]} intensity={1} />
      <pointLight position={[-10, -10, -10]} intensity={0.4} color="#6366f1" />
      <Stars radius={60} depth={40} count={2000} factor={3} saturation={0.6} fade speed={1} />
      <Sparkles count={120} scale={[14, 10, 8]} size={2} speed={0.4} color="#818cf8" />

      {/* Connections */}
      {visibleFiles.map(f =>
        f.connections
          .filter(cid => visibleIds.has(cid) && fileMap[cid])
          .map(cid => (
            <Connection key={`${f.id}-${cid}`} from={f} to={fileMap[cid]} />
          ))
      )}

      {/* File nodes */}
      {visibleFiles.map(f => <FileNode3D key={f.id} file={f} />)}

      {/* Agent entities */}
      {visibleSessions.map(s => <AgentEntity key={s.id} session={s} files={repository.files} />)}

      {/* Ground reference grid */}
      <gridHelper args={[20, 20, '#1e293b', '#0f172a']} position={[0, -4, 0]} />

      <OrbitControls enableDamping dampingFactor={0.08} minDistance={3} maxDistance={30} />
    </>
  );
}

export function Universe() {
  return (
    <div style={{ position: 'fixed', inset: 0, background: '#020617' }}>
      <Canvas
        camera={{ position: [10, 8, 12], fov: 50 }}
        dpr={[1, 2]}
        onPointerMissed={() => useStore.getState().setSelectedNode(null)}
      >
        <UniverseScene />
      </Canvas>
    </div>
  );
}