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
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
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
		Construct: func(ctx context.Context, decl json.RawMessage, m cellmethod.View) (cellfill.ConstructedAlpha, error) {
			*constructions++
			return inner.Construct(ctx, decl, m)
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

// seatConstructions counts what each side's recording constructor was asked to
// build. "Refused before construction" is a claim about what did NOT run, so a
// recorded call is the only thing that can falsify it.
type seatConstructions struct{ alpha, beta int }

// betaSubjectRequiringRegistry re-registers `cdd.stub` on the ASSESSING side as
// a fill that declares it cannot act without a subject, leaving the producing
// side exactly as it ships — declaring nothing. That pairing is the divergence
// this test exists for: a beta that needs the subject beside an alpha that does
// not. Both constructors record, because the alpha's silence must not be what
// lets the run through.
func betaSubjectRequiringRegistry(built *seatConstructions) cellfill.Registry {
	reg := cellfill.CddFills()
	innerAlpha := reg.Alpha[cellfill.FillStubAlpha]
	reg.Alpha[cellfill.FillStubAlpha] = cellfill.AlphaFill{
		Construct: func(ctx context.Context, decl json.RawMessage, m cellmethod.View) (cellfill.ConstructedAlpha, error) {
			built.alpha++
			return innerAlpha.Construct(ctx, decl, m)
		},
	}
	innerBeta := reg.Beta[cellfill.FillStubBeta]
	reg.Beta[cellfill.FillStubBeta] = cellfill.BetaFill{
		Construct: func(ctx context.Context, decl json.RawMessage, m cellmethod.View) (cellfill.ConstructedBeta, error) {
			built.beta++
			return innerBeta.Construct(ctx, decl, m)
		},
		NeedsSubject: true,
	}
	return reg
}

// refuseBetaConstruction decorates the recording beta so that running it AT ALL
// fails the test. The count alone would prove the same thing, but only after
// the build returned; this states the prohibition at the exact point it would
// be broken, which is where a reader looks for it.
func refuseBetaConstruction(t *testing.T, reg cellfill.Registry) cellfill.Registry {
	t.Helper()
	recording := reg.Beta[cellfill.FillStubBeta]
	reg.Beta[cellfill.FillStubBeta] = cellfill.BetaFill{
		Construct: func(ctx context.Context, decl json.RawMessage, m cellmethod.View) (cellfill.ConstructedBeta, error) {
			t.Errorf("the beta constructor ran for a binding its own fill declared it cannot act on")
			return recording.Construct(ctx, decl, m)
		},
		NeedsSubject: recording.NeedsSubject,
	}
	return reg
}

// The assessing side's declared requirement decides the run BEFORE its
// constructor runs — the same property the alpha side already holds. This is
// the `cds.assess` case: a beta that reconstructs the candidate from the pinned
// subject, paired with an alpha that needs none. Undeclared, such a cell built
// both seats and failed inside Review as an episode malfunction; declared, it
// is a spec refusal with nothing constructed.
func TestADeclaredBetaSubjectRequirementIsRefusedBeforeConstruction(t *testing.T) {
	s, err := Parse([]byte(stubCell))
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}

	var built seatConstructions
	_, _, err = r.Build(context.Background(),
		refuseBetaConstruction(t, betaSubjectRequiringRegistry(&built)), Binding{})
	if err == nil {
		t.Fatal("a beta that cannot act without a subject must not be built without one")
	}
	if built != (seatConstructions{}) {
		t.Fatalf("constructors ran before the refusal: %+v", built)
	}
	// The message names the SIDE, the fill and the input that is missing: a
	// refusal that said only "requires contract.subject" would leave an operator
	// looking at two seats.
	for _, want := range []string{"beta", cellfill.FillStubBeta, "contract.subject", "--input"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("the refusal does not name %q: %v", want, err)
		}
	}

	// ...and with a subject the SAME declaration constructs both seats. Without
	// this the assertions above would hold for a registry that refused
	// everything. Only the refusing decoration is dropped: the fill still
	// declares the requirement, and the requirement is now satisfied.
	var admitted seatConstructions
	if _, _, err := r.Build(context.Background(), betaSubjectRequiringRegistry(&admitted),
		Binding{Subject: json.RawMessage(`{"opaque":"subject"}`)}); err != nil {
		t.Fatalf("a supplied subject must satisfy the requirement: %v", err)
	}
	if admitted != (seatConstructions{alpha: 1, beta: 1}) {
		t.Fatalf("constructions = %+v, want exactly one build of each seat", admitted)
	}
}

// --- one methodology, loaded once, before any seat -------------------------

