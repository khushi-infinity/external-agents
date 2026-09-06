# Freebuff x Entire — External Agent Integration (E3)

## One-sentence summary
An Entire external-agent plugin (`entire-agent-freebuff`) that brings Entire checkpoints, transcripts and the Entire Graph to the Freebuff coding agent — capturing what Freebuff sessions do and why — plus a 3D Agent Activity Universe and Databricks risk analytics that make the captured checkpoint context explorable.

## Problem, intended user and why it matters
**User:** a developer (or team) doing serious work with Freebuff — a free AI coding agent used in this very buildathon instead of paid agents like Claude or Codex.

**Problem:** Entire's checkpoint superpowers only work for agents it natively supports (Claude Code, Codex, OpenCode, …). Freebuff — a genuinely popular free agent built on the manicode engine — produced **no checkpoints**: commits made by Freebuff sessions were invisible to Entire, so the *why* behind every change (intent, prompts, files touched, failures) was lost. The only workaround was manually switching to Codex to make milestone commits.

**Why it matters:** teams that standardize on free agents deserve the same safety net as paid-agent teams: checkpoints that preserve intent, rewind, and a resume path. This plugin makes Freebuff a first-class Entire agent so free-agent work stops being a blind spot.

## Selected Entire track and why Entire is essential
**Track E3 — Bring Entire to a New Agent or Workflow.**

We extend Entire to a **new coding agent** (Freebuff) through the external-agent protocol. Entire is essential by construction: the deliverable *is* an Entire integration — session capture, hook lifecycle, and checkpoint writing all run through Entire CLI. The companion 3D "Agent Activity Universe" (built in the same designated fork) consumes the checkpoint/session data this integration captures, and the Databricks layer scores it; neither works without Entire checkpoints existing in the first place.

## Architecture and main workflow

```
Freebuff (manicode engine) sessions
   ~/.config/manicode/projects/<project>/chats/<session-id>/
     chat-messages.json  (original format)   session.jsonl (new format)
        ↓  entire-agent-freebuff (external agent plugin, this repo)
   hook lifecycle:  entire hooks freebuff session-start|prompt-submit|stop|session-end
   transcript analysis: prompts, files changed, summaries — BOTH formats
        ↓  Entire CLI + Checkpoints + Graph
   checkpoints capture Freebuff intent per commit (verified live)
        ↓  adapter (agent-universe/src/data/adapter.ts)
   Agent Activity Universe (3D) + Databricks risk analytics
```

The plugin lives at `agents/entire-agent-freebuff/` with the protocol
surface in `internal/protocol` and Freebuff logic in `internal/freebuff`
(session layout, hooks, and a single dual-format transcript parser).

## What the project is doing — the full picture

This repository is the **E3 submission**, and it is one end-to-end system that makes agent work visible and explainable:

**1. Capture — Entire checkpoints for a previously unsupported agent.**
The product starts where agent work happens: in a coding-agent session. Freebuff (the free agent used to build this very project) had **no Entire support**, so its commits produced no checkpoints and the *why* behind them was lost. `agents/entire-agent-freebuff/` closes that gap through the documented External Agent protocol (see `AGENT.md`): `entire enable --agent freebuff` installs a 4-hook lifecycle registry (`.freebuff/entire-hooks.json`), and the plugin maps Freebuff's on-disk sessions (`~/.config/manicode/projects/<project>/chats/<sid>/`) into Entire sessions, transcripts and — on commit — **checkpoints with intent, prompts, files changed and open questions**. This was verified live today: the milestone commits of this very build were checkpointed from real Freebuff sessions (`614fa84595b2`, `5531a5f11fa4`, `26d2aaa0ddf0`).

