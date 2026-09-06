# ✦ Agent Activity Universe

**BTW Buildathon 2026 — Track E3: Bring Entire to a New Agent or Workflow**
**+ Best Use of Databricks (optional award)**

An interactive 3D visualization of AI coding agents working across a codebase. It uses **Entire Checkpoints** to capture real development context (prompts, files changed, sessions) and **Databricks** to analyze agent behavior, activity, and development risk — rendered as a navigable 3D universe.

> The 3D visualization is not the product. The product is: *"Understand what AI coding agents are doing across your codebase, why they are doing it, and where their work is creating risk."*

---

## What It Does

| Layer | Source | What the user sees |
|-------|--------|--------------------|
| **What did the agent do?** | Entire Checkpoints | Sessions, prompts, files changed, tool calls, tokens |
| **Where does the work connect?** | Entire Graph / code structure | 3D nodes + dependency connections |
| **What patterns & risks emerge?** | Databricks analytics | Risk overlays, hotspots, failure patterns, velocity |
| **Our product** | React Three Fiber | Interactive 3D universe |

### Demo Flow (for judges)
1. Open the app → a 3D repository universe appears (files = nodes, dependencies = lines)
2. Agent entities (◆ Claude, ● Codex, ▲ Copilot, ■ Aider) float near the files they touched
3. Click any file node → inspector panel shows: agent, checkpoint, prompt, files changed, risk score
4. Toggle **Risk View** → Databricks-derived risk colors/glows highlight risky areas
5. Open **Analytics** → agent performance, file hotspots, failure patterns, velocity

---

## Quick Start

```bash
cd agent-universe
npm install
npm run dev
# open http://localhost:5173
```

Production build:

```bash
npm run build
npm run preview
```

---

## Project Structure

```
agent-universe/
├── src/
│   ├── App.tsx                 # Main app (Universe + Controls + SidePanel)
│   ├── store.ts                # Zustand global state
│   ├── data/
│   │   ├── types.ts            # Data model (FileNode, AgentSession, Checkpoint, Analytics)
│   │   ├── sample-data.ts      # Fallback demo dataset (synthetic)
│   │   └── adapter.ts          # Entire data loader (falls back to sample)
│   ├── components/
│   │   ├── Universe.tsx        # 3D scene (Canvas, stars, nodes, agents, connections)
│   │   ├── FileNode.tsx        # 3D file node (risk color, hover, click, replay pulse)
│   │   ├── AgentEntity.tsx     # Floating agent marker with trail
│   │   ├── Connection.tsx      # Dependency line between files
│   │   ├── SidePanel.tsx       # Inspector for selected node
│   │   ├── ReplayPanel.tsx     # Activity Replay timeline scrubber
│   │   └── Controls.tsx        # Top bar + filters + replay + analytics modal
│   └── analytics/
│       └── databricks.ts       # Databricks integration + local risk fallback
└── index.html
```

---

## How Entire Powers This (E3 Fit)

The data adapter (`src/data/adapter.ts`) is designed to read **real Entire checkpoint/session data**. Data loading priority:
1. **Drop a real export at `public/data/entire-export.json`** (schema: `{ name, files: [{ path, risk_score, change_count, last_agent, last_checkpoint, last_prompt, connections }], sessions, checkpoints }`) — the app loads it automatically.
2. Live backend at `/api/repository`.
3. Sample data fallback (always works).

The `importEntireExport()` function also accepts JSON exports from the Entire dashboard. In a full deployment the adapter parses `entire status --json` output, checkpoint refs in `.entire/`, and session transcripts.

**Why this qualifies for E3:** it brings Entire into a new workflow — *development intelligence/observability* — where Entire checkpoint context is an essential input. Removing Entire removes the product's core data. It also uses Entire Graph concepts (code relationships drive the 3D connections, not arbitrary placement).

## How Databricks Powers This

`src/analytics/databricks.ts` defines the pipeline: export normalized events (NDJSON) → ingest into Databricks → run analytics (agent performance, file hotspots, failure patterns, velocity, risk scoring) → feed results back to the 3D view's **Risk View** and **Analytics** modal. The MVP ships with precomputed analytics so the demo never blocks; the live endpoint (`/api/analytics`) is wired and used when available.

## Activity Replay

Click **▶ Activity Replay** to watch a session timeline animate in the 3D scene: session start → file changes (affected nodes pulse white) → checkpoints. The timeline is built from session/checkpoint timestamps via `buildTimeline()` in `src/data/adapter.ts`.

## Noon Curveball Strategy

The layered architecture (`Entire adapter → normalized data → product`) makes adaptation cheap:
- "Support multiple repositories" → extend the adapter, UI unchanged
- "Support a new agent" → add an agent adapter
- "Show human developer activity" → add an event type
- "Work offline" → swap the data-loading layer

## Tech Stack

React 19 · TypeScript · Vite · Three.js · React Three Fiber · @react-three/drei · Zustand