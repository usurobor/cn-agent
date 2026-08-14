package cdsassess

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellcheck"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
	"github.com/usurobor/cnos/src/go/internal/celltest"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// adversarial is the projection this seat is constructed under, built through
// cellmethod rather than hand-written so the test cannot hold a view the
// runtime would never produce.
func adversarial(t *testing.T) cellmethod.View {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "cnos.eng", "skills", "eng", "go")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# body of eng/go\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	decl, err := json.Marshal(map[string]any{"kind": cellmethod.Kind, "skills": []string{"cnos.eng:eng/go"}})
	if err != nil {
		t.Fatal(err)
	}
	b, bodies, err := cellmethod.Load(cellskill.Tree{Root: root}, decl)
	if err != nil {
		t.Fatal(err)
	}
	return cellmethod.Adversarial(b, bodies)
}

func construct(t *testing.T, provider, model string) cellfill.ConstructedBeta {
	t.Helper()
	decl := json.RawMessage(fmt.Sprintf(`{"fill":"cds.assess","cognition":{"provider":%q,"model":%q}}`, provider, model))
	c, err := Factory()(context.Background(), decl, adversarial(t))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	return c
}

// AC4, this package's half: the reviewing seat receives no alpha state, and the
// assertion is on the CONSTRUCTOR'S INPUTS and on the seat's own fields rather
// than on a promise in a comment.
//
// The field list is read reflectively because that is the whole claim: a seat
// that later gained a worktree handle, a repository path, a session, or alpha's
// result would still satisfy any behavioural test written today, and would fail
// this one on the day the field appeared.
func TestTheSeatHoldsNoProducerState(t *testing.T) {
	seat, ok := construct(t, "fake", "").Seat.(AssessBeta)
	if !ok {
		t.Fatal("cds.assess did not construct an AssessBeta")
	}
	rt := reflect.TypeOf(seat)
	want := map[string]string{
		"judge":  "cdsassess.judge",
		"method": "cellmethod.View",
	}
	if rt.NumField() != len(want) {
		t.Fatalf("the reviewing seat has %d fields, want exactly %v", rt.NumField(), want)
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if got := f.Type.String(); want[f.Name] != got {
			t.Errorf("seat field %q is %q; the reviewing seat may hold only %v", f.Name, got, want)
		}
	}

	// The constructor's inputs are (declaration, projection). The declaration's
	// key language is closed to exactly the two keys, so a workspace, a
	// repository, a base, a transcript or a skills list cannot be declared into
	// this seat either.
	for _, bad := range []string{
		`{"fill":"cds.assess","cognition":{"provider":"fake","model":""},"workspace":{"kind":"git-worktree"}}`,
		`{"fill":"cds.assess","cognition":{"provider":"fake","model":""},"skills":["cnos.eng:eng/go"]}`,
		`{"fill":"cds.assess","cognition":{"provider":"fake","model":"","argv":["--x"]}}`,
	} {
		if _, err := Factory()(context.Background(), json.RawMessage(bad), adversarial(t)); err == nil {
			t.Errorf("the closed declaration admitted %s", bad)
		}
	}
}

// The seat is held to the cell's methodology and cannot choose it: a cell that
// declares none, or a wiring that handed this seat the producing projection,
// fails construction rather than reviewing against nothing.
func TestTheSeatRequiresTheAdversarialProjection(t *testing.T) {
	decl := json.RawMessage(`{"fill":"cds.assess","cognition":{"provider":"fake","model":""}}`)
	if _, err := Factory()(context.Background(), decl, cellmethod.View{}); err == nil ||
		!strings.Contains(err.Error(), "declares none") {
		t.Fatalf("a cell with no methodology must be refused, got %v", err)
	}
	wrong := adversarial(t)
	wrong.Role = cellmethod.RoleConstructive
	if _, err := Factory()(context.Background(), decl, wrong); err == nil ||
		!strings.Contains(err.Error(), "adversarial projection") {
		t.Fatalf("the constructive projection must be refused, got %v", err)
	}
}

// The mode a cds.assess seat reports is the truth about what held it: nothing
// rented is mechanical, a real provider is cognitive. A run behind the
// deterministic judge must never be receipted as though cognition happened.
func TestConstructedModeIsTruthful(t *testing.T) {
	if got := construct(t, "fake", "").Mode; got != cellkernel.ModeMechanical {
		t.Fatalf("the deterministic judge reported mode %q", got)
	}
	if got := construct(t, "claude-cli", "a-model-selector").Mode; got != cellkernel.ModeCognitive {
		t.Fatalf("a rented judge reported mode %q", got)
	}
	// The resolved declaration records the projection that held the seat, so a
	// reader can ask which methodology it was, and no skills of its own.
	var rd ResolvedDecl
	if err := json.Unmarshal(construct(t, "fake", "").Decl, &rd); err != nil {
		t.Fatal(err)
	}
	if rd.Methodology.Role != string(cellmethod.RoleAdversarial) || len(rd.Methodology.SHA256) != 64 {
		t.Fatalf("the resolved declaration records no adversarial projection: %+v", rd)
	}
}

