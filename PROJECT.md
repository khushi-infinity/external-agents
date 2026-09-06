# 📘 Agent Activity Universe — Full Project Documentation

**BTW Buildathon 2026 · Track E3 · Best Use of Databricks**

An interactive 3D visualization of AI coding agents working across a codebase. It uses **Entire Checkpoints** to capture real development context (prompts, files changed, sessions) and **Databricks** to analyze agent behavior, activity, and development risk — rendered as a navigable 3D universe.

> *The 3D visualization is not the product. The product is: "Understand what AI coding agents are doing across your codebase, why they are doing it, and where their work is creating risk."*

---

## Table of Contents

1. [What This Project Does](#1-what-this-project-does)
2. [Architecture Overview](#2-architecture-overview)
3. [Component Deep-Dive](#3-component-deep-dive)
4. [Data Flow](#4-data-flow)
5. [The Noon Curveball Response](#5-the-noon-curveball-response)
6. [Databricks Integration](#6-databricks-integration)
7. [Live Verification Evidence](#7-live-verification-evidence)
8. [Tech Stack](#8-tech-stack)
9. [Known Limitations & Next Steps](#9-known-limitations--next-steps)

---

## 1. What This Project Does

### The Problem
AI coding agents (Claude, Codex, Copilot, Aider, Freebuff) are becoming the primary authors of code changes. But teams have **no centralized, structured view** of what all these agents are doing across their codebase:

- **What** did each agent change?
- **Why** did it make that change (what was the intent)?
- **Where** is the work creating risk (hotspot files, failure patterns)?
- **How** do agent behaviors compare (sessions, retries, success rate)?

Entire provides checkpoint infrastructure for tracking agent work — but only for agents it natively supports. Freebuff, the free agent used to build this project, had **no Entire support at all**, making its work invisible.

### Our Solution
We built **three interconnected layers** that make agent work visible, explainable, and actionable:

| Layer | Source | What the User Sees |
|-------|--------|--------------------|
| **Capture** | Entire Checkpoints + Freebuff Plugin | Sessions, prompts, files changed, tool calls, tokens |
| **Analyze** | Databricks Analytics | Risk scoring, hotspots, failure patterns, velocity |
| **Explore** | Agent Activity Universe (3D) | Interactive 3D universe with risk overlays |

### Who This Is For
- **Solo developers** using multiple AI agents who want to understand what's happening in their codebase
- **Teams** standardizing on free agents (like Freebuff) who need the same safety net as paid-agent users
- **Engineering managers** who need visibility into AI-assisted development risk

---

## 2. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                      CAPTURE LAYER                                  │
│  Freebuff Sessions ──► entire-agent-freebuff ──► Entire Checkpoints │
│  (manicode engine)    (external agent plugin)    (sessions, hooks,  │
│                                                transcripts, graph)  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ Entire Graph    │  │ Entire CLI      │  │ Entire Hooks    │
│ (dependencies,  │  │ (status,        │  │ (session-start, │
│  code structure)│  │  checkpoints)   │  │  prompt-submit, │
│                 │  │                 │  │  stop, end)     │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      ANALYZE LAYER                                  │
│  Databricks Notebook ──► Risk Scoring ──► risk_map table           │
│  (ingest_and_score.py)   (agent perf,    (per-file risk scores)    │
│                           hotspots,                               │
│                           failure patterns)                        │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      EXPLORE LAYER                                  │
│  Agent Activity Universe (React Three Fiber)                       │
│  ├── 3D Scene (universe with file nodes + agent markers)          │
│  ├── Inspector Panel (agent, checkpoint, intent, risk)            │
│  ├── Risk View (recolors universe from Databricks risk map)       │
│  ├── Analytics Panel (Databricks-derived metrics)                 │
│  ├── Activity Replay (session timeline animation)                 │
│  └── Agent/Module Filters                                          │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 3. Component Deep-Dive

### 3.1 Freebuff External Agent Plugin (`agents/entire-agent-freebuff/`)

**Purpose:** Makes Freebuff a first-class Entire agent by implementing the External Agent protocol.

**Key Files:**
- `cmd/entire-agent-freebuff/main.go` — CLI binary entry point
- `internal/protocol/` — External Agent protocol (info, detect, hooks, transcript, resume)
- `internal/freebuff/agent.go` — Freebuff session resolution + protocol wiring
- `internal/freebuff/transcript.go` — Dual-format transcript parser (original + JSONL)
- `internal/freebuff/hooks.go` — Hook lifecycle (session-start, prompt-submit, stop, session-end)
- `internal/freebuff/paths.go` — Freebuff session layout resolution (`~/.config/manicode/projects/...`)
- `internal/freebuff/types.go` — Data types for both transcript formats

**Protocol Subcommands:**
- `info` — Returns agent name, version, capabilities
- `detect` — Detects if the repo has Freebuff sessions
- `hooks` — Resolves hook lifecycle events
- `read-session`, `get-transcript` — Session/transcript reading
- `extract-prompts`, `extract-modified-files`, `extract-summary` — Transcript analysis
- `format-resume-command` — Resume command generation

**Hook Registry (`.freebuff/entire-hooks.json`):**
```json
{
  "commands": [
    {"hook": "session-start", "command": "entire hooks freebuff session-start"},
    {"hook": "prompt-submit", "command": "entire hooks freebuff prompt-submit"},
    {"hook": "stop",          "command": "entire hooks freebuff stop"},
    {"hook": "session-end",   "command": "entire hooks freebuff session-end"}
  ]
}
```

**Installation:**
```bash
cd agents/entire-agent-freebuff
mise run build && cp entire-agent-freebuff ~/.local/bin/
```

**Enable in a repo:**
```bash
entire enable --agent freebuff --local --telemetry=false
```

---

### 3.2 Agent Activity Universe (`agent-universe/`)

**Purpose:** Interactive 3D visualization that turns checkpoint context into a navigable universe.

**Tech Stack:**
- React 19 + TypeScript
- React Three Fiber (Three.js)
- Zustand (state management)
- Vite (build)
- Vitest (testing)

**Key Components:**
- `Universe.tsx` — 3D scene (Canvas, stars, file nodes, agent markers, connections)
- `InspectorPanel.tsx` — Click-to-inspect panel (agent, checkpoint, prompt, risk)
- `RiskView.tsx` — Toggleable risk overlay (recolors from Databricks data)
- `AnalyticsPanel.tsx` — Databricks-derived analytics (performance, hotspots, velocity)
- `ActivityReplay.tsx` — Session timeline animation
- `AgentFilter.tsx` — Agent toggle buttons (claude, codex, copilot, aider)
- `ModuleFilter.tsx` — Module filter pills (src/auth, src/payments, etc.)

**Data Flow:**
```
src/data/adapter.ts
  ├── Tries: /data/entire-export.json (live Entire export)
  ├── Tries: /api/repository (live API)
  └── Falls back: sample-data.ts (synthetic demo dataset)
```

**Offline Demo:** `agent-universe/demo.html` (single-file build, no network needed)

---

### 3.3 Databricks Analytics (`databricks/`)

**Purpose:** Risk scoring and analytics engine that powers the Risk View and Analytics panel.

**Pipeline:**
1. **Ingest** NDJSON events (classic `events.ndjson` + new JSONL format)
2. **Normalize** column names (`last_agent` → `agent`), fill defaults, derive module
3. **Compute** analytics:
   - File hotspots (most changed files + risk)
   - Module failure patterns (high-risk files per module)
   - Agent performance (sessions, avg time, success rate)
   - Development velocity (avg files/session, retries)
   - Per-file risk map (SAFE/LOW/MEDIUM/HIGH)
4. **Persist** `risk_map` table for the app's `/api/analytics`

**Key Files:**
- `ingest_and_score.py` — Main notebook (FileStore-based workspaces)
- `ingest_and_score_uc.py` — Unity Catalog variant (no `/FileStore`)
- `events.ndjson` — Classic format dataset
- `events-new-format.ndjson` — Curveball JSONL fixture (18 lines + 1 unknown event)

**Data Contract (`/api/analytics`):**
```json
{
  "agent_performance": {
    "claude": { "sessions": 42, "success_rate": 0.91, "avg_duration": 23 },
    "codex": { "sessions": 27, "success_rate": 0.84, "avg_duration": 19 }
  },
  "file_hotspots": [
    { "path": "src/auth/service.ts", "changes": 17, "risk": 82 }
  ],
  "failure_patterns": [
    { "module": "payments", "failures": 7, "risk": 71 }
  ],
  "velocity": {
    "avg_task_minutes": 23,
    "avg_files_per_session": 8.4,
    "avg_retries": 1.8
  },
  "risk_map": { "src/auth/service.ts": 82 }
}
```

---

## 4. Data Flow

### 4.1 Capture → Entire

```
Freebuff session starts
  → .freebuff/entire-hooks.json fires "session-start"
    → entire hooks freebuff session-start
      → Entire records session metadata

User submits prompt
  → "prompt-submit" hook fires
    → Entire records prompt + timestamp

Agent finishes turn
  → "stop" hook fires
    → Entire records completion + file changes

Session ends
  → "session-end" hook fires
    → Entire finalizes session transcript

User makes commit
  → Entire commit hook auto-creates checkpoint
    → Checkpoint = snapshot of intent + files + prompts
```

### 4.2 Entire → Databricks

```
Entire checkpoints export NDJSON events
  → databricks/events.ndjson (classic format)
    → Databricks notebook ingests
      → risk_map table (per-file risk scores)
        → /api/analytics endpoint
```

### 4.3 Databricks → Universe

```
Universe loads data via adapter
  → Tries /data/entire-export.json
  → Tries /api/repository
  → Falls back to sample-data.ts

User toggles Risk View
  → fetches /api/analytics
    → reads risk_map data
      → recolors 3D nodes by risk score

User opens Analytics panel
  → reads agent_performance, file_hotspots, etc.
    → displays metrics
```

---

## 5. The Noon Curveball Response

### What Changed
At 12:00 PM, the Curveball revealed: "The agent changed its format." Freebuff released a **new JSONL session-event format** alongside the original `chat-messages.json` array.

### Our Response

**Dual-Format Transcript Parser** (`internal/freebuff/transcript.go`):
- Sniffs content (leading `[` vs `{`) to detect format
- Normalizes both formats into one `turn` model
- Unknown JSONL events: skipped + counted, never fatal
- Incomplete transcripts: salvaged to longest valid prefix, flagged `partial`

**Test Coverage** (`internal/freebuff/transcript_test.go`):
- ✅ Original format (chat-messages.json array)
- ✅ New JSONL format (session events)
- ✅ Unknown events (counted, not crashed)
- ✅ Incomplete input (partial results)

**Databricks Adaptation** (cells 7b-7c):
- Tolerant struct JSON decode (nested objects never break a line)
- Column normalization at ingest (`last_agent` → `agent`)
- Unknown event tolerance (counted + skipped, never crash)
- Summary persistence (`pipeline_run_summary` table)

**Live Evidence:**
- `pipeline_run_summary`: 18 lines → 17 known → 1 unknown → **94.4% recognized**
- `risk_map`: 14 rows (classic files + Curveball fixture files merged)

---

## 6. Databricks Integration

### Live Run (Serverless, 6 Sep 2026)

**Workspace:** `https://dbc-e11b1b7c-7876.cloud.databricks.com` (org `7474646952123744`)

**Tables Created:**
- `workspace.agent_universe.risk_map` — 14 rows (per-file risk scores)
- `workspace.agent_universe.pipeline_run_summary` — Curveball tolerance evidence

**Evidence Queries:**
```sql
SELECT * FROM workspace.agent_universe.pipeline_run_summary;
-- curveball_ingest | 18 | 17 | 1 | 94.4

SELECT file, risk_score, change_count
FROM workspace.agent_universe.risk_map
ORDER BY risk_score DESC;
-- src/checkout/apply_coupon.ts | 100 | 27  (NEW format)
-- src/auth/service.ts         | 82  | 17  (classic)
```

**Successful Runs:**
1. `628954233613995` (fixed schema bug)
2. `128947458769841` (with evidence persistence)

---

## 7. Live Verification Evidence

### Entire Checkpoints
| Milestone | ID | Session | Agent | Commit |
|---|---|---|---|---|
| Initial | `e6e841e2def8` | `sess-buildathon-20260906-1015` | Claude Code | `ca06a8e` |
| Pre-noon | `2adb77572071` | `01a07557-…` | Codex | `10cb63c` |
| Curveball | `614fa84595b2` | `fb-curveball-001` | Freebuff | `dcfa571` |
| Final | `26d2aaa0ddf0` | `fb-curveball-001` | Freebuff | `0619daa` |
| Bonus | `9d720fabaf1e` | `fb-curveball-001` | Freebuff | `655a526` |

### Freebuff Plugin
```bash
entire checkpoint list     # 60 checkpoints, synced to origin
go test ./...              # all green (original + new + unknown + incomplete)
```

### 3D App
```bash
npm run build              # ✓ built
npm test                   # 7/7 pass
```

### Databricks
```bash
entire status              # Checkpoints sync to: origin (no pending)
```

---

## 8. Tech Stack

| Layer | Technology | Purpose |
|---|---|---|
| External Agent | Go 1.22+ | Freebuff plugin |
| Protocol | Entire External Agent Protocol | Agent ↔ Entire communication |
| Frontend | React 19 + TypeScript + Vite | 3D universe app |
| 3D Rendering | React Three Fiber (Three.js) | Navigable 3D scene |
| State | Zustand | Global state management |
| Analytics | Databricks (PySpark) | Risk scoring + analytics |
| Testing | Vitest (app), Go test (plugin) | Unit + integration tests |
| Build | mise (Go), npm (frontend) | Tooling + dependencies |

---

## 9. Known Limitations & Next Steps

### Limitations
1. **Freebuff engine hook wiring** — Plugin writes declarative registry; engine auto-fire is a Freebuff-side next step. Hooks can be driven manually today (verified live).
2. **Headless Freebuff** — No `freebuff -p` flag yet; lifecycle e2e needs an interactive session.
3. **3D layout** — Uses deterministic sample positions (not force-directed).
4. **Databricks wiring** — App's `/api/analytics` endpoint wired to sample data; live `risk_map` wiring is post-demo.

### Next Steps
1. Freebuff engine hook auto-loading (plugin + engine agree on `.freebuff/entire-hooks.json`)
2. Live Entire checkpoint → 3D universe adapter (real `.entire/` export)
3. Databricks → app live analytics endpoint (`/api/analytics` → `risk_map`)
4. Headless Freebuff (`freebuff -p`) for fully automated lifecycle e2e
5. Force-directed 3D layout based on real dependency graph

---

## Quick Start

```bash
# 1. Clone the repo
git clone https://github.com/khushi-infinity/external-agents.git
cd external-agents

# 2. Build + install the Freebuff plugin
cd agents/entire-agent-freebuff
mise run build && cp entire-agent-freebuff ~/.local/bin/

# 3. Enable in a repo
entire enable --agent freebuff --local

# 4. Run the 3D app
cd agent-universe && npm install && npm run dev
# → http://localhost:5173

# 5. Run tests
go test ./...            # plugin tests
cd agent-universe && npm test  # app tests

# 6. Checkpoints
entire checkpoint list   # see the 5 milestones
```

**Offline demo:** open `agent-universe/demo.html` directly in any browser.
