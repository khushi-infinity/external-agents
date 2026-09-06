# 🏁 BTW Buildathon 2026 — Submission Card & Checklist (copy-paste ready)

**Project title:** Agent Activity Universe — Freebuff × Entire
**Track:** E3 · Bring Entire to a New Agent or Workflow (+ Best Use of Databricks)
**One-liner:** An Entire external-agent plugin that makes the free coding agent **Freebuff** a first-class Entire agent (checkpoints, transcripts, graph), visualized as an interactive 3D **Agent Activity Universe** and scored for risk by **Databricks** — and hardened for the Noon Curveball: both transcript formats, unknown events never crash, incomplete transcripts yield partial results.

---

## 1. Repos — copy these into the submission form

| Field | Value |
|---|---|
| **GitHub fork (primary)** | `https://github.com/khushi-infinity/external-agents` |
| Branch | `main` |
| Final SHA | `fe5e454` |
| **Entire mirror** | `entire://aws-ap-south-1.entire.io/gh/khushi-infinity/external-agents` |
| Branch | `agent-activity-universe` |
| Final SHA | `ac2463a` |
| Mirror ID | `01M1THHNBW0K7KFTZGE9ETPD9K` (region `aws-ap-south-1`, India) |
| BUILDATHON.md | repo root (11 sections — full build story) |
| PROGRESS.md | repo root (status + evidence log) |
| Show guide | `SHOWING.md` (how to run + timed judge script) |

> Both branches contain identical content (the mirror is Entire's copy of the same repo). If the form accepts one URL only, use the **GitHub fork** and mention the mirror branch + SHA in the notes.

---

## 2. Milestone checkpoints (Entire)

| # | Milestone | Checkpoint ID | Commit |
|---|---|---|---|
| 1 | Initial understanding & intended architecture | `e6e841e2def8` | `ca06a8e` |
| 2 | Pre-noon stable state | `2adb77572071` | `10cb63c` |
| 3 | **Curveball response** — Freebuff plugin, dual-format transcripts | `614fa84595b2` | `dcfa571` |
| 4 | Final verification | `26d2aaa0ddf0` | `0619daa` |
| Bonus | Post-verification Databricks/docs update (auto-created by the commit hook from the live session — proves Freebuff→checkpoint flow) | `9d720fabaf1e` | `655a526` |

Verify on the mirror (`external-agents-mirror`):
```bash
entire checkpoint list
entire checkpoint explain 614fa84595b2    # → session fb-curveball-001, intent captured
entire status                              # "Checkpoints sync to: origin" (no pending)
```

**Checkpoint links (as the form usually expects them):** open the mirror `entire://…` repo in the Entire app and open each checkpoint, or paste the checkpoint IDs + commit SHAs above — each links to its commit. Representative explain output (curveball checkpoint `614fa84595b2`, session `fb-curveball-001`, author Khushi Sarawagi):

> Intent: *"Checkpoints are not added for Freebuff, so Entire does not store checkpoints done by Freebuff. Add them: build the external agent plugin … The noon curveball says the agent changed its format: support both the original transcript format and the new JSONL event format, unknown events must not crash, incomplete transcripts must produce partial results … Add tests for original format, new format, unknown events, and incomplete input."*

---

## 3. Databricks evidence (Best Use of Databricks)

**Workspace:** `https://dbc-e11b1b7c-7876.cloud.databricks.com` (org `7474646952123744`)

| Item | Value |
|---|---|
| Imported notebook | `/Workspace/Users/khushikhush006@gmail.com/Agent Activity Universe/ingest_and_score_uc` |
| Successful run 1 (fixed schema bug) | `#job/768048764575808/run/628954233613995` |
| Successful run 2 (with evidence persistence) | `#job/299946552182750/run/128947458769841` |
| Result tables | `workspace.agent_universe.risk_map` · `workspace.agent_universe.pipeline_run_summary` |

Evidence queries (both return SUCCESS):
```sql
SELECT * FROM workspace.agent_universe.pipeline_run_summary;
-- curveball_ingest | 18 lines | 17 known | 1 unknown skipped | 94.4% recognized
SELECT file, risk_score, change_count FROM workspace.agent_universe.risk_map ORDER BY risk_score DESC;
-- 14 rows: classic files + Curveball fixture files
-- src/checkout/apply_coupon.ts | 100 | 27    ← came from the NEW JSONL format
```
Repo files: `databricks/ingest_and_score.py` (FileStore) + `ingest_and_score_uc.py` (Unity Catalog) + `events.ndjson` + `events-new-format.ndjson` + `README.md` (live-run evidence §).

---

## 4. Demo (what judges run)

```bash
# 3D app (product surface)
cd agent-universe && npm install && npm run dev     # → http://localhost:5173
# no network? open agent-universe/demo.html directly (single-file offline demo)

# E3 proof
cd agents/entire-agent-freebuff && go test ./...    # 4 curveball groups green
cd .. && entire checkpoint list                     # checkpoints above
```

Demo script beats (full script in `SHOWING.md`): ① problem (Freebuff had no checkpoints) → ② 3D universe → ③ click a file node → ④ Risk View + Analytics (Databricks) → ⑤ Curveball fixture + tolerance proof → ⑥ live Freebuff → checkpoint.

---

## 5. Screenshots / recordings to attach (fallback evidence)
1. 3D universe — Risk View ON.
2. 📊 Analytics panel (agent performance, hotspots, failure patterns, velocity).
3. `entire checkpoint list` (the 5 checkpoints above).
4. Databricks run page (either SUCCESS run) + `pipeline_run_summary` row.
5. `go test ./...` green in `agents/entire-agent-freebuff`.

---

## 6. Final checklist (do before 3:00 PM IST)
- [ ] Fork pushed & clean: `github.com/khushi-infinity/external-agents` `main`
- [ ] Mirror pushed & clean: branch `agent-activity-universe`, checkpoints synced
- [ ] This card's SHAs match the live tips (they were updated after the final push)
- [ ] Screenshots saved (see §5)
- [ ] Submit form filled: fork URL + SHA, mirror URL + SHA, 5 checkpoint IDs/commits, BUILDATHON.md, demo access, Databricks workspace link
