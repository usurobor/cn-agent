package cellcheck

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// A fixture is a real repository with a real Go module at the path the recipe
// names. Nothing here stubs the toolchain: the steps this package claims to run
// are the steps these tests make it run, or the observations prove nothing.
const (
	goMod = "module cnosfixture\n\ngo 1.24\n"
	okGo  = "package ok\n\nfunc One() int { return 1 }\n"
	// Valid Go, deliberately not gofmt'd. It must COMPILE and pass vet, or the
	// format step could never be the step that fails.
	messyGo = "package ok\n\nfunc  Two()  int  {  return  2  }\n"
)

func write(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for path, body := range files {
		abs := filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

// candidate builds a repository whose BASE commit is `base` and whose
// uncommitted working tree adds `change` — which is exactly the shape a
// materialized candidate has: a worktree at the pinned base with the matter
// applied on top.
// candidate builds a repository at `base`, records the commit it pins, then
// applies `change` on top. It returns both, because the checker measures the
// change against the PINNED BASE and not against wherever the candidate left
// HEAD.
func candidate(t *testing.T, base, change map[string]string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	write(t, dir, base)
	git(t, dir, "init", "-q", "-b", "main")
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-qm", "base")
	head := strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))
	write(t, dir, change)
	return dir, head
}

func baseTree() map[string]string {
	return map[string]string{"src/go/go.mod": goMod, "src/go/ok/ok.go": okGo}
}

func names(obs Observation) []string {
	out := make([]string, 0, len(obs.Steps))
	for _, s := range obs.Steps {
		out = append(out, s.Name)
	}
	return out
}

func step(t *testing.T, obs Observation, name string) Step {
	t.Helper()
	for _, s := range obs.Steps {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("no %q step in %v", name, names(obs))
	return Step{}
}

// AC3, the passing path: all four steps run, in order, and the observation is
// `pass`.
func TestTheFourStepsRunInOrderAndPass(t *testing.T) {
	dir, base := candidate(t, baseTree(), map[string]string{
		"src/go/ok/added.go": "package ok\n\nfunc Three() int { return 3 }\n",
	})
	obs := Run(context.Background(), dir, base)
	if obs.Status != Pass {
		t.Fatalf("status = %q, want pass; steps: %+v", obs.Status, obs.Steps)
	}
	if got := strings.Join(names(obs), " "); got != "build vet test format" {
		t.Fatalf("steps ran as %q, want \"build vet test format\"", got)
	}
	for _, s := range obs.Steps {
		if s.Status != Pass || s.Exit != 0 {
			t.Fatalf("step %q = %q exit %d", s.Name, s.Status, s.Exit)
		}
	}
	if obs.Recipe != RecipeID {
		t.Fatalf("recipe = %q, want %q", obs.Recipe, RecipeID)
	}
	// Run is handed a directory, not a view, so it cannot know the candidate's
	// identity and must not invent one.
	if obs.Candidate != "" {
		t.Fatalf("Run assigned a candidate identity %q it has no way to know", obs.Candidate)
	}
}

// AC3, `fail` naming `build`: the recipe stops there, so the step list itself
// says where it stopped.
func TestACompileErrorFailsNamingBuild(t *testing.T) {
	dir, base := candidate(t, baseTree(), map[string]string{
		"src/go/ok/broken.go": "package ok\n\nfunc Four() int { return \"not an int\" }\n",
	})
	obs := Run(context.Background(), dir, base)
	if obs.Status != Fail {
		t.Fatalf("status = %q, want fail; steps: %+v", obs.Status, obs.Steps)
	}
	if got := strings.Join(names(obs), " "); got != "build" {
		t.Fatalf("steps ran as %q, want the run to stop at \"build\"", got)
	}
	b := step(t, obs, "build")
	if b.Status != Fail || b.Exit == 0 {
		t.Fatalf("build step = %q exit %d, want fail with a non-zero exit", b.Status, b.Exit)
	}
	if !strings.Contains(b.Tail, "broken.go") {
		t.Fatalf("the build tail does not name the offending file: %q", b.Tail)
	}
}

// AC3, `fail` naming `test`: build and vet pass first, so the test step is
// reached and is the one that fails.
func TestAFailingTestFailsNamingTest(t *testing.T) {
	dir, base := candidate(t, baseTree(), map[string]string{
		"src/go/ok/ok_test.go": "package ok\n\nimport \"testing\"\n\n" +
			"func TestOne(t *testing.T) {\n\tif One() == 1 {\n\t\tt.Fatal(\"deliberate failure\")\n\t}\n}\n",
	})
	obs := Run(context.Background(), dir, base)
	if obs.Status != Fail {
		t.Fatalf("status = %q, want fail; steps: %+v", obs.Status, obs.Steps)
	}
	if got := strings.Join(names(obs), " "); got != "build vet test" {
		t.Fatalf("steps ran as %q, want the run to stop at \"test\"", got)
	}
	if s := step(t, obs, "test"); s.Status != Fail || !strings.Contains(s.Tail, "deliberate failure") {
		t.Fatalf("test step = %q, tail %q", s.Status, s.Tail)
	}
}

// AC3, `fail` naming `format`: an unformatted file the candidate CHANGED.
// This is the step whose failure predicate is not the exit code — `gofmt -l`
// exits 0 and reports by listing — so a passing status here would be the
// vacuous case.
func TestAnUnformattedChangedFileFailsNamingFormat(t *testing.T) {
	dir, base := candidate(t, baseTree(), map[string]string{"src/go/ok/messy.go": messyGo})
	obs := Run(context.Background(), dir, base)
	if obs.Status != Fail {
		t.Fatalf("status = %q, want fail; steps: %+v", obs.Status, obs.Steps)
	}
	if got := strings.Join(names(obs), " "); got != "build vet test format" {
		t.Fatalf("steps ran as %q, want all four with format failing", got)
	}
	f := step(t, obs, "format")
	if f.Status != Fail {
		t.Fatalf("format step = %q, tail %q", f.Status, f.Tail)
	}
	if !strings.Contains(f.Tail, "messy.go") {
		t.Fatalf("the format tail does not name the unformatted file: %q", f.Tail)
	}
}

// AC4. `format` is scoped to the paths the candidate changed: a candidate that
// touches one well-formatted file passes even though the tree it sits in
// contains an unformatted file it never touched.
//
// This is the clause that stops the checker being red forever. The base commit
// below carries `messy.go`, and the second half of this test proves a repo-wide
// `gofmt -l src/go` in this very fixture WOULD list it — without that, the pass
// above could mean the fixture was clean all along.
func TestFormatIsScopedToTheChangedPaths(t *testing.T) {
	tree := baseTree()
	tree["src/go/ok/messy.go"] = messyGo
	dir, base := candidate(t, tree, map[string]string{
		"src/go/ok/added.go": "package ok\n\nfunc Five() int { return 5 }\n",
	})

	obs := Run(context.Background(), dir, base)
	if obs.Status != Pass {
		t.Fatalf("status = %q, want pass; steps: %+v", obs.Status, obs.Steps)
	}
	if f := step(t, obs, "format"); f.Status != Pass {
		t.Fatalf("format step = %q, tail %q — an untouched file was judged", f.Status, f.Tail)
	}

	// The untouched file really is unformatted, in this tree, right now.
	cmd := exec.Command("gofmt", "-l", "src/go")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("gofmt -l src/go: %v", err)
	}
	if !strings.Contains(string(out), "messy.go") {
		t.Fatalf("a repo-wide gofmt lists %q: the fixture proves nothing about scoping", out)
	}
}

