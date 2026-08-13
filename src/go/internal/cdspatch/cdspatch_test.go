package cdspatch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// testRepo builds a one-commit git repository and returns its path and HEAD.
func testRepo(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) string {
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
	run("init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("base\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("add", "-A")
	run("commit", "-qm", "base")
	return dir, run("rev-parse", "HEAD")
}

// skillTree writes a minimal installed package tree with the given skills.
func skillTree(t *testing.T, refs ...string) cellskill.Tree {
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

// pinnedSubject is the contract slot a station reads its repository and base
// from — the ONLY place either is stated. Written with cellwork.Subject rather
// than a literal so this file cannot drift from the language cellwork owns.
func pinnedSubject(t *testing.T, repo, base string) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: base})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// contractFor is the frozen contract an episode hands the seat. The subject is
// part of it because the subject is what says where the work happens.
// fixtureIssue is the corpus's own admissible issue. Read from the shared
// corpus rather than hand-built here: an issue this package invented could
// drift from what the door admits, and then the seat would be tested against a
// document no run can produce.
func fixtureIssue(t *testing.T) json.RawMessage {
	t.Helper()
	raw, err := os.ReadFile("../../../../schemas/cds/fixtures/issue/valid-issue.json")
	if err != nil {
		t.Fatalf("read the issue corpus fixture: %v", err)
	}
	if _, err := cdsissue.Admit(raw); err != nil {
		t.Fatalf("the corpus fixture is not admissible: %v", err)
	}
	return json.RawMessage(raw)
}

func contractFor(t *testing.T, repo, base string) cellkernel.Contract {
	t.Helper()
	return cellkernel.Contract{
		ID:      "cds-code",
		Goal:    "add a NOTES file",
		Issue:   fixtureIssue(t),
		Subject: pinnedSubject(t, repo, base),
		RequiredEvidence: []cellkernel.RequiredRef{
			{ID: DiffArtifactID, Kind: DiffArtifactKind, Producer: cellkernel.RoleAlpha},
		},
	}
}