**2. Analyze — Databricks risk scoring.**
`databricks/` is the analytics brain. The notebook ingests development events (classic `events.ndjson` and — after the Noon Curveball — the new JSONL session-event format `events-new-format.ndjson`), and computes agent performance, file hotspots, module failure patterns, velocity and a **per-file `risk_map`**. Risk bands (SAFE/LOW/MEDIUM/HIGH) power the product's Risk View. Two serverless runs finished `SUCCESS` today and persisted `workspace.agent_universe.risk_map` + `pipeline_run_summary` (evidence in §9).

**3. Explore — the Agent Activity Universe (the demo surface).**
`agent-universe/` is a React Three Fiber app that turns checkpoint context into a navigable 3D space: files are nodes, dependencies and activity are connections, and agents (◆ Claude, ● Codex, ▲ Copilot, ■ Aider) float near the files they touched with live `N files · M calls` tags. Clicking a node opens an inspector (agent, checkpoint intent, prompt, files changed, risk score); filters narrow by agent or module; **Activity Replay** animates the session timeline; **Risk View** recolors the universe from the Databricks risk map; and the **Analytics** panel shows agent performance, hotspots, failure patterns and velocity from the same `risk_map`. Its data adapter (`src/data/adapter.ts`) tries a live Entire export → `/api/analytics` → and finally the labelled sample dataset, so the demo never blocks.

**Why the whole is greater than the parts:** the 3D universe and the Databricks scoring both *depend on Entire checkpoints existing in the first place* — which, for Freebuff, only became true because of the plugin in section 3. Capture (E3 deliverable) → analyze (Databricks) → explore (universe) is a closed loop: the tool we built produced the very checkpoints this project now visualizes.

## Entire Graph findings and verification

Graph impact analysis ran **before** the Curveball implementation (required step). Live output, repo commit `2f47e69`:

- `entire graph impact --symbol ParseHook --repo .` → the lifecycle-handler entry point exists once per agent module (`agents/entire-agent-amp/.../hooks.go`, `goose`, `grok`, …) plus the shared `hookParser.ParseHook` interface in each `internal/protocol` package — the exact surface a new agent must implement.
- `entire graph search --query "external agent jsonl transcript parser extract modified files prompts" --repo .` → ranked `Agent.ExtractModifiedFiles` (kiro `transcript.go:718`), `decodeTranscript` (kilo `session_jsonl.go:40` — the existing JSONL-decode precedent), `modifiedFilesFromMessages` (kilo), `modifiedFiles` (omp) — proving every agent's transcript analyzer funnels through a single `parseTranscript`-style parser, which is the code path the new-format Curveball affects. Verify commands were suggested per hit (`go test ./internal/kiro`).
- Graph output is treated as evidence, not an oracle: every finding above was verified against the source and the new agent's own tests before recording.
- Final semantic diff of the Curveball response: `entire graph diff --base 2f47e69 --head dcfa571 --repo .` → reports the new `.freebuff/entire-hooks.json` registry sections, the rewritten BUILDATHON.md/PROGRESS.md sections, and the agent-universe doc churn, each tagged `(0 dependents)` where nothing else references them — the expected isolated blast radius of an additive external-agent plugin. (Graph results are verified against source and tests before recording.)

## Noon Curveball: what changed and how we adapted

**Constraint received (Track 3, 12:00):** "The agent changed its format." The agent/workflow we integrate (Freebuff) released a **new transcript and lifecycle event format** (a JSONL event stream). The integration must support **both** the original and the new format, unknown events must **never crash** it, an **incomplete transcript** must produce a **partial result** (never a corrupted or discarded session), and existing Checkpoint behaviour must stay compatible. A JSONL fixture representing the new format was attached.

**Assumption that changed:** we assumed Freebuff transcripts are *only* the original `chat-messages.json` array (messages with `variant: user|ai` and `blocks[]`). The new format is a line-oriented JSONL event stream (`session_started`, `user_prompt`, `agent_response`, `file_changed`, `checkpoint_created`, `usage`, `session_ended`, …) that shares no shape with the array.

