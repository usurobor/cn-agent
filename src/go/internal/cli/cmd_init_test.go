package cli

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

// cnos#612 AC3: init must refuse an unrecognized "--flag" instead of
// treating it as the positional hub name. Before this fix, validHubName
// allowed '-', so "cn init --help" scaffolded a stray "cn---help/"
// directory (cnos#606 C3).
func TestInit_RejectsUnrecognizedFlag(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	cmd := &InitCmd{}
	var stdout, stderr bytes.Buffer
	err := cmd.Run(context.Background(), Invocation{
		Args:   []string{"--bogus"},
		Stdout: &stdout,
		Stderr: &stderr,
	})

	if err == nil {
		t.Fatal("expected error for unrecognized flag")
	}
	if !strings.Contains(stderr.String(), "--bogus") {
		t.Errorf("stderr should name the unrecognized flag, got: %q", stderr.String())
	}

	entries, rerr := os.ReadDir(dir)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(entries) != 0 {
		t.Errorf("must not scaffold anything for an unrecognized flag; found: %v", entries)
	}
}

// The specific footgun from cnos#606 C3: "--help" is a flag, not a hub
// name, and must not produce a "cn---help/" directory.
func TestInit_HelpFlagDoesNotScaffold(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	cmd := &InitCmd{}
	var stdout, stderr bytes.Buffer
	_ = cmd.Run(context.Background(), Invocation{
		Args:   []string{"--help"},
		Stdout: &stdout,
		Stderr: &stderr,
	})

	if _, err := os.Stat("cn---help"); err == nil {
		t.Fatal("must not scaffold cn---help/ for the --help flag")
	}
	entries, rerr := os.ReadDir(dir)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(entries) != 0 {
		t.Errorf("must not scaffold anything for --help; found: %v", entries)
	}
}

// Regression: a plain positional hub name must still work.
func TestInit_AcceptsPositionalHubName(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	cmd := &InitCmd{}
	var stdout, stderr bytes.Buffer
	err := cmd.Run(context.Background(), Invocation{
		Args:   []string{"myhub"},
		Stdout: &stdout,
		Stderr: &stderr,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v (stderr: %s)", err, stderr.String())
	}
	if _, serr := os.Stat("cn-myhub"); serr != nil {
		t.Errorf("expected cn-myhub/ to be created: %v", serr)
	}
}
