package agents

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Freebuff is registered only when the caller opts in with FREEBDUFF_E2E=1
// AND E2E_AGENT=freebuff, mirroring the Qwen gate. Freebuff's CLI is a
// terminal UI (no headless `-p` prompt flag yet), so lifecycle scenarios
// must run interactively under tmux against a logged-in engine.
func init() {
	if os.Getenv("FREEBDUFF_E2E") != "1" || os.Getenv("E2E_AGENT") != "freebuff" {
		return
	}
	Register(&Freebuff{})
	RegisterGate("freebuff", 1)
}

type Freebuff struct{}

const freebuffTimeoutMultiplier = 2.0

func (f *Freebuff) Name() string               { return "freebuff" }
func (f *Freebuff) Binary() string             { return "freebuff" }
func (f *Freebuff) EntireAgent() string        { return "freebuff" }
func (f *Freebuff) PromptPattern() string      { return `(?i)freebuff` }
func (f *Freebuff) TimeoutMultiplier() float64 { return freebuffTimeoutMultiplier }
func (f *Freebuff) IsExternalAgent() bool      { return true }

func (f *Freebuff) IsTransientError(out Output, _ error) bool {
	combined := strings.ToLower(out.Stdout + out.Stderr)
	for _, pattern := range []string{"overloaded", "rate limit", "429", "503", "econnreset", "etimedout", "timeout"} {
		if strings.Contains(combined, pattern) {
			return true
		}
	}
	return false
}

func (f *Freebuff) Bootstrap() error {
	return nil
}

// RunPrompt fails fast: the Freebuff CLI exposes no non-interactive prompt
// flag. Lifecycle coverage for Freebuff therefore uses StartSession (tmux).
func (f *Freebuff) RunPrompt(_ context.Context, _ string, _ string, _ ...Option) (Output, error) {
	return Output{}, errors.New("freebuff has no headless prompt flag; use StartSession for interactive lifecycle tests")
}

func (f *Freebuff) StartSession(_ context.Context, dir string) (Session, error) {
	if _, err := exec.LookPath(f.Binary()); err != nil {
		return nil, fmt.Errorf("freebuff not in PATH: %w", err)
	}
	name := fmt.Sprintf("freebuff-test-%d", time.Now().UnixNano())
	s, err := NewTmuxSession(
		name,
		dir,
		[]string{"ENTIRE_TEST_TTY"},
		f.Binary(),
		"--cwd", dir,
	)
	if err != nil {
		return nil, err
	}
	if _, err := s.WaitFor(f.PromptPattern(), 60*time.Second); err != nil {
		_ = s.Close()
		return nil, fmt.Errorf("waiting for freebuff prompt: %w", err)
	}
	return &freebuffSession{TmuxSession: s}, nil
}

type freebuffSession struct {
	*TmuxSession
}

func (s *freebuffSession) WaitFor(pattern string, timeout time.Duration) (string, error) {
	return s.TmuxSession.WaitFor(pattern, time.Duration(float64(timeout)*freebuffTimeoutMultiplier))
}