// A candidate that changes no .go path has nothing to format. The step must
// pass rather than invoke `gofmt -l` with no operands, which reads stdin and
// would block forever.
func TestNoChangedGoPathsPassesFormat(t *testing.T) {
	dir, base := candidate(t, baseTree(), map[string]string{"NOTES.md": "a note\n"})
	obs := Run(context.Background(), dir, base)
	if obs.Status != Pass {
		t.Fatalf("status = %q, want pass; steps: %+v", obs.Status, obs.Steps)
	}
	if f := step(t, obs, "format"); f.Status != Pass || !strings.Contains(f.Tail, "no changed .go paths") {
		t.Fatalf("format step = %q, tail %q", f.Status, f.Tail)
	}
}

// AC3, `unavailable`: a step that cannot START is not a candidate's failure.
// The whole observation is unavailable and nothing downstream may read it as
// "the candidate is bad".
func TestAMissingGoIsUnavailable(t *testing.T) {
	dir, base := candidate(t, baseTree(), map[string]string{
		"src/go/ok/added.go": "package ok\n\nfunc Six() int { return 6 }\n",
	})
	t.Setenv("PATH", "")
	obs := Run(context.Background(), dir, base)
	if obs.Status != Unavailable {
		t.Fatalf("status = %q, want unavailable; steps: %+v", obs.Status, obs.Steps)
	}
	if got := strings.Join(names(obs), " "); got != "build" {
		t.Fatalf("steps ran as %q, want the run to stop at \"build\"", got)
	}
	b := step(t, obs, "build")
	if b.Status != Unavailable || b.Exit != -1 {
		t.Fatalf("build step = %q exit %d, want unavailable with exit -1", b.Status, b.Exit)
	}
	if b.Tail == "" {
		t.Fatal("an unavailable step must say why it could not start")
	}
}

