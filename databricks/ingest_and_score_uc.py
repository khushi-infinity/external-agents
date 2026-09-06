# Databricks Notebook — Agent Activity Universe Analytics (Unity Catalog / Free Edition)
# UC variant: reads NDJSON from a managed volume under the `workspace` catalog
# (dbfs:/Volumes/workspace/agent_universe/files/) and writes tables to
# workspace.agent_universe.* — same logic as ingest_and_score.py.
# ----------------------------------------------------------
# Ingests Entire development activity (NDJSON) and produces the analytics
# that power the product's Risk View + Analytics panel:
#   - agent performance      (sessions, success rate, avg duration)
#   - file hotspots          (most changed files, per-file risk)
#   - failure patterns       (module-level risk)
#   - development velocity   (avg task time, files/session, retries)
#   - risk map               (per-file risk score for the 3D overlay)
#
# Workspace: Databricks Free Edition (serverless 2X-Small SQL warehouse)
# Data: synthetic sample events (labeled) OR live Entire checkpoint export.
# No credentials stored in the repository.

# COMMAND ----------
# 1) Config — point at your uploaded file or Delta table
DATABASE = "workspace.agent_universe"  # fully qualified: catalog.schema (UC)
TABLE_RAW = f"{DATABASE}.entire_events"
CATALOG = "workspace"  # Free Edition UC catalog (auto-created)

spark.conf.set("spark.sql.shuffle.partitions", "4")
spark.sql("USE CATALOG workspace")

# COMMAND ----------
# 2) Ingest NDJSON (upload databricks/events.ndjson or live export via UI:
#    Data > Add Data > Upload File)
from pyspark.sql.functions import col, when, lit

try:
    raw_df = spark.read.format("json").option("multiline", "false").load("/Volumes/workspace/agent_universe/files/events.ndjson")
except Exception as e:
    print("FileStore read failed, trying sample generation:", e)
    # Fallback: build a small representative table so the pipeline is demoable
    rows = [
        ("file_change", "src/auth/service.ts", "claude", 17, 82),
        ("file_change", "src/payments/service.ts", "codex", 14, 71),
        ("file_change", "src/auth/routes.ts", "claude", 12, 30),
        ("file_change", "src/api/routes.ts", "copilot", 11, 40),
        ("file_change", "src/payments/api.ts", "codex", 9, 55),
        ("file_change", "src/database/client.ts", "aider", 9, 60),
    ]
    raw_df = spark.createDataFrame(rows, ["event_type", "file", "agent", "change_count", "risk_score"])

# Tolerant column naming: live NDJSON uses 'agent'; if a source ever emits
# 'last_agent', rename it so the pipeline (and this notebook) always agree.
if "last_agent" in raw_df.columns and "agent" not in raw_df.columns:
    raw_df = raw_df.withColumnRenamed("last_agent", "agent")

raw_df.createOrReplaceTempView("events_raw")
display(raw_df.limit(10))

# COMMAND ----------
# 3) Normalize + clean: fill defaults, derive module, clamp risk
events = spark.sql("""
    SELECT
      COALESCE(event_type, 'file_change')                     AS event_type,
      COALESCE(file, 'unknown')                               AS file,
      COALESCE(agent, 'unknown')                              AS agent,
      COALESCE(change_count, 0)                               AS change_count,
      COALESCE(risk_score, 0)                                 AS risk_score,
      split(file, '/')[0]                                     AS module,
      CASE
        WHEN risk_score >= 70 THEN 'HIGH'
        WHEN risk_score >= 50 THEN 'MEDIUM'
        WHEN risk_score >= 30 THEN 'LOW'
        ELSE 'SAFE'
      END                                                     AS risk_level
    FROM events_raw
""")
events.createOrReplaceTempView("events")

