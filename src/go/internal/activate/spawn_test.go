//go:build unix

package activate

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

// recordingExec is the spawn-test stand-in for syscall.Exec: it records the
// argv0 / argv / envv that the production code would have passed to
// syscall.Exec, then returns nil so the test continues.
type recordingExec struct {
	calls int
	argv0 string
	argv  []string
	envv  []string
}

func (r *recordingExec) fn(argv0 string, argv []string, envv []string) error {
	r.calls++
	r.argv0 = argv0
	r.argv = argv
	r.envv = envv
	return nil
}

// fakeLookPath returns a deterministic resolved path for binary or an error
// when the requested name is not in the map. The path it returns is not
// required to exist on disk — the test never runs the file; it only checks
// what argv would have been handed to syscall.Exec.
func fakeLookPath(installed map[string]string) lookPathFunc {
	return func(name string) (string, error) {
		if p, ok := installed[name]; ok {
			return p, nil
		}
		return "", &exec.Error{Name: name, Err: exec.ErrNotFound}
	}
}

func TestSpawnWith_ClaudeArgvShape(t *testing.T) {
	const fakePath = "/fake/bin/claude"
	const prompt = "You are activating a cnos hub.\n\nHub path: /demo\n"
	rec := &recordingExec{}
	lp := fakeLookPath(map[string]string{"claude": fakePath})

	if err := spawnWith("claude", prompt, lp, rec.fn); err != nil {
		t.Fatalf("spawnWith: unexpected error: %v", err)
	}
	if rec.calls != 1 {
		t.Fatalf("expected exactly one exec call, got %d", rec.calls)
	}
	if rec.argv0 != fakePath {
		t.Errorf("argv0 (resolved path passed to syscall.Exec) = %q, want %q", rec.argv0, fakePath)
	}
	// γ scaffold §Failure modes #2 — argv after argv[0] must be exactly
	// [prompt]: not ["-p", prompt], not ["--print", prompt]. The argv
	// element at index 0 is the binary name (Unix exec convention).
	wantArgv := []string{"claude", prompt}
	if len(rec.argv) != len(wantArgv) {
		t.Fatalf("argv length = %d, want %d (got %#v)", len(rec.argv), len(wantArgv), rec.argv)
	}
	for i, want := range wantArgv {
		if rec.argv[i] != want {
			t.Errorf("argv[%d] = %q, want %q", i, rec.argv[i], want)
		}
	}
}

func TestSpawnWith_CodexArgvShape(t *testing.T) {
	const fakePath = "/fake/bin/codex"
	const prompt = "You are activating a cnos hub.\n"
	rec := &recordingExec{}
	lp := fakeLookPath(map[string]string{"codex": fakePath})

	if err := spawnWith("codex", prompt, lp, rec.fn); err != nil {
		t.Fatalf("spawnWith: unexpected error: %v", err)
	}
	if rec.argv0 != fakePath {
		t.Errorf("argv0 = %q, want %q", rec.argv0, fakePath)
	}
	// γ scaffold §Failure modes #2 / AC2 negative — argv after argv[0]
	// must be exactly [prompt], NOT [exec, prompt]. The `codex exec` form
	// is the explicit non-interactive subcommand and would defeat the
	// cycle's purpose; bare-positional codex "$prompt" is the canonical
	// interactive launch.
	if len(rec.argv) != 2 {
		t.Fatalf("codex argv length = %d, want 2 (got %#v) — would imply [codex, exec, prompt] or similar", len(rec.argv), rec.argv)
	}
	if rec.argv[0] != "codex" {
		t.Errorf("codex argv[0] = %q, want %q", rec.argv[0], "codex")
	}
	if rec.argv[1] != prompt {
		t.Errorf("codex argv[1] = %q, want the rendered prompt verbatim", rec.argv[1])
	}
	if rec.argv[1] == "exec" {
		t.Error("codex argv[1] must NOT be 'exec' — `codex exec <prompt>` is the non-interactive subcommand")
	}
}

func TestSpawnWith_LookPathError(t *testing.T) {
	rec := &recordingExec{}
	lp := fakeLookPath(map[string]string{}) // empty PATH — nothing resolves

	err := spawnWith("claude", "irrelevant", lp, rec.fn)
	if err == nil {
		t.Fatal("expected spawnWith to fail when binary is not in PATH")
	}
	if rec.calls != 0 {
		t.Errorf("exec must not run when LookPath fails (calls = %d)", rec.calls)
	}
	if !strings.Contains(err.Error(), "claude") {
		t.Errorf("error must name the binary, got %q", err)
	}
}

func TestCheckSpawnBinary_Present(t *testing.T) {
	lp := fakeLookPath(map[string]string{"claude": "/fake/bin/claude"})
	if err := checkSpawnBinary("claude", "--claude", lp); err != nil {
		t.Errorf("expected nil for installed binary, got %v", err)
	}
}

func TestCheckSpawnBinary_MissingNamesBinaryAndFlag(t *testing.T) {
	lp := fakeLookPath(map[string]string{}) // empty PATH

	err := checkSpawnBinary("claude", "--claude", lp)
	if err == nil {
		t.Fatal("expected error when binary missing from PATH")
	}
	msg := err.Error()
	// γ scaffold §Failure modes #5 / AC5 partial-fail: error must name
	// BOTH the missing binary AND the flag that requested it, plus an
	// actionable installation / PATH-fix hint.
	if !strings.Contains(msg, "claude") {
		t.Errorf("error must name missing binary 'claude', got %q", msg)
	}
	if !strings.Contains(msg, "--claude") {
		t.Errorf("error must name requesting flag '--claude', got %q", msg)
	}
	if !strings.Contains(msg, "PATH") {
		t.Errorf("error must mention PATH (actionable guidance), got %q", msg)
	}
	if !strings.Contains(msg, "install") {
		t.Errorf("error must mention 'install' (actionable guidance), got %q", msg)
	}
}

func TestCheckSpawnBinary_CodexErrorShape(t *testing.T) {
	lp := fakeLookPath(map[string]string{})

	err := checkSpawnBinary("codex", "--codex", lp)
	if err == nil {
		t.Fatal("expected error when codex binary missing from PATH")
	}
	msg := err.Error()
	if !strings.Contains(msg, "codex") {
		t.Errorf("error must name missing binary 'codex', got %q", msg)
	}
	if !strings.Contains(msg, "--codex") {
		t.Errorf("error must name requesting flag '--codex', got %q", msg)
	}
}

// TestSpawn_DefaultsAreWired confirms the production Spawn / CheckSpawnBinary
// entrypoints route through exec.LookPath + syscall.Exec via the unexported
// hook variables. Verified by swapping defaultLookPath to a stub that
// observes the binary name and returns a known sentinel error.
func TestSpawn_DefaultsAreWired(t *testing.T) {
	saved := defaultLookPath
	t.Cleanup(func() { defaultLookPath = saved })

	sentinel := errors.New("look-path observed")
	var observed string
	defaultLookPath = func(name string) (string, error) {
		observed = name
		return "", sentinel
	}

	if err := CheckSpawnBinary("claude", "--claude"); err == nil {
		t.Fatal("expected error from sentinel lookPath")
	}
	if observed != "claude" {
		t.Errorf("default lookPath was not invoked with the binary name (got %q)", observed)
	}
}
