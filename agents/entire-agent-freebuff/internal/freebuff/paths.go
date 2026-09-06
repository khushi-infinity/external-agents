package freebuff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/entireio/external-agents/agents/entire-agent-freebuff/internal/protocol"
)

const (
	// engineConfigDirName is the Freebuff engine (manicode) config directory
	// name. Verified on macOS: ~/.config/manicode/projects/<project>/chats/.
	engineConfigDirName = "manicode"
	chatsDirName        = "chats"
)

// configRoot returns the Freebuff engine config root. Order of preference:
//  1. FREEBDUFF_CONFIG_DIR (explicit override)
//  2. $XDG_CONFIG_HOME/manicode on Unix
//  3. ~/.config/manicode on Unix (verified layout for the Freebuff engine)
//  4. %APPDATA%\manicode on Windows
func configRoot() string {
	if dir := os.Getenv("FREEBDUFF_CONFIG_DIR"); dir != "" {
		return dir
	}
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, engineConfigDirName)
		}
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(engineConfigDirName, "projects")
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, engineConfigDirName)
}

func projectsRoot() string {
	return filepath.Join(configRoot(), "projects")
}

// runStateRaw is a lenient view of run-state.json; only projectRoot matters
// for repo matching.
type runStateRaw struct {
	SessionState struct {
		FileContext struct {
			ProjectRoot string `json:"projectRoot"`
		} `json:"fileContext"`
	} `json:"sessionState"`
}

// projectDirForRepo finds the Freebuff project directory whose chat
// run-state references the given repo path. Falls back to a project whose
// name matches the repo base name, then to any project, then to a
// deterministic <projects>/<base-name> path (even if not yet created).
func projectDirForRepo(repoPath string) string {
	if repoPath == "" {
		return filepath.Join(projectsRoot(), "default")
	}
	abs, err := filepath.Abs(repoPath)
	if err != nil {
		abs = repoPath
	}
	abs = filepath.Clean(abs)

	base := sanitizeProjectName(filepath.Base(abs))

	projects, err := os.ReadDir(projectsRoot())
	if err != nil {
		return filepath.Join(projectsRoot(), base)
	}

	type candidate struct {
		path string
	}
	var matched []candidate
	var nameMatch string
	var anyProject string

	for _, entry := range projects {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(projectsRoot(), entry.Name())
		if anyProject == "" {
			anyProject = dir
		}
		if entry.Name() == base {
			nameMatch = dir
		}
		if stateHasRepoRoot(dir, abs) {
			matched = append(matched, candidate{path: dir})
		}
	}

	if len(matched) > 0 {
		// Most-recently touched project wins when several chats reference
		// the same repo (multiple checkouts).
		sort.Slice(matched, func(i, j int) bool {
			return dirMtime(matched[i].path) > dirMtime(matched[j].path)
		})
		return matched[0].path
	}
	if nameMatch != "" {
		return nameMatch
	}
	if anyProject != "" {
		return anyProject
	}
	return filepath.Join(projectsRoot(), base)
}

// stateHasRepoRoot reports whether any chat run-state in the project dir
// references the repo root.
func stateHasRepoRoot(projectDir, repoRoot string) bool {
	chatRoot := filepath.Join(projectDir, chatsDirName)
	entries, err := os.ReadDir(chatRoot)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(chatRoot, entry.Name(), "run-state.json"))
		if err != nil {
			continue
		}
		var raw runStateRaw
		if err := json.Unmarshal(data, &raw); err != nil {
			continue
		}
		if filepath.Clean(raw.SessionState.FileContext.ProjectRoot) == repoRoot {
			return true
		}
	}
	return false
}