const methodCell = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id": "m1", "goal": "hold the seat to something"},
  "protocol_id": "p",
  "params": {"language": {"required": true}},
  "methodology": {"kind": "skills.methodology.v0",
                  "skills": ["cnos.eng:eng/code", "$language"]},
  "alpha": {"fill": "cdd.stub"},
  "beta": {"fill": "cdd.stub"}
}`

// installedTree writes a package tree with the given skills installed.
func installedTree(t *testing.T, refs ...string) cellskill.Tree {
	t.Helper()
	root := t.TempDir()
	for _, ref := range refs {
		pkg, path, _ := strings.Cut(ref, ":")
		dir := filepath.Join(root, pkg, "skills", filepath.FromSlash(path))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# body of "+ref+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return cellskill.Tree{Root: root}
}

// viewCapturingRegistry records the projection the alpha constructor receives.
func viewCapturingRegistry(skills cellskill.Resolver, got *cellmethod.View, constructions *int) cellfill.Registry {
	reg := cellfill.CddFills()
	reg.Skills = skills
	inner := reg.Alpha[cellfill.FillStubAlpha]
	reg.Alpha[cellfill.FillStubAlpha] = cellfill.AlphaFill{
		Construct: func(ctx context.Context, decl json.RawMessage, m cellmethod.View) (cellfill.ConstructedAlpha, error) {
			*constructions++
			*got = m
			return inner.Construct(ctx, decl, m)
		},
	}
	return reg
}

// AC1, this package's half: the cell's ONE bundle is loaded here, once, and
// what reaches the producing seat is a projection of it — carrying the digest
// of the ordered (ref, body-digest) list and the bodies of the skills the CELL
// declared, with the hole filled.
func TestTheCellsOneMethodologyReachesTheProducingSeat(t *testing.T) {
	tree := installedTree(t, "cnos.eng:eng/code", "cnos.eng:eng/go")
	s, err := Parse([]byte(methodCell))
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Resolve(map[string]string{"language": "cnos.eng:eng/go"})
	if err != nil {
		t.Fatal(err)
	}
	var got cellmethod.View
	var constructions int
	if _, _, err := r.Build(context.Background(),
		viewCapturingRegistry(tree, &got, &constructions), Binding{}); err != nil {
		t.Fatalf("build: %v", err)
	}
	if got.Role != cellmethod.RoleConstructive {
		t.Fatalf("the seat received the %q projection, want constructive", got.Role)
	}
	bundle, _, err := cellmethod.Load(tree, r.Methodology)
	if err != nil {
		t.Fatal(err)
	}
	if got.SHA256 != bundle.SHA256 {
		t.Fatalf("the seat received digest %q, want the cell's bundle %q", got.SHA256, bundle.SHA256)
	}
	for _, ref := range []string{"cnos.eng:eng/code", "cnos.eng:eng/go"} {
		if !strings.Contains(got.Text, "# body of "+ref) {
			t.Fatalf("the projection does not carry the body of %q", ref)
		}
	}
}

// A methodology that cannot load fails THE CELL, before any seat exists. It is
// not a per-fill question: a bundle naming an uninstalled skill is a broken
// cell whatever its seats happen to want, and discovering it inside a
// constructor would report it as that seat's problem.
func TestAnUnloadableMethodologyFailsBeforeAnySeatIsConstructed(t *testing.T) {
	tree := installedTree(t, "cnos.eng:eng/code") // eng/go is NOT installed
	s, err := Parse([]byte(methodCell))
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Resolve(map[string]string{"language": "cnos.eng:eng/go"})
	if err != nil {
		t.Fatal(err)
	}
	var got cellmethod.View
	var constructions int
	_, _, err = r.Build(context.Background(),
		viewCapturingRegistry(tree, &got, &constructions), Binding{})
	if err == nil {
		t.Fatal("a methodology naming an uninstalled skill must fail the cell")
	}
	if !strings.Contains(err.Error(), "cnos.eng:eng/go") {
		t.Fatalf("the failure does not name the missing skill: %v", err)
	}
	if constructions != 0 {
		t.Fatalf("the alpha constructor ran %d time(s) before the refusal", constructions)
	}
	// ...and the SAME cell over a tree that has the skill builds, or the
	// assertions above would hold for a registry that refused everything.
	if _, _, err := r.Build(context.Background(),
		viewCapturingRegistry(installedTree(t, "cnos.eng:eng/code", "cnos.eng:eng/go"), &got, &constructions),
		Binding{}); err != nil {
		t.Fatalf("an installed methodology must build: %v", err)
	}
}

// A cell that declares NO methodology still builds. Every cell did before the
// bundle existed, and the generic loader must not have acquired a rule that
// belongs to the fills that need one.
func TestACellWithNoMethodologyStillBuilds(t *testing.T) {
	var got cellmethod.View
	var constructions int
	s, err := Parse([]byte(stubCell))
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	// The registry carries no skill resolver either: nothing to load, nothing
	// to resolve it against.
	if _, _, err := r.Build(context.Background(),
		viewCapturingRegistry(nil, &got, &constructions), Binding{}); err != nil {
		t.Fatalf("a cell declaring no methodology must still build: %v", err)
	}
	if !got.Empty() {
		t.Fatalf("a cell declaring no methodology projected %+v", got)
	}
}
