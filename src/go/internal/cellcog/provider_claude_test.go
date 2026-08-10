package cellcog

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// fakeBin writes an executable stand-in for a provider CLI so the adapter's
// own contract — stdin delivery, working dir, bounds, timeout, exit codes —
// is tested without renting cognition.
func fakeBin(t *testing.T, script string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fake-cli")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestNewClosedProviderSet(t *testing.T) {
	if _, mode, err := New(Config{Provider: "fake"}); err != nil || mode != ModeMechanical {
		t.Fatalf("fake: mode=%q err=%v", mode, err)
	}
	if _, mode, err := New(Config{Provider: "claude-cli", Model: "m"}); err != nil || mode != ModeCognitive {
		t.Fatalf("claude-cli: mode=%q err=%v", mode, err)
	}
	if _, mode, err := New(Config{Provider: "codex-cli", Model: "m"}); err != nil || mode != ModeCognitive {
		t.Fatalf("codex-cli: mode=%q err=%v", mode, err)
	}
	if _, _, err := New(Config{Provider: "clyde", Model: "m"}); err == nil {
		t.Fatal("unknown provider must fail construction")
	}
	if _, _, err := New(Config{Provider: "claude-cli"}); err == nil {
		t.Fatal("a real provider without an exact model must fail construction")
	}
}

func TestWorkDeliversPromptInDir(t *testing.T) {
	dir := t.TempDir()
	// The stand-in proves stdin content and working directory both arrive.
	bin := fakeBin(t, `cat > seen.txt; pwd >> seen.txt`)
	if err := (ClaudeCLI{Model: "m", Bin: bin}).Work(context.Background(), dir, "the prompt"); err != nil {
		t.Fatalf("work: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "seen.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "the prompt") || !strings.Contains(string(got), dir) {
		t.Fatalf("prompt/dir did not reach the child: %q", got)
	}
}

func TestWorkFailsClosed(t *testing.T) {
	if err := (ClaudeCLI{Model: "m", Bin: "no-such-binary"}).Work(context.Background(), t.TempDir(), "p"); err == nil {
		t.Fatal("missing binary must fail closed")
	}
	if err := (ClaudeCLI{Model: "m"}).Work(context.Background(), "", "p"); err == nil {
		t.Fatal("empty dir must fail closed")
	}
	bin := fakeBin(t, "cat >/dev/null; echo boom >&2; exit 7")
	err := (ClaudeCLI{Model: "m", Bin: bin}).Work(context.Background(), t.TempDir(), "p")
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("want exit failure with stderr, got %v", err)
	}
}

// A hung provider is bounded: WaitDelay closes the inherited output pipe so a
// grandchild cannot pin Wait past the timeout.
func TestWorkHonorsTimeout(t *testing.T) {
	bin := fakeBin(t, "cat >/dev/null; sleep 30")
	start := time.Now()
	err := (ClaudeCLI{Model: "m", Bin: bin, Timeout: 150 * time.Millisecond}).Work(context.Background(), t.TempDir(), "p")
	if err == nil {
		t.Fatal("a hung provider must time out")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("timeout did not fire promptly: %s", elapsed)
	}
}

func TestWorkHonorsCallerCancellation(t *testing.T) {
	bin := fakeBin(t, "cat >/dev/null; sleep 30")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := (ClaudeCLI{Model: "m", Bin: bin}).Work(ctx, t.TempDir(), "p"); err == nil {
		t.Fatal("a cancelled context must abort the provider")
	}
}

func TestCodexArgsAreTypedNotSmuggled(t *testing.T) {
	dir := t.TempDir()
	bin := fakeBin(t, `printf '%s\n' "$@" > "`+dir+`/args.txt"; cat >/dev/null`)
	if err := (CodexCLI{Model: "m1", Bin: bin}).Work(context.Background(), dir, "p"); err != nil {
		t.Fatalf("work: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "args.txt"))
	for _, want := range []string{"exec", "--model", "m1", "--sandbox", "workspace-write", "--cd"} {
		if !strings.Contains(string(got), want) {
			t.Fatalf("codex argv missing %q: %q", want, got)
		}
	}
}
