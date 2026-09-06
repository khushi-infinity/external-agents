package freebuff

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Transcript analysis for Entire.
//
// The analyzer accepts BOTH Freebuff transcript formats (the Curveball
// requirement): the original chat-messages.json array AND the new JSONL
// event stream. Unknown JSONL events are skipped without crashing, and an
// incomplete transcript (truncated file, corrupt trailing record) yields a
// partial result instead of discarding the session.

// JSONL event names in the new Freebuff event format. Events that carry no
// analyzer content (tool_call, tool_result, file_read, usage) are still
// recognized so they are not mistaken for unknown events.
const (
	jsonlEventSessionStarted    = "session_started"
	jsonlEventUserPrompt        = "user_prompt"
	jsonlEventAgentResponse     = "agent_response"
	jsonlEventFileChanged       = "file_changed"
	jsonlEventCheckpointCreated = "checkpoint_created"
	jsonlEventSessionEnded      = "session_ended"
	jsonlEventToolCall          = "tool_call"
	jsonlEventToolResult        = "tool_result"
	jsonlEventFileRead          = "file_read"
	jsonlEventUsage             = "usage"
)

// fileModifyingTools are Freebuff/manicode tool names that change files.
// Original-format transcripts record tool blocks with these names.
var fileModifyingTools = map[string]struct{}{
	"write":       {},
	"write_file":  {},
	"edit":        {},
	"file_edit":   {},
	"str_replace": {},
	"multi_edit":  {},
	"apply_patch": {},
	"patch":       {},
	"create_file": {},
	"delete_file": {},
	"fs_write":    {},
	"fs_edit":     {},
}

func (a *Agent) ReadTranscript(sessionRef string) ([]byte, error) {
	// Accept both a file and a chat directory as the session reference.
	if info, err := os.Stat(sessionRef); err == nil && info.IsDir() {
		file := selectTranscriptFile(sessionRef)
		if file == "" {
			return nil, fmt.Errorf("no transcript file found in chat dir %s", sessionRef)
		}
		sessionRef = file
	}
	return os.ReadFile(sessionRef)
}

