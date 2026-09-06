package protocol_test

import (
	"testing"

	"github.com/entireio/external-agents/agents/entire-agent-freebuff/internal/freebuff"
	"github.com/entireio/external-agents/agents/entire-agent-freebuff/internal/protocol"
)

func TestInfoResponseShape(t *testing.T) {
	info := freebuff.New().Info()

	if info.ProtocolVersion != protocol.ProtocolVersion {
		t.Fatalf("protocol_version = %d, want %d", info.ProtocolVersion, protocol.ProtocolVersion)
	}
	if info.Name != "freebuff" {
		t.Fatalf("name = %q, want %q", info.Name, "freebuff")
	}
	if info.Type != "Freebuff" {
		t.Fatalf("type = %q, want %q", info.Type, "Freebuff")
	}
	if info.Description == "" {
		t.Fatal("description should not be empty")
	}
	if !info.IsPreview {
		t.Fatal("is_preview should be true for scaffold")
	}
	if len(info.ProtectedDirs) != 1 || info.ProtectedDirs[0] != ".freebuff" {
		t.Fatalf("protected_dirs = %#v, want [.freebuff]", info.ProtectedDirs)
	}
	wantHooks := []string{
		"session-start",
		"prompt-submit",
		"stop",
		"session-end",
	}
	if len(info.HookNames) != len(wantHooks) {
		t.Fatalf("hook_names len = %d, want %d", len(info.HookNames), len(wantHooks))
	}
	for i, want := range wantHooks {
		if info.HookNames[i] != want {
			t.Fatalf("hook_names[%d] = %q, want %q", i, info.HookNames[i], want)
		}
	}
	if !info.Capabilities.Hooks {
		t.Fatal("hooks capability should be true")
	}
	if !info.Capabilities.TranscriptAnalyzer {
		t.Fatal("transcript_analyzer capability should be true")
	}
	if !info.Capabilities.UsesTerminal {
		t.Fatal("uses_terminal capability should be true")
	}
	if info.Capabilities.TranscriptPreparer || info.Capabilities.TokenCalculator ||
		info.Capabilities.TextGenerator || info.Capabilities.HookResponseWriter ||
		info.Capabilities.SubagentAwareExtractor || info.Capabilities.CompactTranscript {
		t.Fatalf("unexpected extra capabilities enabled: %#v", info.Capabilities)
	}
}
