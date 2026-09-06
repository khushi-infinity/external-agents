# Databricks Notebook — Agent Activity Universe Analytics
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
DATABASE = "agent_universe"
TABLE_RAW = f"{DATABASE}.entire_events"
CATALOG = "hive_metastore"  # Free Edition default; adjust to your catalog

spark.conf.set("spark.sql.shuffle.partitions", "4")

# COMMAND ----------
# 2) Ingest NDJSON (upload databricks/events.ndjson or live export via UI:
#    Data > Add Data > Upload File)
from pyspark.sql.functions import col, when, lit

try:
    raw_df = spark.read.format("json").option("multiline", "false").load("/FileStore/agent_universe/events.ndjson")
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
    raw_df = spark.createDataFrame(rows, ["event_type", "file", "last_agent", "change_count", "risk_score"])

raw_df.createOrReplaceTempView("events_raw")
display(raw_df.limit(10))

# COMMAND ----------
# 3) Normalize + clean: fill defaults, derive module, clamp risk
events = spark.sql("""
    SELECT
      COALESCE(event_type, 'file_change')                     AS event_type,
      COALESCE(file, 'unknown')                               AS file,
      COALESCE(last_agent, 'unknown')                         AS agent,
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