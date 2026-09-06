package freebuff

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// --- Original format fixtures (chat-messages.json array) ---

const (
	originalDivider = `{"id":"divider-1","variant":"ai","content":"","blocks":[{"type":"mode-divider","mode":"LITE"}],"timestamp":"09:00 AM"}`
	originalUser1   = `{"id":"user-1","variant":"user","content":"Add coupon validation to checkout. Add tests.","timestamp":"09:00 AM"}`
	originalAI1     = `{"id":"ai-1","variant":"ai","content":"","timestamp":"09:01 AM","blocks":[` +
		`{"type":"text","content":"I'll inspect the checkout flow.","textType":"text"},` +
		`{"type":"text","content":"(reasoning about precedence)","textType":"reasoning"},` +
		`{"type":"tool","toolCallId":"t1","toolName":"write","input":{"path":"src/checkout/apply_coupon.ts"},"output":""},` +
		`{"type":"tool","toolCallId":"t2","toolName":"write_file","input":{"path":"tests/checkout/apply_coupon.test.ts"},"output":""}]}`
	originalUser2 = `{"id":"user-2","variant":"user","content":"Fix the failing test ordering.","timestamp":"09:02 AM"}`
	originalAI2   = `{"id":"ai-2","variant":"ai","content":"","timestamp":"09:03 AM","blocks":[` +
		`{"type":"tool","toolName":"str_replace","input":{"path":"src/checkout/apply_coupon.ts"},"output":""},` +
		`{"type":"text","content":"Expiry now takes precedence.","textType":"text"}]}`
)

func originalFixture() string {
	return "[" + strings.Join([]string{originalDivider, originalUser1, originalAI1, originalUser2, originalAI2}, ",") + "]"
}

func wantOriginalFiles() []string {
	return []string{"src/checkout/apply_coupon.ts", "tests/checkout/apply_coupon.test.ts"}
}

// --- New format fixture (JSONL event stream, mirroring the Track 3 card) ---

const newFormatFixture = `{"timestamp": "2026-09-06T09:00:00.000+05:30", "event": "session_started", "session_id": "btw-track3-demo-001", "agent": {"name": "AcmeCode", "version": "1.4.2"}, "repository": "github.com/example/checkout-service", "branch": "feature/add-coupon-validation"}
{"timestamp": "2026-09-06T09:00:08.214+05:30", "event": "user_prompt", "session_id": "btw-track3-demo-001", "message_id": "msg-001", "text": "Add coupon validation to checkout. Coupons should be rejected if expired, disabled, or below the minimum cart value. Add tests."}
{"timestamp": "2026-09-06T09:00:09.048+05:30", "event": "agent_response", "session_id": "btw-track3-demo-001", "message_id": "msg-002", "parent_message_id": "msg-001", "text": "I'll inspect the checkout flow and existing promotion models, then implement validation and tests."}
{"timestamp": "2026-09-06T09:02:31.901+05:30", "event": "file_changed", "session_id": "btw-track3-demo-001", "path": "src/checkout/apply_coupon.ts", "change": "modified", "summary": "Added validation for enabled state, expiry, and minimum cart value.", "lines_added": 24, "lines_removed": 3}
{"timestamp": "2026-09-06T09:03:12.432+05:30", "event": "file_changed", "session_id": "btw-track3-demo-001", "path": "tests/checkout/apply_coupon.test.ts", "change": "modified", "summary": "Added tests for expired, disabled, low-cart-value, and valid coupons.", "lines_added": 71, "lines_removed": 0}
{"timestamp": "2026-09-06T09:04:33.778+05:30", "event": "usage", "session_id": "btw-track3-demo-001", "model": "acmecode-pro", "input_tokens": 8421, "output_tokens": 2194, "tool_calls": 3}
{"timestamp": "2026-09-06T09:04:49.600+05:30", "event": "checkpoint_created", "session_id": "btw-track3-demo-001", "checkpoint_id": "cp-001", "git_commit": "8d34f70c1e9fd62c1b5dc4fbbbf5013db2817ae1", "summary": "Coupon validation implemented and tested.", "intent": "Reject expired, disabled, or minimum-cart-value-ineligible coupons during checkout.", "open_questions": ["Should expiry or disabled state take precedence in user-facing errors?"]}
{"timestamp": "2026-09-06T09:05:02.314+05:30", "event": "session_ended", "session_id": "btw-track3-demo-001", "status": "completed", "duration_seconds": 302}
`