**How the design changed:** instead of writing a second parser path that duplicates the pipeline, `internal/freebuff/transcript.go` now sniffs the content and normalizes **both** formats into one `turn` model before the analyzer stage:
- Format detection is content-based (leading `[` vs `{`), so every analyzer entry point (`extract-*`, `get-transcript-position`) works unchanged.
- New-format events map onto the same turns: `user_prompt` → prompt, `file_changed` → files, `checkpoint_created.summary` → summary (preferred over assistant text).
- Unknown JSONL events are skipped and counted (`unknownEvents`, flagged `partial`) — never fatal; known-but-inert events (`tool_call`, `tool_result`, `file_read`, `usage`) are recognized as no-ops so they are not mistaken for unknown ones.
- Incomplete transcripts degrade safely: a truncated JSONL stream or truncated chat array is salvaged to the longest valid prefix and flagged `partial`; an empty file yields an empty session, not an error.
- Unknown lifecycle hook names and malformed hook payloads parse to "ignore" in `hooks.go`, keeping older Checkpoint behaviour compatible.

**Why the new result is safe:** the four Curveball-mandated test cases are automated in `internal/freebuff/transcript_test.go` (original format, new format, unknown events, incomplete input), all passing; the committed fixture `testdata/track-3-agent-session.jsonl` is the exact attached card, and the binary was exercised against it live through the CLI (`extract-prompts`, `extract-modified-files`, `extract-summary` → correct outputs, and truncated input → partial results).

## Checkpoint links and what each checkpoint proves

**Repo:** fork `github.com/khushi-infinity/external-agents` · **Mirror:** `entire://aws-ap-south-1.entire.io/gh/khushi-infinity/external-agents` (India, Mirror ID `01M1THHNBW0K7KFTZGE9ETPD9K`) · implementation branch listed in PROGRESS.md.

| Milestone | Checkpoint | What it proves |
|-----------|-----------|----------------|
| Initial understanding & intended architecture | `e6e841e2def8` (commit `ca06a8e`) | Build plan: E3 track, Entire-central architecture, PDF scope |
| Pre-noon stable state (11:45) | `2adb77572071` (commit `10cb63c`) | Runnable product before the Curveball, intent/architecture/risks recorded |
| **Curveball response — Freebuff plugin (12:00+)** | `614fa84595b2` (commit `dcfa571`) | Freebuff external agent added to Entire and enabled (`entire enable --agent freebuff` → 4 hooks); dual-format transcript support; the four Curveball tests pass; checkpoint created by attaching the live **Freebuff session** `fb-curveball-001` (`entire session attach --agent freebuff fb-curveball-001` → “Created checkpoint 614fa84595b2”) — proof Freebuff work now produces Entire checkpoints |
| Final implementation & verification | `26d2aaa0ddf0` (commit `0619daa`, tip of the pushed branch) — preceded by `5531a5f11fa4` (commit `02f6095`) | Final state: both final docs commits were checkpointed **automatically from the active Freebuff session** (`fb-curveball-001`) by the Entire commit hook — Freebuff sessions now produce Entire checkpoints exactly like natively supported agents; tests green, graph evidence recorded, pushed |
| Post-verification Databricks + docs update (13:26) | `9d720fabaf1e` (commit `655a526`, pushed to `agent-activity-universe`) | Live Databricks run evidence recorded (see §9) after the live run caught and fixed two real notebook bugs; auto-checkpointed from the active session by the Entire commit hook — the same commit → checkpoint flow shown repeatedly |

## Setup, run and test instructions