// latestChatDirForRepo returns the most recently modified chat session
// directory whose run-state references the given repo. Used when no session
// id is otherwise known (e.g. a hook fired without an explicit session id).
func latestChatDirForRepo(repoPath string) string {
	abs, err := filepath.Abs(repoPath)
	if err != nil {
		abs = repoPath
	}
	abs = filepath.Clean(abs)

	projects, err := os.ReadDir(projectsRoot())
	if err != nil {
		return ""
	}

	var best string
	var bestMtime int64 = -1
	for _, project := range projects {
		if !project.IsDir() {
			continue
		}
		chatRoot := filepath.Join(projectsRoot(), project.Name(), chatsDirName)
		chats, err := os.ReadDir(chatRoot)
		if err != nil {
			continue
		}
		for _, chat := range chats {
			if !chat.IsDir() {
				continue
			}
			chatDir := filepath.Join(chatRoot, chat.Name())
			// Only chats that reference this repository count; otherwise a
			// chat from an unrelated repo could be misattributed.
			statePath := filepath.Join(chatDir, "run-state.json")
			state, err := os.ReadFile(statePath)
			if err != nil {
				continue
			}
			var raw runStateRaw
			if err := json.Unmarshal(state, &raw); err != nil {
				continue
			}
			if filepath.Clean(raw.SessionState.FileContext.ProjectRoot) != abs {
				continue
			}
			info, err := os.Stat(chatDir)
			if err != nil {
				continue
			}
			if mtime := info.ModTime().UnixMilli(); mtime > bestMtime {
				best = chatDir
				bestMtime = mtime
			}
		}
	}
	return best
}

// chatSessionIDFromDir extracts the session id (chat directory name) from a
// chat directory path.
func chatSessionIDFromDir(chatDir string) string {
	return filepath.Base(chatDir)
}

// originalTranscriptFile is the primary transcript file inside a chat dir.
func originalTranscriptFile(chatDir string) string {
	return filepath.Join(chatDir, "chat-messages.json")
}

// jsonlTranscriptFile is the new-format JSONL event stream inside a chat dir,
// if the engine produced one (e.g. "session.jsonl").
func jsonlTranscriptFile(chatDir string) string {
	return filepath.Join(chatDir, "session.jsonl")
}

// chatDirFromSessionFile returns the chat dir that owns a transcript file.
func chatDirFromSessionFile(sessionFile string) string {
	return filepath.Dir(sessionFile)
}

// selectTranscriptFile returns the existing transcript inside a chat dir,
// preferring the original chat-messages.json, then the new-format JSONL.
// Returns "" when the chat dir does not exist.
func selectTranscriptFile(chatDir string) string {
	if chatDir == "" {
		return ""
	}
	if _, err := os.Stat(originalTranscriptFile(chatDir)); err == nil {
		return originalTranscriptFile(chatDir)
	}
	if _, err := os.Stat(jsonlTranscriptFile(chatDir)); err == nil {
		return jsonlTranscriptFile(chatDir)
	}
	// Any *.jsonl sidecar is a candidate for the new event format.
	entries, err := os.ReadDir(chatDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".jsonl") {
			return filepath.Join(chatDir, entry.Name())
		}
	}
	return ""
}

func sanitizeProjectName(name string) string {
	if name == "" || name == "." || name == string(filepath.Separator) {
		return "default"
	}
	name = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			return r
		default:
			return '_'
		}
	}, name)
	if name == "" {
		return "default"
	}
	return name
}

func dirMtime(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.ModTime().UnixMilli()
}

// The following two methods satisfy the protocol sessionDirResolver and
// sessionFileResolver interfaces.

func (a *Agent) GetSessionDir(repoPath string) (string, error) {
	if repoPath == "" {
		repoPath = protocol.RepoRoot()
	}
	return projectDirForRepo(repoPath), nil
}

func (a *Agent) ResolveSessionFile(sessionDir, sessionID string) string {
	if sessionID == "" {
		sessionID = stubSessionID
	}
	chatDir := filepath.Join(sessionDir, chatsDirName, sessionID)
	if existing := selectTranscriptFile(chatDir); existing != "" {
		return existing
	}
	return originalTranscriptFile(chatDir)
}
