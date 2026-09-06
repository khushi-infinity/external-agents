# 📋 Project Status — Agent Activity Universe

**Event:** BTW Buildathon 2026 (6 September 2026) · **Track:** E3 — Bring Entire to a New Agent or Workflow · **Optional:** Best Use of Databricks (opted in)

> **Last updated:** 6 September 2026 (during the build day)
>
> ⚠️ **SOURCE OF TRUTH #1:** The PDF **`/Users/khushisarawagi/Downloads/What we're building.pdf`** (33 pages) is the authoritative build plan.
> ⚠️ **SOURCE OF TRUTH #2 (STRICT RULES):** The PDF **`/Users/khushisarawagi/Downloads/BTW Buildathon 2026 - Participant Guide.pdf`** (9 pages) is the operating guide and MUST be followed strictly. Key rules: (1) **SUBMISSION DEADLINE IS 3:00 PM IST, NOT 4:00 PM**; (2) all implementation must happen in the clone created through the Entire mirror workflow; (3) required checkpoints at 4 milestones (initial understanding → pre-noon stable → curveball response → final verification); (4) Entire Graph activation required (`entire plugin install graph`, `entire graph init-agents --repo .`); (5) final checklist requires tests covering critical + Curveball behavior; (6) BUILDATHON.md must follow their exact 10-section outline; (7) fork only after official start, mirror with India region; (8) demo owner must be able to sign in and run the critical path; (9) no secrets anywhere; (10) fallback screenshot/recording for fragile live steps. Key directives: (1) the 3D visualization is NOT the product — the product is *"understand what AI coding agents are doing, why, and where risk is created"*; (2) Entire checkpoint/session data MUST be central — *"do not build a generic Three.js dashboard"*; (3) build the smallest useful version with sample data FIRST, then connect real Entire data, then polish, then Databricks; (4) Activity Replay is a priority demo feature; (5) the hero screen is the UNIVERSE, dashboard is secondary; (6) do NOT build login/registration/databases/billing/dashboards/chat. PDF phases: inspect repo → inspect Entire → architecture → backend/data adapter → 3D frontend → connect real data → test → polish → checkpoint + commit.

---

## 1. What the Product Is

**Agent Activity Universe** — an interactive 3D visualization of AI coding agents working across a codebase. It turns Entire checkpoint/session context (prompts, files changed, sessions) into a navigable 3D universe, with Databricks-powered risk/analytics layered on top.

**The product statement:** *"Understand what AI coding agents are doing across your codebase, why they are doing it, and where their work is creating risk."*

| Layer | Source | Role |
|-------|--------|------|
| What did the agent do? | Entire Checkpoints | sessions, prompts, files changed, tool calls |
| Where does the work connect? | Entire Graph / code structure | 3D nodes + dependency connections |
| What patterns & risks emerge? | Databricks analytics | risk overlays, hotspots, failure patterns, velocity |
| Our product | React Three Fiber | interactive 3D universe |

---

## 2. ✅ DONE — Everything Built (1008 LOC)

### Project structure
```
buildathon/
├── BUILDATHON.md                  # Submission doc + judge demo script (§9)
├── cli/                           # (existing) Entire CLI fork — competition repo
└── agent-universe/                # OUR APP
    ├── demo.html                  # Standalone offline demo (1.2 MB, no npm needed)
    ├── index.html                 # Vite entry (title updated)
    ├── package.json               # React 19, TS, Vite 8, Three.js, R3F, drei, zustand
    ├── public/
    │   └── data/
    │       └── entire-export.example.json   # Schema example for real data drop-in
    └── src/
        ├── main.tsx               # React root
        ├── App.tsx                # Wires Universe + Controls + SidePanel + Replay
        ├── index.css              # Dark space theme, reset
        ├── store.ts               # Zustand: data, selection, filters, risk view, replay
        ├── data/
        │   ├── types.ts           # FileNode, AgentSession, Checkpoint, Repository, Analytics
        │   ├── sample-data.ts     # Demo dataset (auth/payments/database/api/test clusters)
        │   └── adapter.ts         # Data loader + importEntireExport + buildTimeline
        ├── components/
        │   ├── Universe.tsx       # 3D scene: stars, sparkles, grid, orbit controls
        │   ├── FileNode.tsx       # File nodes (risk color, hover, click, replay pulse)
        │   ├── AgentEntity.tsx    # Floating agent markers with motion trails
        │   ├── Connection.tsx     # Dependency lines (dashed, risk-tinted)
        │   ├── SidePanel.tsx      # Inspector: file / session / checkpoint details
        │   ├── Controls.tsx       # Top bar + agent/module filters + analytics modal
        │   └── ReplayPanel.tsx    # Activity Replay timeline scrubber
        └── analytics/
            └── databricks.ts      # Databricks pipeline docs + local risk fallback + NDJSON export
```

