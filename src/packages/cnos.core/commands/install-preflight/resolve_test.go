package installpreflight

import (
	"context"
	"os/exec"
	"testing"
)

func runGitPreflightTest(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func TestResolveRepoFromGitRemote_HTTPSForm(t *testing.T) {
	dir := t.TempDir()
	runGitPreflightTest(t, dir, "init", "-q")
	runGitPreflightTest(t, dir, "remote", "add", "origin", "https://github.com/acme/widgets.git")

	got, err := resolveRepoFromGitRemote(context.Background(), dir)
	if err != nil {
		t.Fatalf("resolveRepoFromGitRemote: %v", err)
	}
	if got != "acme/widgets" {
		t.Errorf("got %q, want acme/widgets", got)
	}
}

func TestResolveRepoFromGitRemote_SSHForm(t *testing.T) {
	dir := t.TempDir()
	runGitPreflightTest(t, dir, "init", "-q")
	runGitPreflightTest(t, dir, "remote", "add", "origin", "git@github.com:acme/widgets.git")

	got, err := resolveRepoFromGitRemote(context.Background(), dir)
	if err != nil {
		t.Fatalf("resolveRepoFromGitRemote: %v", err)
	}
	if got != "acme/widgets" {
		t.Errorf("got %q, want acme/widgets", got)
	}
}

func TestResolveRepoFromGitRemote_NoRemote_NamedError(t *testing.T) {
	dir := t.TempDir()
	runGitPreflightTest(t, dir, "init", "-q")

	_, err := resolveRepoFromGitRemote(context.Background(), dir)
	if err == nil {
		t.Fatal("expected an error with no origin remote")
	}
	if got := err.Error(); got == "" {
		t.Error("expected a non-empty, actionable error message")
	}
}