// The declaration is fill + cognition, and nothing else. There is no skills
// key here and there cannot be one: what holds this seat is the cell's
// methodology, projected and handed in.
func declJSON(provider string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{
		"fill": "cds.patch",
		"cognition": {"provider": %q, "model": ""}
	}`, provider))
}

var testSkills = []string{"cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"}

// methodology loads the cell's ONE bundle from an installed tree and projects
// it constructively — the same two calls cellspec.Build makes before it
// constructs a seat. Built through cellmethod rather than hand-written, so this
// file cannot drift from the projection the runtime actually hands over.
func methodology(t *testing.T, tree cellskill.Tree, refs ...string) cellmethod.View {
	t.Helper()
	decl, err := json.Marshal(map[string]any{"kind": cellmethod.Kind, "skills": refs})
	if err != nil {
		t.Fatal(err)
	}
	b, bodies, err := cellmethod.Load(tree, decl)
	if err != nil {
		t.Fatalf("load methodology: %v", err)
	}
	return cellmethod.Constructive(b, bodies)
}

func construct(t *testing.T, provider string) cellfill.ConstructedAlpha {
	t.Helper()
	a, err := Factory()(context.Background(), declJSON(provider),
		methodology(t, skillTree(t, testSkills...), testSkills...))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	return a
}

// AC1, this package's half: ONE bundle reaches the producing seat. The
// declaration records the digest of the projection it was handed and the role
// of that projection, and the prompt carries the bundle's skill BODIES — naming
// a skill is not loading it.
func TestTheCellsMethodologyReachesTheSeatAndIsRecorded(t *testing.T) {
	repo, head := testRepo(t)
	tree := skillTree(t, testSkills...)
	view := methodology(t, tree, testSkills...)
	a, err := Factory()(context.Background(), declJSON("fake"), view)
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	var rd ResolvedDecl
	if err := json.Unmarshal(a.Decl, &rd); err != nil {
		t.Fatalf("resolved decl: %v", err)
	}
	if rd.Methodology.SHA256 != view.SHA256 || len(rd.Methodology.SHA256) != 64 {
		t.Fatalf("recorded methodology digest %q != the projected bundle's %q",
			rd.Methodology.SHA256, view.SHA256)
	}
	if rd.Methodology.Role != string(cellmethod.RoleConstructive) {
		t.Fatalf("recorded role = %q, want constructive", rd.Methodology.Role)
	}

	seat := a.Seat.(PatchAlpha)
	c := contractFor(t, repo, head)
	issue, err := cdsissue.Admit(c.Issue)
	if err != nil {
		t.Fatalf("the fixture contract carries no admissible issue: %v", err)
	}
	prompt := RenderPrompt(c, issue, seat.method)
	// The producing seat is judged against the acceptance criteria, so it has
	// to be shown them. A seat given only the one-line goal writes against a
	// summary of its contract and is then reviewed against the contract — which
	// is the defect the typed issue exists to end, and it does not end at the
	// door.
	for _, crit := range issue.Acceptance {
		if !strings.Contains(prompt, crit.ID) || !strings.Contains(prompt, crit.Statement) {
			t.Fatalf("the producing prompt omits acceptance criterion %q", crit.ID)
		}
		if !strings.Contains(prompt, crit.Verification) {
			t.Fatalf("the producing prompt omits the verification route for %q", crit.ID)
		}
	}
	for _, ref := range testSkills {
		if !strings.Contains(prompt, "# body of "+ref) {
			t.Fatalf("body of %q was not injected into the prompt", ref)
		}
	}
	if !strings.Contains(prompt, view.SHA256) {
		t.Fatal("the prompt does not state which methodology holds the seat")
	}
	if a.Mode != cellkernel.ModeMechanical {
		t.Fatalf("fake provider mode = %q, want mechanical", a.Mode)
	}

	// ...and the recorded digest MOVES with the bodies. Without this the
	// equality above would hold for a digest of the ref list alone, which would
	// record which files were named rather than what they said.
	changed := skillTree(t, testSkills...)
	pkg, path, _ := strings.Cut(testSkills[2], ":")
	body := filepath.Join(changed.Root, pkg, "skills", filepath.FromSlash(path), "SKILL.md")
	if err := os.WriteFile(body, []byte("# body of "+testSkills[2]+"\n!"), 0o600); err != nil {
		t.Fatal(err)
	}
	other, err := Factory()(context.Background(), declJSON("fake"),
		methodology(t, changed, testSkills...))
	if err != nil {
		t.Fatal(err)
	}
	var rd2 ResolvedDecl
	if err := json.Unmarshal(other.Decl, &rd2); err != nil {
		t.Fatal(err)
	}
	if rd2.Methodology.SHA256 == rd.Methodology.SHA256 {
		t.Fatal("a one-byte change to a skill body did not change the recorded methodology digest")
	}
}

// AC1, the deletion: `cds.patch` has no `skills` key. It refuses to be told one
// — at any case spelling — and it puts none in the record. Both halves matter:
// the first says the fill cannot be given a second methodology, the second says
// it does not mint one.
func TestTheDeclarationCannotCarryItsOwnSkills(t *testing.T) {
	view := methodology(t, skillTree(t, testSkills...), testSkills...)
	const cog = `"cognition":{"provider":"fake","model":""}`
	for name, decl := range map[string]string{
		"skills list":        `{"fill":"cds.patch",` + cog + `,"skills":["cnos.eng:eng/go"]}`,
		"skills, mixed case": `{"fill":"cds.patch",` + cog + `,"Skills":["cnos.eng:eng/go"]}`,
		"empty skills list":  `{"fill":"cds.patch",` + cog + `,"skills":[]}`,
		"a methodology of its own": `{"fill":"cds.patch",` + cog +
			`,"methodology":{"kind":"skills.methodology.v0","skills":["cnos.eng:eng/go"]}}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Factory()(context.Background(), json.RawMessage(decl), view)
			if err == nil {
				t.Fatal("a declaration carrying its own skills must not construct")
			}
			if !strings.Contains(err.Error(), "unknown key") {
				t.Fatalf("rejected for the wrong reason: %v", err)
			}
		})
	}

	// Checked on the raw canonical bytes, not on ResolvedDecl: a struct with no
	// field for `skills` could not report one if the fill emitted it.
	a := construct(t, "fake")
	var recorded map[string]json.RawMessage
	if err := json.Unmarshal(a.Decl, &recorded); err != nil {
		t.Fatal(err)
	}
	if _, ok := recorded["skills"]; ok {
		t.Fatalf("the recorded declaration still carries a skills key: %s", a.Decl)
	}
}

