# BTW Buildathon 2026 — Submission Package

## Project: Agent Activity Universe — Freebuff × Entire

**Track:** E3 — Bring Entire to a New Agent or Workflow
**Optional Award:** Best Use of Databricks
**Team Lead:** Khushi Sarawagi
**Date:** 6 September 2026

---

## One-Line Summary

An Entire external-agent plugin that makes the free coding agent **Freebuff** a first-class Entire agent (checkpoints, transcripts, graph), visualized as an interactive 3D **Agent Activity Universe** and scored for risk by **Databricks**.

---

## Repositories

| Field | Value |
|---|---|
| GitHub Fork | https://github.com/khushi-infinity/external-agents |
| Fork Branch | main |
| Fork Final SHA | 7240cd4 |
| Entire Mirror | entire://aws-ap-south-1.entire.io/gh/khushi-infinity/external-agents |
| Mirror Branch | agent-activity-universe |
| Mirror Final SHA | 64303cf |
| Mirror ID | 01M1THHNBW0K7KFTZGE9ETPD9K |
| Mirror Region | aws-ap-south-1 (India) |

---

## Checkpoint Links

| # | Milestone | Checkpoint ID | Agent | Commit |
|---|---|---|---|---|
| 1 | Initial understanding | e6e841e2def8 | Claude Code | ca06a8e |
| 2 | Pre-noon stable | 2adb77572071 | Codex | 10cb63c |
| 3 | Curveball response (Freebuff plugin) | 614fa84595b2 | Freebuff | dcfa571 |
| 4 | Final verification | 26d2aaa0ddf0 | Freebuff | 0619daa |
| 5 | Post-verification docs update (auto-created by commit hook) | 9d720fabaf1e | Freebuff | 655a526 |

Verify: run `entire checkpoint list` in the mirror clone. Export: `entire checkpoint explain <id> --json`.

---

## Databricks Evidence

| Item | Value |
|---|---|
| Workspace | https://dbc-e11b1b7c-7876.cloud.databricks.com (org 7474646952123744) |
| Notebook | /Workspace/Users/khushikhush006@gmail.com/Agent Activity Universe/ingest_and_score_uc |
| Successful Run 1 | https://dbc-e11b1b7c-7876.cloud.databricks.com/?o=7474646952123744#job/768048764575808/run/628954233613995 |
| Successful Run 2 | https://dbc-e11b1b7c-7876.cloud.databricks.com/?o=7474646952123744#job/299946552182750/run/128947458769841 |
| Tables | workspace.agent_universe.risk_map (14 rows) |
| | workspace.agent_universe.pipeline_run_summary (tolerance evidence) |

Evidence SQL:

```sql
SELECT * FROM workspace.agent_universe.pipeline_run_summary;
-- Result: curveball_ingest | 18 lines | 17 known | 1 unknown skipped | 94.4% recognized

SELECT file, risk_score, change_count FROM workspace.agent_universe.risk_map ORDER BY risk_score DESC;
-- Result: 14 rows — classic files + Curveball fixture files merged
```

---

## Demo Instructions

### Run the 3D App

```bash
cd ~/Desktop/buildathon/agent-universe
npm run dev
# Opens http://localhost:5173
```

Offline version: double-click `agent-universe/demo.html` (no server, no network).

### Run the Plugin Tests

```bash
cd ~/Desktop/buildathon/external-agents-mirror/agents/entire-agent-freebuff
go test ./...
# All green: original + new format + unknown events + incomplete input
```

### Show the Checkpoints

```bash
cd ~/Desktop/buildathon/external-agents-mirror
entire checkpoint list
# 60 checkpoints — 5 milestones listed above
```

### Judge Demo Script (8-10 minutes)

