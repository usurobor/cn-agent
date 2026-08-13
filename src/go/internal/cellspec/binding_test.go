package cellspec

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// pinnedSHA is the shape a RECORDED base must have, transcribed here rather
// than reached for inside cellwork: this test is about what the loader binds,
// so it must be able to state the expectation without borrowing the code that
// produces it.
var pinnedSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)

// oneCommitRepo builds a throwaway repository and returns its path. Kept local
// and minimal: cellwork's own tests need the same thing, and a shared test
// helper package for fourteen lines would be a dependency between two test
// suites that are otherwise independent.
func oneCommitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"},
		{"add", "-A"},
		{"commit", "-qm", "base", "--allow-empty"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	return dir
}

// AC2, the half this package owns: an AUTHORED `HEAD` reaches the bound
// contract as an exact commit. Pinning itself is cellwork's (and is proven
// there against a moving HEAD); what is proven here is that the pinned value
// is what Build freezes — the two halves together are "pins once, into the
// contract".
func TestAnAuthoredHeadReachesTheBoundContractPinned(t *testing.T) {
	repo := oneCommitRepo(t)
	authored, err := json.Marshal(cellwork.Subject{
		Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: "HEAD",
	})
	if err != nil {
		t.Fatal(err)
	}
	pinned, err := cellwork.Pin(context.Background(), authored)
	if err != nil {
		t.Fatalf("pin: %v", err)
	}

	s, err := Parse([]byte(stubCell))
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	kspec, _, err := r.Build(context.Background(), cellfill.CddFills(), Binding{
		Issue:   json.RawMessage(`{"opaque":"issue"}`),
		Design:  json.RawMessage(`{"opaque":"design"}`),
		Subject: pinned,
	})
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	bound, err := cellwork.AdmitSubject(kspec.Contract.Subject)
	if err != nil {
		t.Fatalf("the bound subject must be admissible at a station: %v", err)
	}
	if !pinnedSHA.MatchString(bound.BaseSHA) {
		t.Fatalf("bound base_sha = %q, want 40 hex", bound.BaseSHA)
	}
	if !filepath.IsAbs(bound.Repo) {
		t.Fatalf("bound repo %q is not absolute", bound.Repo)
	}
	// The authored document really was the moving-name case, or the assertion
	// above would hold for an input that arrived pinned.
	if !strings.Contains(string(authored), `"base_sha":"HEAD"`) {
		t.Fatalf("the authored subject was not the moving-revision case: %s", authored)
	}
}

// The loader carries the binding VERBATIM and reads none of it. Byte equality
// is the whole assertion: a normalizing pass here would make the record a
// record of what this package understood, and would move the digest for
// documents no author changed.
func TestBindingIsCarriedVerbatimAndNotDecoded(t *testing.T) {
	bind := Binding{
		// Deliberately not canonical JSON: odd spacing, and key order that a
		// re-serialization would sort.
		Issue:   json.RawMessage(`{ "z":1,  "a":"issue" }`),
		Design:  json.RawMessage(`{ "z":2,  "a":"design" }`),
		Subject: json.RawMessage(`{ "z":3,  "a":"subject" }`),
	}
	s, err := Parse([]byte(stubCell))
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	kspec, _, err := r.Build(context.Background(), cellfill.CddFills(), bind)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	for name, pair := range map[string][2]string{
		"issue":   {string(kspec.Contract.Issue), string(bind.Issue)},
		"design":  {string(kspec.Contract.Design), string(bind.Design)},
		"subject": {string(kspec.Contract.Subject), string(bind.Subject)},
	} {
		if pair[0] != pair[1] {
			t.Errorf("bound %s is not the supplied bytes:\n got: %s\nwant: %s", name, pair[0], pair[1])
		}
	}
}

// A zero Binding leaves the contract exactly as it was before run inputs
// existed. Every cell in the corpus that carries no run input depends on this,
// and `omitempty` only keeps the slots out of the canonical bytes if they are
// nil rather than empty.
func TestAZeroBindingBindsNothing(t *testing.T) {
	kspec, _ := buildCell(t, stubCell, nil)
	if kspec.Contract.Issue != nil || kspec.Contract.Design != nil || kspec.Contract.Subject != nil {
		t.Fatalf("a zero binding put something in the contract: %+v", kspec.Contract)
	}
}

// subjectRequiringRegistry is `cdd.stub` re-registered as a fill that DECLARES
// it cannot act without a subject, over a factory that records being called.
// The counter is the assertion instrument: "refused before construction" is a
// claim about what did not run, and only a recorded call can falsify it.
func subjectRequiringRegistry(constructions *int) cellfill.Registry {
	reg := cellfill.CddFills()
	inner := reg.Alpha[cellfill.FillStubAlpha]
	reg.Alpha[cellfill.FillStubAlpha] = cellfill.AlphaFill{
		Construct: func(ctx context.Context, decl json.RawMessage) (cellfill.ConstructedAlpha, error) {
			*constructions++
			return inner.Construct(ctx, decl)
		},
		NeedsSubject: true,
	}
	return reg
}

// The declared requirement decides the run BEFORE the constructor runs. This is
// the property the `cds.patch` regression violated: with the requirement
// undeclared, a subjectless run built the provider adapter and read every skill
// body, and then reported a decisive inadmissibility as a station malfunction
// from inside the episode.
func TestADeclaredSubjectRequirementIsRefusedBeforeConstruction(t *testing.T) {
	s, err := Parse([]byte(stubCell))
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}

	var constructions int
	_, _, err = r.Build(context.Background(), subjectRequiringRegistry(&constructions), Binding{})
	if err == nil {
		t.Fatal("a fill that cannot act without a subject must not be built without one")
	}
	if constructions != 0 {
		t.Fatalf("the constructor ran %d time(s) before the refusal", constructions)
	}
	// The message names the fill and the input it is missing, so the operator
	// is told what to supply rather than that something malfunctioned.
	for _, want := range []string{cellfill.FillStubAlpha, "contract.subject", "--input"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("the refusal does not name %q: %v", want, err)
		}
	}

	// ...and with a subject the SAME registry constructs. Without this the
	// assertions above would hold for a registry that refused everything.
	if _, _, err := r.Build(context.Background(), subjectRequiringRegistry(&constructions),
		Binding{Subject: json.RawMessage(`{"opaque":"subject"}`)}); err != nil {
		t.Fatalf("a supplied subject must satisfy the requirement: %v", err)
	}
	if constructions != 1 {
		t.Fatalf("constructions = %d, want exactly the one admitted build", constructions)
	}
}

// The requirement is per-fill, not a rule the loader applies to everyone: the
// same cell over the UNMODIFIED registry, where `cdd.stub` declares nothing,
// still builds with an empty binding. That is every corpus cell carrying no run
// input, and TestAZeroBindingBindsNothing above is its standing witness.
