package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/activate"
)

// makeHub builds a minimal valid hub at t.TempDir() so activate.Run renders
// successfully. Tests that need the default rendering path use this; tests
// proving pre-render ordering pass a nonexistent HUB_DIR instead, so that
// observing a renderer error proves rendering happened (failure mode the
// pre-render checks must prevent).
func makeHub(t *testing.T) string {
	t.Helper()
	hub := t.TempDir()
	if err := os.MkdirAll(filepath.Join(hub, ".cn"), 0o755); err != nil {
		t.Fatalf("mkdir .cn: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(hub, ".cn", "config.json"),
		[]byte(`{"name":"test-hub","version":"1.0.0","created":"2026-01-01T00:00:00Z"}`),
		0o644,
	); err != nil {
		t.Fatalf("write config.json: %v", err)
	}
	return hub
}

// emptyPath sets $PATH to a single empty temp directory so neither claude
// nor codex resolves via exec.LookPath. Required for the missing-binary AC.
func emptyPath(t *testing.T) {
	t.Helper()
	t.Setenv("PATH", t.TempDir())
}

func runActivate(t *testing.T, args []string) (string, string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := &ActivateCmd{}
	err := cmd.Run(context.Background(), Invocation{
		Args:   args,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	return stdout.String(), stderr.String(), err
}

// AC3: default no-flag path writes prompt to stdout and exits 0.
func TestActivate_DefaultNoFlag_WritesPromptToStdout(t *testing.T) {
	hub := makeHub(t)

	stdout, stderr, err := runActivate(t, []string{hub})
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "You are activating a cnos hub.") {
		t.Errorf("stdout missing activation header:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Hub name: test-hub") {
		t.Errorf("stdout missing hub name:\n%s", stdout)
	}
	// AC3 negative — no spawn diagnostics on stdout (operator pipes / redirects
	// depend on stdout containing only the rendered prompt).
	if strings.Contains(stdout, "✗") {
		t.Errorf("stdout must not carry error/diagnostic glyphs:\n%s", stdout)
	}
}

// AC4 oracle 1+2: --claude and --codex together exit non-zero and the
// stderr error names both flags + the conflict phrase.
func TestActivate_MutualExclusion_Errors(t *testing.T) {
	hub := makeHub(t)

	stdout, stderr, err := runActivate(t, []string{"--claude", "--codex", hub})
	if err == nil {
		t.Fatal("expected non-zero exit on --claude --codex")
	}
	if stdout != "" {
		t.Errorf("stdout must stay empty when mutual-exclusion fires (got %q)", stdout)
	}
	if !strings.Contains(stderr, "--claude") || !strings.Contains(stderr, "--codex") {
		t.Errorf("stderr must name both conflicting flags, got %q", stderr)
	}
	if !strings.Contains(stderr, "mutually exclusive") {
		t.Errorf("stderr must declare the conflict, got %q", stderr)
	}
}

// AC4 oracle 3: mutual-exclusion check runs BEFORE render. Proof: pass a
// HUB_DIR that does not exist. If rendering had been attempted, stderr
// would carry "Hub path not found"; the mutually-exclusive error must
// take precedence.
func TestActivate_MutualExclusion_FiresBeforeRender(t *testing.T) {
	nonexistent := "/this/path/must/not/exist/cnos-380-test"

	_, stderr, err := runActivate(t, []string{"--claude", "--codex", nonexistent})
	if err == nil {
		t.Fatal("expected non-zero exit on --claude --codex")
	}
	if !strings.Contains(stderr, "mutually exclusive") {
		t.Errorf("expected 'mutually exclusive' diagnostic, got: %q", stderr)
	}
	if strings.Contains(stderr, "Hub path not found") {
		t.Error("rendering ran before mutual-exclusion check — pre-render ordering violated")
	}
	if strings.Contains(stderr, "Generating activation prompt") {
		t.Error("renderer 'Generating activation prompt' diagnostic appeared — pre-render ordering violated")
	}
}

// AC5 oracle 1+2+3: missing-binary error names the binary AND the flag AND
// offers actionable guidance.
func TestActivate_MissingBinary_NamesBinaryAndFlag_Claude(t *testing.T) {
	emptyPath(t)
	hub := makeHub(t)

	_, stderr, err := runActivate(t, []string{"--claude", hub})
	if err == nil {
		t.Fatal("expected non-zero exit when claude is not on PATH")
	}
	if !strings.Contains(stderr, "claude") {
		t.Errorf("stderr must name 'claude' (the missing binary), got: %q", stderr)
	}
	if !strings.Contains(stderr, "--claude") {
		t.Errorf("stderr must name '--claude' (the requesting flag), got: %q", stderr)
	}
	if !strings.Contains(stderr, "PATH") {
		t.Errorf("stderr must mention PATH (actionable), got: %q", stderr)
	}
}

// AC5 oracle 4: same shape for --codex.
func TestActivate_MissingBinary_NamesBinaryAndFlag_Codex(t *testing.T) {
	emptyPath(t)
	hub := makeHub(t)

	_, stderr, err := runActivate(t, []string{"--codex", hub})
	if err == nil {
		t.Fatal("expected non-zero exit when codex is not on PATH")
	}
	if !strings.Contains(stderr, "codex") {
		t.Errorf("stderr must name 'codex', got: %q", stderr)
	}
	if !strings.Contains(stderr, "--codex") {
		t.Errorf("stderr must name '--codex', got: %q", stderr)
	}
}

// AC5 oracle 5: missing-binary check fires BEFORE render. Same proof as
// AC4 — nonexistent HUB_DIR; the LookPath error must take precedence over
// any hub-path-not-found error.
func TestActivate_MissingBinary_FiresBeforeRender(t *testing.T) {
	emptyPath(t)
	nonexistent := "/this/path/must/not/exist/cnos-380-test"

	_, stderr, err := runActivate(t, []string{"--claude", nonexistent})
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if !strings.Contains(stderr, "claude") || !strings.Contains(stderr, "--claude") {
		t.Errorf("expected missing-binary diagnostic, got: %q", stderr)
	}
	if strings.Contains(stderr, "Hub path not found") {
		t.Error("rendering ran before LookPath check — pre-render ordering violated")
	}
	if strings.Contains(stderr, "Generating activation prompt") {
		t.Error("renderer diagnostic appeared on stderr — pre-render ordering violated")
	}
}

// AC1 / AC2 oracle: flag is parsed and surfaces in --help output naming the
// interactive-spawn behavior.
func TestActivate_HelpFlag_DocumentsClaudeAndCodex(t *testing.T) {
	stdout, _, err := runActivate(t, []string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "--claude") {
		t.Errorf("--help output missing --claude flag:\n%s", stdout)
	}
	if !strings.Contains(stdout, "--codex") {
		t.Errorf("--help output missing --codex flag:\n%s", stdout)
	}
	if !strings.Contains(stdout, "interactive") {
		t.Errorf("--help output must name the interactive-spawn behavior:\n%s", stdout)
	}
}

// AC3 (bytes-equal): the default no-flag CLI path must produce stdout
// byte-identical to a direct activate.Run invocation against the same hub.
//
// 3.78.0's cn activate HUB_DIR was a thin wrapper around activate.Run. The
// cycle/380 refactor adds --claude / --codex handling but the no-flag arm
// must keep that same property: cmd_activate.Run forwards inv.Stdout to
// activate.Run unchanged. This test fails if the refactor accidentally
// routes the default path through a buffer or adds any new bytes to stdout.
func TestActivate_DefaultNoFlag_BytesEqualToDirectRun(t *testing.T) {
	hub := makeHub(t)

	var cliStdout, cliStderr bytes.Buffer
	cmd := &ActivateCmd{}
	if err := cmd.Run(context.Background(), Invocation{
		Args:   []string{hub},
		Stdout: &cliStdout,
		Stderr: &cliStderr,
	}); err != nil {
		t.Fatalf("cli run: %v", err)
	}

	var directStdout, directStderr bytes.Buffer
	if err := activate.Run(context.Background(), activate.Options{
		HubPath: hub,
		Stdout:  &directStdout,
		Stderr:  &directStderr,
	}); err != nil {
		t.Fatalf("direct run: %v", err)
	}

	if cliStdout.String() != directStdout.String() {
		t.Errorf("default no-flag stdout != direct activate.Run stdout — pipe / redirect consumers would regress\n--- cli stdout ---\n%s\n--- direct stdout ---\n%s",
			cliStdout.String(), directStdout.String())
	}
}

func TestActivate_UnknownFlag_Errors(t *testing.T) {
	hub := makeHub(t)

	_, stderr, err := runActivate(t, []string{"--unknown-flag", hub})
	if err == nil {
		t.Fatal("expected non-zero exit on unknown flag")
	}
	if !strings.Contains(stderr, "--unknown-flag") {
		t.Errorf("stderr must name the unknown flag, got: %q", stderr)
	}
}