1. **Problem (30s):** "Freebuff sessions had zero Entire checkpoints. We fixed that."
2. **3D Universe (1min):** Click a file node → agent, checkpoint, intent shown.
3. **Risk View (1min):** Toggle ON → files recolor from Databricks risk map.
4. **Analytics (1min):** Open panel → agent performance, hotspots, velocity.
5. **Curveball (2min):** Show `events-new-format.ndjson` → tolerance summary: 18/17/1/94.4%.
6. **Live checkpoint (1min):** `entire checkpoint list` → Freebuff-backed milestones.

---

## What the Project Does

### The Problem
AI coding agents produce work, but teams have no structured view of what agents are doing, why, and where the work creates risk. Freebuff — a free agent used in this buildathon — had no Entire support at all.

### The Solution
Three interconnected layers:

**Capture** — The `entire-agent-freebuff` plugin makes Freebuff a first-class Entire agent. Sessions, hooks, transcripts, and checkpoints all flow through Entire.

**Analyze** — Databricks ingests development events, computes risk scores per file and per agent, and persists a `risk_map` table.

**Explore** — The Agent Activity Universe (React Three Fiber) renders checkpoints as a navigable 3D space. Risk View recolors from Databricks. Analytics panel shows metrics.

### The Curveball
At noon, Freebuff changed its data format (JSONL event stream). Our parser reads both formats, unknown events are counted and skipped, truncated transcripts return partial results. Databricks ingested the new format live — 94.4% recognition rate.

---

## Architecture

```
Freebuff Sessions → entire-agent-freebuff (plugin) → Entire Checkpoints
                                   ↓
                    Databricks (risk scoring) → risk_map
                                   ↓
                    Agent Activity Universe (3D) → Risk View + Analytics
```

### Component Map

| Component | Location | Purpose |
|---|---|---|
| Freebuff Plugin | agents/entire-agent-freebuff/ | Entire integration for Freebuff |
| 3D App | agent-universe/ | Interactive universe visualization |
| Databricks Notebook | databricks/ingest_and_score_uc.py | Risk scoring pipeline |
| Sample Events | databricks/events.ndjson | Classic format dataset |
| Curveball Fixture | databricks/events-new-format.ndjson | New JSONL format + unknown event |
| Demo Guide | SHOWING.md | How to run + judge script |
| Full Docs | BUILDATHON.md | 11-section project story |
| Status | PROGRESS.md | Verified status log |
| This File | SUBMISSION.md | Copy-paste submission card |

---

## Tech Stack

| Layer | Technology |
|---|---|
| External Agent Plugin | Go 1.22+, Entire External Agent Protocol |
| 3D Frontend | React 19, TypeScript, Vite, React Three Fiber, Zustand |
| Analytics | Databricks (PySpark), Unity Catalog |
| Testing | Vitest (7 tests), Go test (plugin) |
| Tooling | mise (Go builds), npm (frontend) |

---

## Key Numbers

| Metric | Value |
|---|---|
| Checkpoints created | 60 (5 milestones) |
| Plugin tests | All green (4 curveball groups) |
| App tests | 7/7 pass |
| Databricks risk_map rows | 14 |
| Curveball tolerance | 94.4% (18 lines, 1 known, 1 unknown) |
| Agents visualized | 4 (Claude, Codex, Copilot, Aider) |
| File nodes | 14 |
| Module filters | 6 |

---

## Screenshots to Attach

1. 3D universe with Risk View ON (files recolored by risk)
2. Analytics panel (agent performance, hotspots, velocity)
3. `entire checkpoint list` output (5 milestones visible)
4. Databricks run page (SUCCESS) + pipeline_summary row
5. `go test ./...` green output

---

## Checklist Before Submitting

- [ ] Fork pushed: github.com/khushi-infinity/external-agents, main, 7240cd4
- [ ] Mirror pushed: agent-activity-universe, 64303cf
- [ ] Entire status: checkpoints synced to origin
- [ ] Screenshots saved (5 items above)
- [ ] Submit form filled: repos, SHAs, checkpoints, BUILDATHON.md, demo access

---

Generated for BTW Buildathon 2026 · Track E3