```bash
# 1. Build + install the Freebuff external agent
cd agents/entire-agent-freebuff
mise run build
cp entire-agent-freebuff ~/.local/bin/

# 2. Enable in a repository that Freebuff works in
entire enable --agent freebuff --local --telemetry=false   # installs hook registry
entire agent list                                            # Freebuff discoverable

# 3. Make a Freebuff change and checkpoint it
#    (hooks fire session lifecycle; a commit then creates the checkpoint)
entire hooks freebuff session-start                          # or let the wrapper fire
git commit -am "my freebuff change"
entire checkpoint list

# Manual fallback (attach a Freebuff chat after a commit):
entire session attach --agent freebuff <chat-session-id>

# 4. Tests (Curveball-critical)
cd agents/entire-agent-freebuff
go test ./...          # original + new format + unknown events + incomplete input

# 5. Agent Activity Universe companion app
cd agent-universe && npm install && npm test && npm run build
```

**Offline demo:** open `agent-universe/demo.html`.
**Full judge runbook (open, run, what to show, fallbacks):** `SHOWING.md`.

## Databricks use, data sources and limitations

Databricks is the risk/analytics layer of the companion universe: it ingests development events (from Entire checkpoints) and produces agent performance, file hotspots, failure patterns, velocity and a per-file risk map that drive the **Risk View** overlay and **Analytics** panel. Removing Databricks removes the risk-scoring layer (award criteria: "must power a meaningful part of the workflow").

**Pipeline (`databricks/`):** `ingest_and_score.py` ingests NDJSON → auto-creates `agent_universe` catalog/schema → computes analytics → writes `agent_universe.risk_map`. Notebook cells 7b–7c are the **Curveball adaptation**: the new JSONL session-event format (`databricks/events-new-format.ndjson`, the attached Track 3 fixture plus a deliberately unknown `model_switched` event) is parsed with unknown-event tolerance (skipped + counted, never crash), normalized into the shared event schema with a derived risk heuristic, and merged into the **same** `risk_map` table — classic `events.ndjson` behaviour untouched.

**Data sources:** synthetic sample events (labeled in `agent-universe/src/data/sample-data.ts`) and the Curveball fixture; no personal/customer data; no credentials in the repo. Transformations are traceable in the notebook.

**Live run — DONE (6 Sep, serverless).** The workspace has no DBFS/FileStore (Unity Catalog), so the run used `databricks/ingest_and_score_uc.py` (imported to the workspace, fully qualified `workspace.agent_universe`) with both NDJSON files uploaded to a UC volume. Two runs finished `SUCCESS` (links in SHOWING.md §2c). Persisted results — `workspace.agent_universe.risk_map` (14 rows: classic files + Curveball fixture files) and `pipeline_run_summary` (`curveball_ingest | 18 lines | 17 known | 1 unknown skipped | 94.4% recognized`). Two real bugs were caught and fixed during the live run: the notebook referenced a `last_agent` column the live NDJSON never contains (now normalised to `agent` at ingest) and the new-format JSON decode used a `map<string,string>` schema that nested objects would break (now a tolerant `struct` schema). Full evidence and queries: `databricks/README.md` §Live run evidence and `SHOWING.md` §2c/§3 Beat 5.

## Known limitations and next steps
**Limitations:**
- Freebuff's engine does not yet expose a user-configurable hook registry; the plugin installs a declarative registry (`.freebuff/entire-hooks.json`) and implements `parse-hook` for every declared event. Firing those commands from the engine automatically is the Freebuff-side step (hooks can be driven manually today — verified live).
- Lifecycle e2e scenarios need an interactive, logged-in Freebuff session (`FREEBDUFF_E2E=1 E2E_AGENT=freebuff`); the Freebuff CLI has no headless prompt flag yet.
- Chat start times derive from chat-dir mtime (the engine stores display-only timestamps).
- 3D layout of the companion universe uses deterministic sample positions.

**Next steps:**
1. Freebuff engine hook wiring (plugin auto-loads `.freebuff/entire-hooks.json` registry on session start).
2. Run the Databricks notebook in the workspace and point the app's `/api/analytics` at `risk_map` live.
3. Live Entire checkpoint parsing into the universe adapter (real `.entire/` export).
4. Headless `freebuff -p` support for fully automated lifecycle e2e.
