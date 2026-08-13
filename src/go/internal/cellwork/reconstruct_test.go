package cellwork

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// changed applies fn to a fresh worktree of repo and returns the measured
// diff — the exact bytes a producing seat's matter carries, produced the exact
// way it produces them. Reconstruction is then fed real matter rather than a
// hand-written patch that only resembles one.
func changed(t *testing.T, repo, base string, fn func(dir string)) string {
	t.Helper()
	wt, release, err := Materialize(context.Background(), repo, base)
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	defer release()
	fn(wt.Dir)
	diff, err := wt.Diff(context.Background())
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if diff == "" {
		t.Fatal("the fixture change measured as nothing, so the reconstruction below would prove nothing")
	}
	return diff
}

func write(t *testing.T, dir, path, body string) {
	t.Helper()
	full := filepath.Join(dir, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func file(t *testing.T, v View, path string) FileState {
	t.Helper()
	for _, f := range v.Files {
		if f.Path == path {
			return f
		}
	}
	t.Fatalf("the view carries no %q; it has %d files", path, len(v.Files))
	return FileState{}
}

// AC3: the view is a function of (subject, matter). The content it reports is
// what the patch produces — checked against a real repository, a real measured
// diff, and the file the same change leaves on disk.
func TestReconstructReproducesThePostApplicationState(t *testing.T) {
	repo, head := testRepo(t)
	const added = "# notes\n\nthe seat wrote this\n"
	const rewritten = "base\nand a second line\n"
	matter := changed(t, repo, head, func(dir string) {
		write(t, dir, "NOTES.md", added)
		write(t, dir, "README.md", rewritten)
	})

	view, err := Reconstruct(context.Background(), repo, head, matter)
	if err != nil {
		t.Fatalf("reconstruct: %v", err)
	}
	if len(view.Files) != 2 {
		t.Fatalf("view covers %d paths, want exactly the two the matter touches: %+v", len(view.Files), view.Files)
	}
	if view.Truncated {
		t.Fatal("a two-file view must not report truncation")
	}
	if got := file(t, view, "NOTES.md"); got.Content != added || got.Status != FileAdded {
		t.Fatalf("added file: %+v", got)
	}
	if got := file(t, view, "README.md"); got.Content != rewritten || got.Status != FileModified {
		t.Fatalf("modified file: %+v", got)
	}
	// Ordered by path, so the same inputs cannot yield two orderings.
	if view.Files[0].Path != "NOTES.md" || view.Files[1].Path != "README.md" {
		t.Fatalf("view is not ordered by path: %+v", view.Files)
	}

	same, err := Reconstruct(context.Background(), repo, head, matter)
	if err != nil {
		t.Fatal(err)
	}
	if len(same.Files) != len(view.Files) {
		t.Fatal("the same subject and matter gave two different views")
	}
	for i := range same.Files {
		if same.Files[i] != view.Files[i] {
			t.Fatalf("the same subject and matter gave two different views at %d:\n %+v\n %+v",
				i, view.Files[i], same.Files[i])
		}
	}
}

// Nothing escapes but the value: the reconstruction worktree is gone, the
// repository is untouched, and git has forgotten the checkout.
func TestReconstructLeavesNothingBehind(t *testing.T) {
	repo, head := testRepo(t)
	matter := changed(t, repo, head, func(dir string) { write(t, dir, "NOTES.md", "x\n") })

	before := gitIn(t, repo, "worktree", "list")
	if _, err := Reconstruct(context.Background(), repo, head, matter); err != nil {
		t.Fatalf("reconstruct: %v", err)
	}
	if after := gitIn(t, repo, "worktree", "list"); after != before {
		t.Fatalf("a reconstruction worktree outlived the call:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if _, err := os.Stat(filepath.Join(repo, "NOTES.md")); !os.IsNotExist(err) {
		t.Fatal("reconstruction wrote into the subject repository")
	}
	if got := gitIn(t, repo, "rev-parse", "HEAD"); got != head {
		t.Fatalf("repository HEAD moved to %q", got)
	}
}

// The defect this whole operation exists for. A patch's hunks show only
// the lines around what changed; a criterion like "the symbol is imported" is
// decided by the FILE. Here the import sits far outside every hunk, so the
// matter alone cannot answer the question and the view can.
func TestViewCarriesWhatHunksCannot(t *testing.T) {
	repo, _ := testRepo(t)
	// A file whose import block is far from the line the change touches.
	var b strings.Builder
	b.WriteString("package widget\n\nimport (\n\t\"bytes\"\n\t\"fmt\"\n)\n\n")
	for i := 0; i < 120; i++ {
		b.WriteString("// filler line so the import is nowhere near the change\n")
	}
	b.WriteString("func Render() string {\n\tvar out bytes.Buffer\n\tfmt.Fprint(&out, \"old\")\n\treturn out.String()\n}\n")
	write(t, repo, "widget.go", b.String())
	gitIn(t, repo, "add", "-A")
	gitIn(t, repo, "commit", "-qm", "widget")
	head := gitIn(t, repo, "rev-parse", "HEAD")

	matter := changed(t, repo, head, func(dir string) {
		write(t, dir, "widget.go", strings.Replace(b.String(), `"old"`, `"new"`, 1))
	})

	// The premise: the matter really does not show the import. Without this the
	// test could pass on a diff that happened to include it.
	if strings.Contains(matter, `"bytes"`) {
		t.Fatalf("the fixture diff shows the import, so it proves nothing:\n%s", matter)
	}
	view, err := Reconstruct(context.Background(), repo, head, matter)
	if err != nil {
		t.Fatalf("reconstruct: %v", err)
	}
	got := file(t, view, "widget.go")
	if !strings.Contains(got.Content, "\t\"bytes\"\n") {
		t.Fatal("the reconstructed view does not carry the import the hunks omit")
	}
	if !strings.Contains(got.Content, `"new"`) {
		t.Fatal("the reconstructed view is the base file, not the post-application state")
	}
}

// AC6: every degraded path is reported as ITS OWN fact.
func TestReconstructDegradedPathsAreDistinct(t *testing.T) {
	repo, head := testRepo(t)
	matter := changed(t, repo, head, func(dir string) { write(t, dir, "NOTES.md", "x\n") })
	// A patch whose context does not exist in the base. It has to touch an
	// EXISTING file to have context at all: a patch that only creates a file
	// has none, so any content applies and nothing is proven.
	edit := changed(t, repo, head, func(dir string) { write(t, dir, "README.md", "rewritten\n") })
	stale := strings.Replace(edit, "-base", "-a line the base never had", 1)
	if stale == edit {
		t.Fatal("the fixture patch has no context line to spoil")
	}

	cases := map[string]struct {
		repo, base, matter string
		want               string
	}{
		// Both of these assert cellwork's OWN sentence, not git's. git prints
		// "error: README.md: patch does not apply" itself, so a bare "does not
		// apply" passed even with the wrapper collapsed to a generic message —
		// the wrapper was untested by the test written to test it.
		"matter does not apply": {repo, head, stale, "the matter does not apply to " + head},
		"unresolvable base":     {repo, "no-such-rev", matter, "does not resolve"},
		"not a repository":      {t.TempDir(), head, matter, "is not a git repository"},
		"matter is not a patch": {repo, head,
			"I refactored the parser and everything passes now.\n", "the matter does not apply to " + head},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			v, err := Reconstruct(context.Background(), tc.repo, tc.base, tc.matter)
			if err == nil {
				t.Fatalf("must fail, got a view of %d files", len(v.Files))
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("failed for the wrong reason: got %v, want mention of %q", err, tc.want)
			}
		})
	}
}

// AC6's fourth degraded path: the substrate cannot give the reconstruction a
// place to work. Nothing is wrong with the inputs here — the repository
// resolves, the base resolves and the matter applies — so a runtime that
// folded this into the other three would tell an operator to fix a patch that
// is fine.
//
// Stated exactly, because the earlier name overstated it: this exercises the
// TEMP-DIRECTORY allocation (`create worktree dir`), not `git worktree add`.
// The latter's own failure branch has no witness here; reaching it needs a
// repository that resolves a base and then refuses a worktree, which no
// fixture in this package produces.
func TestReconstructReportsATempDirItCannotAllocate(t *testing.T) {
	repo, head := testRepo(t)
	matter := changed(t, repo, head, func(dir string) { write(t, dir, "NOTES.md", "x\n") })

	// Materialize allocates the worktree's parent with os.MkdirTemp, which
	// reads TMPDIR at each call. Pointing it at a path that does not exist is
	// the one way to fail that step that does not depend on the process's
	// privileges: these tests run as root in CI, where a directory with the
	// write bit cleared is still writable. Set after the fixture above, which
	// needs a working temp dir of its own.
	t.Setenv("TMPDIR", filepath.Join(t.TempDir(), "no-such-directory"))

	v, err := Reconstruct(context.Background(), repo, head, matter)
	if err == nil {
		t.Fatalf("must fail, got a view of %d files", len(v.Files))
	}
	if !strings.Contains(err.Error(), "create worktree dir") {
		t.Fatalf("failed for the wrong reason: got %v, want mention of %q", err, "create worktree dir")
	}
}

// A deletion, a rename and a binary file each report distinctly, and none of
// them is reported as an empty text file.
func TestReconstructReportsDeletionRenameAndBinaryDistinctly(t *testing.T) {
	repo, _ := testRepo(t)
	write(t, repo, "DOOMED.md", "delete me\n")
	write(t, repo, "OLD-NAME.md", strings.Repeat("stable content\n", 20))
	write(t, repo, "EMPTY.md", "")
	gitIn(t, repo, "add", "-A")
	gitIn(t, repo, "commit", "-qm", "fixtures")
	head := gitIn(t, repo, "rev-parse", "HEAD")

	matter := changed(t, repo, head, func(dir string) {
		if err := os.Remove(filepath.Join(dir, "DOOMED.md")); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(filepath.Join(dir, "OLD-NAME.md"), filepath.Join(dir, "NEW-NAME.md")); err != nil {
			t.Fatal(err)
		}
		write(t, dir, "logo.bin", "\x00\x01\x02\xff\xfe")
	})

	view, err := Reconstruct(context.Background(), repo, head, matter)
	if err != nil {
		t.Fatalf("reconstruct: %v", err)
	}
	deleted := file(t, view, "DOOMED.md")
	if deleted.Status != FileDeleted || deleted.Omitted || deleted.Binary {
		t.Fatalf("a deletion must read as deleted and nothing else: %+v", deleted)
	}
	binary := file(t, view, "logo.bin")
	if !binary.Binary || !binary.Omitted || binary.Content != "" {
		t.Fatalf("a binary path must be named with its content omitted: %+v", binary)
	}
	// The empty file is the discriminator: it is the ONLY case where empty
	// content means an empty file, so if it were flagged the three would be
	// indistinguishable again.
	renamed := file(t, view, "NEW-NAME.md")
	if renamed.Status != FileRenamed || renamed.From != "OLD-NAME.md" {
		t.Fatalf("a rename must name where it came from: %+v", renamed)
	}
	if !strings.Contains(renamed.Content, "stable content") {
		t.Fatalf("a renamed file must still carry its content: %+v", renamed)
	}
	for _, f := range view.Files {
		if f.Path == "EMPTY.md" {
			t.Fatal("EMPTY.md was not touched by the matter and must not be in the view")
		}
	}
}

// An empty file that IS touched reads as empty content with no flags — the
// case that makes `Omitted` load-bearing rather than decorative.
func TestReconstructDistinguishesAnEmptyFileFromAnOmittedOne(t *testing.T) {
	repo, head := testRepo(t)
	matter := changed(t, repo, head, func(dir string) { write(t, dir, "PLACEHOLDER.md", "") })

	view, err := Reconstruct(context.Background(), repo, head, matter)
	if err != nil {
		t.Fatalf("reconstruct: %v", err)
	}
	got := file(t, view, "PLACEHOLDER.md")
	if got.Content != "" || got.Omitted || got.Binary || got.Status != FileAdded {
		t.Fatalf("an added empty file must read as added with empty content: %+v", got)
	}
}

// The bound is REPORTED, not silently dropped: content past it is omitted, the
// path is still named, and the view says it is incomplete.
// The bound is on the VIEW, not on the matter, so the fixture is two large
// files ALREADY in the base and a one-line change to each: a small diff whose
// post-application state is larger than a view may carry. That asymmetry is
// the whole reason the view needs its own bound.
func TestReconstructReportsTruncation(t *testing.T) {
	repo, _ := testRepo(t)
	big := strings.Repeat("a line that is not evidence anybody asked for\n", 12000) // ~540 KiB
	write(t, repo, "a-big.txt", big)
	write(t, repo, "b-big.txt", big)
	gitIn(t, repo, "add", "-A")
	gitIn(t, repo, "commit", "-qm", "two large files")
	head := gitIn(t, repo, "rev-parse", "HEAD")

	matter := changed(t, repo, head, func(dir string) {
		write(t, dir, "a-big.txt", big+"one more line\n")
		write(t, dir, "b-big.txt", big+"one more line\n")
		write(t, dir, "c-small.txt", "small\n")
	})
	if len(matter) > maxViewBytes {
		t.Fatalf("the fixture matter is %d bytes; the bound under test is the view's", len(matter))
	}

	view, err := Reconstruct(context.Background(), repo, head, matter)
	if err != nil {
		t.Fatalf("reconstruct: %v", err)
	}
	if !view.Truncated {
		t.Fatal("a view past its byte bound must report truncation")
	}
	if len(view.Files) != 3 {
		t.Fatalf("truncation dropped paths instead of content: %+v", view.Files)
	}
	if first := file(t, view, "a-big.txt"); first.Content == "" {
		t.Fatal("truncation started before the bound was reached")
	}
	if second := file(t, view, "b-big.txt"); !second.Omitted || second.Content != "" || second.Binary {
		t.Fatalf("the path past the bound must be named with its content omitted: %+v", second)
	}
}

// D4. Reconstruct has no production caller, and the package doc says so in the
// present tense. A sentence in a doc comment decays the moment someone wires
// the operation up, so the claim is made a fact about the code: no non-test
// source outside this package names it.
//
// Not a purity rule — it is a scope rule with an expiry. The reviewing seat
// that consumes the view is a later increment, and when it lands this test is
// what must be deleted, deliberately, alongside the doc sentence it guards.
func TestReconstructHasNoProductionCaller(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	// thisFile: <root>/src/go/internal/cellwork/reconstruct_test.go
	goRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	thisPkg := filepath.Dir(thisFile)

	var callers []string
	scanned := 0
	err := filepath.WalkDir(goRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if filepath.Dir(path) == thisPkg {
			return nil
		}
		scanned++
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), "Reconstruct(") {
			callers = append(callers, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk the Go tree: %v", err)
	}
	// A walk that reached nothing proves nothing. Without this the test passes
	// whenever goRoot resolves somewhere empty — the same vacuity the sibling
	// dispatchAllowList guard exists to close, held to the same standard.
	if scanned == 0 {
		t.Fatalf("scanned no non-test Go files under %s: this guard proved nothing", goRoot)
	}
	if len(callers) > 0 {
		t.Fatalf("Reconstruct has acquired production callers, so the package doc is now false: %v", callers)
	}
}

// parseNameStatus is the pure half, so its record shapes are checked without a
// repository: a rename carries two paths and everything else carries one, and
// mistaking which is how a view would name the wrong file.
func TestParseNameStatus(t *testing.T) {
	got, err := parseNameStatus("A\x00added.md\x00R096\x00old.md\x00new.md\x00D\x00gone.md\x00M\x00kept.md\x00")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := []change{
		{status: FileAdded, path: "added.md"},
		{status: FileRenamed, from: "old.md", path: "new.md"},
		{status: FileDeleted, path: "gone.md"},
		{status: FileModified, path: "kept.md"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d records, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("record %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
	if empty, err := parseNameStatus(""); err != nil || len(empty) != 0 {
		t.Fatalf("empty output must decode to no records, got %+v (%v)", empty, err)
	}
	for _, bad := range []string{"A\x00", "R100\x00only-one-path.md\x00", "X\x00mystery.md\x00"} {
		if _, err := parseNameStatus(bad); err == nil {
			t.Errorf("%q must not decode into a record", bad)
		}
	}
}