### Features implemented (all verified working)
| # | Feature | Status |
|---|---------|--------|
| 1 | 3D repository universe (files = nodes, deps = lines) | ✅ renders |
| 2 | Agent entities floating near touched files | ✅ |
| 3 | Click node → inspector panel (agent, checkpoint, prompt, risk, connections) | ✅ |
| 4 | Risk View toggle (nodes + connections recolor by risk) | ✅ verified red "ON" state |
| 5 | Analytics modal (agent performance, hotspots, failure patterns, velocity) | ✅ |
| 6 | Agent filters (claude / codex / copilot / aider) | ✅ |
| 7 | Module filters (src/auth, src/payments, src/database, src/api, tests, src) | ✅ |
| 8 | **Activity Replay** — animated timeline, nodes pulse white | ✅ verified advancing 5→7 |
| 9 | Standalone offline demo (`demo.html`) | ✅ HTTP 200 |
| 10 | Drop-in real Entire data (`public/data/entire-export.json`) | ✅ wired, schema documented |

### Verification evidence
```bash
cd buildathon/agent-universe
npx tsc -b --noEmit   # ✅ zero errors
npm run build         # ✅ production build passes (Three.js chunk-size warning only)
npm run dev           # ✅ serves on :5173
open demo.html        # ✅ works fully offline
```

---

## 3. 🗺️ Architecture

```
AI Agents (Claude / Codex / Copilot / Aider)
    ↓
Entire CLI + Checkpoints + Graph        ← captures development context
    ↓
DATA ADAPTER (src/data/adapter.ts)      ← normalizes to common model
    ├─ 1) public/data/entire-export.json (real export drop-in)
    ├─ 2) /api/repository (live backend)
    └─ 3) sample-data.ts (offline fallback)
    ↓
DATABRICKS (src/analytics/databricks.ts) ← risk scoring, hotspots, velocity
    ↓
3D UNIVERSE (React Three Fiber)          ← files, agents, connections, risk layers
    ↓
Developer / Judge
```

**Curveball resilience:** the `adapter → normalized data → product` separation means:
- "Support multiple repos" → extend adapter
- "Support a new agent" → add agent adapter
- "Show human activity" → add event type
- "Work offline" → data-loading layer swap

---

## 4. 🔜 TO IMPLEMENT NEXT (technical)

### Priority order (highest value first)

**P0 — Must do before 4 PM code freeze:**
1. **Move/commit app into the designated Entire fork.** E3 rule: *"your submitted implementation must live in the designated Entire fork."* The app currently lives in `buildathon/agent-universe/` — copy it into `buildathon/cli/` (e.g. `cli/agent-universe/`) and commit, or confirm with mentors.
2. **Mirror fork on entire.io, India region** (user task — competition rule).
3. **Pre-noon commit + checkpoint (~11:45):** `git add . && git commit -m "feat: pre-noon stable agent activity universe"` then `entire status` to confirm checkpoint.
4. **Noon Curveball:** stop, fresh agent session, fill `BUILDATHON.md` §5–6.

**P1 — Demo strength (do if time before noon or in Curveball window):**
5. **Real Entire data export:** parse `entire status --json` / `.entire/` from the cli fork → generate real `entire-export.json`. Replaces the synthetic sample dataset with genuine competition data — huge judge credibility.
6. **Databricks ingestion job:** export events as NDJSON → load into Databricks SQL Warehouse/notebook → run risk/hotspot scoring → serve via `/api/analytics`. Completes the "meaningful Databricks use" story (currently precomputed analytics).

**P2 — Nice-to-have:**
7. Natural-language risk explainer in inspector ("this file is risky because: 17 changes, 4 failed tests, 14 downstream dependents").
8. Time-slider scrubbing across the whole day, not just single sessions.
9. Diff preview panel when clicking a checkpoint.
10. Code-split the bundle (dynamic import of Three.js) to remove the chunk-size warning.

---

## 5. 🔴 YOUR SIDE — Tasks Only You Can Do

| # | Task | When | How |
|---|------|------|-----|
| 1 | Mirror `khushi-infinity/cli` on entire.io | NOW / before noon | entire.io → import GitHub fork → select **India region** → clone from mirror |
| 2 | Databricks Free Edition signup | NOW / during break | databricks.com/try-databricks → keep workspace URL + credentials for demo |
| 3 | Verify `entire status` shows checkpoints | before noon | `cd buildathon/cli && entire status` |
| 4 | Commit pre-noon stable state | 11:45 | `git commit -m "feat: pre-noon stable agent activity universe"` + verify checkpoint |
| 5 | Receive Curveball, start fresh session | 12:00 | STOP current agent; new session uses checkpoint + Graph context |
| 6 | Fill submission form fields | 3:30–4:00 | project name, track E3, repo URL, commit SHA, demo access, BUILDATHON.md, Databricks opt-in fields |

---

## 6. 📊 Rubric Mapping (how we score)