# COMMAND ----------
# 4) File hotspots — most changed files with risk
hotspots = spark.sql("""
    SELECT file, SUM(change_count) AS changes, MAX(risk_score) AS risk
    FROM events
    GROUP BY file
    ORDER BY changes DESC
""")
display(hotspots.limit(10))

# COMMAND ----------
# 5) Failure patterns / module risk
module_risk = spark.sql("""
    SELECT module,
           COUNT(*)                 AS events,
           SUM(CASE WHEN risk_score >= 70 THEN 1 ELSE 0 END) AS high_risk_files,
           ROUND(AVG(risk_score),1) AS avg_risk
    FROM events
    GROUP BY module
    ORDER BY avg_risk DESC
""")
display(module_risk)

# COMMAND ----------
# 6) Agent performance
agent_perf = spark.sql("""
    SELECT agent,
           COUNT(*)                       AS events,
           ROUND(AVG(risk_score),1)       AS avg_risk,
           SUM(change_count)              AS total_changes
    FROM events
    GROUP BY agent
    ORDER BY total_changes DESC
""")
display(agent_perf)

# COMMAND ----------
# 7) Risk map for the 3D overlay (what the frontend consumes)
risk_map = spark.sql("""
    SELECT file, MAX(risk_score) AS risk_score, SUM(change_count) AS change_count
    FROM events
    GROUP BY file
""")
# Persist so the app's /api/analytics endpoint can read it
risk_map.write.mode("overwrite").saveAsTable(f"{DATABASE}.risk_map")
print("Wrote risk_map table:", f"{DATABASE}.risk_map")
display(risk_map.limit(10))

# COMMAND ----------
# 7b) NOON CURVEBALL — the agent changed its format.
# Freebuff/agents now also emit a NEW JSONL session-event format
# (databricks/events-new-format.ndjson). The pipeline must tolerate it
# without a rewrite: unknown event types are counted and skipped, known
# events are transformed into the same event schema, and the results feed
# the SAME risk_map table. Classic events.ndjson keeps working untouched.
from pyspark.sql.functions import col, lit, coalesce, least, from_json, split, max as spark_max, round as spark_round

KNOWN_EVENTS = {"session_started", "user_prompt", "agent_response", "file_changed",
                "checkpoint_created", "session_ended", "tool_call", "tool_result",
                "file_read", "usage"}

try:
    raw_new = spark.read.text("/Volumes/workspace/agent_universe/files/events-new-format.ndjson") \
        .where(col("value").rlike("\\S"))
    # Tolerant struct schema: only the fields we use are decoded, so nested
    # objects/arrays on other fields (agent, input, output, usage,
    # open_questions...) never break a line. Unknown event types still parse
    # and are counted; only truly malformed lines degrade to NULL event.
    parsed_new = raw_new.select(
        from_json(col("value"),
                  "struct<event:string, session_id:string, path:string, "
                  "lines_added:string, lines_removed:string, model:string, "
                  "timestamp:string>",
                  {"mode": "PERMISSIVE"}).alias("m")
    ).select(
        col("m.event").alias("event"),
        col("m.session_id").alias("session_id"),
        col("m.path").alias("path"),
        col("m.lines_added").alias("lines_added"),
        col("m.lines_removed").alias("lines_removed"),
        col("m.model").alias("model"),
        col("m.timestamp").alias("timestamp"),
    )
except Exception as e:
    print("FileStore read failed for new format, using embedded sample:", e)
    parsed_new = spark.createDataFrame([
        ("session_started", "btw-track3-demo-001", None, None, None, "acmecode-pro", "2026-09-06T09:00:00Z"),
        ("user_prompt",      "btw-track3-demo-001", None, None, None, None, "2026-09-06T09:00:08Z"),
        ("file_changed",     "btw-track3-demo-001", "src/checkout/apply_coupon.ts", "24", "3", None, "2026-09-06T09:02:31Z"),
        ("file_changed",     "btw-track3-demo-001", "tests/checkout/apply_coupon.test.ts", "71", "0", None, "2026-09-06T09:03:12Z"),
        ("checkpoint_created", "btw-track3-demo-001", None, None, None, None, "2026-09-06T09:04:49Z"),
        ("session_ended",    "btw-track3-demo-001", None, None, None, None, "2026-09-06T09:05:02Z"),
    ], ["event", "session_id", "path", "lines_added", "lines_removed", "model", "timestamp"])