// A patch alpha cannot act with no methodology at all, and says so itself: the
// generic loader must not acquire the rule, because a cell with no methodology
// is legitimate for other fills.
func TestNoMethodologyIsARefusalAtTheFill(t *testing.T) {
	_, err := Factory()(context.Background(), declJSON("fake"), cellmethod.View{})
	if err == nil {
		t.Fatal("a patch alpha with no methodology must not construct")
	}
	if !strings.Contains(err.Error(), "declares none") {
		t.Fatalf("rejected for the wrong reason: %v", err)
	}
	// A producing seat must not be handed the falsification projection.
	tree := skillTree(t, testSkills...)
	b, bodies, err := cellmethod.Load(tree, []byte(`{"kind":"`+cellmethod.Kind+
		`","skills":["cnos.eng:eng/go"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Factory()(context.Background(), declJSON("fake"),
		cellmethod.Adversarial(b, bodies)); err == nil ||
		!strings.Contains(err.Error(), "constructive projection") {
		t.Fatalf("the adversarial projection must not construct a producing seat: %v", err)
	}
}

// F1. The declaration cannot name a repository — not as a workspace block, not
// as a stray key, not at any case spelling. Both halves are asserted, because
// they fail differently: a rejected key proves the fill refuses to be told, and
// the absent field in the RECORDED declaration proves it does not put one there
// itself. Together they are "the record has one repository declaration", which
// is the property; while there were two, a closure naming a repository the
// episode never touched still self-verified.
func TestTheDeclarationCannotNameARepository(t *testing.T) {
	f := Factory()
	view := methodology(t, skillTree(t, testSkills...), testSkills...)
	const cog = `"cognition":{"provider":"fake","model":""}`
	for name, decl := range map[string]string{
		"workspace block": `{"fill":"cds.patch",` + cog +
			`,"workspace":{"kind":"git-worktree","repo":".","base_sha":"HEAD"}}`,
		"workspace, mixed case": `{"fill":"cds.patch",` + cog +
			`,"Workspace":{"kind":"git-worktree","repo":".","base_sha":"HEAD"}}`,
		"bare repo key":     `{"fill":"cds.patch",` + cog + `,"repo":"."}`,
		"bare base_sha key": `{"fill":"cds.patch",` + cog + `,"base_sha":"HEAD"}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := f(context.Background(), json.RawMessage(decl), view); err == nil {
				t.Fatal("a declaration naming a repository must not construct")
			} else if !strings.Contains(err.Error(), "unknown key") {
				t.Fatalf("rejected for the wrong reason: %v", err)
			}
		})
	}

	// ...and the canonical declaration the record carries states nothing about
	// a repository either. Checked on the raw canonical bytes, not on
	// ResolvedDecl: a struct with no field for it could not report one.
	a := construct(t, "fake")
	var recorded map[string]json.RawMessage
	if err := json.Unmarshal(a.Decl, &recorded); err != nil {
		t.Fatal(err)
	}
	for key := range recorded {
		switch key {
		case "workspace", "repo", "base_sha":
			t.Fatalf("the recorded declaration names a repository: %s", a.Decl)
		}
	}
}

// F1's other half: what the seat MEASURES is what the contract's subject
// pinned. There is one repository declaration in the record, so the base the
// episode actually stood on and the base the record says it stood on are one
// resolution, not two that happened to agree.
func TestTheMeasuredBaseIsTheContractSubjectsBase(t *testing.T) {
	repo, head := testRepo(t)
	seat := construct(t, "fake").Seat.(PatchAlpha)
	contract := contractFor(t, repo, head)

	out, err := seat.Produce(context.Background(), cellkernel.AlphaInput{Contract: contract})
	if err != nil {
		t.Fatalf("produce: %v", err)
	}
	var measured string
	for _, c := range out.Artifacts {
		if c.ID == BaseArtifactID {
			measured = c.Text
		}
	}
	subject, err := cellwork.AdmitSubject(contract.Subject)
	if err != nil {
		t.Fatal(err)
	}
	if measured != subject.BaseSHA {
		t.Fatalf("measured base %q != contract.subject.base_sha %q", measured, subject.BaseSHA)
	}
	// The assertion above would hold vacuously if the seat could measure
	// anything at all: pointed at a subject naming a commit that does not
	// exist, it must fail rather than fall back to some other revision.
	absent := contractFor(t, repo, strings.Repeat("0", 40))
	if _, err := seat.Produce(context.Background(), cellkernel.AlphaInput{Contract: absent}); err == nil {
		t.Fatal("a subject naming an unresolvable commit must fail the seat")
	}
}

// A cds.patch seat with no subject cannot invent one. Before the deletion it
// carried its own repository and would happily run against it while the
// contract said nothing — which is how a record could describe work on a
// repository the contract never named.
func TestNoSubjectIsARefusalNotADefault(t *testing.T) {
	seat := construct(t, "fake").Seat.(PatchAlpha)
	_, err := seat.Produce(context.Background(), cellkernel.AlphaInput{
		Contract: cellkernel.Contract{ID: "c", Goal: "g"},
	})
	if err == nil {
		t.Fatal("a contract with no subject must not produce")
	}
	if !strings.Contains(err.Error(), "carries no subject") {
		t.Fatalf("rejected for the wrong reason: %v", err)
	}
}