**Entire main challenge (100 pts):**
| Criterion | Pts | Our evidence |
|-----------|-----|--------------|
| Problem and innovation | 20 | Developer visibility gap → 3D universe; strong demo story |
| Technical implementation | 25 | Working React/Three.js app, typed, builds, tests via verification |
| Response to Noon Curveball | 15 | Layered architecture; fill §5–6 on the day |
| Use of Entire Checkpoints | 15 | Adapter reads checkpoint context; inspector shows prompts/checkpoints |
| Use of Entire Graph | 15 | Connections derive from code relationships (extendable to real Graph data) |
| Demonstration & future potential | 10 | Judge script (§9), replay feature, continuation path documented |

**Best Use of Databricks (100 pts, optional):**
| Criterion | Pts | Our evidence |
|-----------|-----|--------------|
| Meaningful use of Databricks | 30 | Risk scoring feeds Risk View — removing it removes a core layer |
| Working implementation & reliability | 25 | analytics module + NDJSON export; live job is P1 |
| User value & product decisions | 20 | Hotspots/failures/velocity inform developer decisions |
| Data quality & provenance | 15 | Synthetic data documented; real export is P1 |
| Response to Curveball | 10 | Data pipeline adapts with adapter swap |

---

## 7. ⚠️ Known Limitations (honest, for BUILDATHON.md §8)

1. MVP uses synthetic sample data; real Entire export path exists but needs the actual `.entire/` output wired.
2. Databricks integration is a documented pipeline + precomputed analytics — live ingestion job is the next milestone.
3. 3D layout uses deterministic sample positions; production layout should derive from Entire Graph relationships.
4. No automated unit tests yet (verification is build + manual + tsc).
5. Bundle is 1.27 MB (Three.js) — acceptable for demo, code-splitting is a P2.

---

## 8. Competition Day Timeline — CORRECTED (participant guide: submit by 3:00 PM)

| Time | Milestone | Status |
|------|-----------|--------|
| 8:00–9:00 | Breakfast, check-in, setup | ✅ |
| 9:00–12:00 | Kickoff + build session | ✅ core built (mirror, graph, hooks ready) |
| 11:45 | Preserve stable state + pre-noon checkpoint | 🔲 user: commit via agent session |
| 12:00–1:00 | ⚡ Noon Curveball + lunch | 🔲 receive, stop, fresh session |
| 1:00–3:00 | Implement constraint, test, final checkpoint, finish BUILDATHON.md | 🔲 |
| **3:00 PM** | **SUBMISSION DEADLINE (not 4 PM!)** | 🔲 |
| 3:00–5:00 | Judging + winner announcement | 🔲 |

## 9. Strict Participant-Guide Compliance Status

| Guide requirement | Status |
|---|---|---|
| Fork only after official start | ✅ khushi-infinity/cli fork |
| `entire login` | ✅ logged in (India, in.auth.entire.io) |
| `entire repo mirror create` + India region | ✅ khushi-infinity/cli on aws-ap-south-1, ready |
| Clone through Entire mirror | ✅ buildathon/cli-mirror (origin = entire://aws-ap-south-1.entire.io) |
| `entire enable` checkpoints | ✅ enabled, sync to origin; hooks installed via `entire doctor --force` |
| `entire plugin install graph` | ✅ graph v0.4.0 installed |
| `entire graph init-agents --repo .` | ✅ wrote .entire/graph-agent.md + AGENTS.md |
| Implementation lives in the designated fork | ✅ agent-universe/ inside cli-mirror, pushed to branch `agent-activity-universe` (SHA 6c5cfc2) |
| Checkpoint 1: initial understanding | 🔲 requires a commit during an ACTIVE agent session (Codex is installed + logged in) |
| Checkpoints 2–4 (pre-noon, curveball, final) | 🔲 user action at milestones |
| Tests covering critical behavior | ✅ 7 vitest tests passing (adapter, timeline, risk, NDJSON) |
| BUILDATHON.md in guide's 10-section outline | ✅ at cli-mirror/BUILDATHON.md |
| Databricks Free Edition account | 🔲 user: sign up before/at venue |
| Fallback screenshot/recording | 🔲 user: record demo.html walkthrough as backup |

**⚠️ CRITICAL USER ACTION — CHECKPOINTS:** Entire checkpoints are created when a git commit happens DURING an active agent session. The build agent (Freebuff) is not natively supported, so milestone commits must be made from inside a supported agent session. **Codex is already installed and logged in** (`~/.npm-global/bin/codex`, ChatGPT plan). The user must:
```bash
cd ~/Desktop/buildathon/cli-mirror
codex    # start interactive session (approve the 7 hooks on first run)
# inside codex: make the milestone commit, e.g.
#   "commit the current state as the pre-noon stable milestone"
# then verify:
entire checkpoint list
```
Checkpoint quality > quantity: capture decisions, rejected options, failures, assumptions, open risks.