func writeTempTranscript(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "transcript.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func TestOriginalFormatAnalyzers(t *testing.T) {
	a := New()
	path := writeTempTranscript(t, originalFixture())

	position, err := a.GetTranscriptPosition(path)
	if err != nil || position != 2 {
		t.Fatalf("position = %d, %v; want 2", position, err)
	}

	files, pos, err := a.ExtractModifiedFiles(path, 0)
	if err != nil || pos != 2 {
		t.Fatalf("extract files = %v, %d, %v", files, pos, err)
	}
	if !reflect.DeepEqual(files, wantOriginalFiles()) {
		t.Fatalf("files = %#v, want %#v", files, wantOriginalFiles())
	}

	prompts, err := a.ExtractPrompts(path, 0)
	if err != nil || len(prompts) != 2 {
		t.Fatalf("prompts = %#v, %v; want 2 prompts", prompts, err)
	}
	if !strings.Contains(prompts[0], "coupon validation") {
		t.Fatalf("prompts[0] = %q", prompts[0])
	}

	summary, has, err := a.ExtractSummary(path)
	if err != nil || !has {
		t.Fatalf("summary = %q, %v, %v", summary, has, err)
	}
	if !strings.Contains(summary, "Expiry now takes precedence") {
		t.Fatalf("summary = %q, want the final assistant response", summary)
	}

	// Offset extraction: nothing new after the final position.
	filesAfter, posAfter, err := a.ExtractModifiedFiles(path, pos)
	if err != nil || posAfter != 2 || len(filesAfter) != 0 {
		t.Fatalf("offset extraction = %#v, %d, %v", filesAfter, posAfter, err)
	}
}

func TestNewJSONLFormatAnalyzers(t *testing.T) {
	a := New()
	path := writeTempTranscript(t, newFormatFixture)

	position, err := a.GetTranscriptPosition(path)
	if err != nil || position != 1 {
		t.Fatalf("position = %d, %v; want 1 turn", position, err)
	}

	files, _, err := a.ExtractModifiedFiles(path, 0)
	if err != nil {
		t.Fatalf("extract files: %v", err)
	}
	if !reflect.DeepEqual(files, wantOriginalFiles()) {
		t.Fatalf("files = %#v, want %#v", files, wantOriginalFiles())
	}

	prompts, err := a.ExtractPrompts(path, 0)
	if err != nil || len(prompts) != 1 {
		t.Fatalf("prompts = %#v, %v", prompts, err)
	}
	if !strings.Contains(prompts[0], "expired, disabled") {
		t.Fatalf("prompts[0] = %q", prompts[0])
	}

	// The new format carries an explicit checkpoint summary — it must win
	// over the assistant response text.
	summary, has, err := a.ExtractSummary(path)
	if err != nil || !has {
		t.Fatalf("summary = %q, %v, %v", summary, has, err)
	}
	if summary != "Coupon validation implemented and tested." {
		t.Fatalf("summary = %q, want checkpoint summary", summary)
	}
}

func TestUnknownJSONLEventsAreSkippedWithoutCrashing(t *testing.T) {
	unknownLines := []string{
		`{"timestamp": "2026-09-06T09:00:05.000+05:30", "event": "model_switched", "session_id": "btw-track3-demo-001", "model": "acmecode-pro"}`,
		`{"timestamp": "2026-09-06T09:00:06.000+05:30", "event": "agent_resumed", "session_id": "btw-track3-demo-001", "reason": "user reopened the session"}`,
		`{"timestamp": "2026-09-06T09:01:00.000+05:30", "event": "tool_policy_hint", "session_id": "btw-track3-demo-001", "hint": {"allow": ["shell"]}}`,
		`{"timestamp": "2026-09-06T09:05:10.000+05:30", "session_id": "btw-track3-demo-001", "note": "future event without an event field"}`,
	}
	content := newFormatFixture + "\n" + strings.Join(unknownLines, "\n") + "\n"

	// Analyzer surface must never crash on unknown events.
	a := New()
	path := writeTempTranscript(t, content)

	files, position, err := a.ExtractModifiedFiles(path, 0)
	if err != nil || position != 1 {
		t.Fatalf("files = %#v pos=%d err=%v; unknown events must not crash", files, position, err)
	}
	if !reflect.DeepEqual(files, wantOriginalFiles()) {
		t.Fatalf("files = %#v, want %#v (unknown events must not corrupt results)", files, wantOriginalFiles())
	}
	prompts, err := a.ExtractPrompts(path, 0)
	if err != nil || len(prompts) != 1 {
		t.Fatalf("prompts = %#v, %v", prompts, err)
	}

	// Direct parse view: the unknown records are counted, not fatal.
	parsed, err := parseTranscript([]byte(content))
	if err != nil {
		t.Fatalf("parseTranscript: %v", err)
	}
	if parsed.unknownEvents != len(unknownLines) {
		t.Fatalf("unknownEvents = %d, want %d", parsed.unknownEvents, len(unknownLines))
	}
	if !parsed.partial {
		t.Fatal("partial should be true when unknown events were skipped")
	}
}

