# 🎬 Agent Activity Universe — How to Open, Run & Show It

**Project:** Freebuff × Entire (Track E3) + Agent Activity Universe + Databricks risk analytics.
**Deadline:** submit by 3:00 PM IST. This guide covers (1) running everything, (2) a timed judge demo script, (3) the evidence to show, and (4) the offline fallback.

---

## 1. What you are showing (the 60-second story)

> “Every AI coding agent in our repo — including **Freebuff**, a free agent Entire didn’t support — now produces **Entire checkpoints** that capture *what* it changed and *why*. We visualize that captured activity as an interactive **3D Agent Activity Universe**, and **Databricks** scores each file and agent for **risk**. Watch: toggling Risk View recolors the universe from real analytics.”

Three beats to hit:
1. **Entire integration (the E3 deliverable)** — Freebuff sessions now produce Entire checkpoints (4 verified checkpoints, hooks live).
2. **The product surface** — 3D universe: agents float near the files they touched; click anything for checkpoints, prompts, risk.
3. **Databricks (analytics layer)** — risk scores/hotspots come from the live `risk_map` pipeline that tolerates the Noon Curveball’s new JSONL event format.

---

## 2. Open & run everything

### 2a. The 3D app (primary demo surface)
```bash
cd agent-universe
npm install          # first time only
npm test             # 7 tests — expect all passing
npm run dev          # → open http://localhost:5173
```
Production/static version: `npm run build && npm run preview`, or simply open the **single-file offline demo** `agent-universe/demo.html` in any browser (no server, no network needed — it falls back to the built-in sample dataset, shown in the console as `[adapter] Using sample dataset for demo`).

### 2b. The Entire ↔ Freebuff integration (the E3 proof)
```bash
# inside the submitted repo (external-agents-mirror / external-agents)
cd agents/entire-agent-freebuff
mise run build && cp entire-agent-freebuff ~/.local/bin/
go test ./...        # Curveball-critical: original + new format + unknown events + incomplete input

# in the repo you demo in:
entire status                          # shows enabled + sessions
entire checkpoint list                 # the 4 milestone checkpoints
entire hooks freebuff session-start    # fire a Freebuff lifecycle hook manually
entire session attach --agent freebuff <chat-session-id>   # manual checkpoint fallback
```

### 2c. Databricks (analytics evidence)
Already **run live** today on workspace `https://dbc-e11b1b7c-7876.cloud.databricks.com` (org `7474646952123744`):

| Item | Location |
|---|---|
| Imported notebook | `/Workspace/Users/khushikhush006@gmail.com/Agent Activity Universe/ingest_and_score_uc` |
| Successful runs | `#job/768048764575808/run/628954233613995` and `#job/299946552182750/run/128947458769841` (both `SUCCESS`) |
| Data files | UC volume `dbfs:/Volumes/workspace/agent_universe/files/` → `events.ndjson`, `events-new-format.ndjson` |
| Result tables | `workspace.agent_universe.risk_map` (14 rows), `workspace.agent_universe.pipeline_run_summary` |

Re-run from scratch (new workspace): import `databricks/ingest_and_score.py` (FileStore) or `databricks/ingest_and_score_uc.py` (Unity Catalog), upload both NDJSON files, run all cells, and open the two tables.

---

## 3. Judge demo script (~8–10 minutes, on `npm run dev` or `demo.html`)

### Beat 1 — Land the problem (30 s, no screen needed)
“Entire checkpoints only existed for agents Entire supports natively. Our team used **Freebuff**, a free agent — so every Freebuff commit was invisible: no intent, no prompts, no checkpoint, no way to rewind. That is the blind spot we removed.”

### Beat 2 — The universe appears (1 min)
Show the 3D scene. Point out:
- **Agent markers** (◆ Claude orange, ● Codex green, ▲ Copilot purple, ■ Aider pink) each tagged `N files · M calls` — agents float near the files they touched.
- File nodes + dotted dependency/activity lines = the **Entire Graph** visualization.
- Bottom pills filter by **AGENT** and **MODULE** — click one (e.g. `claude`, then `src/auth`) and the scene filters live.