// An UNPINNED subject is refused at the seat too. The subject is pinned once,
// before either station exists; a station that re-resolved `HEAD` is exactly
// the two-resolutions failure in a different place.
func TestAnUnpinnedSubjectIsRefusedAtTheSeat(t *testing.T) {
	repo, _ := testRepo(t)
	seat := construct(t, "fake").Seat.(PatchAlpha)
	_, err := seat.Produce(context.Background(), cellkernel.AlphaInput{
		Contract: contractFor(t, repo, "HEAD"),
	})
	if err == nil {
		t.Fatal("an unpinned base must not reach a station")
	}
	if !strings.Contains(err.Error(), "is not pinned") {
		t.Fatalf("rejected for the wrong reason: %v", err)
	}
}

// The registry canonicalizes what a fill returns, so the recorded bytes — and
// therefore the digest — do not depend on how the fill happened to serialize,
// and constructing twice yields byte-identical declarations.
func TestConstructionIsCanonicalAndStable(t *testing.T) {
	reg := cellfill.CddFills()
	reg.Alpha[Fill] = cellfill.AlphaFill{Construct: Factory(), NeedsSubject: true}
	view := methodology(t, skillTree(t, testSkills...), testSkills...)
	first, err := reg.ConstructAlpha(context.Background(), declJSON("fake"), view)
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	second, err := reg.ConstructAlpha(context.Background(), declJSON("fake"), view)
	if err != nil {
		t.Fatalf("construct again: %v", err)
	}
	if string(first.Decl) != string(second.Decl) {
		t.Fatalf("construction is not stable:\n%s\n%s", first.Decl, second.Decl)
	}
	// Canonical form: keys sorted, no insignificant whitespace.
	if !strings.HasPrefix(string(first.Decl), `{"cognition":`) {
		t.Fatalf("declaration is not canonical: %s", first.Decl)
	}
}

func TestConstructionFailsClosed(t *testing.T) {
	f := Factory()
	view := methodology(t, skillTree(t, testSkills...), testSkills...)
	bad := map[string]string{
		"unknown key":      `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"Extra":1}`,
		"unknown provider": `{"fill":"cds.patch","cognition":{"provider":"clyde","model":"m"}}`,
		"modelless claude": `{"fill":"cds.patch","cognition":{"provider":"claude-cli","model":""}}`,
		"no cognition":     `{"fill":"cds.patch"}`,
	}
	for name, decl := range bad {
		t.Run(name, func(t *testing.T) {
			if _, err := f(context.Background(), json.RawMessage(decl), view); err == nil {
				t.Fatalf("%s must fail construction", name)
			}
		})
	}
}

// The load-bearing property: the change is measured from the worktree. A
// coder that changed nothing produces no diff, so a contract requiring one
// closes needs_repair — false completion (#514/#516) is unrepresentable.
type idleCoder struct{}

func (idleCoder) Name() string                               { return "idle" }
func (idleCoder) Work(context.Context, string, string) error { return nil }

func TestIdleCoderCannotFalselyComplete(t *testing.T) {
	repo, head := testRepo(t)
	a := construct(t, "fake").Seat.(PatchAlpha)
	a.coder = idleCoder{}
	out, err := a.Produce(context.Background(), cellkernel.AlphaInput{Contract: contractFor(t, repo, head)})
	if err != nil {
		t.Fatalf("produce: %v", err)
	}
	for _, c := range out.Artifacts {
		if c.ID == DiffArtifactID {
			t.Fatal("a diff artifact was manufactured from no change")
		}
	}
	if !strings.Contains(out.Matter.Data, "no change was made") {
		t.Fatalf("matter must state the truth: %q", out.Matter.Data)
	}
}

