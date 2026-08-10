package cellcog

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// fakeBin writes an executable stand-in for the provider CLI so the adapter's
// own contract — stdin delivery, output capture, bounds, timeout, exit codes —
// is tested without renting cognition.
func fakeBin(t *testing.T, script string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fake-claude")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestClaudeCLIDeliversPromptAndCapturesAnswer(t *testing.T) {
	// Echo the prompt back so we can prove stdin reached the child.
	bin := fakeBin(t, `cat`)
	got, err := ClaudeCLI{Bin: bin}.Complete(context.Background(), "the prompt")
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if got != "the prompt" {
		t.Fatalf("adapter did not round-trip the prompt: %q", got)
	}
}

func TestClaudeCLIFailsOnMissingBinary(t *testing.T) {
	if _, err := (ClaudeCLI{Bin: "definitely-not-a-real-binary"}).Complete(context.Background(), "p"); err == nil {
		t.Fatal("missing binary must fail closed")
	}
}

func TestClaudeCLISurfacesExitFailure(t *testing.T) {
	bin := fakeBin(t, "cat >/dev/null; echo 'boom' >&2; exit 7")
	_, err := ClaudeCLI{Bin: bin}.Complete(context.Background(), "p")
	if err == nil {
		t.Fatal("non-zero exit must fail")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("diagnostic lost the provider's stderr: %v", err)
	}
}

// A runaway provider is bounded rather than trusted: over-long output is an
// error, not a truncated answer silently handed to the parser.
func TestClaudeCLIBoundsTheAnswer(t *testing.T) {
	bin := fakeBin(t, "cat >/dev/null; yes ABCDEFGHIJ | head -c 5000")
	_, err := ClaudeCLI{Bin: bin, MaxBytes: 1024}.Complete(context.Background(), "p")
	if err == nil || !strings.Contains(err.Error(), "more than 1024 bytes") {
		t.Fatalf("want bound error, got %v", err)
	}
}

func TestClaudeCLIHonorsTimeout(t *testing.T) {
	bin := fakeBin(t, "cat >/dev/null; sleep 30")
	start := time.Now()
	_, err := ClaudeCLI{Bin: bin, Timeout: 150 * time.Millisecond}.Complete(context.Background(), "p")
	if err == nil {
		t.Fatal("a hung provider must time out")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("timeout did not fire promptly: %s", elapsed)
	}
}

func TestClaudeCLIHonorsCallerCancellation(t *testing.T) {
	bin := fakeBin(t, "cat >/dev/null; sleep 30")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (ClaudeCLI{Bin: bin}).Complete(ctx, "p"); err == nil {
		t.Fatal("a cancelled context must abort the provider")
	}
}