// A directory that is not a repository cannot answer "what did the candidate
// change", so the format step is unavailable rather than silently repo-wide or
// silently empty. Asserted through git's absence of a repo, not through a
// missing binary, so it is a different path from the test above.
func TestANonRepositoryMakesFormatUnavailable(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, baseTree())
	obs := Run(context.Background(), dir, "0000000000000000000000000000000000000000")
	if obs.Status != Unavailable {
		t.Fatalf("status = %q, want unavailable; steps: %+v", obs.Status, obs.Steps)
	}
	if f := step(t, obs, "format"); f.Status != Unavailable || f.Exit != -1 {
		t.Fatalf("format step = %q exit %d", f.Status, f.Exit)
	}
}

// AC5, this package's half: the checker has no production caller, and this
// test fails when it gains one — at which point the package header's claim is
// what must be rewritten.
//
// The scan is proven able to find callers before it is trusted to report none:
// `cellskill.` is imported and used by production code, and a scan that could
// not see it would report every package as uncalled.
func TestCheckerHasNoProductionCaller(t *testing.T) {
	if callers := productionReferences(t, "cellcheck."); len(callers) != 0 {
		t.Fatalf("the checker now has production callers %v — it is no longer "+
			"consumer-less, and the package header says it is", callers)
	}
	if len(productionReferences(t, "cellskill.")) == 0 {
		t.Fatal("the scan found no reference to cellskill either: it is not looking at production sources")
	}
}

// productionReferences returns the non-test Go files under the module that
// mention `needle`, excluding this package's own directory.
func productionReferences(t *testing.T, needle string) []string {
	t.Helper()
	root, err := filepath.Abs("../..") // src/go, the module root
	if err != nil {
		t.Fatal(err)
	}
	self, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	var found []string
	seen := 0
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if filepath.Dir(path) == self {
			return nil
		}
		seen++
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), needle) {
			rel, _ := filepath.Rel(root, path)
			found = append(found, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	if seen < 50 {
		t.Fatalf("the walk saw only %d production files under %s", seen, root)
	}
	return found
}

// The seat has a shell and therefore git, so it may commit its own work. Once
// it does, `git status` is clean — and a HEAD-relative changed set is empty, so
// this step would report "no changed .go paths" and PASS with an unformatted
// file sitting in the change. A gate that goes green because the seat tidied up
// is worse than no gate, and this is the same defect cellwork.Diff already
// carries a witness for.
func TestACommittedCandidateIsStillMeasuredAgainstTheBase(t *testing.T) {
	dir, base := candidate(t, baseTree(), map[string]string{"src/go/ok/messy.go": messyGo})
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-qm", "the seat commits its own work")
	if head := strings.TrimSpace(git(t, dir, "rev-parse", "HEAD")); head == base {
		t.Fatal("the commit left HEAD at the base, so this proves nothing")
	}
	if status := strings.TrimSpace(git(t, dir, "status", "--porcelain")); status != "" {
		t.Fatalf("git status is not clean, so a HEAD-relative check would still see the change: %q", status)
	}

	obs := Run(context.Background(), dir, base)
	if obs.Status != Fail {
		t.Fatalf("status = %q, want fail; steps: %+v", obs.Status, obs.Steps)
	}
	if f := step(t, obs, "format"); f.Status != Fail || !strings.Contains(f.Tail, "messy.go") {
		t.Fatalf("format = %+v, want a failure naming messy.go", f)
	}
}

// git reports paths relative to the REPOSITORY ROOT, which need not be the
// directory handed to Run. Resolving that wrongly drops every path and the step
// reports "no changed .go paths" — a pass, from having looked in the wrong
// place.
func TestASubdirectoryIsRefusedRatherThanSilentlyEmpty(t *testing.T) {
	dir, base := candidate(t, baseTree(), map[string]string{"src/go/ok/messy.go": messyGo})
	sub := filepath.Join(dir, "src")
	obs := Run(context.Background(), sub, base)
	// build runs first and cannot find src/go below src, so the recipe stops
	// there — which is itself the point: the whole recipe means "this is the
	// candidate repository root". What must never happen is a green format step
	// computed from an empty changed set.
	for _, st := range obs.Steps {
		if st.Name == "format" && st.Status == Pass {
			t.Fatalf("format passed from a subdirectory: %+v", st)
		}
	}
	if obs.Status == Pass {
		t.Fatalf("a subdirectory must not pass the recipe: %+v", obs.Steps)
	}
	// And the changed-set derivation itself refuses rather than returning none.
	if _, err := changedGoFiles(context.Background(), sub, base); err == nil {
		t.Fatal("changedGoFiles must refuse a directory that is not the repository root")
	} else if !strings.Contains(err.Error(), "not the candidate repository root") {
		t.Fatalf("refused for the wrong reason: %v", err)
	}
}
