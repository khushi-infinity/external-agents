# Freebuff — External Agent Research

## Verdict: COMPATIBLE (transcript + session layout verified on macOS)

Freebuff (the "free AI coding agent", built on the manicode engine) stores
repo-scoped chat sessions on disk with rich, parseable transcripts. Entire
can record Freebuff sessions via the external-agent protocol.

## Static Checks

| Check | Result | Notes |
|-------|--------|-------|
| Binary present | PASS | `~/.npm-global/bin/freebuff` (v0.0.171) + engine binary `~/.config/manicode/freebuff` |
| Help available | PASS | `freebuff --help` (usage: `freebuff [options] [command]`) |
| Version info | PASS | `freebuff --version` → `0.0.171` |
| Hook keywords | FAIL | No native hook/plugin/lifecycle config surface exposed by the CLI |
| Session keywords | PASS | `--continue [conversation-id]`, `--cwd <directory>` |
| Config directory | PASS | `~/.config/manicode/projects/<project>/chats/<session-id>/` (verified) |
| Documentation | WARN | No public Freebuff extension docs; engine layout verified empirically |

## Binary

- Name: `freebuff` (npm launcher) → manicode engine binary
- Install: `npm i -g freebuff` or desktop app (mac arm64 `.dmg`)
- Flags: `--continue [id]`, `--cwd <dir>`; TUI only (no headless prompt flag)

## Hook Mechanism

- Native hook config: none exposed. Freebuff is a terminal/desktop TUI.
- Adapter contract: `entire hooks freebuff <name>` with a JSON payload on
  stdin, declared in the repo-local `.freebuff/entire-hooks.json` registry
  installed by `install-hooks`. Freebuff wrappers can fire these on:
  session start, prompt submit, stop (turn end), session end.
- Malformed or unknown payloads/hook names are ignored (never crash).

## Session Management

- Session directory: `~/.config/manicode/projects/<project>/chats/<session-id>/`
- Repo mapping: `run-state.json` → `sessionState.fileContext.projectRoot`
- Session ID: the chat directory name (a UTC timestamp, e.g.
  `2026-09-06T09-24-30.025Z`)
- Project key: engine writes chats under a project dir; the adapter matches
  by repo path, then project name, then recency.

## Transcript

- Primary file: `chat-messages.json` — JSON array of messages:
  - `variant: "user"` → `content` is the user prompt
  - `variant: "ai"` → `blocks[]` holds `text` blocks (`content`,
    `textType: text|reasoning`) and `tool` blocks (`toolName`, `input.path`)
- New format (Curveball): `session.jsonl` (or any `*.jsonl` in the chat
  dir) — one `{event, ...}` record per line. Recognized events:
  `session_started`, `user_prompt`, `agent_response`, `file_changed`,
  `checkpoint_created`, `session_ended`, plus inert `tool_call`,
  `tool_result`, `file_read`, `usage`. Unknown events are skipped.
- File-modifying tools (original format): `write`, `write_file`, `edit`,
  `file_edit`, `str_replace`, `multi_edit`, `apply_patch`, `patch`,
  `create_file`, `delete_file`, `fs_write`, `fs_edit`; paths come from
  `input.path` / `input.file_path` / `input.paths`.
- Summary source: checkpoint summary (new format) else last assistant text.

## Data Storage Verification

- Chat transcripts contain real assistant text, reasoning blocks, and tool
  blocks with paths — NOT placeholders (verified on live chats).
- No secondary storage required; the chat dir is the source of truth.
- `log.jsonl` inside a chat dir is engine process logging (pino), not a
  conversation transcript — never used as transcript input.

## Protocol Mapping

| Subcommand | Implementation |
|-----------|---------------|
| `info` / `detect` | Static metadata; repo `.freebuff` marker or matching engine chat |
| `get-session-id` | Explicit payload id → `.entire/tmp/freebuff-active-session` cache → most recent chat for repo → stub |
| `get-session-dir` | `~/.config/manicode/projects/<project>` matched by repo path |
| `resolve-session-file` | `<project>/chats/<id>/chat-messages.json` (or the `.jsonl` when only the new format exists) |
| `read-session` / `write-session` | Read/write the chat transcript as native bytes |
| `read-transcript` | Raw bytes of the resolved transcript (file or chat dir) |
| `chunk/reassemble-transcript` | Generic base64 chunking |
| `format-resume-command` | `freebuff --continue [id]` |
| `parse-hook` | session-start(1), prompt-submit(2), stop(3), session-end(5) |
| `install/uninstall/are-hooks-installed` | `.freebuff/entire-hooks.json` registry (marker-verified) |
| `get-transcript-position` / `extract-*` | Normalized turn count / prompts / files / summary over both formats |

## Selected Capabilities

| Capability | Declared | Justification |
|-----------|----------|---------------|
| hooks | true | Declarative registry + payload parser; engine wrapper wiring is the Freebuff-side step |
| transcript_analyzer | true | Both formats parsed to one turn model; partial-safe |
| transcript_preparer / token_calculator / text_generator / hook_response_writer / subagent_aware_extractor / compact_transcript | false | Not needed for the verified flows |

## Gaps & Limitations

- No native Freebuff hook surface yet — `.freebuff/entire-hooks.json` is a
  declarative contract, not an engine-loaded config.
- Session start time derives from chat-dir mtime (display timestamps lack a
  date).
- Interactive lifecycle e2e requires a logged-in TUI session.

## E2E Test Prerequisites

- Entire CLI binary: `entire` from PATH / `E2E_ENTIRE_BIN`
- Agent CLI binary: `freebuff` (npm global launcher), logged in
- Non-interactive prompt command: none yet (no headless flag)
- Interactive mode: `freebuff --cwd <dir>` under tmux; gate
  `FREEBDUFF_E2E=1 E2E_AGENT=freebuff`
