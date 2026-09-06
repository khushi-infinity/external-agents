import { useEffect } from 'react';
import { Universe } from './components/Universe';
import { Controls } from './components/Controls';
import { SidePanel } from './components/SidePanel';
import { ReplayPanel } from './components/ReplayPanel';
import { loadRepositoryData } from './data/adapter';
import { useStore } from './store';

export default function App() {
  const replayVisible = useStore(s => s.replayPlaying || s.replayStep > 0);

  useEffect(() => {
    loadRepositoryData().then(repo => useStore.setState({ repository: repo }));
  }, []);

  return (
    <div style={{ width: '100vw', height: '100vh', overflow: 'hidden', background: '#020617' }}>
      <Universe />
      <Controls />
      <SidePanel />
      {replayVisible ? <ReplayPanel /> : null}
    </div>
  );
}