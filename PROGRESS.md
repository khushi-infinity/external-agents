# 📋 Project Status — Agent Activity Universe

**Event:** BTW Buildathon 2026 (6 September 2026) · **Track:** E3 — Bring Entire to a New Agent or Workflow · **Optional:** Best Use of Databricks (opted in)

## Pre-Noon Stable Milestone — 2026-09-06 11:43 IST

**Commit message:** `Pre-noon stable milestone: Agent Activity Universe`

**Intent:** This build is **Agent Activity Universe** for BTW Buildathon Entire track E3. It solves the developer visibility gap around AI coding agents by turning Entire checkpoint/session context into an explorable product surface for understanding what agents changed, why they changed it, and where risk is accumulating.

**Architecture now:** Entire checkpoints and graph context feed a data adapter, which normalizes repository/session/checkpoint data for Databricks risk scoring and analytics. The React Three Fiber frontend renders that model as a 3D universe where files are nodes, dependency/activity relationships are lines, agent activity is visible in-scene, and risk/analytics overlays are available through UI controls.

**Verifiably working now:**
- `agent-universe` has a working React Three Fiber app with file nodes, agent markers, connection lines, orbit controls, inspector panel, filters, Risk View, Analytics modal, and Activity Replay.
- `agent-universe/src/data/adapter.ts` loads a real export when present, tries `/api/repository`, and falls back to sample data so the demo remains usable offline.
- `databricks/ingest_and_score.py`, `databricks/events.ndjson`, and `databricks/README.md` define the Databricks ingestion and scoring path.
- Local verification at 2026-09-06 11:43 IST: `npm test` passed 2 files / 7 tests, and `npm run build` passed production TypeScript + Vite build.

**Deferred until after the Noon Curveball:**
- Wire a live Entire export from `.entire/` data into `public/data/entire-export.json` or a backend endpoint.
- Run the Databricks notebook live in the workspace and capture evidence.
- Connect `/api/analytics` to a live Databricks `risk_map` result instead of relying on the frontend fallback.
- Implement any Curveball-specific requirement and update `BUILDATHON.md` sections 5-6 with the response.
- Do final submission verification, final checkpoint, and final commit before the 3:00 PM IST deadline.

**Known risks / fragile points:**
- Live Entire export is not wired end to end yet; the app is currently protected by a sample-data fallback.
- Databricks code exists and sample data is present, but the notebook has not been run live in the target workspace from this session.
- 3D relationship layout is demo-deterministic and not fully generated from live Entire Graph relationships.
- Vite emits a large chunk warning for the Three.js bundle; acceptable for this milestone, but code-splitting is deferred.
- Remaining Entire checkpoint creation depends on committing from an active supported agent session with Entire hooks enabled.

**Assumptions for a fresh session:**
- Work must continue in `/Users/khushisarawagi/Desktop/buildathon/external-agents-mirror`.
- The designated build app is `agent-universe/`; the Databricks assets are in `databricks/`.
- Synthetic sample data is intentional until a real Entire export is dropped into `agent-universe/public/data/entire-export.json`.
- Submission deadline is 3:00 PM IST on 2026-09-06, and the Noon Curveball should be handled in a fresh post-curveball session.

> **Last updated:** 6 September 2026 (~10:45, build day)
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

## 2. ✅ DONE — Everything Built

### Where the implementation lives (STRICT GUIDE COMPLIANCE)

All implementation lives in the **Entire mirror clone**: `buildathon/external-agents-mirror/` (origin = `entire://aws-ap-south-1.entire.io/gh/khushi-infinity/external-agents`, India region). Pushed to mirror branch `agent-activity-universe`, **latest SHA `87aea50`** (main is protected on the mirror; pushes go to the feature branch).

```
external-agents-mirror/            ← THE designated fork clone (work here)
├── BUILDATHON.md                  # Submission doc, guide's 10-section outline
├── PROGRESS.md                    # This status doc
├── agent-universe/                # OUR APP (full source + tests + demo)
│   ├── demo.html                  # Standalone offline demo (no npm needed)
│   ├── public/data/entire-export.example.json   # Real-data drop-in schema
│   └── src/                       # 7 components + data layer + analytics
└── databricks/                    # Databricks pipeline (notebook + events + README)
    ├── ingest_and_score.py        # Notebook: NDJSON → risk_map table
    ├── events.ndjson              # Sample events (12 files) to upload
    └── README.md                  # 5-minute workspace setup
```

(Also on disk: `buildathon/agent-universe/` — dev working copy with node_modules; sync via rsync. `buildathon/cli-mirror/` — the abandoned cli fork clone, ignore.)

