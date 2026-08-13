package cellwork

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// gitIn runs one git command in dir under a fixed identity, so a test can move
// a repository the way a caller would — including committing after a subject
// has already been pinned.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

// testRepo builds a one-commit git repository and returns its path and HEAD.
func testRepo(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) string {
		t.Helper()
		return gitIn(t, dir, args...)
	}
	run("init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("base\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("add", "-A")
	run("commit", "-qm", "base")
	return dir, run("rev-parse", "HEAD")
}

// Materialize cuts a worktree at the resolved base SHA, not at a moving name.
func TestMaterializeResolvesBase(t *testing.T) {
	repo, head := testRepo(t)
	wt, release, err := Materialize(context.Background(), repo, "HEAD")
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	defer release()
	if wt.BaseSHA != head {
		t.Fatalf("base sha = %q, want resolved HEAD %q", wt.BaseSHA, head)
	}
	if _, err := os.Stat(filepath.Join(wt.Dir, "README.md")); err != nil {
		t.Fatalf("worktree missing base content: %v", err)
	}
}

// Diff reports a file the seat created, and nothing when the seat changed
// nothing — the caller must not manufacture evidence from an empty diff.
func TestDiffReportsCreatedFileAndIsEmptyWhenUnchanged(t *testing.T) {
	repo, _ := testRepo(t)
	wt, release, err := Materialize(context.Background(), repo, "HEAD")
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	defer release()

	empty, err := wt.Diff(context.Background())
	if err != nil {
		t.Fatalf("diff of untouched worktree: %v", err)
	}
	if empty != "" {
		t.Fatalf("diff of untouched worktree = %q, want empty", empty)
	}

	if err := os.WriteFile(filepath.Join(wt.Dir, "NOTES.md"), []byte("hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	diff, err := wt.Diff(context.Background())
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if !strings.Contains(diff, "NOTES.md") || !strings.Contains(diff, "+hello") {
		t.Fatalf("diff does not report the created file: %q", diff)
	}
}

// An unresolvable base fails before any worktree is cut.
func TestMaterializeUnresolvableBaseFails(t *testing.T) {
	repo, _ := testRepo(t)
	if _, _, err := Materialize(context.Background(), repo, "no-such-rev"); err == nil {
		t.Fatal("an unresolvable base must fail")
	}
}

// A non-repository path fails closed rather than silently doing nothing.
func TestMaterializeNonRepositoryFails(t *testing.T) {
	if _, _, err := Materialize(context.Background(), t.TempDir(), "HEAD"); err == nil {
		t.Fatal("a non-repository path must fail")
	}
}
