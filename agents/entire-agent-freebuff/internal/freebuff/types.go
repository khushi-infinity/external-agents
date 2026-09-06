package freebuff

import "encoding/json"

// Native Freebuff (manicode engine) session transcript formats.
//
// Freebuff sessions are stored under the engine config directory:
//
//	<config>/manicode/projects/<project>/chats/<session-id>/
//	    chat-messages.json   original format (JSON array of chat messages)
//	    chat-meta.json       session metadata
//	    run-state.json       engine runtime state (embeds projectRoot)
//
// The engine may additionally emit a *new* JSONL event stream (one JSON
// object per line, each carrying an "event" field such as session_started,
// user_prompt, file_changed, checkpoint_created). The parser in transcript.go
// accepts both formats, skips unknown events, and degrades incomplete
// transcripts to partial results instead of discarding the session.

// turn is the normalized unit of a parsed Freebuff session: everything that
// happened while the agent worked on one user prompt (plus, for the new
// JSONL format, checkpoint/summary context).
type turn struct {
	prompt        string
	model         string
	files         []string
	newFiles      []string
	deletedFiles  []string
	response      string
	summary       string
	hasSummary    bool
	intent        string
	openQuestions []string
}

// parsedSession is the normalized result of parsing one transcript file.
type parsedSession struct {
	sessionID     string
	partial       bool // input was truncated or contained skipped records
	unknownEvents int  // records skipped because they were unknown/corrupt
	turns         []turn
}

// chatMessageRaw is one message in the original Freebuff chat-messages.json
// format. Fields are matched leniently because the engine may add fields.
type chatMessageRaw struct {
	ID        string         `json:"id"`
	Variant   string         `json:"variant"`
	Content   string         `json:"content"`
	Blocks    []chatBlockRaw `json:"blocks"`
	Timestamp string         `json:"timestamp"`
}

type chatBlockRaw struct {
	Type     string          `json:"type"` // text | tool | mode-divider
	Content  string          `json:"content"`
	TextType string          `json:"textType"`
	ToolName string          `json:"toolName"`
	Input    json.RawMessage `json:"input"`
	Output   json.RawMessage `json:"output"`
	Mode     string          `json:"mode"`
}

// jsonlEventRaw is the envelope of one record in the new JSONL event format.
// Only "event" is guaranteed; everything else depends on the event type.
type jsonlEventRaw struct {
	Event     string `json:"event"`
	SessionID string `json:"session_id"`
	Timestamp string `json:"timestamp"`
	MessageID string `json:"message_id"`
	ParentID  string `json:"parent_message_id"`
	Text      string `json:"text"`
	Path      string `json:"path"`
	Change    string `json:"change"`
	Summary   string `json:"summary"`
	Intent    string `json:"intent"`
	Tool      string `json:"tool"`
	Model     string `json:"model"`
	Agent     struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"agent"`
	OpenQuestions []string        `json:"open_questions"`
	Raw           json.RawMessage `json:"-"`
}