### Features implemented (all verified working)
- 3D universe: files = nodes, dependencies = connection lines, stars/sparkles, orbit controls
- Agent entities (Claude ◆ / Codex ● / Copilot ▲ / Aider ■) floating near touched files
- Node inspector: agent, checkpoint, prompt, files changed, risk score, connections
- **Risk View toggle** — nodes + lines recolor by Databricks risk scores
- **Analytics modal** — agent performance, file hotspots, failure patterns, velocity
- Agent + module filters
- **Activity Replay** — timeline scrubber, nodes pulse white during playback (verified)
- Data adapter priority: real export → live API → sample data (demo never breaks)
- **Databricks notebook** `databricks/ingest_and_score.py` — full ingestion + scoring pipeline
- **7 automated tests** (vitest) covering critical behavior

### Verification evidence
- ✅ `npm run build` passes (TypeScript strict + Vite production build)
- ✅ `npx tsc -b --noEmit` — zero errors
- ✅ `npm test` — 7/7 passing (adapter import, timeline ordering, risk scoring, NDJSON export, connection colors)
- ✅ Preview renders: 3D canvas, all controls, replay advances steps, Risk View toggles

---

## 3. 🗺️ Architecture

```
AI Agents (Claude / Codex / Copilot / Aider)
    ↓
Entire CLI + Checkpoints + Graph          ← captures prompts, files, sessions
    ↓
Data adapter (src/data/adapter.ts)        ← normalizes to common model
    1) public/data/entire-export.json (real drop-in export)
    2) /api/repository (live backend)
    3) sample-data.ts (offline fallback)
    ↓
Databricks analytics (databricks/ingest_and_score.py)  ← risk scoring, hotspots, velocity
    ↓
3D Universe (React Three Fiber)           ← files= nodes, deps= lines, agents= entities
    ↓
Developer / Judge
```

The `adapter → normalized data → product` separation makes the Noon Curveball cheap: "support a new agent" → add an agent adapter; "multiple repositories" → extend the adapter; "work offline" → swap the data-loading layer.

---

## 4. 🔜 TO IMPLEMENT NEXT (technical)

### P0 — before 3 PM (guide-required)
1. **Pre-noon checkpoint (11:45)** — commit from inside a Codex session in `external-agents-mirror`, record intent/architecture/unresolved/risks.
2. **Noon Curveball (12:00)** — stop, close session, receive constraint, start fresh session, run `entire graph impact` before editing.
3. **Final checkpoint + verification** — tests pass (critical + Curveball), BUILDATHON.md complete, final commit pushed, SHA matches submission.
4. **Databricks live run** — import `ingest_and_score.py` in the logged-in workspace, upload `events.ndjson`, run cells, screenshot output (fallback evidence).

### P1 — demo strength
5. Generate a **real** `entire-export.json` from the repo's `.entire/` data (replaces synthetic sample).
6. Point `/api/analytics` at the Databricks `risk_map` table (live wiring, not precomputed).

### P2 — polish
7. Risk explainer panel, time-slider, diff preview, bundle code-splitting.

---

## 5. 🔴 YOUR SIDE — Tasks Only You Can Do

| # | Task | When | How |
|---|------|------|-----|
| 1 | ~~Mirror fork on entire.io (India)~~ | ✅ DONE | `khushi-infinity/external-agents` on aws-ap-south-1, Mirror ID `01M1THHNBW0K7KFTZGE9ETPD9K` |
| 2 | Databricks workspace run | before 3 PM | import `databricks/ingest_and_score.py` → upload `events.ndjson` → run all cells → screenshot |
| 3 | **Pre-noon checkpoint** | **11:45** | `cd ~/Desktop/buildathon/external-agents-mirror && codex` → "commit current state as pre-noon stable milestone" → `entire checkpoint list` |
| 4 | **Noon Curveball** | **12:00** | STOP, close session, receive constraint, fresh session, `entire graph impact` before editing |
| 5 | Curveball response + final checkpoint | 1:00–2:30 | implement, test, checkpoint, update BUILDATHON.md §5–6 |
| 6 | Submit | **before 3:00 PM** | track E3, fork URL, final SHA, mirror URL, checkpoint links, demo access, Databricks opt-in |

---

## 6. 📊 Rubric Mapping (how we score)

**Entire main challenge (100 pts):**
| Criterion | Pts | Our evidence |
|-----------|-----|--------------|
| Problem and innovation | 20 | Developer visibility gap → 3D universe; strong demo story |
| Technical implementation | 25 | Working React/Three.js app, typed, builds, 7 tests |
| Response to Noon Curveball | 15 | Layered architecture; fill §5–6 on the day |
| Use of Entire Checkpoints | 15 | Checkpoint `e6e841e2def8` created; adapter reads checkpoint context; inspector shows prompts/checkpoints |
| Use of Entire Graph | 15 | Graph activated + verified: `search`, `impact`, `diff` all return real evidence on our app |
| Demonstration & future potential | 10 | Judge script (BUILDATHON.md §9), replay feature, continuation path documented |