func (a *Agent) ChunkTranscript(content []byte, maxSize int) ([][]byte, error) {
	if maxSize <= 0 {
		return nil, errors.New("max-size must be greater than zero")
	}
	if len(content) == 0 {
		return [][]byte{{}}, nil
	}
	var chunks [][]byte
	for start := 0; start < len(content); start += maxSize {
		end := start + maxSize
		if end > len(content) {
			end = len(content)
		}
		chunk := make([]byte, end-start)
		copy(chunk, content[start:end])
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

func (a *Agent) ReassembleTranscript(chunks [][]byte) ([]byte, error) {
	return bytes.Join(chunks, nil), nil
}

// parseTranscript detects the transcript format and parses it into the
// normalized turn model. It never discards a recoverable session: corrupt
// or unknown records degrade to a partial result.
func parseTranscript(data []byte) (*parsedSession, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		// Empty transcript — a partial session with nothing recorded yet.
		return &parsedSession{}, nil
	}

	switch trimmed[0] {
	case '[':
		return parseOriginalFormat(trimmed)
	case '{':
		return parseJSONLEvents(trimmed)
	default:
		// Unknown envelope. If it is line-oriented JSON, still try the new
		// format so a stream without a leading '{' is not discarded.
		if bytes.ContainsRune(trimmed, '\n') {
			return parseJSONLEvents(trimmed)
		}
		return nil, fmt.Errorf("unsupported Freebuff transcript format (first byte %q)", trimmed[0])
	}
}

// parseOriginalFormat parses the original chat-messages.json array.
// A truncated array (partial write, interrupted session) is salvaged by
// decoding the longest valid prefix and flagged partial.
func parseOriginalFormat(data []byte) (*parsedSession, error) {
	session := &parsedSession{}
	var messages []chatMessageRaw
	if err := json.Unmarshal(data, &messages); err != nil {
		messages = salvageMessageArray(data)
		if messages == nil {
			// Nothing recoverable — keep the session as an empty partial
			// rather than discarding it or crashing the integration.
			session.partial = true
			return session, nil
		}
		session.partial = true
	}

	for _, msg := range messages {
		switch msg.Variant {
		case "user":
			session.turns = append(session.turns, turn{prompt: strings.TrimSpace(msg.Content)})
		case "ai":
			// Skip cosmetic messages (mode dividers, empty placeholders) so
			// they never create spurious turns or skew transcript position.
			if !hasMeaningfulContent(msg) {
				continue
			}
			if len(session.turns) == 0 {
				session.turns = append(session.turns, turn{})
			}
			t := session.turns[len(session.turns)-1]
			collectBlocks(&t, msg.Blocks)
			session.turns[len(session.turns)-1] = t
		default:
			// divider and unknown variants are not conversation content.
		}
	}
	return session, nil
}

// salvageMessageArray attempts to recover message objects from a truncated
// JSON array. For each trailing '}' it tries closing the array after that
// element; for a trailing ']' it tries decoding the content up to it.
// Returns nil when nothing is decodable. Transcripts are small (hundreds of
// KB), so the bounded backward scan is cheap.
func salvageMessageArray(data []byte) []chatMessageRaw {
	limitStart := 0
	if len(data) > 1<<20 {
		limitStart = len(data) - (1 << 20) // cap the backward scan at 1 MiB
	}
	for i := len(data) - 1; i >= limitStart; i-- {
		switch data[i] {
		case '}':
			// The last complete element ends here; the array close is lost.
			candidate := make([]byte, 0, len(data)-i+1)
			candidate = append(candidate, data[:i+1]...)
			candidate = append(candidate, ']')
			var messages []chatMessageRaw
			if json.Unmarshal(candidate, &messages) == nil && len(messages) > 0 {
				return messages
			}
		case ']':
			var messages []chatMessageRaw
			if json.Unmarshal(data[:i+1], &messages) == nil && len(messages) > 0 {
				return messages
			}
		}
	}
	return nil
}

// parseJSONLEvents parses the new JSONL event format. Every line is parsed
// independently: malformed or unknown lines are skipped (counted), and the
// recovered records still produce a partial session.
func parseJSONLEvents(data []byte) (*parsedSession, error) {
	session := &parsedSession{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var ev jsonlEventRaw
		if err := json.Unmarshal(line, &ev); err != nil {
			// Corrupt or truncated record — skip, keep the session partial.
			session.partial = true
			session.unknownEvents++
			continue
		}
		applyJSONLEvent(session, ev)
	}
	if err := scanner.Err(); err != nil {
		// Truncated final line; scanner returns ErrBufferLimit only for
		// oversized lines. Treat a read error as partial, not fatal.
		session.partial = true
	}
	return session, nil
}

func applyJSONLEvent(session *parsedSession, ev jsonlEventRaw) {
	switch ev.Event {
	case jsonlEventSessionStarted:
		if ev.SessionID != "" {
			session.sessionID = ev.SessionID
		}
		if session.turns == nil {
			session.turns = []turn{}
		}
	case jsonlEventUserPrompt:
		session.turns = append(session.turns, turn{
			prompt: strings.TrimSpace(ev.Text),
			model:  ev.Model,
		})
	case jsonlEventAgentResponse:
		t := lastTurnOrAppend(session)
		text := strings.TrimSpace(ev.Text)
		if text != "" {
			if t.response != "" {
				t.response += "\n" + text
			} else {
				t.response = text
			}
		}
		session.turns[len(session.turns)-1] = t
	case jsonlEventFileChanged:
		t := lastTurnOrAppend(session)
		path := strings.TrimSpace(ev.Path)
		if path == "" {
			return
		}
		switch strings.ToLower(strings.TrimSpace(ev.Change)) {
		case "created", "added", "new":
			t.newFiles = appendUnique(t.newFiles, path)
		case "deleted", "removed":
			t.deletedFiles = appendUnique(t.deletedFiles, path)
		default: // modified, edited, updated, or unspecified
			t.files = appendUnique(t.files, path)
		}
		session.turns[len(session.turns)-1] = t
	case jsonlEventCheckpointCreated:
		t := lastTurnOrAppend(session)
		if strings.TrimSpace(ev.Summary) != "" {
			t.summary = strings.TrimSpace(ev.Summary)
			t.hasSummary = true
		}
		if strings.TrimSpace(ev.Intent) != "" {
			t.intent = strings.TrimSpace(ev.Intent)
		}
		if len(ev.OpenQuestions) > 0 {
			t.openQuestions = append(t.openQuestions, ev.OpenQuestions...)
		}
		session.turns[len(session.turns)-1] = t
	case jsonlEventSessionEnded, jsonlEventToolCall, jsonlEventToolResult,
		jsonlEventFileRead, jsonlEventUsage:
		// Recognized events that carry no analyzer content.
		return
	default:
		// Unknown event type: the agent changed its format again. Never
		// crash — record and continue so older checkpoints stay compatible.
		session.partial = true
		session.unknownEvents++
	}
}

// hasMeaningfulContent reports whether an AI message carries conversation
// content worth recording (text or tool blocks).
func hasMeaningfulContent(msg chatMessageRaw) bool {
	if strings.TrimSpace(msg.Content) != "" {
		return true
	}
	for _, block := range msg.Blocks {
		if block.Type == "text" {
			return true
		}
		if block.Type == "tool" && block.ToolName != "" {
			return true
		}
	}
	return false
}

// lastTurnOrAppend returns the last turn, appending an empty one if the
// session has no turns yet (events before the first user prompt).
func lastTurnOrAppend(session *parsedSession) turn {
	if len(session.turns) == 0 {
		session.turns = append(session.turns, turn{})
	}
	return session.turns[len(session.turns)-1]
}

func appendUnique(list []string, value string) []string {
	for _, existing := range list {
		if existing == value {
			return list
		}
	}
	return append(list, value)
}

// collectBlocks folds the original-format AI message blocks into the
// current turn: text blocks become the response, tool blocks with file
// paths become modified files.
func collectBlocks(t *turn, blocks []chatBlockRaw) {
	for _, block := range blocks {
		switch block.Type {
		case "text":
			if block.TextType == "reasoning" {
				continue
			}
			text := strings.TrimSpace(block.Content)
			if text == "" {
				continue
			}
			if t.response != "" {
				t.response += "\n" + text
			} else {
				t.response = text
			}
		case "tool":
			name := strings.ToLower(strings.TrimSpace(block.ToolName))
			if _, ok := fileModifyingTools[name]; !ok {
				continue
			}
			for _, path := range pathsFromToolInput(block.Input) {
				t.files = appendUnique(t.files, path)
			}
		}
	}
}

// pathsFromToolInput extracts file paths from a tool block input object.
// Supports {"path": "..."}, {"file_path": "..."}, {"paths": [...]} and
// {"file_paths": [...]} shapes.
func pathsFromToolInput(input json.RawMessage) []string {
	if len(input) == 0 {
		return nil
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(input, &raw); err != nil {
		return nil
	}
	var paths []string
	add := func(value string) {
		value = strings.TrimSpace(value)
		if value != "" {
			paths = append(paths, value)
		}
	}
	for _, key := range []string{"path", "file_path", "filePath"} {
		if rawValue, ok := raw[key]; ok {
			var s string
			if json.Unmarshal(rawValue, &s) == nil {
				add(s)
			}
		}
	}
	for _, key := range []string{"paths", "file_paths", "files"} {
		if rawValue, ok := raw[key]; ok {
			var list []string
			if json.Unmarshal(rawValue, &list) == nil {
				for _, s := range list {
					add(s)
				}
			}
		}
	}
	return paths
}

// readAndParse opens and parses a transcript at the given path.
func readAndParse(path string) (*parsedSession, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseTranscript(data)
}

func (a *Agent) GetTranscriptPosition(path string) (int, error) {
	session, err := readAndParse(path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return len(session.turns), nil
}

func (a *Agent) ExtractModifiedFiles(path string, offset int) ([]string, int, error) {
	session, err := readAndParse(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	total := len(session.turns)
	if offset < 0 {
		offset = 0
	}
	if offset > total {
		offset = total
	}
	var files []string
	for _, t := range session.turns[offset:] {
		files = append(files, t.files...)
	}
	return dedupeAndSort(files), total, nil
}

func (a *Agent) ExtractPrompts(sessionRef string, offset int) ([]string, error) {
	session, err := readAndParse(sessionRef)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if offset < 0 {
		offset = 0
	}
	var prompts []string
	for _, t := range session.turns[offset:] {
		if t.prompt != "" {
			prompts = append(prompts, t.prompt)
		}
	}
	return prompts, nil
}

func (a *Agent) ExtractSummary(sessionRef string) (string, bool, error) {
	session, err := readAndParse(sessionRef)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	// Prefer the most recent checkpoint summary (new format), then the
	// latest assistant response text (both formats).
	for i := len(session.turns) - 1; i >= 0; i-- {
		if session.turns[i].hasSummary && session.turns[i].summary != "" {
			return session.turns[i].summary, true, nil
		}
	}
	for i := len(session.turns) - 1; i >= 0; i-- {
		if session.turns[i].response != "" {
			return session.turns[i].response, true, nil
		}
	}
	return "", false, nil
}

func dedupeAndSort(list []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, value := range list {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
