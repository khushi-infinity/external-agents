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
Databricks analytics (src/analytics/databricks.ts)  ← risk scoring, hotspots, velocity
    ↓
3D Universe (React Three Fiber)           ← files= nodes, deps= lines, agents= entities
    ↓
Developer / Judge
```

The `adapter → normalized data → product` separation makes the Noon Curveball cheap: "support a new agent" → add an agent adapter; "support multiple repositories" → extend the adapter; "work offline" → swap the data-loading layer.

## Entire Graph findings and verification
*(To be completed with live graph output from the mirror clone during the build.)*

- Graph search / definition lookup: TBD
- Relationship / impact analysis before a high-risk change: TBD
- Final semantic-diff analysis of the submitted implementation: TBD

Graph results are evidence, not an oracle — each finding will be verified against source code and tests before being recorded here.

## Noon Curveball: what changed and how we adapted
*(To be filled after 12:00 noon reveal.)*

- Constraint received: TBD
- Assumptions revisited: TBD
- What changed / what stayed intact: TBD
- Test proving the revised behavior: TBD

## Checkpoint links and what each checkpoint proves

**Repo:** GitHub fork `github.com/khushi-infinity/cli` · **Mirror:** `entire://aws-ap-south-1.entire.io/gh/khushi-infinity/cli` (India region) · **Branch:** `agent-activity-universe` · **Pushed SHA:** `3ef2f1c12`

Checkpoints are created when a commit happens during an ACTIVE agent session (hooks installed, Codex authed). Milestone commits must be made from inside a supported agent session (`codex` in the mirror clone), then verified with `entire checkpoint list`.

| Milestone | Commit | Checkpoint | What it proves |
|-----------|--------|------------|----------------|
| Initial understanding & intended architecture | `07520df06` | TBD | Design decisions, scope, data model, why E3 |
| Pre-noon stable state (11:45) | TBD | TBD | Runnable product before the Curveball |
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

**Role of Databricks:** the risk/analytics layer. Databricks analyzes accumulated development events (agent sessions, file changes, checkpoints, test outcomes) to produce agent performance, file hotspots, failure patterns, development velocity, and per-file risk scores that drive the **Risk View** overlay and the **Analytics** panel. Removing Databricks removes the risk layer — a core workflow of the product.

**Pipeline (designed):** Entire events → NDJSON export (`exportEventsNDJSON`) → Databricks (notebook/SQL warehouse) → scoring → `/api/analytics` → 3D view. The MVP ships precomputed analytics so the demo never blocks; the live ingestion job is the current next milestone.

**Data sources:** synthetic sample dataset (clearly labeled in `src/data/sample-data.ts`); no private, customer, or personal data used. Live Entire data is loaded from the repository's own checkpoint context when an export is provided. No credentials are stored in the repository.

## Known limitations and next steps
**Limitations:**
- MVP renders a synthetic sample dataset; the live Entire export path (`public/data/entire-export.json`) is wired and documented but needs the repository's actual `.entire/` output to be exercised.
- Databricks integration is a documented pipeline with precomputed analytics; the live ingestion → scoring job is the next milestone.
- 3D layout uses deterministic sample positions; production layout should derive from Entire Graph relationships.
- No Curlball-verified behavior yet — will be added at noon.

**Next steps:**
1. Wire live Entire checkpoint parsing into the adapter.
2. Stand up the Databricks ingestion + scoring job and point `/api/analytics` at it.
3. Activity replay over real sessions (the animated demo feature is already built for the sample data).
4. Multi-repository support via the adapter.