package cellwork

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"unicode/utf8"
)

// applyIndependently cuts a second worktree of repo at base and applies matter
// to it WITHOUT going through Reconstruct — a plain `git apply` of the patch
// written to a file. It is the oracle: if the view and this tree disagree, the
// view is describing something the patch does not produce.
//
// Deliberately not sharing gitInput with the code under test. A parity check
// whose two sides run the same helper only proves the helper is consistent
// with itself.
func applyIndependently(t *testing.T, repo, base, matter string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "independent")
	gitIn(t, repo, "worktree", "add", "-q", "--detach", dir, base)
	patch := filepath.Join(t.TempDir(), "matter.patch")
	if err := os.WriteFile(patch, []byte(matter), 0o600); err != nil {
		t.Fatal(err)
	}
	gitIn(t, dir, "apply", "--whitespace=nowarn", patch)
	return dir
}

// changedPaths returns the paths the independent application actually changed,
// read out of git NUL-separated so a path containing a newline is one record
// rather than two. Split here rather than by parseNameStatus, which is the
// code this test exists to check.
func changedPaths(t *testing.T, dir, base string) []string {
	t.Helper()
	gitIn(t, dir, "add", "-A")
	out := gitIn(t, dir, "diff", "--cached", "--no-color", "-z", "--name-only", base)
	var paths []string
	for _, p := range strings.Split(out, "\x00") {
		if p != "" {
			paths = append(paths, p)
		}
	}
	sort.Strings(paths)
	return paths
}

