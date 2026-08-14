package cellwork

import (
	"context"
	"github.com/usurobor/cnos/src/go/internal/celltest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Materialize cuts a worktree at the resolved base SHA, not at a moving name.
func TestMaterializeResolvesBase(t *testing.T) {
	repo, head := celltest.Repo(t)
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
	repo, _ := celltest.Repo(t)
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

// The measurement must be taken against the PINNED BASE, not against wherever
// the seat left HEAD. A seat has a shell and therefore git; if it commits its
// work, the index equals the worktree's HEAD and a HEAD-relative diff is
// empty — the runtime would then record "no change was made" on real work.
// Both directions are proven here: a committed change still measures, and an
// untouched worktree still measures as nothing.
func TestDiffMeasuresAgainstPinnedBaseAfterSeatCommits(t *testing.T) {
	repo, _ := celltest.Repo(t)
	wt, release, err := Materialize(context.Background(), repo, "HEAD")
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	defer release()

	// Direction (b) first: measuring against the base must not manufacture a
	// diff out of the base commit itself.
	empty, err := wt.Diff(context.Background())
	if err != nil {
		t.Fatalf("diff of untouched worktree: %v", err)
	}
	if empty != "" {
		t.Fatalf("diff of untouched worktree = %q, want empty", empty)
	}

	// Direction (a): the seat does what a seat with git can do.
	if err := os.WriteFile(filepath.Join(wt.Dir, "NOTES.md"), []byte("hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	celltest.Git(t, wt.Dir, "add", "-A")
	celltest.Git(t, wt.Dir, "commit", "-qm", "seat commits its own work")
	if head := celltest.Git(t, wt.Dir, "rev-parse", "HEAD"); head == wt.BaseSHA {
		t.Fatal("the seat's commit left HEAD at the base, so this proves nothing")
	}

	diff, err := wt.Diff(context.Background())
	if err != nil {
		t.Fatalf("diff after the seat committed: %v", err)
	}
	if !strings.Contains(diff, "NOTES.md") || !strings.Contains(diff, "+hello") {
		t.Fatalf("a committed change was measured as %q, want a diff naming NOTES.md", diff)
	}
}

// The degraded path: a worktree with no pinned base has nothing to measure
// against, and Diff must refuse rather than measure something else. Proven by
// a SIDE EFFECT rather than by the error alone, because git would also fail on
// an empty revision argument and that failure would prove nothing about the
// guard: with the guard the index is never touched, so the seat's file is
// still untracked afterwards; without it, `git add -A` has already staged it.
func TestDiffWithoutPinnedBaseRefusesToMeasure(t *testing.T) {
	repo, _ := celltest.Repo(t)
	wt, release, err := Materialize(context.Background(), repo, "HEAD")
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	defer release()

	if err := os.WriteFile(filepath.Join(wt.Dir, "NOTES.md"), []byte("hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	wt.BaseSHA = ""

	diff, err := wt.Diff(context.Background())
	if err == nil {
		t.Fatalf("an unpinned worktree measured %q, want a refusal", diff)
	}
	if diff != "" {
		t.Fatalf("a refused measurement still returned %q", diff)
	}
	if status := celltest.Git(t, wt.Dir, "status", "--porcelain"); !strings.HasPrefix(status, "??") {
		t.Fatalf("the refusal touched the index: status = %q, want the file still untracked", status)
	}
}

// An unresolvable base fails before any worktree is cut.
func TestMaterializeUnresolvableBaseFails(t *testing.T) {
	repo, _ := celltest.Repo(t)
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
