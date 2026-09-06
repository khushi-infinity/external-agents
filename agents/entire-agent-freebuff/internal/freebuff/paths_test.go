package freebuff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/entireio/external-agents/agents/entire-agent-freebuff/internal/protocol"
)

// seedChatDir writes a realistic Freebuff chat layout under the hermetic
// config root (FREEBDUFF_CONFIG_DIR already points at the manicode config
// root, so projects live directly under it) and returns the chat dir path.
func seedChatDir(t *testing.T, configRoot, project, sessionID, repoRoot string) string {
	t.Helper()
	chatDir := filepath.Join(configRoot, "projects", project, chatsDirName, sessionID)
	if err := os.MkdirAll(chatDir, 0o700); err != nil {
		t.Fatalf("mkdir chat dir: %v", err)
	}
	runState := `{"sessionState":{"fileContext":{"projectRoot":"` + repoRoot + `","cwd":"` + repoRoot + `"}}}`
	if err := os.WriteFile(filepath.Join(chatDir, "run-state.json"), []byte(runState), 0o600); err != nil {
		t.Fatalf("write run-state: %v", err)
	}
	if err := os.WriteFile(filepath.Join(chatDir, "chat-messages.json"), []byte(originalFixture()), 0o600); err != nil {
		t.Fatalf("write chat-messages: %v", err)
	}
	return chatDir
}

func TestGetSessionDirResolvesRepoProject(t *testing.T) {
	_, configRoot := withHermeticEnv(t)
	repoRoot := t.TempDir()
	// NOTE: withHermeticEnv already set ENTIRE_REPO_ROOT to another dir; set
	// the hermetic repo here to the one the chat references.
	t.Setenv("ENTIRE_REPO_ROOT", repoRoot)

	seedChatDir(t, configRoot, "checkout", "chat-123", repoRoot)

	a := New()
	dir, err := a.GetSessionDir(repoRoot)
	if err != nil {
		t.Fatalf("get-session-dir: %v", err)
	}
	if !strings.HasSuffix(dir, filepath.Join("projects", "checkout")) {
		t.Fatalf("session dir = %q, want .../projects/checkout", dir)
	}

	ref := a.ResolveSessionFile(dir, "chat-123")
	want := filepath.Join(dir, chatsDirName, "chat-123", "chat-messages.json")
	if ref != want {
		t.Fatalf("session file = %q, want %q", ref, want)
	}
	if _, err := os.Stat(ref); err != nil {
		t.Fatalf("resolved session file missing: %v", err)
	}
}

func TestResolveSessionFilePrefersNewJSONLWhenPresent(t *testing.T) {
	_, configRoot := withHermeticEnv(t)
	repoRoot := t.TempDir()
	t.Setenv("ENTIRE_REPO_ROOT", repoRoot)
	chatDir := seedChatDir(t, configRoot, "checkout", "chat-456", repoRoot)

	// Simulate the engine switching to the new JSONL event format: the
	// original file disappears and a session.jsonl appears.
	if err := os.Remove(filepath.Join(chatDir, "chat-messages.json")); err != nil {
		t.Fatalf("remove original: %v", err)
	}
	if err := os.WriteFile(filepath.Join(chatDir, "session.jsonl"), []byte(newFormatFixture), 0o600); err != nil {
		t.Fatalf("write jsonl: %v", err)
	}

	a := New()
	dir, _ := a.GetSessionDir(repoRoot)
	ref := a.ResolveSessionFile(dir, "chat-456")
	want := filepath.Join(chatDir, "session.jsonl")
	if ref != want {
		t.Fatalf("session file = %q, want %q", ref, want)
	}

	// The analyzer reads the new format through the resolved reference.
	summary, has, err := a.ExtractSummary(ref)
	if err != nil || !has {
		t.Fatalf("summary = %q, %v, %v", summary, has, err)
	}
	if summary != "Coupon validation implemented and tested." {
		t.Fatalf("summary = %q", summary)
	}
}

func TestReadSessionReturnsChatTranscript(t *testing.T) {
	_, configRoot := withHermeticEnv(t)
	repoRoot := t.TempDir()
	t.Setenv("ENTIRE_REPO_ROOT", repoRoot)
	seedChatDir(t, configRoot, "checkout", "chat-123", repoRoot)

	a := New()
	session, err := a.ReadSession(&protocol.HookInputJSON{SessionID: "chat-123"})
	if err != nil {
		t.Fatalf("read-session: %v", err)
	}
	if session.SessionID != "chat-123" {
		t.Fatalf("session id = %q", session.SessionID)
	}
	if session.AgentName != "freebuff" {
		t.Fatalf("agent name = %q", session.AgentName)
	}
	if session.SessionRef == "" || !strings.Contains(session.SessionRef, "chat-messages.json") {
		t.Fatalf("session ref = %q", session.SessionRef)
	}
	if len(session.NativeData) == 0 {
		t.Fatal("native data should contain the chat transcript")
	}
	if !strings.Contains(string(session.NativeData), "coupon validation") {
		t.Fatal("native data does not look like a chat transcript")
	}
}

func TestDetect(t *testing.T) {
	withHermeticEnv(t)
	agent := New()

	// A repo-local .freebuff marker makes detection true regardless of
	// whether the Freebuff CLI happens to be installed on this machine.
	repoRoot := t.TempDir()
	t.Setenv("ENTIRE_REPO_ROOT", repoRoot)
	if agent.Detect().Present {
		t.Fatal("detect should be false with no marker and no engine config")
	}

	if err := os.MkdirAll(filepath.Join(repoRoot, repoMarkerDir), 0o700); err != nil {
		t.Fatalf("mkdir marker: %v", err)
	}
	if !agent.Detect().Present {
		t.Fatal("detect should be true when .freebuff exists")
	}
}