// AC3. The view must equal an independently applied patch, compared file by
// file, and every class of change git can express must survive the round trip:
// added, modified, deleted, renamed, mode change, CRLF, a unicode path, a path
// with an embedded newline, binary, symlink, and an empty file.
//
// The classes are not decoration. `-z` on the path listing is what keeps the
// unicode and embedded-newline paths from arriving quoted (so the reported path
// would differ from the path on disk); `--binary` on the measurement is what
// makes the binary file reproducible at all; Lstat is what keeps the symlink
// from being followed; and CRLF is where a line-oriented reconstruction would
// silently normalize bytes the patch preserved.
func TestReconstructEqualsAnIndependentlyAppliedPatch(t *testing.T) {
	const (
		unicodePath = "документ ▲.md"
		newlinePath = "two\nlines.txt"
		crlfAfter   = "new\r\nlines\r\nand more\r\n"
	)
	repo, _ := testRepo(t)
	write(t, repo, "tool.sh", "#!/bin/sh\necho hi\n")
	write(t, repo, "crlf.txt", "old\r\nlines\r\n")
	write(t, repo, "DOOMED.md", "delete me\n")
	write(t, repo, "OLD-NAME.md", strings.Repeat("stable content\n", 20))
	gitIn(t, repo, "add", "-A")
	gitIn(t, repo, "commit", "-qm", "every class of thing a change can be")
	base := gitIn(t, repo, "rev-parse", "HEAD")

	matter := changed(t, repo, base, func(dir string) {
		write(t, dir, "NOTES.md", "added\n")                               // added
		write(t, dir, "README.md", "base\nand a second line\n")            // modified
		write(t, dir, "crlf.txt", crlfAfter)                               // CRLF preserved
		write(t, dir, unicodePath, "unicode path\n")                       // unicode path
		write(t, dir, newlinePath, "embedded newline path\n")              // path with \n
		write(t, dir, "logo.bin", "\x00\x01\x02\xff\xfe")                  // binary
		write(t, dir, "PLACEHOLDER.md", "")                                // empty file
		if err := os.Remove(filepath.Join(dir, "DOOMED.md")); err != nil { // deleted
			t.Fatal(err)
		}
		if err := os.Rename(filepath.Join(dir, "OLD-NAME.md"), // renamed
			filepath.Join(dir, "NEW-NAME.md")); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(filepath.Join(dir, "tool.sh"), 0o755); err != nil { // mode change
			t.Fatal(err)
		}
		if err := os.Symlink("../nowhere/shim", filepath.Join(dir, "shim")); err != nil { // symlink
			t.Fatal(err)
		}
	})

	view, err := Reconstruct(context.Background(), repo, base, matter)
	if err != nil {
		t.Fatalf("reconstruct: %v", err)
	}
	oracle := applyIndependently(t, repo, base, matter)

	// (1) The same paths, and only those. A view that carried more would be
	// reporting state the patch does not produce; one that carried fewer would
	// be hiding part of the candidate.
	var seen []string
	for _, f := range view.Files {
		seen = append(seen, f.Path)
	}
	sort.Strings(seen)
	want := changedPaths(t, oracle, base)
	if strings.Join(seen, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("view paths differ from the independently applied patch:\n view %q\noracle %q", seen, want)
	}

	// (2) File by file, against the bytes the oracle tree actually holds.
	for _, f := range view.Files {
		abs := filepath.Join(oracle, filepath.FromSlash(f.Path))
		switch {
		case f.Status == FileDeleted:
			if _, err := os.Lstat(abs); !os.IsNotExist(err) {
				t.Errorf("%q is reported deleted but the applied patch left it in place (%v)", f.Path, err)
			}
		case f.Symlink:
			target, err := os.Readlink(abs)
			if err != nil {
				t.Errorf("%q is reported as a link but the applied patch left no link: %v", f.Path, err)
				continue
			}
			if f.Content != target {
				t.Errorf("%q: view carries target %q, the applied patch has %q", f.Path, f.Content, target)
			}
		case f.Omitted:
			data, err := os.ReadFile(abs)
			if err != nil {
				t.Errorf("%q: %v", f.Path, err)
				continue
			}
			if f.Content != "" {
				t.Errorf("%q is omitted but still carries %d bytes of content", f.Path, len(f.Content))
			}
			// The omission must be justified by the file, not asserted: this
			// view is small enough that binary is the only reason available.
			if !f.Binary || utf8.Valid(data) {
				t.Errorf("%q is omitted as binary but the applied patch left valid UTF-8", f.Path)
			}
		default:
			data, err := os.ReadFile(abs)
			if err != nil {
				t.Errorf("%q: %v", f.Path, err)
				continue
			}
			if f.Content != string(data) {
				t.Errorf("%q: view content differs from the applied patch:\n view %q\noracle %q",
					f.Path, f.Content, string(data))
			}
		}
	}
	if view.Truncated {
		t.Error("a view of small files must not report truncation, or the comparison above is weaker than it reads")
	}

	// (3) Each class is NAMED, not merely equal. Equality alone would pass if
	// the view reported everything as `modified` with the right bytes.
	for path, want := range map[string]FileStatus{
		"NOTES.md":       FileAdded,
		"README.md":      FileModified,
		"DOOMED.md":      FileDeleted,
		"NEW-NAME.md":    FileRenamed,
		"tool.sh":        FileModified, // a mode change with no content change
		"crlf.txt":       FileModified,
		unicodePath:      FileAdded,
		newlinePath:      FileAdded,
		"logo.bin":       FileAdded,
		"shim":           FileAdded,
		"PLACEHOLDER.md": FileAdded,
	} {
		if got := file(t, view, path).Status; got != want {
			t.Errorf("%q is reported %q, want %q", path, got, want)
		}
	}
	if from := file(t, view, "NEW-NAME.md").From; from != "OLD-NAME.md" {
		t.Errorf("the rename does not name where it came from: %q", from)
	}
	// The mode change is the reason tool.sh is in the view at all: its content
	// never changed, so a reconstruction that ignored modes would omit it.
	if got := file(t, view, "tool.sh").Content; got != "#!/bin/sh\necho hi\n" {
		t.Errorf("tool.sh content changed under a mode-only change: %q", got)
	}
	if info, err := os.Lstat(filepath.Join(oracle, "tool.sh")); err != nil {
		t.Error(err)
	} else if info.Mode().Perm()&0o111 == 0 {
		t.Error("the fixture is not a mode change: the applied patch left tool.sh non-executable")
	}
	// CRLF is carried, not normalized. Asserted on the view directly because a
	// reconstruction that stripped \r would still equal an oracle read the same
	// wrong way — which it is not here, but the claim is worth pinning alone.
	if got := file(t, view, "crlf.txt").Content; got != crlfAfter {
		t.Errorf("CRLF content was not preserved: %q, want %q", got, crlfAfter)
	}
	// An empty file is empty content with no flags — the only case where empty
	// Content means an empty file.
	if got := file(t, view, "PLACEHOLDER.md"); got.Content != "" || got.Omitted || got.Binary {
		t.Errorf("an added empty file must read as empty with no flags: %+v", got)
	}
}
