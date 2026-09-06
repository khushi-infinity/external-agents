# entire-agent-freebuff

External agent binary that adds **Freebuff** support to Entire CLI.

Freebuff is the free AI coding agent built on the manicode engine (CLI +
desktop). This adapter teaches Entire to:

- **Create checkpoints** for work done in Freebuff sessions, capturing what
  the agent did and why.
- **Analyze Freebuff transcripts** — user prompts, modified files, and
  summaries — in **two transcript formats** (see Curveball below).
- **Install hooks** so Freebuff lifecycle events (session start, prompt
  submit, stop, session end) flow through Entire via `entire hooks freebuff`.

## Requirements

- `entire` on `PATH` with external-agent discovery enabled
  (`external_agents: true` in the untracked `.entire/settings.local.json`).
- Freebuff (CLI and/or desktop) with at least one chat in the repository.

## Installation

```bash
cd agents/entire-agent-freebuff
mise run build
cp entire-agent-freebuff ~/.local/bin/
```

Or with Go directly:

```bash
cd agents/entire-agent-freebuff
go build -o entire-agent-freebuff ./cmd/entire-agent-freebuff
cp entire-agent-freebuff ~/.local/bin/
```

## Enable in a Repository

```bash
entire enable --agent freebuff --telemetry=false
entire agent list            # should show Freebuff
entire status                # Agents: ... Freebuff
```

`entire enable` installs the hook registry at `.freebuff/entire-hooks.json`
(declaring the `session-start`, `prompt-submit`, `stop`, and `session-end`
commands a Freebuff wrapper must fire) and enables checkpoints.

### Capturing a Freebuff session

Freebuff sessions live under the engine config directory:

```
~/.config/manicode/projects/<project>/chats/<session-id>/
    chat-messages.json     # original transcript format
    session.jsonl          # new JSONL event format (when produced)
    run-state.json         # maps the chat to its repository
```

The adapter resolves the repo → project → chat mapping from `run-state.json`
`projectRoot`, so a session started in this repository is found
automatically.

Two ways to turn a Freebuff session into a checkpoint:

1. **Hooks** (automatic) — when a Freebuff wrapper fires the installed hook
   commands (`entire hooks freebuff session-start` on launch, `... stop` on
   turn end), commits made during the session create checkpoints that include
   the Freebuff transcript context.
2. **Attach** (manual fallback) — attach the most recent Freebuff chat for
   the repository after a commit:

   ```bash
   entire session attach --agent freebuff <session-id>
   ```

   `<session-id>` is the chat directory name (for example
   `2026-09-06T09-24-30.025Z`). Run `entire session attach` without `--agent`
   to let Entire auto-detect the Freebuff agent from the transcript.

## Behavior

| Capability | Behavior |
|---|---|
| `hooks` | Parses Freebuff lifecycle payloads into Entire events. Unknown hook names and malformed payloads are ignored (never crash). Installs/uninstalls the repo-local hook registry under `.freebuff/`. |
| `transcript_analyzer` | Reads user prompts, file-modifying tool calls, assistant text, and checkpoint summaries from **both** transcript formats below. Truncated or partially written transcripts degrade to partial results. |

## Noon Curveball: The Agent Changed Its Format

Freebuff released a new transcript/lifecycle format (a JSONL event stream).
Existing users still produce the original format, so the adapter supports
**both** and never duplicates the pipeline:

| | Original format | New format |
|---|---|---|
| File | `chat-messages.json` | `session.jsonl` (or any `*.jsonl` in the chat dir) |
| Shape | JSON array of `{variant: user\|ai, content, blocks[]}` | JSONL: one `{event, ...}` record per line |
| Events | `text`/`tool` blocks in AI messages | `session_started`, `user_prompt`, `agent_response`, `file_changed`, `checkpoint_created`, `usage`, `session_ended`, … |

The parser in `internal/freebuff/transcript.go` normalizes both formats into
one turn model, and honors the Curveball guarantees:

- **Both formats supported** — detection by content, so the same analyzer
  entry points work unchanged.
- **Unknown events never crash** — unknown JSONL event types (and unknown
  hook names) are skipped and counted; known-but-inert events (`tool_call`,
  `tool_result`, `file_read`, `usage`) are recognized and ignored.
- **Incomplete transcripts produce partial results** — a truncated JSONL
  stream or a truncated chat-messages.json array yields the records that
  parsed cleanly (flagged partial) instead of discarding the session.

Tests in `internal/freebuff/transcript_test.go` cover all four required
cases: original format, new format, unknown events, and incomplete input.
The Curveball fixture is committed at
`testdata/track-3-agent-session.jsonl`.

## Resume

`freebuff --continue [session-id]` resumes a conversation; the adapter's
`format-resume-command` returns that command so Entire can offer a
one-command resume.

## Limitations

- Freebuff's engine does not yet expose a user-configurable hook registry.
  The adapter installs a declarative hook registry (`.freebuff/entire-hooks.json`)
  and implements `parse-hook` for every event it declares; wiring the engine
  to fire those commands is the remaining Freebuff-side step. Until then use
  `entire session attach --agent freebuff` after commits (path 2 above).
- The engine stores display timestamps in chats; session start time is taken
  from the chat directory mtime.
- Lifecycle e2e scenarios require an interactive Freebuff session
  (`FREEBDUFF_E2E=1 E2E_AGENT=freebuff`); the CLI has no headless prompt
  flag yet.