// recordingJudge captures the prompt it was given and answers with whatever the
// test wants, so the seat's whole realization — admit, reconstruct, check,
// decide, reconcile — is exercised against a real repository without renting
// anything.
type recordingJudge struct {
	prompt string
	answer func(Catalogue) Assessment
	err    error
}

func (j *recordingJudge) Name() string { return "recording" }

func (j *recordingJudge) Judge(_ context.Context, prompt string, c Catalogue) (Assessment, error) {
	j.prompt = prompt
	if j.err != nil {
		return Assessment{}, j.err
	}
	return j.answer(c), nil
}

// matterFor produces the matter a cds.patch alpha would have measured, by
// making a real change in a real worktree of the repository. Hand-written
// patches only resemble matter; this IS matter.
func matterFor(t *testing.T, repo, base string, write map[string]string) cellkernel.Matter {
	t.Helper()
	wt, release, err := cellwork.Materialize(context.Background(), repo, base)
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	defer release()
	for path, body := range write {
		full := filepath.Join(wt.Dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	diff, err := wt.Diff(context.Background())
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if strings.TrimSpace(diff) == "" {
		t.Fatal("the fixture change measured as nothing, so nothing below would be reviewing a change")
	}
	return cellkernel.Matter{Data: diff}
}

func contractFor(t *testing.T, repo, base string) cellkernel.Contract {
	t.Helper()
	iss, err := os.ReadFile(fixture(t, "issue", "valid-issue.json"))
	if err != nil {
		t.Fatal(err)
	}
	subject, err := json.Marshal(cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: base})
	if err != nil {
		t.Fatal(err)
	}
	return cellkernel.Contract{ID: "cds-code", Goal: "make the change", Issue: iss, Subject: subject}
}

// fixture names a file in the shared CDS corpus. The issue this seat reviews
// against is a corpus document, so the catalogue under test is derived from an
// issue BOTH authorities admit — not from one this test invented.
func fixture(t *testing.T, parts ...string) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate this test file")
	}
	root := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
	return filepath.Join(append([]string{root, "schemas", "cds", "fixtures"}, parts...)...)
}

// The whole realization, against a real repository: the seat admits the issue,
// reconstructs the candidate from (subject, matter), runs the closed checker
// against that same directory, decides the two check units itself, and
// reconciles what the judge returned.
func TestReviewReconstructsChecksAndAssesses(t *testing.T) {
	repo, head := celltest.Repo(t)
	matter := matterFor(t, repo, head, map[string]string{"README.md": "changed\n", "NOTES.md": "new\n"})

	j := &recordingJudge{answer: func(c Catalogue) Assessment { return mustJudge(t, c) }}
	seat := AssessBeta{judge: j, method: adversarial(t)}
	out, err := seat.Review(context.Background(), cellkernel.BetaInput{
		Contract: contractFor(t, repo, head),
		Matter:   matter,
	})
	if err != nil {
		t.Fatalf("review: %v", err)
	}

	// Coverage is the corpus issue's two criteria plus the two check units.
	if len(out.Review.Assessment) != 4 {
		t.Fatalf("assessment covers %d units: %+v", len(out.Review.Assessment), out.Review.Assessment)
	}
	byUnit := map[string]cellkernel.UnitResult{}
	for _, u := range out.Review.Assessment {
		byUnit[u.Unit] = u
	}
	if got := byUnit[UnitMatterNonEmpty]; got.Disposition != cellkernel.DispositionPass {
		t.Fatalf("a real diff must pass the matter unit: %+v", got)
	}
	// The candidate is a one-file repository, so this project's recipe cannot
	// build it — a `finding`, decided by the runtime and not by the judge, which
	// answered `unverified` for everything it was allowed to answer.
	if got := byUnit[UnitProjectVerify]; got.Disposition != cellkernel.DispositionFinding ||
		!strings.Contains(got.Reason, cellcheck.RecipeID) {
		t.Fatalf("the checker did not decide its unit: %+v", got)
	}
	if out.Review.Pass {
		t.Fatal("an assessment carrying a finding must not derive a passing review")
	}
	if len(out.Artifacts) != 0 {
		t.Fatalf("beta mints no artifacts; evidence is gamma/V's channel: %+v", out.Artifacts)
	}

	// The prompt carried the runtime-derived view and said whose it is, carried
	// the checker's observation, and carried the catalogue it must dispose of.
	for _, want := range []string{
		"RECONSTRUCTED VIEW", "NOTES.md", "It is not the\nproducer's account",
		"CHECKER OBSERVATION", cellcheck.RecipeID,
		"ASSESSMENT CATALOGUE", UnitProjectVerify, "ALREADY MEASURED BY THE RUNTIME",
		"unverified", "adversarial projection",
	} {
		if !strings.Contains(j.prompt, want) {
			t.Errorf("the prompt does not carry %q", want)
		}
	}
	// ...and nothing about alpha. The seat has no channel to it, so a prompt
	// mentioning a worktree path would mean one had leaked through the matter.
	if strings.Contains(j.prompt, "cnos-cell-worktree-") {
		t.Error("the prompt names a worktree directory: the seat was shown a producer path")
	}
}

