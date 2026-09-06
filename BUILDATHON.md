# Agent Activity Universe

## One-sentence summary
An interactive 3D visualization that turns Entire checkpoint/session context and Databricks risk analytics into a navigable universe of what AI coding agents are doing across a codebase, why, and where risk is accumulating.

## Problem, intended user and why it matters
**User:** a developer or team lead working with AI coding agents (Claude Code, Codex, Copilot, Aider).

**Problem:** AI coding agents generate code faster than developers can follow. Developers have poor visibility into what agents are doing across the codebase, why they are doing it (the intent behind changes), how work connects across sessions and checkpoints, which files are repeatedly touched, and where development risk is accumulating. Git shows *what* changed; Entire checkpoints preserve *why*; nothing made both legible at a glance.

**Why it matters:** teams cannot review, hand off, or trust agent work they cannot see. Our product turns invisible agent activity into an inspectable, risk-aware 3D development universe.

## Selected Entire track and why Entire is essential
**Track E3 — Bring Entire to a New Agent or Workflow.**

We bring Entire into a *development-intelligence / observability workflow*: the product's core input is actual Entire checkpoint/session context (prompts, files changed, sessions, tool calls). The 3D visualization is not the product by itself — the product is understanding agent work and its risk. Removing Entire removes the product's data source entirely; this is not a generic Three.js dashboard bolted onto Entire.

The app also demonstrates the Entire Graph integration path: file relationships drive the 3D connections (via the data adapter's import/parse layer), and the app is designed to consume `entire status --json` / `.entire/` checkpoint output as its live data feed.

## Architecture and main workflow

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

The `adapter → normalized data → product` separation makes the Noon Curveball cheap: "support a new agent" → add an agent adapter; "support multiple repositories" → extend the adapter; "work offline" → swap the data-loading layer.

## Entire Graph findings and verification
Live output captured from the mirror clone (`entire graph ... --repo .`), commit `ca06a8e`:

- **Graph search:** `entire graph search --query "compute risk score for a file from session analytics"` → ranked `computeRiskMap` in `agent-universe/src/analytics/databricks.ts:35` (score 37.7) as top hit, then `sample-data.ts` and `Connection.tsx` risk coloring — proving the graph indexes our app code semantically, not just the Go repo.
- **Relationship / impact analysis before a high-risk change:** `entire graph impact --symbol computeRiskMap` → 1 direct caller (`databricks.test.ts`), 0 callees, 2 type consumers (`Repository` in `src/data/types.ts:46`) — the exact blast-radius picture needed before touching the risk-scoring function.
- **Final semantic-diff analysis of the submitted implementation:** `entire graph diff --base dcd5c4d --head ca06a8e` → entity-level change list (e.g., `PROGRESS.md body changed, 0 dependents`), confirming the semantic-diff workflow used for the Curveball response.

Graph results are evidence, not an oracle — each finding is verified against source code and tests before being recorded here. Run the same commands in the mirror clone to reproduce.

## Noon Curveball: what changed and how we adapted
*(To be filled after 12:00 noon reveal.)*

- Constraint received: TBD
- Assumptions revisited: TBD
- What changed / what stayed intact: TBD
- Test proving the revised behavior: TBD

## Checkpoint links and what each checkpoint proves

**Repo:** GitHub fork `github.com/khushi-infinity/external-agents` · **Mirror:** `entire://aws-ap-south-1.entire.io/gh/khushi-infinity/external-agents` (India region, Mirror ID `01M1THHNBW0K7KFTZGE9ETPD9K`) · **Branch:** `agent-activity-universe` · **Pushed SHA:** `10cb63c` (feature branch; `main` is protected on the mirror, so pushes go to the feature branch and are verified with `entire checkpoint list`).

Checkpoints are created when a commit happens during an ACTIVE agent session (hooks installed, agent authed). Milestone commits are made from inside a supported agent session (`codex` in the mirror clone), then verified with `entire checkpoint list`.

| Milestone | Commit | Checkpoint | What it proves |
|-----------|--------|------------|----------------|
| Initial understanding & intended architecture | `a081600` (amended `ca06a8e`) | `e6e841e2def8` | Captures the build plan: E3 track, Entire-central architecture, PDF scope decisions |
| Pre-noon stable state (11:45) | `10cb63c` | `2adb77572071` | Runnable product before the Curveball |
| Curveball response (12:00+) | TBD | TBD | Adaptation implemented, tested, explained |
| Final implementation & verification (before 3 PM) | TBD | TBD | Tests pass, BUILDATHON.md complete, graph evidence |

## Setup, run and test instructions

```bash
# The app lives in agent-universe/ inside this fork
cd agent-universe

# Install
npm install

# Run in development
npm run dev
# open http://localhost:5173

# Run tests
npm test

# Production build
npm run build
```

**Offline demo:** open `agent-universe/demo.html` in any browser — a standalone single-file build that works without npm or internet.

## Databricks use, data sources and limitations

**Role of Databricks:** the risk/analytics layer, and the reason the product can *score* agent work rather than just display it. Databricks analyzes accumulated development events (agent sessions, file changes, checkpoints, test outcomes) to produce agent performance, file hotspots, failure patterns, development velocity, and per-file risk scores that drive the **Risk View** overlay and the **Analytics** panel. Removing Databricks removes the risk-scoring layer — a core workflow of the product (award criteria: "Databricks must power a meaningful part of the data/AI workflow; removing it materially reduces functionality").

**Pipeline (implemented in `databricks/`):**
1. `databricks/events.ndjson` — sample development events (12 files) in the exact NDJSON schema the notebook consumes.
2. `databricks/ingest_and_score.py` — Databricks notebook: ingest NDJSON → auto-create `agent_universe` catalog/schema → compute agent performance, file hotspots, failure patterns, velocity, and a per-file risk map → write `agent_universe.risk_map` table.
3. App side: `exportEventsNDJSON()` in `src/analytics/databricks.ts` emits the same schema from the app's data, so the notebook consumes *real* app state; the app's `computeRiskMap` mirrors the notebook's scoring for offline demo parity.

**Live run:** import `ingest_and_score.py` into a Databricks workspace (Create → Import), upload `events.ndjson` (Data → Add Data), start the serverless 2X-Small SQL warehouse, run all cells, screenshot the output table for the submission's fallback evidence.

**Data sources:** synthetic sample dataset (clearly labeled in `src/data/sample-data.ts`); no private, customer, or personal data used. Live Entire data is loaded from the repository's own checkpoint context when an export is provided. No credentials are stored in the repository.

## Known limitations and next steps
**Limitations:**
- MVP renders a synthetic sample dataset; the live Entire export path (`public/data/entire-export.json`) is wired and documented but needs the repository's actual `.entire/` output to be exercised.
- Databricks pipeline is fully implemented and runnable (notebook + events + schema); the app consumes the same scoring logic client-side for offline demo parity — live notebook → app wiring is the remaining hook-up.
- 3D layout uses deterministic sample positions; production layout should derive from Entire Graph relationships.
- No Curveball-verified behavior yet — will be added at noon.

**Next steps:**
1. Wire live Entire checkpoint parsing into the adapter.
2. Run the Databricks notebook in the team workspace and point `/api/analytics` at the `risk_map` table.
3. Activity replay over real sessions (the animated demo feature is already built for the sample data).
4. Multi-repository support via the adapter.
