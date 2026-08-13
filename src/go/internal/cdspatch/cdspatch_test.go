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

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
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
func contractFor(t *testing.T, repo, base string) cellkernel.Contract {
	t.Helper()
	return cellkernel.Contract{
		ID:      "cds-code",
		Goal:    "add a NOTES file",
		Subject: pinnedSubject(t, repo, base),
		RequiredEvidence: []cellkernel.RequiredRef{
			{ID: DiffArtifactID, Kind: DiffArtifactKind, Producer: cellkernel.RoleAlpha},
		},
	}
}

func declJSON(provider string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{
		"fill": "cds.patch",
		"cognition": {"provider": %q, "model": ""},
		"skills": ["cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"]
	}`, provider))
}

var testSkills = []string{"cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"}

func construct(t *testing.T, provider string) cellfill.ConstructedAlpha {
	t.Helper()
	f := Factory(skillTree(t, testSkills...))
	a, err := f(context.Background(), declJSON(provider))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	return a
}

// The constructor resolves and LOADS skill bodies: the resolved declaration
// records ordered refs + content digests, and the prompt carries the bodies.
func TestConstructionLoadsSkillsAndCanonicalizes(t *testing.T) {
	repo, head := testRepo(t)
	a := construct(t, "fake")
	var rd ResolvedDecl
	if err := json.Unmarshal(a.Decl, &rd); err != nil {
		t.Fatalf("resolved decl: %v", err)
	}
	if len(rd.Skills) != 4 || rd.Skills[0].Ref != "cnos.eng:eng/code" || rd.Skills[2].Ref != "cnos.eng:eng/go" {
		t.Fatalf("resolved skills wrong/unordered: %+v", rd.Skills)
	}
	for _, s := range rd.Skills {
		if len(s.SHA256) != 64 {
			t.Fatalf("skill %q has no content digest", s.Ref)
		}
	}
	seat := a.Seat.(PatchAlpha)
	prompt := RenderPrompt(contractFor(t, repo, head), seat.skills)
	if !strings.Contains(prompt, "# body of cnos.eng:eng/go") {
		t.Fatal("skill BODY was not injected into the prompt — naming is not loading")
	}
	if a.Mode != cellkernel.ModeMechanical {
		t.Fatalf("fake provider mode = %q, want mechanical", a.Mode)
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
	f := Factory(skillTree(t, testSkills...))
	const skills = `"skills":["cnos.eng:eng/go"]`
	const cog = `"cognition":{"provider":"fake","model":""}`
	for name, decl := range map[string]string{
		"workspace block": `{"fill":"cds.patch",` + cog + `,` + skills +
			`,"workspace":{"kind":"git-worktree","repo":".","base_sha":"HEAD"}}`,
		"workspace, mixed case": `{"fill":"cds.patch",` + cog + `,` + skills +
			`,"Workspace":{"kind":"git-worktree","repo":".","base_sha":"HEAD"}}`,
		"bare repo key":     `{"fill":"cds.patch",` + cog + `,` + skills + `,"repo":"."}`,
		"bare base_sha key": `{"fill":"cds.patch",` + cog + `,` + skills + `,"base_sha":"HEAD"}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := f(context.Background(), json.RawMessage(decl)); err == nil {
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
	reg.Alpha[Fill] = cellfill.AlphaFill{Construct: Factory(skillTree(t, testSkills...)), NeedsSubject: true}
	first, err := reg.ConstructAlpha(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	second, err := reg.ConstructAlpha(context.Background(), declJSON("fake"))
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
	f := Factory(skillTree(t, testSkills...))
	bad := map[string]string{
		"unknown key":       `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"skills":["cnos.eng:eng/go"],"Extra":1}`,
		"unknown provider":  `{"fill":"cds.patch","cognition":{"provider":"clyde","model":"m"},"skills":["cnos.eng:eng/go"]}`,
		"modelless claude":  `{"fill":"cds.patch","cognition":{"provider":"claude-cli","model":""},"skills":["cnos.eng:eng/go"]}`,
		"no skills":         `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"skills":[]}`,
		"uninstalled skill": `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"skills":["cnos.eng:eng/nope"]}`,
	}
	for name, decl := range bad {
		t.Run(name, func(t *testing.T) {
			if _, err := f(context.Background(), json.RawMessage(decl)); err == nil {
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
	b, err := betas.ConstructBeta(context.Background(), json.RawMessage(`{"fill":"cdd.mechanical-unmet"}`))
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
		"cognition": {"provider": "fake"},
		"skills": ["cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"]
	}`)
	a, err := Factory(skillTree(t, testSkills...))(context.Background(), decl)
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