// mustJudge answers `unverified` for every unit the seat is allowed to answer
// and repeats every measured one, which is what an honest judge with no
// evidence would do.
func mustJudge(t *testing.T, c Catalogue) Assessment {
	t.Helper()
	a, err := refusingJudge{}.Judge(context.Background(), "", c)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

// A matter with no change in it reconstructs nothing, and that is a verdict
// rather than a malfunction: the matter unit takes the finding, the checker
// unit is unverified because the recipe never ran, and no worktree is cut.
func TestAMatterWithNoChangeReconstructsNothing(t *testing.T) {
	repo, head := celltest.Repo(t)
	j := &recordingJudge{answer: func(c Catalogue) Assessment { return mustJudge(t, c) }}
	seat := AssessBeta{judge: j, method: adversarial(t)}
	out, err := seat.Review(context.Background(), cellkernel.BetaInput{
		Contract: contractFor(t, repo, head),
		Matter:   cellkernel.Matter{Data: "no change was made to " + repo + " at " + head},
	})
	if err != nil {
		t.Fatalf("an unreviewable matter is a verdict, not an error: %v", err)
	}
	byUnit := map[string]cellkernel.UnitResult{}
	for _, u := range out.Review.Assessment {
		byUnit[u.Unit] = u
	}
	if got := byUnit[UnitMatterNonEmpty]; got.Disposition != cellkernel.DispositionFinding {
		t.Fatalf("matter unit = %+v, want a finding", got)
	}
	if got := byUnit[UnitProjectVerify]; got.Disposition != cellkernel.DispositionUnverified {
		t.Fatalf("checker unit = %+v, want unverified: the recipe never ran", got)
	}
	if !strings.Contains(j.prompt, "THERE IS NO VIEW") {
		t.Error("the prompt does not say that no candidate was reconstructed")
	}
}

// A judge that does not review is a malfunction, not a failing review. The seat
// must not manufacture `pass:false` out of an unusable answer — that would put
// a judgement nobody made into the record.
func TestANonReviewingJudgeIsAMalfunction(t *testing.T) {
	repo, head := celltest.Repo(t)
	matter := matterFor(t, repo, head, map[string]string{"README.md": "changed\n"})
	in := cellkernel.BetaInput{Contract: contractFor(t, repo, head), Matter: matter}

	partial := &recordingJudge{answer: func(c Catalogue) Assessment {
		full := mustJudge(t, c)
		return Assessment{Units: full.Units[1:]} // one obligation quietly dropped
	}}
	if _, err := (AssessBeta{judge: partial, method: adversarial(t)}).Review(context.Background(), in); err == nil ||
		!strings.Contains(err.Error(), "did not review") {
		t.Fatalf("a non-covering answer must be a malfunction, got %v", err)
	}

	broken := &recordingJudge{err: fmt.Errorf("the provider fell over")}
	if _, err := (AssessBeta{judge: broken, method: adversarial(t)}).Review(context.Background(), in); err == nil ||
		!strings.Contains(err.Error(), "fell over") {
		t.Fatalf("a failing judge must surface as a malfunction, got %v", err)
	}
}

// A contract this seat cannot read is refused before anything is reconstructed
// or rented: an unadmissible issue means there are no criteria to decide, and a
// subject that is not pinned means two stations could have measured different
// trees.
func TestReviewRefusesAnUnreadableContract(t *testing.T) {
	repo, head := celltest.Repo(t)
	seat := AssessBeta{judge: &recordingJudge{answer: func(c Catalogue) Assessment { return mustJudge(t, c) }}, method: adversarial(t)}
	good := contractFor(t, repo, head)

	noIssue := good
	noIssue.Issue = nil
	if _, err := seat.Review(context.Background(), cellkernel.BetaInput{Contract: noIssue}); err == nil {
		t.Error("a contract with no issue must be refused")
	}
	unpinned := good
	unpinned.Subject = json.RawMessage(`{"kind":"git.snapshot/0.1","repo":"` + repo + `","base_sha":"HEAD"}`)
	if _, err := seat.Review(context.Background(), cellkernel.BetaInput{Contract: unpinned}); err == nil ||
		!strings.Contains(err.Error(), "not pinned") {
		t.Errorf("an unpinned subject must be refused, got %v", err)
	}
	noJudge := AssessBeta{method: adversarial(t)}
	if _, err := noJudge.Review(context.Background(), cellkernel.BetaInput{Contract: good}); err == nil {
		t.Error("a seat with no judge must fail closed")
	}
}
