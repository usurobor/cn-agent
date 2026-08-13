package cellwork

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A symlink's content in git is its TARGET PATH. Reading through it would put
// the bytes of any host file the runtime can read into the reviewing seat's
// prompt, under a header promising the view came from (subject, matter) alone —
// so this asserts the SECRET IS ABSENT, not merely that the target is present.
func TestReconstructReportsSymlinksWithoutFollowingThem(t *testing.T) {
	repo, base := testRepo(t)
	secret := filepath.Join(t.TempDir(), "host-secret.txt")
	if err := os.WriteFile(secret, []byte("HOST-SECRET-CONTENT\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	wt, release, err := Materialize(context.Background(), repo, base)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(secret, filepath.Join(wt.Dir, "link.txt")); err != nil {
		t.Fatal(err)
	}
	matter, err := wt.Diff(context.Background())
	release()
	if err != nil {
		t.Fatal(err)
	}

	view, err := Reconstruct(context.Background(), repo, base, matter)
	if err != nil {
		t.Fatalf("a matter adding a symlink is a legal patch: %v", err)
	}
	var link *FileState
	for i := range view.Files {
		if view.Files[i].Path == "link.txt" {
			link = &view.Files[i]
		}
	}
	if link == nil {
		t.Fatalf("the view does not carry link.txt: %+v", view.Files)
	}
	if !link.Symlink {
		t.Error("a symlink must be reported as one; a reader must not read its content as file bytes")
	}
	if link.Content != secret {
		t.Errorf("content = %q, want the link target %q", link.Content, secret)
	}
	if strings.Contains(link.Content, "HOST-SECRET-CONTENT") {
		t.Fatal("the reconstruction followed the link: host bytes outside (subject, matter) reached the view")
	}
}

// A relative link whose target is not in the change is a legal patch. Following
// it ended the whole episode on a file that was never read.
func TestReconstructTakesADanglingSymlinkInStride(t *testing.T) {
	repo, base := testRepo(t)
	wt, release, err := Materialize(context.Background(), repo, base)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../nowhere/shim", filepath.Join(wt.Dir, "shim")); err != nil {
		t.Fatal(err)
	}
	matter, err := wt.Diff(context.Background())
	release()
	if err != nil {
		t.Fatal(err)
	}
	view, err := Reconstruct(context.Background(), repo, base, matter)
	if err != nil {
		t.Fatalf("a dangling link must not end the episode: %v", err)
	}
	for _, f := range view.Files {
		if f.Path == "shim" && f.Symlink && f.Content == "../nowhere/shim" {
			return
		}
	}
	t.Fatalf("the dangling link is not reported as a link: %+v", view.Files)
}

// The bound must be decided from metadata. A one-line patch to a file already
// large in the base would otherwise read the whole file before anything checked
// it — "a limit checked on a fully buffered result is not a limit".
func TestViewBoundIsDecidedBeforeReading(t *testing.T) {
	repo, _ := testRepo(t)
	// Many short lines, not one long one: the patch must stay small while the
	// FILE is large, which is exactly the case a post-read bound cannot catch.
	big := strings.Repeat("padding line to make this file large\n", (maxViewBytes+(4<<20))/37)
	if err := os.WriteFile(filepath.Join(repo, "big.txt"), []byte(big), 0o600); err != nil {
		t.Fatal(err)
	}
	gitIn(t, repo, "add", "-A")
	gitIn(t, repo, "commit", "-qm", "a file already large in the base")
	base := gitIn(t, repo, "rev-parse", "HEAD")

	wt, release, err := Materialize(context.Background(), repo, base)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wt.Dir, "big.txt"), []byte(big+"one more line\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	matter, err := wt.Diff(context.Background())
	release()
	if err != nil {
		t.Fatal(err)
	}

	var before, after runtimeMem
	readMem(&before)
	view, err := Reconstruct(context.Background(), repo, base, matter)
	readMem(&after)
	if err != nil {
		t.Fatal(err)
	}
	if !view.Truncated {
		t.Fatal("an over-bound file must mark the view truncated")
	}
	// The whole point: the oversize file was never materialized in memory.
	if grew := after.total - before.total; grew > uint64(maxViewBytes)*3 {
		t.Fatalf("reconstruction allocated %d bytes for a file it omitted; the bound was decided after reading", grew)
	}
}
