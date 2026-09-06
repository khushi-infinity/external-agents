package freebuff

import (
	"os"
	"path/filepath"
	"testing"
)

// withHermeticEnv points ENTIRE_REPO_ROOT and FREEBDUFF_CONFIG_DIR at temp
// dirs so hook and path tests never touch the real Freebuff config.
func withHermeticEnv(t *testing.T) (repoRoot, configRoot string) {
	t.Helper()
	repoRoot = t.TempDir()
	configRoot = t.TempDir()
	t.Setenv("ENTIRE_REPO_ROOT", repoRoot)
	t.Setenv("FREEBDUFF_CONFIG_DIR", configRoot)
	return repoRoot, configRoot
}

func TestParseHookLifecycleMapping(t *testing.T) {
	repoRoot, _ := withHermeticEnv(t)
	a := New()

	// session-start
	ev, err := a.ParseHook(HookNameSessionStart, []byte(`{"session_id":"chat-123","model":"free-lite"}`))
	if err != nil {
		t.Fatalf("session-start: %v", err)
	}
	if ev == nil || ev.Type != 1 || ev.SessionID != "chat-123" {
		t.Fatalf("session-start event = %#v", ev)
	}
	if got := readRepoSessionCache(repoRoot); got != "chat-123" {
		t.Fatalf("active session cache = %q, want chat-123", got)
	}

	// prompt-submit
	ev, err = a.ParseHook(HookNamePromptSubmit, []byte(`{"session_id":"chat-123","prompt":"Fix the login bug"}`))
	if err != nil {
		t.Fatalf("prompt-submit: %v", err)
	}
	if ev == nil || ev.Type != 2 || ev.Prompt != "Fix the login bug" {
		t.Fatalf("prompt-submit event = %#v", ev)
	}

	// stop
	ev, err = a.ParseHook(HookNameStop, []byte(`{"session_id":"chat-123"}`))
	if err != nil {
		t.Fatalf("stop: %v", err)
	}
	if ev == nil || ev.Type != 3 || ev.SessionID != "chat-123" {
		t.Fatalf("stop event = %#v", ev)
	}
	if ev.SessionRef == "" {
		t.Fatal("stop event should carry a session_ref")
	}

	// session-end clears the cache
	ev, err = a.ParseHook(HookNameSessionEnd, []byte(`{"session_id":"chat-123"}`))
	if err != nil {
		t.Fatalf("session-end: %v", err)
	}
	if ev == nil || ev.Type != 5 {
		t.Fatalf("session-end event = %#v", ev)
	}
	if got := readRepoSessionCache(repoRoot); got != "" {
		t.Fatalf("cache = %q after session-end, want cleared", got)
	}
}

func TestParseHookUnknownAndMalformedInputsAreIgnored(t *testing.T) {
	_, _ = withHermeticEnv(t)
	a := New()

	// An unknown hook name must not crash the lifecycle pipeline.
	ev, err := a.ParseHook("context-compaction", []byte(`{"session_id":"chat-123"}`))
	if err != nil || ev != nil {
		t.Fatalf("unknown hook = %#v, %v; want nil, nil", ev, err)
	}

	// A malformed payload is treated as an absent event.
	ev, err = a.ParseHook(HookNamePromptSubmit, []byte(`{not json`))
	if err != nil || ev != nil {
		t.Fatalf("malformed payload = %#v, %v; want nil, nil", ev, err)
	}

	// An empty payload (hooks fired with no stdin) is also fine.
	ev, err = a.ParseHook(HookNameSessionStart, nil)
	if err != nil {
		t.Fatalf("empty payload: %v", err)
	}
	if ev == nil {
		t.Fatal("session-start with no payload should mint/attach a session")
	}
}

func TestInstallUninstallHooksRoundTrip(t *testing.T) {
	repoRoot, _ := withHermeticEnv(t)
	a := New()

	if a.AreHooksInstalled() {
		t.Fatal("hooks should not be installed in a clean repo")
	}

	count, err := a.InstallHooks(false, false)
	if err != nil {
		t.Fatalf("install hooks: %v", err)
	}
	if count != len(hookRegistryEntries()) {
		t.Fatalf("installed = %d, want %d", count, len(hookRegistryEntries()))
	}
	if !a.AreHooksInstalled() {
		t.Fatal("hooks should be installed after install")
	}

	// Idempotent.
	if again, err := a.InstallHooks(false, false); err != nil || again != 0 {
		t.Fatalf("reinstall = %d, %v; want 0, nil", again, err)
	}

	regPath := filepath.Join(repoRoot, hooksRegistryDir, hooksRegistryFile)
	if _, err := os.Stat(regPath); err != nil {
		t.Fatalf("registry not written: %v", err)
	}

	if err := a.UninstallHooks(); err != nil {
		t.Fatalf("uninstall hooks: %v", err)
	}
	if a.AreHooksInstalled() {
		t.Fatal("hooks should be gone after uninstall")
	}
	if _, err := os.Stat(regPath); !os.IsNotExist(err) {
		t.Fatalf("registry should be removed, stat err = %v", err)
	}
}