func TestIncompleteTranscriptsDegradeToPartial(t *testing.T) {
	a := New()

	t.Run("truncated jsonl event stream", func(t *testing.T) {
		lines := strings.Split(strings.TrimSuffix(newFormatFixture, "\n"), "\n")
		last := lines[len(lines)-1]
		truncated := strings.Join(lines[:len(lines)-1], "\n") + "\n" + last[:len(last)/2]
		path := writeTempTranscript(t, truncated)

		files, position, err := a.ExtractModifiedFiles(path, 0)
		if err != nil {
			t.Fatalf("truncated JSONL must not error: %v", err)
		}
		if position != 1 {
			t.Fatalf("position = %d, want 1 (records before the cut)", position)
		}
		if !reflect.DeepEqual(files, wantOriginalFiles()) {
			t.Fatalf("files = %#v", files)
		}
		parsed, _ := parseTranscript([]byte(truncated))
		if !parsed.partial {
			t.Fatal("truncated JSONL should be flagged partial")
		}
	})

	t.Run("truncated original json array", func(t *testing.T) {
		// Cut right after ai-1: the array close and the last two messages
		// are lost. The salvage path must recover the complete prefix.
		truncated := "[" + strings.Join([]string{originalDivider, originalUser1, originalAI1}, ",")
		path := writeTempTranscript(t, truncated)

		files, position, err := a.ExtractModifiedFiles(path, 0)
		if err != nil {
			t.Fatalf("truncated chat array must not error: %v", err)
		}
		if position != 1 {
			t.Fatalf("position = %d, want 1", position)
		}
		if !reflect.DeepEqual(files, wantOriginalFiles()) {
			t.Fatalf("files = %#v, want %#v (salvaged prefix)", files, wantOriginalFiles())
		}
		prompts, _ := a.ExtractPrompts(path, 0)
		if len(prompts) != 1 || !strings.Contains(prompts[0], "coupon validation") {
			t.Fatalf("prompts = %#v, want the first prompt only", prompts)
		}
	})

	t.Run("empty transcript file", func(t *testing.T) {
		path := writeTempTranscript(t, "")
		position, err := a.GetTranscriptPosition(path)
		if err != nil || position != 0 {
			t.Fatalf("position = %d, %v; want 0", position, err)
		}
		files, _, err := a.ExtractModifiedFiles(path, 0)
		if err != nil || len(files) != 0 {
			t.Fatalf("files = %#v, %v", files, err)
		}
	})

	t.Run("session started but no activity yet", func(t *testing.T) {
		content := `{"timestamp": "2026-09-06T09:00:00.000+05:30", "event": "session_started", "session_id": "btw-track3-demo-001", "agent": {"name": "AcmeCode", "version": "1.4.2"}}` + "\n"
		path := writeTempTranscript(t, content)
		position, err := a.GetTranscriptPosition(path)
		if err != nil || position != 0 {
			t.Fatalf("position = %d, %v; want 0", position, err)
		}
		prompts, err := a.ExtractPrompts(path, 0)
		if err != nil || len(prompts) != 0 {
			t.Fatalf("prompts = %#v, %v; want empty", prompts, err)
		}
	})
}

func TestFormatDetectionAndDocumentRoundTrip(t *testing.T) {
	// Both formats must round-trip through the same analyzer entry point.
	original := parse(t, originalFixture())
	if original.sessionID != "" {
		t.Fatalf("original sessionID = %q, want empty", original.sessionID)
	}

	newFormat := parse(t, newFormatFixture)
	if newFormat.sessionID != "btw-track3-demo-001" {
		t.Fatalf("new format sessionID = %q", newFormat.sessionID)
	}

	// A malformed leading byte reports a clear error but a JSONL stream
	// without a leading brace is still attempted.
	if _, err := parseTranscript([]byte("<html>not a transcript</html>")); err == nil {
		t.Fatal("expected an error for an unsupported format")
	}
}

// TestCommittedCurveballFixture guards the exact Track 3 Curveball fixture
// (the new JSONL format attached at noon) through the full analyzer surface.
func TestCommittedCurveballFixture(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "track-3-agent-session.jsonl")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("missing committed fixture: %v", err)
	}

	a := New()
	files, position, err := a.ExtractModifiedFiles(path, 0)
	if err != nil {
		t.Fatalf("extract modified files: %v", err)
	}
	if position != 1 {
		t.Fatalf("position = %d, want 1", position)
	}
	if !reflect.DeepEqual(files, wantOriginalFiles()) {
		t.Fatalf("files = %#v, want %#v", files, wantOriginalFiles())
	}
	prompts, err := a.ExtractPrompts(path, 0)
	if err != nil || len(prompts) != 1 {
		t.Fatalf("prompts = %#v, %v", prompts, err)
	}
	summary, has, err := a.ExtractSummary(path)
	if err != nil || !has {
		t.Fatalf("summary = %q, %v, %v", summary, has, err)
	}
	if summary != "Coupon validation implemented and tested." {
		t.Fatalf("summary = %q, want checkpoint summary", summary)
	}
}

func parse(t *testing.T, content string) *parsedSession {
	t.Helper()
	session, err := parseTranscript([]byte(content))
	if err != nil {
		t.Fatalf("parseTranscript: %v", err)
	}
	return session
}