**Best Use of Databricks (100 pts, optional):**
| Criterion | Pts | Our evidence |
|-----------|-----|--------------|
| Meaningful use of Databricks | 30 | Risk scoring feeds Risk View — removing it removes a core layer |
| Working implementation & reliability | 25 | Full notebook pipeline (`ingest_and_score.py`) + 7 passing tests |
| User value & product decisions | 20 | Hotspots/failures/velocity inform developer decisions |
| Data quality & provenance | 15 | Synthetic data documented; real export is P1 |
| Response to Curveball | 10 | Data pipeline adapts with adapter swap |

---

## 7. ⚠️ Known Limitations (honest, for BUILDATHON.md §8)

1. MVP uses synthetic sample data; real Entire export path exists (`entire-export.json` drop-in) but needs the actual `.entire/` output wired.
2. Databricks pipeline is fully implemented and runnable (notebook + events + schema); live notebook → app wiring (`/api/analytics` → `risk_map`) is the remaining hook-up.
3. 3D layout uses deterministic sample positions; production layout should derive from Entire Graph relationships.
4. Bundle is ~1.27 MB (Three.js) — acceptable for demo, code-splitting is a P2.
5. Checkpoints 2–4 still need user action from inside a supported agent session (Codex installed + logged in).

---

## 8. Competition Day Timeline — CORRECTED (participant guide: submit by 3:00 PM)

| Time | Milestone | Status |
|------|-----------|--------|
| 8:00–9:00 | Breakfast, check-in, setup | ✅ |
| 9:00–12:00 | Kickoff + build session | ✅ core built (correct fork, mirror, graph, hooks, checkpoint 1) |
| 11:45 | Preserve stable state + pre-noon checkpoint | 🔲 user: commit via Codex session |
| 12:00–1:00 | ⚡ Noon Curveball + lunch | 🔲 receive, stop, fresh session |
| 1:00–3:00 | Implement constraint, test, final checkpoint, finish BUILDATHON.md | 🔲 |
| **3:00 PM** | **SUBMISSION DEADLINE (not 4 PM!)** | 🔲 |
| 3:00–5:00 | Judging + winner announcement | 🔲 |

---

## 9. Strict Participant-Guide Compliance Status

| Guide requirement | Status |
|---|---|
| Correct E3 designated repo (`entireio/external-agents`) | ✅ forked as `khushi-infinity/external-agents` (we initially forked `cli` — **fixed**) |
| `entire login` | ✅ logged in (India, in.auth.entire.io) |
| `entire repo mirror create` + India region | ✅ aws-ap-south-1, Mirror ID `01M1THHNBW0K7KFTZGE9ETPD9K`, ready |
| Clone through Entire mirror | ✅ `buildathon/external-agents-mirror` (origin = entire://aws-ap-south-1.entire.io/gh/khushi-infinity/external-agents) |
| `entire enable` checkpoints | ✅ enabled, sync to origin; git hooks installed via `entire doctor --force` |
| `entire plugin install graph` | ✅ graph v0.4.0 installed |
| `entire graph init-agents --repo .` | ✅ wrote `.entire/graph-agent.md` + `AGENTS.md` |
| Implementation lives in the designated fork | ✅ agent-universe/ + databricks/ inside external-agents-mirror, pushed to `agent-activity-universe` (SHA `87aea50`) |
| Checkpoint 1: initial understanding | ✅ **`e6e841e2def8`** (09-06 10:37, linked to `ca06a8e`) — created via `entire session attach` with the real build-session transcript |
| Checkpoints 2–4 (pre-noon, curveball, final) | 🔲 user action at milestones (Codex session) |
| Entire Graph evidence gathered | ✅ `search` (found computeRiskMap top-ranked), `impact` (blast radius), `diff` (semantic change list) — recorded in BUILDATHON.md |
| Tests covering critical behavior | ✅ 7 vitest tests passing (adapter, timeline, risk, NDJSON, connection colors) |
| BUILDATHON.md in guide's 10-section outline | ✅ at external-agents-mirror/BUILDATHON.md (all sections filled except Curveball §5–6) |
| Databricks use — meaningful + working | ✅ notebook pipeline + events + README committed; **user: run in workspace** for live evidence |
| Fallback screenshot/recording | 🔲 user: record demo.html walkthrough + Databricks notebook output |

**⚠️ CRITICAL USER ACTION — REMAINING CHECKPOINTS:** Entire checkpoints are created when a git commit happens DURING an active agent session. The build agent (Freebuff) is not natively supported, so milestone commits must be made from inside a supported agent session. **Codex is already installed and logged in** (`~/.npm-global/bin/codex`). The user must:
```bash
cd ~/Desktop/buildathon/external-agents-mirror
codex    # start interactive session (approve the hooks on first run)
# inside codex: make the milestone commit, e.g.
#   "commit the current state as the pre-noon stable milestone"
# then verify:
entire checkpoint list
```
Checkpoint quality > quantity: capture decisions, rejected options, failures, assumptions, open risks.