// Case-2 honesty end to end: even a REAL measured change closes needs_repair,
// because the mechanical-unmet beta cannot judge the goal and says so. The
// diff is preserved in the closure for the independent review Case 3 adds.
// (The fake coder writes a file unrelated to the NOTES goal — exactly the
// change a passing mechanical beta would have wrongly blessed.)
func TestMeasuredChangeAwaitsIndependentReview(t *testing.T) {
	repo, head := testRepo(t)
	a := construct(t, "fake")
	patchContract := contractFor(t, repo, head)
	betas := cellfill.CddFills()
	b, err := betas.ConstructBeta(context.Background(), json.RawMessage(`{"fill":"cdd.mechanical-unmet"}`), cellmethod.View{})
	if err != nil {
		t.Fatalf("beta: %v", err)
	}
	meta := cellkernel.RunMeta{
		ExecutionMode: cellfill.CombineModes(a.Mode, b.Mode),
		ResolvedSpec: cellkernel.ResolvedSpec{
			Version: "cnos.cellspec.v0", DeclaredProtocol: "cnos.cdd.cds.receipt.v1",
			Alpha: a.Decl, Beta: b.Decl,
		},
	}
	cl, err := cellkernel.RunEpisode(context.Background(),
		cellkernel.Spec{Contract: patchContract, Alpha: a.Seat, Beta: b.Seat}, meta)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != cellkernel.NeedsRepair {
		t.Fatalf("a mechanically unreviewed change closed %q, want needs_repair", cl.Status)
	}
	var diff, base string
	for _, c := range cl.Receipt.Record.Alpha.Artifacts {
		switch c.ID {
		case DiffArtifactID:
			diff = c.Text
		case BaseArtifactID:
			base = c.Text
		}
	}
	if !strings.Contains(diff, "CELL-FAKE-CHANGE.txt") {
		t.Fatal("the measured diff was not preserved in the closure")
	}
	if base != head {
		t.Fatalf("base_sha = %q, want resolved HEAD %q", base, head)
	}
	if err := cellkernel.VerifyClosure(patchContract, meta, cl); err != nil {
		t.Fatalf("closure must verify: %v", err)
	}
}

// The worktree is disposable and the repository is left untouched.
func TestWorktreeIsDisposableAndRepoUntouched(t *testing.T) {
	repo, head := testRepo(t)
	a := construct(t, "fake").Seat.(PatchAlpha)
	if _, err := a.Produce(context.Background(), cellkernel.AlphaInput{Contract: contractFor(t, repo, head)}); err != nil {
		t.Fatalf("produce: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "CELL-FAKE-CHANGE.txt")); !os.IsNotExist(err) {
		t.Fatal("the seat wrote into the caller's repository")
	}
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(out)) != head {
		t.Fatalf("repository HEAD moved: %s (%v)", out, err)
	}
	cmd = exec.Command("git", "worktree", "list")
	cmd.Dir = repo
	list, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(list), "\n") != 1 {
		t.Fatalf("a worktree outlived the episode:\n%s", list)
	}
}

// The seat fails closed on every input it cannot honour. The repository and
// base cases arrive through the CONTRACT now, which is the only place they are
// stated — reaching into the seat's fields to corrupt them is no longer
// possible, and that is the fix rather than a limitation of this test.
func TestPatchAlphaFailsClosed(t *testing.T) {
	repo, head := testRepo(t)
	seat := construct(t, "fake").Seat.(PatchAlpha)

	noCoder := seat
	noCoder.coder = nil
	if _, err := noCoder.Produce(context.Background(),
		cellkernel.AlphaInput{Contract: contractFor(t, repo, head)}); err == nil {
		t.Fatal("nil coder must fail closed")
	}
	if _, err := seat.Produce(context.Background(),
		cellkernel.AlphaInput{Contract: contractFor(t, t.TempDir(), head)}); err == nil {
		t.Fatal("a subject naming a non-repository must fail closed")
	}
}

// A fake ignores the model, so an authored declaration may simply omit the
// key — Go's decoder yields the empty model either way. The authored CUE
// shape admits the same omission; this is the Go half of that parity
// (Pi #57 C1), and the RESOLVED declaration must still record `model: ""`
// so a receipt says what held the seat rather than what the author typed.
func TestFakeMayOmitModelAndStillRecordsIt(t *testing.T) {
	decl := json.RawMessage(`{
		"fill": "cds.patch",
		"cognition": {"provider": "fake"}
	}`)
	a, err := Factory()(context.Background(), decl,
		methodology(t, skillTree(t, testSkills...), testSkills...))
	if err != nil {
		t.Fatalf("a fake omitting its meaningless model must construct: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(a.Decl, &raw); err != nil {
		t.Fatal(err)
	}
	var cog map[string]any
	if err := json.Unmarshal(raw["cognition"], &cog); err != nil {
		t.Fatal(err)
	}
	if m, present := cog["model"]; !present || m != "" {
		t.Fatalf("resolved cognition must record model:\"\", got %v (present=%v)", m, present)
	}
}
