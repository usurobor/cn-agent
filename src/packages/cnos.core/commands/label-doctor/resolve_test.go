package labeldoctor

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func TestResolveRepoFromGitRemote_HTTPSWithDotGit(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "remote", "add", "origin", "https://github.com/usurobor/cnos.git")

	got, err := resolveRepoFromGitRemote(context.Background(), dir)
	if err != nil {
		t.Fatalf("resolveRepoFromGitRemote: %v", err)
	}
	if got != "usurobor/cnos" {
		t.Errorf("got %q, want usurobor/cnos", got)
	}
}

func TestResolveRepoFromGitRemote_HTTPSNoDotGit(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "remote", "add", "origin", "https://github.com/usurobor/cnos")

	got, err := resolveRepoFromGitRemote(context.Background(), dir)
	if err != nil {
		t.Fatalf("resolveRepoFromGitRemote: %v", err)
	}
	if got != "usurobor/cnos" {
		t.Errorf("got %q, want usurobor/cnos", got)
	}
}

func TestResolveRepoFromGitRemote_SSHForm(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "remote", "add", "origin", "git@github.com:usurobor/cnos.git")

	got, err := resolveRepoFromGitRemote(context.Background(), dir)
	if err != nil {
		t.Fatalf("resolveRepoFromGitRemote: %v", err)
	}
	if got != "usurobor/cnos" {
		t.Errorf("got %q, want usurobor/cnos", got)
	}
}

// TestResolveRepoFromGitRemote_NotAGitRepo mirrors the real failure mode
// repoinstall.go's ensureCanonicalDispatchLabels hits in
// repoinstall_test.go's dispatch-cds tests, whose RepoRoot is a plain
// t.TempDir() with no git init: the error must be a named, actionable
// one, not a silent zero value.
func TestResolveRepoFromGitRemote_NotAGitRepo(t *testing.T) {
	dir := t.TempDir()
	if _, err := resolveRepoFromGitRemote(context.Background(), dir); err == nil {
		t.Fatal("expected an error when dir is not a git repository")
	} else if !strings.Contains(err.Error(), "could not resolve target repo") {
		t.Errorf("error should name the resolution failure, got: %v", err)
	}
}

// TestResolveRepoFromGitRemote_NoOriginRemote mirrors
// cmd_repo_install_test.go's initGitRepo helper, which git-inits a repo
// but never configures an "origin" remote.
func TestResolveRepoFromGitRemote_NoOriginRemote(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init", "-q")
	if _, err := resolveRepoFromGitRemote(context.Background(), dir); err == nil {
		t.Fatal("expected an error when no \"origin\" remote is configured")
	} else if !strings.Contains(err.Error(), "could not resolve target repo") {
		t.Errorf("error should name the resolution failure, got: %v", err)
	}
}
