package freebuff

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/entireio/external-agents/agents/entire-agent-freebuff/internal/protocol"
)

const (
	// stubSessionID is used when a session id cannot be discovered (e.g. a
	// bare protocol round-trip in the compliance suite).
	stubSessionID = "freebuff-session-000"

	repoMarkerDir = ".freebuff" // repo-local marker written by Freebuff
)

// Agent implements the Entire external-agent protocol for the Freebuff
// coding agent (the manicode-engine CLI and desktop app).
type Agent struct{}

func New() *Agent {
	return &Agent{}
}

func (a *Agent) Info() protocol.InfoResponse {
	return protocol.InfoResponse{
		ProtocolVersion: protocol.ProtocolVersion,
		Name:            "freebuff",
		Type:            "Freebuff",
		Description:     "Freebuff - external agent plugin for Entire CLI (dual-format transcripts)",
		IsPreview:       true,
		ProtectedDirs:   []string{repoMarkerDir},
		HookNames: []string{
			HookNameSessionStart,
			HookNamePromptSubmit,
			HookNameStop,
			HookNameSessionEnd,
		},
		Capabilities: protocol.DeclaredCapabilities{
			Hooks:              true,
			TranscriptAnalyzer: true,
			UsesTerminal:       true,
		},
	}
}

func (a *Agent) Detect() protocol.DetectResponse {
	repoRoot := protocol.RepoRoot()
	// A .freebuff marker in the repo (written when Freebuff opens the repo)
	// is the primary repo-local signal.
	if _, err := os.Stat(filepath.Join(repoRoot, repoMarkerDir)); err == nil {
		return protocol.DetectResponse{Present: true}
	}
	// A Freebuff engine chat that references this repo is also a signal.
	if chatDir := latestChatDirForRepo(repoRoot); chatDir != "" {
		return protocol.DetectResponse{Present: true}
	}
	return protocol.DetectResponse{Present: false}
}

// resolveSessionID discovers the Freebuff chat session id for the repo:
//  1. explicit session id from the hook payload
//  2. repo-local active-session cache (written when hooks fire)
//  3. the most recently modified chat dir for the repo
//  4. stub
func (a *Agent) resolveSessionID(repoRoot, explicit string) string {
	if explicit != "" {
		return explicit
	}
	if cached := readRepoSessionCache(repoRoot); cached != "" {
		return cached
	}
	if chatDir := latestChatDirForRepo(repoRoot); chatDir != "" {
		return chatSessionIDFromDir(chatDir)
	}
	return stubSessionID
}

// discoverChatDir locates the chat directory for a session id under a repo's
// project, searching all project dirs.
func discoverChatDir(repoRoot, sessionID string) string {
	if sessionID == "" || sessionID == stubSessionID {
		return ""
	}
	projects, err := os.ReadDir(projectsRoot())
	if err != nil {
		return ""
	}
	for _, project := range projects {
		if !project.IsDir() {
			continue
		}
		dir := filepath.Join(projectsRoot(), project.Name(), chatsDirName, sessionID)
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		return dir
	}
	return ""
}

func (a *Agent) GetSessionID(input *protocol.HookInputJSON) string {
	if input != nil && input.SessionID != "" {
		return input.SessionID
	}
	repoRoot := protocol.RepoRoot()
	if cached := readRepoSessionCache(repoRoot); cached != "" {
		return cached
	}
	if chatDir := latestChatDirForRepo(repoRoot); chatDir != "" {
		return chatSessionIDFromDir(chatDir)
	}
	return stubSessionID
}

func (a *Agent) ReadSession(input *protocol.HookInputJSON) (protocol.AgentSessionJSON, error) {
	repoRoot := protocol.RepoRoot()
	sessionID := a.GetSessionID(input)

	sessionDir, err := a.GetSessionDir(repoRoot)
	if err != nil {
		return protocol.AgentSessionJSON{}, err
	}

	var sessionRef string
	if input != nil && input.SessionRef != "" {
		sessionRef = input.SessionRef
	} else {
		// Prefer an existing transcript over the default path so attached
		// and captured sessions read real Freebuff content.
		sessionRef = a.ResolveSessionFile(sessionDir, sessionID)
	}

	var nativeData []byte
	if sessionRef != "" {
		if _, statErr := os.Stat(sessionRef); statErr == nil {
			data, readErr := os.ReadFile(sessionRef)
			if readErr != nil {
				return protocol.AgentSessionJSON{}, fmt.Errorf("failed to read transcript: %w", readErr)
			}
			nativeData = data
		}
	}

	startTime := time.Now().UTC().Format(time.RFC3339)
	if chatDir := discoverChatDir(repoRoot, sessionID); chatDir != "" {
		if info, err := os.Stat(chatDir); err == nil {
			startTime = info.ModTime().UTC().Format(time.RFC3339)
		}
	}

	return protocol.AgentSessionJSON{
		SessionID:     sessionID,
		AgentName:     "freebuff",
		RepoPath:      repoRoot,
		SessionRef:    sessionRef,
		StartTime:     startTime,
		NativeData:    nativeData,
		ModifiedFiles: []string{},
		NewFiles:      []string{},
		DeletedFiles:  []string{},
	}, nil
}

func (a *Agent) WriteSession(session protocol.AgentSessionJSON) error {
	ref := session.SessionRef
	if ref == "" {
		sessionDir, err := a.GetSessionDir(session.RepoPath)
		if err != nil {
			return err
		}
		ref = a.ResolveSessionFile(sessionDir, session.SessionID)
	}
	if ref == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(ref), 0o700); err != nil {
		return err
	}
	return os.WriteFile(ref, session.NativeData, 0o600)
}

func (a *Agent) FormatResumeCommand(sessionID string) string {
	if sessionID == "" || sessionID == stubSessionID {
		return "freebuff --continue"
	}
	return "freebuff --continue " + sessionID
}
