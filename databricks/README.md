# Databricks — Agent Activity Universe Analytics

This folder contains the Databricks side of the product. Databricks is **essential** to the core workflow: it analyzes accumulated development events (from Entire checkpoints) and produces the risk scoring + analytics that drive the product's **Risk View** and **Analytics** panel. Removing Databricks removes a core workflow of the product.

## Pipeline

```
Entire checkpoints/sessions
        ↓  (app exports NDJSON via exportEventsNDJSON / databricks/events.ndjson)
Development event dataset (JSON/NDJSON) — classic OR new JSONL event format
        ↓  (this notebook: ingest → clean → aggregate, tolerates both formats)
Databricks analytics: agent performance, file hotspots,
failure patterns, velocity, per-file risk map
        ↓
App /api/analytics → 3D Risk View + Analytics panel
```

### Noon Curveball — the agent changed its format

Agents now also emit a **new JSONL session-event format** (`session_started`, `user_prompt`, `file_changed`, `checkpoint_created`, `session_ended`, …). The notebook adapts without a rewrite:

- Upload **both** `databricks/events.ndjson` (original) and `databricks/events-new-format.ndjson` (new format; the Track 3 Curveball fixture plus a deliberately unknown `model_switched` event).
- Cells 7b–7c parse the new format, **count and skip unknown event types** (never crash), normalize `file_changed` records into the same event schema with a derived risk heuristic, and merge both formats into the **same** `agent_universe.risk_map` table.
- Classic behaviour is untouched: cells 1–7 run exactly as before on original-format data.
- The run prints the tolerance evidence (unknown events skipped / recognized %) so the adapted workflow is verifiable.

## Setup (5 minutes, in your Databricks workspace)

1. **Sign in** to Databricks (Free Edition). Create a workspace (one per team; the deployment owner must be able to sign in before judging).
2. **Create a serverless SQL warehouse** (Free Edition: one 2X-Small). Start it.
3. **Import the notebook:**
   - Workspace → Create → Import → upload `databricks/ingest_and_score.py`
4. **Upload the events data:**
   - Option A (recommended): from the app, export real events:
     ```js
     // app: window.__exportEvents() → databricks/events.ndjson
     // or use the sample at databricks/events.ndjson
     ```
   - Data → Add Data → upload `databricks/events.ndjson` **and** `databricks/events-new-format.ndjson` to `/FileStore/agent_universe/`
5. **Run all cells** in the notebook (serverless compute). It writes the `agent_universe.risk_map` table.
6. **Verify:** the notebook displays file hotspots, module risk, agent performance, velocity, and the risk map.

## The analytics contract (what the frontend consumes)

`GET /api/analytics` should return (or the app falls back to sample analytics):

```json
{
  "agent_performance": { "claude": { "sessions": 42, "success_rate": 0.91, "avg_duration": 23 }, "..." : { } },
  "file_hotspots": [ { "path": "src/auth/service.ts", "changes": 17, "risk": 82 } ],
  "failure_patterns": [ { "module": "payments", "failures": 7, "risk": 71 } ],
  "velocity": { "avg_task_minutes": 23, "avg_files_per_session": 8.4, "avg_retries": 1.8 },
  "risk_map": { "src/auth/service.ts": 82 }
}
```

The app's `src/analytics/databricks.ts` handles this: it tries `fetch('/api/analytics')` first, falls back to precomputed analytics so the demo never blocks.

## Data provenance & responsible use

- Data is **synthetic sample data** (clearly labeled in `src/data/sample-data.ts`) unless you upload a real Entire checkpoint export.
- No personal data, customer data, or credentials are uploaded. Keep credentials out of notebooks and the repository.
- Transformations are traceable: the notebook shows every aggregate (hotspots, module risk, agent perf, velocity, risk map).

## Free Edition constraints (designed around)

- One serverless 2X-Small SQL warehouse; up to 5 concurrent job tasks.
- One active pipeline per type; up to 3 Apps; one workspace/metastore per account.
- Quota exhaustion can block compute for the rest of the day → **run the notebook early**, save outputs, keep the local fallback demo ready.
- Keep a screenshot/recording of the notebook results as fallback evidence.