### Beat 3 — Inspect one unit of work (1.5 min)
Click a file node → inspector panel shows agent, checkpoint, prompt/intent, files changed, risk score.
**Say:** “Every card here came out of an Entire checkpoint — this is the *why* behind the change, captured at commit time.”

### Beat 4 — Risk View = Databricks (1.5 min)
Click **Risk View** (top-right) → the button becomes `⦿ Risk View ON` and files glow by risk (red high / amber medium / green low). Then open **📊 Analytics** and scroll the four panels:
- Agent Performance (sessions · avg time · success %)
- File Hotspots (changes · risk)
- Failure Patterns by module
- Development Velocity + Modules
**Say:** “These numbers come from the Databricks `risk_map` pipeline. The notebook also proves it can survive agents changing their data format — that is the Noon Curveball we had to pass.”

### Beat 5 — The Curveball / Freebuff checkpoint proof (2–3 min)
This is the E3 moment. Run (or show pre-captured output of):
```bash
cd agents/entire-agent-freebuff && go test ./...        # all green
entire checkpoint list                                  # e6e841e2def8 → 2adb77572071 → 614fa84595b2 → 26d2aaa0ddf0
```
Open `databricks/events-new-format.ndjson` in an editor — it is the **exact Curveball fixture** (the `session_started … session_ended` JSONL events) plus one deliberately **unknown** `model_switched` line. Then show the persisted tolerance result:
```sql
SELECT * FROM workspace.agent_universe.pipeline_run_summary;
-- curveball_ingest | 18 | 17 | 1 | 94.4   (1 unknown skipped, never crashed)
SELECT file, risk_score, change_count FROM workspace.agent_universe.risk_map ORDER BY risk_score DESC;
```
**Say:** “The agent changed its format at noon. Our parser reads **both** formats, unknown events are counted and skipped, truncated transcripts return partial results — four automated test groups prove it, and Databricks ingested the new format into the same risk map live.”

### Beat 6 — Freebuff → checkpoint, live (1–2 min, optional but impressive)
```bash
entire hooks freebuff session-start     # Freebuff session appears in `entire status`
git commit -am "demo change made by Freebuff"
entire checkpoint list                  # new checkpoint auto-created from the Freebuff session
```
**Say:** “That commit just became an Entire checkpoint because Freebuff is now a first-class agent — the exact thing that was impossible this morning.”

---

## 4. Screenshots / recordings to keep as fallback evidence
1. 3D universe (Risk View ON) — the app on screen.
2. Analytics panel (all four sections).
3. `entire checkpoint list` output (the four milestone checkpoints).
4. Databricks run page (either `SUCCESS` run link above) + `pipeline_run_summary` row.
5. `go test ./...` green output for the Freebuff plugin.

---

## 5. If anything breaks (offline fallback)
- **No node / no network?** Open `agent-universe/demo.html` directly — single file, self-contained, uses sample data (console: `[adapter] Using sample dataset for demo`). Two 404s for `/data/…` and `/api/repository` in the console are the *intended* fallback chain, not errors.
- **No Databricks UI?** Show the persisted tables via the CLI:
  ```bash
  databricks tables list workspace agent_universe
  # risk_map, pipeline_run_summary
  ```
- **Freebuff hooks not firing?** The plugin writes `.freebuff/entire-hooks.json`; drive hooks manually with `entire hooks freebuff …` (verified working) — engine auto-fire is a documented next step, not a blocker for the demo.

---

## 6. Submission checklist (before 3:00 PM IST)
- [ ] Fork pushed: `github.com/khushi-infinity/external-agents` (main) — final SHA recorded.
- [ ] Mirror pushed: `entire://aws-ap-south-1.entire.io/gh/khushi-infinity/external-agents` (branch `agent-activity-universe`) — checkpoints synced to origin (`entire status` shows no pending).
- [ ] BUILDATHON.md + PROGRESS.md final, in repo root and this workspace top level.
- [ ] Screenshots saved (see §4) + this guide available next to the demo.
- [ ] Submit form: fork URL + SHA, mirror URL, 4 checkpoint links (`e6e841e2def8`, `2adb77572071`, `614fa84595b2`, `26d2aaa0ddf0`), BUILDATHON.md, demo access.