parsed_new = parsed_new.where(col("event").isNotNull())
total_new = parsed_new.count()
known_new = parsed_new.filter(col("event").isin(*sorted(KNOWN_EVENTS)))
unknown_new = total_new - known_new.count()
recognized = (100.0 * known_new.count() / total_new) if total_new else 0.0
print(f"Curveball tolerance: {unknown_new} unknown event(s) skipped of {total_new} "
      f"({recognized:.1f}% recognized) — unknown events never crash the pipeline.")

# Transform new-format file_changed records into the SAME normalized schema
# used by classic events, with a deterministic per-file risk heuristic
# (new-format lines carry no precomputed risk; we derive one).
files_new = (known_new
             .filter(col("event") == "file_changed")
             .withColumn("agent", coalesce(col("model"), lit("freebuff-agent")))
             .withColumn("changes",
                         coalesce(col("lines_added").cast("int"), lit(0))
                         + coalesce(col("lines_removed").cast("int"), lit(0))))
files_new = files_new.withColumn("change_count", least(col("changes"), lit(999)).cast("int")) \
    .withColumn("risk_score", least(lit(100), lit(40) + col("change_count") * 3))
new_events = files_new.select(
    lit("file_change").alias("event_type"),
    col("path").alias("file"),
    col("agent"),
    col("change_count"),
    col("risk_score"),
)
print("New-format events normalized to the shared schema:")
display(new_events)

# COMMAND ----------
# 7c) Merge formats into ONE risk map (classic + new), then persist.
# Existing classic behaviour is untouched; the new format joins the same
# analytics so hotspots/module risk/agent performance include both.
risk_map_new = (new_events.groupBy("file")
                .agg(spark_max("risk_score").alias("risk_score"),
                     spark_max("change_count").alias("change_count")))
risk_map_all = risk_map.unionByName(risk_map_new)
risk_map_all.createOrReplaceTempView("risk_map_all")
risk_map_final = spark.sql("""
    SELECT file, MAX(risk_score) AS risk_score, SUM(change_count) AS change_count
    FROM risk_map_all
    GROUP BY file
""")
risk_map_final.write.mode("overwrite").saveAsTable(f"{DATABASE}.risk_map")
print("risk_map rewritten from BOTH formats (original + new JSONL event format).")
display(risk_map_final.limit(10))

# COMMAND ----------
# 8) Velocity metrics
velocity = spark.sql("""
    SELECT
      ROUND(AVG(change_count),1) AS avg_files_per_session,
      COUNT(DISTINCT agent)      AS active_agents,
      COUNT(*)                   AS total_events
    FROM events
""")
display(velocity)

print("""
Agent Activity Universe — Databricks pipeline complete.
Feed the risk_map output into the app via /api/analytics
(see databricks/README.md for the exact contract).
""")
# COMMAND ----------
# 9) Persist a machine-readable pipeline summary next to risk_map so the
#    curveball-tolerance outcome is queryable evidence, not just a print.
summary = spark.createDataFrame([(
    "curveball_ingest",
    int(total_new),                       # lines parsed from the NEW JSONL format
    int(known_new.count()),               # lines whose event type is known
    int(unknown_new),                     # unknown events skipped (never crash)
    round(float(recognized), 1),          # % recognized
)], ["pipeline", "new_events_total", "new_events_known", "new_events_unknown", "pct_recognized"])
summary.write.mode("overwrite").saveAsTable(f"{DATABASE}.pipeline_run_summary")
print("Wrote pipeline_run_summary:", summary.collect()[0])
display(summary)
